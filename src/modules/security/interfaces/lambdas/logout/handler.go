package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"os"

	"github.com/TebanMT/smartGou/src/modules/security/app"
	securityDomain "github.com/TebanMT/smartGou/src/modules/security/domain"
	"github.com/TebanMT/smartGou/src/modules/security/infrastructure/cognito"
	"github.com/TebanMT/smartGou/src/shared/utils"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type LogoutRequest struct {
	AccessToken string `json:"access_token" validate:"required"`
}

type LogoutResponse struct {
	Success bool `json:"success"`
}

var cognitoService *cognito.CognitoService
var err error

func init() {
	cognitoService, err = cognito.NewCognitoService(os.Getenv("COGNITO_USER_POOL_ID"), os.Getenv("COGNITO_USER_POOL_CLIENT_ID"))
	if err != nil {
		log.Fatal("Error initializing cognito service: ", err)
	}
}

// @Summary Logout
// @Description This endpoint is used to logout a user. It invalidates the refresh token not the access token.
// @Description The access token could be used until its expiration but you can´t create a new access token with the same refresh token.
// @Tags security
// @Accept  json
// @Produce  json
// @Param access_token body LogoutRequest true "Access token"
// @Success 200 {object} utils.Response[LogoutResponse]
// @Failure 400 {object} utils.Response[any]
// @Failure 401 {object} utils.Response[any]
// @Failure 404 {object} utils.Response[any]
// @Failure 500 {object} utils.Response[any]
// @Router /auth/sessions [delete]
func logoutLambdaHandler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var logoutRequest LogoutRequest
	err := json.Unmarshal([]byte(request.Body), &logoutRequest)
	if err != nil {
		return utils.JsonResponse[any](400, "", nil, err.Error())
	}

	err = utils.ValidateRequest(logoutRequest)
	if err != nil {
		return utils.JsonResponse[any](400, "", nil, err.Error())
	}

	success, err := app.NewLogoutUseCase(cognitoService).Logout(ctx, logoutRequest.AccessToken)
	switch true {
	case err == nil:
		return utils.JsonResponse(200, "", LogoutResponse{
			Success: success,
		}, "")
	case errors.Is(err, securityDomain.ErrInvalidAccessToken):
		return utils.JsonResponse[any](400, "", nil, err.Error())
	case errors.Is(err, securityDomain.ErrUserNotFoundException):
		return utils.JsonResponse[any](404, "", nil, err.Error())
	default:
		return utils.JsonResponse[any](500, "", nil, err.Error())
	}
}

func main() {
	lambda.Start(logoutLambdaHandler)
}
