package test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/aws"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/retry"
	"github.com/gruntwork-io/terratest/modules/ssh"
	"github.com/gruntwork-io/terratest/modules/terraform"
	test_structure "github.com/gruntwork-io/terratest/modules/test-structure"
)

// An example of how to test the Terraform module in examples/terraform-ssh-password-example using Terratest. The test
// also shows an example of how to break a test down into "stages" so you can skip stages by setting environment
// variables (e.g., skip stage "teardown" by setting the environment variable "SKIP_teardown=true"), which speeds up
// iteration when running this test over and over again locally.
func TestTerraformSshPasswordExample(t *testing.T) {
	t.Parallel()

	exampleFolder := test_structure.CopyTerraformFolderToTemp(t, "../", "examples/terraform-ssh-password-example")

	// At the end of the test, run `terraform destroy` to clean up any resources that were created.
	defer test_structure.RunTestStage(t, "teardown", func() {
		terraformOptions := test_structure.LoadTerraformOptions(t, exampleFolder)
		terraform.Destroy(t, terraformOptions)
	})

	// Deploy the example.
	test_structure.RunTestStage(t, "setup", func() {
		terraformOptions := configureTerraformSshPasswordOptions(t, exampleFolder)

		// Save the options so later test stages can use them.
		test_structure.SaveTerraformOptions(t, exampleFolder, terraformOptions)

		// This will run `terraform init` and `terraform apply` and fail the test if there are any errors.
		terraform.InitAndApply(t, terraformOptions)
	})

	// Make sure we can SSH to the public instance directly from the public internet.
	test_structure.RunTestStage(t, "validate", func() {
		terraformOptions := test_structure.LoadTerraformOptions(t, exampleFolder)

		testSSHPasswordToPublicHost(t, terraformOptions)
	})
}

func configureTerraformSshPasswordOptions(t *testing.T, exampleFolder string) *terraform.Options {
	// A unique ID we can use to namespace resources so we don't clash with anything already in the AWS account or
	// tests running in parallel.
	uniqueID := random.UniqueId()

	// Give this EC2 instance and other resources in the Terraform code a name with a unique ID so it doesn't clash
	// with anything else in the AWS account.
	instanceName := fmt.Sprintf("terratest-ssh-password-example-%s", uniqueID)

	// Pick a random AWS region to test in. This helps ensure your code works in all regions.
	awsRegion := aws.GetRandomStableRegion(t, nil, nil)

	// Some AWS regions are missing certain instance types, so pick an available type based on the region we picked
	instanceType := aws.GetRecommendedInstanceType(t, awsRegion, []string{"t2.micro", "t3.micro"})

	// Create a random password that we can use for SSH access.
	password := random.UniqueId()

	// Construct the terraform options with default retryable errors to handle the most common retryable errors in
	// terraform testing.
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		// The path to where our Terraform code is located.
		TerraformDir: exampleFolder,

		// Variables to pass to our Terraform code using -var options.
		Vars: map[string]interface{}{
			"aws_region":         awsRegion,
			"instance_name":      instanceName,
			"instance_type":      instanceType,
			"terratest_password": password,
		},
	})

	return terraformOptions
}

func testSSHPasswordToPublicHost(t *testing.T, terraformOptions *terraform.Options) {
	// Run `terraform output` to get the value of an output variable.
	publicInstanceIP := terraform.Output(t, terraformOptions, "public_instance_ip")

	// We're going to try to SSH to the instance IP, using the username and password that will be set up (by
	// Terraform's user_data script) in the instance.
	publicHost := ssh.Host{
		Hostname:    publicInstanceIP,
		Password:    terraformOptions.Vars["terratest_password"].(string),
		SshUserName: "terratest",
	}

	// It can take a minute or so for the instance to boot up, so retry a few times.
	maxRetries := 30
	timeBetweenRetries := 5 * time.Second
	description := fmt.Sprintf("SSH to public host %s", publicInstanceIP)

	// Run a simple echo command on the server.
	expectedText := "Hello, World"
	command := fmt.Sprintf("echo -n '%s'", expectedText)

	// Verify that we can SSH to the instance and run commands.
	retry.DoWithRetry(t, description, maxRetries, timeBetweenRetries, func() (string, error) {
		actualText, err := ssh.CheckSshCommandE(t, publicHost, command)

		if err != nil {
			return "", err
		}

		if strings.TrimSpace(actualText) != expectedText {
			return "", fmt.Errorf("Expected SSH command to return '%s' but got '%s'", expectedText, actualText)
		}

		return "", nil
	})

	// Run a command on the server that results in an error.
	expectedText = "Hello, World"
	command = fmt.Sprintf("echo -n '%s' && exit 1", expectedText)
	description = fmt.Sprintf("SSH to public host %s with error command", publicInstanceIP)

	// Verify that we can SSH to the instance, run the command which forces an error, and see the output.
	retry.DoWithRetry(t, description, maxRetries, timeBetweenRetries, func() (string, error) {
		actualText, err := ssh.CheckSshCommandE(t, publicHost, command)

		if err == nil {
			return "", fmt.Errorf("Expected SSH command to return an error but got none")
		}

		if strings.TrimSpace(actualText) != expectedText {
			return "", fmt.Errorf("Expected SSH command to return '%s' but got '%s'", expectedText, actualText)
		}

		return "", nil
	})
}
