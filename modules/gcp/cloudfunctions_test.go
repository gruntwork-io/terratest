package gcp_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gruntwork-io/terratest/modules/gcp/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/cloudfunctions/v2"
	"google.golang.org/api/option"
)

// newFakeCloudFunctionsService points a real Cloud Functions client at a local test server, so the Google transport is
// exercised rather than a hand-written stand-in for a type we do not own.
func newFakeCloudFunctionsService(t *testing.T, handler http.Handler) *cloudfunctions.Service {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	service, err := cloudfunctions.NewService(context.Background(),
		option.WithEndpoint(server.URL), option.WithoutAuthentication())
	require.NoError(t, err)

	return service
}

func TestCloudFunctionAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a function the terraform-google-
	// serverless module created, not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Contains(t, r.URL.Path, "/locations/us-central1/functions/gw-library-test", "unexpected path")
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/functions/gw-library-test","description":"created by terratest","state":"ACTIVE","labels":{"purpose":"terratest"},"buildConfig":{"runtime":"go123","entryPoint":"Hello"},"serviceConfig":{"availableMemory":"256M","timeoutSeconds":30,"maxInstanceCount":2}}`))
	})

	function, err := gcp.GetCloudFunctionV2AttrsWithClient(context.Background(), newFakeCloudFunctionsService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "ACTIVE", function.State)
	require.NotNil(t, function.BuildConfig)
	assert.Equal(t, "go123", function.BuildConfig.Runtime)
	assert.Equal(t, "Hello", function.BuildConfig.EntryPoint)
	require.NotNil(t, function.ServiceConfig)
	assert.Equal(t, "256M", function.ServiceConfig.AvailableMemory)
}

func TestCloudFunctionAttrsWithClientMissingFunction(t *testing.T) {
	t.Parallel()

	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	// The error names the function and everything that identifies it, and each of those is a value
	// no other part of the message contains, or its check could not fail.
	_, err := gcp.GetCloudFunctionV2AttrsWithClient(context.Background(), newFakeCloudFunctionsService(t, handler), "gw-library-test-project", "us-central1", "gone")
	require.ErrorContains(t, err, "does not exist")
	require.ErrorContains(t, err, "gone")
	require.ErrorContains(t, err, "us-central1")
	require.ErrorContains(t, err, "gw-library-test-project")
}
