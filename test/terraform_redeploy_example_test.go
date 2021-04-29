package test

import (
	"crypto/tls"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/aws"
	http_helper "github.com/gruntwork-io/terratest/modules/http-helper"
	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/retry"
	"github.com/gruntwork-io/terratest/modules/terraform"
	test_structure "github.com/gruntwork-io/terratest/modules/test-structure"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// An example of how to test the Terraform module in examples/terraform-redeploy-example using Terratest. We deploy the
// Terraform code, check that the load balancer returns the expected response, redeploy the code, and check that the
// entire time during the redeploy, the load balancer continues returning a valid response and never returns an error
// (i.e., we validate that zero-downtime deployment works).
//
// The test is broken into "stages" so you can skip stages by setting environment variables (e.g., skip stage
// "deploy_initial" by setting the environment variable "SKIP_deploy_initial=true"), which speeds up iteration when
// running this test over and over again locally.
func TestTerraformRedeployExample(t *testing.T) {
	t.Parallel()

	// The folder where we have our Terraform code
	workingDir := "../examples/terraform-redeploy-example"

	// Pick a random AWS region to test in. This helps ensure your code works in all regions.
	test_structure.RunTestStage(t, "pick_region", func() {
		awsRegion := aws.GetRandomStableRegion(t, nil, nil)
		// Save the region, so that we reuse the same region when we skip stages
		test_structure.SaveString(t, workingDir, "region", awsRegion)
	})

	// At the end of the test, clean up all the resources we created
	defer test_structure.RunTestStage(t, "teardown", func() {
		terraformOptions := test_structure.LoadTerraformOptions(t, workingDir)
		terraform.Destroy(t, terraformOptions)
	})

	// At the end of the test, fetch the logs from each Instance. This can be useful for
	// debugging issues without having to manually SSH to the server.
	defer test_structure.RunTestStage(t, "logs", func() {
		awsRegion := test_structure.LoadString(t, workingDir, "region")
		fetchSyslogForAsg(t, awsRegion, workingDir)
		fetchFilesFromAsg(t, awsRegion, workingDir)
	})

	// Deploy the web app
	test_structure.RunTestStage(t, "deploy_initial", func() {
		awsRegion := test_structure.LoadString(t, workingDir, "region")
		initialDeploy(t, awsRegion, workingDir)
	})

	// Validate that the ASG deployed and is responding to HTTP requests
	test_structure.RunTestStage(t, "validate_initial", func() {
		awsRegion := test_structure.LoadString(t, workingDir, "region")
		validateAsgRunningWebServer(t, awsRegion, workingDir)
	})

	// Validate that we can deploy a change to the ASG with zero downtime
	test_structure.RunTestStage(t, "validate_redeploy", func() {
		validateAsgRedeploy(t, workingDir)
	})
}

// Do the initial deployment of the terraform-redeploy-example
func initialDeploy(t *testing.T, awsRegion string, workingDir string) {
	// A unique ID we can use to namespace resources so we don't clash with anything already in the AWS account or
	// tests running in parallel
	uniqueID := random.UniqueId()

	// Create a KeyPair we can use later to SSH to each Instance
	keyPair := aws.CreateAndImportEC2KeyPair(t, awsRegion, uniqueID)
	test_structure.SaveEc2KeyPair(t, workingDir, keyPair)

	// Give the ASG and other resources in the Terraform code a name with a unique ID so it doesn't clash
	// with anything else in the AWS account.
	name := fmt.Sprintf("redeploy-test-%s", uniqueID)

	// Specify the text the ASG will return when we make HTTP requests to it.
	text := fmt.Sprintf("Hello, %s!", uniqueID)

	// Some AWS regions are missing certain instance types, so pick an available type based on the region we picked
	instanceType := aws.GetRecommendedInstanceType(t, awsRegion, []string{"t2.micro", "t3.micro"})

	// Construct the terraform options with default retryable errors to handle the most common retryable errors in
	// terraform testing.
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: workingDir,

		// Variables to pass to our Terraform code using -var options
		Vars: map[string]interface{}{
			"aws_region":    awsRegion,
			"instance_name": name,
			"instance_text": text,
			"instance_type": instanceType,
			"key_pair_name": keyPair.Name,
		},
	})

	// Save the Terraform Options struct so future test stages can use it
	test_structure.SaveTerraformOptions(t, workingDir, terraformOptions)

	// This will run `terraform init` and `terraform apply` and fail the test if there are any errors
	terraform.InitAndApply(t, terraformOptions)
}

// Validate the ASG has been deployed and is working
func validateAsgRunningWebServer(t *testing.T, awsRegion string, workingDir string) {
	// Load the Terraform Options saved by the earlier deploy_terraform stage
	terraformOptions := test_structure.LoadTerraformOptions(t, workingDir)

	// Run `terraform output` to get the value of an output variable
	url := terraform.Output(t, terraformOptions, "url")
	asgName := terraform.OutputRequired(t, terraformOptions, "asg_name")

	// Setup a TLS configuration to submit with the helper, a blank struct is acceptable
	tlsConfig := tls.Config{}

	// Wait and verify the ASG is scaled to the desired capacity. It can take a few minutes for the ASG to boot up, so
	// retry a few times.
	maxRetries := 30
	timeBetweenRetries := 5 * time.Second
	aws.WaitForCapacity(t, asgName, awsRegion, maxRetries, timeBetweenRetries)
	capacityInfo := aws.GetCapacityInfoForAsg(t, asgName, awsRegion)
	assert.Equal(t, capacityInfo.DesiredCapacity, int64(3))
	assert.Equal(t, capacityInfo.CurrentCapacity, int64(3))

	// Figure out what text the ASG should return for each request
	expectedText, _ := terraformOptions.Vars["instance_text"].(string)

	// Verify that we get back a 200 OK with the expectedText
	// It can take a few minutes for the ALB to boot up, so retry a few times
	http_helper.HttpGetWithRetry(t, url, &tlsConfig, 200, expectedText, maxRetries, timeBetweenRetries)
}

// Validate we can deploy an update to the ASG with zero downtime for users accessing the ALB
func validateAsgRedeploy(t *testing.T, workingDir string) {
	// Load the Terraform Options saved by the earlier deploy_terraform stage
	terraformOptions := test_structure.LoadTerraformOptions(t, workingDir)

	// Figure out what text the ASG was returning for each request
	originalText, _ := terraformOptions.Vars["instance_text"].(string)

	// New text for the ASG to return for each request
	newText := fmt.Sprintf("%s-redeploy", originalText)
	terraformOptions.Vars["instance_text"] = newText

	// Save the updated Terraform Options struct
	test_structure.SaveTerraformOptions(t, workingDir, terraformOptions)

	// Run `terraform output` to get the value of an output variable
	url := terraform.Output(t, terraformOptions, "url")

	// Setup a TLS configuration to submit with the helper, a blank struct is acceptable
	tlsConfig := tls.Config{}

	// Check once per second that the ELB returns a proper response to make sure there is no downtime during deployment
	elbChecks := retry.DoInBackgroundUntilStopped(t, fmt.Sprintf("Check URL %s", url), 1*time.Second, func() {
		http_helper.HttpGetWithCustomValidation(t, url, &tlsConfig, func(statusCode int, body string) bool {
			return statusCode == 200 && (body == originalText || body == newText)
		})
	})

	// Redeploy the cluster
	terraform.Apply(t, terraformOptions)

	// Stop checking the ELB
	elbChecks.Done()
}

// (Deprecated) See the fetchFilesFromAsg method below for a more powerful solution.
//
// Fetch the most recent syslogs for the instances in the ASG. This is a handy way to see what happened on each
// Instance as part of your test log output, without having to re-run the test and manually SSH to the Instances.
func fetchSyslogForAsg(t *testing.T, awsRegion string, workingDir string) {
	// Load the Terraform Options saved by the earlier deploy_terraform stage
	terraformOptions := test_structure.LoadTerraformOptions(t, workingDir)

	asgName := terraform.OutputRequired(t, terraformOptions, "asg_name")
	asgLogs := aws.GetSyslogForInstancesInAsg(t, asgName, awsRegion)

	logger.Logf(t, "===== First few hundred bytes of syslog for instances in ASG %s =====\n\n", asgName)

	for instanceID, logs := range asgLogs {
		logger.Logf(t, "Most recent syslog for Instance %s:\n\n%s\n", instanceID, logs)
	}
}

// Default syslog location on Ubuntu
const syslogPathUbuntu = "/var/log/syslog"

// Default location where the User Data script generates an index.html on Ubuntu
const indexHtmlUbuntu = "/index.html"

// This size is configured in the terraform-redeploy-example itself
const asgSize = 3

func fetchFilesFromAsg(t *testing.T, awsRegion string, workingDir string) {
	// Load the Terraform Options and Key Pair saved by the earlier deploy_terraform stage
	terraformOptions := test_structure.LoadTerraformOptions(t, workingDir)
	keyPair := test_structure.LoadEc2KeyPair(t, workingDir)

	asgName := terraform.OutputRequired(t, terraformOptions, "asg_name")
	instanceIdToFilePathToContents := aws.FetchContentsOfFilesFromAsg(t, awsRegion, "ubuntu", keyPair, asgName, true, syslogPathUbuntu, indexHtmlUbuntu)

	require.Len(t, instanceIdToFilePathToContents, asgSize)

	// Check that the index.html file on each Instance contains the expected text
	expectedText := terraformOptions.Vars["instance_text"]
	for instanceID, filePathToContents := range instanceIdToFilePathToContents {
		require.Contains(t, filePathToContents, indexHtmlUbuntu)
		assert.Equal(t, expectedText, strings.TrimSpace(filePathToContents[indexHtmlUbuntu]), "Expected %s on instance %s to contain %s", indexHtmlUbuntu, instanceID, expectedText)
	}

	logger.Logf(t, "===== Full contents of syslog for instances in ASG %s =====\n\n", asgName)

	// Print out the FULL contents of syslog (unlike the deprecated GetSyslogForInstancesInAsg, which only returns the
	// first few hundred bytes)
	for instanceID, filePathToContents := range instanceIdToFilePathToContents {
		require.Contains(t, filePathToContents, syslogPathUbuntu)
		logger.Logf(t, "Full syslog for Instance %s:\n\n%s\n", instanceID, filePathToContents[syslogPathUbuntu])
	}
}
