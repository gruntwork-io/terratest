package azure

import (
	"context"
	"net"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2019-09-01/network"
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// VirtualNetworkExistsContext indicates whether the specified Azure Virtual Network exists.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func VirtualNetworkExistsContext(t testing.TestingT, ctx context.Context, vnetName string, resGroupName string, subscriptionID string) bool {
	exists, err := VirtualNetworkExistsContextE(ctx, vnetName, resGroupName, subscriptionID)
	require.NoError(t, err)

	return exists
}

// VirtualNetworkExistsContextE indicates whether the specified Azure Virtual Network exists.
// The ctx parameter supports cancellation and timeouts.
func VirtualNetworkExistsContextE(ctx context.Context, vnetName string, resGroupName string, subscriptionID string) (bool, error) {
	_, err := GetVirtualNetworkContextE(ctx, vnetName, resGroupName, subscriptionID)
	if err != nil {
		if ResourceNotFoundErrorExists(err) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

// VirtualNetworkExists indicates whether the specified Azure Virtual Network exists.
// This function would fail the test if there is an error.
//
// Deprecated: Use [VirtualNetworkExistsContext] instead.
func VirtualNetworkExists(t testing.TestingT, vnetName string, resGroupName string, subscriptionID string) bool {
	return VirtualNetworkExistsContext(t, context.Background(), vnetName, resGroupName, subscriptionID) //nolint:staticcheck
}

// VirtualNetworkExistsE indicates whether the specified Azure Virtual Network exists.
//
// Deprecated: Use [VirtualNetworkExistsContextE] instead.
func VirtualNetworkExistsE(vnetName string, resGroupName string, subscriptionID string) (bool, error) {
	return VirtualNetworkExistsContextE(context.Background(), vnetName, resGroupName, subscriptionID)
}

// SubnetExistsContext indicates whether the specified Azure Virtual Network Subnet exists.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func SubnetExistsContext(t testing.TestingT, ctx context.Context, subnetName string, vnetName string, resGroupName string, subscriptionID string) bool {
	exists, err := SubnetExistsContextE(ctx, subnetName, vnetName, resGroupName, subscriptionID)
	require.NoError(t, err)

	return exists
}

// SubnetExistsContextE indicates whether the specified Azure Virtual Network Subnet exists.
// The ctx parameter supports cancellation and timeouts.
func SubnetExistsContextE(ctx context.Context, subnetName string, vnetName string, resGroupName string, subscriptionID string) (bool, error) {
	_, err := GetSubnetContextE(ctx, subnetName, vnetName, resGroupName, subscriptionID)
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
//
// Deprecated: Use [SubnetExistsContext] instead.
func SubnetExists(t testing.TestingT, subnetName string, vnetName string, resGroupName string, subscriptionID string) bool {
	return SubnetExistsContext(t, context.Background(), subnetName, vnetName, resGroupName, subscriptionID) //nolint:staticcheck
}

// SubnetExistsE indicates whether the specified Azure Virtual Network Subnet exists.
//
// Deprecated: Use [SubnetExistsContextE] instead.
func SubnetExistsE(subnetName string, vnetName string, resGroupName string, subscriptionID string) (bool, error) {
	return SubnetExistsContextE(context.Background(), subnetName, vnetName, resGroupName, subscriptionID)
}

// CheckSubnetContainsIPContext checks if the Private IP is contained in the Subnet Address Range.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func CheckSubnetContainsIPContext(t testing.TestingT, ctx context.Context, ipAddress string, subnetName string, vnetName string, resGroupName string, subscriptionID string) bool {
	inRange, err := CheckSubnetContainsIPContextE(ctx, ipAddress, subnetName, vnetName, resGroupName, subscriptionID)
	require.NoError(t, err)

	return inRange
}

// CheckSubnetContainsIPContextE checks if the Private IP is contained in the Subnet Address Range.
// The ctx parameter supports cancellation and timeouts.
func CheckSubnetContainsIPContextE(ctx context.Context, ipAddress string, subnetName string, vnetName string, resGroupName string, subscriptionID string) (bool, error) {
	// Convert the IP to a net IP address
	ip := net.ParseIP(ipAddress)
	if ip == nil {
		return false, NewFailedToParseError("IP Address", ipAddress)
	}

	// Get Subnet
	subnet, err := GetSubnetContextE(ctx, subnetName, vnetName, resGroupName, subscriptionID)
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

// CheckSubnetContainsIP checks if the Private IP is contained in the Subnet Address Range.
// This function would fail the test if there is an error.
//
// Deprecated: Use [CheckSubnetContainsIPContext] instead.
func CheckSubnetContainsIP(t testing.TestingT, ipAddress string, subnetName string, vnetName string, resGroupName string, subscriptionID string) bool {
	return CheckSubnetContainsIPContext(t, context.Background(), ipAddress, subnetName, vnetName, resGroupName, subscriptionID) //nolint:staticcheck
}

// CheckSubnetContainsIPE checks if the Private IP is contained in the Subnet Address Range.
//
// Deprecated: Use [CheckSubnetContainsIPContextE] instead.
func CheckSubnetContainsIPE(ipAddress string, subnetName string, vnetName string, resGroupName string, subscriptionID string) (bool, error) {
	return CheckSubnetContainsIPContextE(context.Background(), ipAddress, subnetName, vnetName, resGroupName, subscriptionID)
}

// GetVirtualNetworkSubnetsContext gets all Subnet names and their respective address prefixes in the
// specified Virtual Network. This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetVirtualNetworkSubnetsContext(t testing.TestingT, ctx context.Context, vnetName string, resGroupName string, subscriptionID string) map[string]string {
	subnets, err := GetVirtualNetworkSubnetsContextE(ctx, vnetName, resGroupName, subscriptionID)
	require.NoError(t, err)

	return subnets
}

// GetVirtualNetworkSubnetsContextE gets all Subnet names and their respective address prefixes in the specified Virtual Network.
// Returning both the name and prefix together helps reduce calls for these frequently accessed properties.
// The ctx parameter supports cancellation and timeouts.
func GetVirtualNetworkSubnetsContextE(ctx context.Context, vnetName string, resGroupName string, subscriptionID string) (map[string]string, error) {
	subNetDetails := map[string]string{}

	client, err := GetSubnetClientE(subscriptionID)
	if err != nil {
		return subNetDetails, err
	}

	subnets, err := client.List(ctx, resGroupName, vnetName)
	if err != nil {
		return subNetDetails, err
	}

	for _, v := range subnets.Values() {
		subnetName := v.Name
		subNetAddressPrefix := v.AddressPrefix

		subNetDetails[*subnetName] = *subNetAddressPrefix
	}

	return subNetDetails, nil
}

// GetVirtualNetworkSubnets gets all Subnet names and their respective address prefixes in the
// specified Virtual Network. This function would fail the test if there is an error.
//
// Deprecated: Use [GetVirtualNetworkSubnetsContext] instead.
func GetVirtualNetworkSubnets(t testing.TestingT, vnetName string, resGroupName string, subscriptionID string) map[string]string {
	return GetVirtualNetworkSubnetsContext(t, context.Background(), vnetName, resGroupName, subscriptionID) //nolint:staticcheck
}

// GetVirtualNetworkSubnetsE gets all Subnet names and their respective address prefixes in the specified Virtual Network.
// Returning both the name and prefix together helps reduce calls for these frequently accessed properties.
//
// Deprecated: Use [GetVirtualNetworkSubnetsContextE] instead.
func GetVirtualNetworkSubnetsE(vnetName string, resGroupName string, subscriptionID string) (map[string]string, error) {
	return GetVirtualNetworkSubnetsContextE(context.Background(), vnetName, resGroupName, subscriptionID)
}

// GetVirtualNetworkDNSServerIPsContext gets a list of all Virtual Network DNS server IPs.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetVirtualNetworkDNSServerIPsContext(t testing.TestingT, ctx context.Context, vnetName string, resGroupName string, subscriptionID string) []string {
	vnetDNSIPs, err := GetVirtualNetworkDNSServerIPsContextE(ctx, vnetName, resGroupName, subscriptionID)
	require.NoError(t, err)

	return vnetDNSIPs
}

// GetVirtualNetworkDNSServerIPsContextE gets a list of all Virtual Network DNS server IPs.
// The ctx parameter supports cancellation and timeouts.
func GetVirtualNetworkDNSServerIPsContextE(ctx context.Context, vnetName string, resGroupName string, subscriptionID string) ([]string, error) {
	vnet, err := GetVirtualNetworkContextE(ctx, vnetName, resGroupName, subscriptionID)
	if err != nil {
		return nil, err
	}

	return *vnet.DhcpOptions.DNSServers, nil
}

// GetVirtualNetworkDNSServerIPs gets a list of all Virtual Network DNS server IPs.
// This function would fail the test if there is an error.
//
// Deprecated: Use [GetVirtualNetworkDNSServerIPsContext] instead.
func GetVirtualNetworkDNSServerIPs(t testing.TestingT, vnetName string, resGroupName string, subscriptionID string) []string {
	return GetVirtualNetworkDNSServerIPsContext(t, context.Background(), vnetName, resGroupName, subscriptionID) //nolint:staticcheck
}

// GetVirtualNetworkDNSServerIPsE gets a list of all Virtual Network DNS server IPs.
//
// Deprecated: Use [GetVirtualNetworkDNSServerIPsContextE] instead.
func GetVirtualNetworkDNSServerIPsE(vnetName string, resGroupName string, subscriptionID string) ([]string, error) {
	return GetVirtualNetworkDNSServerIPsContextE(context.Background(), vnetName, resGroupName, subscriptionID)
}

// GetSubnetContextE gets a subnet.
// The ctx parameter supports cancellation and timeouts.
func GetSubnetContextE(ctx context.Context, subnetName string, vnetName string, resGroupName string, subscriptionID string) (*network.Subnet, error) {
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
	subnet, err := client.Get(ctx, resGroupName, vnetName, subnetName, "")
	if err != nil {
		return nil, err
	}

	return &subnet, nil
}

// GetSubnetE gets a subnet.
//
// Deprecated: Use [GetSubnetContextE] instead.
func GetSubnetE(subnetName string, vnetName string, resGroupName string, subscriptionID string) (*network.Subnet, error) {
	return GetSubnetContextE(context.Background(), subnetName, vnetName, resGroupName, subscriptionID)
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

// GetVirtualNetworkContextE gets Virtual Network in the specified Azure Resource Group.
// The ctx parameter supports cancellation and timeouts.
func GetVirtualNetworkContextE(ctx context.Context, vnetName string, resGroupName string, subscriptionID string) (*network.VirtualNetwork, error) {
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
	vnet, err := client.Get(ctx, resGroupName, vnetName, "")
	if err != nil {
		return nil, err
	}

	return &vnet, nil
}

// GetVirtualNetworkE gets Virtual Network in the specified Azure Resource Group.
//
// Deprecated: Use [GetVirtualNetworkContextE] instead.
func GetVirtualNetworkE(vnetName string, resGroupName string, subscriptionID string) (*network.VirtualNetwork, error) {
	return GetVirtualNetworkContextE(context.Background(), vnetName, resGroupName, subscriptionID)
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
