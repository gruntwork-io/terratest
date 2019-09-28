// +build azure

// NOTE: We use build tags to differentiate azure testing because we currently do not have azure access setup for
// CircleCI.

package azure

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	subscriptionID        = ""
	rgName                = "terratest-rg"
	location              = "northeurope"
	vnetName              = "terratest-vnet"
	expectedSubnets       = map[string]string{"terratest-subnet1": "172.17.0.0/24", "terratest-subnet2": "172.17.1.0/24"}
	expectedVnetTags      = map[string]string{"terratest": "true", "environment": "dev"}
	expectedNsgs          = map[string]string{"terratest-subnet1": "terratest-nsg1", "terratest-subnet2": "terratest-nsg2"}
	expectedSecurityRules = map[string][]string{"terratest-subnet1": []string{"Allow_HTTPS_Inbound", "Allow_SSH_Inbound"}, "terratest-subnet2": []string{}}
	publicIPName          = "terratest-vm-ip"
)

func TestGetVirtualNetworkClient(t *testing.T) {
	t.Parallel()

	vnetClient, err := GetVirtualNetworkClient(subscriptionID)
	require.NoError(t, err)
	assert.NotEmpty(t, *vnetClient)
}

func TestGetVirtualNetwork(t *testing.T) {
	t.Parallel()

	vnet := GetVirtualNetwork(t, vnetName, rgName, subscriptionID)
	assert.NotEmpty(t, vnet)
}

func TestAssertVirtualNetworkExists(t *testing.T) {
	t.Parallel()

	AssertVirtualNetworkExists(t, vnetName, rgName, subscriptionID)
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

func TestGetSubnetsForVirtualNetwork(t *testing.T) {
	t.Parallel()

	subnets := GetSubnetsForVirtualNetwork(t, vnetName, rgName, subscriptionID)
	for subnetName, expectedAddressPrefix := range expectedSubnets {
		t.Run(subnetName, func(t *testing.T) {
			actualSubnet := filterSubnetsByName(subnets, subnetName)

			assert.Equal(t, expectedAddressPrefix, actualSubnet.AddressPrefix)
		})
	}
}

func TestGetTagsForVirtualNetwork(t *testing.T) {
	t.Parallel()

	tags := GetTagsForVirtualNetwork(t, vnetName, rgName, subscriptionID)

	assert.Equal(t, expectedVnetTags["terratest"], tags["terratest"])
	assert.Equal(t, expectedVnetTags["environment"], tags["environment"])
}

func TestGetNetworkSecurityGroupForSubnet(t *testing.T) {
	t.Parallel()

	for subnetName, expectedNsgName := range expectedNsgs {
		t.Run(subnetName, func(t *testing.T) {
			nsg := GetNetworkSecurityGroupForSubnet(t, subnetName, vnetName, rgName, subscriptionID)

			assert.Equal(t, expectedNsgName, nsg.Name)
		})
	}
}

func TestGetNetworkSecurityGroupForSubnetWithNetworkSecurityRules(t *testing.T) {
	t.Parallel()

	for subnetName, expectedRuleNames := range expectedSecurityRules {
		t.Run(subnetName, func(t *testing.T) {
			nsg := GetNetworkSecurityGroupForSubnet(t, subnetName, vnetName, rgName, subscriptionID)

			assert.ElementsMatch(t, expectedRuleNames, nsg.SecurityRulesNames)
		})
	}
}

func TestGetPublicIPClient(t *testing.T) {
	t.Parallel()

	ipClient, err := GetPublicIPClient(subscriptionID)
	require.NoError(t, err)

	assert.NotEmpty(t, *ipClient)
}

func TestCheckPublicDNSNameAvailability(t *testing.T) {
	t.Parallel()

	inUse := "terratest487639"
	inUseRes := CheckPublicDNSNameAvailability(t, location, inUse, subscriptionID)
	res := CheckPublicDNSNameAvailability(t, location, "ku2dhtmk97qzx", subscriptionID)

	assert.False(t, inUseRes)
	assert.True(t, res)
}

func TestGetPublicIP(t *testing.T) {
	t.Parallel()

	expectedIPAddress := "94.245.92.160"
	expectedIPAddressVersion := "IPv4"
	expectedAllocationMethod := "Dynamic"
	expectedFullDNSName := "terratest487639.northeurope.cloudapp.azure.com"
	actual := GetPublicIP(t, rgName, publicIPName, subscriptionID)

	assert.Equal(t, expectedIPAddress, actual.IPAddress)
	assert.Equal(t, expectedIPAddressVersion, actual.IPAddressVersion)
	assert.Equal(t, expectedAllocationMethod, actual.AllocationMethod)
	assert.Equal(t, expectedFullDNSName, actual.FullDNSName)
}
