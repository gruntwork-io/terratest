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

// GetNetAppActiveDirectoryAttrs returns the settings Google Cloud holds for the given NetApp Active Directory connection, so a test can assert on what was
// actually created rather than only that it exists.
// A volume can only speak SMB if the pool it is in points at a directory, so the domain and the DNS it uses are the whole point.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetNetAppActiveDirectoryAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, id string) *netapp.ActiveDirectory {
	attrs, err := GetNetAppActiveDirectoryAttrsE(t, ctx, projectID, location, id)
	require.NoError(t, err)

	return attrs
}

// GetNetAppActiveDirectoryAttrsE returns the settings Google Cloud holds for the given NetApp Active Directory connection.
// The ctx parameter supports cancellation and timeouts.
func GetNetAppActiveDirectoryAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, id string) (*netapp.ActiveDirectory, error) {
	logger.Default.Logf(t, "Getting settings for NetApp Active Directory connection %s in %s in project %s", id, location, projectID)

	service, err := NewNetAppServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetNetAppActiveDirectoryAttrsWithClient(ctx, service, projectID, location, id)
}

// GetNetAppActiveDirectoryAttrsWithClient returns the settings Google Cloud holds for the given NetApp Active Directory connection using the supplied
// *netapp.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see netapp_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetNetAppActiveDirectoryAttrsWithClient(ctx context.Context, service *netapp.Service, projectID string, location string, id string) (*netapp.ActiveDirectory, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/activeDirectories/%s", projectID, location, id)

	attrs, err := service.Projects.Locations.ActiveDirectories.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the NetApp Active Directory connection %s in %s in project %s does not exist", id, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for NetApp Active Directory connection %s in %s in project %s: %w", id, location, projectID, err)
	}

	return attrs, nil
}

// GetNetAppStoragePoolAttrs returns the settings Google Cloud holds for the given NetApp storage pool, so a test can assert on what was
// actually created rather than only that it exists.
// A pool buys the capacity and service level the volumes inside it draw on, so those two numbers decide what any volume in it can do.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetNetAppStoragePoolAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, id string) *netapp.StoragePool {
	attrs, err := GetNetAppStoragePoolAttrsE(t, ctx, projectID, location, id)
	require.NoError(t, err)

	return attrs
}

// GetNetAppStoragePoolAttrsE returns the settings Google Cloud holds for the given NetApp storage pool.
// The ctx parameter supports cancellation and timeouts.
func GetNetAppStoragePoolAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, id string) (*netapp.StoragePool, error) {
	logger.Default.Logf(t, "Getting settings for NetApp storage pool %s in %s in project %s", id, location, projectID)

	service, err := NewNetAppServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetNetAppStoragePoolAttrsWithClient(ctx, service, projectID, location, id)
}

// GetNetAppStoragePoolAttrsWithClient returns the settings Google Cloud holds for the given NetApp storage pool using the supplied
// *netapp.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see netapp_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetNetAppStoragePoolAttrsWithClient(ctx context.Context, service *netapp.Service, projectID string, location string, id string) (*netapp.StoragePool, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/storagePools/%s", projectID, location, id)

	attrs, err := service.Projects.Locations.StoragePools.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the NetApp storage pool %s in %s in project %s does not exist", id, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for NetApp storage pool %s in %s in project %s: %w", id, location, projectID, err)
	}

	return attrs, nil
}

// GetNetAppVolumeAttrs returns the settings Google Cloud holds for the given NetApp volume, so a test can assert on what was
// actually created rather than only that it exists.
// The volume is the share clients mount, so its size, its protocol and who may mount it are what callers care about.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetNetAppVolumeAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, id string) *netapp.Volume {
	attrs, err := GetNetAppVolumeAttrsE(t, ctx, projectID, location, id)
	require.NoError(t, err)

	return attrs
}

// GetNetAppVolumeAttrsE returns the settings Google Cloud holds for the given NetApp volume.
// The ctx parameter supports cancellation and timeouts.
func GetNetAppVolumeAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, id string) (*netapp.Volume, error) {
	logger.Default.Logf(t, "Getting settings for NetApp volume %s in %s in project %s", id, location, projectID)

	service, err := NewNetAppServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetNetAppVolumeAttrsWithClient(ctx, service, projectID, location, id)
}

// GetNetAppVolumeAttrsWithClient returns the settings Google Cloud holds for the given NetApp volume using the supplied
// *netapp.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see netapp_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetNetAppVolumeAttrsWithClient(ctx context.Context, service *netapp.Service, projectID string, location string, id string) (*netapp.Volume, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/volumes/%s", projectID, location, id)

	attrs, err := service.Projects.Locations.Volumes.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the NetApp volume %s in %s in project %s does not exist", id, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for NetApp volume %s in %s in project %s: %w", id, location, projectID, err)
	}

	return attrs, nil
}

// GetNetAppVolumeSnapshotAttrs returns the settings Google Cloud holds for the given NetApp volume snapshot, so a test can assert on what was
// actually created rather than only that it exists.
// A snapshot is a point a volume can be restored to, so whether it exists at all is the thing worth asserting.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetNetAppVolumeSnapshotAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, volumeID string, id string) *netapp.Snapshot {
	attrs, err := GetNetAppVolumeSnapshotAttrsE(t, ctx, projectID, location, volumeID, id)
	require.NoError(t, err)

	return attrs
}

// GetNetAppVolumeSnapshotAttrsE returns the settings Google Cloud holds for the given NetApp volume snapshot.
// The ctx parameter supports cancellation and timeouts.
func GetNetAppVolumeSnapshotAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, volumeID string, id string) (*netapp.Snapshot, error) {
	logger.Default.Logf(t, "Getting settings for NetApp volume snapshot %s of volume %s in %s in project %s", id, volumeID, location, projectID)

	service, err := NewNetAppServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetNetAppVolumeSnapshotAttrsWithClient(ctx, service, projectID, location, volumeID, id)
}

// GetNetAppVolumeSnapshotAttrsWithClient returns the settings Google Cloud holds for the given NetApp volume snapshot using the supplied
// *netapp.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see netapp_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetNetAppVolumeSnapshotAttrsWithClient(ctx context.Context, service *netapp.Service, projectID string, location string, volumeID string, id string) (*netapp.Snapshot, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/volumes/%s/snapshots/%s", projectID, location, volumeID, id)

	attrs, err := service.Projects.Locations.Volumes.Snapshots.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the NetApp volume snapshot %s of volume %s in %s in project %s does not exist", id, volumeID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for NetApp volume snapshot %s of volume %s in %s in project %s: %w", id, volumeID, location, projectID, err)
	}

	return attrs, nil
}

// GetNetAppVolumeReplicationAttrs returns the settings Google Cloud holds for the given NetApp volume replication, so a test can assert on what was
// actually created rather than only that it exists.
// A replication copies a volume to another region, so its schedule and its mirror state say whether the copy is current.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetNetAppVolumeReplicationAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, volumeID string, id string) *netapp.Replication {
	attrs, err := GetNetAppVolumeReplicationAttrsE(t, ctx, projectID, location, volumeID, id)
	require.NoError(t, err)

	return attrs
}

// GetNetAppVolumeReplicationAttrsE returns the settings Google Cloud holds for the given NetApp volume replication.
// The ctx parameter supports cancellation and timeouts.
func GetNetAppVolumeReplicationAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, volumeID string, id string) (*netapp.Replication, error) {
	logger.Default.Logf(t, "Getting settings for NetApp volume replication %s of volume %s in %s in project %s", id, volumeID, location, projectID)

	service, err := NewNetAppServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetNetAppVolumeReplicationAttrsWithClient(ctx, service, projectID, location, volumeID, id)
}

// GetNetAppVolumeReplicationAttrsWithClient returns the settings Google Cloud holds for the given NetApp volume replication using the supplied
// *netapp.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see netapp_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetNetAppVolumeReplicationAttrsWithClient(ctx context.Context, service *netapp.Service, projectID string, location string, volumeID string, id string) (*netapp.Replication, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/volumes/%s/replications/%s", projectID, location, volumeID, id)

	attrs, err := service.Projects.Locations.Volumes.Replications.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the NetApp volume replication %s of volume %s in %s in project %s does not exist", id, volumeID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for NetApp volume replication %s of volume %s in %s in project %s: %w", id, volumeID, location, projectID, err)
	}

	return attrs, nil
}

// GetNetAppVolumeQuotaRuleAttrs returns the settings Google Cloud holds for the given NetApp volume quota rule, so a test can assert on what was
// actually created rather than only that it exists.
// A rule caps how much of a volume one user or group may fill, so the target it names and the limit are the point.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetNetAppVolumeQuotaRuleAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, volumeID string, id string) *netapp.QuotaRule {
	attrs, err := GetNetAppVolumeQuotaRuleAttrsE(t, ctx, projectID, location, volumeID, id)
	require.NoError(t, err)

	return attrs
}

// GetNetAppVolumeQuotaRuleAttrsE returns the settings Google Cloud holds for the given NetApp volume quota rule.
// The ctx parameter supports cancellation and timeouts.
func GetNetAppVolumeQuotaRuleAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, volumeID string, id string) (*netapp.QuotaRule, error) {
	logger.Default.Logf(t, "Getting settings for NetApp volume quota rule %s of volume %s in %s in project %s", id, volumeID, location, projectID)

	service, err := NewNetAppServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetNetAppVolumeQuotaRuleAttrsWithClient(ctx, service, projectID, location, volumeID, id)
}

// GetNetAppVolumeQuotaRuleAttrsWithClient returns the settings Google Cloud holds for the given NetApp volume quota rule using the supplied
// *netapp.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see netapp_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetNetAppVolumeQuotaRuleAttrsWithClient(ctx context.Context, service *netapp.Service, projectID string, location string, volumeID string, id string) (*netapp.QuotaRule, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/volumes/%s/quotaRules/%s", projectID, location, volumeID, id)

	attrs, err := service.Projects.Locations.Volumes.QuotaRules.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the NetApp volume quota rule %s of volume %s in %s in project %s does not exist", id, volumeID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for NetApp volume quota rule %s of volume %s in %s in project %s: %w", id, volumeID, location, projectID, err)
	}

	return attrs, nil
}

// GetNetAppBackupAttrs returns the settings Google Cloud holds for the given NetApp backup, so a test can assert on what was
// actually created rather than only that it exists.
// A backup lives in a vault rather than beside the volume, so the volume it came from and its state are what say it is usable.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetNetAppBackupAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, vaultID string, id string) *netapp.Backup {
	attrs, err := GetNetAppBackupAttrsE(t, ctx, projectID, location, vaultID, id)
	require.NoError(t, err)

	return attrs
}

// GetNetAppBackupAttrsE returns the settings Google Cloud holds for the given NetApp backup.
// The ctx parameter supports cancellation and timeouts.
func GetNetAppBackupAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, vaultID string, id string) (*netapp.Backup, error) {
	logger.Default.Logf(t, "Getting settings for NetApp backup %s in vault %s in %s in project %s", id, vaultID, location, projectID)

	service, err := NewNetAppServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetNetAppBackupAttrsWithClient(ctx, service, projectID, location, vaultID, id)
}

// GetNetAppBackupAttrsWithClient returns the settings Google Cloud holds for the given NetApp backup using the supplied
// *netapp.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see netapp_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetNetAppBackupAttrsWithClient(ctx context.Context, service *netapp.Service, projectID string, location string, vaultID string, id string) (*netapp.Backup, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/backupVaults/%s/backups/%s", projectID, location, vaultID, id)

	attrs, err := service.Projects.Locations.BackupVaults.Backups.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the NetApp backup %s in vault %s in %s in project %s does not exist", id, vaultID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for NetApp backup %s in vault %s in %s in project %s: %w", id, vaultID, location, projectID, err)
	}

	return attrs, nil
}

// GetNetAppHostGroupAttrs returns the settings Google Cloud holds for the given NetApp host group, so a test can assert on
// what was actually created rather than only that it exists. A host group names the initiators a volume may be exported to; it exports nothing by itself.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetNetAppHostGroupAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, groupID string) *netapp.HostGroup {
	result, err := GetNetAppHostGroupAttrsE(t, ctx, projectID, location, groupID)
	require.NoError(t, err)

	return result
}

// GetNetAppHostGroupAttrsE returns the settings Google Cloud holds for the given NetApp host group.
// The ctx parameter supports cancellation and timeouts.
func GetNetAppHostGroupAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, groupID string) (*netapp.HostGroup, error) {
	logger.Default.Logf(t, "Getting settings for NetApp host group %s in project %s", groupID, projectID)

	service, err := NewNetAppServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetNetAppHostGroupAttrsWithClient(ctx, service, projectID, location, groupID)
}

// GetNetAppHostGroupAttrsWithClient returns the settings Google Cloud holds for the given NetApp host group using the
// supplied *netapp.Service. Prefer this variant in unit tests where the service is backed by an
// httptest fake server (see netapp_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetNetAppHostGroupAttrsWithClient(ctx context.Context, service *netapp.Service, projectID string, location string, groupID string) (*netapp.HostGroup, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/hostGroups/%s", projectID, location, groupID)

	result, err := service.Projects.Locations.HostGroups.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the NetApp host group %s in project %s does not exist", groupID, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for NetApp host group %s in project %s: %w", groupID, projectID, err)
	}

	return result, nil
}

// GetNetAppKmsConfigAttrs returns the settings Google Cloud holds for the given NetApp KMS config, so a test can assert on
// what was actually created rather than only that it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetNetAppKmsConfigAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, configID string) *netapp.KmsConfig {
	result, err := GetNetAppKmsConfigAttrsE(t, ctx, projectID, location, configID)
	require.NoError(t, err)

	return result
}

// GetNetAppKmsConfigAttrsE returns the settings Google Cloud holds for the given NetApp KMS config.
// The ctx parameter supports cancellation and timeouts.
func GetNetAppKmsConfigAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, configID string) (*netapp.KmsConfig, error) {
	logger.Default.Logf(t, "Getting settings for NetApp KMS config %s in project %s", configID, projectID)

	service, err := NewNetAppServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetNetAppKmsConfigAttrsWithClient(ctx, service, projectID, location, configID)
}

// GetNetAppKmsConfigAttrsWithClient returns the settings Google Cloud holds for the given NetApp KMS config using the
// supplied *netapp.Service. Prefer this variant in unit tests where the service is backed by an
// httptest fake server (see netapp_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetNetAppKmsConfigAttrsWithClient(ctx context.Context, service *netapp.Service, projectID string, location string, configID string) (*netapp.KmsConfig, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/kmsConfigs/%s", projectID, location, configID)

	result, err := service.Projects.Locations.KmsConfigs.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the NetApp KMS config %s in project %s does not exist", configID, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for NetApp KMS config %s in project %s: %w", configID, projectID, err)
	}

	return result, nil
}

// NewNetAppServiceE creates a NetApp Volumes service authenticated the same way every other client in
// this module is.
// The ctx parameter supports cancellation and timeouts.
func NewNetAppServiceE(t testing.TestingT, ctx context.Context) (*netapp.Service, error) {
	return netapp.NewService(ctx, append(withOptions(), option.WithScopes(netapp.CloudPlatformScope))...)
}
