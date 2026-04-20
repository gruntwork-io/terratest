package azure

import (
	"context"
	"errors"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v6"
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// PublicAddressExists indicates whether the specified Azure Public Address exists.
// This function would fail the test if there is an error.
//
// Deprecated: Use [PublicAddressExistsContext] instead.
func PublicAddressExists(t testing.TestingT, publicAddressName string, resGroupName string, subscriptionID string) bool {
	t.Helper()

	return PublicAddressExistsContext(t, context.Background(), publicAddressName, resGroupName, subscriptionID)
}

// PublicAddressExistsE indicates whether the specified Azure Public Address exists.
//
// Deprecated: Use [PublicAddressExistsContextE] instead.
func PublicAddressExistsE(publicAddressName string, resGroupName string, subscriptionID string) (bool, error) {
	return PublicAddressExistsContextE(context.Background(), publicAddressName, resGroupName, subscriptionID)
}

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

// GetIPOfPublicIPAddressByName gets the IP of the specified Public IP Address.
// This function would fail the test if there is an error.
//
// Deprecated: Use [GetIPOfPublicIPAddressByNameContext] instead.
func GetIPOfPublicIPAddressByName(t testing.TestingT, publicAddressName string, resGroupName string, subscriptionID string) string {
	t.Helper()

	return GetIPOfPublicIPAddressByNameContext(t, context.Background(), publicAddressName, resGroupName, subscriptionID)
}

// GetIPOfPublicIPAddressByNameE gets the IP of the specified Public IP Address.
//
// Deprecated: Use [GetIPOfPublicIPAddressByNameContextE] instead.
func GetIPOfPublicIPAddressByNameE(publicAddressName string, resGroupName string, subscriptionID string) (string, error) {
	return GetIPOfPublicIPAddressByNameContextE(context.Background(), publicAddressName, resGroupName, subscriptionID)
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

	return ExtractIPOfPublicIPAddress(pip)
}

// CheckPublicDNSNameAvailability checks whether a domain name in the cloudapp.azure.com zone
// is available for use. This function would fail the test if there is an error.
//
// Deprecated: Use [CheckPublicDNSNameAvailabilityContext] instead.
func CheckPublicDNSNameAvailability(t testing.TestingT, location string, domainNameLabel string, subscriptionID string) bool {
	t.Helper()

	return CheckPublicDNSNameAvailabilityContext(t, context.Background(), location, domainNameLabel, subscriptionID)
}

// CheckPublicDNSNameAvailabilityE checks whether a domain name in the cloudapp.azure.com zone
// is available for use.
//
// Deprecated: Use [CheckPublicDNSNameAvailabilityContextE] instead.
func CheckPublicDNSNameAvailabilityE(location string, domainNameLabel string, subscriptionID string) (bool, error) {
	return CheckPublicDNSNameAvailabilityContextE(context.Background(), location, domainNameLabel, subscriptionID)
}

// CheckPublicDNSNameAvailabilityContext checks whether a domain name in the cloudapp.azure.com zone
// is available for use. This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func CheckPublicDNSNameAvailabilityContext(t testing.TestingT, ctx context.Context, location string, domainNameLabel string, subscriptionID string) bool {
	t.Helper()

	available, err := CheckPublicDNSNameAvailabilityContextE(ctx, location, domainNameLabel, subscriptionID)
	require.NoError(t, err)

	return available
}

// CheckPublicDNSNameAvailabilityContextE checks whether a domain name in the cloudapp.azure.com zone
// is available for use.
// The ctx parameter supports cancellation and timeouts.
func CheckPublicDNSNameAvailabilityContextE(ctx context.Context, location string, domainNameLabel string, subscriptionID string) (bool, error) {
	client, err := CreateNetworkManagementClientContextE(ctx, subscriptionID)
	if err != nil {
		return false, err
	}

	res, err := client.CheckDNSNameAvailability(ctx, location, domainNameLabel, nil)
	if err != nil {
		return false, err
	}

	if res.Available == nil {
		return false, nil
	}

	return *res.Available, nil
}

// GetPublicIPAddressE gets a Public IP Address in the specified Azure Resource Group.
//
// Deprecated: Use [GetPublicIPAddressContextE] instead.
func GetPublicIPAddressE(publicIPAddressName string, resGroupName string, subscriptionID string) (*armnetwork.PublicIPAddress, error) {
	return GetPublicIPAddressContextE(context.Background(), publicIPAddressName, resGroupName, subscriptionID)
}

// GetPublicIPAddressContextE gets a Public IP Address in the specified Azure Resource Group.
// The ctx parameter supports cancellation and timeouts.
func GetPublicIPAddressContextE(ctx context.Context, publicIPAddressName string, resGroupName string, subscriptionID string) (*armnetwork.PublicIPAddress, error) {
	resGroupName, err := getTargetAzureResourceGroupName(resGroupName)
	if err != nil {
		return nil, err
	}

	client, err := GetPublicIPAddressClientContextE(ctx, subscriptionID)
	if err != nil {
		return nil, err
	}

	return GetPublicIPAddressWithClient(ctx, client, resGroupName, publicIPAddressName)
}

// GetPublicIPAddressWithClient gets a Public IP Address using the provided PublicIPAddressesClient.
func GetPublicIPAddressWithClient(ctx context.Context, client *armnetwork.PublicIPAddressesClient, resGroupName string, publicIPAddressName string) (*armnetwork.PublicIPAddress, error) {
	resp, err := client.Get(ctx, resGroupName, publicIPAddressName, nil)
	if err != nil {
		return nil, err
	}

	return &resp.PublicIPAddress, nil
}

// ExtractIPOfPublicIPAddress gets the IP string from a PublicIPAddress.
func ExtractIPOfPublicIPAddress(pip *armnetwork.PublicIPAddress) (string, error) {
	if pip == nil {
		return "", errors.New("public IP address is nil")
	}

	if pip.Properties == nil || pip.Properties.IPAddress == nil {
		name := "<unknown>"
		if pip.Name != nil {
			name = *pip.Name
		}

		return "", fmt.Errorf("public IP address %q has no IP address assigned", name)
	}

	return *pip.Properties.IPAddress, nil
}

// GetPublicIPAddressClientContextE creates a Public IP Addresses client in the specified Azure Subscription.
// The ctx parameter supports cancellation and timeouts.
func GetPublicIPAddressClientContextE(ctx context.Context, subscriptionID string) (*armnetwork.PublicIPAddressesClient, error) {
	return CreatePublicIPAddressesClientContextE(ctx, subscriptionID)
}

// GetPublicIPAddressClientE creates a Public IP Addresses client in the specified Azure Subscription.
//
// Deprecated: Use [GetPublicIPAddressClientContextE] instead.
func GetPublicIPAddressClientE(subscriptionID string) (*armnetwork.PublicIPAddressesClient, error) {
	return GetPublicIPAddressClientContextE(context.Background(), subscriptionID)
}
