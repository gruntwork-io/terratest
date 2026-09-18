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

// NewLoggingServiceE creates a Cloud Logging service authenticated the same way every other client
// in this module is.
// The ctx parameter supports cancellation and timeouts.
func NewLoggingServiceE(t testing.TestingT, ctx context.Context) (*logging.Service, error) {
	return logging.NewService(ctx, append(withOptions(), option.WithScopes(logging.CloudPlatformScope))...)
}
