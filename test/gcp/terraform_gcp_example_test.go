//go:build gcp
// +build gcp

// NOTE: We use build tags to differentiate GCP testing for better isolation and parallelism when executing our tests.

package test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/gcp"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/retry"
	"github.com/gruntwork-io/terratest/modules/ssh"
	"github.com/gruntwork-io/terratest/modules/terraform"
	test_structure "github.com/gruntwork-io/terratest/modules/test-structure"
	"github.com/stretchr/testify/assert"
)

func TestTerraformGcpExample(t *testing.T) {
	t.Parallel()

	exampleDir := test_structure.CopyTerraformFolderToTemp(t, "../../", "examples/terraform-gcp-example")

	// Get the Project Id to use
	projectId := gcp.GetGoogleProjectIDFromEnvVar(t)

	// Create all resources in the following zone
	zone := "us-east1-b"

	// Give the example bucket a unique name so we can distinguish it from any other bucket in your GCP account
	expectedBucketName := fmt.Sprintf("terratest-gcp-example-%s", strings.ToLower(random.UniqueId()))

	// Also give the example instance a unique name
	expectedInstanceName := fmt.Sprintf("terratest-gcp-example-%s", strings.ToLower(random.UniqueId()))

	// website::tag::1::Configure Terraform setting path to Terraform code, bucket name, and instance name. Construct
	// the terraform options with default retryable errors to handle the most common retryable errors in terraform
	// testing.
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: exampleDir,

		// Variables to pass to our Terraform code using -var options
		Vars: map[string]interface{}{
			"gcp_project_id": projectId,
			"zone":           zone,
			"instance_name":  expectedInstanceName,
			"bucket_name":    expectedBucketName,
		},
	})

	// website::tag::5::At the end of the test, run `terraform destroy` to clean up any resources that were created
	defer terraform.Destroy(t, terraformOptions)

	// website::tag::2::This will run `terraform init` and `terraform apply` and fail the test if there are any errors
	terraform.InitAndApply(t, terraformOptions)

	// Run `terraform output` to get the value of some of the output variables
	bucketURL := terraform.Output(t, terraformOptions, "bucket_url")
	instanceName := terraform.Output(t, terraformOptions, "instance_name")

	// website::tag::3::Verify that the new bucket url matches the expected url
	expectedURL := fmt.Sprintf("gs://%s", expectedBucketName)
	assert.Equal(t, expectedURL, bucketURL)

	// Verify that the Storage Bucket exists
	gcp.AssertStorageBucketExists(t, expectedBucketName)

	// Add a tag to the Compute Instance
	instance := gcp.FetchInstance(t, projectId, instanceName)
	instance.SetLabels(t, map[string]string{"testing": "testing-tag-value2"})

	// Check for the labels within a retry loop as it can sometimes take a while for the
	// changes to propagate.
	maxRetries := 12
	timeBetweenRetries := 5 * time.Second
	expectedText := "testing-tag-value2"

	// website::tag::4::Check if the GCP instance contains a given tag.
	retry.DoWithRetry(t, fmt.Sprintf("Checking Instance %s for labels", instanceName), maxRetries, timeBetweenRetries, func() (string, error) {
		// Look up the tags for the given Instance ID
		instance := gcp.FetchInstance(t, projectId, instanceName)
		instanceLabels := instance.GetLabels(t)

		testingTag, containsTestingTag := instanceLabels["testing"]
		actualText := strings.TrimSpace(testingTag)
		if !containsTestingTag {
			return "", fmt.Errorf("Expected the tag 'testing' to exist")
		}

		if actualText != expectedText {
			return "", fmt.Errorf("Expected GetLabelsForComputeInstanceE to return '%s' but got '%s'", expectedText, actualText)
		}

		return "", nil
	})
}

// Create a Compute Instance, and attempt to SSH in and run a command.
func TestSshAccessToComputeInstance(t *testing.T) {
	t.Parallel()

	exampleDir := test_structure.CopyTerraformFolderToTemp(t, "../../", "examples/terraform-gcp-example")

	// Setup values for our Terraform apply
	projectID := gcp.GetGoogleProjectIDFromEnvVar(t)
	randomValidGcpName := gcp.RandomValidGcpName()
	zone := gcp.GetRandomZone(t, projectID, ZonesThatSupportF1Micro, nil, nil)

	terraformOptions := &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: exampleDir,

		// Variables to pass to our Terraform code using -var options
		Vars: map[string]interface{}{
			"gcp_project_id": projectID,
			"instance_name":  randomValidGcpName,
			"bucket_name":    randomValidGcpName,
			"zone":           zone,
		},
	}

	// At the end of the test, run `terraform destroy` to clean up any resources that were created
	defer terraform.Destroy(t, terraformOptions)

	// This will run `terraform init` and `terraform apply` and fail the test if there are any errors
	terraform.InitAndApply(t, terraformOptions)

	// Run `terraform output` to get the value of an output variable
	publicIp := terraform.Output(t, terraformOptions, "public_ip")

	// Attempt to SSH and execute the command
	instance := gcp.FetchInstance(t, projectID, randomValidGcpName)

	sampleText := "Hello World"
	sshUsername := "terratest"

	keyPair := ssh.GenerateRSAKeyPair(t, 2048)
	instance.AddSshKey(t, sshUsername, keyPair.PublicKey)

	host := ssh.Host{
		Hostname:    publicIp,
		SshKeyPair:  keyPair,
		SshUserName: sshUsername,
	}

	maxRetries := 20
	sleepBetweenRetries := 3 * time.Second

	retry.DoWithRetry(t, "Attempting to SSH", maxRetries, sleepBetweenRetries, func() (string, error) {
		output, err := ssh.CheckSshCommandE(t, host, fmt.Sprintf("echo '%s'", sampleText))
		if err != nil {
			return "", err
		}

		if strings.TrimSpace(sampleText) != strings.TrimSpace(output) {
			return "", fmt.Errorf("Expected: %s. Got: %s\n", sampleText, output)
		}

		return "", nil
	})
}
