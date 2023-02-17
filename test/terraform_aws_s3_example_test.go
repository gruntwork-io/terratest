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

// An example of how to test the Terraform module in examples/terraform-aws-s3-example using Terratest.
func TestTerraformAwsS3Example(t *testing.T) {
	t.Parallel()

	// Give this S3 Bucket a unique ID for a name tag so we can distinguish it from any other Buckets provisioned
	// in your AWS account
	expectedName := fmt.Sprintf("terratest-aws-s3-example-%s", strings.ToLower(random.UniqueId()))

	// Give this S3 Bucket an environment to operate as a part of for the purposes of resource tagging
	expectedEnvironment := "Automated Testing"

	// Pick a random AWS region to test in. This helps ensure your code works in all regions.
	awsRegion := aws.GetRandomStableRegion(t, nil, nil)

	// Construct the terraform options with default retryable errors to handle the most common retryable errors in
	// terraform testing.
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: "../examples/terraform-aws-s3-example",

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
