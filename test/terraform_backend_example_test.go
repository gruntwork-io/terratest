package test_test

import (
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/aws"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/require"
)

// An example of how to test the Terraform module in examples/terraform-backend-example using Terratest.
func TestTerraformBackendExample(t *testing.T) {
	t.Parallel()

	awsRegion := aws.GetRandomRegionContext(t, t.Context(), nil, nil)
	uniqueID := random.UniqueID()

	// Create an S3 bucket where we can store state
	bucketName := "test-terraform-backend-example-" + strings.ToLower(uniqueID)
	defer cleanupS3Bucket(t, awsRegion, bucketName)

	aws.CreateS3BucketContext(t, t.Context(), awsRegion, bucketName)

	key := uniqueID + "/terraform.tfstate"
	data := "data-for-test-" + uniqueID

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

	defer terraform.DestroyContext(t, t.Context(), terraformOptions)

	terraform.InitAndApplyContext(t, t.Context(), terraformOptions)

	// Check a state file actually got stored and contains our data in it somewhere (since that data is used in an
	// output of the Terraform code)
	contents := aws.GetS3ObjectContentsContext(t, t.Context(), awsRegion, bucketName, key)
	require.Contains(t, contents, data)

	// The module doesn't really *do* anything, so we just check a dummy output here and move on
	foo := terraform.OutputRequiredContext(t, t.Context(), terraformOptions, "foo")
	require.Equal(t, data, foo)
}

func cleanupS3Bucket(t *testing.T, awsRegion string, bucketName string) {
	t.Helper()

	aws.EmptyS3BucketContext(t, t.Context(), awsRegion, bucketName)
	aws.DeleteS3BucketContext(t, t.Context(), awsRegion, bucketName)
}
