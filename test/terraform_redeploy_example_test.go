package test

import (
	"fmt"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/aws"
	"github.com/gruntwork-io/terratest/modules/http-helper"
	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/retry"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/gruntwork-io/terratest/modules/test-structure"
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

	// Pick a random AWS region to test in. This helps ensure your code works in all regions.
	awsRegion := aws.GetRandomRegion(t, nil, nil)

	// The folder where we have our Terraform code
	workingDir := "../examples/terraform-redeploy-example"

	// At the end of the test, clean up all the resources we created
	defer test_structure.RunTestStage(t, "teardown", func() {
		terraformOptions := test_structure.LoadTerraformOptions(t, workingDir)
		terraform.Destroy(t, terraformOptions)
	})

	// At the end of the test, fetch the most recent syslog entries from each Instance. This can be useful for
	// debugging issues without having to manually SSH to the server.
	defer test_structure.RunTestStage(t, "logs", func() {
		fetchSyslogForAsg(t, awsRegion, workingDir)
	})

	// Deploy the web app
	test_structure.RunTestStage(t, "deploy_initial", func() {
		initialDeploy(t, awsRegion, workingDir)
	})

	// Validate that the ASG deployed and is responding to HTTP requests
	test_structure.RunTestStage(t, "validate_initial", func() {
		validateAsgRunningWebServer(t, workingDir)
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

	// Give the ASG and other resources in the Terraform code a name with a unique ID so it doesn't clash
	// with anything else in the AWS account.
	name := fmt.Sprintf("redeploy-test-%s", uniqueID)

	// Specify the text the ASG will return when we make HTTP requests to it.
	text := fmt.Sprintf("Hello, %s!", uniqueID)

	terraformOptions := &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: workingDir,

		// Variables to pass to our Terraform code using -var options
		Vars: map[string]interface{}{
			"aws_region":    awsRegion,
			"instance_name": name,
			"instance_text": text,
		},
	}

	// Save the Terraform Options struct so future test stages can use it
	test_structure.SaveTerraformOptions(t, workingDir, terraformOptions)

	// This will run `terraform init` and `terraform apply` and fail the test if there are any errors
	terraform.InitAndApply(t, terraformOptions)
}

// Validate the ASG has been deployed and is working
func validateAsgRunningWebServer(t *testing.T, workingDir string) {
	// Load the Terraform Options saved by the earlier deploy_terraform stage
	terraformOptions := test_structure.LoadTerraformOptions(t, workingDir)

	// Run `terraform output` to get the value of an output variable
	url := terraform.Output(t, terraformOptions, "url")

	// Figure out what text the ASG should return for each request
	expectedText, _ := terraformOptions.Vars["instance_text"].(string)

	// It can take a few minutes for the ASG and ALB to boot up, so retry a few times
	maxRetries := 30
	timeBetweenRetries := 5 * time.Second

	// Verify that we get back a 200 OK with the expectedText
	http_helper.HttpGetWithRetry(t, url, 200, expectedText, maxRetries, timeBetweenRetries)
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

	// Check once per second that the ELB returns a proper response to make sure there is no downtime during deployment
	elbChecks := retry.DoInBackgroundUntilStopped(t, fmt.Sprintf("Check URL %s", url), 1*time.Second, func() {
		http_helper.HttpGetWithCustomValidation(t, url, func(statusCode int, body string) bool {
			return statusCode == 200 && (body == originalText || body == newText)
		})
	})

	// Redeploy the cluster
	terraform.Apply(t, terraformOptions)

	// Stop checking the ELB
	elbChecks.Done()
}

// Fetch the most recent syslogs for the instances in the ASG. This is a handy way to see what happened on each
// Instance as part of your test log output, without having to re-run the test and manually SSH to the Instances.
func fetchSyslogForAsg(t *testing.T, awsRegion string, workingDir string) {
	// Load the Terraform Options saved by the earlier deploy_terraform stage
	terraformOptions := test_structure.LoadTerraformOptions(t, workingDir)

	asgName := terraform.OutputRequired(t, terraformOptions, "asg_name")
	asgLogs := aws.GetSyslogForInstancesInAsg(t, asgName, awsRegion)

	logger.Logf(t, "===== Syslog for instances in ASG %s =====\n\n", asgName)

	for instanceID, logs := range asgLogs {
		logger.Logf(t, "Most recent syslog for Instance %s:\n\n%s\n", instanceID, logs)
	}
}
