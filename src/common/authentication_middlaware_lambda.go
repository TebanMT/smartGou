package common

import (
	"context"
	"strings"

	securityDomain "github.com/TebanMT/smartGou/src/modules/security/domain"
	"github.com/aws/aws-lambda-go/events"
)

type LambdaHandler func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)

func AuthenticationMiddleware(s securityDomain.TokenManager, next LambdaHandler) LambdaHandler {
	return func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		token, err := extractToken(request)
		if err != nil {
			return JsonResponse[any](401, "", nil, err.Error())
		}
		claims, err := s.ParseTokenAndValidate(ctx, token)
		if err != nil {
			return JsonResponse[any](401, "", nil, err.Error())
		}
		request.RequestContext.Authorizer = map[string]interface{}{
			"claims": claims,
		}
		return next(ctx, request)
	}
}

func extractToken(request events.APIGatewayProxyRequest) (string, error) {
	token := request.Headers["authorization"]
	if strings.HasPrefix(token, "Bearer ") {
		token = strings.TrimPrefix(token, "Bearer ")
		return token, nil
	}

	if token == "" {
		return "", securityDomain.ErrUnauthorized
	}
	return "", securityDomain.ErrInvalidToken
}
