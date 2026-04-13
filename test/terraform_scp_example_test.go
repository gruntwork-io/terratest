package test_test

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/aws"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/retry"
	"github.com/gruntwork-io/terratest/modules/ssh"
	"github.com/gruntwork-io/terratest/modules/terraform"
	test_structure "github.com/gruntwork-io/terratest/modules/test-structure"
	"github.com/stretchr/testify/assert"
)

func TestTerraformScpExample(t *testing.T) {
	t.Parallel()

	exampleFolder := test_structure.CopyTerraformFolderToTemp(t, "../", "examples/terraform-asg-scp-example")

	// At the end of the test, run `terraform destroy` to clean up any resources that were created
	defer test_structure.RunTestStage(t, "teardown", func() {
		terraformOptions := test_structure.LoadTerraformOptions(t, exampleFolder)
		terraform.DestroyContext(t, t.Context(), terraformOptions)

		keyPair := test_structure.LoadEc2KeyPair(t, exampleFolder)
		aws.DeleteEC2KeyPairContext(t, t.Context(), keyPair)
	})

	// Deploy the example
	test_structure.RunTestStage(t, "setup", func() {
		terraformOptions, keyPair := createTerraformOptions(t, exampleFolder)

		// Save the options and key pair so later test stages can use them
		test_structure.SaveTerraformOptions(t, exampleFolder, terraformOptions)
		test_structure.SaveEc2KeyPair(t, exampleFolder, keyPair)

		// This will run `terraform init` and `terraform apply` and fail the test if there are any errors
		terraform.InitAndApplyContext(t, t.Context(), terraformOptions)
	})

	// Make sure we can SCP a file from an EC2 instance to our local box
	test_structure.RunTestStage(t, "validate_file", func() {
		terraformOptions := test_structure.LoadTerraformOptions(t, exampleFolder)
		keyPair := test_structure.LoadEc2KeyPair(t, exampleFolder)

		testScpFromHost(t, terraformOptions, keyPair)
	})

	// Make sure we can SCP all files in a given remote dir from an EC2 instance to our local box
	test_structure.RunTestStage(t, "validate_dir", func() {
		terraformOptions := test_structure.LoadTerraformOptions(t, exampleFolder)
		keyPair := test_structure.LoadEc2KeyPair(t, exampleFolder)

		testScpDirFromHost(t, terraformOptions, keyPair)
	})

	// Make sure we can SCP all files in a given remote dir from an EC2 instance to our local box
	test_structure.RunTestStage(t, "validate_asg", func() {
		terraformOptions := test_structure.LoadTerraformOptions(t, exampleFolder)
		keyPair := test_structure.LoadEc2KeyPair(t, exampleFolder)

		testScpFromAsg(t, terraformOptions, keyPair, exampleFolder)
	})
}

func createTerraformOptions(t *testing.T, exampleFolder string) (*terraform.Options, *aws.Ec2Keypair) {
	t.Helper()

	// A unique ID we can use to namespace resources so we don't clash with anything already in the AWS account or
	// tests running in parallel
	uniqueID := random.UniqueID()

	// Give this EC2 Instance and other resources in the Terraform code a name with a unique ID so it doesn't clash
	// with anything else in the AWS account.
	instanceName := "terratest-asg-scp-example-" + uniqueID

	// Pick a random AWS region to test in. This helps ensure your code works in all regions.
	awsRegion := aws.GetRandomStableRegionContext(t, t.Context(), nil, nil)

	// Some AWS regions are missing certain instance types, so pick an available type based on the region we picked
	instanceType := aws.GetRecommendedInstanceTypeContext(t, t.Context(), awsRegion, []string{"t2.micro, t3.micro", "t2.small", "t3.small"})

	// Create an EC2 KeyPair that we can use for SSH access
	keyPairName := "terratest-asg-scp-example-" + uniqueID
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
			"key_pair_name": keyPairName,
			"instance_type": instanceType,
		},
	})

	return terraformOptions, keyPair
}

func testScpDirFromHost(t *testing.T, terraformOptions *terraform.Options, keyPair *aws.Ec2Keypair) {
	t.Helper()

	// Run `terraform output` to get the value of an output variable
	awsRegion := terraformOptions.Vars["aws_region"].(string)
	asgName := terraform.OutputContext(t, t.Context(), terraformOptions, "asg_name")
	instanceIds := aws.GetInstanceIdsForAsgContext(t, t.Context(), asgName, awsRegion)
	publicInstanceIP := aws.GetPublicIPOfEc2InstanceContext(t, t.Context(), instanceIds[0], awsRegion)

	// We're going to try to SSH to the instance IP, using the Key Pair we created earlier, and the user "ubuntu",
	// as we know the Instance is running an Ubuntu AMI that has such a user
	sshUserName := "ubuntu"
	publicHost := ssh.Host{
		Hostname:    publicInstanceIP,
		SshKeyPair:  keyPair.KeyPair,
		SshUserName: sshUserName,
	}

	_, remoteTempFilePath := writeSampleDataToInstance(t, publicInstanceIP, sshUserName, keyPair)
	remoteTempFolder := filepath.Dir(remoteTempFilePath)

	defer cleanup(t, publicInstanceIP, sshUserName, keyPair, remoteTempFolder)

	localDestDir := "/tmp/tempFolder"

	testcases := []struct {
		name          string
		options       ssh.SCPDownloadOptions
		expectedFiles int
	}{
		{
			"GrabAllFiles",
			ssh.SCPDownloadOptions{RemoteHost: publicHost, RemoteDir: remoteTempFolder, LocalDir: filepath.Join(localDestDir, random.UniqueID())},
			2,
		},
		{
			"GrabAllFilesExplicit",
			ssh.SCPDownloadOptions{RemoteHost: publicHost, RemoteDir: remoteTempFolder, LocalDir: filepath.Join(localDestDir, random.UniqueID()), FileNameFilters: []string{"*"}},
			2,
		},
		{
			"GrabFilesWithFilter",
			ssh.SCPDownloadOptions{RemoteHost: publicHost, RemoteDir: remoteTempFolder, LocalDir: filepath.Join(localDestDir, random.UniqueID()), FileNameFilters: []string{"*.baz"}},
			1,
		},
	}

	for i := range testcases {
		testCase := testcases[i]

		t.Run(testCase.name, func(t *testing.T) {
			err := ssh.SCPDirFromContextE(t, t.Context(), &testCase.options, false)
			if err != nil {
				t.Fatalf("Error copying from remote: %s", err.Error())
			}

			expectedNumFiles := testCase.expectedFiles

			fileInfos, err := os.ReadDir(testCase.options.LocalDir)
			if err != nil {
				t.Fatalf("Error reading from local dir: %s, due to: %s", testCase.options.LocalDir, err.Error())
			}

			actualNumFilesCopied := len(fileInfos)

			if len(fileInfos) != expectedNumFiles {
				t.Fatalf("Error: expected %d files to be copied. Only found %d", expectedNumFiles, actualNumFilesCopied)
			}

			// Clean up the temp file we created
			os.RemoveAll(testCase.options.LocalDir)
		})
	}
}

func testScpFromHost(t *testing.T, terraformOptions *terraform.Options, keyPair *aws.Ec2Keypair) {
	t.Helper()

	// Run `terraform output` to get the value of an output variable
	awsRegion := terraformOptions.Vars["aws_region"].(string)
	asgName := terraform.OutputContext(t, t.Context(), terraformOptions, "asg_name")
	instanceIds := aws.GetInstanceIdsForAsgContext(t, t.Context(), asgName, awsRegion)
	publicInstanceIP := aws.GetPublicIPOfEc2InstanceContext(t, t.Context(), instanceIds[0], awsRegion)

	// We're going to try to SSH to the instance IP, using the Key Pair we created earlier, and the user "ubuntu",
	// as we know the Instance is running an Ubuntu AMI that has such a user
	sshUserName := "ubuntu"
	publicHost := ssh.Host{
		Hostname:    publicInstanceIP,
		SshKeyPair:  keyPair.KeyPair,
		SshUserName: sshUserName,
	}

	randomData, remoteTempFilePath := writeSampleDataToInstance(t, publicInstanceIP, sshUserName, keyPair)
	remoteTempFolder := filepath.Base(remoteTempFilePath)

	defer cleanup(t, publicInstanceIP, sshUserName, keyPair, remoteTempFolder)

	localTempFileName := "/tmp/test.out"
	localFile, err := os.Create(localTempFileName)

	// Clean up the temp file we created
	defer os.Remove(localTempFileName)

	if err != nil {
		t.Fatalf("Error: creating local temp file: %s", err.Error())
	}

	ssh.SCPFileFromContextE(t, t.Context(), &publicHost, remoteTempFilePath, localFile, false)

	buf, err := os.ReadFile(localTempFileName)
	if err != nil {
		t.Fatalf("Error: Unable to read local file from disk: %s", err.Error())
	}

	localFileContents := string(buf)

	if !strings.Contains(localFileContents, randomData) {
		t.Fatalf("Error: unable to find %s in the local file. Local file's contents were: %s", randomData, localFileContents)
	}
}

func testScpFromAsg(t *testing.T, terraformOptions *terraform.Options, keyPair *aws.Ec2Keypair, exampleFolder string) {
	t.Helper()

	// Run `terraform output` to get the value of an output variable
	awsRegion := terraformOptions.Vars["aws_region"].(string)
	asgName := terraform.OutputContext(t, t.Context(), terraformOptions, "asg_name")
	instanceIds := aws.GetInstanceIdsForAsgContext(t, t.Context(), asgName, awsRegion)
	publicInstanceIP := aws.GetPublicIPOfEc2InstanceContext(t, t.Context(), instanceIds[0], awsRegion)

	// This is where we'll store the logs from the remote server
	localDestinationDirectory := filepath.Join(exampleFolder, "logs")
	sshUserName := "ubuntu"

	randomData, remoteTempFilePath := writeSampleDataToInstance(t, publicInstanceIP, sshUserName, keyPair)
	remoteTempFolder, remoteTempFileName := filepath.Split(remoteTempFilePath)

	defer cleanup(t, publicInstanceIP, sshUserName, keyPair, remoteTempFolder)

	// This is where we will look for the downloaded syslog
	localSyslogLocation := filepath.Join(localDestinationDirectory, publicInstanceIP, "testFolder", remoteTempFileName)

	// Create a RemoteFileSpecification for our test ASG
	// We will specify that we'd like to grab /var/log/syslog
	// and store that locally.
	spec := aws.RemoteFileSpecification{
		SshUser:  sshUserName,
		UseSudo:  true,
		KeyPair:  keyPair,
		AsgNames: []string{asgName},
		RemotePathToFileFilter: map[string][]string{
			remoteTempFolder: {remoteTempFileName},
		},
		LocalDestinationDir: localDestinationDirectory,
	}

	// Go and SCP the test file from EC2 instance
	aws.FetchFilesFromAsgsPContextE(t, t.Context(), awsRegion, &spec)

	// Clean up the temp file we created
	defer os.RemoveAll(localDestinationDirectory)

	// Read the locally copied syslog
	buf, err := os.ReadFile(localSyslogLocation)
	if err != nil {
		t.Fatalf("Error: Unable to read local file from disk: %s", err.Error())
	}

	localFileContents := string(buf)

	assert.Contains(t, localFileContents, randomData)
}

func writeSampleDataToInstance(t *testing.T, publicInstanceIP string, sshUserName string, keyPair *aws.Ec2Keypair) (string, string) {
	t.Helper()

	// We're going to try to SSH to the instance IP, using the Key Pair we created earlier, and the user "ubuntu",
	// as we know the Instance is running an Ubuntu AMI that has such a user
	publicHost := ssh.Host{
		Hostname:    publicInstanceIP,
		SshKeyPair:  keyPair.KeyPair,
		SshUserName: sshUserName,
	}

	// It can take a minute or so for the Instance to boot up, so retry a few times
	maxRetries := 30
	timeBetweenRetries := 5 * time.Second
	description := "SSH to public host " + publicInstanceIP
	remoteTempFolder := "/tmp/testFolder"
	remoteTempFilePath := filepath.Join(remoteTempFolder, "test.foo")
	remoteTempFilePath2 := filepath.Join(remoteTempFolder, "test.baz")
	randomData := random.UniqueID()

	// Verify that we can SSH to the Instance and run commands
	retry.DoWithRetryContext(t, t.Context(), description, maxRetries, timeBetweenRetries, func() (string, error) {
		_, err := ssh.CheckSSHCommandContextE(t, t.Context(), &publicHost, fmt.Sprintf("mkdir -p %s && touch %s && touch %s && echo \"%s\" >> %s", remoteTempFolder, remoteTempFilePath, remoteTempFilePath2, randomData, remoteTempFilePath))
		if err != nil {
			return "", err
		}

		return "", nil
	})

	return randomData, remoteTempFilePath
}

func cleanup(t *testing.T, publicInstanceIP string, sshUserName string, keyPair *aws.Ec2Keypair, folderToClean string) {
	t.Helper()

	publicHost := ssh.Host{
		Hostname:    publicInstanceIP,
		SshKeyPair:  keyPair.KeyPair,
		SshUserName: sshUserName,
	}

	maxRetries := 30
	timeBetweenRetries := 5 * time.Second
	description := "SSH to public host " + publicInstanceIP

	// clean up the remote folder as we want may want to run another test case
	defer retry.DoWithRetryContext(t, t.Context(), description, maxRetries, timeBetweenRetries, func() (string, error) {
		_, err := ssh.CheckSSHCommandContextE(t, t.Context(), &publicHost, "rm -rf "+folderToClean)
		if err != nil {
			return "", err
		}

		return "", nil
	})
}
