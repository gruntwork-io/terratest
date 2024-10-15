package azure

import (
	"context"
	"testing"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/privatedns/mgmt/privatedns"
	"github.com/stretchr/testify/require"
)

// PrivateDNSZoneExists indicates whether the specified private DNS zone exists.
// This function would fail the test if there is an error.
func PrivateDNSZoneExists(t *testing.T, zoneName string, resourceGroupName string, subscriptionID string) bool {
	exists, err := PrivateDNSZoneExistsE(zoneName, resourceGroupName, subscriptionID)
	require.NoError(t, err)

	return exists
}

// PrivateDNSZoneExistsE indicates whether the specified private DNS zone exists.
func PrivateDNSZoneExistsE(zoneName string, resourceGroupName string, subscriptionID string) (bool, error) {
	_, err := GetPrivateDNSZoneE(zoneName, resourceGroupName, subscriptionID)
	if err != nil {
		if ResourceNotFoundErrorExists(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// GetPrivateDNSZone gets the private DNS zone object
// This function would fail the test if there is an error.
func GetPrivateDNSZone(t *testing.T, zoneName string, resGroupName string, subscriptionID string) *privatedns.PrivateZone {
	zone, err := GetPrivateDNSZoneE(zoneName, resGroupName, subscriptionID)

	require.NoError(t, err)

	return zone
}

// GetPrivateDNSZoneE gets the private DNS zone object
func GetPrivateDNSZoneE(zoneName string, resGroupName string, subscriptionID string) (*privatedns.PrivateZone, error) {
	rgName, err := getTargetAzureResourceGroupName(resGroupName)
	if err != nil {
		return nil, err
	}

	client, err := CreatePrivateDnsZonesClientE(subscriptionID)
	if err != nil {
		return nil, err
	}

	zone, err := client.Get(context.Background(), rgName, zoneName)
	if err != nil {
		return nil, err
	}

	return &zone, nil
}
