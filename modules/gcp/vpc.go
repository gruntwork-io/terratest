package gcp

import (
	"context"
	"fmt"

	"github.com/gruntwork-io/terratest/modules/core/v2/logger"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/compute/v1"
)

// FetchSubnetworkContext queries GCP to return the settings it holds for the given VPC subnetwork,
// so a test can assert on what was actually created rather than only that it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func FetchSubnetworkContext(t testing.TestingT, ctx context.Context, projectID string, region string, name string) *compute.Subnetwork {
	subnetwork, err := FetchSubnetworkContextE(t, ctx, projectID, region, name)
	require.NoError(t, err)

	return subnetwork
}

// FetchSubnetworkContextE queries GCP to return the settings it holds for the given VPC subnetwork.
// The ctx parameter supports cancellation and timeouts.
func FetchSubnetworkContextE(t testing.TestingT, ctx context.Context, projectID string, region string, name string) (*compute.Subnetwork, error) {
	logger.Default.Logf(t, "Getting VPC subnetwork %s in region %s", name, region)

	service, err := NewComputeServiceContextE(t, ctx)
	if err != nil {
		return nil, err
	}

	return FetchSubnetworkWithClient(ctx, service, projectID, region, name)
}

// FetchSubnetworkWithClient queries GCP to return the settings it holds for the given VPC
// subnetwork using the supplied *compute.Service. Prefer this variant in unit tests where the
// service is backed by an httptest fake server (see vpc_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func FetchSubnetworkWithClient(ctx context.Context, service *compute.Service, projectID string, region string, name string) (*compute.Subnetwork, error) {
	subnetwork, err := service.Subnetworks.Get(projectID, region, name).Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("Subnetworks.Get(%s, %s, %s) got error: %w", projectID, region, name, err)
	}

	return subnetwork, nil
}

// FetchRouterContext queries GCP to return the settings it holds for the given Cloud Router, so a
// test can assert on what was actually created rather than only that it exists. A Cloud NAT gateway
// is part of its router rather than a resource of its own, so its settings come back in the
// router's Nats field.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func FetchRouterContext(t testing.TestingT, ctx context.Context, projectID string, region string, name string) *compute.Router {
	router, err := FetchRouterContextE(t, ctx, projectID, region, name)
	require.NoError(t, err)

	return router
}

// FetchRouterContextE queries GCP to return the settings it holds for the given Cloud Router.
// The ctx parameter supports cancellation and timeouts.
func FetchRouterContextE(t testing.TestingT, ctx context.Context, projectID string, region string, name string) (*compute.Router, error) {
	logger.Default.Logf(t, "Getting Cloud Router %s in region %s", name, region)

	service, err := NewComputeServiceContextE(t, ctx)
	if err != nil {
		return nil, err
	}

	return FetchRouterWithClient(ctx, service, projectID, region, name)
}

// FetchRouterWithClient queries GCP to return the settings it holds for the given Cloud Router
// using the supplied *compute.Service. Prefer this variant in unit tests where the service is
// backed by an httptest fake server (see vpc_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func FetchRouterWithClient(ctx context.Context, service *compute.Service, projectID string, region string, name string) (*compute.Router, error) {
	router, err := service.Routers.Get(projectID, region, name).Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("Routers.Get(%s, %s, %s) got error: %w", projectID, region, name, err)
	}

	return router, nil
}

// FetchAddressContext queries GCP to return the settings it holds for the given regional IP address
// reservation, so a test can assert on what was actually created rather than only that it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func FetchAddressContext(t testing.TestingT, ctx context.Context, projectID string, region string, name string) *compute.Address {
	address, err := FetchAddressContextE(t, ctx, projectID, region, name)
	require.NoError(t, err)

	return address
}

// FetchAddressContextE queries GCP to return the settings it holds for the given regional IP
// address reservation.
// The ctx parameter supports cancellation and timeouts.
func FetchAddressContextE(t testing.TestingT, ctx context.Context, projectID string, region string, name string) (*compute.Address, error) {
	logger.Default.Logf(t, "Getting regional IP address reservation %s in region %s", name, region)

	service, err := NewComputeServiceContextE(t, ctx)
	if err != nil {
		return nil, err
	}

	return FetchAddressWithClient(ctx, service, projectID, region, name)
}

// FetchAddressWithClient queries GCP to return the settings it holds for the given regional IP
// address reservation using the supplied *compute.Service. Prefer this variant in unit tests where
// the service is backed by an httptest fake server (see vpc_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func FetchAddressWithClient(ctx context.Context, service *compute.Service, projectID string, region string, name string) (*compute.Address, error) {
	address, err := service.Addresses.Get(projectID, region, name).Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("Addresses.Get(%s, %s, %s) got error: %w", projectID, region, name, err)
	}

	return address, nil
}

// FetchGlobalAddressContext queries GCP to return the settings it holds for the given global IP
// address reservation, so a test can assert on what was actually created rather than only that it
// exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func FetchGlobalAddressContext(t testing.TestingT, ctx context.Context, projectID string, name string) *compute.Address {
	address, err := FetchGlobalAddressContextE(t, ctx, projectID, name)
	require.NoError(t, err)

	return address
}

// FetchGlobalAddressContextE queries GCP to return the settings it holds for the given global IP
// address reservation.
// The ctx parameter supports cancellation and timeouts.
func FetchGlobalAddressContextE(t testing.TestingT, ctx context.Context, projectID string, name string) (*compute.Address, error) {
	logger.Default.Logf(t, "Getting global IP address reservation %s", name)

	service, err := NewComputeServiceContextE(t, ctx)
	if err != nil {
		return nil, err
	}

	return FetchGlobalAddressWithClient(ctx, service, projectID, name)
}

// FetchGlobalAddressWithClient queries GCP to return the settings it holds for the given global IP
// address reservation using the supplied *compute.Service. Prefer this variant in unit tests where
// the service is backed by an httptest fake server (see vpc_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func FetchGlobalAddressWithClient(ctx context.Context, service *compute.Service, projectID string, name string) (*compute.Address, error) {
	address, err := service.GlobalAddresses.Get(projectID, name).Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("GlobalAddresses.Get(%s, %s) got error: %w", projectID, name, err)
	}

	return address, nil
}
