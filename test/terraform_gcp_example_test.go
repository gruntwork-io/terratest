package test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/gcp"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/retry"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestTerraformGcpExample(t *testing.T) {
	t.Parallel()

	// Get the Project Id to use
	projectId := gcp.GetGoogleProjectIDFromEnvVar()

	// Create all resources in the following zone
	zone := "us-east1-b"

	// Give the example bucket a unique name so we can distinguish it from any other bucket in your GCP account
	expectedBucketName := fmt.Sprintf("terratest-gcp-example-%s", strings.ToLower(random.UniqueId()))

	// Also give the example instance a unique name
	expectedInstanceName := fmt.Sprintf("terratest-gcp-example-%s", strings.ToLower(random.UniqueId()))

	terraformOptions := &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: "../examples/terraform-gcp-example",

		// Variables to pass to our Terraform code using -var options
		Vars: map[string]interface{}{
			"zone":          zone,
			"instance_name": expectedInstanceName,
			"bucket_name":   expectedBucketName,
		},
	}

	// At the end of the test, run `terraform destroy` to clean up any resources that were created
	defer terraform.Destroy(t, terraformOptions)

	// This will run `terraform init` and `terraform apply` and fail the test if there are any errors
	terraform.InitAndApply(t, terraformOptions)

	// Run `terraform output` to get the value of some of the output variables
	bucketURL := terraform.Output(t, terraformOptions, "bucket_url")
	instanceID := terraform.Output(t, terraformOptions, "instance_id")

	// Verify that the new bucket url matches the expected url
	expectedURL := fmt.Sprintf("gs://%s", expectedBucketName)
	assert.Equal(t, expectedURL, bucketURL)

	// Verify that the Storage Bucket exists
	gcp.AssertStorageBucketExists(t, expectedBucketName)

	// Add a tag to the Compute Instance
	gcp.AddLabelsToInstance(t, projectId, zone, instanceID, map[string]string{"testing": "testing-tag-value2"})

	// Check for the labels within a retry loop as it can sometimes take a while for the
	// changes to propagate.
	maxRetries := 12
	timeBetweenRetries := 5 * time.Second
	expectedText := "testing-tag-value2"

	retry.DoWithRetry(t, fmt.Sprintf("Checking Instance %s for labels", instanceID), maxRetries, timeBetweenRetries, func() (string, error) {
		// Look up the tags for the given Instance ID
		instanceLabels, err := gcp.GetLabelsForComputeInstanceE(t, projectId, zone, instanceID)

		if err != nil {
			return "", err
		}

		testingTag, containsTestingTag := instanceLabels["testing"]
		actualText := strings.TrimSpace(testingTag)
		if containsTestingTag != true {
			return "", fmt.Errorf("Expected the tag 'testing' to exist")
		}

		if actualText != expectedText {
			return "", fmt.Errorf("Expeced GetLabelsForComputeInstanceE to return '%s' but got '%s'", expectedText, actualText)
		}

		return "", nil
	})
}
