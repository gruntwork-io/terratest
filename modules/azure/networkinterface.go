package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v6"
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

// GetNetworkInterfacePrivateIPsContext gets a list of the Private IPs of a Network Interface configs.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetNetworkInterfacePrivateIPsContext(t testing.TestingT, ctx context.Context, nicName string, resGroupName string, subscriptionID string) []string {
	t.Helper()

	IPs, err := GetNetworkInterfacePrivateIPsContextE(ctx, nicName, resGroupName, subscriptionID)
	require.NoError(t, err)

	return IPs
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
	for _, IPConfiguration := range nic.Properties.IPConfigurations {
		privateIPs = append(privateIPs, *IPConfiguration.Properties.PrivateIPAddress)
	}

	return privateIPs, nil
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
	for _, IPConfiguration := range nic.Properties.IPConfigurations {
		// Iterate each config, for successful configurations check for a Public Address reference.
		// Not failing on errors as this is an optimistic accumulator.
		nicConfig, err := GetNetworkInterfaceConfigurationContextE(ctx, nicName, *IPConfiguration.Name, resGroupName, subscriptionID)
		if err == nil {
			if nicConfig.Properties.PublicIPAddress != nil {
				publicAddressID := GetNameFromResourceID(*nicConfig.Properties.PublicIPAddress.ID)

				publicIP, err := GetIPOfPublicIPAddressByNameContextE(ctx, publicAddressID, resGroupName, subscriptionID)
				if err == nil {
					publicIPs = append(publicIPs, publicIP)
				}
			}
		}
	}

	return publicIPs, nil
}

// GetNetworkInterfaceConfigurationContextE gets a Network Interface Configuration in the specified Azure Resource Group.
// The ctx parameter supports cancellation and timeouts.
func GetNetworkInterfaceConfigurationContextE(ctx context.Context, nicName string, nicConfigName string, resGroupName string, subscriptionID string) (*armnetwork.InterfaceIPConfiguration, error) {
	// Validate Azure Resource Group
	resGroupName, err := getTargetAzureResourceGroupName(resGroupName)
	if err != nil {
		return nil, err
	}

	// Get the client reference
	client, err := GetNetworkInterfaceConfigurationClientContextE(ctx, subscriptionID)
	if err != nil {
		return nil, err
	}

	// Get the Network Interface
	resp, err := client.Get(ctx, resGroupName, nicName, nicConfigName, nil)
	if err != nil {
		return nil, err
	}

	return &resp.InterfaceIPConfiguration, nil
}

// GetNetworkInterfaceConfigurationClientContextE creates a new Network Interface Configuration client in the specified Azure Subscription.
// The ctx parameter supports cancellation and timeouts.
func GetNetworkInterfaceConfigurationClientContextE(ctx context.Context, subscriptionID string) (*armnetwork.InterfaceIPConfigurationsClient, error) {
	return CreateNewNetworkInterfaceIPConfigurationClientContextE(ctx, subscriptionID)
}

// GetNetworkInterfaceConfigurationClientE creates a new Network Interface Configuration client in the specified Azure Subscription.
//
// Deprecated: Use [GetNetworkInterfaceConfigurationClientContextE] instead.
func GetNetworkInterfaceConfigurationClientE(subscriptionID string) (*armnetwork.InterfaceIPConfigurationsClient, error) {
	return GetNetworkInterfaceConfigurationClientContextE(context.Background(), subscriptionID)
}

// GetNetworkInterfaceContextE gets a Network Interface in the specified Azure Resource Group.
// The ctx parameter supports cancellation and timeouts.
func GetNetworkInterfaceContextE(ctx context.Context, nicName string, resGroupName string, subscriptionID string) (*armnetwork.Interface, error) {
	// Validate Azure Resource Group
	resGroupName, err := getTargetAzureResourceGroupName(resGroupName)
	if err != nil {
		return nil, err
	}

	// Get the client reference
	client, err := GetNetworkInterfaceClientContextE(ctx, subscriptionID)
	if err != nil {
		return nil, err
	}

	// Get the Network Interface
	resp, err := client.Get(ctx, resGroupName, nicName, nil)
	if err != nil {
		return nil, err
	}

	return &resp.Interface, nil
}

// GetNetworkInterfaceClientContextE creates a new Network Interface client in the specified Azure Subscription.
// The ctx parameter supports cancellation and timeouts.
func GetNetworkInterfaceClientContextE(ctx context.Context, subscriptionID string) (*armnetwork.InterfacesClient, error) {
	return CreateNewNetworkInterfacesClientContextE(ctx, subscriptionID)
}

// GetNetworkInterfaceClientE creates a new Network Interface client in the specified Azure Subscription.
//
// Deprecated: Use [GetNetworkInterfaceClientContextE] instead.
func GetNetworkInterfaceClientE(subscriptionID string) (*armnetwork.InterfacesClient, error) {
	return GetNetworkInterfaceClientContextE(context.Background(), subscriptionID)
}
