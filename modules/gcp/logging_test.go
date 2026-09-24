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

func TestGetLogSinkAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The values are the ones the terraform-google-observability project sink module sets, because
	// the point of reading settings back is asserting a module configured the sink it was asked for.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/sinks/gw-library-test"))
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"name":"gw-library-test",
			"description":"created by terratest",
			"destination":"logging.googleapis.com/projects/gw-library-test-project/locations/global/buckets/_Default",
			"filter":"severity>=ERROR",
			"disabled":true,
			"writerIdentity":"serviceAccount:service-123@gcp-sa-logging.iam.gserviceaccount.com"
		}`))
	})

	sink, err := gcp.GetLogSinkAttrsWithClient(context.Background(), newFakeLoggingService(t, handler), "gw-library-test-project", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "created by terratest", sink.Description)
	assert.Equal(t, "logging.googleapis.com/projects/gw-library-test-project/locations/global/buckets/_Default", sink.Destination)
	assert.Equal(t, "severity>=ERROR", sink.Filter)
	assert.True(t, sink.Disabled)
	assert.Contains(t, sink.WriterIdentity, "serviceAccount:")
}

func TestGetLogSinkAttrsWithClientMissingSink(t *testing.T) {
	t.Parallel()

	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	// The error names the sink and the project as well as saying it is absent, so all three are
	// asserted rather than only the phrase.
	_, err := gcp.GetLogSinkAttrsWithClient(context.Background(), newFakeLoggingService(t, handler), "gw-library-test-project", "gone")
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

func TestGetLogExclusionAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for an exclusion the terraform-google-observability project exclusion module created, not a
	// copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/v2/projects/gw-library-test-project/exclusions/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"gw-library-test","description":"created by terratest","filter":"severity < WARNING","disabled":true}`))
	})

	exclusion, err := gcp.GetLogExclusionAttrsWithClient(context.Background(), newFakeLoggingService(t, handler), "gw-library-test-project", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "created by terratest", exclusion.Description)
	assert.Equal(t, "severity < WARNING", exclusion.Filter)
	assert.True(t, exclusion.Disabled)
}

func TestGetLogViewAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a view the terraform-google-observability log view module created, not a
	// copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/v2/projects/gw-library-test-project/locations/global/buckets/gw-library-test/views/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/global/buckets/gw-library-test/views/gw-library-test","description":"created by terratest","filter":"LOG_ID(\"terratest\")"}`))
	})

	view, err := gcp.GetLogViewAttrsWithClient(context.Background(), newFakeLoggingService(t, handler), "gw-library-test-project", "global", "gw-library-test", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "created by terratest", view.Description)
	assert.Equal(t, `LOG_ID("terratest")`, view.Filter)
}

func TestGetLogLinkAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a link the terraform-google-observability linked dataset module created, not a
	// copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/v2/projects/gw-library-test-project/locations/global/buckets/gw-library-test/links/gw_library_test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/global/buckets/gw-library-test/links/gw_library_test","description":"created by terratest","lifecycleState":"ACTIVE","bigqueryDataset":{"datasetId":"bigquery.googleapis.com/projects/gw-library-test-project/datasets/gw_library_test"}}`))
	})

	link, err := gcp.GetLogLinkAttrsWithClient(context.Background(), newFakeLoggingService(t, handler), "gw-library-test-project", "global", "gw-library-test", "gw_library_test")
	require.NoError(t, err)

	assert.Equal(t, "created by terratest", link.Description)
	assert.Equal(t, "ACTIVE", link.LifecycleState)
	require.NotNil(t, link.BigqueryDataset)
	assert.Equal(t, "bigquery.googleapis.com/projects/gw-library-test-project/datasets/gw_library_test", link.BigqueryDataset.DatasetId)
}

func TestGetLogScopeAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a scope the terraform-google-observability log scope module created, not a
	// copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/v2/projects/gw-library-test-project/locations/global/logScopes/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/global/logScopes/gw-library-test","description":"created by terratest","resourceNames":["projects/gw-library-test-project"]}`))
	})

	scope, err := gcp.GetLogScopeAttrsWithClient(context.Background(), newFakeLoggingService(t, handler), "gw-library-test-project", "global", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "created by terratest", scope.Description)
	assert.Equal(t, []string{"projects/gw-library-test-project"}, scope.ResourceNames)
}

func TestGetLogSavedQueryAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a saved query the terraform-google-observability saved query module created, not a
	// copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/v2/projects/gw-library-test-project/locations/global/savedQueries/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/global/savedQueries/gw-library-test","displayName":"terratest query","description":"created by terratest","visibility":"PRIVATE","loggingQuery":{"filter":"severity >= ERROR","summaryFieldStart":0,"summaryFieldEnd":20}}`))
	})

	query, err := gcp.GetLogSavedQueryAttrsWithClient(context.Background(), newFakeLoggingService(t, handler), "gw-library-test-project", "global", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "terratest query", query.DisplayName)
	assert.Equal(t, "PRIVATE", query.Visibility)
	require.NotNil(t, query.LoggingQuery)
	assert.Equal(t, "severity >= ERROR", query.LoggingQuery.Filter)
	assert.Equal(t, int64(20), query.LoggingQuery.SummaryFieldEnd)
}

func TestGetLogViewIamPolicyAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a view the terraform-google-observability log view IAM policy module created, not a
	// copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/v2/projects/gw-library-test-project/locations/global/buckets/gw-library-test/views/gw-library-test:getIamPolicy"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"version":1,"etag":"BwXhqw==","bindings":[{"role":"roles/logging.viewAccessor","members":["serviceAccount:gw-library-test@gw-library-test-project.iam.gserviceaccount.com"]}]}`))
	})

	policy, err := gcp.GetLogViewIamPolicyAttrsWithClient(context.Background(), newFakeLoggingService(t, handler), "gw-library-test-project", "global", "gw-library-test", "gw-library-test")
	require.NoError(t, err)

	require.Len(t, policy.Bindings, 1)
	assert.Equal(t, "roles/logging.viewAccessor", policy.Bindings[0].Role)
}
