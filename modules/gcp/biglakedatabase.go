package gcp

import (
	"context"
	"errors"
	"fmt"

	"github.com/gruntwork-io/terratest/modules/core/v2/logger"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/biglake/v1"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
)

// GetBigLakeDatabaseAttrs returns the settings Google Cloud holds for the given BigLake database, so a test can assert on
// what was actually created rather than only that it exists. A database belongs to a catalog, so the caller names both.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetBigLakeDatabaseAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, catalogID string, databaseID string) *biglake.Database {
	result, err := GetBigLakeDatabaseAttrsE(t, ctx, projectID, location, catalogID, databaseID)
	require.NoError(t, err)

	return result
}

// GetBigLakeDatabaseAttrsE returns the settings Google Cloud holds for the given BigLake database.
// The ctx parameter supports cancellation and timeouts.
func GetBigLakeDatabaseAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, catalogID string, databaseID string) (*biglake.Database, error) {
	logger.Default.Logf(t, "Getting settings for BigLake database %s in project %s", databaseID, projectID)

	service, err := NewBigLakeDatabaseServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetBigLakeDatabaseAttrsWithClient(ctx, service, projectID, location, catalogID, databaseID)
}

// GetBigLakeDatabaseAttrsWithClient returns the settings Google Cloud holds for the given BigLake database using the
// supplied *biglake.Service. Prefer this variant in unit tests where the service is backed by an
// httptest fake server (see biglakedatabase_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetBigLakeDatabaseAttrsWithClient(ctx context.Context, service *biglake.Service, projectID string, location string, catalogID string, databaseID string) (*biglake.Database, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/catalogs/%s/databases/%s", projectID, location, catalogID, databaseID)

	result, err := service.Projects.Locations.Catalogs.Databases.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the BigLake database %s in project %s does not exist", databaseID, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for BigLake database %s in project %s: %w", databaseID, projectID, err)
	}

	return result, nil
}

// GetBigLakeTableAttrs returns the settings Google Cloud holds for the given BigLake table, so a test can assert on
// what was actually created rather than only that it exists. A table belongs to a database inside a catalog, so the caller names all three.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetBigLakeTableAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, catalogID string, databaseID string, tableID string) *biglake.Table {
	result, err := GetBigLakeTableAttrsE(t, ctx, projectID, location, catalogID, databaseID, tableID)
	require.NoError(t, err)

	return result
}

// GetBigLakeTableAttrsE returns the settings Google Cloud holds for the given BigLake table.
// The ctx parameter supports cancellation and timeouts.
func GetBigLakeTableAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, catalogID string, databaseID string, tableID string) (*biglake.Table, error) {
	logger.Default.Logf(t, "Getting settings for BigLake table %s in project %s", tableID, projectID)

	service, err := NewBigLakeDatabaseServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetBigLakeTableAttrsWithClient(ctx, service, projectID, location, catalogID, databaseID, tableID)
}

// GetBigLakeTableAttrsWithClient returns the settings Google Cloud holds for the given BigLake table using the
// supplied *biglake.Service. Prefer this variant in unit tests where the service is backed by an
// httptest fake server (see biglakedatabase_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetBigLakeTableAttrsWithClient(ctx context.Context, service *biglake.Service, projectID string, location string, catalogID string, databaseID string, tableID string) (*biglake.Table, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/catalogs/%s/databases/%s/tables/%s", projectID, location, catalogID, databaseID, tableID)

	result, err := service.Projects.Locations.Catalogs.Databases.Tables.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the BigLake table %s in project %s does not exist", tableID, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for BigLake table %s in project %s: %w", tableID, projectID, err)
	}

	return result, nil
}

// NewBigLakeDatabaseServiceE creates a BigLake service authenticated the same way every other client in this
// module is.
// The ctx parameter supports cancellation and timeouts.
func NewBigLakeDatabaseServiceE(t testing.TestingT, ctx context.Context) (*biglake.Service, error) {
	return biglake.NewService(ctx, append(withOptions(), option.WithScopes(biglake.CloudPlatformScope))...)
}
