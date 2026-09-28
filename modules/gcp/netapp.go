package gcp

import (
	"context"
	"errors"
	"fmt"

	"github.com/gruntwork-io/terratest/modules/core/v2/logger"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/netapp/v1"
	"google.golang.org/api/option"
)

// GetNetAppBackupVaultAttrs returns the settings Google Cloud holds for the NetApp backup vault, so a test can assert on what
// was actually created rather than only that it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetNetAppBackupVaultAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, vaultID string) *netapp.BackupVault {
	result, err := GetNetAppBackupVaultAttrsE(t, ctx, projectID, location, vaultID)
	require.NoError(t, err)

	return result
}

// GetNetAppBackupVaultAttrsE returns the settings Google Cloud holds for the NetApp backup vault.
// The ctx parameter supports cancellation and timeouts.
func GetNetAppBackupVaultAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, vaultID string) (*netapp.BackupVault, error) {
	logger.Default.Logf(t, "Getting settings for NetApp backup vault %s in %s in project %s", vaultID, location, projectID)

	service, err := NewNetAppServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetNetAppBackupVaultAttrsWithClient(ctx, service, projectID, location, vaultID)
}

// GetNetAppBackupVaultAttrsWithClient returns the settings Google Cloud holds for the NetApp backup vault using the supplied
// *netapp.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see netapp_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetNetAppBackupVaultAttrsWithClient(ctx context.Context, service *netapp.Service, projectID string, location string, vaultID string) (*netapp.BackupVault, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/backupVaults/%s", projectID, location, vaultID)

	result, err := service.Projects.Locations.BackupVaults.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the NetApp backup vault %s in %s in project %s does not exist", vaultID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for NetApp backup vault %s in %s in project %s: %w", vaultID, location, projectID, err)
	}

	return result, nil
}

// GetNetAppBackupPolicyAttrs returns the settings Google Cloud holds for the NetApp backup policy, so a test can assert on what
// was actually created rather than only that it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetNetAppBackupPolicyAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, policyID string) *netapp.BackupPolicy {
	result, err := GetNetAppBackupPolicyAttrsE(t, ctx, projectID, location, policyID)
	require.NoError(t, err)

	return result
}

// GetNetAppBackupPolicyAttrsE returns the settings Google Cloud holds for the NetApp backup policy.
// The ctx parameter supports cancellation and timeouts.
func GetNetAppBackupPolicyAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, policyID string) (*netapp.BackupPolicy, error) {
	logger.Default.Logf(t, "Getting settings for NetApp backup policy %s in %s in project %s", policyID, location, projectID)

	service, err := NewNetAppServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetNetAppBackupPolicyAttrsWithClient(ctx, service, projectID, location, policyID)
}

// GetNetAppBackupPolicyAttrsWithClient returns the settings Google Cloud holds for the NetApp backup policy using the supplied
// *netapp.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see netapp_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetNetAppBackupPolicyAttrsWithClient(ctx context.Context, service *netapp.Service, projectID string, location string, policyID string) (*netapp.BackupPolicy, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/backupPolicies/%s", projectID, location, policyID)

	result, err := service.Projects.Locations.BackupPolicies.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the NetApp backup policy %s in %s in project %s does not exist", policyID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for NetApp backup policy %s in %s in project %s: %w", policyID, location, projectID, err)
	}

	return result, nil
}

// NewNetAppServiceE creates a NetApp Volumes service authenticated the same way every other client in
// this module is.
// The ctx parameter supports cancellation and timeouts.
func NewNetAppServiceE(t testing.TestingT, ctx context.Context) (*netapp.Service, error) {
	return netapp.NewService(ctx, append(withOptions(), option.WithScopes(netapp.CloudPlatformScope))...)
}
