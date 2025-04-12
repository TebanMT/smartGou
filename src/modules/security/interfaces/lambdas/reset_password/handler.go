package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"os"

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
	"gorm.io/gorm"
)

type ResetPasswordRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
	Code     string `json:"code" validate:"required"`
}

type ResetPasswordResponse struct {
	Success bool `json:"success"`
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

// @Summary Reset password
// @Description Reset password. If the code, password and email are valid, the password will be reset.
// @Tags security
// @Accept  json
// @Produce  json
// @Param payload body ResetPasswordRequest true "Reset password request"
// @Success 200 {object} common.Response[ResetPasswordResponse]
// @Failure 400 {object} common.Response[any]
// @Failure 404 {object} common.Response[any]
// @Failure 401 with code expired {object} common.Response[any]
// @Failure 500 {object} common.Response[any]
// @Router /auth/recovery-password [patch]
func ResetPasswordHandler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var resetPasswordRequest ResetPasswordRequest
	err := json.Unmarshal([]byte(request.Body), &resetPasswordRequest)
	if err != nil {
		return utils.JsonResponse[any](400, "", nil, err.Error())
	}

	err = utils.ValidateRequest(resetPasswordRequest)
	if err != nil {
		return utils.JsonResponse[any](400, "", nil, err.Error())
	}

	passwordRecoveryUseCase := app.NewPasswordRecoveryUseCase(cognitoService, userRepository, unitOfWork)
	success, err := passwordRecoveryUseCase.ResetPassword(ctx, resetPasswordRequest.Email, resetPasswordRequest.Password, resetPasswordRequest.Code)
	switch true {
	case err == nil:
		return utils.JsonResponse(200, "", ResetPasswordResponse{Success: success}, "")
	case errors.Is(err, securityDomain.ErrUserNotFoundException), errors.Is(err, userDomain.ErrInvalidEmail):
		return utils.JsonResponse[any](404, "", nil, err.Error())
	case errors.Is(err, securityDomain.ErrInvalidOTP), errors.Is(err, userDomain.ErrInvalidEmail),
		errors.Is(err, userDomain.ErrEmailRequired), errors.Is(err, userDomain.ErrPasswordRequired),
		errors.Is(err, userDomain.ErrPasswordTooShort), errors.Is(err, userDomain.ErrInvalidPassword):
		return utils.JsonResponse[any](400, "", nil, err.Error())
	case errors.Is(err, securityDomain.ErrExpiredOTP):
		return utils.JsonResponse[any](401, "", nil, err.Error())
	default:
		return utils.JsonResponse[any](500, "", nil, err.Error())
	}
}

func main() {
	lambda.Start(ResetPasswordHandler)
}
