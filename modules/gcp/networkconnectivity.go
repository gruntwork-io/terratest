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

// GetNetworkConnectivityMulticloudDataTransferConfigAttrs returns the settings Google Cloud holds for the given multicloud data transfer config, so a test can assert on what was
// actually created rather than only that it exists.
// The config is what enables charged transfer to a named cloud, so which services it covers is what it bills for.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetNetworkConnectivityMulticloudDataTransferConfigAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, id string) *networkconnectivity.MulticloudDataTransferConfig {
	attrs, err := GetNetworkConnectivityMulticloudDataTransferConfigAttrsE(t, ctx, projectID, location, id)
	require.NoError(t, err)

	return attrs
}

// GetNetworkConnectivityMulticloudDataTransferConfigAttrsE returns the settings Google Cloud holds for the given multicloud data transfer config.
// The ctx parameter supports cancellation and timeouts.
func GetNetworkConnectivityMulticloudDataTransferConfigAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, id string) (*networkconnectivity.MulticloudDataTransferConfig, error) {
	logger.Default.Logf(t, "Getting settings for multicloud data transfer config %s in %s in project %s", id, location, projectID)

	service, err := NewNetworkConnectivityServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetNetworkConnectivityMulticloudDataTransferConfigAttrsWithClient(ctx, service, projectID, location, id)
}

// GetNetworkConnectivityMulticloudDataTransferConfigAttrsWithClient returns the settings Google Cloud holds for the given multicloud data transfer config using the supplied
// *networkconnectivity.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see networkconnectivity_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetNetworkConnectivityMulticloudDataTransferConfigAttrsWithClient(ctx context.Context, service *networkconnectivity.Service, projectID string, location string, id string) (*networkconnectivity.MulticloudDataTransferConfig, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/multicloudDataTransferConfigs/%s", projectID, location, id)

	attrs, err := service.Projects.Locations.MulticloudDataTransferConfigs.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the multicloud data transfer config %s in %s in project %s does not exist", id, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for multicloud data transfer config %s in %s in project %s: %w", id, location, projectID, err)
	}

	return attrs, nil
}

// GetNetworkConnectivityRegionalEndpointAttrs returns the settings Google Cloud holds for the given regional endpoint, so a test can assert on what was
// actually created rather than only that it exists.
// A regional endpoint is the address a Google API is reached on from inside one network, so the access type and the address are the point.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetNetworkConnectivityRegionalEndpointAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, id string) *networkconnectivity.RegionalEndpoint {
	attrs, err := GetNetworkConnectivityRegionalEndpointAttrsE(t, ctx, projectID, location, id)
	require.NoError(t, err)

	return attrs
}

// GetNetworkConnectivityRegionalEndpointAttrsE returns the settings Google Cloud holds for the given regional endpoint.
// The ctx parameter supports cancellation and timeouts.
func GetNetworkConnectivityRegionalEndpointAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, id string) (*networkconnectivity.RegionalEndpoint, error) {
	logger.Default.Logf(t, "Getting settings for regional endpoint %s in %s in project %s", id, location, projectID)

	service, err := NewNetworkConnectivityServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetNetworkConnectivityRegionalEndpointAttrsWithClient(ctx, service, projectID, location, id)
}

// GetNetworkConnectivityRegionalEndpointAttrsWithClient returns the settings Google Cloud holds for the given regional endpoint using the supplied
// *networkconnectivity.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see networkconnectivity_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetNetworkConnectivityRegionalEndpointAttrsWithClient(ctx context.Context, service *networkconnectivity.Service, projectID string, location string, id string) (*networkconnectivity.RegionalEndpoint, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/regionalEndpoints/%s", projectID, location, id)

	attrs, err := service.Projects.Locations.RegionalEndpoints.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the regional endpoint %s in %s in project %s does not exist", id, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for regional endpoint %s in %s in project %s: %w", id, location, projectID, err)
	}

	return attrs, nil
}

// GetNetworkConnectivityServiceConnectionPolicyAttrs returns the settings Google Cloud holds for the given service connection policy, so a test can assert on what was
// actually created rather than only that it exists.
// The policy says which subnetworks a producer may put its endpoints in, so those subnetworks and the limit on connections are the point.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetNetworkConnectivityServiceConnectionPolicyAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, id string) *networkconnectivity.ServiceConnectionPolicy {
	attrs, err := GetNetworkConnectivityServiceConnectionPolicyAttrsE(t, ctx, projectID, location, id)
	require.NoError(t, err)

	return attrs
}

// GetNetworkConnectivityServiceConnectionPolicyAttrsE returns the settings Google Cloud holds for the given service connection policy.
// The ctx parameter supports cancellation and timeouts.
func GetNetworkConnectivityServiceConnectionPolicyAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, id string) (*networkconnectivity.ServiceConnectionPolicy, error) {
	logger.Default.Logf(t, "Getting settings for service connection policy %s in %s in project %s", id, location, projectID)

	service, err := NewNetworkConnectivityServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetNetworkConnectivityServiceConnectionPolicyAttrsWithClient(ctx, service, projectID, location, id)
}

// GetNetworkConnectivityServiceConnectionPolicyAttrsWithClient returns the settings Google Cloud holds for the given service connection policy using the supplied
// *networkconnectivity.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see networkconnectivity_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetNetworkConnectivityServiceConnectionPolicyAttrsWithClient(ctx context.Context, service *networkconnectivity.Service, projectID string, location string, id string) (*networkconnectivity.ServiceConnectionPolicy, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/serviceConnectionPolicies/%s", projectID, location, id)

	attrs, err := service.Projects.Locations.ServiceConnectionPolicies.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the service connection policy %s in %s in project %s does not exist", id, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for service connection policy %s in %s in project %s: %w", id, location, projectID, err)
	}

	return attrs, nil
}

// GetNetworkConnectivityTransportAttrs returns the settings Google Cloud holds for the given transport, so a test can assert on what was
// actually created rather than only that it exists.
// A transport is the pair of links a multicloud connection runs over, so its bandwidth and profile describe what it can carry.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetNetworkConnectivityTransportAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, id string) *networkconnectivity.Transport {
	attrs, err := GetNetworkConnectivityTransportAttrsE(t, ctx, projectID, location, id)
	require.NoError(t, err)

	return attrs
}

// GetNetworkConnectivityTransportAttrsE returns the settings Google Cloud holds for the given transport.
// The ctx parameter supports cancellation and timeouts.
func GetNetworkConnectivityTransportAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, id string) (*networkconnectivity.Transport, error) {
	logger.Default.Logf(t, "Getting settings for transport %s in %s in project %s", id, location, projectID)

	service, err := NewNetworkConnectivityServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetNetworkConnectivityTransportAttrsWithClient(ctx, service, projectID, location, id)
}

// GetNetworkConnectivityTransportAttrsWithClient returns the settings Google Cloud holds for the given transport using the supplied
// *networkconnectivity.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see networkconnectivity_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetNetworkConnectivityTransportAttrsWithClient(ctx context.Context, service *networkconnectivity.Service, projectID string, location string, id string) (*networkconnectivity.Transport, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/transports/%s", projectID, location, id)

	attrs, err := service.Projects.Locations.Transports.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the transport %s in %s in project %s does not exist", id, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for transport %s in %s in project %s: %w", id, location, projectID, err)
	}

	return attrs, nil
}

// GetNetworkConnectivityHubGroupAttrs returns the settings Google Cloud holds for the given hub group, so a test can assert on what was
// actually created rather than only that it exists.
// A group is how spokes on a hub are split into sets that may or may not reach each other, so its auto-accept list is what it decides.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetNetworkConnectivityHubGroupAttrs(t testing.TestingT, ctx context.Context, projectID string, hubID string, id string) *networkconnectivity.Group {
	attrs, err := GetNetworkConnectivityHubGroupAttrsE(t, ctx, projectID, hubID, id)
	require.NoError(t, err)

	return attrs
}

// GetNetworkConnectivityHubGroupAttrsE returns the settings Google Cloud holds for the given hub group.
// The ctx parameter supports cancellation and timeouts.
func GetNetworkConnectivityHubGroupAttrsE(t testing.TestingT, ctx context.Context, projectID string, hubID string, id string) (*networkconnectivity.Group, error) {
	logger.Default.Logf(t, "Getting settings for hub group %s in %s in project %s", id, hubID, projectID)

	service, err := NewNetworkConnectivityServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetNetworkConnectivityHubGroupAttrsWithClient(ctx, service, projectID, hubID, id)
}

// GetNetworkConnectivityHubGroupAttrsWithClient returns the settings Google Cloud holds for the given hub group using the supplied
// *networkconnectivity.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see networkconnectivity_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetNetworkConnectivityHubGroupAttrsWithClient(ctx context.Context, service *networkconnectivity.Service, projectID string, hubID string, id string) (*networkconnectivity.Group, error) {
	name := fmt.Sprintf("projects/%s/locations/global/hubs/%s/groups/%s", projectID, hubID, id)

	attrs, err := service.Projects.Locations.Global.Hubs.Groups.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the hub group %s in %s in project %s does not exist", id, hubID, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for hub group %s in %s in project %s: %w", id, hubID, projectID, err)
	}

	return attrs, nil
}

// GetNetworkConnectivityTransferDestinationAttrs returns the settings Google Cloud holds for the given transfer destination, so a test can assert on what was
// actually created rather than only that it exists.
// A destination is the address range on the other cloud that charged transfer applies to, so the range and its state are the point.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetNetworkConnectivityTransferDestinationAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, configID string, id string) *networkconnectivity.Destination {
	attrs, err := GetNetworkConnectivityTransferDestinationAttrsE(t, ctx, projectID, location, configID, id)
	require.NoError(t, err)

	return attrs
}

// GetNetworkConnectivityTransferDestinationAttrsE returns the settings Google Cloud holds for the given transfer destination.
// The ctx parameter supports cancellation and timeouts.
func GetNetworkConnectivityTransferDestinationAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, configID string, id string) (*networkconnectivity.Destination, error) {
	logger.Default.Logf(t, "Getting settings for transfer destination %s in %s in %s in project %s", id, configID, location, projectID)

	service, err := NewNetworkConnectivityServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetNetworkConnectivityTransferDestinationAttrsWithClient(ctx, service, projectID, location, configID, id)
}

// GetNetworkConnectivityTransferDestinationAttrsWithClient returns the settings Google Cloud holds for the given transfer destination using the supplied
// *networkconnectivity.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see networkconnectivity_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetNetworkConnectivityTransferDestinationAttrsWithClient(ctx context.Context, service *networkconnectivity.Service, projectID string, location string, configID string, id string) (*networkconnectivity.Destination, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/multicloudDataTransferConfigs/%s/destinations/%s", projectID, location, configID, id)

	attrs, err := service.Projects.Locations.MulticloudDataTransferConfigs.Destinations.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the transfer destination %s in %s in %s in project %s does not exist", id, configID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for transfer destination %s in %s in %s in project %s: %w", id, configID, location, projectID, err)
	}

	return attrs, nil
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

// GetInternalRangeAttrs returns the settings Google Cloud holds for the given internal range, so a
// test can assert on what was actually created rather than only that it exists. An internal range
// reserves address space inside a network so nothing else is given it, which is only worth anything if
// the range and its usage came out as asked.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetInternalRangeAttrs(t testing.TestingT, ctx context.Context, projectID string, rangeID string) *networkconnectivity.InternalRange {
	internalRange, err := GetInternalRangeAttrsE(t, ctx, projectID, rangeID)
	require.NoError(t, err)

	return internalRange
}

// GetInternalRangeAttrsE returns the settings Google Cloud holds for the given internal range.
// The ctx parameter supports cancellation and timeouts.
func GetInternalRangeAttrsE(t testing.TestingT, ctx context.Context, projectID string, rangeID string) (*networkconnectivity.InternalRange, error) {
	logger.Default.Logf(t, "Getting settings for internal range %s in project %s", rangeID, projectID)

	service, err := NewNetworkConnectivityServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetInternalRangeAttrsWithClient(ctx, service, projectID, rangeID)
}

// GetInternalRangeAttrsWithClient returns the settings Google Cloud holds for the given internal range
// using the supplied *networkconnectivity.Service. Prefer this variant in unit tests where the service
// is backed by an httptest fake server (see networkconnectivity_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetInternalRangeAttrsWithClient(ctx context.Context, service *networkconnectivity.Service, projectID string, rangeID string) (*networkconnectivity.InternalRange, error) {
	name := fmt.Sprintf("projects/%s/locations/global/internalRanges/%s", projectID, rangeID)

	internalRange, err := service.Projects.Locations.InternalRanges.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the internal range %s does not exist in project %s", rangeID, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for internal range %s in project %s: %w", rangeID, projectID, err)
	}

	return internalRange, nil
}

// NewNetworkConnectivityServiceE creates a Network Connectivity service authenticated the same way
// every other client in this module is.
// The ctx parameter supports cancellation and timeouts.
func NewNetworkConnectivityServiceE(t testing.TestingT, ctx context.Context) (*networkconnectivity.Service, error) {
	return networkconnectivity.NewService(ctx, append(withOptions(), option.WithScopes(networkconnectivity.CloudPlatformScope))...)
}
