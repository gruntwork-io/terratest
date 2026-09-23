package gcp

import (
	"context"
	"fmt"

	"github.com/gruntwork-io/terratest/modules/core/v2/logger"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/compute/v1"
)

// FetchHAVPNGateway queries GCP to return the settings it holds for the given HA VPN gateway, so a test can
// assert on what was actually created rather than only that it exists. An HA VPN gateway is a different resource from the classic one, which FetchTargetVPNGateway reads.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func FetchHAVPNGateway(t testing.TestingT, ctx context.Context, projectID string, region string, name string) *compute.VpnGateway {
	gateway, err := FetchHAVPNGatewayE(t, ctx, projectID, region, name)
	require.NoError(t, err)

	return gateway
}

// FetchHAVPNGatewayE queries GCP to return the settings it holds for the given HA VPN gateway.
// The ctx parameter supports cancellation and timeouts.
func FetchHAVPNGatewayE(t testing.TestingT, ctx context.Context, projectID string, region string, name string) (*compute.VpnGateway, error) {
	logger.Default.Logf(t, "Getting HA VPN gateway %s in region %s", name, region)

	service, err := NewComputeServiceContextE(t, ctx)
	if err != nil {
		return nil, err
	}

	return FetchHAVPNGatewayWithClient(ctx, service, projectID, region, name)
}

// FetchHAVPNGatewayWithClient queries GCP to return the settings it holds for the given HA VPN gateway using the
// supplied *compute.Service. Prefer this variant in unit tests where the service is backed by an
// httptest fake server (see vpn_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func FetchHAVPNGatewayWithClient(ctx context.Context, service *compute.Service, projectID string, region string, name string) (*compute.VpnGateway, error) {
	gateway, err := service.VpnGateways.Get(projectID, region, name).Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("VpnGateways.Get(%s, %s, %s) got error: %w", projectID, region, name, err)
	}

	return gateway, nil
}

// FetchTargetVPNGateway queries GCP to return the settings it holds for the given classic VPN gateway, so a test can
// assert on what was actually created rather than only that it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func FetchTargetVPNGateway(t testing.TestingT, ctx context.Context, projectID string, region string, name string) *compute.TargetVpnGateway {
	gateway, err := FetchTargetVPNGatewayE(t, ctx, projectID, region, name)
	require.NoError(t, err)

	return gateway
}

// FetchTargetVPNGatewayE queries GCP to return the settings it holds for the given classic VPN gateway.
// The ctx parameter supports cancellation and timeouts.
func FetchTargetVPNGatewayE(t testing.TestingT, ctx context.Context, projectID string, region string, name string) (*compute.TargetVpnGateway, error) {
	logger.Default.Logf(t, "Getting classic VPN gateway %s in region %s", name, region)

	service, err := NewComputeServiceContextE(t, ctx)
	if err != nil {
		return nil, err
	}

	return FetchTargetVPNGatewayWithClient(ctx, service, projectID, region, name)
}

// FetchTargetVPNGatewayWithClient queries GCP to return the settings it holds for the given classic VPN gateway using the
// supplied *compute.Service. Prefer this variant in unit tests where the service is backed by an
// httptest fake server (see vpn_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func FetchTargetVPNGatewayWithClient(ctx context.Context, service *compute.Service, projectID string, region string, name string) (*compute.TargetVpnGateway, error) {
	gateway, err := service.TargetVpnGateways.Get(projectID, region, name).Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("TargetVpnGateways.Get(%s, %s, %s) got error: %w", projectID, region, name, err)
	}

	return gateway, nil
}

// FetchExternalVPNGateway queries GCP to return the settings it holds for the given external VPN gateway, so a test can
// assert on what was actually created rather than only that it exists. An external gateway describes the far end of a tunnel, which is not in this project, so it is global.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func FetchExternalVPNGateway(t testing.TestingT, ctx context.Context, projectID string, name string) *compute.ExternalVpnGateway {
	gateway, err := FetchExternalVPNGatewayE(t, ctx, projectID, name)
	require.NoError(t, err)

	return gateway
}

// FetchExternalVPNGatewayE queries GCP to return the settings it holds for the given external VPN gateway.
// The ctx parameter supports cancellation and timeouts.
func FetchExternalVPNGatewayE(t testing.TestingT, ctx context.Context, projectID string, name string) (*compute.ExternalVpnGateway, error) {
	logger.Default.Logf(t, "Getting external VPN gateway %s", name)

	service, err := NewComputeServiceContextE(t, ctx)
	if err != nil {
		return nil, err
	}

	return FetchExternalVPNGatewayWithClient(ctx, service, projectID, name)
}

// FetchExternalVPNGatewayWithClient queries GCP to return the settings it holds for the given external VPN gateway using the
// supplied *compute.Service. Prefer this variant in unit tests where the service is backed by an
// httptest fake server (see vpn_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func FetchExternalVPNGatewayWithClient(ctx context.Context, service *compute.Service, projectID string, name string) (*compute.ExternalVpnGateway, error) {
	gateway, err := service.ExternalVpnGateways.Get(projectID, name).Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("ExternalVpnGateways.Get(%s, %s) got error: %w", projectID, name, err)
	}

	return gateway, nil
}

// FetchVPNTunnel queries GCP to return the settings it holds for the given VPN tunnel, so a test can
// assert on what was actually created rather than only that it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func FetchVPNTunnel(t testing.TestingT, ctx context.Context, projectID string, region string, name string) *compute.VpnTunnel {
	tunnel, err := FetchVPNTunnelE(t, ctx, projectID, region, name)
	require.NoError(t, err)

	return tunnel
}

// FetchVPNTunnelE queries GCP to return the settings it holds for the given VPN tunnel.
// The ctx parameter supports cancellation and timeouts.
func FetchVPNTunnelE(t testing.TestingT, ctx context.Context, projectID string, region string, name string) (*compute.VpnTunnel, error) {
	logger.Default.Logf(t, "Getting VPN tunnel %s in region %s", name, region)

	service, err := NewComputeServiceContextE(t, ctx)
	if err != nil {
		return nil, err
	}

	return FetchVPNTunnelWithClient(ctx, service, projectID, region, name)
}

// FetchVPNTunnelWithClient queries GCP to return the settings it holds for the given VPN tunnel using the
// supplied *compute.Service. Prefer this variant in unit tests where the service is backed by an
// httptest fake server (see vpn_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func FetchVPNTunnelWithClient(ctx context.Context, service *compute.Service, projectID string, region string, name string) (*compute.VpnTunnel, error) {
	tunnel, err := service.VpnTunnels.Get(projectID, region, name).Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("VpnTunnels.Get(%s, %s, %s) got error: %w", projectID, region, name, err)
	}

	return tunnel, nil
}

// FetchRouterRoutePolicy queries GCP to return the settings it holds for the given route policy on
// a Cloud Router, so a test can assert on what was actually created rather than only that it
// exists. A route policy belongs to a router rather than standing on its own, so it is asked for by
// router and by name.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func FetchRouterRoutePolicy(t testing.TestingT, ctx context.Context, projectID string, region string, router string, name string) *compute.RoutePolicy {
	policy, err := FetchRouterRoutePolicyE(t, ctx, projectID, region, router, name)
	require.NoError(t, err)

	return policy
}

// FetchRouterRoutePolicyE queries GCP to return the settings it holds for the given route policy on
// a Cloud Router.
// The ctx parameter supports cancellation and timeouts.
func FetchRouterRoutePolicyE(t testing.TestingT, ctx context.Context, projectID string, region string, router string, name string) (*compute.RoutePolicy, error) {
	logger.Default.Logf(t, "Getting route policy %s on router %s in region %s", name, router, region)

	service, err := NewComputeServiceContextE(t, ctx)
	if err != nil {
		return nil, err
	}

	return FetchRouterRoutePolicyWithClient(ctx, service, projectID, region, router, name)
}

// FetchRouterRoutePolicyWithClient queries GCP to return the settings it holds for the given route
// policy on a Cloud Router using the supplied *compute.Service. Prefer this variant in unit tests
// where the service is backed by an httptest fake server (see vpn_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func FetchRouterRoutePolicyWithClient(ctx context.Context, service *compute.Service, projectID string, region string, router string, name string) (*compute.RoutePolicy, error) {
	// The policy is named by a query parameter rather than by a path segment, so it cannot be
	// fetched the way every other resource in this file is.
	response, err := service.Routers.GetRoutePolicy(projectID, region, router).Policy(name).Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("Routers.GetRoutePolicy(%s, %s, %s, %s) got error: %w", projectID, region, router, name, err)
	}

	if response.Resource == nil {
		return nil, fmt.Errorf("router %s in region %s in project %s holds no route policy named %s", router, region, projectID, name)
	}

	return response.Resource, nil
}
