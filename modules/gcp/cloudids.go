package gcp

import (
	"context"
	"errors"
	"fmt"

	"github.com/gruntwork-io/terratest/modules/core/v2/logger"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/ids/v1"
	"google.golang.org/api/option"
)

// GetCloudIDSEndpointAttrs returns the settings Google Cloud holds for the Cloud IDS endpoint, so a test can assert on what was
// actually created rather than only that it exists.
// The severity an endpoint inspects at and the network it watches are the resource, and the endpoint
// is zonal even though the collection reads like a region.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetCloudIDSEndpointAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, endpointID string) *ids.Endpoint {
	result, err := GetCloudIDSEndpointAttrsE(t, ctx, projectID, location, endpointID)
	require.NoError(t, err)

	return result
}

// GetCloudIDSEndpointAttrsE returns the settings Google Cloud holds for the Cloud IDS endpoint.
// The ctx parameter supports cancellation and timeouts.
func GetCloudIDSEndpointAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, endpointID string) (*ids.Endpoint, error) {
	logger.Default.Logf(t, "Getting settings for Cloud IDS endpoint %s in %s in project %s", endpointID, location, projectID)

	service, err := NewCloudIDSServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetCloudIDSEndpointAttrsWithClient(ctx, service, projectID, location, endpointID)
}

// GetCloudIDSEndpointAttrsWithClient returns the settings Google Cloud holds for the Cloud IDS endpoint using the supplied
// *ids.Service. Prefer this variant in unit tests where the service is backed by an httptest
// fake server (see cloudids_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetCloudIDSEndpointAttrsWithClient(ctx context.Context, service *ids.Service, projectID string, location string, endpointID string) (*ids.Endpoint, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/endpoints/%s", projectID, location, endpointID)

	result, err := service.Projects.Locations.Endpoints.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Cloud IDS endpoint %s in %s in project %s does not exist", endpointID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Cloud IDS endpoint %s in %s in project %s: %w", endpointID, location, projectID, err)
	}

	return result, nil
}

// NewCloudIDSServiceE creates a Cloud IDS service authenticated the same way every other client in this
// module is.
// Google serves Cloud IDS from the ids API.
// The ctx parameter supports cancellation and timeouts.
func NewCloudIDSServiceE(t testing.TestingT, ctx context.Context) (*ids.Service, error) {
	return ids.NewService(ctx, append(withOptions(), option.WithScopes(ids.CloudPlatformScope))...)
}
