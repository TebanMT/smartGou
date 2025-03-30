// custom_auth/verify_auth_challenge_response/main.go
package main

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"
)

type Event struct {
	Version       string `json:"version"`
	UserName      string `json:"userName"`
	Region        string `json:"region"`
	UserPoolId    string `json:"userPoolId"`
	TriggerSource string `json:"triggerSource"`
	Request       struct {
		PrivateChallengeParameters map[string]string `json:"privateChallengeParameters"`
		ChallengeAnswer            string            `json:"challengeAnswer"`
	} `json:"request"`
	Response struct {
		AnswerCorrect bool `json:"answerCorrect"`
	} `json:"response"`
}

func handler(ctx context.Context, event Event) (Event, error) {
	expected := event.Request.PrivateChallengeParameters["answer"]
	actual := event.Request.ChallengeAnswer

	event.Response.AnswerCorrect = (expected == actual)
	return event, nil
}

func main() {
	lambda.Start(handler)
}
