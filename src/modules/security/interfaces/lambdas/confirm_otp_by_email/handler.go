package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"os"

	"github.com/TebanMT/smartGou/infraestructure/db"
	"github.com/TebanMT/smartGou/src/common"
	"github.com/TebanMT/smartGou/src/modules/security/app"
	securityDomain "github.com/TebanMT/smartGou/src/modules/security/domain"
	"github.com/TebanMT/smartGou/src/modules/security/infrastructure/cognito"
	"github.com/TebanMT/smartGou/src/modules/users/infraestructure/db/repositories"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ConfirmOtpRequest struct {
	UserID uuid.UUID `json:"user_id" validate:"required"`
	Code   string    `json:"code" validate:"required"`
}

var cognitoService *cognito.CognitoService
var err error
var dbInstance *gorm.DB

func init() {
	cognitoService, err = cognito.NewCognitoService(os.Getenv("COGNITO_USER_POOL_ID"), os.Getenv("COGNITO_USER_POOL_CLIENT_ID"))
	if err != nil {
		log.Fatal("Error initializing cognito service: ", err)
	}
	dbInstance = db.InitConnection()
}

// @Summary Login with email and password
// @Description  This endpoint verifies a OTP by email.
// @Description  Verify the code given a userID and a code. The code is the code that the user received in the email.
// @Description  Also, the email will be verified as true in provider (cognito) and database.
// @Tags security
// @Accept  json
// @Produce  json
// @Param payload body ConfirmOtpRequest true "Request body"
// @Success 200 {object} common.Response[any]
// @Failure 400 {object} common.Response[any]
// @Failure 404 {object} common.Response[any]
// @Failure 409 {object} common.Response[any]
// @Failure 500 {object} common.Response[any]
// @Router /auth [patch]
func confirmOtpLambdaHandler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	var confirmOtpRequest ConfirmOtpRequest
	err := json.Unmarshal([]byte(request.Body), &confirmOtpRequest)
	if err != nil {
		return common.JsonResponse[any](400, "", nil, err.Error())
	}

	err = common.ValidateRequest(confirmOtpRequest)
	if err != nil {
		return common.JsonResponse[any](400, "", nil, err.Error())
	}

	unitOfWork := common.NewUnitOfWork(dbInstance)
	userRepository := repositories.NewUserRepository()
	err = app.NewVerifyOTPByEmail(cognitoService, userRepository, unitOfWork).VerifyOTPByEmail(ctx, confirmOtpRequest.UserID, confirmOtpRequest.Code)
	switch true {
	case err == nil:
		return common.JsonResponse[any](200, "", nil, "")
	case errors.Is(err, securityDomain.ErrUserNotFoundException):
		return common.JsonResponse[any](404, "", nil, err.Error())
	case errors.Is(err, securityDomain.ErrInvalidOTP):
		return common.JsonResponse[any](400, "", nil, err.Error())
	case errors.Is(err, securityDomain.ErrUserAlreadyConfirmed):
		return common.JsonResponse[any](409, "", nil, err.Error())
	default:
		return common.JsonResponse[any](500, "", nil, err.Error())
	}
}

func main() {
	lambda.Start(confirmOtpLambdaHandler)
}
