package cognito

import (
	"context"
	"log"
	"strings"

	securityDomain "github.com/TebanMT/smartGou/src/modules/security/domain"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
	"github.com/google/uuid"
)

type CognitoService struct {
	client           *cognitoidentityprovider.Client
	userPoolId       string
	userPoolClientId string
}

func NewCognitoService(userPoolId string, userPoolClientId string) (*CognitoService, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	client := cognitoidentityprovider.NewFromConfig(cfg)

	return &CognitoService{client: client, userPoolId: userPoolId, userPoolClientId: userPoolClientId}, nil
}

func (s *CognitoService) SendOTPToPhone(ctx context.Context, userID uuid.UUID) (*securityDomain.LoginChallengeByPhone, error) {
	// Send a code to the user's phone number
	result, err := s.client.InitiateAuth(ctx, &cognitoidentityprovider.InitiateAuthInput{
		AuthFlow: types.AuthFlowTypeCustomAuth,
		ClientId: aws.String(s.userPoolClientId),
		AuthParameters: map[string]string{
			"USERNAME": userID.String(),
		},
	})
	if err != nil {
		return nil, err
	}

	session := *result.Session
	code := result.ChallengeParameters["answer"]

	return securityDomain.NewLoginChallengeByPhone(code, session, userID), nil

}

func (s *CognitoService) VerifyOTPFromPhone(ctx context.Context, loginChallenge *securityDomain.LoginChallengeByPhone) (*securityDomain.VerifyOTPByPhoneResponse, error) {
	verifyOTPByPhoneResponse := &securityDomain.VerifyOTPByPhoneResponse{
		LoginChallenge: nil,
		TokenEntity:    nil,
	}
	result, err := s.client.RespondToAuthChallenge(ctx, &cognitoidentityprovider.RespondToAuthChallengeInput{
		ChallengeName: types.ChallengeNameTypeCustomChallenge,
		ChallengeResponses: map[string]string{
			"USERNAME": loginChallenge.UserId.String(),
			"ANSWER":   loginChallenge.Code,
		},
		ClientId: aws.String(s.userPoolClientId),
		Session:  aws.String(loginChallenge.Session),
	})
	if err != nil {
		if strings.Contains(err.Error(), "NotAuthorizedException: Invalid session for the user.") ||
			strings.Contains(err.Error(), "CodeMismatchException: Invalid code/session provided") ||
			strings.Contains(err.Error(), "InvalidParameterException") {
			return verifyOTPByPhoneResponse, securityDomain.ErrInvalidSession
		}
		if strings.Contains(err.Error(), "UserNotFoundException") {
			return verifyOTPByPhoneResponse, securityDomain.ErrUserNotFoundException
		}
		if strings.Contains(err.Error(), "NotAuthorizedException: Incorrect username or password.") {
			return verifyOTPByPhoneResponse, securityDomain.ErrMaxAttemptsReached
		}
		return verifyOTPByPhoneResponse, err
	}
	if result.AuthenticationResult == nil {
		loginChallenge.Session = *result.Session
		verifyOTPByPhoneResponse.LoginChallenge = loginChallenge
		return verifyOTPByPhoneResponse, securityDomain.ErrInvalidOTP

	}

	s.client.AdminUpdateUserAttributes(ctx, &cognitoidentityprovider.AdminUpdateUserAttributesInput{
		UserPoolId: aws.String(s.userPoolId),
		Username:   aws.String(loginChallenge.UserId.String()),
		UserAttributes: []types.AttributeType{
			{
				Name:  aws.String("phone_number_verified"),
				Value: aws.String("true"),
			},
		},
	})

	tokenEntity := securityDomain.NewTokenEntity(
		*result.AuthenticationResult.AccessToken,
		*result.AuthenticationResult.RefreshToken,
		*result.AuthenticationResult.IdToken,
		int(result.AuthenticationResult.ExpiresIn),
	)
	verifyOTPByPhoneResponse.TokenEntity = tokenEntity

	return verifyOTPByPhoneResponse, nil
}

func (s *CognitoService) RegisterWithPhoneNumber(ctx context.Context, phoneNumber string, userID uuid.UUID) error {
	// Create a new user in Cognito without password
	input := &cognitoidentityprovider.AdminCreateUserInput{
		Username:   aws.String(userID.String()),
		UserPoolId: aws.String(s.userPoolId),
		DesiredDeliveryMediums: []types.DeliveryMediumType{
			types.DeliveryMediumTypeSms,
		},
		MessageAction: types.MessageActionTypeSuppress,
		UserAttributes: []types.AttributeType{
			{
				Name:  aws.String("phone_number"),
				Value: aws.String(phoneNumber),
			},
		},
	}

	_, err := s.client.AdminCreateUser(context.TODO(), input)
	if err != nil {
		if !strings.Contains(err.Error(), "UsernameExistsException") {
			return err
		}
	}

	return nil
}

func (s *CognitoService) RegisterWithEmail(ctx context.Context, email string, password string, userID uuid.UUID) error {
	input := &cognitoidentityprovider.SignUpInput{
		ClientId: aws.String(s.userPoolClientId),
		Username: aws.String(userID.String()),
		Password: aws.String(password),
		UserAttributes: []types.AttributeType{
			{
				Name:  aws.String("email"),
				Value: aws.String(email),
			},
		},
	}

	_, err := s.client.SignUp(ctx, input)
	if err != nil {
		return err
	}

	return nil
}

func (s *CognitoService) ConfirmOtpByEmail(ctx context.Context, userID uuid.UUID, code string) error {
	_, err := s.client.ConfirmSignUp(ctx, &cognitoidentityprovider.ConfirmSignUpInput{
		Username:         aws.String(userID.String()),
		ClientId:         aws.String(s.userPoolClientId),
		ConfirmationCode: aws.String(code),
	})
	if err != nil {
		if strings.Contains(err.Error(), "UserNotFoundException") {
			return securityDomain.ErrUserNotFoundException
		}
		if strings.Contains(err.Error(), "CodeMismatchException") {
			return securityDomain.ErrInvalidOTP
		}
		if strings.Contains(err.Error(), "NotAuthorizedException") && strings.Contains(err.Error(), "CONFIRMED") {
			return securityDomain.ErrUserAlreadyConfirmed
		}
		if strings.Contains(err.Error(), "UsernameExistsException") {
			return securityDomain.ErrUserAlreadyExists
		}
		return err
	}

	return nil
}

func (s *CognitoService) ResendOtpByEmail(ctx context.Context, userID uuid.UUID) error {
	_, err := s.client.ResendConfirmationCode(ctx, &cognitoidentityprovider.ResendConfirmationCodeInput{
		Username: aws.String(userID.String()),
		ClientId: aws.String(s.userPoolClientId),
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *CognitoService) LoginWithEmail(ctx context.Context, email string, password string) (*securityDomain.TokenEntity, error) {
	result, err := s.client.InitiateAuth(ctx, &cognitoidentityprovider.InitiateAuthInput{
		AuthFlow: types.AuthFlowTypeUserPasswordAuth,
		ClientId: aws.String(s.userPoolClientId),
		AuthParameters: map[string]string{
			"EMAIL":    email,
			"PASSWORD": password,
		},
	})
	if err != nil {
		if strings.Contains(err.Error(), "UserNotFoundException") {
			return nil, securityDomain.ErrUserNotFoundException
		}
		if strings.Contains(err.Error(), "NotAuthorizedException") {
			return nil, securityDomain.ErrInvalidCredentials
		}
		if strings.Contains(err.Error(), "UserNotConfirmedException") {
			return nil, securityDomain.ErrUserNotConfirmed
		}
		return nil, err
	}

	tokenEntity := securityDomain.NewTokenEntity(*result.AuthenticationResult.AccessToken, *result.AuthenticationResult.RefreshToken, *result.AuthenticationResult.IdToken, 0)
	return tokenEntity, nil

}

func (s *CognitoService) RefreshToken(ctx context.Context, refreshToken string) (*securityDomain.TokenEntity, error) {
	result, err := s.client.InitiateAuth(context.TODO(), &cognitoidentityprovider.InitiateAuthInput{
		AuthFlow: types.AuthFlowTypeRefreshTokenAuth,
		ClientId: aws.String(s.userPoolClientId),
		AuthParameters: map[string]string{
			"REFRESH_TOKEN": refreshToken,
		},
	})
	if err != nil {
		if strings.Contains(err.Error(), "NotAuthorizedException") {
			return nil, securityDomain.ErrRefreshTokenExpired
		}
		if strings.Contains(err.Error(), "InvalidParameterException ") {
			return nil, securityDomain.ErrInvalidRefreshToken
		}
		if strings.Contains(err.Error(), "UserNotFoundException") {
			return nil, securityDomain.ErrUserNotFoundException
		}
		return nil, err
	}

	tokenEntity := securityDomain.NewTokenEntity(*result.AuthenticationResult.AccessToken, *result.AuthenticationResult.RefreshToken, *result.AuthenticationResult.IdToken, 0)
	return tokenEntity, nil

}

func (s *CognitoService) Logout(ctx context.Context, accessToken string) (bool, error) {
	_, err := s.client.GlobalSignOut(context.TODO(), &cognitoidentityprovider.GlobalSignOutInput{
		AccessToken: aws.String(accessToken),
	})
	if err != nil {
		if strings.Contains(err.Error(), "NotAuthorizedException") {
			return false, securityDomain.ErrInvalidAccessToken
		}
		return false, err
	}

	return true, nil

}

/*
func (s *CognitoService) RequestSignUp(ctx context.Context, phoneNumber string, userID uuid.UUID) (string, error) {

	// Create a new user in Cognito without password
	input := &cognitoidentityprovider.AdminCreateUserInput{
		Username:   aws.String(userID.String()),
		UserPoolId: aws.String(s.userPoolId),
		DesiredDeliveryMediums: []types.DeliveryMediumType{
			types.DeliveryMediumTypeSms,
		},
		MessageAction: types.MessageActionTypeSuppress,
		UserAttributes: []types.AttributeType{
			{
				Name:  aws.String("phone_number"),
				Value: aws.String(phoneNumber),
			},
		},
	}

	_, err := s.client.AdminCreateUser(context.TODO(), input)
	if err != nil {
		if !strings.Contains(err.Error(), "UsernameExistsException") {
			return "", err
		}
	}

	// Send a code to the user's phone number
	result, err := s.client.InitiateAuth(context.TODO(), &cognitoidentityprovider.InitiateAuthInput{
		AuthFlow: types.AuthFlowTypeCustomAuth,
		ClientId: aws.String(s.userPoolClientId),
		AuthParameters: map[string]string{
			"USERNAME": userID.String(),
		},
	})
	if err != nil {
		return "", err
	}

	session := *result.Session

	return session, nil

}

func (s *CognitoService) VerifyRequestSignUp(ctx context.Context, code string, session string, userID uuid.UUID) (*securityDomain.TokenEntity, error) {
	result, err := s.client.RespondToAuthChallenge(context.TODO(), &cognitoidentityprovider.RespondToAuthChallengeInput{
		ChallengeName: types.ChallengeNameTypeCustomChallenge,
		ChallengeResponses: map[string]string{
			"USERNAME": userID.String(),
			"ANSWER":   code,
		},
		ClientId: aws.String(s.userPoolClientId),
		Session:  aws.String(session),
	})
	if err != nil {
		if strings.Contains(err.Error(), "NotAuthorizedException: Invalid session for the user.") {
			return nil, securityDomain.ErrInvalidSession
		}
		if strings.Contains(err.Error(), "UserNotFoundException") {
			return nil, securityDomain.ErrUserNotFoundException
		}
		return nil, err
	}
	if result.AuthenticationResult == nil {
		return nil, securityDomain.ErrInvalidOTP
	}

	s.client.AdminUpdateUserAttributes(ctx, &cognitoidentityprovider.AdminUpdateUserAttributesInput{
		UserPoolId: aws.String(s.userPoolId),
		Username:   aws.String(userID.String()),
		UserAttributes: []types.AttributeType{
			{
				Name:  aws.String("phone_number_verified"),
				Value: aws.String("true"),
			},
		},
	})

	tokenEntity := securityDomain.NewTokenEntity(
		*result.AuthenticationResult.AccessToken,
		*result.AuthenticationResult.RefreshToken,
		*result.AuthenticationResult.IdToken,
		int(result.AuthenticationResult.ExpiresIn),
	)

	return tokenEntity, nil
}

func (s *CognitoService) CompleteSignUp(ctx context.Context, userID uuid.UUID, email string, password string, username string) error {
	_, err := s.client.AdminUpdateUserAttributes(context.TODO(), &cognitoidentityprovider.AdminUpdateUserAttributesInput{
		Username:   aws.String(userID.String()),
		UserPoolId: aws.String(s.userPoolId),
		UserAttributes: []types.AttributeType{
			{
				Name:  aws.String("email"),
				Value: aws.String(email),
			},
			{
				Name:  aws.String("email_verified"),
				Value: aws.String("true"),
			},
			{
				Name:  aws.String("custom:external_id"),
				Value: aws.String(userID.String()),
			},
			{
				Name:  aws.String("custom:username"),
				Value: aws.String(username),
			},
		},
	})
	if err != nil {
		if strings.Contains(err.Error(), "UsernameExistsException") {
			return securityDomain.ErrUserAlreadyExists
		}
		return err
	}
	_, err = s.client.AdminSetUserPassword(context.TODO(), &cognitoidentityprovider.AdminSetUserPasswordInput{
		Username:   aws.String(userID.String()),
		UserPoolId: aws.String(s.userPoolId),
		Password:   aws.String(password),
		Permanent:  true,
	})
	if err != nil {
		return err
	}

	return nil
}


func (s *CognitoService) SendOtp(ctx context.Context, phoneNumber string) (string, error) {
	var password string
	randomUsername := uuid.New().String()
	lowerCharSet := "abcdedfghijklmnopqrst"
	upperCharSet := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	specialCharSet := "!@#$%&*"
	numberSet := "0123456789"
	for i := 0; i < 5; i++ {
		randomLowerChar := rand.Intn(len(lowerCharSet))
		randomUpperChar := rand.Intn(len(upperCharSet))
		randomSpecialChar := rand.Intn(len(specialCharSet))
		randomNumber := rand.Intn(len(numberSet))
		password += string(lowerCharSet[randomLowerChar])
		password += string(upperCharSet[randomUpperChar])
		password += string(specialCharSet[randomSpecialChar])
		password += string(numberSet[randomNumber])
	}
	username := "NoneUserName-" + randomUsername
	input := &cognitoidentityprovider.SignUpInput{
		Username: aws.String(username),
		ClientId: aws.String(s.userPoolClientId),
		Password: aws.String(password),
		UserAttributes: []types.AttributeType{
			{
				Name:  aws.String("phone_number"),
				Value: aws.String(phoneNumber),
			},
		},
	}

	_, err := s.client.SignUp(context.TODO(), input)
	if err != nil {
		log.Fatal("Error signing up user in cognito: ", err)
		return "", fmt.Errorf("error signing up user in cognito: %v", err)
	}

	return username, nil
}

func (s *CognitoService) SignUpUser(ctx context.Context, username string, email string, phoneNumber string, password string) error {
	input := &cognitoidentityprovider.SignUpInput{
		ClientId: aws.String(s.userPoolClientId),
		Username: aws.String(username),
		Password: aws.String(password),
		UserAttributes: []types.AttributeType{
			{Name: aws.String("email"), Value: aws.String(email)},
			{Name: aws.String("phone_number"), Value: aws.String(phoneNumber)},
		},
	}

	_, err := s.client.SignUp(context.TODO(), input)
	if err != nil {
		if strings.Contains(err.Error(), "UsernameExistsException") {
			return securityDomain.ErrUserAlreadyExists
		}
		return err
	}

	return nil
}

func (s *CognitoService) ConfirmOtp(ctx context.Context, username string, code string) error {
	_, err := s.client.ConfirmSignUp(context.TODO(), &cognitoidentityprovider.ConfirmSignUpInput{
		Username:         aws.String(username),
		ClientId:         aws.String(s.userPoolClientId),
		ConfirmationCode: aws.String(code),
	})
	if err != nil {
		if strings.Contains(err.Error(), "UserNotFoundException") {
			return securityDomain.ErrUserNotFoundException
		}
		if strings.Contains(err.Error(), "CodeMismatchException") {
			return securityDomain.ErrInvalidOTP
		}
		if strings.Contains(err.Error(), "NotAuthorizedException") && strings.Contains(err.Error(), "CONFIRMED") {
			return securityDomain.ErrUserAlreadyConfirmed
		}
		if strings.Contains(err.Error(), "UsernameExistsException") {
			return securityDomain.ErrUserAlreadyExists
		}
		return err
	}

	return nil
}

func (s *CognitoService) LoginWithEmail(ctx context.Context, email string, password string) (*securityDomain.TokenEntity, error) {
	result, err := s.client.InitiateAuth(context.TODO(), &cognitoidentityprovider.InitiateAuthInput{
		AuthFlow: types.AuthFlowTypeUserPasswordAuth,
		ClientId: aws.String(s.userPoolClientId),
		AuthParameters: map[string]string{
			"EMAIL":    email,
			"PASSWORD": password,
		},
	})
	if err != nil {
		if strings.Contains(err.Error(), "UserNotFoundException") {
			return nil, securityDomain.ErrUserNotFoundException
		}
		if strings.Contains(err.Error(), "UserNotConfirmedException") {
			return nil, securityDomain.ErrUserNotConfirmed
		}
		if strings.Contains(err.Error(), "NotAuthorizedException") {
			return nil, securityDomain.ErrInvalidCredentials
		}
		if strings.Contains(err.Error(), "UsernameExistsException") {
			return nil, securityDomain.ErrUserAlreadyExists
		}
		return nil, err
	}

	tokenEntity := securityDomain.NewTokenEntity(*result.AuthenticationResult.AccessToken, *result.AuthenticationResult.RefreshToken, *result.AuthenticationResult.IdToken, 0)
	return tokenEntity, nil
}

func (s *CognitoService) RefreshToken(ctx context.Context, refreshToken string) (*securityDomain.TokenEntity, error) {
	result, err := s.client.InitiateAuth(context.TODO(), &cognitoidentityprovider.InitiateAuthInput{
		AuthFlow: types.AuthFlowTypeRefreshTokenAuth,
		ClientId: aws.String(s.userPoolClientId),
		AuthParameters: map[string]string{
			"REFRESH_TOKEN": refreshToken,
		},
	})
	if err != nil {
		if strings.Contains(err.Error(), "NotAuthorizedException") {
			return nil, securityDomain.ErrRefreshTokenExpired
		}
		if strings.Contains(err.Error(), "InvalidParameterException ") {
			return nil, securityDomain.ErrInvalidRefreshToken
		}
		if strings.Contains(err.Error(), "UserNotFoundException") {
			return nil, securityDomain.ErrUserNotFoundException
		}
		return nil, err
	}

	tokenEntity := securityDomain.NewTokenEntity(*result.AuthenticationResult.AccessToken, *result.AuthenticationResult.RefreshToken, *result.AuthenticationResult.IdToken, 0)
	return tokenEntity, nil
}

func (s *CognitoService) Logout(ctx context.Context, accessToken string) (bool, error) {
	_, err := s.client.GlobalSignOut(context.TODO(), &cognitoidentityprovider.GlobalSignOutInput{
		AccessToken: aws.String(accessToken),
	})
	if err != nil {
		if strings.Contains(err.Error(), "NotAuthorizedException") {
			return false, securityDomain.ErrInvalidAccessToken
		}
		return false, err
	}

	return true, nil
}
*/
