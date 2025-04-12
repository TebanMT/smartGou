package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"os"
	"time"

	"github.com/TebanMT/smartGou/infraestructure/db"
	"github.com/TebanMT/smartGou/src/modules/security/app"
	securityDomain "github.com/TebanMT/smartGou/src/modules/security/domain"
	"github.com/TebanMT/smartGou/src/modules/security/infrastructure/cognito"
	userDomain "github.com/TebanMT/smartGou/src/modules/users/domain"
	"github.com/TebanMT/smartGou/src/modules/users/infraestructure/db/repositories"
	commonDomain "github.com/TebanMT/smartGou/src/shared/domain"
	"github.com/TebanMT/smartGou/src/shared/utils"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// RequestSignUpRequest is the request body for the sign up and login by phone
type RequestSignUpRequest struct {
	PhoneNumber string `json:"phone_number" validate:"required"`
	DailingCode string `json:"dailing_code" validate:"required"`
}

// RequestSignUpResponse is the response body for the sign up and login by phone
type RequestSignUpResponse struct {
	UserID      uuid.UUID `json:"user_id"`
	Session     string    `json:"session"`
	MaxAttempts int       `json:"max_attempts"`
	ExpiresAt   time.Time `json:"expires_at"`
}

var cognitoService *cognito.CognitoService
var err error

var dbInstance *gorm.DB
var unitOfWork commonDomain.UnitOfWork
var userRepository userDomain.UserRepository

func init() {
	dbInstance = db.InitConnection()
	unitOfWork = commonDomain.NewUnitOfWork(dbInstance)
	userRepository = repositories.NewUserRepository()
	cognitoService, err = cognito.NewCognitoService(os.Getenv("COGNITO_USER_POOL_ID"), os.Getenv("COGNITO_USER_POOL_CLIENT_ID"))
	if err != nil {
		log.Fatal("Error initializing cognito service: ", err)
	}
}

// @Summary Request sign up / login by phone
// @Description  Try to sign up or login by phone number.
// @Description  The user will be created in the database and inderity provider (cognito) if it doesn't exist.
// @Description  If the user exists, the code will be sent to the phone number.
// @Description  The user will be created without password and return the session.
// @Description  The session is specific context for the user to login in the cognito user pool. It is used to verify the user's phone number.
// @Description  For others implementations, this session could be not needed.
// @Description  The user ID (uuid type) is the user's ID in the database and also the username in the cognito user pool.
// @Tags security
// @Accept  json
// @Produce  json
// @Param payload body RequestSignUpRequest true "Request body"
// @Success 200 {object} common.Response[RequestSignUpResponse]
// @Failure 400 {object} common.Response[any]
// @Failure 409 {object} common.Response[any]
// @Failure 500 {object} common.Response[any]
// @Router /auth/otp [post]
func RequestSignUpHandler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var requestSignUpRequest RequestSignUpRequest
	err := json.Unmarshal([]byte(request.Body), &requestSignUpRequest)
	if err != nil {
		return utils.JsonResponse[any](400, "", nil, err.Error())
	}

	err = utils.ValidateRequest(requestSignUpRequest)
	if err != nil {
		return utils.JsonResponse[any](400, "", nil, err.Error())
	}

	loginChallenge, err := app.NewRequestOTPByPhone(cognitoService, userRepository, unitOfWork).RequestOTPByPhone(ctx, requestSignUpRequest.PhoneNumber, requestSignUpRequest.DailingCode)

	switch {
	case err == nil:
		return utils.JsonResponse(200, "Verification code sent successfully", RequestSignUpResponse{
			UserID:      loginChallenge.UserId,
			Session:     loginChallenge.Session,
			MaxAttempts: loginChallenge.MaxAttempts,
			ExpiresAt:   loginChallenge.ExpiresAt,
		}, "")
	case errors.Is(err, securityDomain.ErrInvalidPhoneOrDailingCode),
		errors.Is(err, securityDomain.ErrInvalidPhoneLength),
		errors.Is(err, securityDomain.ErrInvalidDailingCodeLength),
		errors.Is(err, securityDomain.ErrInvalidPhoneFormat),
		errors.Is(err, securityDomain.ErrInvalidDailingCodeFormat):
		return utils.JsonResponse[any](400, "", nil, err.Error())
	case errors.Is(err, userDomain.ErrPhoneAlreadyExists),
		errors.Is(err, userDomain.ErrEmailAlreadyExists),
		errors.Is(err, userDomain.ErrUsernameAlreadyExists),
		errors.Is(err, securityDomain.ErrPhoneAlreadyVerified),
		errors.Is(err, securityDomain.ErrUserAlreadyExists):
		return utils.JsonResponse[any](409, "", nil, err.Error())
	default:
		return utils.JsonResponse[any](500, "", nil, err.Error())
	}
}

func main() {
	lambda.Start(RequestSignUpHandler)
}
