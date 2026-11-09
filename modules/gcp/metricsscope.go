package gcp

import (
	"context"
	"errors"
	"fmt"

	"github.com/gruntwork-io/terratest/modules/core/v2/logger"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/googleapi"
	monitoringv1 "google.golang.org/api/monitoring/v1"
	"google.golang.org/api/option"
)

// GetMonitoringMetricsScopeAttrs returns the settings Google Cloud holds for the metrics scope of the given scoping project, so a test can assert on
// what was actually created rather than only that it exists.
// A monitored project is an entry in the scoping project's metrics scope rather than a resource with a
// read of its own, so a test asserts on the list this returns.
// The scoping project may be given as its id or as its number, and the entries come back named by
// number.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetMonitoringMetricsScopeAttrs(t testing.TestingT, ctx context.Context, scopingProject string) *monitoringv1.MetricsScope {
	result, err := GetMonitoringMetricsScopeAttrsE(t, ctx, scopingProject)
	require.NoError(t, err)

	return result
}

// GetMonitoringMetricsScopeAttrsE returns the settings Google Cloud holds for the metrics scope of the given scoping project.
// The ctx parameter supports cancellation and timeouts.
func GetMonitoringMetricsScopeAttrsE(t testing.TestingT, ctx context.Context, scopingProject string) (*monitoringv1.MetricsScope, error) {
	logger.Default.Logf(t, "Getting the metrics scope of project %s", scopingProject)

	service, err := NewMetricsScopeServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetMonitoringMetricsScopeAttrsWithClient(ctx, service, scopingProject)
}

// GetMonitoringMetricsScopeAttrsWithClient returns the settings Google Cloud holds for the metrics scope of the given scoping project using the supplied
// *monitoringv1.Service. Prefer this variant in unit tests where the service is backed by an httptest fake server
// (see metricsscope_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetMonitoringMetricsScopeAttrsWithClient(ctx context.Context, service *monitoringv1.Service, scopingProject string) (*monitoringv1.MetricsScope, error) {
	name := fmt.Sprintf("locations/global/metricsScopes/%s", scopingProject)

	result, err := service.Locations.Global.MetricsScopes.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("project %s has no metrics scope", scopingProject)
		}

		return nil, fmt.Errorf("failed to get the metrics scope of project %s: %w", scopingProject, err)
	}

	return result, nil
}

// NewMetricsScopeServiceE creates a metrics scope service authenticated the same way every other client in
// this module is.
// Metrics scopes are served by the first version of the Monitoring API, the same one the dashboards
// in monitoringdashboard.go come from.
// The ctx parameter supports cancellation and timeouts.
func NewMetricsScopeServiceE(t testing.TestingT, ctx context.Context) (*monitoringv1.Service, error) {
	return monitoringv1.NewService(ctx, append(withOptions(), option.WithScopes(monitoringv1.CloudPlatformScope))...)
}
