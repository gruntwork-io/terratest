package gcp

import (
	"context"
	"errors"
	"fmt"

	"github.com/gruntwork-io/terratest/modules/core/v2/logger"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/networkservices/v1"
	"google.golang.org/api/option"
)

// GetNetworkServicesMeshAttrs returns the settings Google Cloud holds for the given service mesh, so a test can assert on
// what was actually created rather than only that it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetNetworkServicesMeshAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, meshID string) *networkservices.Mesh {
	mesh, err := GetNetworkServicesMeshAttrsE(t, ctx, projectID, location, meshID)
	require.NoError(t, err)

	return mesh
}

// GetNetworkServicesMeshAttrsE returns the settings Google Cloud holds for the given service mesh.
// The ctx parameter supports cancellation and timeouts.
func GetNetworkServicesMeshAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, meshID string) (*networkservices.Mesh, error) {
	logger.Default.Logf(t, "Getting settings for mesh %s in %s in project %s", meshID, location, projectID)

	service, err := NewNetworkServicesServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetNetworkServicesMeshAttrsWithClient(ctx, service, projectID, location, meshID)
}

// GetNetworkServicesMeshAttrsWithClient returns the settings Google Cloud holds for the given service mesh using the supplied
// *networkservices.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see networkservices_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetNetworkServicesMeshAttrsWithClient(ctx context.Context, service *networkservices.Service, projectID string, location string, meshID string) (*networkservices.Mesh, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/meshes/%s", projectID, location, meshID)

	mesh, err := service.Projects.Locations.Meshes.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the mesh %s does not exist in %s in project %s", meshID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for mesh %s in %s in project %s: %w", meshID, location, projectID, err)
	}

	return mesh, nil
}

// GetNetworkServicesGRPCRouteAttrs returns the settings Google Cloud holds for the given gRPC route, so a test can assert on
// what was actually created rather than only that it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetNetworkServicesGRPCRouteAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, routeID string) *networkservices.GrpcRoute {
	route, err := GetNetworkServicesGRPCRouteAttrsE(t, ctx, projectID, location, routeID)
	require.NoError(t, err)

	return route
}

// GetNetworkServicesGRPCRouteAttrsE returns the settings Google Cloud holds for the given gRPC route.
// The ctx parameter supports cancellation and timeouts.
func GetNetworkServicesGRPCRouteAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, routeID string) (*networkservices.GrpcRoute, error) {
	logger.Default.Logf(t, "Getting settings for gRPC route %s in %s in project %s", routeID, location, projectID)

	service, err := NewNetworkServicesServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetNetworkServicesGRPCRouteAttrsWithClient(ctx, service, projectID, location, routeID)
}

// GetNetworkServicesGRPCRouteAttrsWithClient returns the settings Google Cloud holds for the given gRPC route using the supplied
// *networkservices.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see networkservices_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetNetworkServicesGRPCRouteAttrsWithClient(ctx context.Context, service *networkservices.Service, projectID string, location string, routeID string) (*networkservices.GrpcRoute, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/grpcRoutes/%s", projectID, location, routeID)

	route, err := service.Projects.Locations.GrpcRoutes.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the gRPC route %s does not exist in %s in project %s", routeID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for gRPC route %s in %s in project %s: %w", routeID, location, projectID, err)
	}

	return route, nil
}

// GetNetworkServicesHTTPRouteAttrs returns the settings Google Cloud holds for the given HTTP route, so a test can assert on
// what was actually created rather than only that it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetNetworkServicesHTTPRouteAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, routeID string) *networkservices.HttpRoute {
	route, err := GetNetworkServicesHTTPRouteAttrsE(t, ctx, projectID, location, routeID)
	require.NoError(t, err)

	return route
}

// GetNetworkServicesHTTPRouteAttrsE returns the settings Google Cloud holds for the given HTTP route.
// The ctx parameter supports cancellation and timeouts.
func GetNetworkServicesHTTPRouteAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, routeID string) (*networkservices.HttpRoute, error) {
	logger.Default.Logf(t, "Getting settings for HTTP route %s in %s in project %s", routeID, location, projectID)

	service, err := NewNetworkServicesServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetNetworkServicesHTTPRouteAttrsWithClient(ctx, service, projectID, location, routeID)
}

// GetNetworkServicesHTTPRouteAttrsWithClient returns the settings Google Cloud holds for the given HTTP route using the supplied
// *networkservices.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see networkservices_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetNetworkServicesHTTPRouteAttrsWithClient(ctx context.Context, service *networkservices.Service, projectID string, location string, routeID string) (*networkservices.HttpRoute, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/httpRoutes/%s", projectID, location, routeID)

	route, err := service.Projects.Locations.HttpRoutes.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the HTTP route %s does not exist in %s in project %s", routeID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for HTTP route %s in %s in project %s: %w", routeID, location, projectID, err)
	}

	return route, nil
}

// GetNetworkServicesTCPRouteAttrs returns the settings Google Cloud holds for the given TCP route, so a test can assert on
// what was actually created rather than only that it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetNetworkServicesTCPRouteAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, routeID string) *networkservices.TcpRoute {
	route, err := GetNetworkServicesTCPRouteAttrsE(t, ctx, projectID, location, routeID)
	require.NoError(t, err)

	return route
}

// GetNetworkServicesTCPRouteAttrsE returns the settings Google Cloud holds for the given TCP route.
// The ctx parameter supports cancellation and timeouts.
func GetNetworkServicesTCPRouteAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, routeID string) (*networkservices.TcpRoute, error) {
	logger.Default.Logf(t, "Getting settings for TCP route %s in %s in project %s", routeID, location, projectID)

	service, err := NewNetworkServicesServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetNetworkServicesTCPRouteAttrsWithClient(ctx, service, projectID, location, routeID)
}

// GetNetworkServicesTCPRouteAttrsWithClient returns the settings Google Cloud holds for the given TCP route using the supplied
// *networkservices.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see networkservices_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetNetworkServicesTCPRouteAttrsWithClient(ctx context.Context, service *networkservices.Service, projectID string, location string, routeID string) (*networkservices.TcpRoute, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/tcpRoutes/%s", projectID, location, routeID)

	route, err := service.Projects.Locations.TcpRoutes.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the TCP route %s does not exist in %s in project %s", routeID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for TCP route %s in %s in project %s: %w", routeID, location, projectID, err)
	}

	return route, nil
}

// GetNetworkServicesTLSRouteAttrs returns the settings Google Cloud holds for the given TLS route, so a test can assert on
// what was actually created rather than only that it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetNetworkServicesTLSRouteAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, routeID string) *networkservices.TlsRoute {
	route, err := GetNetworkServicesTLSRouteAttrsE(t, ctx, projectID, location, routeID)
	require.NoError(t, err)

	return route
}

// GetNetworkServicesTLSRouteAttrsE returns the settings Google Cloud holds for the given TLS route.
// The ctx parameter supports cancellation and timeouts.
func GetNetworkServicesTLSRouteAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, routeID string) (*networkservices.TlsRoute, error) {
	logger.Default.Logf(t, "Getting settings for TLS route %s in %s in project %s", routeID, location, projectID)

	service, err := NewNetworkServicesServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetNetworkServicesTLSRouteAttrsWithClient(ctx, service, projectID, location, routeID)
}

// GetNetworkServicesTLSRouteAttrsWithClient returns the settings Google Cloud holds for the given TLS route using the supplied
// *networkservices.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see networkservices_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetNetworkServicesTLSRouteAttrsWithClient(ctx context.Context, service *networkservices.Service, projectID string, location string, routeID string) (*networkservices.TlsRoute, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/tlsRoutes/%s", projectID, location, routeID)

	route, err := service.Projects.Locations.TlsRoutes.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the TLS route %s does not exist in %s in project %s", routeID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for TLS route %s in %s in project %s: %w", routeID, location, projectID, err)
	}

	return route, nil
}

// GetNetworkServicesEndpointPolicyAttrs returns the settings Google Cloud holds for the given endpoint policy, so a test can assert on
// what was actually created rather than only that it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetNetworkServicesEndpointPolicyAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, policyID string) *networkservices.EndpointPolicy {
	policy, err := GetNetworkServicesEndpointPolicyAttrsE(t, ctx, projectID, location, policyID)
	require.NoError(t, err)

	return policy
}

// GetNetworkServicesEndpointPolicyAttrsE returns the settings Google Cloud holds for the given endpoint policy.
// The ctx parameter supports cancellation and timeouts.
func GetNetworkServicesEndpointPolicyAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, policyID string) (*networkservices.EndpointPolicy, error) {
	logger.Default.Logf(t, "Getting settings for endpoint policy %s in %s in project %s", policyID, location, projectID)

	service, err := NewNetworkServicesServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetNetworkServicesEndpointPolicyAttrsWithClient(ctx, service, projectID, location, policyID)
}

// GetNetworkServicesEndpointPolicyAttrsWithClient returns the settings Google Cloud holds for the given endpoint policy using the supplied
// *networkservices.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see networkservices_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetNetworkServicesEndpointPolicyAttrsWithClient(ctx context.Context, service *networkservices.Service, projectID string, location string, policyID string) (*networkservices.EndpointPolicy, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/endpointPolicies/%s", projectID, location, policyID)

	policy, err := service.Projects.Locations.EndpointPolicies.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the endpoint policy %s does not exist in %s in project %s", policyID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for endpoint policy %s in %s in project %s: %w", policyID, location, projectID, err)
	}

	return policy, nil
}

// GetNetworkServicesAuthzExtensionAttrs returns the settings Google Cloud holds for the given authorization extension, so a test can assert on what was
// actually created rather than only that it exists.
// An extension hands an authorization decision to a service of our own, so the service it calls and how long it waits are the point.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetNetworkServicesAuthzExtensionAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, id string) *networkservices.AuthzExtension {
	attrs, err := GetNetworkServicesAuthzExtensionAttrsE(t, ctx, projectID, location, id)
	require.NoError(t, err)

	return attrs
}

// GetNetworkServicesAuthzExtensionAttrsE returns the settings Google Cloud holds for the given authorization extension.
// The ctx parameter supports cancellation and timeouts.
func GetNetworkServicesAuthzExtensionAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, id string) (*networkservices.AuthzExtension, error) {
	logger.Default.Logf(t, "Getting settings for authorization extension %s in %s in project %s", id, location, projectID)

	service, err := NewNetworkServicesServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetNetworkServicesAuthzExtensionAttrsWithClient(ctx, service, projectID, location, id)
}

// GetNetworkServicesAuthzExtensionAttrsWithClient returns the settings Google Cloud holds for the given authorization extension using the supplied
// *networkservices.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see networkservices_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetNetworkServicesAuthzExtensionAttrsWithClient(ctx context.Context, service *networkservices.Service, projectID string, location string, id string) (*networkservices.AuthzExtension, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/authzExtensions/%s", projectID, location, id)

	attrs, err := service.Projects.Locations.AuthzExtensions.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the authorization extension %s in %s in project %s does not exist", id, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for authorization extension %s in %s in project %s: %w", id, location, projectID, err)
	}

	return attrs, nil
}

// GetNetworkServicesGatewayAttrs returns the settings Google Cloud holds for the given gateway, so a test can assert on what was
// actually created rather than only that it exists.
// A gateway is where a mesh accepts traffic from outside it, so which ports it listens on and which type it is decide what can reach it.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetNetworkServicesGatewayAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, id string) *networkservices.Gateway {
	attrs, err := GetNetworkServicesGatewayAttrsE(t, ctx, projectID, location, id)
	require.NoError(t, err)

	return attrs
}

// GetNetworkServicesGatewayAttrsE returns the settings Google Cloud holds for the given gateway.
// The ctx parameter supports cancellation and timeouts.
func GetNetworkServicesGatewayAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, id string) (*networkservices.Gateway, error) {
	logger.Default.Logf(t, "Getting settings for gateway %s in %s in project %s", id, location, projectID)

	service, err := NewNetworkServicesServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetNetworkServicesGatewayAttrsWithClient(ctx, service, projectID, location, id)
}

// GetNetworkServicesGatewayAttrsWithClient returns the settings Google Cloud holds for the given gateway using the supplied
// *networkservices.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see networkservices_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetNetworkServicesGatewayAttrsWithClient(ctx context.Context, service *networkservices.Service, projectID string, location string, id string) (*networkservices.Gateway, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/gateways/%s", projectID, location, id)

	attrs, err := service.Projects.Locations.Gateways.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the gateway %s in %s in project %s does not exist", id, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for gateway %s in %s in project %s: %w", id, location, projectID, err)
	}

	return attrs, nil
}

// GetNetworkServicesLbEdgeExtensionAttrs returns the settings Google Cloud holds for the given load balancer edge extension, so a test can assert on what was
// actually created rather than only that it exists.
// An edge extension runs at the edge before a request is routed, so which forwarding rules it attaches to is what it affects.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetNetworkServicesLbEdgeExtensionAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, id string) *networkservices.LbEdgeExtension {
	attrs, err := GetNetworkServicesLbEdgeExtensionAttrsE(t, ctx, projectID, location, id)
	require.NoError(t, err)

	return attrs
}

// GetNetworkServicesLbEdgeExtensionAttrsE returns the settings Google Cloud holds for the given load balancer edge extension.
// The ctx parameter supports cancellation and timeouts.
func GetNetworkServicesLbEdgeExtensionAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, id string) (*networkservices.LbEdgeExtension, error) {
	logger.Default.Logf(t, "Getting settings for load balancer edge extension %s in %s in project %s", id, location, projectID)

	service, err := NewNetworkServicesServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetNetworkServicesLbEdgeExtensionAttrsWithClient(ctx, service, projectID, location, id)
}

// GetNetworkServicesLbEdgeExtensionAttrsWithClient returns the settings Google Cloud holds for the given load balancer edge extension using the supplied
// *networkservices.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see networkservices_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetNetworkServicesLbEdgeExtensionAttrsWithClient(ctx context.Context, service *networkservices.Service, projectID string, location string, id string) (*networkservices.LbEdgeExtension, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/lbEdgeExtensions/%s", projectID, location, id)

	attrs, err := service.Projects.Locations.LbEdgeExtensions.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the load balancer edge extension %s in %s in project %s does not exist", id, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for load balancer edge extension %s in %s in project %s: %w", id, location, projectID, err)
	}

	return attrs, nil
}

// GetNetworkServicesLbRouteExtensionAttrs returns the settings Google Cloud holds for the given load balancer route extension, so a test can assert on what was
// actually created rather than only that it exists.
// A route extension can pick the backend for a request, so the extension chain it runs is what decides where traffic lands.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetNetworkServicesLbRouteExtensionAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, id string) *networkservices.LbRouteExtension {
	attrs, err := GetNetworkServicesLbRouteExtensionAttrsE(t, ctx, projectID, location, id)
	require.NoError(t, err)

	return attrs
}

// GetNetworkServicesLbRouteExtensionAttrsE returns the settings Google Cloud holds for the given load balancer route extension.
// The ctx parameter supports cancellation and timeouts.
func GetNetworkServicesLbRouteExtensionAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, id string) (*networkservices.LbRouteExtension, error) {
	logger.Default.Logf(t, "Getting settings for load balancer route extension %s in %s in project %s", id, location, projectID)

	service, err := NewNetworkServicesServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetNetworkServicesLbRouteExtensionAttrsWithClient(ctx, service, projectID, location, id)
}

// GetNetworkServicesLbRouteExtensionAttrsWithClient returns the settings Google Cloud holds for the given load balancer route extension using the supplied
// *networkservices.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see networkservices_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetNetworkServicesLbRouteExtensionAttrsWithClient(ctx context.Context, service *networkservices.Service, projectID string, location string, id string) (*networkservices.LbRouteExtension, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/lbRouteExtensions/%s", projectID, location, id)

	attrs, err := service.Projects.Locations.LbRouteExtensions.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the load balancer route extension %s in %s in project %s does not exist", id, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for load balancer route extension %s in %s in project %s: %w", id, location, projectID, err)
	}

	return attrs, nil
}

// GetNetworkServicesLbTrafficExtensionAttrs returns the settings Google Cloud holds for the given load balancer traffic extension, so a test can assert on what was
// actually created rather than only that it exists.
// A traffic extension sees requests and responses as they pass, so which chain runs and on which rules is the whole point.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetNetworkServicesLbTrafficExtensionAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, id string) *networkservices.LbTrafficExtension {
	attrs, err := GetNetworkServicesLbTrafficExtensionAttrsE(t, ctx, projectID, location, id)
	require.NoError(t, err)

	return attrs
}

// GetNetworkServicesLbTrafficExtensionAttrsE returns the settings Google Cloud holds for the given load balancer traffic extension.
// The ctx parameter supports cancellation and timeouts.
func GetNetworkServicesLbTrafficExtensionAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, id string) (*networkservices.LbTrafficExtension, error) {
	logger.Default.Logf(t, "Getting settings for load balancer traffic extension %s in %s in project %s", id, location, projectID)

	service, err := NewNetworkServicesServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetNetworkServicesLbTrafficExtensionAttrsWithClient(ctx, service, projectID, location, id)
}

// GetNetworkServicesLbTrafficExtensionAttrsWithClient returns the settings Google Cloud holds for the given load balancer traffic extension using the supplied
// *networkservices.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see networkservices_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetNetworkServicesLbTrafficExtensionAttrsWithClient(ctx context.Context, service *networkservices.Service, projectID string, location string, id string) (*networkservices.LbTrafficExtension, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/lbTrafficExtensions/%s", projectID, location, id)

	attrs, err := service.Projects.Locations.LbTrafficExtensions.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the load balancer traffic extension %s in %s in project %s does not exist", id, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for load balancer traffic extension %s in %s in project %s: %w", id, location, projectID, err)
	}

	return attrs, nil
}

// GetNetworkServicesWasmPluginAttrs returns the settings Google Cloud holds for the given Wasm plugin, so a test can assert on what was
// actually created rather than only that it exists.
// A plugin is the code an extension runs, so which version is the main one decides what actually executes.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetNetworkServicesWasmPluginAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, id string) *networkservices.WasmPlugin {
	attrs, err := GetNetworkServicesWasmPluginAttrsE(t, ctx, projectID, location, id)
	require.NoError(t, err)

	return attrs
}

// GetNetworkServicesWasmPluginAttrsE returns the settings Google Cloud holds for the given Wasm plugin.
// The ctx parameter supports cancellation and timeouts.
func GetNetworkServicesWasmPluginAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, id string) (*networkservices.WasmPlugin, error) {
	logger.Default.Logf(t, "Getting settings for Wasm plugin %s in %s in project %s", id, location, projectID)

	service, err := NewNetworkServicesServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetNetworkServicesWasmPluginAttrsWithClient(ctx, service, projectID, location, id)
}

// GetNetworkServicesWasmPluginAttrsWithClient returns the settings Google Cloud holds for the given Wasm plugin using the supplied
// *networkservices.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see networkservices_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetNetworkServicesWasmPluginAttrsWithClient(ctx context.Context, service *networkservices.Service, projectID string, location string, id string) (*networkservices.WasmPlugin, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/wasmPlugins/%s", projectID, location, id)

	attrs, err := service.Projects.Locations.WasmPlugins.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Wasm plugin %s in %s in project %s does not exist", id, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Wasm plugin %s in %s in project %s: %w", id, location, projectID, err)
	}

	return attrs, nil
}

// NewNetworkServicesServiceE creates a Network Services service authenticated the same way every other client in this
// module is.
// The ctx parameter supports cancellation and timeouts.
func NewNetworkServicesServiceE(t testing.TestingT, ctx context.Context) (*networkservices.Service, error) {
	return networkservices.NewService(ctx, append(withOptions(), option.WithScopes(networkservices.CloudPlatformScope))...)
}
