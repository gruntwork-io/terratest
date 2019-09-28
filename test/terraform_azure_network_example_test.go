package test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/mathieubuisson/terratest/modules/azure"
	// "github.com/gruntwork-io/terratest/modules/azure"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

// An example of how to test the Terraform module in examples/terraform-azure-network-example using Terratest.
func TestTerraformAzureNetworkExample(t *testing.T) {
	t.Parallel()

	// subscriptionID is overridden by the environment variable "ARM_SUBSCRIPTION_ID"
	subscriptionID := ""
	region := azure.GetRandomStableRegion(t, []string{}, []string{}, subscriptionID)
	firstSubnetName := "terratest-subnet1"
	secondSubnetName := "terratest-subnet2"
	domainNameLabel := fmt.Sprintf("terratest-example-dnslabel-%s", strings.ToLower(random.UniqueId()))
	publicIPName := "terratest-example-ip"

	terraformOptions := &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: "../examples/terraform-azure-network-example",

		// Variables to pass to our Terraform code using -var options
		Vars: map[string]interface{}{
			"location":                    region,
			"public_ip_domain_name_label": domainNameLabel,
			"first_subnet_name":           firstSubnetName,
			"second_subnet_name":          secondSubnetName,
			"public_ip_name":              publicIPName,
		},
	}

	// At the end of the test, run `terraform destroy` to clean up any resources that were created
	defer terraform.Destroy(t, terraformOptions)

	// This will run `terraform init` and `terraform apply` and fail the test if there are any errors
	terraform.InitAndApply(t, terraformOptions)

	// Run `terraform output` to get the value of an output variable
	rgName := terraform.Output(t, terraformOptions, "resource_group_name")
	vnetName := terraform.Output(t, terraformOptions, "virtual_network_name")
	firstSubnetAddressPrefix := terraform.Output(t, terraformOptions, "first_subnet_address")
	secondSubnetAddressPrefix := terraform.Output(t, terraformOptions, "second_subnet_address")
	publicIPAddress := terraform.Output(t, terraformOptions, "public_ip_address")
	publicIPFqdn := terraform.Output(t, terraformOptions, "public_ip_fqdn")

	// Verify that the Virtual Network is now present in the Resource Group
	AssertVirtualNetworkExists(t, vnetName, rgName, subscriptionID)

	// Lookup the subnets in the given Virtual Network
	actualSubnets := GetSubnetsForVirtualNetwork(t, vnetName, rgName, subscriptionID)

	// Verify that each subnet has the expected address prefix
	actualFirstSubnet = filterSubnetsByName(actualSubnets, firstSubnetName)
	assert.Equal(t, firstSubnetAddressPrefix, actualFirstSubnet.AddressPrefix)
	actualSecondSubnet = filterSubnetsByName(actualSubnets, secondSubnetName)
	assert.Equal(t, secondSubnetAddressPrefix, actualSecondSubnet.AddressPrefix)

	// Verify that each subnet has the expected Network Security Group
	actualFirstSubnetNsg := GetNetworkSecurityGroupForSubnet(t, firstSubnetName, vnetName, rgName, subscriptionID)
	assert.Equal(t, expectedNsgName, actualFirstSubnetNsg.Name)
	actualSecondSubnetNsg := GetNetworkSecurityGroupForSubnet(t, secondSubnetName, vnetName, rgName, subscriptionID)
	assert.Equal(t, expectedNsgName, actualSecondSubnetNsg.Name)

	// Lookup the actual Public IP resource
	actualIP := GetPublicIP(t, rgName, publicIPName, subscriptionID)

	// Verify that the Public IP resource has the expected properties
	assert.Equal(t, publicIPAddress, actualIP.IPAddress)
	assert.Equal(t, publicIPFqdn, actualIP.FullDNSName)

	// Verify that our Public IP's domain name label is not available anymore
	available := CheckPublicDNSNameAvailability(t, region, domainNameLabel, subscriptionID)
	assert.False(t, available)
}

func filterSubnetsByName(subnets []Subnet, name string) Subnet {
	out := Subnet{}
	for _, s := range subnets {
		if s.Name == name {
			out = s
			break
		}
	}
	return out
}
