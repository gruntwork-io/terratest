package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/privatedns/mgmt/privatedns"
)

// PrivateDNSZoneExistsContextE indicates whether the specified private DNS zone exists.
// The ctx parameter supports cancellation and timeouts.
func PrivateDNSZoneExistsContextE(ctx context.Context, zoneName string, resourceGroupName string, subscriptionID string) (bool, error) {
	_, err := GetPrivateDNSZoneContextE(ctx, zoneName, resourceGroupName, subscriptionID)
	if err != nil {
		if ResourceNotFoundErrorExists(err) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

// PrivateDNSZoneExistsE indicates whether the specified private DNS zone exists.
//
// Deprecated: Use [PrivateDNSZoneExistsContextE] instead.
func PrivateDNSZoneExistsE(zoneName string, resourceGroupName string, subscriptionID string) (bool, error) {
	return PrivateDNSZoneExistsContextE(context.Background(), zoneName, resourceGroupName, subscriptionID)
}

// GetPrivateDNSZoneContextE gets the specified private DNS zone object.
// The ctx parameter supports cancellation and timeouts.
func GetPrivateDNSZoneContextE(ctx context.Context, zoneName string, resGroupName string, subscriptionID string) (*privatedns.PrivateZone, error) {
	rgName, err := getTargetAzureResourceGroupName(resGroupName)
	if err != nil {
		return nil, err
	}

	client, err := CreatePrivateDNSZonesClientE(subscriptionID)
	if err != nil {
		return nil, err
	}

	zone, err := client.Get(ctx, rgName, zoneName)
	if err != nil {
		return nil, err
	}

	return &zone, nil
}

// GetPrivateDNSZoneE gets the private DNS zone object.
//
// Deprecated: Use [GetPrivateDNSZoneContextE] instead.
func GetPrivateDNSZoneE(zoneName string, resGroupName string, subscriptionID string) (*privatedns.PrivateZone, error) {
	return GetPrivateDNSZoneContextE(context.Background(), zoneName, resGroupName, subscriptionID)
}
