// custom_auth/create_auth_challenge/main.go
package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(ctx context.Context, event events.CognitoEventUserPoolsCustomMessage) (events.CognitoEventUserPoolsCustomMessage, error) {
	switch event.TriggerSource {
	case "CustomMessage_SignUp":
		event.Response.EmailSubject = "Welcome to SmartGo!"
		event.Response.EmailMessage = fmt.Sprintf("Your verification code is: %s", event.Request.CodeParameter)

	case "CustomMessage_ForgotPassword":
		event.Response.EmailSubject = "Reset your password in SmartGo"
		event.Response.EmailMessage = fmt.Sprintf("Your password reset code is: %s", event.Request.CodeParameter)

	case "CustomMessage_ResendCode":
		event.Response.EmailSubject = "Resend verification code in SmartGo"
		event.Response.EmailMessage = fmt.Sprintf("Your verification code is: %s", event.Request.CodeParameter)
	}

	return event, nil
}

func main() {
	lambda.Start(handler)
}
