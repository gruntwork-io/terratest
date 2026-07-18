package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v6"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
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
	nic, err := GetNetworkInterfaceContextE(ctx, nicName, resGroupName, subscriptionID)
	if err != nil {
		return nil, err
	}

	return ExtractNetworkInterfacePrivateIPs(nic), nil
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

	nic, err := GetNetworkInterfaceContextE(ctx, nicName, resGroupName, subscriptionID)
	if err != nil {
		return publicIPs, err
	}

	if nic == nil || nic.Properties == nil {
		return publicIPs, nil
	}

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

// GetNetworkInterfaceConfigurationContextE gets a Network Interface Configuration in the specified Azure Resource Group.
// The ctx parameter supports cancellation and timeouts.
func GetNetworkInterfaceConfigurationContextE(ctx context.Context, nicName string, nicConfigName string, resGroupName string, subscriptionID string) (*armnetwork.InterfaceIPConfiguration, error) {

	resGroupName, err := getTargetAzureResourceGroupName(resGroupName)
	if err != nil {
		return nil, err
	}

	client, err := GetNetworkInterfaceConfigurationClientContextE(ctx, subscriptionID)
	if err != nil {
		return nil, err
	}

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

// GetNetworkInterfaceContextE gets a Network Interface in the specified Azure Resource Group.
// The ctx parameter supports cancellation and timeouts.
func GetNetworkInterfaceContextE(ctx context.Context, nicName string, resGroupName string, subscriptionID string) (*armnetwork.Interface, error) {

	resGroupName, err := getTargetAzureResourceGroupName(resGroupName)
	if err != nil {
		return nil, err
	}

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
