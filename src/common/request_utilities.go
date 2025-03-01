package common

import (
	"encoding/json"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/go-playground/validator"
)

type Response struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
	Data       any    `json:"data"`
}

var validate = validator.New()
var accessControlAllowOrigin = os.Getenv("ACCESS_CONTROL_ALLOW_ORIGIN")

func ValidateRequest(request any) error {
	err := validate.Struct(request)
	if err != nil {
		return err
	}
	return nil
}

func JsonResponse(statusCode int, message string, data any) (events.APIGatewayProxyResponse, error) {
	response := Response{
		StatusCode: statusCode,
		Message:    message,
		Data:       data,
	}
	responseBody, _ := json.Marshal(response)
	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Headers: map[string]string{
			"Access-Control-Allow-Origin": accessControlAllowOrigin,
		},
		Body: string(responseBody),
	}, nil
}
