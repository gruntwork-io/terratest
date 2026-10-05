package gcp

import (
	"context"
	"errors"
	"fmt"

	"github.com/gruntwork-io/terratest/modules/core/v2/logger"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/networkmanagement/v1"
	"google.golang.org/api/option"
)

// GetConnectivityTestAttrs returns the settings Google Cloud holds for the given connectivity test, so a test can assert on
// what was actually created rather than only that it exists. A connectivity test is a stored question about whether one
// endpoint can reach another, and its answer is the reachability verdict Google last computed.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetConnectivityTestAttrs(t testing.TestingT, ctx context.Context, projectID string, testID string) *networkmanagement.ConnectivityTest {
	connectivityTest, err := GetConnectivityTestAttrsE(t, ctx, projectID, testID)
	require.NoError(t, err)

	return connectivityTest
}

// GetConnectivityTestAttrsE returns the settings Google Cloud holds for the given connectivity test.
// The ctx parameter supports cancellation and timeouts.
func GetConnectivityTestAttrsE(t testing.TestingT, ctx context.Context, projectID string, testID string) (*networkmanagement.ConnectivityTest, error) {
	logger.Default.Logf(t, "Getting settings for connectivity test %s in project %s", testID, projectID)

	service, err := NewNetworkManagementServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetConnectivityTestAttrsWithClient(ctx, service, projectID, testID)
}

// GetConnectivityTestAttrsWithClient returns the settings Google Cloud holds for the given connectivity test using the supplied
// *networkmanagement.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see networkmanagement_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetConnectivityTestAttrsWithClient(ctx context.Context, service *networkmanagement.Service, projectID string, testID string) (*networkmanagement.ConnectivityTest, error) {
	name := fmt.Sprintf("projects/%s/locations/global/connectivityTests/%s", projectID, testID)

	connectivityTest, err := service.Projects.Locations.Global.ConnectivityTests.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the connectivity test %s does not exist in project %s", testID, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for connectivity test %s in project %s: %w", testID, projectID, err)
	}

	return connectivityTest, nil
}

// NewNetworkManagementServiceE creates a Network Management service authenticated the same way every other client in this
// module is.
// The ctx parameter supports cancellation and timeouts.
func NewNetworkManagementServiceE(t testing.TestingT, ctx context.Context) (*networkmanagement.Service, error) {
	return networkmanagement.NewService(ctx, append(withOptions(), option.WithScopes(networkmanagement.CloudPlatformScope))...)
}
