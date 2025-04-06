package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"os"

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
	"gorm.io/gorm"
)

type RequestRecoveryPassword struct {
	Email string `json:"email" validate:"required,email"`
}

type RequestRecoveryPasswordResponse struct {
	Success bool `json:"success"`
}

var cognitoService *cognito.CognitoService
var err error
var dbInstance *gorm.DB
var unitOfWork commonDomain.UnitOfWork
var userRepository userDomain.UserRepository

func init() {
	cognitoService, err = cognito.NewCognitoService(os.Getenv("COGNITO_USER_POOL_ID"), os.Getenv("COGNITO_USER_POOL_CLIENT_ID"))
	if err != nil {
		log.Fatal("Error initializing cognito service: ", err)
	}
	dbInstance = db.InitConnection()
	unitOfWork = common.NewUnitOfWork(dbInstance)
	userRepository = repositories.NewUserRepository()
}

// @Summary Request password recovery
// @Description Request password recovery. If the email is found, the user will receive an email with a code to reset the password.
// @Description In otherwise, the user will receive an error.
// @Tags security
// @Accept  json
// @Produce  json
// @Param payload body RequestRecoveryPassword true "Request password recovery"
// @Success 200 {object} common.Response[RequestRecoveryPasswordResponse]
// @Failure 400 {object} common.Response[any]
// @Failure 404 {object} common.Response[any]
// @Failure 500 {object} common.Response[any]
// @Router /auth/recovery-password [post]
func RequestRecoveryPasswordHandler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var request RequestRecoveryPassword
	err := json.Unmarshal([]byte(event.Body), &request)
	if err != nil {
		return common.JsonResponse[any](400, "", nil, err.Error())
	}
	err = common.ValidateRequest(request)
	if err != nil {
		return common.JsonResponse[any](400, "", nil, err.Error())
	}
	passwordRecoveryUseCase := app.NewPasswordRecoveryUseCase(cognitoService, userRepository, unitOfWork)
	success, err := passwordRecoveryUseCase.RequestPasswordRecovery(ctx, request.Email)
	switch true {
	case err == nil:
		return common.JsonResponse(200, "", RequestRecoveryPasswordResponse{Success: success}, "")
	case errors.Is(err, securityDomain.ErrUserNotFoundException), errors.Is(err, userDomain.ErrUserNotFound):
		return common.JsonResponse[any](404, "", nil, err.Error())
	case errors.Is(err, userDomain.ErrInvalidEmail), errors.Is(err, userDomain.ErrEmailRequired):
		return common.JsonResponse[any](400, "", nil, err.Error())
	default:
		return common.JsonResponse[any](500, "", nil, err.Error())
	}
}

func main() {
	lambda.Start(RequestRecoveryPasswordHandler)
}
