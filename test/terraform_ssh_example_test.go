package test_test

import (
	"errors"
	"fmt"
	"os"
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

const expectedTextSSH = "Hello, World"

// An example of how to test the Terraform module in examples/terraform-ssh-example using Terratest. The test also
// shows an example of how to break a test down into "stages" so you can skip stages by setting environment variables
// (e.g., skip stage "teardown" by setting the environment variable "SKIP_teardown=true"), which speeds up iteration
// when running this test over and over again locally.
func TestTerraformSshExample(t *testing.T) {
	t.Parallel()

	exampleFolder := test_structure.CopyTerraformFolderToTemp(t, "../", "examples/terraform-ssh-example")

	// At the end of the test, run `terraform destroy` to clean up any resources that were created
	defer test_structure.RunTestStage(t, "teardown", func() {
		terraformOptions := test_structure.LoadTerraformOptions(t, exampleFolder)
		terraform.DestroyContext(t, t.Context(), terraformOptions)

		keyPair := test_structure.LoadEc2KeyPair(t, exampleFolder)
		aws.DeleteEC2KeyPairContext(t, t.Context(), keyPair)
	})

	// Deploy the example
	test_structure.RunTestStage(t, "setup", func() {
		terraformOptions, keyPair := configureTerraformOptions(t, exampleFolder)

		// Save the options and key pair so later test stages can use them
		test_structure.SaveTerraformOptions(t, exampleFolder, terraformOptions)
		test_structure.SaveEc2KeyPair(t, exampleFolder, keyPair)

		// This will run `terraform init` and `terraform apply` and fail the test if there are any errors
		terraform.InitAndApplyContext(t, t.Context(), terraformOptions)
	})

	// Make sure we can SSH to the public Instance directly from the public Internet and the private Instance by using
	// the public Instance as a jump host
	test_structure.RunTestStage(t, "validate", func() {
		terraformOptions := test_structure.LoadTerraformOptions(t, exampleFolder)
		keyPair := test_structure.LoadEc2KeyPair(t, exampleFolder)

		testSSHToPublicHost(t, terraformOptions, keyPair)
		testSSHToPrivateHost(t, terraformOptions, keyPair)
		testSSHAgentToPublicHost(t, terraformOptions, keyPair)
		testSSHAgentToPrivateHost(t, terraformOptions, keyPair)
		testSCPToPublicHost(t, terraformOptions, keyPair)
	})
}

func configureTerraformOptions(t *testing.T, exampleFolder string) (*terraform.Options, *aws.Ec2Keypair) {
	t.Helper()

	// A unique ID we can use to namespace resources so we don't clash with anything already in the AWS account or
	// tests running in parallel
	uniqueID := random.UniqueID()

	// Give this EC2 Instance and other resources in the Terraform code a name with a unique ID so it doesn't clash
	// with anything else in the AWS account.
	instanceName := "terratest-ssh-example-" + uniqueID

	// Pick a random AWS region to test in. This helps ensure your code works in all regions.
	awsRegion := aws.GetRandomStableRegionContext(t, t.Context(), nil, nil)

	// Some AWS regions are missing certain instance types, so pick an available type based on the region we picked
	instanceType := aws.GetRecommendedInstanceTypeContext(t, t.Context(), awsRegion, []string{"t2.micro, t3.micro", "t2.small", "t3.small"})

	// Create an EC2 KeyPair that we can use for SSH access
	keyPairName := "terratest-ssh-example-" + uniqueID
	keyPair := aws.CreateAndImportEC2KeyPairContext(t, t.Context(), awsRegion, keyPairName)

	// Construct the terraform options with default retryable errors to handle the most common retryable errors in
	// terraform testing.
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: exampleFolder,

		// Variables to pass to our Terraform code using -var options
		Vars: map[string]interface{}{
			"aws_region":    awsRegion,
			"instance_name": instanceName,
			"instance_type": instanceType,
			"key_pair_name": keyPairName,
		},
	})

	return terraformOptions, keyPair
}

func testSSHToPublicHost(t *testing.T, terraformOptions *terraform.Options, keyPair *aws.Ec2Keypair) {
	t.Helper()

	// Run `terraform output` to get the value of an output variable
	publicInstanceIP := terraform.OutputContext(t, t.Context(), terraformOptions, "public_instance_ip")

	// We're going to try to SSH to the instance IP, using the Key Pair we created earlier, and the user "ubuntu",
	// as we know the Instance is running an Ubuntu AMI that has such a user
	publicHost := ssh.Host{
		Hostname:    publicInstanceIP,
		SshKeyPair:  keyPair.KeyPair,
		SshUserName: "ubuntu",
	}

	// It can take a minute or so for the Instance to boot up, so retry a few times
	maxRetries := 30
	timeBetweenRetries := 5 * time.Second
	description := "SSH to public host " + publicInstanceIP

	// Run a simple echo command on the server
	command := fmt.Sprintf("echo -n '%s'", expectedTextSSH)

	// Verify that we can SSH to the Instance and run commands
	retry.DoWithRetryContext(t, t.Context(), description, maxRetries, timeBetweenRetries, func() (string, error) {
		actualText, err := ssh.CheckSSHCommandContextE(t, t.Context(), &publicHost, command)
		if err != nil {
			return "", err
		}

		if strings.TrimSpace(actualText) != expectedTextSSH {
			return "", fmt.Errorf("Expected SSH command to return '%s' but got '%s'", expectedTextSSH, actualText)
		}

		return "", nil
	})

	// Run a command on the server that results in an error,
	command = fmt.Sprintf("echo -n '%s' && exit 1", expectedTextSSH)
	description = "SSH to public host " + publicInstanceIP + " with error command"

	// Verify that we can SSH to the Instance, run the command and see the output
	retry.DoWithRetryContext(t, t.Context(), description, maxRetries, timeBetweenRetries, func() (string, error) {
		actualText, err := ssh.CheckSSHCommandContextE(t, t.Context(), &publicHost, command)
		if err == nil {
			return "", errors.New("Expected SSH command to return an error but got none")
		}

		if strings.TrimSpace(actualText) != expectedTextSSH {
			return "", fmt.Errorf("Expected SSH command to return '%s' but got '%s'", expectedTextSSH, actualText)
		}

		return "", nil
	})
}

func testSSHToPrivateHost(t *testing.T, terraformOptions *terraform.Options, keyPair *aws.Ec2Keypair) {
	t.Helper()

	// Run `terraform output` to get the value of an output variable
	publicInstanceIP := terraform.OutputContext(t, t.Context(), terraformOptions, "public_instance_ip")

	// Get IP of private instance from AWS helper function instead of Terraform output
	privateInstanceID := terraform.OutputContext(t, t.Context(), terraformOptions, "private_instance_id")
	deployedAWSRegion := terraformOptions.Vars["aws_region"].(string)
	privateInstanceIP := aws.GetPrivateIPOfEc2InstanceContext(t, t.Context(), privateInstanceID, deployedAWSRegion)

	sshToPrivateHost(t, publicInstanceIP, privateInstanceIP, keyPair)
}

func sshToPrivateHost(t *testing.T, publicInstanceIP string, privateInstanceIP string, keyPair *aws.Ec2Keypair) {
	t.Helper()

	// We're going to try to SSH to the private instance using the public instance as a jump host. For both instances,
	// we are using the Key Pair we created earlier, and the user "ubuntu", as we know the Instances are running an
	// Ubuntu AMI that has such a user
	publicHost := ssh.Host{
		Hostname:    publicInstanceIP,
		SshKeyPair:  keyPair.KeyPair,
		SshUserName: "ubuntu",
	}
	privateHost := ssh.Host{
		Hostname:    privateInstanceIP,
		SshKeyPair:  keyPair.KeyPair,
		SshUserName: "ubuntu",
	}

	// It can take a minute or so for the Instance to boot up, so retry a few times
	maxRetries := 30
	timeBetweenRetries := 5 * time.Second
	description := fmt.Sprintf("SSH to private host %s via public host %s", privateInstanceIP, publicInstanceIP)

	// Run a simple echo command on the server
	command := fmt.Sprintf("echo -n '%s'", expectedTextSSH)

	// Verify that we can SSH to the Instance and run commands
	retry.DoWithRetryContext(t, t.Context(), description, maxRetries, timeBetweenRetries, func() (string, error) {
		actualText, err := ssh.CheckPrivateSSHConnectionContextE(t, t.Context(), &publicHost, &privateHost, command)
		if err != nil {
			return "", err
		}

		if strings.TrimSpace(actualText) != expectedTextSSH {
			return "", fmt.Errorf("Expected SSH command to return '%s' but got '%s'", expectedTextSSH, actualText)
		}

		return "", nil
	})
}

func testSCPToPublicHost(t *testing.T, terraformOptions *terraform.Options, keyPair *aws.Ec2Keypair) {
	t.Helper()

	// Run `terraform output` to get the value of an output variable
	publicInstanceIP := terraform.OutputContext(t, t.Context(), terraformOptions, "public_instance_ip")

	// We're going to try to SSH to the instance IP, using the Key Pair we created earlier, and the user "ubuntu",
	// as we know the Instance is running an Ubuntu AMI that has such a user
	publicHost := ssh.Host{
		Hostname:    publicInstanceIP,
		SshKeyPair:  keyPair.KeyPair,
		SshUserName: "ubuntu",
	}

	// It can take a minute or so for the Instance to boot up, so retry a few times
	maxRetries := 10
	timeBetweenRetries := 1 * time.Second
	description := "SCP file to public host " + publicInstanceIP

	// Verify that we can SSH to the Instance and run commands
	retry.DoWithRetryContext(t, t.Context(), description, maxRetries, timeBetweenRetries, func() (string, error) {
		err := ssh.SCPFileToContextE(t, t.Context(), &publicHost, os.FileMode(0644), "/tmp/test.txt", expectedTextSSH)
		if err != nil {
			return "", err
		}

		actualText, err := ssh.FetchContentsOfFileContextE(t, t.Context(), &publicHost, false, "/tmp/test.txt")
		if err != nil {
			return "", err
		}

		if strings.TrimSpace(actualText) != expectedTextSSH {
			return "", fmt.Errorf("Expected SSH command to return '%s' but got '%s'", expectedTextSSH, actualText)
		}

		return "", nil
	})
}

func testSSHAgentToPublicHost(t *testing.T, terraformOptions *terraform.Options, keyPair *aws.Ec2Keypair) {
	t.Helper()

	// Run `terraform output` to get the value of an output variable
	publicInstanceIP := terraform.OutputContext(t, t.Context(), terraformOptions, "public_instance_ip")

	// start the ssh agent
	sshAgent := ssh.SSHAgentWithKeyPair(t, t.Context(), keyPair.KeyPair)
	defer sshAgent.Stop()

	// We're going to try to SSH to the instance IP, using the Key Pair we created earlier. Instead of
	// directly using the SSH key in the SSH connection, we're going to rely on an existing SSH agent that we
	// programatically emulate within this test. We're going to use the user "ubuntu" as we know the Instance
	// is running an Ubuntu AMI that has such a user
	publicHost := ssh.Host{
		Hostname:         publicInstanceIP,
		SshUserName:      "ubuntu",
		OverrideSshAgent: sshAgent,
	}

	// It can take a minute or so for the Instance to boot up, so retry a few times
	maxRetries := 30
	timeBetweenRetries := 5 * time.Second
	description := "SSH with Agent to public host " + publicInstanceIP

	// Run a simple echo command on the server
	command := fmt.Sprintf("echo -n '%s'", expectedTextSSH)

	// Verify that we can SSH to the Instance and run commands
	retry.DoWithRetryContext(t, t.Context(), description, maxRetries, timeBetweenRetries, func() (string, error) {
		actualText, err := ssh.CheckSSHCommandContextE(t, t.Context(), &publicHost, command)
		if err != nil {
			return "", err
		}

		if strings.TrimSpace(actualText) != expectedTextSSH {
			return "", fmt.Errorf("Expected SSH command to return '%s' but got '%s'", expectedTextSSH, actualText)
		}

		return "", nil
	})
}

func testSSHAgentToPrivateHost(t *testing.T, terraformOptions *terraform.Options, keyPair *aws.Ec2Keypair) {
	t.Helper()

	// Run `terraform output` to get the value of an output variable
	publicInstanceIP := terraform.OutputContext(t, t.Context(), terraformOptions, "public_instance_ip")
	privateInstanceIP := terraform.OutputContext(t, t.Context(), terraformOptions, "private_instance_ip")

	// start the ssh agent
	sshAgent := ssh.SSHAgentWithKeyPair(t, t.Context(), keyPair.KeyPair)
	defer sshAgent.Stop()

	// We're going to try to SSH to the private instance using the public instance as a jump host. Instead of
	// directly using the SSH key in the SSH connection, we're going to rely on an existing SSH agent that we
	// programatically emulate within this test. For both instances, we are using the Key Pair we created earlier,
	// and the user "ubuntu", as we know the Instances are running an Ubuntu AMI that has such a user
	publicHost := ssh.Host{
		Hostname:         publicInstanceIP,
		SshUserName:      "ubuntu",
		OverrideSshAgent: sshAgent,
	}
	privateHost := ssh.Host{
		Hostname:         privateInstanceIP,
		SshUserName:      "ubuntu",
		OverrideSshAgent: sshAgent,
	}

	// It can take a minute or so for the Instance to boot up, so retry a few times
	maxRetries := 30
	timeBetweenRetries := 5 * time.Second
	description := fmt.Sprintf("SSH with Agent to private host %s via public host %s", privateInstanceIP, publicInstanceIP)

	// Run a simple echo command on the server
	command := fmt.Sprintf("echo -n '%s'", expectedTextSSH)

	// Verify that we can SSH to the Instance and run commands
	retry.DoWithRetryContext(t, t.Context(), description, maxRetries, timeBetweenRetries, func() (string, error) {
		actualText, err := ssh.CheckPrivateSSHConnectionContextE(t, t.Context(), &publicHost, &privateHost, command)
		if err != nil {
			return "", err
		}

		if strings.TrimSpace(actualText) != expectedTextSSH {
			return "", fmt.Errorf("Expected SSH command to return '%s' but got '%s'", expectedTextSSH, actualText)
		}

		return "", nil
	})
}
