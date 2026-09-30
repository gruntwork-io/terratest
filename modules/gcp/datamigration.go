package gcp

import (
	"context"
	"errors"
	"fmt"

	"github.com/gruntwork-io/terratest/modules/core/v2/logger"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/datamigration/v1"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
)

// GetConnectionProfileAttrs returns the settings Google Cloud holds for the given Database Migration
// connection profile, so a test can assert on which database it points at rather than only that it
// exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetConnectionProfileAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, profileID string) *datamigration.ConnectionProfile {
	profile, err := GetConnectionProfileAttrsE(t, ctx, projectID, location, profileID)
	require.NoError(t, err)

	return profile
}

// GetConnectionProfileAttrsE returns the settings Google Cloud holds for the given Database
// Migration connection profile.
// The ctx parameter supports cancellation and timeouts.
func GetConnectionProfileAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, profileID string) (*datamigration.ConnectionProfile, error) {
	logger.Default.Logf(t, "Getting settings for connection profile %s in %s in project %s", profileID, location, projectID)

	service, err := NewDatabaseMigrationServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetConnectionProfileAttrsWithClient(ctx, service, projectID, location, profileID)
}

// GetConnectionProfileAttrsWithClient returns the settings Google Cloud holds for the given Database
// Migration connection profile using the supplied *datamigration.Service. Prefer this variant in
// unit tests where the service is backed by an httptest fake server (see datamigration_test.go for
// the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetConnectionProfileAttrsWithClient(ctx context.Context, service *datamigration.Service, projectID string, location string, profileID string) (*datamigration.ConnectionProfile, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/connectionProfiles/%s", projectID, location, profileID)

	profile, err := service.Projects.Locations.ConnectionProfiles.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the connection profile %s in %s in project %s does not exist", profileID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for connection profile %s in %s in project %s: %w", profileID, location, projectID, err)
	}

	return profile, nil
}

// NewDatabaseMigrationServiceE creates a Database Migration service authenticated the same way every
// other client in this module is.
// The ctx parameter supports cancellation and timeouts.
func NewDatabaseMigrationServiceE(t testing.TestingT, ctx context.Context) (*datamigration.Service, error) {
	return datamigration.NewService(ctx, append(withOptions(), option.WithScopes(datamigration.CloudPlatformScope))...)
}
