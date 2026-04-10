//go:build azure || (azureslim && network)
// +build azure azureslim,network

// NOTE: We use build tags to differentiate azure testing because we currently do not have azure access setup for
// CircleCI.

package test_test

import (
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/azure"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestTerraformAzureNetworkExample(t *testing.T) {
	t.Parallel()

	// Create values for Terraform
	subscriptionID := ""               // subscriptionID is overridden by the environment variable "ARM_SUBSCRIPTION_ID"
	uniquePostfix := random.UniqueID() // "resource" - switch for terratest or manual terraform deployment
	expectedLocation := "eastus2"
	expectedSubnetRange := "10.0.20.0/24"
	expectedPrivateIP := "10.0.20.5"
	expectedDNSIP01 := "10.0.0.5"
	expectedDNSIP02 := "10.0.0.6"
	exectedDNSLabel := "dns-terratest-" + strings.ToLower(uniquePostfix) // only lowercase, numeric and hyphens chars allowed for DNS

	// Configure Terraform setting up a path to Terraform code.
	terraformOptions := &terraform.Options{
		// Relative path to the Terraform dir
		TerraformDir: "../../examples/azure/terraform-azure-network-example",

		// Variables to pass to our Terraform code using -var options.
		Vars: map[string]interface{}{
			"postfix":           uniquePostfix,
			"subnet_prefix":     expectedSubnetRange,
			"private_ip":        expectedPrivateIP,
			"dns_ip_01":         expectedDNSIP01,
			"dns_ip_02":         expectedDNSIP02,
			"location":          expectedLocation,
			"domain_name_label": exectedDNSLabel,
		},
	}

	// At the end of the test, run `terraform destroy` to clean up any resources that were created
	t.Cleanup(func() {
		terraform.DestroyContext(t, t.Context(), terraformOptions)
	})

	// Run `terraform init` and `terraform apply`. Fail the test if there are any errors
	terraform.InitAndApplyContext(t, t.Context(), terraformOptions)

	// Run `terraform output` to get the values of output variables
	expectedRgName := terraform.OutputContext(t, t.Context(), terraformOptions, "resource_group_name")
	expectedVNetName := terraform.OutputContext(t, t.Context(), terraformOptions, "virtual_network_name")
	expectedSubnetName := terraform.OutputContext(t, t.Context(), terraformOptions, "subnet_name")
	expectedPublicAddressName := terraform.OutputContext(t, t.Context(), terraformOptions, "public_address_name")
	expectedPrivateNicName := terraform.OutputContext(t, t.Context(), terraformOptions, "network_interface_internal")
	expectedPublicNicName := terraform.OutputContext(t, t.Context(), terraformOptions, "network_interface_external")

	// Tests are separated into subtests to differentiate integrated tests and pure resource tests

	// Integrated network resource tests
	t.Run("VirtualNetwork_Subnet", func(t *testing.T) {
		t.Parallel()

		// Check the Subnet exists in the Virtual Network Subnets with the expected Address Prefix
		actualVnetSubnets := azure.GetVirtualNetworkSubnetsContext(t, t.Context(), expectedVNetName, expectedRgName, subscriptionID)
		assert.NotNil(t, actualVnetSubnets[expectedSubnetName])
		assert.Equal(t, expectedSubnetRange, actualVnetSubnets[expectedSubnetName])
	})

	t.Run("NIC_PublicAddress", func(t *testing.T) {
		t.Parallel()

		// Check the internal network interface does NOT have a public IP
		actualPrivateIPOnly := azure.GetNetworkInterfacePublicIPsContext(t, t.Context(), expectedPrivateNicName, expectedRgName, subscriptionID)
		assert.Empty(t, actualPrivateIPOnly)

		// Check the external network interface has a public IP
		actualPublicIPs := azure.GetNetworkInterfacePublicIPsContext(t, t.Context(), expectedPublicNicName, expectedRgName, subscriptionID)
		assert.Len(t, actualPublicIPs, 1)
	})

	t.Run("Subnet_NIC", func(t *testing.T) {
		t.Parallel()

		// Check the private IP is in the subnet range
		checkPrivateIPInSubnet := azure.CheckSubnetContainsIPContext(t, t.Context(), expectedPrivateIP, expectedSubnetName, expectedVNetName, expectedRgName, subscriptionID)
		assert.True(t, checkPrivateIPInSubnet)
	})

	// Test for resource presence
	t.Run("Exists", func(t *testing.T) {
		t.Parallel()

		// Check the Virtual Network exists
		assert.True(t, azure.VirtualNetworkExistsContext(t, t.Context(), expectedVNetName, expectedRgName, subscriptionID))

		// Check the Subnet exists
		assert.True(t, azure.SubnetExistsContext(t, t.Context(), expectedSubnetName, expectedVNetName, expectedRgName, subscriptionID))

		// Check the Network Interfaces exist
		assert.True(t, azure.NetworkInterfaceExistsContext(t, t.Context(), expectedPrivateNicName, expectedRgName, subscriptionID))
		assert.True(t, azure.NetworkInterfaceExistsContext(t, t.Context(), expectedPublicNicName, expectedRgName, subscriptionID))

		// Check Network Interface that does not exist in the Resource Group
		assert.False(t, azure.NetworkInterfaceExistsContext(t, t.Context(), "negative-test", expectedRgName, subscriptionID))

		// Check Public Address exists
		assert.True(t, azure.PublicAddressExistsContext(t, t.Context(), expectedPublicAddressName, expectedRgName, subscriptionID))
	})

	// Tests for useful network properties
	t.Run("Network", func(t *testing.T) {
		t.Parallel()

		// Check the Virtual Network DNS server IPs
		actualDNSIPs := azure.GetVirtualNetworkDNSServerIPsContext(t, t.Context(), expectedVNetName, expectedRgName, subscriptionID)
		assert.Contains(t, actualDNSIPs, expectedDNSIP01)
		assert.Contains(t, actualDNSIPs, expectedDNSIP02)

		// Check the Network Interface private IP
		actualPrivateIPs := azure.GetNetworkInterfacePrivateIPsContext(t, t.Context(), expectedPrivateNicName, expectedRgName, subscriptionID)
		assert.Contains(t, actualPrivateIPs, expectedPrivateIP)

		// Check the Public Address's Public IP is allocated
		actualPublicIP := azure.GetIPOfPublicIPAddressByNameContext(t, t.Context(), expectedPublicAddressName, expectedRgName, subscriptionID)
		assert.NotEmpty(t, actualPublicIP)

		// Check DNS created for this example is reserved
		actualDNSNotAvailable := azure.CheckPublicDNSNameAvailabilityContext(t, t.Context(), expectedLocation, exectedDNSLabel, subscriptionID)
		assert.False(t, actualDNSNotAvailable)

		// Check new randomized DNS is available
		newDNSLabel := "dns-terratest-" + strings.ToLower(random.UniqueID())
		actualDNSAvailable := azure.CheckPublicDNSNameAvailabilityContext(t, t.Context(), expectedLocation, newDNSLabel, subscriptionID)
		assert.True(t, actualDNSAvailable)
	})
}
