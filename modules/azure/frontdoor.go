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

// FrontDoorExists indicates whether the Front Door exists for the subscription.
// This function would fail the test if there is an error.
//
// Deprecated: Use [FrontDoorExistsContext] instead.
func FrontDoorExists(t testing.TestingT, frontDoorName string, resourceGroupName string, subscriptionID string) bool {
	t.Helper()

	return FrontDoorExistsContext(t, context.Background(), frontDoorName, resourceGroupName, subscriptionID)
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

// GetFrontDoor gets a Front Door by name if it exists for the subscription.
// This function would fail the test if there is an error.
//
// Deprecated: Use [GetFrontDoorContext] instead.
func GetFrontDoor(t testing.TestingT, frontDoorName string, resourceGroupName string, subscriptionID string) *armfrontdoor.FrontDoor {
	t.Helper()

	return GetFrontDoorContext(t, context.Background(), frontDoorName, resourceGroupName, subscriptionID)
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

// FrontDoorFrontendEndpointExists indicates whether the frontend endpoint exists for the provided Front Door.
// This function would fail the test if there is an error.
//
// Deprecated: Use [FrontDoorFrontendEndpointExistsContext] instead.
func FrontDoorFrontendEndpointExists(t testing.TestingT, endpointName string, frontDoorName string, resourceGroupName string, subscriptionID string) bool {
	t.Helper()

	return FrontDoorFrontendEndpointExistsContext(t, context.Background(), endpointName, frontDoorName, resourceGroupName, subscriptionID)
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

// GetFrontDoorFrontendEndpoint gets a frontend endpoint by name for the provided Front Door if it exists for the subscription.
// This function would fail the test if there is an error.
//
// Deprecated: Use [GetFrontDoorFrontendEndpointContext] instead.
func GetFrontDoorFrontendEndpoint(t testing.TestingT, endpointName string, frontDoorName string, resourceGroupName string, subscriptionID string) *armfrontdoor.FrontendEndpoint {
	t.Helper()

	return GetFrontDoorFrontendEndpointContext(t, context.Background(), endpointName, frontDoorName, resourceGroupName, subscriptionID)
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

// FrontDoorExistsE indicates whether the specified Front Door exists.
//
// Deprecated: Use [FrontDoorExistsContextE] instead.
func FrontDoorExistsE(frontDoorName string, resourceGroupName string, subscriptionID string) (bool, error) {
	return FrontDoorExistsContextE(context.Background(), frontDoorName, resourceGroupName, subscriptionID)
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

// FrontDoorFrontendEndpointExistsE indicates whether the specified endpoint exists for the provided Front Door.
//
// Deprecated: Use [FrontDoorFrontendEndpointExistsContextE] instead.
func FrontDoorFrontendEndpointExistsE(endpointName string, frontDoorName string, resourceGroupName string, subscriptionID string) (bool, error) {
	return FrontDoorFrontendEndpointExistsContextE(context.Background(), endpointName, frontDoorName, resourceGroupName, subscriptionID)
}

// GetFrontDoorContextE gets the specified Front Door if it exists.
// The ctx parameter supports cancellation and timeouts.
func GetFrontDoorContextE(ctx context.Context, frontDoorName, resourceGroupName, subscriptionID string) (*armfrontdoor.FrontDoor, error) {
	client, err := CreateFrontDoorClientContextE(ctx, subscriptionID)
	if err != nil {
		return nil, err
	}

	return GetFrontDoorWithClient(ctx, client, resourceGroupName, frontDoorName)
}

// GetFrontDoorE gets the specified Front Door if it exists.
//
// Deprecated: Use [GetFrontDoorContextE] instead.
func GetFrontDoorE(frontDoorName, resourceGroupName, subscriptionID string) (*armfrontdoor.FrontDoor, error) {
	return GetFrontDoorContextE(context.Background(), frontDoorName, resourceGroupName, subscriptionID)
}

// GetFrontDoorWithClient gets the specified Front Door using the provided client.
// This variant is useful for testing with fake clients.
func GetFrontDoorWithClient(ctx context.Context, client *armfrontdoor.FrontDoorsClient, resourceGroupName, frontDoorName string) (*armfrontdoor.FrontDoor, error) {
	resp, err := client.Get(ctx, resourceGroupName, frontDoorName, nil)
	if err != nil {
		return nil, err
	}

	return &resp.FrontDoor, nil
}

// GetFrontDoorFrontendEndpointContextE gets the specified Frontend Endpoint for the provided Front Door if it exists.
// The ctx parameter supports cancellation and timeouts.
func GetFrontDoorFrontendEndpointContextE(ctx context.Context, endpointName, frontDoorName, resourceGroupName, subscriptionID string) (*armfrontdoor.FrontendEndpoint, error) {
	client, err := CreateFrontDoorFrontendEndpointClientContextE(ctx, subscriptionID)
	if err != nil {
		return nil, err
	}

	return GetFrontDoorFrontendEndpointWithClient(ctx, client, resourceGroupName, frontDoorName, endpointName)
}

// GetFrontDoorFrontendEndpointE gets the specified Frontend Endpoint for the provided Front Door if it exists.
//
// Deprecated: Use [GetFrontDoorFrontendEndpointContextE] instead.
func GetFrontDoorFrontendEndpointE(endpointName, frontDoorName, resourceGroupName, subscriptionID string) (*armfrontdoor.FrontendEndpoint, error) {
	return GetFrontDoorFrontendEndpointContextE(context.Background(), endpointName, frontDoorName, resourceGroupName, subscriptionID)
}

// GetFrontDoorFrontendEndpointWithClient gets the specified Frontend Endpoint using the provided client.
// This variant is useful for testing with fake clients.
func GetFrontDoorFrontendEndpointWithClient(ctx context.Context, client *armfrontdoor.FrontendEndpointsClient, resourceGroupName, frontDoorName, endpointName string) (*armfrontdoor.FrontendEndpoint, error) {
	resp, err := client.Get(ctx, resourceGroupName, frontDoorName, endpointName, nil)
	if err != nil {
		return nil, err
	}

	return &resp.FrontendEndpoint, nil
}

// GetFrontDoorClientE returns a Front Door client; otherwise error.
//
// Deprecated: Use [CreateFrontDoorClientContextE] instead.
func GetFrontDoorClientE(subscriptionID string) (*armfrontdoor.FrontDoorsClient, error) {
	return CreateFrontDoorClientContextE(context.Background(), subscriptionID)
}

// GetFrontDoorFrontendEndpointClientE returns a Front Door frontend endpoints client; otherwise error.
//
// Deprecated: Use [CreateFrontDoorFrontendEndpointClientContextE] instead.
func GetFrontDoorFrontendEndpointClientE(subscriptionID string) (*armfrontdoor.FrontendEndpointsClient, error) {
	return CreateFrontDoorFrontendEndpointClientContextE(context.Background(), subscriptionID)
}
