package azure

import (
	"context"
	"net"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2019-09-01/network"
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// VirtualNetworkExists indicates whether the specified Azure Virtual Network exists.
// This function would fail the test if there is an error.
func VirtualNetworkExists(t testing.TestingT, vnetName string, resGroupName string, subscriptionID string) bool {
	exists, err := VirtualNetworkExistsE(vnetName, resGroupName, subscriptionID)
	require.NoError(t, err)
	return exists
}

// VirtualNetworkExistsE indicates whether the specified Azure Virtual Network exists.
func VirtualNetworkExistsE(vnetName string, resGroupName string, subscriptionID string) (bool, error) {
	// Get the Virtual Network
	_, err := GetVirtualNetworkE(vnetName, resGroupName, subscriptionID)
	if err != nil {
		if ResourceNotFoundErrorExists(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// SubnetExists indicates whether the specified Azure Virtual Network Subnet exists.
// This function would fail the test if there is an error.
func SubnetExists(t testing.TestingT, subnetName string, vnetName string, resGroupName string, subscriptionID string) bool {
	exists, err := SubnetExistsE(subnetName, vnetName, resGroupName, subscriptionID)
	require.NoError(t, err)
	return exists
}

// SubnetExistsE indicates whether the specified Azure Virtual Network Subnet exists.
func SubnetExistsE(subnetName string, vnetName string, resGroupName string, subscriptionID string) (bool, error) {
	// Get the Subnet
	_, err := GetSubnetE(subnetName, vnetName, resGroupName, subscriptionID)
	if err != nil {
		if ResourceNotFoundErrorExists(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// CheckSubnetContainsIP checks if the Private IP is contined in the Subnet Address Range.
// This function would fail the test if there is an error.
func CheckSubnetContainsIP(t testing.TestingT, IP string, subnetName string, vnetName string, resGroupName string, subscriptionID string) bool {
	inRange, err := CheckSubnetContainsIPE(IP, subnetName, vnetName, resGroupName, subscriptionID)
	require.NoError(t, err)
	return inRange
}

// CheckSubnetContainsIPE checks if the Private IP is contined in the Subnet Address Range.
func CheckSubnetContainsIPE(ipAddress string, subnetName string, vnetName string, resGroupName string, subscriptionID string) (bool, error) {
	// Convert the IP to a net IP address
	ip := net.ParseIP(ipAddress)
	if ip == nil {
		return false, NewFailedToParseError("IP Address", ipAddress)
	}

	// Get Subnet
	subnet, err := GetSubnetE(subnetName, vnetName, resGroupName, subscriptionID)
	if err != nil {
		return false, err
	}

	// Get Subnet IP range, this required field is never nil therefore no exception handling required.
	subnetPrefix := *subnet.AddressPrefix

	// Check if the IP is in the Subnet Range using the net package
	_, ipNet, err := net.ParseCIDR(subnetPrefix)
	if err != nil {
		return false, NewFailedToParseError("Subnet Range", subnetPrefix)
	}

	return ipNet.Contains(ip), nil
}

// GetVirtualNetworkSubnets gets all Subnet names and their respective address prefixes in the
// specified Virtual Network. This function would fail the test if there is an error.
func GetVirtualNetworkSubnets(t testing.TestingT, vnetName string, resGroupName string, subscriptionID string) map[string]string {
	subnets, err := GetVirtualNetworkSubnetsE(vnetName, resGroupName, subscriptionID)
	require.NoError(t, err)
	return subnets
}

// GetVirtualNetworkSubnetsE gets all Subnet names and their respective address prefixes in the specified Virtual Network.
// Returning both the name and prefix together helps reduce calls for these frequently accessed properties.
func GetVirtualNetworkSubnetsE(vnetName string, resGroupName string, subscriptionID string) (map[string]string, error) {
	subNetDetails := map[string]string{}

	client, err := GetSubnetClientE(subscriptionID)
	if err != nil {
		return subNetDetails, err
	}

	subnets, err := client.List(context.Background(), resGroupName, vnetName)
	if err != nil {
		return subNetDetails, err
	}

	for _, v := range subnets.Values() {
		subnetName := v.Name
		subNetAddressPrefix := v.AddressPrefix

		subNetDetails[string(*subnetName)] = string(*subNetAddressPrefix)
	}
	return subNetDetails, nil
}

// GetVirtualNetworkDNSServerIPs gets a list of all Virtual Network DNS server IPs.
// This function would fail the test if there is an error.
func GetVirtualNetworkDNSServerIPs(t testing.TestingT, vnetName string, resGroupName string, subscriptionID string) []string {
	vnetDNSIPs, err := GetVirtualNetworkDNSServerIPsE(vnetName, resGroupName, subscriptionID)
	require.NoError(t, err)

	return vnetDNSIPs
}

// GetVirtualNetworkDNSServerIPsE gets a list of all Virtual Network DNS server IPs with Error.
func GetVirtualNetworkDNSServerIPsE(vnetName string, resGroupName string, subscriptionID string) ([]string, error) {
	// Get Virtual Network
	vnet, err := GetVirtualNetworkE(vnetName, resGroupName, subscriptionID)
	if err != nil {
		return nil, err
	}

	return *vnet.DhcpOptions.DNSServers, nil
}

// GetSubnetE gets a subnet.
func GetSubnetE(subnetName string, vnetName string, resGroupName string, subscriptionID string) (*network.Subnet, error) {
	// Validate Azure Resource Group
	resGroupName, err := getTargetAzureResourceGroupName(resGroupName)
	if err != nil {
		return nil, err
	}

	// Get the client reference
	client, err := GetSubnetClientE(subscriptionID)
	if err != nil {
		return nil, err
	}

	// Get the Subnet
	subnet, err := client.Get(context.Background(), resGroupName, vnetName, subnetName, "")
	if err != nil {
		return nil, err
	}

	return &subnet, nil
}

// GetSubnetClientE creates a subnet client.
func GetSubnetClientE(subscriptionID string) (*network.SubnetsClient, error) {
	// Create a new Subnet client from client factory
	client, err := CreateNewSubnetClientE(subscriptionID)
	if err != nil {
		return nil, err
	}

	// Create an authorizer
	authorizer, err := NewAuthorizer()
	if err != nil {
		return nil, err
	}
	client.Authorizer = *authorizer

	return client, nil
}

// GetVirtualNetworkE gets Virtual Network in the specified Azure Resource Group.
func GetVirtualNetworkE(vnetName string, resGroupName string, subscriptionID string) (*network.VirtualNetwork, error) {
	// Validate Azure Resource Group
	resGroupName, err := getTargetAzureResourceGroupName(resGroupName)
	if err != nil {
		return nil, err
	}

	// Get the client reference
	client, err := GetVirtualNetworksClientE(subscriptionID)
	if err != nil {
		return nil, err
	}

	// Get the Virtual Network
	vnet, err := client.Get(context.Background(), resGroupName, vnetName, "")
	if err != nil {
		return nil, err
	}
	return &vnet, nil
}

// GetVirtualNetworksClientE creates a virtual network client in the specified Azure Subscription.
func GetVirtualNetworksClientE(subscriptionID string) (*network.VirtualNetworksClient, error) {
	// Create a new Virtual Network client from client factory
	client, err := CreateNewVirtualNetworkClientE(subscriptionID)
	if err != nil {
		return nil, err
	}

	// Create an authorizer
	authorizer, err := NewAuthorizer()
	if err != nil {
		return nil, err
	}
	client.Authorizer = *authorizer

	return client, nil
}
