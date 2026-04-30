package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/cdn/armcdn"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/frontdoor/armfrontdoor"
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// FrontDoorExistsContext indicates whether the Front Door exists for the subscription.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
//
// Note: this helper targets the legacy classic Front Door API (`armfrontdoor`), which Microsoft
// deprecated for new resource creation on April 1, 2025. For new deployments use the CDN
// Front Door helpers (e.g. [CDNFrontDoorProfileExistsContext]). The legacy helpers remain
// for callers operating against existing classic Front Door resources during the deprecation
// window (full retirement is March 31, 2027).
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

// ---------------------------------------------------------------------------------------------------------------------
// CDN Front Door (modern API) helpers
//
// These helpers wrap the `armcdn` SDK and target the modern CDN Front Door resource family
// (`azurerm_cdn_frontdoor_profile`, `azurerm_cdn_frontdoor_endpoint`, ...) that supersedes
// the classic `azurerm_frontdoor` resource. Profiles are referred to interchangeably as
// "Front Door profiles" by the AzureRM provider; the underlying ARM resource type is
// `Microsoft.Cdn/profiles` with an AFD-specific SKU.
// ---------------------------------------------------------------------------------------------------------------------

// CDNFrontDoorProfileExistsContext indicates whether the specified CDN Front Door profile exists.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func CDNFrontDoorProfileExistsContext(t testing.TestingT, ctx context.Context, profileName string, resourceGroupName string, subscriptionID string) bool {
	t.Helper()

	exists, err := CDNFrontDoorProfileExistsContextE(ctx, profileName, resourceGroupName, subscriptionID)
	require.NoError(t, err)

	return exists
}

// CDNFrontDoorProfileExistsContextE indicates whether the specified CDN Front Door profile exists.
// The ctx parameter supports cancellation and timeouts.
func CDNFrontDoorProfileExistsContextE(ctx context.Context, profileName string, resourceGroupName string, subscriptionID string) (bool, error) {
	_, err := GetCDNFrontDoorProfileContextE(ctx, profileName, resourceGroupName, subscriptionID)
	if err != nil {
		if ResourceNotFoundErrorExists(err) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

// GetCDNFrontDoorProfileContext returns the specified CDN Front Door profile.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetCDNFrontDoorProfileContext(t testing.TestingT, ctx context.Context, profileName string, resourceGroupName string, subscriptionID string) *armcdn.Profile {
	t.Helper()

	profile, err := GetCDNFrontDoorProfileContextE(ctx, profileName, resourceGroupName, subscriptionID)
	require.NoError(t, err)

	return profile
}

// GetCDNFrontDoorProfileContextE returns the specified CDN Front Door profile if it exists.
// The ctx parameter supports cancellation and timeouts.
func GetCDNFrontDoorProfileContextE(ctx context.Context, profileName, resourceGroupName, subscriptionID string) (*armcdn.Profile, error) {
	client, err := CreateCDNProfilesClientContextE(ctx, subscriptionID)
	if err != nil {
		return nil, err
	}

	return GetCDNFrontDoorProfileWithClient(ctx, client, resourceGroupName, profileName)
}

// GetCDNFrontDoorProfileWithClient gets the specified CDN Front Door profile using the provided client.
// This variant is useful for testing with fake clients.
func GetCDNFrontDoorProfileWithClient(ctx context.Context, client *armcdn.ProfilesClient, resourceGroupName, profileName string) (*armcdn.Profile, error) {
	resp, err := client.Get(ctx, resourceGroupName, profileName, nil)
	if err != nil {
		return nil, err
	}

	return &resp.Profile, nil
}

// CDNFrontDoorEndpointExistsContext indicates whether the specified CDN Front Door endpoint exists for the given profile.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func CDNFrontDoorEndpointExistsContext(t testing.TestingT, ctx context.Context, endpointName string, profileName string, resourceGroupName string, subscriptionID string) bool {
	t.Helper()

	exists, err := CDNFrontDoorEndpointExistsContextE(ctx, endpointName, profileName, resourceGroupName, subscriptionID)
	require.NoError(t, err)

	return exists
}

// CDNFrontDoorEndpointExistsContextE indicates whether the specified CDN Front Door endpoint exists for the given profile.
// The ctx parameter supports cancellation and timeouts.
func CDNFrontDoorEndpointExistsContextE(ctx context.Context, endpointName, profileName, resourceGroupName, subscriptionID string) (bool, error) {
	_, err := GetCDNFrontDoorEndpointContextE(ctx, endpointName, profileName, resourceGroupName, subscriptionID)
	if err != nil {
		if ResourceNotFoundErrorExists(err) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

// GetCDNFrontDoorEndpointContext returns the specified CDN Front Door endpoint for the given profile.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetCDNFrontDoorEndpointContext(t testing.TestingT, ctx context.Context, endpointName string, profileName string, resourceGroupName string, subscriptionID string) *armcdn.AFDEndpoint {
	t.Helper()

	ep, err := GetCDNFrontDoorEndpointContextE(ctx, endpointName, profileName, resourceGroupName, subscriptionID)
	require.NoError(t, err)

	return ep
}

// GetCDNFrontDoorEndpointContextE returns the specified CDN Front Door endpoint for the given profile if it exists.
// The ctx parameter supports cancellation and timeouts.
func GetCDNFrontDoorEndpointContextE(ctx context.Context, endpointName, profileName, resourceGroupName, subscriptionID string) (*armcdn.AFDEndpoint, error) {
	client, err := CreateCDNAFDEndpointsClientContextE(ctx, subscriptionID)
	if err != nil {
		return nil, err
	}

	return GetCDNFrontDoorEndpointWithClient(ctx, client, resourceGroupName, profileName, endpointName)
}

// GetCDNFrontDoorEndpointWithClient gets the specified CDN Front Door endpoint using the provided client.
// This variant is useful for testing with fake clients.
func GetCDNFrontDoorEndpointWithClient(ctx context.Context, client *armcdn.AFDEndpointsClient, resourceGroupName, profileName, endpointName string) (*armcdn.AFDEndpoint, error) {
	resp, err := client.Get(ctx, resourceGroupName, profileName, endpointName, nil)
	if err != nil {
		return nil, err
	}

	return &resp.AFDEndpoint, nil
}
