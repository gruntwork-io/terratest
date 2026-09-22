package gcp

import (
	"context"
	"errors"
	"fmt"

	"github.com/gruntwork-io/terratest/modules/core/v2/logger"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/dataplex/v1"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
)

// GetDataplexLakeAttrs returns the settings Google Cloud holds for the given Dataplex lake, so a
// test can assert on what was actually created rather than only that it exists. A lake lives in a
// location, which has to be given.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetDataplexLakeAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, lakeID string) *dataplex.GoogleCloudDataplexV1Lake {
	lake, err := GetDataplexLakeAttrsE(t, ctx, projectID, location, lakeID)
	require.NoError(t, err)

	return lake
}

// GetDataplexLakeAttrsE returns the settings Google Cloud holds for the given Dataplex lake.
// The ctx parameter supports cancellation and timeouts.
func GetDataplexLakeAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, lakeID string) (*dataplex.GoogleCloudDataplexV1Lake, error) {
	logger.Default.Logf(t, "Getting settings for Dataplex lake %s in project %s", lakeID, projectID)

	service, err := NewDataplexServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetDataplexLakeAttrsWithClient(ctx, service, projectID, location, lakeID)
}

// GetDataplexLakeAttrsWithClient returns the settings Google Cloud holds for the given Dataplex
// lake using the supplied *dataplex.Service. Prefer this variant in unit tests where the service is
// backed by an httptest fake server (see dataplex_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetDataplexLakeAttrsWithClient(ctx context.Context, service *dataplex.Service, projectID string, location string, lakeID string) (*dataplex.GoogleCloudDataplexV1Lake, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/lakes/%s", projectID, location, lakeID)

	lake, err := service.Projects.Locations.Lakes.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("Dataplex lake %s does not exist in project %s", lakeID, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Dataplex lake %s in project %s: %w", lakeID, projectID, err)
	}

	return lake, nil
}

// NewDataplexServiceE creates a Dataplex service authenticated the same way every other client in
// this module is.
// The ctx parameter supports cancellation and timeouts.
func NewDataplexServiceE(t testing.TestingT, ctx context.Context) (*dataplex.Service, error) {
	return dataplex.NewService(ctx, append(withOptions(), option.WithScopes(dataplex.CloudPlatformScope))...)
}
