package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	sharedDomain "github.com/TebanMT/smartGou/src/shared/domain"
	"github.com/aws/aws-lambda-go/events"
	"github.com/go-playground/validator"
	"github.com/gorilla/schema"
)

var ERROR_CODES = map[int]string{
	// 2xx
	200: "OK",
	201: "CREATED",
	202: "ACCEPTED",
	203: "NON_AUTHORITATIVE_INFORMATION",
	204: "NO_CONTENT",
	205: "RESET_CONTENT",
	206: "PARTIAL_CONTENT",
	// 4xx
	400: "BAD_REQUEST",
	401: "UNAUTHORIZED",
	403: "FORBIDDEN",
	404: "NOT_FOUND",
	405: "METHOD_NOT_ALLOWED",
	406: "NOT_ACCEPTABLE",
	407: "PROXY_AUTHENTICATION_REQUIRED",
	408: "REQUEST_TIMEOUT",
	410: "GONE",
	411: "LENGTH_REQUIRED",
	412: "PRECONDITION_FAILED",
	// 5xx
	500: "INTERNAL_SERVER_ERROR",
	501: "NOT_IMPLEMENTED",
	502: "BAD_GATEWAY",
	503: "SERVICE_UNAVAILABLE",
	504: "GATEWAY_TIMEOUT",
}

var decoder = schema.NewDecoder()
var validate = validator.New()
var accessControlAllowOrigin = os.Getenv("ACCESS_CONTROL_ALLOW_ORIGIN")

type Response[T any] struct {
	StatusCode int     `json:"status_code"`
	Message    string  `json:"message"`
	Data       T       `json:"data"`
	Exception  *string `json:"exception" extensions:"x-nullable=true"`
}

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

func singleToMulti(src map[string]string) map[string][]string {
	dst := make(map[string][]string, len(src))
	for k, v := range src {
		// incluso si v == "", queremos un slice
		dst[k] = []string{v}
	}
	return dst
}

func ValidateQueryStringParameters[T any](request events.APIGatewayProxyRequest, params *T) error {
	multiValueQueryStringParameters := singleToMulti(request.QueryStringParameters)
	fmt.Println("multiValueQueryStringParameters", multiValueQueryStringParameters)
	err := decoder.Decode(params, multiValueQueryStringParameters)
	if err != nil {
		fmt.Println("error Decode", err)
		return err
	}

	if err := validate.Struct(params); err != nil {
		fmt.Println("error Struct", err)
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

func JsonResponse[T any](statusCode int, message string, data T, exception string) (events.APIGatewayProxyResponse, error) {
	var exceptionPtr *string
	if exception != "" {
		exceptionPtr = &exception
	}
	if message == "" {
		message = ERROR_CODES[statusCode]
	}
	response := Response[T]{
		StatusCode: statusCode,
		Message:    message,
		Data:       data,
		Exception:  exceptionPtr,
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

func DomainErrorToHttpCode(err error) *int {
	if err == nil {
		return nil
	}
	var customError sharedDomain.CustomError
	var httpCode int
	if errors.As(err, &customError) {
		switch customError.ErrorCode {
		case sharedDomain.DataValidationCode:
			httpCode = 400
		case sharedDomain.UnexpectedErrors:
			httpCode = 500
		}
	}
	return &httpCode
}
