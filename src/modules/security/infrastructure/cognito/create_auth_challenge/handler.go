// custom_auth/create_auth_challenge/main.go
package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

type Event struct {
	Version       string `json:"version"`
	Region        string `json:"region"`
	UserPoolId    string `json:"userPoolId"`
	UserName      string `json:"userName"`
	TriggerSource string `json:"triggerSource"`
	Request       struct {
		UserAttributes map[string]string `json:"userAttributes"`
	} `json:"request"`
	Response struct {
		PublicChallengeParameters  map[string]string `json:"publicChallengeParameters"`
		PrivateChallengeParameters map[string]string `json:"privateChallengeParameters"`
		ChallengeMetadata          string            `json:"challengeMetadata"`
	} `json:"response"`
}

func handler(ctx context.Context, event Event) (Event, error) {
	rand.Seed(time.Now().UnixNano())
	otp := fmt.Sprintf("%06d", rand.Intn(1000000))
	fmt.Println("otp", otp)
	phone := event.Request.UserAttributes["phone_number"]

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return event, fmt.Errorf("failed to load AWS config: %w", err)
	}

	client := sns.NewFromConfig(cfg)
	_, err = client.Publish(ctx, &sns.PublishInput{
		Message:     awsString(fmt.Sprintf("Bienvenido a SmartGou, tu código de verificación es: %s", otp)),
		PhoneNumber: awsString(phone),
	})
	if err != nil {
		return event, fmt.Errorf("failed to send SMS: %w", err)
	}

	event.Response.PrivateChallengeParameters = map[string]string{"answer": otp}
	event.Response.PublicChallengeParameters = map[string]string{"message": "OTP enviado"}
	event.Response.ChallengeMetadata = "CUSTOM_CHALLENGE"

	return event, nil
}

func awsString(s string) *string {
	return &s
}

func main() {
	lambda.Start(handler)
}
