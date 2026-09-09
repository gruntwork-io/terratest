package gcp

import (
	"context"
	"errors"
	"fmt"

	"github.com/gruntwork-io/terratest/modules/core/v2/logger"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/bigquery/v2"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
)

// GetBigQueryDatasetAttrs returns the settings Google Cloud holds for the given BigQuery dataset,
// so a test can assert on what was actually created rather than only that it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetBigQueryDatasetAttrs(t testing.TestingT, ctx context.Context, projectID string, datasetID string) *bigquery.Dataset {
	dataset, err := GetBigQueryDatasetAttrsE(t, ctx, projectID, datasetID)
	require.NoError(t, err)

	return dataset
}

// GetBigQueryDatasetAttrsE returns the settings Google Cloud holds for the given BigQuery dataset.
// The ctx parameter supports cancellation and timeouts.
func GetBigQueryDatasetAttrsE(t testing.TestingT, ctx context.Context, projectID string, datasetID string) (*bigquery.Dataset, error) {
	logger.Default.Logf(t, "Getting settings for BigQuery dataset %s in project %s", datasetID, projectID)

	service, err := NewBigQueryServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetBigQueryDatasetAttrsWithClient(ctx, service, projectID, datasetID)
}

// GetBigQueryDatasetAttrsWithClient returns the settings Google Cloud holds for the given BigQuery
// dataset using the supplied *bigquery.Service. Prefer this variant in unit tests where the service
// is backed by an httptest fake server (see bigquery_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetBigQueryDatasetAttrsWithClient(ctx context.Context, service *bigquery.Service, projectID string, datasetID string) (*bigquery.Dataset, error) {
	dataset, err := service.Datasets.Get(projectID, datasetID).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("BigQuery dataset %s does not exist in project %s", datasetID, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for BigQuery dataset %s in project %s: %w", datasetID, projectID, err)
	}

	return dataset, nil
}

// NewBigQueryServiceE creates a BigQuery service authenticated the same way every other client in
// this module is.
// The ctx parameter supports cancellation and timeouts.
func NewBigQueryServiceE(t testing.TestingT, ctx context.Context) (*bigquery.Service, error) {
	return bigquery.NewService(ctx, append(withOptions(), option.WithScopes(bigquery.CloudPlatformScope))...)
}
