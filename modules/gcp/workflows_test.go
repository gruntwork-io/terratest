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
	"google.golang.org/api/workflows/v1"
)

// newFakeWorkflowsService points a real *workflows.Service at a local test server, so the Google
// transport is exercised rather than a hand-written stand-in for a type we do not own.
func newFakeWorkflowsService(t *testing.T, handler http.Handler) *workflows.Service {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	service, err := workflows.NewService(context.Background(),
		option.WithEndpoint(server.URL), option.WithoutAuthentication())
	require.NoError(t, err)

	return service
}

func TestGetWorkflowAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The values are the ones the terraform-google-serverless workflow module sets, because the
	// point of reading settings back is asserting a module configured the workflow it was asked
	// for.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/locations/us-central1/workflows/gw-library-test"))
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"name":"projects/gw-library-test-project/locations/us-central1/workflows/gw-library-test",
			"description":"created by terratest",
			"state":"ACTIVE",
			"callLogLevel":"LOG_ERRORS_ONLY",
			"labels":{"purpose":"terratest"}
		}`))
	})

	workflow, err := gcp.GetWorkflowAttrsWithClient(context.Background(), newFakeWorkflowsService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "created by terratest", workflow.Description)
	assert.Equal(t, "ACTIVE", workflow.State)
	assert.Equal(t, "LOG_ERRORS_ONLY", workflow.CallLogLevel)
	assert.Equal(t, map[string]string{"purpose": "terratest"}, workflow.Labels)
}

func TestGetWorkflowAttrsWithClientMissingWorkflow(t *testing.T) {
	t.Parallel()

	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	// The error names the workflow and the project as well as saying it is absent, so all three are
	// asserted rather than only the phrase.
	_, err := gcp.GetWorkflowAttrsWithClient(context.Background(), newFakeWorkflowsService(t, handler), "gw-library-test-project", "us-central1", "gone")
	require.ErrorContains(t, err, "does not exist")
	require.ErrorContains(t, err, "gone")
	require.ErrorContains(t, err, "gw-library-test-project")
}
