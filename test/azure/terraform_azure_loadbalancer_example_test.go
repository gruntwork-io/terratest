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

func TestTerraformAzureLoadBalancerExample(t *testing.T) {
	t.Parallel()

	// subscriptionID is overridden by the environment variable "ARM_SUBSCRIPTION_ID"
	subscriptionID := ""
	uniquePostfix := random.UniqueID()
	privateIPForLB02 := "10.200.2.10"

	// Configure Terraform setting up a path to Terraform code.
	terraformOptions := &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: "../../examples/azure/terraform-azure-loadbalancer-example",

		// Variables to pass to our Terraform code using -var options
		Vars: map[string]interface{}{
			"postfix":       uniquePostfix,
			"lb_private_ip": privateIPForLB02,
			// "location": "East US",
		},
	}

	// At the end of the test, run `terraform destroy` to clean up any resources that were created.
	t.Cleanup(func() {
		terraform.DestroyContext(t, t.Context(), terraformOptions)
	})

	// Run `terraform init` and `terraform apply`. Fail the test if there are any errors.
	terraform.InitAndApplyContext(t, t.Context(), terraformOptions)

	// Run `terraform output` to get the values of output variables
	resourceGroupName := terraform.OutputContext(t, t.Context(), terraformOptions, "resource_group_name")
	expectedLBPublicName := terraform.OutputContext(t, t.Context(), terraformOptions, "lb_public_name")
	expectedLBPrivateName := terraform.OutputContext(t, t.Context(), terraformOptions, "lb_private_name")
	expectedLBNoFEConfigName := terraform.OutputContext(t, t.Context(), terraformOptions, "lb_default_name")
	expectedLBPublicFeConfigName := terraform.OutputContext(t, t.Context(), terraformOptions, "lb_public_fe_config_name")
	expectedLBPrivateFeConfigName := terraform.OutputContext(t, t.Context(), terraformOptions, "lb_private_fe_config_static_name")
	expectedLBPrivateIP := terraform.OutputContext(t, t.Context(), terraformOptions, "lb_private_ip_static")

	actualLBDoesNotExist := azure.LoadBalancerExistsContext(t, t.Context(), "negative-test", resourceGroupName, subscriptionID)
	assert.False(t, actualLBDoesNotExist)

	t.Run("LoadBalancer_Public", func(t *testing.T) {
		t.Parallel()

		// Check Public Load Balancer exists.
		actualLBPublicExists := azure.LoadBalancerExistsContext(t, t.Context(), expectedLBPublicName, resourceGroupName, subscriptionID)
		assert.True(t, actualLBPublicExists)

		// Check Frontend Configuration for Load Balancer.
		actualLBPublicFeConfigNames := azure.GetLoadBalancerFrontendIPConfigNamesContext(t, t.Context(), expectedLBPublicName, resourceGroupName, subscriptionID)
		assert.Contains(t, actualLBPublicFeConfigNames, expectedLBPublicFeConfigName)

		// Check Frontend Configuration Public Address and Public IP assignment
		actualLBPublicIPAddress, actualLBPublicIPType := azure.GetIPOfLoadBalancerFrontendIPConfigContext(t, t.Context(), expectedLBPublicFeConfigName, expectedLBPublicName, resourceGroupName, subscriptionID)
		assert.NotEmpty(t, actualLBPublicIPAddress)
		assert.Equal(t, azure.PublicIP, actualLBPublicIPType)
	})

	t.Run("LoadBalancer_Private", func(t *testing.T) {
		t.Parallel()

		// Check Private Load Balancer exists.
		actualLBPrivateExists := azure.LoadBalancerExistsContext(t, t.Context(), expectedLBPrivateName, resourceGroupName, subscriptionID)
		assert.True(t, actualLBPrivateExists)

		// Check Frontend Configuration for Load Balancer.
		actualLBPrivateFeConfigNames := azure.GetLoadBalancerFrontendIPConfigNamesContext(t, t.Context(), expectedLBPrivateName, resourceGroupName, subscriptionID)
		assert.Len(t, actualLBPrivateFeConfigNames, 2)
		assert.Contains(t, actualLBPrivateFeConfigNames, expectedLBPrivateFeConfigName)

		// Check Frontend Configuration Private IP Type and Address.
		actualLBPrivateIPAddress, actualLBPrivateIPType := azure.GetIPOfLoadBalancerFrontendIPConfigContext(t, t.Context(), expectedLBPrivateFeConfigName, expectedLBPrivateName, resourceGroupName, subscriptionID)
		assert.NotEmpty(t, actualLBPrivateIPAddress)
		assert.Equal(t, expectedLBPrivateIP, actualLBPrivateIPAddress)
		assert.Equal(t, azure.PrivateIP, actualLBPrivateIPType)
	})

	t.Run("LoadBalancer_Default", func(t *testing.T) {
		t.Parallel()

		// Check No Frontend Config Load Balancer exists.
		actualLBNoFEConfigExists := azure.LoadBalancerExistsContext(t, t.Context(), expectedLBNoFEConfigName, resourceGroupName, subscriptionID)
		assert.True(t, actualLBNoFEConfigExists)

		// Check for No Frontend Configuration for Load Balancer.
		actualLBNoFEConfigFeConfigNames := azure.GetLoadBalancerFrontendIPConfigNamesContext(t, t.Context(), expectedLBNoFEConfigName, resourceGroupName, subscriptionID)
		assert.Empty(t, actualLBNoFEConfigFeConfigNames)
	})
}
