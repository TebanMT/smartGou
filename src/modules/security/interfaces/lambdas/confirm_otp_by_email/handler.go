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

func confirmOtpLambdaHandler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	var confirmOtpRequest ConfirmOtpRequest
	err := json.Unmarshal([]byte(request.Body), &confirmOtpRequest)
	if err != nil {
		return common.JsonResponse(400, "Bad Request: "+err.Error(), nil)
	}

	err = common.ValidateRequest(confirmOtpRequest)
	if err != nil {
		return common.JsonResponse(400, "Bad Request: "+err.Error(), nil)
	}

	unitOfWork := common.NewUnitOfWork(dbInstance)
	userRepository := repositories.NewUserRepository()
	err = app.NewVerifyOTPByEmail(cognitoService, userRepository, unitOfWork).VerifyOTPByEmail(ctx, confirmOtpRequest.UserID, confirmOtpRequest.Code)
	switch true {
	case err == nil:
		return common.JsonResponse(200, "OTP confirmed successfully", nil)
	case errors.Is(err, securityDomain.ErrUserNotFoundException):
		return common.JsonResponse(404, "User not found", nil)
	case errors.Is(err, securityDomain.ErrInvalidOTP):
		return common.JsonResponse(400, "Invalid OTP", nil)
	case errors.Is(err, securityDomain.ErrUserAlreadyConfirmed):
		return common.JsonResponse(409, "User already confirmed", nil)
	default:
		return common.JsonResponse(500, "Internal Server Error: "+err.Error(), nil)
	}
}

func main() {
	lambda.Start(confirmOtpLambdaHandler)
}
