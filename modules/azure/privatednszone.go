package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/privatedns/armprivatedns"
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

// GetPrivateDNSZoneContextE gets the specified private DNS zone object.
// The ctx parameter supports cancellation and timeouts.
func GetPrivateDNSZoneContextE(ctx context.Context, zoneName string, resGroupName string, subscriptionID string) (*armprivatedns.PrivateZone, error) {
	rgName, err := getTargetAzureResourceGroupName(resGroupName)
	if err != nil {
		return nil, err
	}

	client, err := CreatePrivateDNSZonesClientE(subscriptionID)
	if err != nil {
		return nil, err
	}

	resp, err := client.Get(ctx, rgName, zoneName, nil)
	if err != nil {
		return nil, err
	}

	return &resp.PrivateZone, nil
}
