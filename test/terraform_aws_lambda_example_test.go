//go:build aws

package test_test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/aws"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/shell"
	"github.com/gruntwork-io/terratest/modules/terraform"
	test_structure "github.com/gruntwork-io/terratest/modules/test-structure"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// An example of how to test the Terraform module in examples/terraform-aws-lambda-example using Terratest.
func TestTerraformAwsLambdaExample(t *testing.T) {
	t.Parallel()

	// Make a copy of the terraform module to a temporary directory. This allows running multiple tests in parallel
	// against the same terraform module.
	exampleFolder := test_structure.CopyTerraformFolderToTemp(t, "../", "examples/terraform-aws-lambda-example")

	err := buildLambdaBinary(t, exampleFolder)
	require.NoError(t, err)

	// Give this lambda function a unique ID for a name so we can distinguish it from any other lambdas
	// in your AWS account
	functionName := "terratest-aws-lambda-example-" + random.UniqueID()

	// Pick a random AWS region to test in. This helps ensure your code works in all regions.
	awsRegion := aws.GetRandomStableRegionContext(t, t.Context(), nil, nil)

	// Construct the terraform options with default retryable errors to handle the most common retryable errors in
	// terraform testing.
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: exampleFolder,

		// Variables to pass to our Terraform code using -var options
		Vars: map[string]interface{}{
			"function_name": functionName,
			"region":        awsRegion,
		},
	})

	// At the end of the test, run `terraform destroy` to clean up any resources that were created
	defer terraform.DestroyContext(t, t.Context(), terraformOptions)

	// This will run `terraform init` and `terraform apply` and fail the test if there are any errors
	terraform.InitAndApplyContext(t, t.Context(), terraformOptions)

	// Invoke the function, so we can test its output
	response := aws.InvokeFunctionContext(t, t.Context(), awsRegion, functionName, ExampleFunctionPayload{ShouldFail: false, Echo: "hi!"})

	// This function just echos it's input as a JSON string when `ShouldFail` is `false``
	assert.Equal(t, `"hi!"`, string(response))

	// Invoke the function, this time causing it to error and capturing the error
	_, err = aws.InvokeFunctionContextE(t, t.Context(), awsRegion, functionName, ExampleFunctionPayload{ShouldFail: true, Echo: "hi!"})

	// Function-specific errors have their own special return
	var functionError *aws.FunctionError

	require.ErrorAs(t, err, &functionError)

	// Make sure the function-specific error comes back
	assert.Contains(t, string(functionError.Payload), "failed to handle")
}

func buildLambdaBinary(t *testing.T, tempDir string) error {
	t.Helper()

	cmd := shell.Command{
		Command: "go",
		Args: []string{
			"build",
			"-o",
			tempDir + "/src/bootstrap",
			tempDir + "/src/bootstrap.go",
		},
		Env: map[string]string{
			"GOOS":        "linux",
			"GOARCH":      "amd64",
			"CGO_ENABLED": "0",
		},
	}

	_, err := shell.RunCommandContextAndGetOutputE(t, t.Context(), &cmd)

	return err
}

// Another example of how to test the Terraform module in
// examples/terraform-aws-lambda-example using Terratest, this time with
// the aws.InvokeFunctionWithParams.
func TestTerraformAwsLambdaWithParamsExample(t *testing.T) {
	t.Parallel()

	// Make a copy of the terraform module to a temporary directory. This allows running multiple tests in parallel
	// against the same terraform module.
	exampleFolder := test_structure.CopyTerraformFolderToTemp(t, "../", "examples/terraform-aws-lambda-example")

	err := buildLambdaBinary(t, exampleFolder)
	require.NoError(t, err)

	// Give this lambda function a unique ID for a name so we can distinguish it from any other lambdas
	// in your AWS account
	functionName := "terratest-aws-lambda-withparams-example-" + random.UniqueID()

	// Pick a random AWS region to test in. This helps ensure your code works in all regions.
	awsRegion := aws.GetRandomStableRegionContext(t, t.Context(), nil, nil)

	// Construct the terraform options with default retryable errors to handle the most common retryable errors in
	// terraform testing.
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: exampleFolder,

		// Variables to pass to our Terraform code using -var options
		Vars: map[string]interface{}{
			"function_name": functionName,
			"region":        awsRegion,
		},
	})

	// At the end of the test, run `terraform destroy` to clean up any resources that were created
	defer terraform.DestroyContext(t, t.Context(), terraformOptions)

	// This will run `terraform init` and `terraform apply` and fail the test if there are any errors
	terraform.InitAndApplyContext(t, t.Context(), terraformOptions)

	// Call InvokeFunctionWithParms with an InvocationType of "DryRun".
	// A "DryRun" invocation does not execute the function, so the example
	// test function will not be checking the payload.
	invocationType := aws.InvocationTypeDryRun

	input := &aws.LambdaOptions{InvocationType: &invocationType}
	out := aws.InvokeFunctionWithParamsContext(t, t.Context(), awsRegion, functionName, input)

	// With "DryRun", there's no message in the output, but there is
	// a status code which will have a value of 204 for a successful
	// invocation.
	assert.Equal(t, 204, int(out.StatusCode))

	// Invoke the function, this time causing the Lambda to error and
	// capturing the error.
	invocationType = aws.InvocationTypeRequestResponse
	input = &aws.LambdaOptions{
		InvocationType: &invocationType,
		Payload:        ExampleFunctionPayload{ShouldFail: true, Echo: "hi!"},
	}
	out, err = aws.InvokeFunctionWithParamsContextE(t, t.Context(), awsRegion, functionName, input)

	// The Lambda executed, but should have failed.
	require.Error(t, err, "Unhandled")

	// Make sure the function-specific error comes back
	assert.Contains(t, string(out.Payload), "failed to handle")

	// Call InvokeFunctionWithParamsE with a LambdaOptions struct that has
	// an unsupported InvocationType.  The function should fail.
	invocationType = "Event"
	input = &aws.LambdaOptions{
		InvocationType: &invocationType,
		Payload:        ExampleFunctionPayload{ShouldFail: false, Echo: "hi!"},
	}
	_, err = aws.InvokeFunctionWithParamsContextE(t, t.Context(), awsRegion, functionName, input)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "LambdaOptions.InvocationType, if specified, must either be \"RequestResponse\" or \"DryRun\"")
}

type ExampleFunctionPayload struct {
	Echo       string
	ShouldFail bool
}
