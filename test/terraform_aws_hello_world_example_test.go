//go:build aws

package test_test

import (
	"testing"
	"time"

	http_helper "github.com/gruntwork-io/terratest/modules/http-helper"

	"github.com/gruntwork-io/terratest/modules/terraform"
)

func TestTerraformAwsHelloWorldExample(t *testing.T) {
	t.Parallel()

	// website::tag::2:: Construct the terraform options with default retryable errors to handle the most common
	// retryable errors in terraform testing.
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		// website::tag::1:: The path to where our Terraform code is located
		TerraformDir: "../examples/terraform-aws-hello-world-example",
	})

	// website::tag::6:: At the end of the test, run `terraform destroy` to clean up any resources that were created.
	defer terraform.DestroyContext(t, t.Context(), terraformOptions)

	// website::tag::3:: Run `terraform init` and `terraform apply`. Fail the test if there are any errors.
	terraform.InitAndApplyContext(t, t.Context(), terraformOptions)

	// website::tag::4:: Run `terraform output` to get the IP of the instance
	publicIP := terraform.OutputContext(t, t.Context(), terraformOptions, "public_ip")

	// website::tag::5:: Make an HTTP request to the instance and make sure we get back a 200 OK with the body "Hello, World!"
	url := "http://" + publicIP + ":8080"
	http_helper.HTTPGetWithRetryContext(t, t.Context(), url, nil, 200, "Hello, World!", 30, 5*time.Second)
}
