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

	subID, err := getTargetAzureSubscription(subscriptionID)
	if err != nil {
		return nil, err
	}

	cred, err := newArmCredential()
	if err != nil {
		return nil, err
	}

	opts, err := newArmClientOptions()
	if err != nil {
		return nil, err
	}

	client, err := armprivatedns.NewPrivateZonesClient(subID, cred, opts)
	if err != nil {
		return nil, err
	}

	return GetPrivateDNSZoneWithClient(ctx, client, rgName, zoneName)
}

// GetPrivateDNSZoneWithClient gets the specified private DNS zone using the provided PrivateZonesClient.
// This variant is useful for testing with fake clients.
func GetPrivateDNSZoneWithClient(ctx context.Context, client *armprivatedns.PrivateZonesClient, resGroupName string, zoneName string) (*armprivatedns.PrivateZone, error) {
	resp, err := client.Get(ctx, resGroupName, zoneName, nil)
	if err != nil {
		return nil, err
	}

	return &resp.PrivateZone, nil
}
