//go:build gcp
// +build gcp

// NOTE: We use build tags to differentiate GCP testing for better isolation

package test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/gcp"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
)

func TestTerraformGcpPubSubExample(t *testing.T) {
	t.Parallel()

	// Get the Project ID from the environment variable.
	projectID := gcp.GetGoogleProjectIDFromEnvVar(t)

	// Create random unique names for our Pub/Sub resources 
	// so multiple tests running simultaneously don't collide
	expectedTopicName := fmt.Sprintf("pubsub-topic-%s", random.UniqueId())
	expectedSubscriptionName := fmt.Sprintf("pubsub-sub-%s", random.UniqueId())

	// Force lowercase to ensure they match valid GCP PubSub naming requirements
	expectedTopicName = strings.ToLower(expectedTopicName)
	expectedSubscriptionName = strings.ToLower(expectedSubscriptionName)

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: "../../examples/terraform-gcp-pubsub-example",

		// Variables to pass to our Terraform code using -var options
		Vars: map[string]interface{}{
			"gcp_project_id":    projectID,
			"topic_name":        expectedTopicName,
			"subscription_name": expectedSubscriptionName,
		},
	})

	// At the end of the test, run `terraform destroy` to clean up any resources that were created
	defer terraform.Destroy(t, terraformOptions)

	// This will run `terraform init` and `terraform apply` and fail the test if there are any errors
	terraform.InitAndApply(t, terraformOptions)

	// Pull out the outputs from the Terraform configuration
	actualTopicName := terraform.Output(t, terraformOptions, "topic_name")
	actualSubscriptionName := terraform.Output(t, terraformOptions, "subscription_name")

	// Call our newly written Custom GCP Go Functions in modules/gcp/pubsub.go
	// They will automatically fail the test if the resources do not exist in the GCP environment
	gcp.AssertTopicExists(t, projectID, actualTopicName)
	gcp.AssertSubscriptionExists(t, projectID, actualSubscriptionName)
}
