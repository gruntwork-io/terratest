package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
)

type Event struct {
	ShouldFail bool   `json:"ShouldFail"`
	Echo       string `json:"Echo"`
}

// Fails if ShouldFail is `true`, otherwise echos the input.
func HandleRequest(ctx context.Context, evnt Event) (string, error) {
	if evnt.ShouldFail {
		return "", fmt.Errorf("Failed to handle %#v", evnt)
	}
	return evnt.Echo, nil
}

func main() {
	lambda.Start(HandleRequest)
}
