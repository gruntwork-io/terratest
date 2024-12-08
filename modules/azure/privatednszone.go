package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/privatedns/mgmt/privatedns"
)

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
