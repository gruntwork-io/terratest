package test

import (
	"fmt"
	"testing"

	"github.com/gruntwork-io/terratest/modules/aws"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// An example of how to test the Terraform module in examples/terraform-aws-lambda-example using Terratest.
func TestTerraformAwsLambdaExample(t *testing.T) {
	t.Parallel()

	// Give this lambda function a unique ID for a name so we can distinguish it from any other lambdas
	// in your AWS account
	functionName := fmt.Sprintf("terratest-aws-lambda-example-%s", random.UniqueId())

	// Pick a random AWS region to test in. This helps ensure your code works in all regions.
	awsRegion := aws.GetRandomStableRegion(t, nil, nil)

	// Construct the terraform options with default retryable errors to handle the most common retryable errors in
	// terraform testing.
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: "../examples/terraform-aws-lambda-example",

		// Variables to pass to our Terraform code using -var options
		Vars: map[string]interface{}{
			"function_name": functionName,
		},

		// Environment variables to set when running Terraform
		EnvVars: map[string]string{
			"AWS_DEFAULT_REGION": awsRegion,
		},
	})

	// At the end of the test, run `terraform destroy` to clean up any resources that were created
	defer terraform.Destroy(t, terraformOptions)

	// This will run `terraform init` and `terraform apply` and fail the test if there are any errors
	terraform.InitAndApply(t, terraformOptions)

	// Invoke the function, so we can test its output
	response := aws.InvokeFunction(t, awsRegion, functionName, ExampleFunctionPayload{ShouldFail: false, Echo: "hi!"})

	// This function just echos it's input as a JSON string when `ShouldFail` is `false``
	assert.Equal(t, `"hi!"`, string(response))

	// Invoke the function, this time causing it to error and capturing the error
	response, err := aws.InvokeFunctionE(t, awsRegion, functionName, ExampleFunctionPayload{ShouldFail: true, Echo: "hi!"})

	// Function-specific errors have their own special return
	functionError, ok := err.(*aws.FunctionError)
	require.True(t, ok)

	// Make sure the function-specific error comes back
	assert.Contains(t, string(functionError.Payload), "Failed to handle")

	// Get function configuration
	functionConfiguration := aws.GetLambdaFunctionConfiguration(t, awsRegion, functionName)

	// Extract some specific values from function configuration
	expectedFunctionHash := *functionConfiguration.CodeSha256
	expectedFunctionArn := *functionConfiguration.FunctionArn

	// Get outputs from terraform module
	actualFunctionHash := terraform.Output(t, terraformOptions, "lambda_function_arn")
	actualFunctionArn := terraform.Output(t, terraformOptions, "lambda_source_code_hash")

	// Make sure that the values match
	assert.Equal(t, expectedFunctionArn, actualFunctionArn)
	assert.Equal(t, expectedFunctionHash, actualFunctionHash)

	// Get configuration of non-existing lambda function
	response, err := aws.GetLambdaFunctionConfigurationE(t, awsRegion, fmt.Sprintf("fake-function-%s", random.UniqueId()))

	// Verify that the error is returned
	assert.NotNil(t, err)

}

type ExampleFunctionPayload struct {
	Echo       string
	ShouldFail bool
}
