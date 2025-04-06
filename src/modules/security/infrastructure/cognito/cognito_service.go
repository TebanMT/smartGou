package cognito

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/MicahParks/keyfunc"
	securityDomain "github.com/TebanMT/smartGou/src/modules/security/domain"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

var (
	jwks     *keyfunc.JWKS
	jwksOnce sync.Once
)

type CognitoService struct {
	client           *cognitoidentityprovider.Client
	userPoolId       string
	userPoolClientId string
	region           string
}

func NewCognitoService(userPoolId string, userPoolClientId string) (*CognitoService, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	client := cognitoidentityprovider.NewFromConfig(cfg)

	return &CognitoService{client: client, userPoolId: userPoolId, userPoolClientId: userPoolClientId, region: os.Getenv("AWS_REGION")}, nil
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
			"USERNAME": email,
			"PASSWORD": password,
		},
	})
	if err != nil {
		fmt.Println("Error: ", err)
		if strings.Contains(err.Error(), "UserNotFoundException") {
			return nil, securityDomain.ErrInvalidCredentials
		}
		if strings.Contains(err.Error(), "NotAuthorizedException: Password attempts exceeded") {
			return nil, securityDomain.ErrMaxAttemptsReached
		}
		if strings.Contains(err.Error(), "NotAuthorizedException") {
			return nil, securityDomain.ErrInvalidCredentials
		}
		if strings.Contains(err.Error(), "UserNotConfirmedException") {
			return nil, securityDomain.ErrUserNotConfirmed
		}
		return nil, err
	}

	tokenEntity := securityDomain.NewTokenEntity(
		*result.AuthenticationResult.AccessToken,
		*result.AuthenticationResult.RefreshToken,
		*result.AuthenticationResult.IdToken,
		int(result.AuthenticationResult.ExpiresIn),
	)
	return tokenEntity, nil

}

func (s *CognitoService) RefreshToken(ctx context.Context, refreshToken string) (*securityDomain.TokenEntity, error) {
	result, err := s.client.InitiateAuth(ctx, &cognitoidentityprovider.InitiateAuthInput{
		AuthFlow: types.AuthFlowTypeRefreshTokenAuth,
		ClientId: aws.String(s.userPoolClientId),
		AuthParameters: map[string]string{
			"REFRESH_TOKEN": refreshToken,
		},
	})
	if err != nil {
		fmt.Println("Error: ", err)
		if strings.Contains(err.Error(), "NotAuthorizedException") {
			return nil, securityDomain.ErrInvalidRefreshToken
		}
		if strings.Contains(err.Error(), "InvalidParameterException ") {
			return nil, securityDomain.ErrInvalidRefreshToken
		}
		if strings.Contains(err.Error(), "UserNotFoundException") {
			return nil, securityDomain.ErrUserNotFoundException
		}
		return nil, err
	}

	tokenEntity := securityDomain.NewTokenEntity(*result.AuthenticationResult.AccessToken, refreshToken, *result.AuthenticationResult.IdToken, int(result.AuthenticationResult.ExpiresIn))
	return tokenEntity, nil

}

func (s *CognitoService) Logout(ctx context.Context, accessToken string) (bool, error) {
	_, err := s.client.GlobalSignOut(ctx, &cognitoidentityprovider.GlobalSignOutInput{
		AccessToken: aws.String(accessToken),
	})
	if err != nil {
		fmt.Println("Error: ", err)
		if strings.Contains(err.Error(), "NotAuthorizedException") {
			return false, securityDomain.ErrInvalidAccessToken
		}
		return false, err
	}

	return true, nil

}

func (s *CognitoService) ParseTokenAndValidate(ctx context.Context, token string) (*securityDomain.TokenClaims, error) {

	type TokenClaims struct {
		Sub      string `json:"sub"`
		Username string `json:"username"`
		Exp      int64  `json:"exp"`
		Iss      string `json:"iss"`
		Aud      string `json:"aud"`
		Iat      int64  `json:"iat"`
		jwt.RegisteredClaims
	}
	jwks, err := s.getJWKS()
	if err != nil {
		fmt.Println("Error getting JWKS: ", err)
		return nil, err
	}
	parsedToken, err := jwt.ParseWithClaims(token, &TokenClaims{}, jwks.Keyfunc)

	if err != nil {
		fmt.Println("Error parsing token: ", err)
		return nil, securityDomain.ErrInvalidToken
	}

	claims, ok := parsedToken.Claims.(*TokenClaims)
	if !ok {
		return nil, securityDomain.ErrInvalidToken
	}

	if !parsedToken.Valid {
		return nil, securityDomain.ErrInvalidToken
	}

	if claims.Exp < time.Now().Unix() {
		return nil, securityDomain.ErrTokenExpired
	}
	fmt.Println("Token kid:", parsedToken.Header["kid"])

	return securityDomain.NewTokenClaims(
		claims.Sub,
		claims.Username,
		claims.Exp,
		claims.Iss,
		claims.Aud,
		claims.Iat,
	), nil
}

func (s *CognitoService) getJWKS() (*keyfunc.JWKS, error) {
	var err error
	jwksOnce.Do(func() {
		jwksURL := fmt.Sprintf("https://cognito-idp.%s.amazonaws.com/%s/.well-known/jwks.json", s.region, s.userPoolId)
		fmt.Println("jwksURL: ", jwksURL)
		jwks, err = keyfunc.Get(jwksURL, keyfunc.Options{
			RefreshInterval:   time.Hour,
			RefreshUnknownKID: true,
		})
	})
	return jwks, err
}

func (s *CognitoService) PasswordRecovery(ctx context.Context, email string) (bool, error) {
	_, err := s.client.ForgotPassword(ctx, &cognitoidentityprovider.ForgotPasswordInput{
		ClientId: aws.String(s.userPoolClientId),
		Username: aws.String(email),
	})
	if err != nil {
		fmt.Println("Error: ", err)
		if strings.Contains(err.Error(), "UserNotFoundException") {
			return false, securityDomain.ErrUserNotFoundException
		}
		return false, err
	}
	return true, nil
}

func (s *CognitoService) PasswordReset(ctx context.Context, userID uuid.UUID, newPassword string, confirmationCode string) (bool, error) {
	_, err := s.client.ConfirmForgotPassword(ctx, &cognitoidentityprovider.ConfirmForgotPasswordInput{
		ClientId:         aws.String(s.userPoolClientId),
		Username:         aws.String(userID.String()),
		Password:         aws.String(newPassword),
		ConfirmationCode: aws.String(confirmationCode),
	})
	if err != nil {
		fmt.Println("Error: ", err)
		if strings.Contains(err.Error(), "CodeMismatchException") {
			return false, securityDomain.ErrInvalidOTP
		}
		if strings.Contains(err.Error(), "UserNotFoundException") {
			return false, securityDomain.ErrUserNotFoundException
		}
		if strings.Contains(err.Error(), "ExpiredCodeException") {
			return false, securityDomain.ErrExpiredOTP
		}
		return false, err
	}
	return true, nil
}
