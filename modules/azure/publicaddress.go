package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2019-09-01/network"
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// PublicAddressExists indicates whether the specified AzurePublic Address exists.
// This function would fail the test if there is an error.
func PublicAddressExists(t testing.TestingT, publicAddressName string, resGroupName string, subscriptionID string) bool {
	exists, err := PublicAddressExistsE(publicAddressName, resGroupName, subscriptionID)
	require.NoError(t, err)
	return exists
}

// PublicAddressExistsE indicates whether the specified AzurePublic Address exists.
func PublicAddressExistsE(publicAddressName string, resGroupName string, subscriptionID string) (bool, error) {
	// Get the Public Address
	_, err := GetPublicIPAddressE(publicAddressName, resGroupName, subscriptionID)
	if err != nil {
		if ResourceNotFoundErrorExists(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// GetIPOfPublicIPAddressByName gets the Public IP of the Public IP Address specified.
// This function would fail the test if there is an error.
func GetIPOfPublicIPAddressByName(t testing.TestingT, publicAddressName string, resGroupName string, subscriptionID string) string {
	IP, err := GetIPOfPublicIPAddressByNameE(publicAddressName, resGroupName, subscriptionID)
	require.NoError(t, err)
	return IP
}

// GetIPOfPublicIPAddressByNameE gets the Public IP of the Public IP Address specified.
func GetIPOfPublicIPAddressByNameE(publicAddressName string, resGroupName string, subscriptionID string) (string, error) {
	// Create a NIC client
	pip, err := GetPublicIPAddressE(publicAddressName, resGroupName, subscriptionID)
	if err != nil {
		return "", err
	}

	return *pip.IPAddress, nil
}

// CheckPublicDNSNameAvailability checks whether a Domain Name in the cloudapp.azure.com zone
// is available for use. This function would fail the test if there is an error.
func CheckPublicDNSNameAvailability(t testing.TestingT, location string, domainNameLabel string, subscriptionID string) bool {
	available, err := CheckPublicDNSNameAvailabilityE(location, domainNameLabel, subscriptionID)
	if err != nil {
		return false
	}
	return available
}

// CheckPublicDNSNameAvailabilityE checks whether a Domain Name in the cloudapp.azure.com zone is available for use.
func CheckPublicDNSNameAvailabilityE(location string, domainNameLabel string, subscriptionID string) (bool, error) {
	client, err := GetPublicIPAddressClientE(subscriptionID)
	if err != nil {
		return false, err
	}

	res, err := client.CheckDNSNameAvailability(context.Background(), location, domainNameLabel)
	if err != nil {
		return false, err
	}

	return *res.Available, nil
}

// GetPublicIPAddressE gets a Public IP Addresses in the specified Azure Resource Group.
func GetPublicIPAddressE(publicIPAddressName string, resGroupName string, subscriptionID string) (*network.PublicIPAddress, error) {
	// Validate resource group name and subscription ID
	resGroupName, err := getTargetAzureResourceGroupName(resGroupName)
	if err != nil {
		return nil, err
	}

	// Get the client reference
	client, err := GetPublicIPAddressClientE(subscriptionID)
	if err != nil {
		return nil, err
	}

	// Get the Public IP Address
	pip, err := client.Get(context.Background(), resGroupName, publicIPAddressName, "")
	if err != nil {
		return nil, err
	}
	return &pip, nil
}

// GetPublicIPAddressClientE creates a Public IP Addresses client in the specified Azure Subscription.
func GetPublicIPAddressClientE(subscriptionID string) (*network.PublicIPAddressesClient, error) {
	// Get the Public IP Address client from clientfactory
	client, err := CreatePublicIPAddressesClientE(subscriptionID)
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
