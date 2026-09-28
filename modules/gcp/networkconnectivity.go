package gcp

import (
	"context"
	"errors"
	"fmt"

	"github.com/gruntwork-io/terratest/modules/core/v2/logger"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/networkconnectivity/v1"
	"google.golang.org/api/option"
)

// GetNetworkConnectivityHubAttrs returns the settings Google Cloud holds for the given Network
// Connectivity Center hub, so a test can assert on what was actually created rather than only
// that it exists. A hub is global, so no location is needed.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetNetworkConnectivityHubAttrs(t testing.TestingT, ctx context.Context, projectID string, hubName string) *networkconnectivity.Hub {
	hub, err := GetNetworkConnectivityHubAttrsE(t, ctx, projectID, hubName)
	require.NoError(t, err)

	return hub
}

// GetNetworkConnectivityHubAttrsE returns the settings Google Cloud holds for the given Network
// Connectivity Center hub.
// The ctx parameter supports cancellation and timeouts.
func GetNetworkConnectivityHubAttrsE(t testing.TestingT, ctx context.Context, projectID string, hubName string) (*networkconnectivity.Hub, error) {
	logger.Default.Logf(t, "Getting settings for Network Connectivity hub %s in project %s", hubName, projectID)

	service, err := NewNetworkConnectivityServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetNetworkConnectivityHubAttrsWithClient(ctx, service, projectID, hubName)
}

// GetNetworkConnectivityHubAttrsWithClient returns the settings Google Cloud holds for the given
// Network Connectivity Center hub using the supplied *networkconnectivity.Service. Prefer this
// variant in unit tests where the service is backed by an httptest fake server (see
// networkconnectivity_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetNetworkConnectivityHubAttrsWithClient(ctx context.Context, service *networkconnectivity.Service, projectID string, hubName string) (*networkconnectivity.Hub, error) {
	name := fmt.Sprintf("projects/%s/locations/global/hubs/%s", projectID, hubName)

	hub, err := service.Projects.Locations.Global.Hubs.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("Network Connectivity hub %s does not exist in project %s", hubName, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Network Connectivity hub %s in project %s: %w", hubName, projectID, err)
	}

	return hub, nil
}

// NewNetworkConnectivityServiceE creates a Network Connectivity service authenticated the same way
// every other client in this module is.
// The ctx parameter supports cancellation and timeouts.
func NewNetworkConnectivityServiceE(t testing.TestingT, ctx context.Context) (*networkconnectivity.Service, error) {
	return networkconnectivity.NewService(ctx, append(withOptions(), option.WithScopes(networkconnectivity.CloudPlatformScope))...)
}

// GetNetworkConnectivitySpokeAttrs returns the settings Google Cloud holds for the given Network
// Connectivity Center spoke, so a test can assert on what was actually created rather than only
// that it exists. A spoke lives in a location, unlike the hub it attaches to, which is global.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetNetworkConnectivitySpokeAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, spokeID string) *networkconnectivity.Spoke {
	spoke, err := GetNetworkConnectivitySpokeAttrsE(t, ctx, projectID, location, spokeID)
	require.NoError(t, err)

	return spoke
}

// GetNetworkConnectivitySpokeAttrsE returns the settings Google Cloud holds for the given Network
// Connectivity Center spoke.
// The ctx parameter supports cancellation and timeouts.
func GetNetworkConnectivitySpokeAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, spokeID string) (*networkconnectivity.Spoke, error) {
	logger.Default.Logf(t, "Getting settings for spoke %s in %s in project %s", spokeID, location, projectID)

	service, err := NewNetworkConnectivityServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetNetworkConnectivitySpokeAttrsWithClient(ctx, service, projectID, location, spokeID)
}

// GetNetworkConnectivitySpokeAttrsWithClient returns the settings Google Cloud holds for the given
// spoke using the supplied *networkconnectivity.Service. Prefer this variant in unit tests where
// the service is backed by an httptest fake server (see networkconnectivity_test.go for the
// pattern).
// The ctx parameter supports cancellation and timeouts.
func GetNetworkConnectivitySpokeAttrsWithClient(ctx context.Context, service *networkconnectivity.Service, projectID string, location string, spokeID string) (*networkconnectivity.Spoke, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/spokes/%s", projectID, location, spokeID)

	spoke, err := service.Projects.Locations.Spokes.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the spoke %s does not exist in %s in project %s", spokeID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for spoke %s in %s in project %s: %w", spokeID, location, projectID, err)
	}

	return spoke, nil
}

// GetPolicyBasedRouteAttrs returns the settings Google Cloud holds for the given policy based
// route, so a test can assert on what was actually created rather than only that it exists. A
// policy based route is global, so no location is needed.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetPolicyBasedRouteAttrs(t testing.TestingT, ctx context.Context, projectID string, routeID string) *networkconnectivity.PolicyBasedRoute {
	route, err := GetPolicyBasedRouteAttrsE(t, ctx, projectID, routeID)
	require.NoError(t, err)

	return route
}

// GetPolicyBasedRouteAttrsE returns the settings Google Cloud holds for the given policy based
// route.
// The ctx parameter supports cancellation and timeouts.
func GetPolicyBasedRouteAttrsE(t testing.TestingT, ctx context.Context, projectID string, routeID string) (*networkconnectivity.PolicyBasedRoute, error) {
	logger.Default.Logf(t, "Getting settings for policy based route %s in project %s", routeID, projectID)

	service, err := NewNetworkConnectivityServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetPolicyBasedRouteAttrsWithClient(ctx, service, projectID, routeID)
}

// GetPolicyBasedRouteAttrsWithClient returns the settings Google Cloud holds for the given policy
// based route using the supplied *networkconnectivity.Service. Prefer this variant in unit tests
// where the service is backed by an httptest fake server (see networkconnectivity_test.go for the
// pattern).
// The ctx parameter supports cancellation and timeouts.
func GetPolicyBasedRouteAttrsWithClient(ctx context.Context, service *networkconnectivity.Service, projectID string, routeID string) (*networkconnectivity.PolicyBasedRoute, error) {
	name := fmt.Sprintf("projects/%s/locations/global/policyBasedRoutes/%s", projectID, routeID)

	route, err := service.Projects.Locations.Global.PolicyBasedRoutes.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the policy based route %s does not exist in project %s", routeID, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for policy based route %s in project %s: %w", routeID, projectID, err)
	}

	return route, nil
}
