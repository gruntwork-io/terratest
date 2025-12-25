package test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/aws/v2"
	"github.com/gruntwork-io/terratest/modules/random/v2"
	"github.com/gruntwork-io/terratest/modules/terraform/v2"
	"github.com/stretchr/testify/require"
)

// An example of how to test the Terraform module in examples/terraform-backend-example using Terratest.
func TestTerraformBackendExample(t *testing.T) {
	t.Parallel()

	awsRegion := aws.GetRandomRegion(t, nil, nil)
	uniqueId := random.UniqueId()

	// Create an S3 bucket where we can store state
	bucketName := fmt.Sprintf("test-terraform-backend-example-%s", strings.ToLower(uniqueId))
	defer cleanupS3Bucket(t, awsRegion, bucketName)
	aws.CreateS3Bucket(t, awsRegion, bucketName)

	key := fmt.Sprintf("%s/terraform.tfstate", uniqueId)
	data := fmt.Sprintf("data-for-test-%s", uniqueId)

	// Deploy the module, configuring it to use the S3 bucket as an S3 backend
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../examples/terraform-backend-example",
		Vars: map[string]interface{}{
			"foo": data,
		},
		BackendConfig: map[string]interface{}{
			"bucket": bucketName,
			"key":    key,
			"region": awsRegion,
		},
	})

	defer terraform.Destroy(t, terraformOptions)
	terraform.InitAndApply(t, terraformOptions)

	// Check a state file actually got stored and contains our data in it somewhere (since that data is used in an
	// output of the Terraform code)
	contents := aws.GetS3ObjectContents(t, awsRegion, bucketName, key)
	require.Contains(t, contents, data)

	// The module doesn't really *do* anything, so we just check a dummy output here and move on
	foo := terraform.OutputRequired(t, terraformOptions, "foo")
	require.Equal(t, data, foo)
}

func cleanupS3Bucket(t *testing.T, awsRegion string, bucketName string) {
	aws.EmptyS3Bucket(t, awsRegion, bucketName)
	aws.DeleteS3Bucket(t, awsRegion, bucketName)
}
