package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"os"
	"time"

	"github.com/TebanMT/smartGou/infraestructure/db"
	"github.com/TebanMT/smartGou/src/common"
	"github.com/TebanMT/smartGou/src/common/domain"
	"github.com/TebanMT/smartGou/src/modules/security/app"
	securityDomain "github.com/TebanMT/smartGou/src/modules/security/domain"
	"github.com/TebanMT/smartGou/src/modules/security/infrastructure/cognito"
	userDomain "github.com/TebanMT/smartGou/src/modules/users/domain"
	"github.com/TebanMT/smartGou/src/modules/users/infraestructure/db/repositories"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RequestSignUpRequest struct {
	PhoneNumber string `json:"phone_number" validate:"required"`
	DailingCode string `json:"dailing_code" validate:"required"`
}

type RequestSignUpResponse struct {
	UserID      uuid.UUID `json:"user_id"`
	Session     string    `json:"session"`
	MaxAttempts int       `json:"max_attempts"`
	ExpiresAt   time.Time `json:"expires_at"`
}

var cognitoService *cognito.CognitoService
var err error

var dbInstance *gorm.DB
var unitOfWork domain.UnitOfWork
var userRepository userDomain.UserRepository

func init() {
	dbInstance = db.InitConnection()
	unitOfWork = common.NewUnitOfWork(dbInstance)
	userRepository = repositories.NewUserRepository()
	cognitoService, err = cognito.NewCognitoService(os.Getenv("COGNITO_USER_POOL_ID"), os.Getenv("COGNITO_USER_POOL_CLIENT_ID"))
	if err != nil {
		log.Fatal("Error initializing cognito service: ", err)
	}
}

func RequestSignUpHandler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var requestSignUpRequest RequestSignUpRequest
	err := json.Unmarshal([]byte(request.Body), &requestSignUpRequest)
	if err != nil {
		return common.JsonResponse(400, err.Error(), nil)
	}

	err = common.ValidateRequest(requestSignUpRequest)
	if err != nil {
		return common.JsonResponse(400, err.Error(), nil)
	}

	loginChallenge, err := app.NewRequestOTPByPhone(cognitoService, userRepository, unitOfWork).RequestOTPByPhone(ctx, requestSignUpRequest.PhoneNumber, requestSignUpRequest.DailingCode)

	switch {
	case err == nil:
		return common.JsonResponse(200, "Verification code sent successfully", RequestSignUpResponse{
			UserID:      loginChallenge.UserId,
			Session:     loginChallenge.Session,
			MaxAttempts: loginChallenge.MaxAttempts,
			ExpiresAt:   loginChallenge.ExpiresAt,
		})
	case errors.Is(err, securityDomain.ErrInvalidPhoneOrDailingCode),
		errors.Is(err, securityDomain.ErrInvalidPhoneLength),
		errors.Is(err, securityDomain.ErrInvalidDailingCodeLength),
		errors.Is(err, securityDomain.ErrInvalidPhoneFormat),
		errors.Is(err, securityDomain.ErrInvalidDailingCodeFormat):
		return common.JsonResponse(400, err.Error(), nil)
	case errors.Is(err, userDomain.ErrPhoneAlreadyExists),
		errors.Is(err, userDomain.ErrEmailAlreadyExists),
		errors.Is(err, userDomain.ErrUsernameAlreadyExists),
		errors.Is(err, securityDomain.ErrPhoneAlreadyVerified),
		errors.Is(err, securityDomain.ErrUserAlreadyExists):
		return common.JsonResponse(409, err.Error(), nil)
	default:
		return common.JsonResponse(500, "Internal Server Error: "+err.Error(), nil)
	}
}

func main() {
	lambda.Start(RequestSignUpHandler)
}
