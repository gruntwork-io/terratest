package azure

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v6"
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// PublicAddressExistsContext indicates whether the specified Azure Public Address exists.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func PublicAddressExistsContext(t testing.TestingT, ctx context.Context, publicAddressName string, resGroupName string, subscriptionID string) bool {
	t.Helper()

	exists, err := PublicAddressExistsContextE(ctx, publicAddressName, resGroupName, subscriptionID)
	require.NoError(t, err)

	return exists
}

// PublicAddressExistsContextE indicates whether the specified Azure Public Address exists.
// The ctx parameter supports cancellation and timeouts.
func PublicAddressExistsContextE(ctx context.Context, publicAddressName string, resGroupName string, subscriptionID string) (bool, error) {
	_, err := GetPublicIPAddressContextE(ctx, publicAddressName, resGroupName, subscriptionID)
	if err != nil {
		if ResourceNotFoundErrorExists(err) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

// GetIPOfPublicIPAddressByNameContext gets the IP of the specified Public IP Address.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetIPOfPublicIPAddressByNameContext(t testing.TestingT, ctx context.Context, publicAddressName string, resGroupName string, subscriptionID string) string {
	t.Helper()

	IP, err := GetIPOfPublicIPAddressByNameContextE(ctx, publicAddressName, resGroupName, subscriptionID)
	require.NoError(t, err)

	return IP
}

// GetIPOfPublicIPAddressByNameContextE gets the IP of the specified Public IP Address.
// The ctx parameter supports cancellation and timeouts.
func GetIPOfPublicIPAddressByNameContextE(ctx context.Context, publicAddressName string, resGroupName string, subscriptionID string) (string, error) {
	pip, err := GetPublicIPAddressContextE(ctx, publicAddressName, resGroupName, subscriptionID)
	if err != nil {
		return "", err
	}

	if pip.Properties == nil || pip.Properties.IPAddress == nil {
		return "", fmt.Errorf("public IP address %q has no IP address assigned", publicAddressName)
	}

	return *pip.Properties.IPAddress, nil
}

// CheckPublicDNSNameAvailabilityContext checks whether a domain name in the cloudapp.azure.com zone
// is available for use. This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func CheckPublicDNSNameAvailabilityContext(t testing.TestingT, ctx context.Context, location string, domainNameLabel string, subscriptionID string) bool {
	t.Helper()

	available, err := CheckPublicDNSNameAvailabilityContextE(ctx, location, domainNameLabel, subscriptionID)
	if err != nil {
		return false
	}

	return available
}

// CheckPublicDNSNameAvailabilityContextE checks whether a domain name in the cloudapp.azure.com zone
// is available for use.
// The ctx parameter supports cancellation and timeouts.
func CheckPublicDNSNameAvailabilityContextE(ctx context.Context, location string, domainNameLabel string, subscriptionID string) (bool, error) {
	client, err := CreateNetworkManagementClientE(subscriptionID)
	if err != nil {
		return false, err
	}

	res, err := client.CheckDNSNameAvailability(ctx, location, domainNameLabel, nil)
	if err != nil {
		return false, err
	}

	return *res.Available, nil
}

// GetPublicIPAddressContextE gets a Public IP Address in the specified Azure Resource Group.
// The ctx parameter supports cancellation and timeouts.
func GetPublicIPAddressContextE(ctx context.Context, publicIPAddressName string, resGroupName string, subscriptionID string) (*armnetwork.PublicIPAddress, error) {
	resGroupName, err := getTargetAzureResourceGroupName(resGroupName)
	if err != nil {
		return nil, err
	}

	client, err := GetPublicIPAddressClientE(subscriptionID)
	if err != nil {
		return nil, err
	}

	resp, err := client.Get(ctx, resGroupName, publicIPAddressName, nil)
	if err != nil {
		return nil, err
	}

	return &resp.PublicIPAddress, nil
}

// GetPublicIPAddressClientE creates a Public IP Addresses client in the specified Azure Subscription.
func GetPublicIPAddressClientE(subscriptionID string) (*armnetwork.PublicIPAddressesClient, error) {
	return CreatePublicIPAddressesClientE(subscriptionID)
}
