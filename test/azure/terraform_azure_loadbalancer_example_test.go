//go:build azure
// +build azure

// NOTE: We use build tags to differentiate azure testing because we currently do not have azure access setup for
// CircleCI.

package test

import (
	"context"
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
	defer terraform.Destroy(t, terraformOptions)

	// Run `terraform init` and `terraform apply`. Fail the test if there are any errors.
	terraform.InitAndApply(t, terraformOptions)

	// Run `terraform output` to get the values of output variables
	resourceGroupName := terraform.Output(t, terraformOptions, "resource_group_name")
	expectedLBPublicName := terraform.Output(t, terraformOptions, "lb_public_name")
	expectedLBPrivateName := terraform.Output(t, terraformOptions, "lb_private_name")
	expectedLBNoFEConfigName := terraform.Output(t, terraformOptions, "lb_default_name")
	expectedLBPublicFeConfigName := terraform.Output(t, terraformOptions, "lb_public_fe_config_name")
	expectedLBPrivateFeConfigName := terraform.Output(t, terraformOptions, "lb_private_fe_config_static_name")
	expectedLBPrivateIP := terraform.Output(t, terraformOptions, "lb_private_ip_static")

	actualLBDoesNotExist := azure.LoadBalancerExistsContext(t, context.Background(), "negative-test", resourceGroupName, subscriptionID)
	assert.False(t, actualLBDoesNotExist)

	t.Run("LoadBalancer_Public", func(t *testing.T) {
		// Check Public Load Balancer exists.
		actualLBPublicExists := azure.LoadBalancerExistsContext(t, context.Background(), expectedLBPublicName, resourceGroupName, subscriptionID)
		assert.True(t, actualLBPublicExists)

		// Check Frontend Configuration for Load Balancer.
		actualLBPublicFeConfigNames := azure.GetLoadBalancerFrontendIPConfigNamesContext(t, context.Background(), expectedLBPublicName, resourceGroupName, subscriptionID)
		assert.Contains(t, actualLBPublicFeConfigNames, expectedLBPublicFeConfigName)

		// Check Frontend Configuration Public Address and Public IP assignment
		actualLBPublicIPAddress, actualLBPublicIPType := azure.GetIPOfLoadBalancerFrontendIPConfigContext(t, context.Background(), expectedLBPublicFeConfigName, expectedLBPublicName, resourceGroupName, subscriptionID)
		assert.NotEmpty(t, actualLBPublicIPAddress)
		assert.Equal(t, azure.PublicIP, actualLBPublicIPType)
	})

	t.Run("LoadBalancer_Private", func(t *testing.T) {
		// Check Private Load Balancer exists.
		actualLBPrivateExists := azure.LoadBalancerExistsContext(t, context.Background(), expectedLBPrivateName, resourceGroupName, subscriptionID)
		assert.True(t, actualLBPrivateExists)

		// Check Frontend Configuration for Load Balancer.
		actualLBPrivateFeConfigNames := azure.GetLoadBalancerFrontendIPConfigNamesContext(t, context.Background(), expectedLBPrivateName, resourceGroupName, subscriptionID)
		assert.Equal(t, 2, len(actualLBPrivateFeConfigNames))
		assert.Contains(t, actualLBPrivateFeConfigNames, expectedLBPrivateFeConfigName)

		// Check Frontend Configuration Private IP Type and Address.
		actualLBPrivateIPAddress, actualLBPrivateIPType := azure.GetIPOfLoadBalancerFrontendIPConfigContext(t, context.Background(), expectedLBPrivateFeConfigName, expectedLBPrivateName, resourceGroupName, subscriptionID)
		assert.NotEmpty(t, actualLBPrivateIPAddress)
		assert.Equal(t, expectedLBPrivateIP, actualLBPrivateIPAddress)
		assert.Equal(t, azure.PrivateIP, actualLBPrivateIPType)
	})

	t.Run("LoadBalancer_Default", func(t *testing.T) {
		// Check No Frontend Config Load Balancer exists.
		actualLBNoFEConfigExists := azure.LoadBalancerExistsContext(t, context.Background(), expectedLBNoFEConfigName, resourceGroupName, subscriptionID)
		assert.True(t, actualLBNoFEConfigExists)

		// Check for No Frontend Configuration for Load Balancer.
		actualLBNoFEConfigFeConfigNames := azure.GetLoadBalancerFrontendIPConfigNamesContext(t, context.Background(), expectedLBNoFEConfigName, resourceGroupName, subscriptionID)
		assert.Equal(t, 0, len(actualLBNoFEConfigFeConfigNames))
	})
}
