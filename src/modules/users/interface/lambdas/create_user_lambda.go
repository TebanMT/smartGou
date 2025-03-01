package lambdas

import (
	"context"
	"encoding/json"

	"github.com/TebanMT/smartGou/infraestructure/db"
	"github.com/TebanMT/smartGou/src/common"
	"github.com/TebanMT/smartGou/src/modules/users/domain"
	"github.com/TebanMT/smartGou/src/modules/users/infraestructure/db/repositories"
	"github.com/aws/aws-lambda-go/events"
	"gorm.io/gorm"
)

type UserRequest struct {
	Username       string `json:"username" validate:"required,min=3,max=50"`
	Name           string `json:"name" validate:"required,min=3,max=50"`
	LastName       string `json:"last_name" validate:"required,min=3,max=50"`
	SecondLastName string `json:"second_last_name" validate:"required,min=3,max=50"`
	Email          string `json:"email" validate:"required,email"`
	DailingCode    string `json:"dailing_code" validate:"required,min=1,max=5"`
	Phone          string `json:"phone" validate:"required"`
}

var dbInstance *gorm.DB

func init() {
	dbInstance = db.InitConnection()
}

func CreateUserLambdaHandler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	var userRequest UserRequest
	err := json.Unmarshal([]byte(request.Body), &userRequest)
	if err != nil {
		return common.JsonResponse(400, "Bad Request: "+err.Error(), nil)
	}

	err = common.ValidateRequest(userRequest)
	if err != nil {
		return common.JsonResponse(400, "Bad Request: "+err.Error(), nil)
	}

	user := domain.User{
		Username:       userRequest.Username,
		Name:           userRequest.Name,
		LastName:       userRequest.LastName,
		SecondLastName: userRequest.SecondLastName,
		Email:          userRequest.Email,
		DailingCode:    userRequest.DailingCode,
		Phone:          userRequest.Phone,
	}

	userRepository := repositories.NewUserRepository(dbInstance)
	err = userRepository.CreateUser(&user)
	if err != nil {
		return common.JsonResponse(500, "Internal Server Error: "+err.Error(), nil)
	}

	return common.JsonResponse(200, "User created successfully", user)
}
