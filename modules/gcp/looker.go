package gcp

import (
	"context"
	"errors"
	"fmt"

	"github.com/gruntwork-io/terratest/modules/core/v2/logger"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/looker/v1"
	"google.golang.org/api/option"
)

// GetLookerInstanceAttrs returns the settings Google Cloud holds for the given Looker instance, so a test can assert on
// what was actually created rather than only that it exists.
// The URI users reach the instance on is assigned by Google, so it can only be asserted on from a
// read.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetLookerInstanceAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, instanceID string) *looker.Instance {
	result, err := GetLookerInstanceAttrsE(t, ctx, projectID, location, instanceID)
	require.NoError(t, err)

	return result
}

// GetLookerInstanceAttrsE returns the settings Google Cloud holds for the given Looker instance.
// The ctx parameter supports cancellation and timeouts.
func GetLookerInstanceAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, instanceID string) (*looker.Instance, error) {
	logger.Default.Logf(t, "Getting settings for Looker instance %s in %s of project %s", instanceID, location, projectID)

	service, err := NewLookerServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetLookerInstanceAttrsWithClient(ctx, service, projectID, location, instanceID)
}

// GetLookerInstanceAttrsWithClient returns the settings Google Cloud holds for the given Looker instance using the supplied
// *looker.Service. Prefer this variant in unit tests where the service is backed by an httptest fake server
// (see looker_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetLookerInstanceAttrsWithClient(ctx context.Context, service *looker.Service, projectID string, location string, instanceID string) (*looker.Instance, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/instances/%s", projectID, location, instanceID)

	result, err := service.Projects.Locations.Instances.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Looker instance %s does not exist in %s of project %s", instanceID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Looker instance %s in %s of project %s: %w", instanceID, location, projectID, err)
	}

	return result, nil
}

// NewLookerServiceE creates a Looker service authenticated the same way every other client in
// this module is.
// The ctx parameter supports cancellation and timeouts.
func NewLookerServiceE(t testing.TestingT, ctx context.Context) (*looker.Service, error) {
	return looker.NewService(ctx, append(withOptions(), option.WithScopes(looker.CloudPlatformScope))...)
}
