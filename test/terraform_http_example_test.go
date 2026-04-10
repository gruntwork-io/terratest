package test_test

import (
	"crypto/tls"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/aws"
	http_helper "github.com/gruntwork-io/terratest/modules/http-helper"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
)

// An example of how to test the Terraform module in examples/terraform-http-example using Terratest.
func TestTerraformHttpExample(t *testing.T) {
	t.Parallel()

	// A unique ID we can use to namespace resources so we don't clash with anything already in the AWS account or
	// tests running in parallel
	uniqueID := random.UniqueID()

	// Give this EC2 Instance and other resources in the Terraform code a name with a unique ID so it doesn't clash
	// with anything else in the AWS account.
	instanceName := "terratest-http-example-" + uniqueID

	// Specify the text the EC2 Instance will return when we make HTTP requests to it.
	instanceText := "Hello, " + uniqueID + "!"

	// Pick a random AWS region to test in. This helps ensure your code works in all regions.
	awsRegion := aws.GetRandomStableRegionContext(t, t.Context(), nil, nil)

	// Some AWS regions are missing certain instance types, so pick an available type based on the region we picked
	instanceType := aws.GetRecommendedInstanceTypeContext(t, t.Context(), awsRegion, []string{"t2.micro, t3.micro", "t2.small", "t3.small"})

	// Construct the terraform options with default retryable errors to handle the most common retryable errors in
	// terraform testing.
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: "../examples/terraform-http-example",

		// Variables to pass to our Terraform code using -var options
		Vars: map[string]interface{}{
			"aws_region":    awsRegion,
			"instance_name": instanceName,
			"instance_text": instanceText,
			"instance_type": instanceType,
		},
	})

	// At the end of the test, run `terraform destroy` to clean up any resources that were created
	defer terraform.DestroyContext(t, t.Context(), terraformOptions)

	// This will run `terraform init` and `terraform apply` and fail the test if there are any errors
	terraform.InitAndApplyContext(t, t.Context(), terraformOptions)

	// Run `terraform output` to get the value of an output variable
	instanceURL := terraform.OutputContext(t, t.Context(), terraformOptions, "instance_url")

	// Setup a TLS configuration to submit with the helper, a blank struct is acceptable
	tlsConfig := tls.Config{}

	// It can take a minute or so for the Instance to boot up, so retry a few times
	maxRetries := 30
	timeBetweenRetries := 5 * time.Second

	// Verify that we get back a 200 OK with the expected instanceText
	http_helper.HTTPGetWithRetryContext(t, t.Context(), instanceURL, &tlsConfig, 200, instanceText, maxRetries, timeBetweenRetries)
}
