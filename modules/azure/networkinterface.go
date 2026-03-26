package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2019-09-01/network"
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// NetworkInterfaceExistsContext indicates whether the specified Azure Network Interface exists.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func NetworkInterfaceExistsContext(t testing.TestingT, ctx context.Context, nicName string, resGroupName string, subscriptionID string) bool {
	t.Helper()

	exists, err := NetworkInterfaceExistsContextE(ctx, nicName, resGroupName, subscriptionID)
	require.NoError(t, err)

	return exists
}

// NetworkInterfaceExists indicates whether the specified Azure Network Interface exists.
// This function would fail the test if there is an error.
//
// Deprecated: Use [NetworkInterfaceExistsContext] instead.
func NetworkInterfaceExists(t testing.TestingT, nicName string, resGroupName string, subscriptionID string) bool {
	t.Helper()

	return NetworkInterfaceExistsContext(t, context.Background(), nicName, resGroupName, subscriptionID) //nolint:staticcheck
}

// NetworkInterfaceExistsContextE indicates whether the specified Azure Network Interface exists.
// The ctx parameter supports cancellation and timeouts.
func NetworkInterfaceExistsContextE(ctx context.Context, nicName string, resGroupName string, subscriptionID string) (bool, error) {
	// Get the Network Interface
	_, err := GetNetworkInterfaceContextE(ctx, nicName, resGroupName, subscriptionID)
	if err != nil {
		if ResourceNotFoundErrorExists(err) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

// NetworkInterfaceExistsE indicates whether the specified Azure Network Interface exists.
//
// Deprecated: Use [NetworkInterfaceExistsContextE] instead.
func NetworkInterfaceExistsE(nicName string, resGroupName string, subscriptionID string) (bool, error) {
	return NetworkInterfaceExistsContextE(context.Background(), nicName, resGroupName, subscriptionID)
}

// GetNetworkInterfacePrivateIPsContext gets a list of the Private IPs of a Network Interface configs.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetNetworkInterfacePrivateIPsContext(t testing.TestingT, ctx context.Context, nicName string, resGroupName string, subscriptionID string) []string {
	t.Helper()

	IPs, err := GetNetworkInterfacePrivateIPsContextE(ctx, nicName, resGroupName, subscriptionID)
	require.NoError(t, err)

	return IPs
}

// GetNetworkInterfacePrivateIPs gets a list of the Private IPs of a Network Interface configs.
// This function would fail the test if there is an error.
//
// Deprecated: Use [GetNetworkInterfacePrivateIPsContext] instead.
func GetNetworkInterfacePrivateIPs(t testing.TestingT, nicName string, resGroupName string, subscriptionID string) []string {
	t.Helper()

	return GetNetworkInterfacePrivateIPsContext(t, context.Background(), nicName, resGroupName, subscriptionID) //nolint:staticcheck
}

// GetNetworkInterfacePrivateIPsContextE gets a list of the Private IPs of a Network Interface configs.
// The ctx parameter supports cancellation and timeouts.
func GetNetworkInterfacePrivateIPsContextE(ctx context.Context, nicName string, resGroupName string, subscriptionID string) ([]string, error) {
	var privateIPs []string

	// Get the Network Interface client
	nic, err := GetNetworkInterfaceContextE(ctx, nicName, resGroupName, subscriptionID)
	if err != nil {
		return privateIPs, err
	}

	// Get the Private IPs from each configuration
	for _, IPConfiguration := range *nic.IPConfigurations {
		privateIPs = append(privateIPs, *IPConfiguration.PrivateIPAddress)
	}

	return privateIPs, nil
}

// GetNetworkInterfacePrivateIPsE gets a list of the Private IPs of a Network Interface configs.
//
// Deprecated: Use [GetNetworkInterfacePrivateIPsContextE] instead.
func GetNetworkInterfacePrivateIPsE(nicName string, resGroupName string, subscriptionID string) ([]string, error) {
	return GetNetworkInterfacePrivateIPsContextE(context.Background(), nicName, resGroupName, subscriptionID)
}

// GetNetworkInterfacePublicIPsContext returns a list of all the Public IPs found in the Network Interface configurations.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetNetworkInterfacePublicIPsContext(t testing.TestingT, ctx context.Context, nicName string, resGroupName string, subscriptionID string) []string {
	t.Helper()

	IPs, err := GetNetworkInterfacePublicIPsContextE(ctx, nicName, resGroupName, subscriptionID)
	require.NoError(t, err)

	return IPs
}

// GetNetworkInterfacePublicIPs returns a list of all the Public IPs found in the Network Interface configurations.
// This function would fail the test if there is an error.
//
// Deprecated: Use [GetNetworkInterfacePublicIPsContext] instead.
func GetNetworkInterfacePublicIPs(t testing.TestingT, nicName string, resGroupName string, subscriptionID string) []string {
	t.Helper()

	return GetNetworkInterfacePublicIPsContext(t, context.Background(), nicName, resGroupName, subscriptionID) //nolint:staticcheck
}

// GetNetworkInterfacePublicIPsContextE returns a list of all the Public IPs found in the Network Interface configurations.
// The ctx parameter supports cancellation and timeouts.
func GetNetworkInterfacePublicIPsContextE(ctx context.Context, nicName string, resGroupName string, subscriptionID string) ([]string, error) {
	var publicIPs []string

	// Get the Network Interface client
	nic, err := GetNetworkInterfaceContextE(ctx, nicName, resGroupName, subscriptionID)
	if err != nil {
		return publicIPs, err
	}

	// Get the Public IPs from each configuration available
	for _, IPConfiguration := range *nic.IPConfigurations {
		// Iterate each config, for successful configurations check for a Public Address reference.
		// Not failing on errors as this is an optimistic accumulator.
		nicConfig, err := GetNetworkInterfaceConfigurationContextE(ctx, nicName, *IPConfiguration.Name, resGroupName, subscriptionID)
		if err == nil {
			if nicConfig.PublicIPAddress != nil {
				publicAddressID := GetNameFromResourceID(*nicConfig.PublicIPAddress.ID)

				publicIP, err := GetIPOfPublicIPAddressByNameContextE(ctx, publicAddressID, resGroupName, subscriptionID)
				if err == nil {
					publicIPs = append(publicIPs, publicIP)
				}
			}
		}
	}

	return publicIPs, nil
}

// GetNetworkInterfacePublicIPsE returns a list of all the Public IPs found in the Network Interface configurations.
//
// Deprecated: Use [GetNetworkInterfacePublicIPsContextE] instead.
func GetNetworkInterfacePublicIPsE(nicName string, resGroupName string, subscriptionID string) ([]string, error) {
	return GetNetworkInterfacePublicIPsContextE(context.Background(), nicName, resGroupName, subscriptionID)
}

// GetNetworkInterfaceConfigurationContextE gets a Network Interface Configuration in the specified Azure Resource Group.
// The ctx parameter supports cancellation and timeouts.
func GetNetworkInterfaceConfigurationContextE(ctx context.Context, nicName string, nicConfigName string, resGroupName string, subscriptionID string) (*network.InterfaceIPConfiguration, error) {
	// Validate Azure Resource Group
	resGroupName, err := getTargetAzureResourceGroupName(resGroupName)
	if err != nil {
		return nil, err
	}

	// Get the client reference
	client, err := GetNetworkInterfaceConfigurationClientE(subscriptionID)
	if err != nil {
		return nil, err
	}

	// Get the Network Interface
	nicConfig, err := client.Get(ctx, resGroupName, nicName, nicConfigName)
	if err != nil {
		return nil, err
	}

	return &nicConfig, nil
}

// GetNetworkInterfaceConfigurationE gets a Network Interface Configuration in the specified Azure Resource Group.
//
// Deprecated: Use [GetNetworkInterfaceConfigurationContextE] instead.
func GetNetworkInterfaceConfigurationE(nicName string, nicConfigName string, resGroupName string, subscriptionID string) (*network.InterfaceIPConfiguration, error) {
	return GetNetworkInterfaceConfigurationContextE(context.Background(), nicName, nicConfigName, resGroupName, subscriptionID)
}

// GetNetworkInterfaceConfigurationClientE creates a new Network Interface Configuration client in the specified Azure Subscription.
func GetNetworkInterfaceConfigurationClientE(subscriptionID string) (*network.InterfaceIPConfigurationsClient, error) {
	// Create a new client from client factory
	client, err := CreateNewNetworkInterfaceIPConfigurationClientE(subscriptionID)
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

// GetNetworkInterfaceContextE gets a Network Interface in the specified Azure Resource Group.
// The ctx parameter supports cancellation and timeouts.
func GetNetworkInterfaceContextE(ctx context.Context, nicName string, resGroupName string, subscriptionID string) (*network.Interface, error) {
	// Validate Azure Resource Group
	resGroupName, err := getTargetAzureResourceGroupName(resGroupName)
	if err != nil {
		return nil, err
	}

	// Get the client reference
	client, err := GetNetworkInterfaceClientE(subscriptionID)
	if err != nil {
		return nil, err
	}

	// Get the Network Interface
	nic, err := client.Get(ctx, resGroupName, nicName, "")
	if err != nil {
		return nil, err
	}

	return &nic, nil
}

// GetNetworkInterfaceE gets a Network Interface in the specified Azure Resource Group.
//
// Deprecated: Use [GetNetworkInterfaceContextE] instead.
func GetNetworkInterfaceE(nicName string, resGroupName string, subscriptionID string) (*network.Interface, error) {
	return GetNetworkInterfaceContextE(context.Background(), nicName, resGroupName, subscriptionID)
}

// GetNetworkInterfaceClientE creates a new Network Interface client in the specified Azure Subscription.
func GetNetworkInterfaceClientE(subscriptionID string) (*network.InterfacesClient, error) {
	// Create new NIC client from client factory
	client, err := CreateNewNetworkInterfacesClientE(subscriptionID)
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
