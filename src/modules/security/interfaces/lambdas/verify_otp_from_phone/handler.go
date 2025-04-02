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
	commonDomain "github.com/TebanMT/smartGou/src/common/domain"
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

type VerifyOTPByPhoneRequest struct {
	Code    string    `json:"code" validate:"required"`
	Session string    `json:"session" validate:"required"`
	UserID  uuid.UUID `json:"user_id" validate:"required"`
}

type VerifyOTPByPhoneResponse struct {
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
	unitOfWork = common.NewUnitOfWork(dbInstance)
	userRepository = repositories.NewUserRepository()
	cognitoService, err = cognito.NewCognitoService(os.Getenv("COGNITO_USER_POOL_ID"), os.Getenv("COGNITO_USER_POOL_CLIENT_ID"))
	if err != nil {
		log.Fatal("Error initializing cognito service: ", err)
	}
}

// @Summary Verify OTP from phone
// @Description  This endpoint verifies a OTP by phone number.
// @Description  If the OTP is valid, the user will be logged in and return the token entity.
// @Description  If the OTP is invalid, the user will be logged in and return the login challenge entity.
// @Description  The login challenge entity contains the session and the userID with which the user can try to verify the OTP again.
// @Tags security
// @Accept  json
// @Produce  json
// @Param payload body VerifyOTPByPhoneRequest true "Request body"
// @Success 200 {object} common.Response[VerifyOTPByPhoneResponse]
// @Failure 400 {object} common.Response[any]
// @Failure 401 {object} common.Response[any]
// @Failure 404 {object} common.Response[any]
// @Failure 429 {object} common.Response[any]
// @Failure 500 {object} common.Response[any]
// @Router /auth/otp [patch]
func VerifyRequestSignUpHandler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var requestSignUpRequest VerifyOTPByPhoneRequest
	err := json.Unmarshal([]byte(request.Body), &requestSignUpRequest)
	if err != nil {
		return common.JsonResponse[any](400, "", nil, err.Error())
	}

	err = common.ValidateRequest(requestSignUpRequest)
	if err != nil {
		return common.JsonResponse[any](400, "", nil, err.Error())
	}

	newUserCase := app.NewVerifyOTPByPhone(cognitoService, userRepository, unitOfWork)
	tokenEntity, loginChallenge, err := newUserCase.VerifyOTPByPhone(ctx, requestSignUpRequest.Code, requestSignUpRequest.Session, requestSignUpRequest.UserID)

	switch true {
	case tokenEntity != nil:
		return common.JsonResponse[any](200, "Verification code verified successfully", tokenEntity, "")
	case errors.Is(err, securityDomain.ErrInvalidOTP):
		return common.JsonResponse[any](400, "", VerifyOTPByPhoneResponse{
			UserID:      loginChallenge.UserId,
			Session:     loginChallenge.Session,
			MaxAttempts: loginChallenge.MaxAttempts,
			ExpiresAt:   loginChallenge.ExpiresAt,
		}, "")
	case errors.Is(err, securityDomain.ErrMaxAttemptsReached):
		return common.JsonResponse[any](429, "", nil, err.Error())
	case errors.Is(err, securityDomain.ErrInvalidSession):
		return common.JsonResponse[any](401, "", nil, err.Error())
	case errors.Is(err, securityDomain.ErrUserNotFoundException):
		return common.JsonResponse[any](404, "", nil, err.Error())
	default:
		return common.JsonResponse[any](500, "", nil, err.Error())
	}
}

func main() {
	lambda.Start(VerifyRequestSignUpHandler)
}
