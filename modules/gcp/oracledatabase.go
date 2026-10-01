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
	"google.golang.org/api/oracledatabase/v1"
)

// GetOracleAutonomousDatabaseAttrs returns the settings Google Cloud holds for the given Oracle autonomous database, so a test can assert on what was
// actually created rather than only that it exists.
// An autonomous database manages itself, so its workload type, its licence model and the network it sits on are what the module actually chose.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetOracleAutonomousDatabaseAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, databaseID string) *oracledatabase.AutonomousDatabase {
	attrs, err := GetOracleAutonomousDatabaseAttrsE(t, ctx, projectID, location, databaseID)
	require.NoError(t, err)

	return attrs
}

// GetOracleAutonomousDatabaseAttrsE returns the settings Google Cloud holds for the given Oracle autonomous database.
// The ctx parameter supports cancellation and timeouts.
func GetOracleAutonomousDatabaseAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, databaseID string) (*oracledatabase.AutonomousDatabase, error) {
	logger.Default.Logf(t, "Getting settings for Oracle autonomous database %s in %s in project %s", databaseID, location, projectID)

	service, err := NewOracleDatabaseServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetOracleAutonomousDatabaseAttrsWithClient(ctx, service, projectID, location, databaseID)
}

// GetOracleAutonomousDatabaseAttrsWithClient returns the settings Google Cloud holds for the given Oracle autonomous database using the supplied
// *oracledatabase.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see oracledatabase_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetOracleAutonomousDatabaseAttrsWithClient(ctx context.Context, service *oracledatabase.Service, projectID string, location string, databaseID string) (*oracledatabase.AutonomousDatabase, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/autonomousDatabases/%s", projectID, location, databaseID)

	attrs, err := service.Projects.Locations.AutonomousDatabases.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Oracle autonomous database %s in %s in project %s does not exist", databaseID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Oracle autonomous database %s in %s in project %s: %w", databaseID, location, projectID, err)
	}

	return attrs, nil
}

// GetOracleCloudExadataInfrastructureAttrs returns the settings Google Cloud holds for the given Oracle Exadata infrastructure, so a test can assert on what was
// actually created rather than only that it exists.
// The infrastructure is the hardware a VM cluster runs on, so its shape and compute count are what is billed.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetOracleCloudExadataInfrastructureAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, infrastructureID string) *oracledatabase.CloudExadataInfrastructure {
	attrs, err := GetOracleCloudExadataInfrastructureAttrsE(t, ctx, projectID, location, infrastructureID)
	require.NoError(t, err)

	return attrs
}

// GetOracleCloudExadataInfrastructureAttrsE returns the settings Google Cloud holds for the given Oracle Exadata infrastructure.
// The ctx parameter supports cancellation and timeouts.
func GetOracleCloudExadataInfrastructureAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, infrastructureID string) (*oracledatabase.CloudExadataInfrastructure, error) {
	logger.Default.Logf(t, "Getting settings for Oracle Exadata infrastructure %s in %s in project %s", infrastructureID, location, projectID)

	service, err := NewOracleDatabaseServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetOracleCloudExadataInfrastructureAttrsWithClient(ctx, service, projectID, location, infrastructureID)
}

// GetOracleCloudExadataInfrastructureAttrsWithClient returns the settings Google Cloud holds for the given Oracle Exadata infrastructure using the supplied
// *oracledatabase.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see oracledatabase_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetOracleCloudExadataInfrastructureAttrsWithClient(ctx context.Context, service *oracledatabase.Service, projectID string, location string, infrastructureID string) (*oracledatabase.CloudExadataInfrastructure, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/cloudExadataInfrastructures/%s", projectID, location, infrastructureID)

	attrs, err := service.Projects.Locations.CloudExadataInfrastructures.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Oracle Exadata infrastructure %s in %s in project %s does not exist", infrastructureID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Oracle Exadata infrastructure %s in %s in project %s: %w", infrastructureID, location, projectID, err)
	}

	return attrs, nil
}

// GetOracleCloudVMClusterAttrs returns the settings Google Cloud holds for the given Oracle Exadata VM cluster, so a test can assert on what was
// actually created rather than only that it exists.
// A VM cluster is where databases actually run, so its CPU count, its licence model and the infrastructure it names decide what it can do.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetOracleCloudVMClusterAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, clusterID string) *oracledatabase.CloudVmCluster {
	attrs, err := GetOracleCloudVMClusterAttrsE(t, ctx, projectID, location, clusterID)
	require.NoError(t, err)

	return attrs
}

// GetOracleCloudVMClusterAttrsE returns the settings Google Cloud holds for the given Oracle Exadata VM cluster.
// The ctx parameter supports cancellation and timeouts.
func GetOracleCloudVMClusterAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, clusterID string) (*oracledatabase.CloudVmCluster, error) {
	logger.Default.Logf(t, "Getting settings for Oracle Exadata VM cluster %s in %s in project %s", clusterID, location, projectID)

	service, err := NewOracleDatabaseServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetOracleCloudVMClusterAttrsWithClient(ctx, service, projectID, location, clusterID)
}

// GetOracleCloudVMClusterAttrsWithClient returns the settings Google Cloud holds for the given Oracle Exadata VM cluster using the supplied
// *oracledatabase.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see oracledatabase_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetOracleCloudVMClusterAttrsWithClient(ctx context.Context, service *oracledatabase.Service, projectID string, location string, clusterID string) (*oracledatabase.CloudVmCluster, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/cloudVmClusters/%s", projectID, location, clusterID)

	attrs, err := service.Projects.Locations.CloudVmClusters.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Oracle Exadata VM cluster %s in %s in project %s does not exist", clusterID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Oracle Exadata VM cluster %s in %s in project %s: %w", clusterID, location, projectID, err)
	}

	return attrs, nil
}

// GetOracleDBSystemAttrs returns the settings Google Cloud holds for the given Oracle database system, so a test can assert on what was
// actually created rather than only that it exists.
// A database system is a single-node Oracle deployment, so its shape and the network it sits on are the point.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetOracleDBSystemAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, systemID string) *oracledatabase.DbSystem {
	attrs, err := GetOracleDBSystemAttrsE(t, ctx, projectID, location, systemID)
	require.NoError(t, err)

	return attrs
}

// GetOracleDBSystemAttrsE returns the settings Google Cloud holds for the given Oracle database system.
// The ctx parameter supports cancellation and timeouts.
func GetOracleDBSystemAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, systemID string) (*oracledatabase.DbSystem, error) {
	logger.Default.Logf(t, "Getting settings for Oracle database system %s in %s in project %s", systemID, location, projectID)

	service, err := NewOracleDatabaseServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetOracleDBSystemAttrsWithClient(ctx, service, projectID, location, systemID)
}

// GetOracleDBSystemAttrsWithClient returns the settings Google Cloud holds for the given Oracle database system using the supplied
// *oracledatabase.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see oracledatabase_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetOracleDBSystemAttrsWithClient(ctx context.Context, service *oracledatabase.Service, projectID string, location string, systemID string) (*oracledatabase.DbSystem, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/dbSystems/%s", projectID, location, systemID)

	attrs, err := service.Projects.Locations.DbSystems.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Oracle database system %s in %s in project %s does not exist", systemID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Oracle database system %s in %s in project %s: %w", systemID, location, projectID, err)
	}

	return attrs, nil
}

// GetOracleExadbVMClusterAttrs returns the settings Google Cloud holds for the given Oracle Exascale VM cluster, so a test can assert on what was
// actually created rather than only that it exists.
// An Exascale cluster draws its storage from a vault rather than from attached disks, so the vault it names is what it runs on.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetOracleExadbVMClusterAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, clusterID string) *oracledatabase.ExadbVmCluster {
	attrs, err := GetOracleExadbVMClusterAttrsE(t, ctx, projectID, location, clusterID)
	require.NoError(t, err)

	return attrs
}

// GetOracleExadbVMClusterAttrsE returns the settings Google Cloud holds for the given Oracle Exascale VM cluster.
// The ctx parameter supports cancellation and timeouts.
func GetOracleExadbVMClusterAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, clusterID string) (*oracledatabase.ExadbVmCluster, error) {
	logger.Default.Logf(t, "Getting settings for Oracle Exascale VM cluster %s in %s in project %s", clusterID, location, projectID)

	service, err := NewOracleDatabaseServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetOracleExadbVMClusterAttrsWithClient(ctx, service, projectID, location, clusterID)
}

// GetOracleExadbVMClusterAttrsWithClient returns the settings Google Cloud holds for the given Oracle Exascale VM cluster using the supplied
// *oracledatabase.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see oracledatabase_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetOracleExadbVMClusterAttrsWithClient(ctx context.Context, service *oracledatabase.Service, projectID string, location string, clusterID string) (*oracledatabase.ExadbVmCluster, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/exadbVmClusters/%s", projectID, location, clusterID)

	attrs, err := service.Projects.Locations.ExadbVmClusters.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Oracle Exascale VM cluster %s in %s in project %s does not exist", clusterID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Oracle Exascale VM cluster %s in %s in project %s: %w", clusterID, location, projectID, err)
	}

	return attrs, nil
}

// GetOracleExascaleDBStorageVaultAttrs returns the settings Google Cloud holds for the given Oracle Exascale storage vault, so a test can assert on what was
// actually created rather than only that it exists.
// A vault is the storage an Exascale cluster draws on, so its size is what every database in it shares.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetOracleExascaleDBStorageVaultAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, vaultID string) *oracledatabase.ExascaleDbStorageVault {
	attrs, err := GetOracleExascaleDBStorageVaultAttrsE(t, ctx, projectID, location, vaultID)
	require.NoError(t, err)

	return attrs
}

// GetOracleExascaleDBStorageVaultAttrsE returns the settings Google Cloud holds for the given Oracle Exascale storage vault.
// The ctx parameter supports cancellation and timeouts.
func GetOracleExascaleDBStorageVaultAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, vaultID string) (*oracledatabase.ExascaleDbStorageVault, error) {
	logger.Default.Logf(t, "Getting settings for Oracle Exascale storage vault %s in %s in project %s", vaultID, location, projectID)

	service, err := NewOracleDatabaseServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetOracleExascaleDBStorageVaultAttrsWithClient(ctx, service, projectID, location, vaultID)
}

// GetOracleExascaleDBStorageVaultAttrsWithClient returns the settings Google Cloud holds for the given Oracle Exascale storage vault using the supplied
// *oracledatabase.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see oracledatabase_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetOracleExascaleDBStorageVaultAttrsWithClient(ctx context.Context, service *oracledatabase.Service, projectID string, location string, vaultID string) (*oracledatabase.ExascaleDbStorageVault, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/exascaleDbStorageVaults/%s", projectID, location, vaultID)

	attrs, err := service.Projects.Locations.ExascaleDbStorageVaults.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Oracle Exascale storage vault %s in %s in project %s does not exist", vaultID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Oracle Exascale storage vault %s in %s in project %s: %w", vaultID, location, projectID, err)
	}

	return attrs, nil
}

// GetOracleODBNetworkAttrs returns the settings Google Cloud holds for the given Oracle ODB network, so a test can assert on what was
// actually created rather than only that it exists.
// An ODB network is the VPC Oracle's own resources are reachable on, so the network it names decides what can connect.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetOracleODBNetworkAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, networkID string) *oracledatabase.OdbNetwork {
	attrs, err := GetOracleODBNetworkAttrsE(t, ctx, projectID, location, networkID)
	require.NoError(t, err)

	return attrs
}

// GetOracleODBNetworkAttrsE returns the settings Google Cloud holds for the given Oracle ODB network.
// The ctx parameter supports cancellation and timeouts.
func GetOracleODBNetworkAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, networkID string) (*oracledatabase.OdbNetwork, error) {
	logger.Default.Logf(t, "Getting settings for Oracle ODB network %s in %s in project %s", networkID, location, projectID)

	service, err := NewOracleDatabaseServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetOracleODBNetworkAttrsWithClient(ctx, service, projectID, location, networkID)
}

// GetOracleODBNetworkAttrsWithClient returns the settings Google Cloud holds for the given Oracle ODB network using the supplied
// *oracledatabase.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see oracledatabase_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetOracleODBNetworkAttrsWithClient(ctx context.Context, service *oracledatabase.Service, projectID string, location string, networkID string) (*oracledatabase.OdbNetwork, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/odbNetworks/%s", projectID, location, networkID)

	attrs, err := service.Projects.Locations.OdbNetworks.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Oracle ODB network %s in %s in project %s does not exist", networkID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Oracle ODB network %s in %s in project %s: %w", networkID, location, projectID, err)
	}

	return attrs, nil
}

// GetOracleODBSubnetAttrs returns the settings Google Cloud holds for the given Oracle ODB subnet, so a test can assert on what was
// actually created rather than only that it exists.
// A subnet carves the range Oracle's resources take addresses from, so the range and its purpose are the point.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetOracleODBSubnetAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, networkID string, subnetID string) *oracledatabase.OdbSubnet {
	attrs, err := GetOracleODBSubnetAttrsE(t, ctx, projectID, location, networkID, subnetID)
	require.NoError(t, err)

	return attrs
}

// GetOracleODBSubnetAttrsE returns the settings Google Cloud holds for the given Oracle ODB subnet.
// The ctx parameter supports cancellation and timeouts.
func GetOracleODBSubnetAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, networkID string, subnetID string) (*oracledatabase.OdbSubnet, error) {
	logger.Default.Logf(t, "Getting settings for Oracle ODB subnet %s in network %s in %s in project %s", subnetID, networkID, location, projectID)

	service, err := NewOracleDatabaseServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetOracleODBSubnetAttrsWithClient(ctx, service, projectID, location, networkID, subnetID)
}

// GetOracleODBSubnetAttrsWithClient returns the settings Google Cloud holds for the given Oracle ODB subnet using the supplied
// *oracledatabase.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see oracledatabase_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetOracleODBSubnetAttrsWithClient(ctx context.Context, service *oracledatabase.Service, projectID string, location string, networkID string, subnetID string) (*oracledatabase.OdbSubnet, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/odbNetworks/%s/odbSubnets/%s", projectID, location, networkID, subnetID)

	attrs, err := service.Projects.Locations.OdbNetworks.OdbSubnets.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Oracle ODB subnet %s in network %s in %s in project %s does not exist", subnetID, networkID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Oracle ODB subnet %s in network %s in %s in project %s: %w", subnetID, networkID, location, projectID, err)
	}

	return attrs, nil
}

// NewOracleDatabaseServiceE creates a Oracle Database service authenticated the same way every other client in this
// module is.
// The ctx parameter supports cancellation and timeouts.
func NewOracleDatabaseServiceE(t testing.TestingT, ctx context.Context) (*oracledatabase.Service, error) {
	return oracledatabase.NewService(ctx, append(withOptions(), option.WithScopes(oracledatabase.CloudPlatformScope))...)
}
