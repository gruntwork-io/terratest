package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/frontdoor/armfrontdoor"
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// FrontDoorExistsContext indicates whether the Front Door exists for the subscription.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func FrontDoorExistsContext(t testing.TestingT, ctx context.Context, frontDoorName string, resourceGroupName string, subscriptionID string) bool {
	t.Helper()

	exists, err := FrontDoorExistsContextE(ctx, frontDoorName, resourceGroupName, subscriptionID)
	require.NoError(t, err)

	return exists
}

// GetFrontDoorContext gets a Front Door by name if it exists for the subscription.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetFrontDoorContext(t testing.TestingT, ctx context.Context, frontDoorName string, resourceGroupName string, subscriptionID string) *armfrontdoor.FrontDoor {
	t.Helper()

	fd, err := GetFrontDoorContextE(ctx, frontDoorName, resourceGroupName, subscriptionID)
	require.NoError(t, err)

	return fd
}

// FrontDoorFrontendEndpointExistsContext indicates whether the frontend endpoint exists for the provided Front Door.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func FrontDoorFrontendEndpointExistsContext(t testing.TestingT, ctx context.Context, endpointName string, frontDoorName string, resourceGroupName string, subscriptionID string) bool {
	t.Helper()

	exists, err := FrontDoorFrontendEndpointExistsContextE(ctx, endpointName, frontDoorName, resourceGroupName, subscriptionID)
	require.NoError(t, err)

	return exists
}

// GetFrontDoorFrontendEndpointContext gets a frontend endpoint by name for the provided Front Door if it exists for the subscription.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetFrontDoorFrontendEndpointContext(t testing.TestingT, ctx context.Context, endpointName string, frontDoorName string, resourceGroupName string, subscriptionID string) *armfrontdoor.FrontendEndpoint {
	t.Helper()

	ep, err := GetFrontDoorFrontendEndpointContextE(ctx, endpointName, frontDoorName, resourceGroupName, subscriptionID)
	require.NoError(t, err)

	return ep
}

// FrontDoorExistsContextE indicates whether the specified Front Door exists.
// The ctx parameter supports cancellation and timeouts.
func FrontDoorExistsContextE(ctx context.Context, frontDoorName string, resourceGroupName string, subscriptionID string) (bool, error) {
	_, err := GetFrontDoorContextE(ctx, frontDoorName, resourceGroupName, subscriptionID)
	if err != nil {
		if ResourceNotFoundErrorExists(err) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

// FrontDoorFrontendEndpointExistsContextE indicates whether the specified endpoint exists for the provided Front Door.
// The ctx parameter supports cancellation and timeouts.
func FrontDoorFrontendEndpointExistsContextE(ctx context.Context, endpointName string, frontDoorName string, resourceGroupName string, subscriptionID string) (bool, error) {
	_, err := GetFrontDoorFrontendEndpointContextE(ctx, endpointName, frontDoorName, resourceGroupName, subscriptionID)
	if err != nil {
		if ResourceNotFoundErrorExists(err) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

// GetFrontDoorContextE gets the specified Front Door if it exists.
// The ctx parameter supports cancellation and timeouts.
func GetFrontDoorContextE(ctx context.Context, frontDoorName, resourceGroupName, subscriptionID string) (*armfrontdoor.FrontDoor, error) {
	client, err := GetFrontDoorClientE(subscriptionID)
	if err != nil {
		return nil, err
	}

	resp, err := client.Get(ctx, resourceGroupName, frontDoorName, nil)
	if err != nil {
		return nil, err
	}

	return &resp.FrontDoor, nil
}

// GetFrontDoorFrontendEndpointContextE gets the specified Frontend Endpoint for the provided Front Door if it exists.
// The ctx parameter supports cancellation and timeouts.
func GetFrontDoorFrontendEndpointContextE(ctx context.Context, endpointName, frontDoorName, resourceGroupName, subscriptionID string) (*armfrontdoor.FrontendEndpoint, error) {
	client, err := GetFrontDoorFrontendEndpointClientE(subscriptionID)
	if err != nil {
		return nil, err
	}

	resp, err := client.Get(ctx, resourceGroupName, frontDoorName, endpointName, nil)
	if err != nil {
		return nil, err
	}

	return &resp.FrontendEndpoint, nil
}

// GetFrontDoorClientE returns a front door client; otherwise error.
func GetFrontDoorClientE(subscriptionID string) (*armfrontdoor.FrontDoorsClient, error) {
	return CreateFrontDoorClientE(subscriptionID)
}

// GetFrontDoorFrontendEndpointClientE returns a front door frontend endpoints client; otherwise error.
func GetFrontDoorFrontendEndpointClientE(subscriptionID string) (*armfrontdoor.FrontendEndpointsClient, error) {
	return CreateFrontDoorFrontendEndpointClientE(subscriptionID)
}
