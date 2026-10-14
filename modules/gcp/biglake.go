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

// GetBigLakeCatalogAttrs returns the settings Google Cloud holds for the BigLake catalog, so a test can assert
// on what was actually created rather than only that it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetBigLakeCatalogAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, catalogID string) *biglake.Catalog {
	result, err := GetBigLakeCatalogAttrsE(t, ctx, projectID, location, catalogID)
	require.NoError(t, err)

	return result
}

// GetBigLakeCatalogAttrsE returns the settings Google Cloud holds for the BigLake catalog.
// The ctx parameter supports cancellation and timeouts.
func GetBigLakeCatalogAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, catalogID string) (*biglake.Catalog, error) {
	logger.Default.Logf(t, "Getting settings for BigLake catalog %s in project %s", catalogID, projectID)

	service, err := NewBigLakeServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetBigLakeCatalogAttrsWithClient(ctx, service, projectID, location, catalogID)
}

// GetBigLakeCatalogAttrsWithClient returns the settings Google Cloud holds for the BigLake catalog using the
// supplied *biglake.Service. Prefer this variant in unit tests where the service is backed by an
// httptest fake server (see biglake_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetBigLakeCatalogAttrsWithClient(ctx context.Context, service *biglake.Service, projectID string, location string, catalogID string) (*biglake.Catalog, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/catalogs/%s", projectID, location, catalogID)

	result, err := service.Projects.Locations.Catalogs.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the BigLake catalog %s in project %s does not exist", catalogID, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for BigLake catalog %s in project %s: %w", catalogID, projectID, err)
	}

	return result, nil
}

// GetBigLakeIcebergCatalogIamPolicyAttrs returns the IAM policy Google Cloud holds for the BigLake Iceberg catalog, so a test can assert on what was
// actually created rather than only that it exists.
// An Iceberg catalog is not the BigLake catalog GetBigLakeCatalogAttrs reads: it lives at
// projects/<project>/catalogs/<catalog> with no location in the name, and the older read would look
// for it under a location and find nothing.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetBigLakeIcebergCatalogIamPolicyAttrs(t testing.TestingT, ctx context.Context, projectID string, catalogID string) *biglake.Policy {
	result, err := GetBigLakeIcebergCatalogIamPolicyAttrsE(t, ctx, projectID, catalogID)
	require.NoError(t, err)

	return result
}

// GetBigLakeIcebergCatalogIamPolicyAttrsE returns the IAM policy Google Cloud holds for the BigLake Iceberg catalog.
// The ctx parameter supports cancellation and timeouts.
func GetBigLakeIcebergCatalogIamPolicyAttrsE(t testing.TestingT, ctx context.Context, projectID string, catalogID string) (*biglake.Policy, error) {
	logger.Default.Logf(t, "Getting the IAM policy for BigLake Iceberg catalog %s in project %s", catalogID, projectID)

	service, err := NewBigLakeServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetBigLakeIcebergCatalogIamPolicyAttrsWithClient(ctx, service, projectID, catalogID)
}

// GetBigLakeIcebergCatalogIamPolicyAttrsWithClient returns the IAM policy Google Cloud holds for the BigLake Iceberg catalog using the supplied
// *biglake.Service. Prefer this variant in unit tests where the service is backed by an httptest
// fake server (see biglake_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetBigLakeIcebergCatalogIamPolicyAttrsWithClient(ctx context.Context, service *biglake.Service, projectID string, catalogID string) (*biglake.Policy, error) {
	name := fmt.Sprintf("projects/%s/catalogs/%s", projectID, catalogID)

	result, err := service.Projects.Catalogs.GetIamPolicy(name).OptionsRequestedPolicyVersion(iamPolicyVersionWithConditions).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the BigLake Iceberg catalog %s in project %s does not exist", catalogID, projectID)
		}

		return nil, fmt.Errorf("failed to get the IAM policy for BigLake Iceberg catalog %s in project %s: %w", catalogID, projectID, err)
	}

	return result, nil
}

// GetBigLakeIcebergTableIamPolicyAttrs returns the IAM policy Google Cloud holds for the BigLake Iceberg table, so a test can assert on what was
// actually created rather than only that it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetBigLakeIcebergTableIamPolicyAttrs(t testing.TestingT, ctx context.Context, projectID string, catalogID string, namespaceID string, tableID string) *biglake.Policy {
	result, err := GetBigLakeIcebergTableIamPolicyAttrsE(t, ctx, projectID, catalogID, namespaceID, tableID)
	require.NoError(t, err)

	return result
}

// GetBigLakeIcebergTableIamPolicyAttrsE returns the IAM policy Google Cloud holds for the BigLake Iceberg table.
// The ctx parameter supports cancellation and timeouts.
func GetBigLakeIcebergTableIamPolicyAttrsE(t testing.TestingT, ctx context.Context, projectID string, catalogID string, namespaceID string, tableID string) (*biglake.Policy, error) {
	logger.Default.Logf(t, "Getting the IAM policy for BigLake Iceberg table %s in namespace %s of catalog %s in project %s", tableID, namespaceID, catalogID, projectID)

	service, err := NewBigLakeServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetBigLakeIcebergTableIamPolicyAttrsWithClient(ctx, service, projectID, catalogID, namespaceID, tableID)
}

// GetBigLakeIcebergTableIamPolicyAttrsWithClient returns the IAM policy Google Cloud holds for the BigLake Iceberg table using the supplied
// *biglake.Service. Prefer this variant in unit tests where the service is backed by an httptest
// fake server (see biglake_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetBigLakeIcebergTableIamPolicyAttrsWithClient(ctx context.Context, service *biglake.Service, projectID string, catalogID string, namespaceID string, tableID string) (*biglake.Policy, error) {
	name := fmt.Sprintf("projects/%s/catalogs/%s/namespaces/%s/tables/%s", projectID, catalogID, namespaceID, tableID)

	result, err := service.Projects.Catalogs.Namespaces.Tables.GetIamPolicy(name).OptionsRequestedPolicyVersion(iamPolicyVersionWithConditions).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the BigLake Iceberg table %s in namespace %s of catalog %s in project %s does not exist", tableID, namespaceID, catalogID, projectID)
		}

		return nil, fmt.Errorf("failed to get the IAM policy for BigLake Iceberg table %s in namespace %s of catalog %s in project %s: %w", tableID, namespaceID, catalogID, projectID, err)
	}

	return result, nil
}

// NewBigLakeServiceE creates a BigLake service authenticated the same way every other client in
// this module is.
// The ctx parameter supports cancellation and timeouts.
func NewBigLakeServiceE(t testing.TestingT, ctx context.Context) (*biglake.Service, error) {
	return biglake.NewService(ctx, append(withOptions(), option.WithScopes(biglake.CloudPlatformScope))...)
}
