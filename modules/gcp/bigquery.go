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

// GetBigQueryTableAttrs returns the settings Google Cloud holds for the given BigQuery table, so a
// test can assert on what was actually created rather than only that it exists. The schema comes
// back with the table, so a caller can check each column.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetBigQueryTableAttrs(t testing.TestingT, ctx context.Context, projectID string, datasetID string, tableID string) *bigquery.Table {
	table, err := GetBigQueryTableAttrsE(t, ctx, projectID, datasetID, tableID)
	require.NoError(t, err)

	return table
}

// GetBigQueryTableAttrsE returns the settings Google Cloud holds for the given BigQuery table.
// The ctx parameter supports cancellation and timeouts.
func GetBigQueryTableAttrsE(t testing.TestingT, ctx context.Context, projectID string, datasetID string, tableID string) (*bigquery.Table, error) {
	logger.Default.Logf(t, "Getting settings for BigQuery table %s.%s in project %s", datasetID, tableID, projectID)

	service, err := NewBigQueryServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetBigQueryTableAttrsWithClient(ctx, service, projectID, datasetID, tableID)
}

// GetBigQueryTableAttrsWithClient returns the settings Google Cloud holds for the given BigQuery
// table using the supplied *bigquery.Service. Prefer this variant in unit tests where the service
// is backed by an httptest fake server (see bigquery_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetBigQueryTableAttrsWithClient(ctx context.Context, service *bigquery.Service, projectID string, datasetID string, tableID string) (*bigquery.Table, error) {
	table, err := service.Tables.Get(projectID, datasetID, tableID).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("BigQuery table %s.%s does not exist in project %s", datasetID, tableID, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for BigQuery table %s.%s in project %s: %w", datasetID, tableID, projectID, err)
	}

	return table, nil
}

// NewBigQueryServiceE creates a BigQuery service authenticated the same way every other client in
// this module is.
// The ctx parameter supports cancellation and timeouts.
func NewBigQueryServiceE(t testing.TestingT, ctx context.Context) (*bigquery.Service, error) {
	return bigquery.NewService(ctx, append(withOptions(), option.WithScopes(bigquery.CloudPlatformScope))...)
}

// GetBigQueryJobAttrs returns the settings Google Cloud holds for the given BigQuery job, so a test
// can assert on what it was asked to run rather than only that it exists. A job is a record of work
// rather than a resource that persists, so this reads it after it has finished. A job that ran in a
// single region is only found when its location is named, so the caller passes the one it ran in; the
// US and EU multi-regions accept an empty location.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetBigQueryJobAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, jobID string) *bigquery.Job {
	job, err := GetBigQueryJobAttrsE(t, ctx, projectID, location, jobID)
	require.NoError(t, err)

	return job
}

// GetBigQueryJobAttrsE returns the settings Google Cloud holds for the given BigQuery job.
// The ctx parameter supports cancellation and timeouts.
func GetBigQueryJobAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, jobID string) (*bigquery.Job, error) {
	logger.Default.Logf(t, "Getting settings for BigQuery job %s in %s in project %s", jobID, location, projectID)

	service, err := NewBigQueryServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetBigQueryJobAttrsWithClient(ctx, service, projectID, location, jobID)
}

// GetBigQueryJobAttrsWithClient returns the settings Google Cloud holds for the given BigQuery job
// using the supplied *bigquery.Service. Prefer this variant in unit tests where the service is
// backed by an httptest fake server (see bigquery_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetBigQueryJobAttrsWithClient(ctx context.Context, service *bigquery.Service, projectID string, location string, jobID string) (*bigquery.Job, error) {
	// This call names the project and the job as separate parameters rather than as one path. A job
	// that ran in a single region is only found when the location goes with it.
	call := service.Jobs.Get(projectID, jobID)
	if location != "" {
		call = call.Location(location)
	}

	job, err := call.Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the BigQuery job %s in %s in project %s does not exist", jobID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for BigQuery job %s in %s in project %s: %w", jobID, location, projectID, err)
	}

	return job, nil
}
