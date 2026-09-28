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
	"google.golang.org/api/option"
	"google.golang.org/api/run/v2"
)

// newFakeCloudRunService points a real Cloud Run client at a local test server, so the Google transport is
// exercised rather than a hand-written stand-in for a type we do not own.
func newFakeCloudRunService(t *testing.T, handler http.Handler) *run.Service {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	service, err := run.NewService(context.Background(),
		option.WithEndpoint(server.URL), option.WithoutAuthentication())
	require.NoError(t, err)

	return service
}

func TestGetCloudRunServiceAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a service the terraform-google-
	// serverless module created, not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Contains(t, r.URL.Path, "/locations/us-central1/services/gw-library-test", "unexpected path")
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/services/gw-library-test","description":"created by terratest","ingress":"INGRESS_TRAFFIC_INTERNAL_ONLY","labels":{"purpose":"terratest"},"template":{"containers":[{"image":"us-docker.pkg.dev/cloudrun/container/hello"}],"maxInstanceRequestConcurrency":8}}`))
	})

	runService, err := gcp.GetCloudRunServiceAttrsWithClient(context.Background(), newFakeCloudRunService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "created by terratest", runService.Description)
	assert.Equal(t, "INGRESS_TRAFFIC_INTERNAL_ONLY", runService.Ingress)
	require.NotNil(t, runService.Template)
	require.Len(t, runService.Template.Containers, 1)
	assert.Equal(t, "us-docker.pkg.dev/cloudrun/container/hello", runService.Template.Containers[0].Image)
}

func TestGetCloudRunServiceAttrsWithClientMissingService(t *testing.T) {
	t.Parallel()

	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	// The error names the service and everything that identifies it, and each of those is a value no
	// other part of the message contains, or its check could not fail.
	_, err := gcp.GetCloudRunServiceAttrsWithClient(context.Background(), newFakeCloudRunService(t, handler), "gw-library-test-project", "us-central1", "gone")
	require.ErrorContains(t, err, "does not exist")
	require.ErrorContains(t, err, "gone")
	require.ErrorContains(t, err, "us-central1")
	require.ErrorContains(t, err, "gw-library-test-project")
}

func TestGetCloudRunJobAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a job the terraform-google-serverless
	// module created, not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Contains(t, r.URL.Path, "/locations/us-central1/jobs/gw-library-test", "unexpected path")
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/jobs/gw-library-test","labels":{"purpose":"terratest"},"template":{"taskCount":3,"template":{"containers":[{"image":"us-docker.pkg.dev/cloudrun/container/job"}],"maxRetries":2}}}`))
	})

	job, err := gcp.GetCloudRunJobAttrsWithClient(context.Background(), newFakeCloudRunService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	require.NotNil(t, job.Template)
	assert.Equal(t, int64(3), job.Template.TaskCount)
	require.NotNil(t, job.Template.Template)
	require.Len(t, job.Template.Template.Containers, 1)
	assert.Equal(t, "us-docker.pkg.dev/cloudrun/container/job", job.Template.Template.Containers[0].Image)
}

func TestGetCloudRunJobAttrsWithClientMissingJob(t *testing.T) {
	t.Parallel()

	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	// The error names the job and everything that identifies it, and each of those is a value no
	// other part of the message contains, or its check could not fail.
	_, err := gcp.GetCloudRunJobAttrsWithClient(context.Background(), newFakeCloudRunService(t, handler), "gw-library-test-project", "us-central1", "gone")
	require.ErrorContains(t, err, "does not exist")
	require.ErrorContains(t, err, "gone")
	require.ErrorContains(t, err, "us-central1")
	require.ErrorContains(t, err, "gw-library-test-project")
}

func TestGetCloudRunWorkerPoolAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a pool the terraform-google-serverless worker pool module created, not a
	// copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/workerPools/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/workerPools/gw-library-test","labels":{"purpose":"terratest"},"scaling":{"manualInstanceCount":0},"template":{"containers":[{"image":"us-docker.pkg.dev/cloudrun/container/worker-pool"}]}}`))
	})

	pool, err := gcp.GetCloudRunWorkerPoolAttrsWithClient(context.Background(), newFakeCloudRunService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "terratest", pool.Labels["purpose"])
	require.NotNil(t, pool.Template)
	require.Len(t, pool.Template.Containers, 1)
}

func TestGetCloudRunServiceIamPolicyAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a service the terraform-google-serverless service IAM policy module granted access
	// on, not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/services/gw-library-test:getIamPolicy"), "unexpected path %s", r.URL.Path)
		// A conditional binding only comes back at version 3, so the read has to ask for it.
		assert.Equal(t, "3", r.URL.Query().Get("options.requestedPolicyVersion"))
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"version":3,"etag":"BwXhqw==","bindings":[{"role":"roles/run.invoker","members":["serviceAccount:gw-library-test@gw-library-test-project.iam.gserviceaccount.com"],"condition":{"title":"until 2030","expression":"request.time < timestamp(\"2030-01-01T00:00:00Z\")"}}]}`))
	})

	policy, err := gcp.GetCloudRunServiceIamPolicyAttrsWithClient(context.Background(), newFakeCloudRunService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	require.Len(t, policy.Bindings, 1)
	assert.Equal(t, int64(3), policy.Version)
	assert.Equal(t, "roles/run.invoker", policy.Bindings[0].Role)
	require.NotNil(t, policy.Bindings[0].Condition, "a conditional binding should keep its condition")
	assert.Equal(t, `request.time < timestamp("2030-01-01T00:00:00Z")`, policy.Bindings[0].Condition.Expression)
}

func TestGetCloudRunJobIamPolicyAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a job the terraform-google-serverless job IAM policy module granted access
	// on, not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/jobs/gw-library-test:getIamPolicy"), "unexpected path %s", r.URL.Path)
		// A conditional binding only comes back at version 3, so the read has to ask for it.
		assert.Equal(t, "3", r.URL.Query().Get("options.requestedPolicyVersion"))
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"version":3,"etag":"BwXhqw==","bindings":[{"role":"roles/run.invoker","members":["serviceAccount:gw-library-test@gw-library-test-project.iam.gserviceaccount.com"],"condition":{"title":"until 2030","expression":"request.time < timestamp(\"2030-01-01T00:00:00Z\")"}}]}`))
	})

	policy, err := gcp.GetCloudRunJobIamPolicyAttrsWithClient(context.Background(), newFakeCloudRunService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	require.Len(t, policy.Bindings, 1)
	assert.Equal(t, int64(3), policy.Version)
	assert.Equal(t, "roles/run.invoker", policy.Bindings[0].Role)
	require.NotNil(t, policy.Bindings[0].Condition, "a conditional binding should keep its condition")
	assert.Equal(t, `request.time < timestamp("2030-01-01T00:00:00Z")`, policy.Bindings[0].Condition.Expression)
}

func TestGetCloudRunWorkerPoolIamPolicyAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a worker pool the terraform-google-serverless worker pool IAM policy module granted access
	// on, not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/workerPools/gw-library-test:getIamPolicy"), "unexpected path %s", r.URL.Path)
		// A conditional binding only comes back at version 3, so the read has to ask for it.
		assert.Equal(t, "3", r.URL.Query().Get("options.requestedPolicyVersion"))
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"version":3,"etag":"BwXhqw==","bindings":[{"role":"roles/run.viewer","members":["serviceAccount:gw-library-test@gw-library-test-project.iam.gserviceaccount.com"],"condition":{"title":"until 2030","expression":"request.time < timestamp(\"2030-01-01T00:00:00Z\")"}}]}`))
	})

	policy, err := gcp.GetCloudRunWorkerPoolIamPolicyAttrsWithClient(context.Background(), newFakeCloudRunService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	require.Len(t, policy.Bindings, 1)
	assert.Equal(t, int64(3), policy.Version)
	assert.Equal(t, "roles/run.viewer", policy.Bindings[0].Role)
	require.NotNil(t, policy.Bindings[0].Condition, "a conditional binding should keep its condition")
	assert.Equal(t, `request.time < timestamp("2030-01-01T00:00:00Z")`, policy.Bindings[0].Condition.Expression)
}
