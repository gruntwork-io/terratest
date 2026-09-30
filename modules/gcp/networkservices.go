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

// NewNetworkServicesServiceE creates a Network Services service authenticated the same way every other client in this
// module is.
// The ctx parameter supports cancellation and timeouts.
func NewNetworkServicesServiceE(t testing.TestingT, ctx context.Context) (*networkservices.Service, error) {
	return networkservices.NewService(ctx, append(withOptions(), option.WithScopes(networkservices.CloudPlatformScope))...)
}
