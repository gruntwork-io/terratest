package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2019-09-01/network"
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

// PublicAddressExists indicates whether the specified Azure Public Address exists.
// This function would fail the test if there is an error.
//
// Deprecated: Use [PublicAddressExistsContext] instead.
func PublicAddressExists(t testing.TestingT, publicAddressName string, resGroupName string, subscriptionID string) bool {
	t.Helper()

	return PublicAddressExistsContext(t, context.Background(), publicAddressName, resGroupName, subscriptionID)
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

// PublicAddressExistsE indicates whether the specified Azure Public Address exists.
//
// Deprecated: Use [PublicAddressExistsContextE] instead.
func PublicAddressExistsE(publicAddressName string, resGroupName string, subscriptionID string) (bool, error) {
	return PublicAddressExistsContextE(context.Background(), publicAddressName, resGroupName, subscriptionID)
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

// GetIPOfPublicIPAddressByName gets the Public IP of the Public IP Address specified.
// This function would fail the test if there is an error.
//
// Deprecated: Use [GetIPOfPublicIPAddressByNameContext] instead.
func GetIPOfPublicIPAddressByName(t testing.TestingT, publicAddressName string, resGroupName string, subscriptionID string) string {
	t.Helper()

	return GetIPOfPublicIPAddressByNameContext(t, context.Background(), publicAddressName, resGroupName, subscriptionID)
}

// GetIPOfPublicIPAddressByNameContextE gets the IP of the specified Public IP Address.
// The ctx parameter supports cancellation and timeouts.
func GetIPOfPublicIPAddressByNameContextE(ctx context.Context, publicAddressName string, resGroupName string, subscriptionID string) (string, error) {
	pip, err := GetPublicIPAddressContextE(ctx, publicAddressName, resGroupName, subscriptionID)
	if err != nil {
		return "", err
	}

	return *pip.IPAddress, nil
}

// GetIPOfPublicIPAddressByNameE gets the Public IP of the Public IP Address specified.
//
// Deprecated: Use [GetIPOfPublicIPAddressByNameContextE] instead.
func GetIPOfPublicIPAddressByNameE(publicAddressName string, resGroupName string, subscriptionID string) (string, error) {
	return GetIPOfPublicIPAddressByNameContextE(context.Background(), publicAddressName, resGroupName, subscriptionID)
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

// CheckPublicDNSNameAvailability checks whether a Domain Name in the cloudapp.azure.com zone
// is available for use. This function would fail the test if there is an error.
//
// Deprecated: Use [CheckPublicDNSNameAvailabilityContext] instead.
func CheckPublicDNSNameAvailability(t testing.TestingT, location string, domainNameLabel string, subscriptionID string) bool {
	t.Helper()

	return CheckPublicDNSNameAvailabilityContext(t, context.Background(), location, domainNameLabel, subscriptionID)
}

// CheckPublicDNSNameAvailabilityContextE checks whether a domain name in the cloudapp.azure.com zone
// is available for use.
// The ctx parameter supports cancellation and timeouts.
func CheckPublicDNSNameAvailabilityContextE(ctx context.Context, location string, domainNameLabel string, subscriptionID string) (bool, error) {
	client, err := GetPublicIPAddressClientE(subscriptionID)
	if err != nil {
		return false, err
	}

	res, err := client.CheckDNSNameAvailability(ctx, location, domainNameLabel)
	if err != nil {
		return false, err
	}

	return *res.Available, nil
}

// CheckPublicDNSNameAvailabilityE checks whether a Domain Name in the cloudapp.azure.com zone is available for use.
//
// Deprecated: Use [CheckPublicDNSNameAvailabilityContextE] instead.
func CheckPublicDNSNameAvailabilityE(location string, domainNameLabel string, subscriptionID string) (bool, error) {
	return CheckPublicDNSNameAvailabilityContextE(context.Background(), location, domainNameLabel, subscriptionID)
}

// GetPublicIPAddressContextE gets a Public IP Address in the specified Azure Resource Group.
// The ctx parameter supports cancellation and timeouts.
func GetPublicIPAddressContextE(ctx context.Context, publicIPAddressName string, resGroupName string, subscriptionID string) (*network.PublicIPAddress, error) {
	resGroupName, err := getTargetAzureResourceGroupName(resGroupName)
	if err != nil {
		return nil, err
	}

	client, err := GetPublicIPAddressClientE(subscriptionID)
	if err != nil {
		return nil, err
	}

	pip, err := client.Get(ctx, resGroupName, publicIPAddressName, "")
	if err != nil {
		return nil, err
	}

	return &pip, nil
}

// GetPublicIPAddressE gets a Public IP Addresses in the specified Azure Resource Group.
//
// Deprecated: Use [GetPublicIPAddressContextE] instead.
func GetPublicIPAddressE(publicIPAddressName string, resGroupName string, subscriptionID string) (*network.PublicIPAddress, error) {
	return GetPublicIPAddressContextE(context.Background(), publicIPAddressName, resGroupName, subscriptionID)
}

// GetPublicIPAddressClientE creates a Public IP Addresses client in the specified Azure Subscription.
func GetPublicIPAddressClientE(subscriptionID string) (*network.PublicIPAddressesClient, error) {
	client, err := CreatePublicIPAddressesClientE(subscriptionID)
	if err != nil {
		return nil, err
	}

	authorizer, err := NewAuthorizer()
	if err != nil {
		return nil, err
	}

	client.Authorizer = *authorizer

	return client, nil
}
