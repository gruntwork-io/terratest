package aws

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

type InvocationTypeOption string

const (
	InvocationTypeRequestResponse InvocationTypeOption = "RequestResponse"
	InvocationTypeDryRun                               = "DryRun"
)

func (itype *InvocationTypeOption) Value() (string, error) {
	if itype != nil {
		switch *itype {
		case
			InvocationTypeRequestResponse,
			InvocationTypeDryRun:
			return string(*itype), nil
		default:
			msg := fmt.Sprintf("LambdaOptions.InvocationType, if specified, must either be \"%s\" or \"%s\"",
				InvocationTypeRequestResponse,
				InvocationTypeDryRun)
			return "", errors.New(msg)
		}
	}
	return string(InvocationTypeRequestResponse), nil
}

// LambdaOptions contains additional parameters for InvokeFunctionWithParams().
// It contains a subset of the fields found in the lambda.InvokeInput struct.
type LambdaOptions struct {
	// InvocationType can be one of InvocationTypeOption values:
	//    * InvocationTypeRequestResponse (default) - Invoke the function
	//      synchronously.  Keep the connection open until the function
	//      returns a response or times out.
	//    * InvocationTypeDryRun - Validate parameter values and verify
	//      that the user or role has permission to invoke the function.
	InvocationType *InvocationTypeOption

	// Lambda function input; will be converted to JSON.
	Payload interface{}
}

// LambdaOutput contains the output from InvokeFunctionWithParams().  The
// fields may or may not have a value depending on the invocation type and
// whether an error occurred or not.
type LambdaOutput struct {
	// The response from the function, or an error object.
	Payload []byte

	// The HTTP status code for a successful request is in the 200 range.
	// For RequestResponse invocation type, the status code is 200.
	// For the DryRun invocation type, the status code is 204.
	StatusCode *int64
}

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

// InvokeFunctionWithParams invokes a lambda function using parameters
// supplied in the LambdaOptions struct and returns values in a LambdaOutput
// struct.  Checks for failure using "require".
func InvokeFunctionWithParams(t testing.TestingT, region, functionName string, input *LambdaOptions) *LambdaOutput {
	out, err := InvokeFunctionWithParamsE(t, region, functionName, input)
	require.NoError(t, err)
	return out
}

// InvokeFunctionWithParamsE invokes a lambda function using parameters
// supplied in the LambdaOptions struct.  Returns the status code and payload
// in a LambdaOutput struct and the error.  A non-nil error will either reflect
// a problem with the parameters supplied to this function or an error returned
// by the Lambda.
func InvokeFunctionWithParamsE(t testing.TestingT, region, functionName string, input *LambdaOptions) (*LambdaOutput, error) {
	lambdaClient, err := NewLambdaClientE(t, region)
	if err != nil {
		return nil, err
	}

	// Verify the InvocationType is one of the allowed values and report
	// an error if it's not.  By default the InvocationType will be
	// "RequestResponse".
	invocationType, err := input.InvocationType.Value()
	if err != nil {
		return nil, err
	}

	invokeInput := &lambda.InvokeInput{
		FunctionName:   &functionName,
		InvocationType: &invocationType,
	}

	if input.Payload != nil {
		payloadJson, err := json.Marshal(input.Payload)
		if err != nil {
			return nil, err
		}
		invokeInput.Payload = payloadJson
	}

	out, err := lambdaClient.Invoke(invokeInput)
	if err != nil {
		return nil, err
	}

	// As this function supports different invocation types, it must
	// then support different combinations of output other than just
	// payload.
	lambdaOutput := LambdaOutput{
		Payload:    out.Payload,
		StatusCode: out.StatusCode,
	}

	if out.FunctionError != nil {
		return &lambdaOutput, errors.New(*out.FunctionError)
	}

	return &lambdaOutput, nil
}

type FunctionError struct {
	Message    string
	StatusCode int64
	Payload    []byte
}

func (err *FunctionError) Error() string {
	return fmt.Sprintf("%q error with status code %d invoking lambda function: %q", err.Message, err.StatusCode, err.Payload)
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
