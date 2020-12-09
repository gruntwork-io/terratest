package aws

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// InvokeFunction invokes a lambda function.
func InvokeFunction(t testing.TestingT, region, functionName string, payload interface{}) []byte {
	out, err := InvokeFunctionE(t, region, functionName, payload)
	require.NoError(t, err)
	return out
}

// InvokeFunctionE invokes a lambda function.
func InvokeFunctionE(t testing.TestingT, region, functionName string, payload interface{}) ([]byte, error) {
	lambdaClient, err := NewLambdaClientE(t, region)
	if err != nil {
		return nil, err
	}

	invokeInput := &lambda.InvokeInput{
		FunctionName: &functionName,
	}

	if payload != nil {
		payloadJson, err := json.Marshal(payload)

		if err != nil {
			return nil, err
		}
		invokeInput.Payload = payloadJson
	}

	out, err := lambdaClient.Invoke(invokeInput)
	require.NoError(t, err)
	if err != nil {
		return nil, err
	}

	if out.FunctionError != nil {
		return out.Payload, &FunctionError{Message: *out.FunctionError, StatusCode: *out.StatusCode, Payload: out.Payload}
	}

	return out.Payload, nil
}

type FunctionError struct {
	Message    string
	StatusCode int64
	Payload    []byte
}

func (err *FunctionError) Error() string {
	return fmt.Sprintf("%s error invoking lambda function: %v", err.Message, err.Payload)
}

//GetLambdaFunctionConfiguration gets ARN for a given function.
func GetLambdaFunctionConfiguration(t testing.TestingT, region, functionName string) lambda.FunctionConfiguration {
	out, err := GetLambdaFunctionArnE(t, region, functionName)
	require.NoError(t, err)
	return out
}

//GetLambdaFunctionConfigurationE gets ARN for a given function.
func GetLambdaFunctionConfigurationE(t testing.TestingT, region, functionName string) (lambda.FunctionConfiguration, error) {
	lambdaClient, err := NewLambdaClient(t, region)

	input := &lambda.GetFunctionInput{
		FunctionName: &functionName
	}

	out, err = lambdaClient.GetFunction(input)
	if err != nil {
		return nil, err
	}

	return *out.Configuration, nil
}

// NewLambdaClient creates a new Lambda client.
func NewLambdaClient(t testing.TestingT, region string) *lambda.Lambda {
	client, err := NewLambdaClientE(t, region)
	require.NoError(t, err)
	return client
}

// NewLambdaClientE creates a new Lambda client.
func NewLambdaClientE(t testing.TestingT, region string) (*lambda.Lambda, error) {
	sess, err := NewAuthenticatedSession(region)
	if err != nil {
		return nil, err
	}

	return lambda.New(sess), nil
}
