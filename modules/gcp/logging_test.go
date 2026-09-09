package gcp_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/gcp/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/logging/v2"
	"google.golang.org/api/option"
)

// newFakeLoggingService points a real *logging.Service at a local test server, so the Google
// transport is exercised rather than a hand-written stand-in for a type we do not own.
func newFakeLoggingService(t *testing.T, handler http.Handler) *logging.Service {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	service, err := logging.NewService(context.Background(),
		option.WithEndpoint(server.URL), option.WithoutAuthentication())
	require.NoError(t, err)

	return service
}

func TestGetLogMetricAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The values are the ones the terraform-google-observability metric module sets, because the
	// point of reading settings back is asserting a module configured the metric it was asked for.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/metrics/gw-library-test"))
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"name":"gw-library-test",
			"description":"created by terratest",
			"filter":"severity >= ERROR",
			"disabled":true
		}`))
	})

	metric, err := gcp.GetLogMetricAttrsWithClient(context.Background(), newFakeLoggingService(t, handler), "gw-library-test-project", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "gw-library-test", metric.Name)
	assert.Equal(t, "created by terratest", metric.Description)
	assert.Equal(t, "severity >= ERROR", metric.Filter)
	assert.True(t, metric.Disabled)
}

func TestGetLogMetricAttrsWithClientMissingMetric(t *testing.T) {
	t.Parallel()

	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	// The error names the metric and the project as well as saying it is absent, so all three are
	// asserted rather than only the phrase.
	_, err := gcp.GetLogMetricAttrsWithClient(context.Background(), newFakeLoggingService(t, handler), "gw-library-test-project", "gone")
	require.ErrorContains(t, err, "does not exist")
	require.ErrorContains(t, err, "gone")
	require.ErrorContains(t, err, "gw-library-test-project")
}
