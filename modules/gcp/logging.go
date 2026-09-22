package gcp

import (
	"context"
	"errors"
	"fmt"

	"github.com/gruntwork-io/terratest/modules/core/v2/logger"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/logging/v2"
	"google.golang.org/api/option"
)

// GetLogMetricAttrs returns the settings Google Cloud holds for the given log-based metric, so a
// test can assert on what was actually created rather than only that it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetLogMetricAttrs(t testing.TestingT, ctx context.Context, projectID string, metricName string) *logging.LogMetric {
	metric, err := GetLogMetricAttrsE(t, ctx, projectID, metricName)
	require.NoError(t, err)

	return metric
}

// GetLogMetricAttrsE returns the settings Google Cloud holds for the given log-based metric.
// The ctx parameter supports cancellation and timeouts.
func GetLogMetricAttrsE(t testing.TestingT, ctx context.Context, projectID string, metricName string) (*logging.LogMetric, error) {
	logger.Default.Logf(t, "Getting settings for log-based metric %s in project %s", metricName, projectID)

	service, err := NewLoggingServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetLogMetricAttrsWithClient(ctx, service, projectID, metricName)
}

// GetLogMetricAttrsWithClient returns the settings Google Cloud holds for the given log-based
// metric using the supplied *logging.Service. Prefer this variant in unit tests where the service
// is backed by an httptest fake server (see logging_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetLogMetricAttrsWithClient(ctx context.Context, service *logging.Service, projectID string, metricName string) (*logging.LogMetric, error) {
	name := fmt.Sprintf("projects/%s/metrics/%s", projectID, metricName)

	metric, err := service.Projects.Metrics.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("log-based metric %s does not exist in project %s", metricName, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for log-based metric %s in project %s: %w", metricName, projectID, err)
	}

	return metric, nil
}

// GetLogSinkAttrs returns the settings Google Cloud holds for the given project log sink, so a
// test can assert on what was actually created rather than only that it exists. The writer
// identity comes back with it, which is the service account a destination has to grant.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetLogSinkAttrs(t testing.TestingT, ctx context.Context, projectID string, sinkName string) *logging.LogSink {
	sink, err := GetLogSinkAttrsE(t, ctx, projectID, sinkName)
	require.NoError(t, err)

	return sink
}

// GetLogSinkAttrsE returns the settings Google Cloud holds for the given project log sink.
// The ctx parameter supports cancellation and timeouts.
func GetLogSinkAttrsE(t testing.TestingT, ctx context.Context, projectID string, sinkName string) (*logging.LogSink, error) {
	logger.Default.Logf(t, "Getting settings for log sink %s in project %s", sinkName, projectID)

	service, err := NewLoggingServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetLogSinkAttrsWithClient(ctx, service, projectID, sinkName)
}

// GetLogSinkAttrsWithClient returns the settings Google Cloud holds for the given project log sink
// using the supplied *logging.Service. Prefer this variant in unit tests where the service is
// backed by an httptest fake server (see logging_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetLogSinkAttrsWithClient(ctx context.Context, service *logging.Service, projectID string, sinkName string) (*logging.LogSink, error) {
	name := fmt.Sprintf("projects/%s/sinks/%s", projectID, sinkName)

	sink, err := service.Projects.Sinks.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("log sink %s does not exist in project %s", sinkName, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for log sink %s in project %s: %w", sinkName, projectID, err)
	}

	return sink, nil
}

// GetLogBucketAttrs returns the settings Google Cloud holds for the given log bucket, so a test can
// assert on what was actually created rather than only that it exists. A bucket is named by its
// location as well as its id, and a deleted one is still returned for seven days with a
// lifecycleState of DELETE_REQUESTED rather than as absent.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetLogBucketAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, bucketID string) *logging.LogBucket {
	bucket, err := GetLogBucketAttrsE(t, ctx, projectID, location, bucketID)
	require.NoError(t, err)

	return bucket
}

// GetLogBucketAttrsE returns the settings Google Cloud holds for the given log bucket.
// The ctx parameter supports cancellation and timeouts.
func GetLogBucketAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, bucketID string) (*logging.LogBucket, error) {
	logger.Default.Logf(t, "Getting settings for log bucket %s in location %s in project %s", bucketID, location, projectID)

	service, err := NewLoggingServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetLogBucketAttrsWithClient(ctx, service, projectID, location, bucketID)
}

// GetLogBucketAttrsWithClient returns the settings Google Cloud holds for the given log bucket using
// the supplied *logging.Service. Prefer this variant in unit tests where the service is backed by an
// httptest fake server (see logging_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetLogBucketAttrsWithClient(ctx context.Context, service *logging.Service, projectID string, location string, bucketID string) (*logging.LogBucket, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/buckets/%s", projectID, location, bucketID)

	bucket, err := service.Projects.Locations.Buckets.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("log bucket %s does not exist in location %s in project %s", bucketID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for log bucket %s in location %s in project %s: %w", bucketID, location, projectID, err)
	}

	return bucket, nil
}

// NewLoggingServiceE creates a Cloud Logging service authenticated the same way every other client
// in this module is.
// The ctx parameter supports cancellation and timeouts.
func NewLoggingServiceE(t testing.TestingT, ctx context.Context) (*logging.Service, error) {
	return logging.NewService(ctx, append(withOptions(), option.WithScopes(logging.CloudPlatformScope))...)
}
