package gcp

import (
	"context"
	"errors"
	"fmt"

	"github.com/gruntwork-io/terratest/modules/core/v2/logger"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
	"google.golang.org/api/spanner/v1"
)

// GetSpannerInstanceAttrs returns the settings Google Cloud holds for the given Spanner instance,
// so a test can assert on what was actually created rather than only that it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetSpannerInstanceAttrs(t testing.TestingT, ctx context.Context, projectID string, instanceID string) *spanner.Instance {
	instance, err := GetSpannerInstanceAttrsE(t, ctx, projectID, instanceID)
	require.NoError(t, err)

	return instance
}

// GetSpannerInstanceAttrsE returns the settings Google Cloud holds for the given Spanner instance.
// The ctx parameter supports cancellation and timeouts.
func GetSpannerInstanceAttrsE(t testing.TestingT, ctx context.Context, projectID string, instanceID string) (*spanner.Instance, error) {
	logger.Default.Logf(t, "Getting settings for Spanner instance %s in project %s", instanceID, projectID)

	service, err := NewSpannerServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetSpannerInstanceAttrsWithClient(ctx, service, projectID, instanceID)
}

// GetSpannerInstanceAttrsWithClient returns the settings Google Cloud holds for the given Spanner
// instance using the supplied *spanner.Service. Prefer this variant in unit tests where the service
// is backed by an httptest fake server (see spanner_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetSpannerInstanceAttrsWithClient(ctx context.Context, service *spanner.Service, projectID string, instanceID string) (*spanner.Instance, error) {
	name := fmt.Sprintf("projects/%s/instances/%s", projectID, instanceID)

	instance, err := service.Projects.Instances.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Spanner instance %s does not exist in project %s", instanceID, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Spanner instance %s in project %s: %w", instanceID, projectID, err)
	}

	return instance, nil
}

// GetSpannerDatabaseAttrs returns the settings Google Cloud holds for the given Spanner database,
// so a test can assert on what was actually created rather than only that it exists. A database is
// named by its instance as well as its own id.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetSpannerDatabaseAttrs(t testing.TestingT, ctx context.Context, projectID string, instanceID string, databaseID string) *spanner.Database {
	database, err := GetSpannerDatabaseAttrsE(t, ctx, projectID, instanceID, databaseID)
	require.NoError(t, err)

	return database
}

// GetSpannerDatabaseAttrsE returns the settings Google Cloud holds for the given Spanner database.
// The ctx parameter supports cancellation and timeouts.
func GetSpannerDatabaseAttrsE(t testing.TestingT, ctx context.Context, projectID string, instanceID string, databaseID string) (*spanner.Database, error) {
	logger.Default.Logf(t, "Getting settings for Spanner database %s in instance %s in project %s", databaseID, instanceID, projectID)

	service, err := NewSpannerServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetSpannerDatabaseAttrsWithClient(ctx, service, projectID, instanceID, databaseID)
}

// GetSpannerDatabaseAttrsWithClient returns the settings Google Cloud holds for the given Spanner
// database using the supplied *spanner.Service. Prefer this variant in unit tests where the service
// is backed by an httptest fake server (see spanner_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetSpannerDatabaseAttrsWithClient(ctx context.Context, service *spanner.Service, projectID string, instanceID string, databaseID string) (*spanner.Database, error) {
	name := fmt.Sprintf("projects/%s/instances/%s/databases/%s", projectID, instanceID, databaseID)

	database, err := service.Projects.Instances.Databases.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Spanner database %s does not exist in instance %s in project %s", databaseID, instanceID, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Spanner database %s in instance %s in project %s: %w", databaseID, instanceID, projectID, err)
	}

	return database, nil
}

// NewSpannerServiceE creates a Spanner service authenticated the same way every other client in
// this module is.
// The ctx parameter supports cancellation and timeouts.
func NewSpannerServiceE(t testing.TestingT, ctx context.Context) (*spanner.Service, error) {
	return spanner.NewService(ctx, append(withOptions(), option.WithScopes(spanner.CloudPlatformScope))...)
}
