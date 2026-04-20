package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v6"
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// NetworkInterfaceExists indicates whether the specified Azure Network Interface exists.
// This function would fail the test if there is an error.
//
// Deprecated: Use [NetworkInterfaceExistsContext] instead.
func NetworkInterfaceExists(t testing.TestingT, nicName string, resGroupName string, subscriptionID string) bool {
	t.Helper()

	return NetworkInterfaceExistsContext(t, context.Background(), nicName, resGroupName, subscriptionID)
}

// NetworkInterfaceExistsE indicates whether the specified Azure Network Interface exists.
//
// Deprecated: Use [NetworkInterfaceExistsContextE] instead.
func NetworkInterfaceExistsE(nicName string, resGroupName string, subscriptionID string) (bool, error) {
	return NetworkInterfaceExistsContextE(context.Background(), nicName, resGroupName, subscriptionID)
}

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

// GetNetworkInterfacePrivateIPs gets a list of the Private IPs of a Network Interface configs.
// This function would fail the test if there is an error.
//
// Deprecated: Use [GetNetworkInterfacePrivateIPsContext] instead.
func GetNetworkInterfacePrivateIPs(t testing.TestingT, nicName string, resGroupName string, subscriptionID string) []string {
	t.Helper()

	return GetNetworkInterfacePrivateIPsContext(t, context.Background(), nicName, resGroupName, subscriptionID)
}

// GetNetworkInterfacePrivateIPsE gets a list of the Private IPs of a Network Interface configs.
//
// Deprecated: Use [GetNetworkInterfacePrivateIPsContextE] instead.
func GetNetworkInterfacePrivateIPsE(nicName string, resGroupName string, subscriptionID string) ([]string, error) {
	return GetNetworkInterfacePrivateIPsContextE(context.Background(), nicName, resGroupName, subscriptionID)
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
	nic, err := GetNetworkInterfaceContextE(ctx, nicName, resGroupName, subscriptionID)
	if err != nil {
		return nil, err
	}

	return ExtractNetworkInterfacePrivateIPs(nic), nil
}

// GetNetworkInterfacePublicIPs returns a list of all the Public IPs found in the Network Interface configurations.
// This function would fail the test if there is an error.
//
// Deprecated: Use [GetNetworkInterfacePublicIPsContext] instead.
func GetNetworkInterfacePublicIPs(t testing.TestingT, nicName string, resGroupName string, subscriptionID string) []string {
	t.Helper()

	return GetNetworkInterfacePublicIPsContext(t, context.Background(), nicName, resGroupName, subscriptionID)
}

// GetNetworkInterfacePublicIPsE returns a list of all the Public IPs found in the Network Interface configurations.
//
// Deprecated: Use [GetNetworkInterfacePublicIPsContextE] instead.
func GetNetworkInterfacePublicIPsE(nicName string, resGroupName string, subscriptionID string) ([]string, error) {
	return GetNetworkInterfacePublicIPsContextE(context.Background(), nicName, resGroupName, subscriptionID)
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

	if nic == nil || nic.Properties == nil {
		return publicIPs, nil
	}

	// Get the Public IPs from each configuration available.
	// Not failing on individual errors as this is an optimistic accumulator —
	// it collects what it can and skips configurations that fail.
	for _, IPConfiguration := range nic.Properties.IPConfigurations {
		if IPConfiguration == nil || IPConfiguration.Name == nil {
			continue
		}

		nicConfig, err := GetNetworkInterfaceConfigurationContextE(ctx, nicName, *IPConfiguration.Name, resGroupName, subscriptionID)
		if err != nil {
			continue
		}

		if nicConfig == nil || nicConfig.Properties == nil || nicConfig.Properties.PublicIPAddress == nil ||
			nicConfig.Properties.PublicIPAddress.ID == nil {
			continue
		}

		publicAddressID := GetNameFromResourceID(*nicConfig.Properties.PublicIPAddress.ID)

		publicIP, err := GetIPOfPublicIPAddressByNameContextE(ctx, publicAddressID, resGroupName, subscriptionID)
		if err != nil {
			continue
		}

		publicIPs = append(publicIPs, publicIP)
	}

	return publicIPs, nil
}

// GetNetworkInterfaceConfigurationE gets a Network Interface IP Configuration in the specified Azure Resource Group.
//
// Deprecated: Use [GetNetworkInterfaceConfigurationContextE] instead.
func GetNetworkInterfaceConfigurationE(nicName string, nicConfigName string, resGroupName string, subscriptionID string) (*armnetwork.InterfaceIPConfiguration, error) {
	return GetNetworkInterfaceConfigurationContextE(context.Background(), nicName, nicConfigName, resGroupName, subscriptionID)
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
	return CreateNetworkInterfaceIPConfigurationClientContextE(ctx, subscriptionID)
}

// GetNetworkInterfaceConfigurationClientE creates a new Network Interface Configuration client in the specified Azure Subscription.
//
// Deprecated: Use [GetNetworkInterfaceConfigurationClientContextE] instead.
func GetNetworkInterfaceConfigurationClientE(subscriptionID string) (*armnetwork.InterfaceIPConfigurationsClient, error) {
	return GetNetworkInterfaceConfigurationClientContextE(context.Background(), subscriptionID)
}

// GetNetworkInterfaceE gets a Network Interface in the specified Azure Resource Group.
//
// Deprecated: Use [GetNetworkInterfaceContextE] instead.
func GetNetworkInterfaceE(nicName string, resGroupName string, subscriptionID string) (*armnetwork.Interface, error) {
	return GetNetworkInterfaceContextE(context.Background(), nicName, resGroupName, subscriptionID)
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

	return GetNetworkInterfaceWithClient(ctx, client, resGroupName, nicName)
}

// GetNetworkInterfaceWithClient gets a Network Interface using the provided InterfacesClient.
func GetNetworkInterfaceWithClient(ctx context.Context, client *armnetwork.InterfacesClient, resGroupName string, nicName string) (*armnetwork.Interface, error) {
	resp, err := client.Get(ctx, resGroupName, nicName, nil)
	if err != nil {
		return nil, err
	}

	return &resp.Interface, nil
}

// ExtractNetworkInterfacePrivateIPs gets a list of the Private IPs from a Network Interface.
func ExtractNetworkInterfacePrivateIPs(nic *armnetwork.Interface) []string {
	if nic == nil || nic.Properties == nil {
		return nil
	}

	privateIPs := make([]string, 0, len(nic.Properties.IPConfigurations))

	for _, ipConfig := range nic.Properties.IPConfigurations {
		if ipConfig == nil || ipConfig.Properties == nil || ipConfig.Properties.PrivateIPAddress == nil {
			continue
		}

		privateIPs = append(privateIPs, *ipConfig.Properties.PrivateIPAddress)
	}

	return privateIPs
}

// GetNetworkInterfaceClientContextE creates a new Network Interface client in the specified Azure Subscription.
// The ctx parameter supports cancellation and timeouts.
func GetNetworkInterfaceClientContextE(ctx context.Context, subscriptionID string) (*armnetwork.InterfacesClient, error) {
	return CreateNetworkInterfacesClientContextE(ctx, subscriptionID)
}

// GetNetworkInterfaceClientE creates a new Network Interface client in the specified Azure Subscription.
//
// Deprecated: Use [GetNetworkInterfaceClientContextE] instead.
func GetNetworkInterfaceClientE(subscriptionID string) (*armnetwork.InterfacesClient, error) {
	return GetNetworkInterfaceClientContextE(context.Background(), subscriptionID)
}
