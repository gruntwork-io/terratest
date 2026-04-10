//go:build azure
// +build azure

// NOTE: We use build tags to differentiate azure testing because we currently do not have azure access setup for
// CircleCI.

package test_test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/azure"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestTerraformAzureFrontDoorExample(t *testing.T) {
	t.Parallel()

	subscriptionID := ""
	uniquePostfix := random.UniqueID()

	// website::tag::1:: Configure Terraform setting up a path to Terraform code.
	terraformOptions := &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: "../../examples/azure/terraform-azure-frontdoor-example",
		Vars: map[string]interface{}{
			"postfix": uniquePostfix,
		},
	}

	// website::tag::4:: At the end of the test, run `terraform destroy` to clean up any resources that were created
	defer terraform.DestroyContext(t, t.Context(), terraformOptions)

	// website::tag::2:: Run `terraform init` and `terraform apply`. Fail the test if there are any errors.
	terraform.InitAndApplyContext(t, t.Context(), terraformOptions)

	// website::tag::3:: Run `terraform output` to get the values of output variables
	resourceGroupName := terraform.OutputContext(t, t.Context(), terraformOptions, "resource_group_name")
	frontDoorName := terraform.OutputContext(t, t.Context(), terraformOptions, "front_door_name")
	frontDoorURL := terraform.OutputContext(t, t.Context(), terraformOptions, "front_door_url")
	frontendEndpointName := terraform.OutputContext(t, t.Context(), terraformOptions, "front_door_endpoint_name")

	// website::tag::4:: Get FrontDoor details and assert them against the terraform output
	// NOTE: the value of subscriptionID can be left blank, it will be replaced by the value
	//       of the environment variable ARM_SUBSCRIPTION_ID

	frontDoorExists := azure.FrontDoorExistsContext(t, t.Context(), frontDoorName, resourceGroupName, subscriptionID)
	assert.True(t, frontDoorExists)

	actualFrontDoorInstance := azure.GetFrontDoorContext(t, t.Context(), frontDoorName, resourceGroupName, subscriptionID)
	assert.Equal(t, frontDoorName, *actualFrontDoorInstance.Name)

	endpointExists := azure.FrontDoorFrontendEndpointExistsContext(t, t.Context(), frontendEndpointName, frontDoorName, resourceGroupName, subscriptionID)
	assert.True(t, endpointExists)

	actualFrontDoorEndpoint := azure.GetFrontDoorFrontendEndpointContext(t, t.Context(), frontendEndpointName, frontDoorName, resourceGroupName, subscriptionID)
	assert.Equal(t, frontDoorURL, *actualFrontDoorEndpoint.Properties.HostName)
}
