package gcp

import (
	"context"
	"errors"
	"fmt"

	"github.com/gruntwork-io/terratest/modules/core/v2/logger"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/backupdr/v1"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
)

// GetBackupDRBackupVaultAttrs returns the settings Google Cloud holds for the Backup and DR backup vault, so a test can assert on what
// was actually created rather than only that it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetBackupDRBackupVaultAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, vaultID string) *backupdr.BackupVault {
	result, err := GetBackupDRBackupVaultAttrsE(t, ctx, projectID, location, vaultID)
	require.NoError(t, err)

	return result
}

// GetBackupDRBackupVaultAttrsE returns the settings Google Cloud holds for the Backup and DR backup vault.
// The ctx parameter supports cancellation and timeouts.
func GetBackupDRBackupVaultAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, vaultID string) (*backupdr.BackupVault, error) {
	logger.Default.Logf(t, "Getting settings for Backup and DR backup vault %s in %s in project %s", vaultID, location, projectID)

	service, err := NewBackupDRServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetBackupDRBackupVaultAttrsWithClient(ctx, service, projectID, location, vaultID)
}

// GetBackupDRBackupVaultAttrsWithClient returns the settings Google Cloud holds for the Backup and DR backup vault using the supplied
// *backupdr.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see backupdr_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetBackupDRBackupVaultAttrsWithClient(ctx context.Context, service *backupdr.Service, projectID string, location string, vaultID string) (*backupdr.BackupVault, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/backupVaults/%s", projectID, location, vaultID)

	result, err := service.Projects.Locations.BackupVaults.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Backup and DR backup vault %s in %s in project %s does not exist", vaultID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Backup and DR backup vault %s in %s in project %s: %w", vaultID, location, projectID, err)
	}

	return result, nil
}

// GetBackupDRBackupPlanAttrs returns the settings Google Cloud holds for the Backup and DR backup plan, so a test can assert on what
// was actually created rather than only that it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetBackupDRBackupPlanAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, planID string) *backupdr.BackupPlan {
	result, err := GetBackupDRBackupPlanAttrsE(t, ctx, projectID, location, planID)
	require.NoError(t, err)

	return result
}

// GetBackupDRBackupPlanAttrsE returns the settings Google Cloud holds for the Backup and DR backup plan.
// The ctx parameter supports cancellation and timeouts.
func GetBackupDRBackupPlanAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, planID string) (*backupdr.BackupPlan, error) {
	logger.Default.Logf(t, "Getting settings for Backup and DR backup plan %s in %s in project %s", planID, location, projectID)

	service, err := NewBackupDRServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetBackupDRBackupPlanAttrsWithClient(ctx, service, projectID, location, planID)
}

// GetBackupDRBackupPlanAttrsWithClient returns the settings Google Cloud holds for the Backup and DR backup plan using the supplied
// *backupdr.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see backupdr_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetBackupDRBackupPlanAttrsWithClient(ctx context.Context, service *backupdr.Service, projectID string, location string, planID string) (*backupdr.BackupPlan, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/backupPlans/%s", projectID, location, planID)

	result, err := service.Projects.Locations.BackupPlans.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Backup and DR backup plan %s in %s in project %s does not exist", planID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Backup and DR backup plan %s in %s in project %s: %w", planID, location, projectID, err)
	}

	return result, nil
}

// NewBackupDRServiceE creates a Backup and DR service authenticated the same way every other client in
// this module is.
// The ctx parameter supports cancellation and timeouts.
func NewBackupDRServiceE(t testing.TestingT, ctx context.Context) (*backupdr.Service, error) {
	return backupdr.NewService(ctx, append(withOptions(), option.WithScopes(backupdr.CloudPlatformScope))...)
}
