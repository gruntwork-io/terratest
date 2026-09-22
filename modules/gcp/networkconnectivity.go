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
