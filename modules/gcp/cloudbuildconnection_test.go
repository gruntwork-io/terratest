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
	cloudbuildv2 "google.golang.org/api/cloudbuild/v2"
	"google.golang.org/api/option"
)

// newFakeCloudBuildV2Service points a real Cloud Build v2 client at a local test server, so the Google
// transport is exercised rather than a hand-written stand-in for a type we do not own.
func newFakeCloudBuildV2Service(t *testing.T, handler http.Handler) *cloudbuildv2.Service {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	service, err := cloudbuildv2.NewService(context.Background(),
		option.WithEndpoint(server.URL), option.WithoutAuthentication())
	require.NoError(t, err)

	return service
}

func TestGetCloudBuildConnectionIamPolicyAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a connection a Gruntwork module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/connections/gw-library-test:getIamPolicy"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"version":3,"etag":"BwXhqw==","bindings":[{"role":"roles/cloudbuild.connectionViewer","members":["serviceAccount:gw-library-test@gw-library-test-project.iam.gserviceaccount.com"],"condition":{"title":"until 2030","expression":"request.time < timestamp(\"2030-01-01T00:00:00Z\")"}}]}`))
	})

	result, err := gcp.GetCloudBuildConnectionIamPolicyAttrsWithClient(context.Background(), newFakeCloudBuildV2Service(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, int64(3), result.Version)
	require.Len(t, result.Bindings, 1)
	assert.Equal(t, "roles/cloudbuild.connectionViewer", result.Bindings[0].Role)
	assert.Equal(t, `request.time < timestamp("2030-01-01T00:00:00Z")`, result.Bindings[0].Condition.Expression)
}

func TestGetCloudBuildConnectionIamPolicyAttrsWithClientReportsAMissingConnection(t *testing.T) {
	t.Parallel()

	// A caller who names a connection that is not there should read a sentence about that connection,
	// not a status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetCloudBuildConnectionIamPolicyAttrsWithClient(context.Background(), newFakeCloudBuildV2Service(t, handler), "gw-library-test-project", "us-central1", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}
