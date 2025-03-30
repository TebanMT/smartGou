package common

import (
	"encoding/json"
	"errors"
	"fmt"
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
		if fe, ok := err.(validator.ValidationErrors); ok {
			for _, fieldError := range fe {
				customMessage := getCustomMessage(fieldError)
				return errors.New(customMessage)
			}
		}
		return err
	}
	return nil
}

func getCustomMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return fe.Field() + ": This field is required"
	case "email":
		return fe.Field() + ": Must be a valid email"
	case "min":
		return fmt.Sprintf("Must be at least %s characters", fe.Param())
	case "max":
		return fmt.Sprintf("Must be at most %s characters", fe.Param())
	case "password":
		return fe.Field() + ": Password must be at least 8 characters long"
	default:
		return fe.Field() + ": Invalid format"
	}
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
