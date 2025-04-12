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

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type RefreshTokenResponse struct {
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

// @Summary Refresh access token
// @Description Refresh access token
// @Tags security
// @Accept  json
// @Produce  json
// @Param refresh_token body RefreshTokenRequest true "Refresh token"
// @Success 200 {object} RefreshTokenResponse
// @Failure 400 {object} common.Response[any]
// @Failure 401 {object} common.Response[any]
// @Failure 404 {object} common.Response[any]
// @Failure 500 {object} common.Response[any]
// @Router /auth/sessions [patch]
func refreshTokenLambdaHandler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var refreshTokenRequest RefreshTokenRequest
	err := json.Unmarshal([]byte(request.Body), &refreshTokenRequest)
	if err != nil {
		return utils.JsonResponse[any](400, "", nil, err.Error())
	}
	err = utils.ValidateRequest(refreshTokenRequest)
	if err != nil {
		return utils.JsonResponse[any](400, "", nil, err.Error())
	}

	token, err := app.NewRefreshTokenUseCase(cognitoService).RefreshToken(ctx, refreshTokenRequest.RefreshToken)
	switch true {
	case err == nil:
		return utils.JsonResponse(200, "", RefreshTokenResponse{
			AccessToken:  token.AccessToken,
			RefreshToken: token.RefreshToken,
			IdToken:      token.IdToken,
			ExpiresIn:    token.ExpiresIn,
		}, "")
	case errors.Is(err, securityDomain.ErrRefreshTokenExpired):
		return utils.JsonResponse[any](401, "", nil, err.Error())
	case errors.Is(err, securityDomain.ErrInvalidRefreshToken):
		return utils.JsonResponse[any](400, "", nil, err.Error())
	case errors.Is(err, securityDomain.ErrUserNotFoundException):
		return utils.JsonResponse[any](404, "", nil, err.Error())
	default:
		return utils.JsonResponse[any](500, "", nil, err.Error())
	}
}

func main() {
	lambda.Start(refreshTokenLambdaHandler)
}
