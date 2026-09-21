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

func TestGetLogBucketAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a log bucket the
	// terraform-google-observability module created, not a copy of any one fixture's values. A bucket
	// is named by its location as well as its id, so both have to reach the request.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/global/buckets/gw-library-test"))
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"name":"projects/gw-library-test-project/locations/global/buckets/gw-library-test",
			"description":"created by terratest",
			"retentionDays":7,
			"locked":false,
			"lifecycleState":"ACTIVE"
		}`))
	})

	bucket, err := gcp.GetLogBucketAttrsWithClient(context.Background(), newFakeLoggingService(t, handler), "gw-library-test-project", "global", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "projects/gw-library-test-project/locations/global/buckets/gw-library-test", bucket.Name)
	assert.Equal(t, "created by terratest", bucket.Description)
	assert.Equal(t, int64(7), bucket.RetentionDays)
	assert.False(t, bucket.Locked)
	assert.Equal(t, "ACTIVE", bucket.LifecycleState)
}

func TestGetLogBucketAttrsWithClientMissingBucket(t *testing.T) {
	t.Parallel()

	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	// The error names the bucket, its location and the project as well as saying it is absent,
	// since a bucket is only identified by all three together.
	_, err := gcp.GetLogBucketAttrsWithClient(context.Background(), newFakeLoggingService(t, handler), "gw-library-test-project", "us-central1", "gone")
	require.ErrorContains(t, err, "does not exist")
	require.ErrorContains(t, err, "gone")
	require.ErrorContains(t, err, "us-central1")
	require.ErrorContains(t, err, "gw-library-test-project")
}
