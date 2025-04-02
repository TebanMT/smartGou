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

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	IdToken      string `json:"id_token"`
	ExpiresIn    int    `json:"expires_in"`
}

var cognitoService *cognito.CognitoService
var err error

func init() {
	cognitoService, err = cognito.NewCognitoService(os.Getenv("COGNITO_USER_POOL_ID"), os.Getenv("COGNITO_USER_POOL_CLIENT_ID"))
	if err != nil {
		log.Fatal("Error initializing cognito service: ", err)
	}
}

// @Summary Login with email and password
// @Description  Login with email
// @Description  This use case is used to login a user with email and password.
// @Description  The user will be logged in if the email and password are correct, otherwise an error will be returned.
// @Description  The token is a JWT token that is used to authenticate the user.
// @Description  The token is valid for 1 hour.
// @Tags security
// @Accept  json
// @Produce  json
// @Param payload body LoginWithEmailRequest true "Request body"
// @Success 200 {object} common.Response[TokenResponse]
// @Failure 400 {object} common.Response[any]
// @Failure 401 {object} common.Response[any]
// @Failure 403 {object} common.Response[any]
// @Failure 500 {object} common.Response[any]
// @Router /auth/sessions [post]
func LoginWithEmailHandler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var loginWithEmailRequest LoginWithEmailRequest
	err := json.Unmarshal([]byte(request.Body), &loginWithEmailRequest)
	if err != nil {
		return common.JsonResponse[any](400, "", nil, err.Error())
	}

	loginWithEmail := app.NewLoginWithEmail(cognitoService)
	token, err := loginWithEmail.LoginWithEmail(ctx, loginWithEmailRequest.Email, loginWithEmailRequest.Password)
	switch true {
	case err == nil:
		return common.JsonResponse(200, "", token, "")
	case errors.Is(err, securityDomain.ErrUserNotFoundException):
		return common.JsonResponse[any](400, "", nil, err.Error())
	case errors.Is(err, securityDomain.ErrInvalidCredentials):
		return common.JsonResponse[any](401, "", nil, err.Error())
	case errors.Is(err, securityDomain.ErrUserNotConfirmed):
		return common.JsonResponse[any](403, "", nil, err.Error())
	case errors.Is(err, securityDomain.ErrMaxAttemptsReached):
		return common.JsonResponse[any](403, "", nil, err.Error())
	default:
		return common.JsonResponse[any](500, "", nil, err.Error())
	}

}

func main() {
	lambda.Start(LoginWithEmailHandler)
}
