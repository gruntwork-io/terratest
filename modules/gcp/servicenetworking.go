package gcp

import (
	"context"
	"fmt"
	"strings"

	"github.com/gruntwork-io/terratest/modules/core/v2/logger"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/option"
	"google.golang.org/api/servicenetworking/v1"
)

// serviceNetworkingParent is the service a private connection is made to. Every connection a
// Terraform module creates is to this one, which is what serves Cloud SQL, Memorystore and the
// other managed products that live in Google's own network.
const serviceNetworkingParent = "services/servicenetworking.googleapis.com"

// GetServiceNetworkingConnectionAttrs returns the settings Google Cloud holds for the private
// connection between the given network and Google's services, so a test can assert on what was
// actually created rather than only that it exists. A connection has no name of its own and cannot
// be fetched one at a time, so this lists the network's connections and returns the only one.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetServiceNetworkingConnectionAttrs(t testing.TestingT, ctx context.Context, projectID string, networkName string) *servicenetworking.Connection {
	connection, err := GetServiceNetworkingConnectionAttrsE(t, ctx, projectID, networkName)
	require.NoError(t, err)

	return connection
}

// GetServiceNetworkingConnectionAttrsE returns the settings Google Cloud holds for the private
// connection between the given network and Google's services.
// The ctx parameter supports cancellation and timeouts.
func GetServiceNetworkingConnectionAttrsE(t testing.TestingT, ctx context.Context, projectID string, networkName string) (*servicenetworking.Connection, error) {
	logger.Default.Logf(t, "Getting settings for the private connection on network %s in project %s", networkName, projectID)

	service, err := NewServiceNetworkingServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetServiceNetworkingConnectionAttrsWithClient(ctx, service, projectID, networkName)
}

// GetServiceNetworkingConnectionAttrsWithClient returns the settings Google Cloud holds for the
// private connection between the given network and Google's services using the supplied
// *servicenetworking.APIService. Prefer this variant in unit tests where the service is backed by an
// httptest fake server (see servicenetworking_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetServiceNetworkingConnectionAttrsWithClient(ctx context.Context, service *servicenetworking.APIService, projectID string, networkName string) (*servicenetworking.Connection, error) {
	network := fmt.Sprintf("projects/%s/global/networks/%s", projectID, networkName)

	response, err := service.Services.Connections.List(serviceNetworkingParent).Network(network).Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("failed to list the private connections on network %s in project %s: %w", networkName, projectID, err)
	}

	// The list is filtered to one network, so anything but a single connection is the test's answer
	// rather than a choice this read should make.
	switch len(response.Connections) {
	case 0:
		return nil, fmt.Errorf("no private connection exists on network %s in project %s", networkName, projectID)
	case 1:
		return response.Connections[0], nil
	default:
		names := make([]string, 0, len(response.Connections))
		for _, connection := range response.Connections {
			names = append(names, connection.Peering)
		}

		return nil, fmt.Errorf("network %s in project %s has %d private connections, %s, and not the one this read expects", networkName, projectID, len(response.Connections), strings.Join(names, ", "))
	}
}

// NewServiceNetworkingServiceE creates a Service Networking service authenticated the same way
// every other client in this module is.
// The ctx parameter supports cancellation and timeouts.
func NewServiceNetworkingServiceE(t testing.TestingT, ctx context.Context) (*servicenetworking.APIService, error) {
	return servicenetworking.NewService(ctx, append(withOptions(), option.WithScopes(servicenetworking.CloudPlatformScope))...)
}
