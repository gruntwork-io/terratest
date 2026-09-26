package gcp

import (
	"context"
	"errors"
	"fmt"

	"github.com/gruntwork-io/terratest/modules/core/v2/logger"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/gkebackup/v1"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
)

// GetGKEBackupPlanAttrs returns the settings Google Cloud holds for the given backup plan, so a test can assert on
// what was actually created rather than only that it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetGKEBackupPlanAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, planID string) *gkebackup.BackupPlan {
	plan, err := GetGKEBackupPlanAttrsE(t, ctx, projectID, location, planID)
	require.NoError(t, err)

	return plan
}

// GetGKEBackupPlanAttrsE returns the settings Google Cloud holds for the given backup plan.
// The ctx parameter supports cancellation and timeouts.
func GetGKEBackupPlanAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, planID string) (*gkebackup.BackupPlan, error) {
	logger.Default.Logf(t, "Getting settings for backup plan %s in %s in project %s", planID, location, projectID)

	service, err := NewGKEBackupServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetGKEBackupPlanAttrsWithClient(ctx, service, projectID, location, planID)
}

// GetGKEBackupPlanAttrsWithClient returns the settings Google Cloud holds for the given backup plan using the supplied
// *gkebackup.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see gkebackup_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetGKEBackupPlanAttrsWithClient(ctx context.Context, service *gkebackup.Service, projectID string, location string, planID string) (*gkebackup.BackupPlan, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/backupPlans/%s", projectID, location, planID)

	plan, err := service.Projects.Locations.BackupPlans.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the backup plan %s does not exist in %s in project %s", planID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for backup plan %s in %s in project %s: %w", planID, location, projectID, err)
	}

	return plan, nil
}

// GetGKERestorePlanAttrs returns the settings Google Cloud holds for the given restore plan, so a test can assert on
// what was actually created rather than only that it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetGKERestorePlanAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, planID string) *gkebackup.RestorePlan {
	plan, err := GetGKERestorePlanAttrsE(t, ctx, projectID, location, planID)
	require.NoError(t, err)

	return plan
}

// GetGKERestorePlanAttrsE returns the settings Google Cloud holds for the given restore plan.
// The ctx parameter supports cancellation and timeouts.
func GetGKERestorePlanAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, planID string) (*gkebackup.RestorePlan, error) {
	logger.Default.Logf(t, "Getting settings for restore plan %s in %s in project %s", planID, location, projectID)

	service, err := NewGKEBackupServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetGKERestorePlanAttrsWithClient(ctx, service, projectID, location, planID)
}

// GetGKERestorePlanAttrsWithClient returns the settings Google Cloud holds for the given restore plan using the supplied
// *gkebackup.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see gkebackup_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetGKERestorePlanAttrsWithClient(ctx context.Context, service *gkebackup.Service, projectID string, location string, planID string) (*gkebackup.RestorePlan, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/restorePlans/%s", projectID, location, planID)

	plan, err := service.Projects.Locations.RestorePlans.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the restore plan %s does not exist in %s in project %s", planID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for restore plan %s in %s in project %s: %w", planID, location, projectID, err)
	}

	return plan, nil
}

// NewGKEBackupServiceE creates a Backup for GKE service authenticated the same way every other client in this
// module is.
// The ctx parameter supports cancellation and timeouts.
func NewGKEBackupServiceE(t testing.TestingT, ctx context.Context) (*gkebackup.Service, error) {
	return gkebackup.NewService(ctx, append(withOptions(), option.WithScopes(gkebackup.CloudPlatformScope))...)
}
