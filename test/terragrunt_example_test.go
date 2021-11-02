package test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestTerragruntExample(t *testing.T) {
	// website::tag::3:: Construct the terraform options with default retryable errors to handle the most common retryable errors in terraform testing.
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		// website::tag::1:: Set the path to the Terragrunt module that will be tested.
		TerraformDir: "../examples/terragrunt-example",
		// website::tag::2:: Set the terraform binary path to terragrunt so that terratest uses terragrunt instead of terraform. You must ensure that you have terragrunt downloaded and available in your PATH.
		TerraformBinary: "terragrunt",
	})

	// website::tag::6:: Clean up resources with "terragrunt destroy" at the end of the test.
	defer terraform.Destroy(t, terraformOptions)

	// website::tag::4:: Run "terragrunt apply". Under the hood, terragrunt will run "terraform init" and "terraform apply". Fail the test if there are any errors.
	terraform.Apply(t, terraformOptions)

	// website::tag::5:: Run `terraform output` to get the values of output variables and check they have the expected values.
	output := terraform.Output(t, terraformOptions, "output")
	assert.Equal(t, "one input another input", output)
}
