package azure

import (
	"context"
	"errors"
	"net"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v6"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
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
	client, err := GetSubnetClientContextE(ctx, subscriptionID)
	if err != nil {
		return false, err
	}

	return CheckSubnetContainsIPWithClient(ctx, client, ipAddress, subnetName, vnetName, resGroupName)
}

// CheckSubnetContainsIPWithClient checks if the Private IP is contained in the Subnet Address Range
// using a pre-built SubnetsClient.
func CheckSubnetContainsIPWithClient(ctx context.Context, client *armnetwork.SubnetsClient, ipAddress string, subnetName string, vnetName string, resGroupName string) (bool, error) {

	ip := net.ParseIP(ipAddress)
	if ip == nil {
		return false, NewFailedToParseError("IP Address", ipAddress)
	}

	subnet, err := GetSubnetWithClient(ctx, client, resGroupName, vnetName, subnetName)
	if err != nil {
		return false, err
	}

	if subnet.Properties == nil || subnet.Properties.AddressPrefix == nil {
		return false, errors.New("subnet has no address prefix")
	}

	_, ipNet, parseErr := net.ParseCIDR(*subnet.Properties.AddressPrefix)
	if parseErr != nil {
		return false, NewFailedToParseError("Subnet Range", *subnet.Properties.AddressPrefix)
	}

	return ipNet.Contains(ip), nil
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
	client, err := GetSubnetClientContextE(ctx, subscriptionID)
	if err != nil {
		return nil, err
	}

	return GetVirtualNetworkSubnetsWithClient(ctx, client, resGroupName, vnetName)
}

// GetVirtualNetworkSubnetsWithClient gets all Subnet names and their respective address prefixes
// using the provided SubnetsClient.
func GetVirtualNetworkSubnetsWithClient(ctx context.Context, client *armnetwork.SubnetsClient, resGroupName string, vnetName string) (map[string]string, error) {
	subNetDetails := map[string]string{}

	pager := client.NewListPager(resGroupName, vnetName, nil)
	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, v := range page.Value {
			if v == nil || v.Name == nil || v.Properties == nil || v.Properties.AddressPrefix == nil {
				continue
			}

			subNetDetails[*v.Name] = *v.Properties.AddressPrefix
		}
	}

	return subNetDetails, nil
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

	return ExtractVirtualNetworkDNSServerIPs(vnet), nil
}

// ExtractVirtualNetworkDNSServerIPs gets a list of all DNS server IPs from a VirtualNetwork.
func ExtractVirtualNetworkDNSServerIPs(vnet *armnetwork.VirtualNetwork) []string {
	if vnet == nil || vnet.Properties == nil || vnet.Properties.DhcpOptions == nil {
		return nil
	}

	dnsServers := make([]string, 0, len(vnet.Properties.DhcpOptions.DNSServers))
	for _, s := range vnet.Properties.DhcpOptions.DNSServers {
		if s == nil {
			continue
		}

		dnsServers = append(dnsServers, *s)
	}

	return dnsServers
}

// GetSubnetContextE gets a subnet.
// The ctx parameter supports cancellation and timeouts.
func GetSubnetContextE(ctx context.Context, subnetName string, vnetName string, resGroupName string, subscriptionID string) (*armnetwork.Subnet, error) {

	resGroupName, err := getTargetAzureResourceGroupName(resGroupName)
	if err != nil {
		return nil, err
	}

	client, err := GetSubnetClientContextE(ctx, subscriptionID)
	if err != nil {
		return nil, err
	}

	return GetSubnetWithClient(ctx, client, resGroupName, vnetName, subnetName)
}

// GetSubnetWithClient gets a subnet using the provided SubnetsClient.
func GetSubnetWithClient(ctx context.Context, client *armnetwork.SubnetsClient, resGroupName string, vnetName string, subnetName string) (*armnetwork.Subnet, error) {
	resp, err := client.Get(ctx, resGroupName, vnetName, subnetName, nil)
	if err != nil {
		return nil, err
	}

	return &resp.Subnet, nil
}

// GetSubnetClientContextE creates a subnet client.
// The ctx parameter supports cancellation and timeouts.
func GetSubnetClientContextE(ctx context.Context, subscriptionID string) (*armnetwork.SubnetsClient, error) {
	return CreateSubnetClientContextE(ctx, subscriptionID)
}

// GetVirtualNetworkContextE gets Virtual Network in the specified Azure Resource Group.
// The ctx parameter supports cancellation and timeouts.
func GetVirtualNetworkContextE(ctx context.Context, vnetName string, resGroupName string, subscriptionID string) (*armnetwork.VirtualNetwork, error) {

	resGroupName, err := getTargetAzureResourceGroupName(resGroupName)
	if err != nil {
		return nil, err
	}

	client, err := GetVirtualNetworksClientContextE(ctx, subscriptionID)
	if err != nil {
		return nil, err
	}

	return GetVirtualNetworkWithClient(ctx, client, resGroupName, vnetName)
}

// GetVirtualNetworkWithClient gets a Virtual Network using the provided VirtualNetworksClient.
func GetVirtualNetworkWithClient(ctx context.Context, client *armnetwork.VirtualNetworksClient, resGroupName string, vnetName string) (*armnetwork.VirtualNetwork, error) {
	resp, err := client.Get(ctx, resGroupName, vnetName, nil)
	if err != nil {
		return nil, err
	}

	return &resp.VirtualNetwork, nil
}

// GetVirtualNetworksClientContextE creates a virtual network client in the specified Azure Subscription.
// The ctx parameter supports cancellation and timeouts.
func GetVirtualNetworksClientContextE(ctx context.Context, subscriptionID string) (*armnetwork.VirtualNetworksClient, error) {
	return CreateVirtualNetworkClientContextE(ctx, subscriptionID)
}
