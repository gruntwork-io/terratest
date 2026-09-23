package gcp

import (
	"context"
	"errors"
	"fmt"

	"github.com/gruntwork-io/terratest/modules/core/v2/logger"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/googleapi"
	monitoringdashboard "google.golang.org/api/monitoring/v1"
	"google.golang.org/api/option"
)

// GetMonitoringDashboardAttrs returns the settings Google Cloud holds for the given Cloud
// Monitoring dashboard, so a test can assert on what was actually created rather than only that it
// exists. The dashboards live behind their own version of the Monitoring API, so this read builds a
// client of its own.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetMonitoringDashboardAttrs(t testing.TestingT, ctx context.Context, projectID string, dashboardID string) *monitoringdashboard.Dashboard {
	dashboard, err := GetMonitoringDashboardAttrsE(t, ctx, projectID, dashboardID)
	require.NoError(t, err)

	return dashboard
}

// GetMonitoringDashboardAttrsE returns the settings Google Cloud holds for the given Cloud
// Monitoring dashboard.
// The ctx parameter supports cancellation and timeouts.
func GetMonitoringDashboardAttrsE(t testing.TestingT, ctx context.Context, projectID string, dashboardID string) (*monitoringdashboard.Dashboard, error) {
	logger.Default.Logf(t, "Getting settings for Cloud Monitoring dashboard %s in project %s", dashboardID, projectID)

	service, err := NewMonitoringDashboardServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetMonitoringDashboardAttrsWithClient(ctx, service, projectID, dashboardID)
}

// GetMonitoringDashboardAttrsWithClient returns the settings Google Cloud holds for the given Cloud
// Monitoring dashboard using the supplied *monitoringdashboard.Service. Prefer this variant in unit
// tests where the service is backed by an httptest fake server (see monitoringdashboard_test.go for
// the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetMonitoringDashboardAttrsWithClient(ctx context.Context, service *monitoringdashboard.Service, projectID string, dashboardID string) (*monitoringdashboard.Dashboard, error) {
	name := fmt.Sprintf("projects/%s/dashboards/%s", projectID, dashboardID)

	dashboard, err := service.Projects.Dashboards.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Cloud Monitoring dashboard %s does not exist in project %s", dashboardID, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Cloud Monitoring dashboard %s in project %s: %w", dashboardID, projectID, err)
	}

	return dashboard, nil
}

// NewMonitoringDashboardServiceE creates a Cloud Monitoring dashboards service authenticated the
// same way every other client in this module is.
// The ctx parameter supports cancellation and timeouts.
func NewMonitoringDashboardServiceE(t testing.TestingT, ctx context.Context) (*monitoringdashboard.Service, error) {
	return monitoringdashboard.NewService(ctx, append(withOptions(), option.WithScopes(monitoringdashboard.CloudPlatformScope))...)
}
