package aws

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// InvocationTypeOption identifies the invocation mode passed to AWS Lambda
// when calling [InvokeFunctionWithParamsContextE]. See the
// InvocationType-prefixed constants for the supported values.
type InvocationTypeOption string

// Supported [InvocationTypeOption] values for Lambda invocation modes.
const (
	// InvocationTypeRequestResponse invokes the function synchronously, keeping
	// the connection open until the function returns a response or times out.
	InvocationTypeRequestResponse InvocationTypeOption = "RequestResponse"
	// InvocationTypeDryRun validates parameter values and verifies that the user
	// or role has permission to invoke the function, without actually invoking it.
	InvocationTypeDryRun InvocationTypeOption = "DryRun"
)

// Value returns the string form of the [InvocationTypeOption], defaulting to
// [InvocationTypeRequestResponse] when itype is nil. It returns an error if
// itype is set to a value that is not one of the supported invocation types.
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
	StatusCode int32
}

// InvokeFunctionContextE invokes a lambda function.
// The ctx parameter supports cancellation and timeouts.
func InvokeFunctionContextE(t testing.TestingT, ctx context.Context, region, functionName string, payload interface{}) ([]byte, error) {
	lambdaClient, err := NewLambdaClientContextE(t, ctx, region)
	if err != nil {
		return nil, err
	}

	invokeInput := &lambda.InvokeInput{
		FunctionName: &functionName,
	}

	if payload != nil {
		payloadJSON, err := json.Marshal(payload)
		if err != nil {
			return nil, err
		}

		invokeInput.Payload = payloadJSON
	}

	out, err := lambdaClient.Invoke(ctx, invokeInput)
	if err != nil {
		return nil, err
	}

	if out.FunctionError != nil {
		return out.Payload, &FunctionError{Message: *out.FunctionError, StatusCode: out.StatusCode, Payload: out.Payload}
	}

	return out.Payload, nil
}

// InvokeFunctionContext invokes a lambda function.
// This function will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func InvokeFunctionContext(t testing.TestingT, ctx context.Context, region, functionName string, payload interface{}) []byte {
	t.Helper()
	out, err := InvokeFunctionContextE(t, ctx, region, functionName, payload)
	require.NoError(t, err)

	return out
}

// InvokeFunction invokes a lambda function.
//
// Deprecated: Use [InvokeFunctionContext] instead.
func InvokeFunction(t testing.TestingT, region, functionName string, payload interface{}) []byte {
	t.Helper()
	return InvokeFunctionContext(t, context.Background(), region, functionName, payload)
}

// InvokeFunctionE invokes a lambda function.
//
// Deprecated: Use [InvokeFunctionContextE] instead.
func InvokeFunctionE(t testing.TestingT, region, functionName string, payload interface{}) ([]byte, error) {
	return InvokeFunctionContextE(t, context.Background(), region, functionName, payload)
}

// InvokeFunctionWithParamsContextE invokes a lambda function using parameters
// supplied in the LambdaOptions struct.  Returns the status code and payload
// in a LambdaOutput struct and the error.  A non-nil error will either reflect
// a problem with the parameters supplied to this function or an error returned
// by the Lambda.
// The ctx parameter supports cancellation and timeouts.
func InvokeFunctionWithParamsContextE(t testing.TestingT, ctx context.Context, region, functionName string, input *LambdaOptions) (*LambdaOutput, error) {
	lambdaClient, err := NewLambdaClientContextE(t, ctx, region)
	if err != nil {
		return nil, err
	}

	// Verify the InvocationType is one of the allowed values and report
	// an error if it's not.  By default, the InvocationType will be
	// "RequestResponse".
	invocationType, err := input.InvocationType.Value()
	if err != nil {
		return nil, err
	}

	invokeInput := &lambda.InvokeInput{
		FunctionName:   &functionName,
		InvocationType: types.InvocationType(invocationType),
	}

	if input.Payload != nil {
		payloadJSON, err := json.Marshal(input.Payload)
		if err != nil {
			return nil, err
		}

		invokeInput.Payload = payloadJSON
	}

	out, err := lambdaClient.Invoke(ctx, invokeInput)
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

// InvokeFunctionWithParamsContext invokes a lambda function using parameters
// supplied in the LambdaOptions struct and returns values in a LambdaOutput
// struct.  Checks for failure using "require".
// The ctx parameter supports cancellation and timeouts.
func InvokeFunctionWithParamsContext(t testing.TestingT, ctx context.Context, region, functionName string, input *LambdaOptions) *LambdaOutput {
	t.Helper()
	out, err := InvokeFunctionWithParamsContextE(t, ctx, region, functionName, input)
	require.NoError(t, err)

	return out
}

// InvokeFunctionWithParams invokes a lambda function using parameters
// supplied in the LambdaOptions struct and returns values in a LambdaOutput
// struct.  Checks for failure using "require".
//
// Deprecated: Use [InvokeFunctionWithParamsContext] instead.
func InvokeFunctionWithParams(t testing.TestingT, region, functionName string, input *LambdaOptions) *LambdaOutput {
	t.Helper()
	return InvokeFunctionWithParamsContext(t, context.Background(), region, functionName, input)
}

// InvokeFunctionWithParamsE invokes a lambda function using parameters
// supplied in the LambdaOptions struct.  Returns the status code and payload
// in a LambdaOutput struct and the error.  A non-nil error will either reflect
// a problem with the parameters supplied to this function or an error returned
// by the Lambda.
//
// Deprecated: Use [InvokeFunctionWithParamsContextE] instead.
func InvokeFunctionWithParamsE(t testing.TestingT, region, functionName string, input *LambdaOptions) (*LambdaOutput, error) {
	return InvokeFunctionWithParamsContextE(t, context.Background(), region, functionName, input)
}

// FunctionError is returned when an AWS Lambda invocation reports a function
// error in its response. It carries the error message, the HTTP status code,
// and the raw payload returned by the function.
type FunctionError struct {
	Message    string
	Payload    []byte
	StatusCode int32
}

func (err *FunctionError) Error() string {
	return fmt.Sprintf("%q error with status code %d invoking lambda function: %q", err.Message, err.StatusCode, err.Payload)
}

// NewLambdaClientContextE creates a new Lambda client.
// The ctx parameter supports cancellation and timeouts.
func NewLambdaClientContextE(t testing.TestingT, ctx context.Context, region string) (*lambda.Client, error) {
	sess, err := NewAuthenticatedSessionContext(ctx, region)
	if err != nil {
		return nil, err
	}

	return lambda.NewFromConfig(*sess), nil
}

// NewLambdaClientContext creates a new Lambda client.
// This function will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func NewLambdaClientContext(t testing.TestingT, ctx context.Context, region string) *lambda.Client {
	t.Helper()
	client, err := NewLambdaClientContextE(t, ctx, region)
	require.NoError(t, err)

	return client
}

// NewLambdaClient creates a new Lambda client.
//
// Deprecated: Use [NewLambdaClientContext] instead.
func NewLambdaClient(t testing.TestingT, region string) *lambda.Client {
	t.Helper()
	return NewLambdaClientContext(t, context.Background(), region)
}

// NewLambdaClientE creates a new Lambda client.
//
// Deprecated: Use [NewLambdaClientContextE] instead.
func NewLambdaClientE(t testing.TestingT, region string) (*lambda.Client, error) {
	return NewLambdaClientContextE(t, context.Background(), region)
}
