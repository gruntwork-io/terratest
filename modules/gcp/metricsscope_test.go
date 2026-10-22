package gcp_test

import (
	"context"
	"net/http"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/gcp/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetMonitoringMetricsScopeAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns once the monitored project module has added a
	// project to a scope, not a copy of any one fixture's values. The entries come back named by
	// project number rather than by project id.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/locations/global/metricsScopes/gw-library-test-project"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"locations/global/metricsScopes/424242424242","monitoredProjects":[{"name":"locations/global/metricsScopes/424242424242/projects/424242424242"},{"name":"locations/global/metricsScopes/424242424242/projects/515151515151"}]}`))
	})

	scope, err := gcp.GetMonitoringMetricsScopeAttrsWithClient(context.Background(), newFakeMonitoringDashboardService(t, handler), "gw-library-test-project")
	require.NoError(t, err)

	require.Len(t, scope.MonitoredProjects, 2)
	assert.Equal(t, "locations/global/metricsScopes/424242424242/projects/515151515151", scope.MonitoredProjects[1].Name)
}

func TestGetMonitoringMetricsScopeAttrsWithClientReportsAMissingScope(t *testing.T) {
	t.Parallel()

	// A project Monitoring has never been enabled on answers 404, and a caller should be told that in
	// those words rather than be handed the transport's own.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error":{"code":404,"message":"Not found."}}`))
	})

	_, err := gcp.GetMonitoringMetricsScopeAttrsWithClient(context.Background(), newFakeMonitoringDashboardService(t, handler), "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "has no metrics scope")
}
