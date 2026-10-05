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

// GetFilestoreSnapshotAttrs returns the settings Google Cloud holds for the given Filestore snapshot, so
// a test can assert on how it was described. A snapshot belongs to an instance, so the caller names both.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetFilestoreSnapshotAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, instanceID string, snapshotID string) *file.Snapshot {
	snapshot, err := GetFilestoreSnapshotAttrsE(t, ctx, projectID, location, instanceID, snapshotID)
	require.NoError(t, err)

	return snapshot
}

// GetFilestoreSnapshotAttrsE returns the settings Google Cloud holds for the given Filestore snapshot.
// The ctx parameter supports cancellation and timeouts.
func GetFilestoreSnapshotAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, instanceID string, snapshotID string) (*file.Snapshot, error) {
	logger.Default.Logf(t, "Getting settings for Filestore snapshot %s of instance %s in %s in project %s", snapshotID, instanceID, location, projectID)

	service, err := NewFilestoreServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetFilestoreSnapshotAttrsWithClient(ctx, service, projectID, location, instanceID, snapshotID)
}

// GetFilestoreSnapshotAttrsWithClient returns the settings Google Cloud holds for the given Filestore
// snapshot using the supplied *file.Service. Prefer this variant in unit tests where the service is
// backed by an httptest fake server (see filestore_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetFilestoreSnapshotAttrsWithClient(ctx context.Context, service *file.Service, projectID string, location string, instanceID string, snapshotID string) (*file.Snapshot, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/instances/%s/snapshots/%s", projectID, location, instanceID, snapshotID)

	snapshot, err := service.Projects.Locations.Instances.Snapshots.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Filestore snapshot %s of instance %s in %s in project %s does not exist", snapshotID, instanceID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Filestore snapshot %s of instance %s in %s in project %s: %w", snapshotID, instanceID, location, projectID, err)
	}

	return snapshot, nil
}

// GetFilestoreBackupAttrs returns the settings Google Cloud holds for the given Filestore backup, so a
// test can assert on which share it was taken from. A backup sits beside the instance rather than inside
// it, so its name carries only a location.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetFilestoreBackupAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, backupID string) *file.Backup {
	backup, err := GetFilestoreBackupAttrsE(t, ctx, projectID, location, backupID)
	require.NoError(t, err)

	return backup
}

// GetFilestoreBackupAttrsE returns the settings Google Cloud holds for the given Filestore backup.
// The ctx parameter supports cancellation and timeouts.
func GetFilestoreBackupAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, backupID string) (*file.Backup, error) {
	logger.Default.Logf(t, "Getting settings for Filestore backup %s in %s in project %s", backupID, location, projectID)

	service, err := NewFilestoreServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetFilestoreBackupAttrsWithClient(ctx, service, projectID, location, backupID)
}

// GetFilestoreBackupAttrsWithClient returns the settings Google Cloud holds for the given Filestore
// backup using the supplied *file.Service. Prefer this variant in unit tests where the service is backed
// by an httptest fake server (see filestore_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetFilestoreBackupAttrsWithClient(ctx context.Context, service *file.Service, projectID string, location string, backupID string) (*file.Backup, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/backups/%s", projectID, location, backupID)

	backup, err := service.Projects.Locations.Backups.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Filestore backup %s in %s in project %s does not exist", backupID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Filestore backup %s in %s in project %s: %w", backupID, location, projectID, err)
	}

	return backup, nil
}

// NewFilestoreServiceE creates a Filestore service authenticated the same way every other client in
// this module is.
// The ctx parameter supports cancellation and timeouts.
func NewFilestoreServiceE(t testing.TestingT, ctx context.Context) (*file.Service, error) {
	return file.NewService(ctx, append(withOptions(), option.WithScopes(file.CloudPlatformScope))...)
}
