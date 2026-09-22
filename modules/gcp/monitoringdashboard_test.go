package gcp_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gruntwork-io/terratest/modules/gcp/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	monitoringdashboard "google.golang.org/api/monitoring/v1"
	"google.golang.org/api/option"
)

// newFakeMonitoringDashboardService points a real Cloud Monitoring dashboards client at a local test server, so the Google transport is
// exercised rather than a hand-written stand-in for a type we do not own.
func newFakeMonitoringDashboardService(t *testing.T, handler http.Handler) *monitoringdashboard.Service {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	service, err := monitoringdashboard.NewService(context.Background(),
		option.WithEndpoint(server.URL), option.WithoutAuthentication())
	require.NoError(t, err)

	return service
}

func TestGetMonitoringDashboardAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a dashboard the terraform-google-
	// observability module created, not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Contains(t, r.URL.Path, "/projects/gw-library-test-project/dashboards/gw-library-test", "unexpected path")
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/123/dashboards/gw-library-test","displayName":"terratest dashboard","gridLayout":{"columns":"2","widgets":[{"title":"requests"}]}}`))
	})

	dashboard, err := gcp.GetMonitoringDashboardAttrsWithClient(context.Background(), newFakeMonitoringDashboardService(t, handler), "gw-library-test-project", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "terratest dashboard", dashboard.DisplayName)
	require.NotNil(t, dashboard.GridLayout)
	assert.Equal(t, int64(2), dashboard.GridLayout.Columns)
	require.Len(t, dashboard.GridLayout.Widgets, 1)
	assert.Equal(t, "requests", dashboard.GridLayout.Widgets[0].Title)
}

func TestGetMonitoringDashboardAttrsWithClientMissingDashboard(t *testing.T) {
	t.Parallel()

	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	// The error names the dashboard and everything that identifies it, and each of those is a value
	// no other part of the message contains, or its check could not fail.
	_, err := gcp.GetMonitoringDashboardAttrsWithClient(context.Background(), newFakeMonitoringDashboardService(t, handler), "gw-library-test-project", "gone")
	require.ErrorContains(t, err, "does not exist")
	require.ErrorContains(t, err, "gone")
	require.ErrorContains(t, err, "gw-library-test-project")
}
