package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/TebanMT/smartGou/infraestructure/db"
	"github.com/TebanMT/smartGou/src/modules/security/infrastructure/cognito"
	userApp "github.com/TebanMT/smartGou/src/modules/users/app"
	"github.com/TebanMT/smartGou/src/modules/users/domain"
	userRepositories "github.com/TebanMT/smartGou/src/modules/users/infraestructure/db/repositories"
	commonDomain "github.com/TebanMT/smartGou/src/shared/domain"
	"github.com/TebanMT/smartGou/src/shared/middleware"
	"github.com/TebanMT/smartGou/src/shared/utils"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/google/uuid"
)

type UserPathRequest struct {
	ID uuid.UUID `json:"id" validate:"required"`
}

type UserResponse struct {
	UserID                uuid.UUID `json:"user_id"`
	Username              *string   `json:"username" extensions:"x-nullable=true"`
	Email                 *string   `json:"email" extensions:"x-nullable=true"`
	Phone                 *string   `json:"phone" extensions:"x-nullable=true"`
	IsOnboardingCompleted bool      `json:"is_onboarding_completed"`
	VerifiedPhone         bool      `json:"verified_phone"`
	VerifiedEmail         bool      `json:"verified_email"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
	Name                  *string   `json:"name" extensions:"x-nullable=true"`
	LastName              *string   `json:"last_name" extensions:"x-nullable=true"`
	SecondLastName        *string   `json:"second_last_name" extensions:"x-nullable=true"`
	DailingCode           *string   `json:"dailing_code" extensions:"x-nullable=true"`
}

var getUserProfileUseCase *userApp.GetUserProfile
var cognitoService *cognito.CognitoService
var unitOfWork commonDomain.UnitOfWork
var err error

func init() {
	dbInstance := db.InitConnection()
	unitOfWork = commonDomain.NewUnitOfWork(dbInstance)
	userRepository := userRepositories.NewUserRepository()
	getUserProfileUseCase = userApp.NewGetUserProfile(userRepository, unitOfWork)
	cognitoService, err = cognito.NewCognitoService(os.Getenv("COGNITO_USER_POOL_ID"), os.Getenv("COGNITO_USER_POOL_CLIENT_ID"))
	if err != nil {
		log.Fatal(err)
	}
}

// @Summary Get user profile
// @Description  This endpoint gets the user profile for a user.
// @Tags users
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param id-user path string true "User ID formatted as UUID"
// @Success 200 {object} utils.Response[UserResponse]
// @Failure 400 {object} utils.Response[any]
// @Failure 401 {object} utils.Response[any]
// @Failure 404 {object} utils.Response[any]
// @Failure 500 {object} utils.Response[any]
// @Router /users/{id-user} [get]
func getUserProfileLambdaHandler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var userPathRequest UserPathRequest
	id_user, err := uuid.Parse(request.PathParameters["id-user"])
	if err != nil {
		return utils.JsonResponse[any](400, "", nil, "User ID is not valid")
	}
	userPathRequest.ID = id_user

	err = utils.ValidateRequest(userPathRequest)
	if err != nil {
		return utils.JsonResponse[any](400, "", nil, err.Error())
	}

	user, err := getUserProfileUseCase.GetUserProfile(ctx, userPathRequest.ID)
	fmt.Println("user", user)
	switch true {
	case err == nil:
		return utils.JsonResponse(200, "", UserResponse{
			UserID:                user.UserID,
			Username:              user.Username,
			Email:                 user.Email,
			Phone:                 user.Phone,
			IsOnboardingCompleted: user.IsOnboardingCompleted,
			VerifiedPhone:         user.VerifiedPhone,
			VerifiedEmail:         user.VerifiedEmail,
			CreatedAt:             user.CreatedAt,
			UpdatedAt:             user.UpdatedAt,
			Name:                  user.Name,
			LastName:              user.LastName,
			SecondLastName:        user.SecondLastName,
			DailingCode:           user.DailingCode,
		}, "")
	case err == domain.ErrUserNotFound:
		return utils.JsonResponse[any](404, "", nil, err.Error())
	default:
		return utils.JsonResponse[any](500, "", nil, err.Error())
	}
}

func main() {
	handler := middleware.AuthenticationMiddleware(cognitoService, getUserProfileLambdaHandler)
	lambda.Start(handler)
}
