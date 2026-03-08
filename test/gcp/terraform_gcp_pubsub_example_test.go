//go:build gcp
// +build gcp

// NOTE: We use build tags to differentiate GCP testing for better isolation

package test

import (
	"context"
	"fmt"
	"testing"

	"github.com/gruntwork-io/terratest/modules/gcp"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	test_structure "github.com/gruntwork-io/terratest/modules/test-structure"
	"github.com/stretchr/testify/assert"
)

func TestTerraformGcpPubSubExample(t *testing.T) {
	t.Parallel()

	// Get the Project ID from the environment variable.
	projectID := gcp.GetGoogleProjectIDFromEnvVar(t)

	// Create random unique names for our Pub/Sub resources
	// so multiple tests running simultaneously don't collide.
	expectedTopicName := fmt.Sprintf("pubsub-topic-%s", random.UniqueId())
	expectedSubscriptionName := fmt.Sprintf("pubsub-sub-%s", random.UniqueId())

	exampleDir := test_structure.CopyTerraformFolderToTemp(t, "../../", "examples/terraform-gcp-pubsub-example")

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: exampleDir,

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

	// Verify the Terraform outputs match what we expected
	assert.Equal(t, expectedTopicName, actualTopicName)
	assert.Equal(t, expectedSubscriptionName, actualSubscriptionName)

	// Verify the topic and subscription exist in GCP
	gcp.AssertTopicExistsContext(t, context.Background(), projectID, actualTopicName)
	gcp.AssertSubscriptionExistsContext(t, context.Background(), projectID, actualSubscriptionName)
}
