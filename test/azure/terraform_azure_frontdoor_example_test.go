//go:build azure
// +build azure

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

	// website::tag::3:: Run `terraform output` to get the values of output variables.
	// The example deploys a CDN Front Door (azurerm_cdn_frontdoor_*) which supersedes the
	// retired classic Front Door resource. The output names are kept stable for callers, but
	// `front_door_name` is now the CDN profile name and `front_door_endpoint_name` is the
	// CDN endpoint name (a child of the profile).
	resourceGroupName := terraform.OutputContext(t, t.Context(), terraformOptions, "resource_group_name")
	frontDoorProfileName := terraform.OutputContext(t, t.Context(), terraformOptions, "front_door_name")
	frontDoorURL := terraform.OutputContext(t, t.Context(), terraformOptions, "front_door_url")
	frontDoorEndpointName := terraform.OutputContext(t, t.Context(), terraformOptions, "front_door_endpoint_name")

	// website::tag::4:: Get CDN Front Door details and assert them against the terraform output.
	// NOTE: the value of subscriptionID can be left blank, it will be replaced by the value
	//       of the environment variable ARM_SUBSCRIPTION_ID

	profileExists := azure.CDNFrontDoorProfileExistsContext(t, t.Context(), frontDoorProfileName, resourceGroupName, subscriptionID)
	assert.True(t, profileExists)

	actualProfile := azure.GetCDNFrontDoorProfileContext(t, t.Context(), frontDoorProfileName, resourceGroupName, subscriptionID)
	assert.Equal(t, frontDoorProfileName, *actualProfile.Name)

	endpointExists := azure.CDNFrontDoorEndpointExistsContext(t, t.Context(), frontDoorEndpointName, frontDoorProfileName, resourceGroupName, subscriptionID)
	assert.True(t, endpointExists)

	actualEndpoint := azure.GetCDNFrontDoorEndpointContext(t, t.Context(), frontDoorEndpointName, frontDoorProfileName, resourceGroupName, subscriptionID)
	assert.Equal(t, frontDoorEndpointName, *actualEndpoint.Name)
	assert.Equal(t, frontDoorURL, *actualEndpoint.Properties.HostName)
}
