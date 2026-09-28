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

// NewBigLakeServiceE creates a BigLake service authenticated the same way every other client in
// this module is.
// The ctx parameter supports cancellation and timeouts.
func NewBigLakeServiceE(t testing.TestingT, ctx context.Context) (*biglake.Service, error) {
	return biglake.NewService(ctx, append(withOptions(), option.WithScopes(biglake.CloudPlatformScope))...)
}
