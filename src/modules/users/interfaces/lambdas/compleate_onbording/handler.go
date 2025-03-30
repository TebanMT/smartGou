package main

import (
	"context"
	"strconv"

	"github.com/TebanMT/smartGou/infraestructure/db"
	"github.com/TebanMT/smartGou/src/common"
	userApp "github.com/TebanMT/smartGou/src/modules/users/app"
	userRepositories "github.com/TebanMT/smartGou/src/modules/users/infraestructure/db/repositories"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type UserPathRequest struct {
	ID int `json:"id" validate:"required"`
}

var completeOnboardingUseCase *userApp.CompleteOnboardingUseCase

func init() {
	dbInstance := db.InitConnection()
	unitOfWork := common.NewUnitOfWork(dbInstance)
	userRepository := userRepositories.NewUserRepository()
	completeOnboardingUseCase = userApp.NewCompleteOnboardingUseCase(userRepository, unitOfWork)
}

func completeOnboardingLambdaHandler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var userPathRequest UserPathRequest
	id_user, err := strconv.Atoi(request.PathParameters["id-user"])
	if err != nil {
		return common.JsonResponse(400, "Bad Request: "+err.Error(), nil)
	}
	userPathRequest.ID = id_user

	err = common.ValidateRequest(userPathRequest)
	if err != nil {
		return common.JsonResponse(400, "Bad Request: "+err.Error(), nil)
	}

	err = completeOnboardingUseCase.CompleteOnboarding(ctx, userPathRequest.ID)
	if err != nil {
		return common.JsonResponse(500, "Internal Server Error: "+err.Error(), nil)
	}

	return common.JsonResponse(200, "Onboarding completed successfully", nil)
}

func main() {
	lambda.Start(completeOnboardingLambdaHandler)
}
