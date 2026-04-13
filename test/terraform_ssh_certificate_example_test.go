package test_test

import (
	"crypto/rand"
	"crypto/rsa"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	stdssh "golang.org/x/crypto/ssh"

	"github.com/gruntwork-io/terratest/modules/aws"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/retry"
	"github.com/gruntwork-io/terratest/modules/ssh"
	"github.com/gruntwork-io/terratest/modules/terraform"
	test_structure "github.com/gruntwork-io/terratest/modules/test-structure"
	"github.com/stretchr/testify/require"
)

const expectedTextSSHCert = "Hello, World"

// An example of how to test the Terraform module in examples/terraform-ssh-certificate-example using Terratest. The test
// also shows an example of how to break a test down into "stages" so you can skip stages by setting environment
// variables (e.g., skip stage "teardown" by setting the environment variable "SKIP_teardown=true"), which speeds up
// iteration when running this test over and over again locally.
func TestTerraformSshCertificateExample(t *testing.T) {
	t.Parallel()

	exampleFolder := test_structure.CopyTerraformFolderToTemp(t, "../", "examples/terraform-ssh-certificate-example")

	// At the end of the test, run `terraform destroy` to clean up any resources that were created.
	defer test_structure.RunTestStage(t, "teardown", func() {
		terraformOptions := test_structure.LoadTerraformOptions(t, exampleFolder)
		terraform.DestroyContext(t, t.Context(), terraformOptions)
	})

	// Deploy the example.
	test_structure.RunTestStage(t, "setup", func() {
		terraformOptions, keyPair := configureTerraformSSHCertificateOptions(t, exampleFolder)

		// Save the options so later test stages can use them.
		test_structure.SaveTerraformOptions(t, exampleFolder, terraformOptions)
		test_structure.SaveSSHKeyPair(t, exampleFolder, keyPair)

		// This will run `terraform init` and `terraform apply` and fail the test if there are any errors.
		terraform.InitAndApplyContext(t, t.Context(), terraformOptions)
	})

	// Make sure we can SSH to the public instance directly from the public internet.
	test_structure.RunTestStage(t, "validate", func() {
		terraformOptions := test_structure.LoadTerraformOptions(t, exampleFolder)
		keypair := test_structure.LoadSSHKeyPair(t, exampleFolder)

		testSSHCertificateToPublicHost(t, terraformOptions, keypair)
	})
}

func configureTerraformSSHCertificateOptions(t *testing.T, exampleFolder string) (*terraform.Options, *ssh.KeyPair) {
	t.Helper()

	// A unique ID we can use to namespace resources so we don't clash with anything already in the AWS account or
	// tests running in parallel.
	uniqueID := random.UniqueID()

	// Give this EC2 instance and other resources in the Terraform code a name with a unique ID so it doesn't clash
	// with anything else in the AWS account.
	instanceName := "terratest-ssh-certificate-example-" + uniqueID

	// Pick a random AWS region to test in. This helps ensure your code works in all regions.
	awsRegion := aws.GetRandomStableRegionContext(t, t.Context(), nil, nil)

	// Some AWS regions are missing certain instance types, so pick an available type based on the region we picked
	instanceType := aws.GetRecommendedInstanceTypeContext(t, t.Context(), awsRegion, []string{"t2.micro, t3.micro", "t2.small", "t3.small"})

	// Create a random ssh certificate
	k, err := rsa.GenerateKey(rand.Reader, 4096)
	require.NoError(t, err)

	caSigner, err := stdssh.NewSignerFromKey(k)
	require.NoError(t, err)

	keyPair := ssh.GenerateRSAKeyPair(t, 2048)
	pub, _, _, _, err := stdssh.ParseAuthorizedKey([]byte(keyPair.PublicKey))
	require.NoError(t, err)

	cert := &stdssh.Certificate{
		Key:             pub,
		Serial:          1,
		CertType:        stdssh.UserCert,
		KeyId:           "terratest",           // identity
		ValidPrincipals: []string{"terratest"}, // allowed usernames
		ValidAfter:      uint64(time.Now().Unix()),
		ValidBefore:     uint64(time.Now().Add(365 * 24 * time.Hour).Unix()), // 1 year
		Permissions: stdssh.Permissions{
			Extensions: map[string]string{
				"permit-pty": "", // allow PTY
			},
		},
	}
	require.NoError(t, cert.SignCert(rand.Reader, caSigner))

	keyPair.PublicKey = string(stdssh.MarshalAuthorizedKey(cert))

	// Construct the terraform options with default retryable errors to handle the most common retryable errors in
	// terraform testing.
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		// The path to where our Terraform code is located.
		TerraformDir: exampleFolder,

		// Variables to pass to our Terraform code using -var options.
		Vars: map[string]interface{}{
			"aws_region":        awsRegion,
			"instance_name":     instanceName,
			"instance_type":     instanceType,
			"ssh_ca_public_key": string(stdssh.MarshalAuthorizedKey(caSigner.PublicKey())),
		},
	})

	return terraformOptions, keyPair
}

func testSSHCertificateToPublicHost(t *testing.T, terraformOptions *terraform.Options, keyPair *ssh.KeyPair) {
	t.Helper()

	// Run `terraform output` to get the value of an output variable.
	publicInstanceIP := terraform.OutputContext(t, t.Context(), terraformOptions, "public_instance_ip")

	// We're going to try to SSH to the instance IP, using the username and password that will be set up (by
	// Terraform's user_data script) in the instance.
	publicHost := ssh.Host{
		Hostname:    publicInstanceIP,
		SshUserName: "terratest",
		SshKeyPair:  keyPair,
	}

	// It can take a minute or so for the instance to boot up, so retry a few times.
	maxRetries := 30
	timeBetweenRetries := 10 * time.Second
	description := "SSH to public host " + publicInstanceIP

	// Run a simple echo command on the server.
	command := fmt.Sprintf("echo -n '%s'", expectedTextSSHCert)

	// Verify that we can SSH to the instance and run commands.
	retry.DoWithRetryContext(t, t.Context(), description, maxRetries, timeBetweenRetries, func() (string, error) {
		actualText, err := ssh.CheckSSHCommandContextE(t, t.Context(), &publicHost, command)
		if err != nil {
			return "", err
		}

		if strings.TrimSpace(actualText) != expectedTextSSHCert {
			return "", fmt.Errorf("Expected SSH command to return '%s' but got '%s'", expectedTextSSHCert, actualText)
		}

		return "", nil
	})

	// Run a command on the server that results in an error.
	command = fmt.Sprintf("echo -n '%s' && exit 1", expectedTextSSHCert)
	description = "SSH to public host " + publicInstanceIP + " with error command"

	// Verify that we can SSH to the instance, run the command which forces an error, and see the output.
	retry.DoWithRetryContext(t, t.Context(), description, maxRetries, timeBetweenRetries, func() (string, error) {
		actualText, err := ssh.CheckSSHCommandContextE(t, t.Context(), &publicHost, command)
		if err == nil {
			return "", errors.New("Expected SSH command to return an error but got none")
		}

		if strings.TrimSpace(actualText) != expectedTextSSHCert {
			return "", fmt.Errorf("Expected SSH command to return '%s' but got '%s'", expectedTextSSHCert, actualText)
		}

		return "", nil
	})
}
