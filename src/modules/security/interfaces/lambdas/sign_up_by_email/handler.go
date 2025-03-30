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
	"github.com/TebanMT/smartGou/src/modules/security/infrastructure/cognito"
	userDomain "github.com/TebanMT/smartGou/src/modules/users/domain"
	"github.com/TebanMT/smartGou/src/modules/users/infraestructure/db/repositories"
	"gorm.io/gorm"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type SignUpByEmailRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
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

func signUpByEmailLambdaHandler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	var signUpByEmailRequest SignUpByEmailRequest
	err := json.Unmarshal([]byte(request.Body), &signUpByEmailRequest)
	if err != nil {
		return common.JsonResponse(400, "Bad Request: "+err.Error(), nil)
	}

	err = common.ValidateRequest(signUpByEmailRequest)
	if err != nil {
		return common.JsonResponse(400, "Bad Request: "+err.Error(), nil)
	}

	userID, err := app.NewSignUpByEmailUseCase(cognitoService, userRepository, unitOfWork).SignUpByEmail(ctx, signUpByEmailRequest.Email, signUpByEmailRequest.Password)
	switch {
	case err == nil:
		return common.JsonResponse(200, "User created successfully", map[string]string{"user_id": userID})
	case errors.Is(err, userDomain.ErrInvalidEmail),
		errors.Is(err, userDomain.ErrInvalidPassword),
		errors.Is(err, userDomain.ErrPasswordTooShort):
		return common.JsonResponse(400, err.Error(), nil)
	case errors.Is(err, userDomain.ErrEmailAlreadyExists):
		return common.JsonResponse(409, err.Error(), nil)
	default:
		return common.JsonResponse(500, "Internal Server Error: "+err.Error(), nil)
	}
}

func main() {
	lambda.Start(signUpByEmailLambdaHandler)
}
