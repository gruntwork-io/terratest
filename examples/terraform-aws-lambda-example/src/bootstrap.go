package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
)

type Event struct {
	Echo       string `json:"Echo"`
	ShouldFail bool   `json:"ShouldFail"`
}

// HandleRequest Fails if ShouldFail is `true`, otherwise echos the input.
func HandleRequest(ctx context.Context, evnt *Event) (string, error) {
	if evnt == nil {
		return "", errors.New("received nil event")
	}

	if evnt.ShouldFail {
		return "", fmt.Errorf("failed to handle %#v", evnt)
	}

	return evnt.Echo, nil
}

func main() {
	lambda.Start(HandleRequest)
}
