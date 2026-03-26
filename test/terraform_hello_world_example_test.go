package test_test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestTerraformHelloWorldExample(t *testing.T) {
	t.Parallel()

	// website::tag::2:: Construct the terraform options with default retryable errors to handle the most common
	// retryable errors in terraform testing.
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		// website::tag::1:: Set the path to the Terraform code that will be tested.
		TerraformDir: "../examples/terraform-hello-world-example",
	})

	// website::tag::5:: Clean up resources with "terraform destroy" at the end of the test.
	defer terraform.DestroyContext(t, t.Context(), terraformOptions)

	// website::tag::3:: Run "terraform init" and "terraform apply". Fail the test if there are any errors.
	terraform.InitAndApplyContext(t, t.Context(), terraformOptions)

	// website::tag::4:: Run `terraform output` to get the values of output variables and check they have the expected values.
	output := terraform.OutputContext(t, t.Context(), terraformOptions, "hello_world")
	assert.Equal(t, "Hello, World!", output)
}
