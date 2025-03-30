package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"os"

	"github.com/TebanMT/smartGou/src/common"
	"github.com/TebanMT/smartGou/src/modules/security/app"
	securityDomain "github.com/TebanMT/smartGou/src/modules/security/domain"
	"github.com/TebanMT/smartGou/src/modules/security/infrastructure/cognito"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type LoginWithEmailRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

var cognitoService *cognito.CognitoService
var err error

func init() {
	cognitoService, err = cognito.NewCognitoService(os.Getenv("COGNITO_USER_POOL_ID"), os.Getenv("COGNITO_USER_POOL_CLIENT_ID"))
	if err != nil {
		log.Fatal("Error initializing cognito service: ", err)
	}
}

func LoginWithEmailHandler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var loginWithEmailRequest LoginWithEmailRequest
	err := json.Unmarshal([]byte(request.Body), &loginWithEmailRequest)
	if err != nil {
		return common.JsonResponse(400, err.Error(), nil)
	}

	loginWithEmail := app.NewLoginWithEmail(cognitoService)
	token, err := loginWithEmail.LoginWithEmail(ctx, loginWithEmailRequest.Email, loginWithEmailRequest.Password)
	switch true {
	case err == nil:
		return common.JsonResponse(200, "Login successful", token)
	case errors.Is(err, securityDomain.ErrUserNotFoundException):
		return common.JsonResponse(400, err.Error(), nil)
	case errors.Is(err, securityDomain.ErrInvalidCredentials):
		return common.JsonResponse(401, err.Error(), nil)
	case errors.Is(err, securityDomain.ErrUserNotConfirmed):
		return common.JsonResponse(403, err.Error(), nil)
	case errors.Is(err, securityDomain.ErrMaxAttemptsReached):
		return common.JsonResponse(403, err.Error(), nil)
	default:
		return common.JsonResponse(500, err.Error(), nil)
	}

}

func main() {
	lambda.Start(LoginWithEmailHandler)
}
