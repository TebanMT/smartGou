package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"os"

	"github.com/TebanMT/smartGou/infraestructure/db"
	"github.com/TebanMT/smartGou/src/modules/security/app"
	"github.com/TebanMT/smartGou/src/modules/security/infrastructure/cognito"
	userDomain "github.com/TebanMT/smartGou/src/modules/users/domain"
	"github.com/TebanMT/smartGou/src/modules/users/infraestructure/db/repositories"
	commonDomain "github.com/TebanMT/smartGou/src/shared/domain"
	"github.com/TebanMT/smartGou/src/shared/utils"
	"gorm.io/gorm"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type SignUpByEmailRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type SignUpByEmailResponse struct {
	UserID string `json:"user_id"`
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

// @Summary Sign up by email and password
// @Description  This endpoint signs up a user by email and password.
// @Description  The user will be created in the database and inderity provider (cognito) if it doesn't exist.
// @Description  If the user exists, an error will be returned.
// @Description  A mail will be sent to the user to verify the email.
// @Description  If the user exists but the email is not verified, a new OTP will be sent to the user.
// @Tags security
// @Accept  json
// @Produce  json
// @Param payload body SignUpByEmailRequest true "Request body"
// @Success 200 {object} utils.Response[SignUpByEmailResponse]
// @Failure 400 {object} utils.Response[any]
// @Failure 409 {object} utils.Response[any]
// @Failure 500 {object} utils.Response[any]
// @Router /auth [post]
func signUpByEmailLambdaHandler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	var signUpByEmailRequest SignUpByEmailRequest
	err := json.Unmarshal([]byte(request.Body), &signUpByEmailRequest)
	if err != nil {
		return utils.JsonResponse[any](400, "", nil, err.Error())
	}

	err = utils.ValidateRequest(signUpByEmailRequest)
	if err != nil {
		return utils.JsonResponse[any](400, "", nil, err.Error())
	}

	userID, err := app.NewSignUpByEmailUseCase(cognitoService, userRepository, unitOfWork).SignUpByEmail(ctx, signUpByEmailRequest.Email, signUpByEmailRequest.Password)
	switch {
	case err == nil:
		return utils.JsonResponse(200, "", SignUpByEmailResponse{UserID: userID}, "")
	case errors.Is(err, userDomain.ErrInvalidEmail),
		errors.Is(err, userDomain.ErrInvalidPassword),
		errors.Is(err, userDomain.ErrPasswordTooShort):
		return utils.JsonResponse[any](400, "", nil, err.Error())
	case errors.Is(err, userDomain.ErrEmailAlreadyExists):
		return utils.JsonResponse[any](409, "", nil, err.Error())
	default:
		return utils.JsonResponse[any](500, "", nil, err.Error())
	}
}

func main() {
	lambda.Start(signUpByEmailLambdaHandler)
}
