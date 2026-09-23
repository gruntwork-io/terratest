package gcp

import (
	"context"
	"errors"
	"fmt"

	"github.com/gruntwork-io/terratest/modules/core/v2/logger"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/file/v1"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
)

// GetFilestoreInstanceAttrs returns the settings Google Cloud holds for the given Filestore
// instance, so a test can assert on what was actually created rather than only that it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetFilestoreInstanceAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, instanceID string) *file.Instance {
	instance, err := GetFilestoreInstanceAttrsE(t, ctx, projectID, location, instanceID)
	require.NoError(t, err)

	return instance
}

// GetFilestoreInstanceAttrsE returns the settings Google Cloud holds for the given Filestore
// instance.
// The ctx parameter supports cancellation and timeouts.
func GetFilestoreInstanceAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, instanceID string) (*file.Instance, error) {
	logger.Default.Logf(t, "Getting settings for Filestore instance %s in location %s in project %s", instanceID, location, projectID)

	service, err := NewFilestoreServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetFilestoreInstanceAttrsWithClient(ctx, service, projectID, location, instanceID)
}

// GetFilestoreInstanceAttrsWithClient returns the settings Google Cloud holds for the given
// Filestore instance using the supplied *file.Service. Prefer this variant in unit tests where the
// service is backed by an httptest fake server (see filestore_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetFilestoreInstanceAttrsWithClient(ctx context.Context, service *file.Service, projectID string, location string, instanceID string) (*file.Instance, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/instances/%s", projectID, location, instanceID)

	instance, err := service.Projects.Locations.Instances.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Filestore instance %s does not exist in location %s in project %s", instanceID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Filestore instance %s in location %s in project %s: %w", instanceID, location, projectID, err)
	}

	return instance, nil
}

// NewFilestoreServiceE creates a Filestore service authenticated the same way every other client in
// this module is.
// The ctx parameter supports cancellation and timeouts.
func NewFilestoreServiceE(t testing.TestingT, ctx context.Context) (*file.Service, error) {
	return file.NewService(ctx, append(withOptions(), option.WithScopes(file.CloudPlatformScope))...)
}
