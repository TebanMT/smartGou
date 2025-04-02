package main

import (
	"context"
	"log"
	"os"

	"github.com/TebanMT/smartGou/infraestructure/db"
	"github.com/TebanMT/smartGou/src/common"
	securityDomain "github.com/TebanMT/smartGou/src/modules/security/domain"
	"github.com/TebanMT/smartGou/src/modules/security/infrastructure/cognito"
	userApp "github.com/TebanMT/smartGou/src/modules/users/app"
	userRepositories "github.com/TebanMT/smartGou/src/modules/users/infraestructure/db/repositories"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/google/uuid"
)

type UserPathRequest struct {
	ID uuid.UUID `json:"id" validate:"required"`
}

var completeOnboardingUseCase *userApp.CompleteOnboardingUseCase
var cognitoService *cognito.CognitoService
var err error

func init() {
	dbInstance := db.InitConnection()
	unitOfWork := common.NewUnitOfWork(dbInstance)
	userRepository := userRepositories.NewUserRepository()
	completeOnboardingUseCase = userApp.NewCompleteOnboardingUseCase(userRepository, unitOfWork)
	cognitoService, err = cognito.NewCognitoService(os.Getenv("COGNITO_USER_POOL_ID"), os.Getenv("COGNITO_USER_POOL_CLIENT_ID"))
	if err != nil {
		log.Fatal("Error initializing cognito service: ", err)
	}
}

// @Summary Complete onboarding
// @Description  This endpoint completes the onboarding process for a user.
// @Tags users
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param id-user path string true "User ID formatted as UUID"
// @Success 200 {object} common.Response[any]
// @Failure 400 {object} common.Response[any]
// @Failure 401 {object} common.Response[any]
// @Failure 500 {object} common.Response[any]
// @Router /users/{id-user}/onboarding [patch]
func completeOnboardingLambdaHandler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var userPathRequest UserPathRequest
	claims := request.RequestContext.Authorizer["claims"].(*securityDomain.TokenClaims)
	id_user_claims, err := uuid.Parse(claims.UserId)
	if err != nil {
		return common.JsonResponse[any](400, "", nil, err.Error())
	}
	id_user, err := uuid.Parse(request.PathParameters["id-user"])
	if err != nil || id_user != id_user_claims {
		return common.JsonResponse[any](400, "", nil, err.Error())
	}
	userPathRequest.ID = id_user

	err = common.ValidateRequest(userPathRequest)
	if err != nil {
		return common.JsonResponse[any](400, "", nil, err.Error())
	}

	err = completeOnboardingUseCase.CompleteOnboarding(ctx, userPathRequest.ID)
	if err != nil {
		return common.JsonResponse[any](500, "", nil, err.Error())
	}

	return common.JsonResponse[any](200, "", nil, "")
}

func main() {
	handler := common.AuthenticationMiddleware(cognitoService, completeOnboardingLambdaHandler)
	lambda.Start(handler)
}
