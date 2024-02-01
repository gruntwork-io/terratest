//go:build mockaws
// +build mockaws

// NOTE: We use build tags to differentiate mockaws testing because this uses
// a local mock AWS service
// (https://docs.getmoto.org/en/latest/docs/server_mode.html)

package test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/aws"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

// An example of how to test the Terraform module in examples/terraform-aws-endpoint-example using Terratest.
func TestTerraformAwsEndpointExample(t *testing.T) {
	// Set a custom endpoint for AWS, and set the keys to dummy keys to
	// pass that check
	t.Setenv("TERRATEST_CUSTOM_AWS_ENDPOINT", "http://localhost:5000")
	t.Setenv("AWS_ACCESS_KEY_ID", "dummy")
	t.Setenv("AWS_SECRET_ACCESS_KEY", "dummy")

	// By default, the golang client will interact with s3 in a way that
	// uses subdomains to specify information about the bucket. If you have
	// issues with subdomains of "localhost" not resolving properly, set
	// this to "1" to tell terratest to use url paths rather than
	// subdomains for this information.
	// See:
	// https://docs.aws.amazon.com/AmazonS3/latest/userguide/VirtualHosting.html#path-style-access
	//
	// Note that AWS is discouraging path style endpoints in the future, so
	// the default is not to use them, see:
	// https://aws.amazon.com/blogs/aws/amazon-s3-path-deprecation-plan-the-rest-of-the-story/
	t.Setenv("TERRATEST_AWS_S3_USE_PATH_STYLE_ENDPOINT", "1")

	// Give this S3 Bucket a unique ID for a name tag so we can distinguish it from any other Buckets provisioned
	// in your AWS account
	expectedName := fmt.Sprintf("terratest-aws-endpoint-example-%s", strings.ToLower(random.UniqueId()))

	// Give this S3 Bucket an environment to operate as a part of for the purposes of resource tagging
	expectedEnvironment := "Automated Testing"

	// Pick a random AWS region to test in. This helps ensure your code works in all regions.
	awsRegion := aws.GetRandomStableRegion(t, nil, nil)

	// Construct the terraform options with default retryable errors to handle the most common retryable errors in
	// terraform testing.
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: "../examples/terraform-aws-endpoint-example",

		// Variables to pass to our Terraform code using -var options
		Vars: map[string]interface{}{
			"tag_bucket_name":        expectedName,
			"tag_bucket_environment": expectedEnvironment,
			"with_policy":            "true",
			"region":                 awsRegion,
		},
	})

	// At the end of the test, run `terraform destroy` to clean up any resources that were created
	defer terraform.Destroy(t, terraformOptions)

	// This will run `terraform init` and `terraform apply` and fail the test if there are any errors
	terraform.InitAndApply(t, terraformOptions)

	// Run `terraform output` to get the value of an output variable
	bucketID := terraform.Output(t, terraformOptions, "bucket_id")

	// Verify that our Bucket has versioning enabled
	actualStatus := aws.GetS3BucketVersioning(t, awsRegion, bucketID)
	expectedStatus := "Enabled"
	assert.Equal(t, expectedStatus, actualStatus)

	// Verify that our Bucket has a policy attached
	aws.AssertS3BucketPolicyExists(t, awsRegion, bucketID)

	// Verify that our bucket has server access logging TargetBucket set to what's expected
	loggingTargetBucket := aws.GetS3BucketLoggingTarget(t, awsRegion, bucketID)
	expectedLogsTargetBucket := fmt.Sprintf("%s-logs", bucketID)
	loggingObjectTargetPrefix := aws.GetS3BucketLoggingTargetPrefix(t, awsRegion, bucketID)
	expectedLogsTargetPrefix := "TFStateLogs/"

	assert.Equal(t, expectedLogsTargetBucket, loggingTargetBucket)
	assert.Equal(t, expectedLogsTargetPrefix, loggingObjectTargetPrefix)
}
