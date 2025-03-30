package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
)

type Challenge struct {
	ChallengeName   string `json:"challengeName"`
	ChallengeResult bool   `json:"challengeResult"`
}

type Event struct {
	Version       string `json:"version"`
	UserName      string `json:"userName"`
	Region        string `json:"region"`
	UserPoolId    string `json:"userPoolId"`
	TriggerSource string `json:"triggerSource"`
	Request       struct {
		Session []Challenge `json:"session"`
	} `json:"request"`
	Response struct {
		ChallengeName      string `json:"challengeName"`
		IssueTokens        bool   `json:"issueTokens"`
		FailAuthentication bool   `json:"failAuthentication"`
	} `json:"response"`
}

func defineAuthChallengeLambdaHandler(ctx context.Context, event Event) (Event, error) {
	session := event.Request.Session
	event.Version = "1"

	fmt.Printf("Intentos previos: %d", len(session))
	fmt.Printf("Session: %v", session)
	//fmt.Printf("Último intento correcto: %v", session[len(session)-1].ChallengeResult)

	if len(session) > 0 && session[len(session)-1].ChallengeResult {
		event.Response.IssueTokens = true
		event.Response.FailAuthentication = false
	} else if len(session) > 3 {
		event.Response.IssueTokens = false
		event.Response.FailAuthentication = true
	} else {
		event.Response.ChallengeName = "CUSTOM_CHALLENGE"
		event.Response.IssueTokens = false
		event.Response.FailAuthentication = false
	}

	fmt.Printf("IssueTokens: %v", event.Response.IssueTokens)
	fmt.Printf("FailAuthentication: %v", event.Response.FailAuthentication)

	return event, nil
}

func main() {
	lambda.Start(defineAuthChallengeLambdaHandler)
}
