package gcp

import (
	"context"
	"errors"
	"fmt"

	"github.com/gruntwork-io/terratest/modules/core/v2/logger"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/alloydb/v1"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
)

// GetAlloyDBClusterAttrs returns the settings Google Cloud holds for the given AlloyDB cluster, so
// a test can assert on what was actually created rather than only that it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetAlloyDBClusterAttrs(t testing.TestingT, ctx context.Context, projectID string, region string, clusterID string) *alloydb.Cluster {
	cluster, err := GetAlloyDBClusterAttrsE(t, ctx, projectID, region, clusterID)
	require.NoError(t, err)

	return cluster
}

// GetAlloyDBClusterAttrsE returns the settings Google Cloud holds for the given AlloyDB cluster.
// The ctx parameter supports cancellation and timeouts.
func GetAlloyDBClusterAttrsE(t testing.TestingT, ctx context.Context, projectID string, region string, clusterID string) (*alloydb.Cluster, error) {
	logger.Default.Logf(t, "Getting settings for AlloyDB cluster %s in region %s in project %s", clusterID, region, projectID)

	service, err := NewAlloyDBServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetAlloyDBClusterAttrsWithClient(ctx, service, projectID, region, clusterID)
}

// GetAlloyDBClusterAttrsWithClient returns the settings Google Cloud holds for the given AlloyDB
// cluster using the supplied *alloydb.Service. Prefer this variant in unit tests where the service
// is backed by an httptest fake server (see alloydb_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetAlloyDBClusterAttrsWithClient(ctx context.Context, service *alloydb.Service, projectID string, region string, clusterID string) (*alloydb.Cluster, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/clusters/%s", projectID, region, clusterID)

	cluster, err := service.Projects.Locations.Clusters.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the AlloyDB cluster %s does not exist in region %s in project %s", clusterID, region, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for AlloyDB cluster %s in region %s in project %s: %w", clusterID, region, projectID, err)
	}

	return cluster, nil
}

// GetAlloyDBBackupAttrs returns the settings Google Cloud holds for the AlloyDB backup, so a test can assert on what was
// actually created rather than only that it exists.
// A backup names the cluster it was taken from, which is the only way a test can tell one backup of
// two clusters from the other.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetAlloyDBBackupAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, backupID string) *alloydb.Backup {
	result, err := GetAlloyDBBackupAttrsE(t, ctx, projectID, location, backupID)
	require.NoError(t, err)

	return result
}

// GetAlloyDBBackupAttrsE returns the settings Google Cloud holds for the AlloyDB backup.
// The ctx parameter supports cancellation and timeouts.
func GetAlloyDBBackupAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, backupID string) (*alloydb.Backup, error) {
	logger.Default.Logf(t, "Getting settings for AlloyDB backup %s in %s in project %s", backupID, location, projectID)

	service, err := NewAlloyDBServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetAlloyDBBackupAttrsWithClient(ctx, service, projectID, location, backupID)
}

// GetAlloyDBBackupAttrsWithClient returns the settings Google Cloud holds for the AlloyDB backup using the supplied
// *alloydb.Service. Prefer this variant in unit tests where the service is backed by an httptest
// fake server (see alloydb_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetAlloyDBBackupAttrsWithClient(ctx context.Context, service *alloydb.Service, projectID string, location string, backupID string) (*alloydb.Backup, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/backups/%s", projectID, location, backupID)

	result, err := service.Projects.Locations.Backups.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the AlloyDB backup %s in %s in project %s does not exist", backupID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for AlloyDB backup %s in %s in project %s: %w", backupID, location, projectID, err)
	}

	return result, nil
}

// GetAlloyDBUserAttrs returns the settings Google Cloud holds for the AlloyDB user, so a test can assert on what was
// actually created rather than only that it exists.
// The database roles a user holds are reported here and nowhere else, and the password never is, so
// this read is how a test asserts on the grant without reading a secret.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetAlloyDBUserAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, clusterID string, userID string) *alloydb.User {
	result, err := GetAlloyDBUserAttrsE(t, ctx, projectID, location, clusterID, userID)
	require.NoError(t, err)

	return result
}

// GetAlloyDBUserAttrsE returns the settings Google Cloud holds for the AlloyDB user.
// The ctx parameter supports cancellation and timeouts.
func GetAlloyDBUserAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, clusterID string, userID string) (*alloydb.User, error) {
	logger.Default.Logf(t, "Getting settings for AlloyDB user %s of cluster %s in %s in project %s", userID, clusterID, location, projectID)

	service, err := NewAlloyDBServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetAlloyDBUserAttrsWithClient(ctx, service, projectID, location, clusterID, userID)
}

// GetAlloyDBUserAttrsWithClient returns the settings Google Cloud holds for the AlloyDB user using the supplied
// *alloydb.Service. Prefer this variant in unit tests where the service is backed by an httptest
// fake server (see alloydb_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetAlloyDBUserAttrsWithClient(ctx context.Context, service *alloydb.Service, projectID string, location string, clusterID string, userID string) (*alloydb.User, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/clusters/%s/users/%s", projectID, location, clusterID, userID)

	result, err := service.Projects.Locations.Clusters.Users.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the AlloyDB user %s of cluster %s in %s in project %s does not exist", userID, clusterID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for AlloyDB user %s of cluster %s in %s in project %s: %w", userID, clusterID, location, projectID, err)
	}

	return result, nil
}

// NewAlloyDBServiceE creates a AlloyDB service authenticated the same way every other client in
// this module is.
// The ctx parameter supports cancellation and timeouts.
func NewAlloyDBServiceE(t testing.TestingT, ctx context.Context) (*alloydb.Service, error) {
	return alloydb.NewService(ctx, append(withOptions(), option.WithScopes(alloydb.CloudPlatformScope))...)
}
