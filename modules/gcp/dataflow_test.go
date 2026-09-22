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
	dataflow "google.golang.org/api/dataflow/v1b3"
	"google.golang.org/api/option"
)

// newFakeDataflowService points a real Dataflow client at a local test server, so the Google transport is
// exercised rather than a hand-written stand-in for a type we do not own.
func newFakeDataflowService(t *testing.T, handler http.Handler) *dataflow.Service {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	service, err := dataflow.NewService(context.Background(),
		option.WithEndpoint(server.URL), option.WithoutAuthentication())
	require.NoError(t, err)

	return service
}

func TestGetDataflowJobAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a job the terraform-google-data-
	// analytics module created, not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/jobs/2026-09-22_00_00_00-1"), "unexpected path %s", r.URL.Path)
		// The view rides in the query string, so asserting on the path alone would pass with the
		// summary view, which carries neither labels nor an environment.
		assert.Equal(t, "JOB_VIEW_ALL", r.URL.Query().Get("view"), "the full job view should be asked for")
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"id":"2026-09-22_00_00_00-1","name":"gw-library-test","type":"JOB_TYPE_BATCH","currentState":"JOB_STATE_RUNNING","labels":{"purpose":"terratest"},"environment":{"serviceAccountEmail":"gw-library-test@gw-library-test-project.iam.gserviceaccount.com"}}`))
	})

	job, err := gcp.GetDataflowJobAttrsWithClient(context.Background(), newFakeDataflowService(t, handler), "gw-library-test-project", "us-central1", "2026-09-22_00_00_00-1")
	require.NoError(t, err)

	assert.Equal(t, "gw-library-test", job.Name)
	assert.Equal(t, "JOB_TYPE_BATCH", job.Type)
	assert.Equal(t, "JOB_STATE_RUNNING", job.CurrentState)
	assert.Equal(t, "terratest", job.Labels["purpose"])
}

func TestGetDataflowJobAttrsWithClientMissingJob(t *testing.T) {
	t.Parallel()

	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	// The error names the job and everything that identifies it, and each of those is a value no
	// other part of the message contains, or its check could not fail.
	_, err := gcp.GetDataflowJobAttrsWithClient(context.Background(), newFakeDataflowService(t, handler), "gw-library-test-project", "us-central1", "gone")
	require.ErrorContains(t, err, "does not exist")
	require.ErrorContains(t, err, "gone")
	require.ErrorContains(t, err, "us-central1")
	require.ErrorContains(t, err, "gw-library-test-project")
}
