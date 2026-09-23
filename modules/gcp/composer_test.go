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
	"google.golang.org/api/composer/v1"
	"google.golang.org/api/option"
)

// newFakeComposerService points a real Cloud Composer client at a local test server, so the Google transport is
// exercised rather than a hand-written stand-in for a type we do not own.
func newFakeComposerService(t *testing.T, handler http.Handler) *composer.Service {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	service, err := composer.NewService(context.Background(),
		option.WithEndpoint(server.URL), option.WithoutAuthentication())
	require.NoError(t, err)

	return service
}

func TestGetComposerEnvironmentAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a environment the terraform-google-
	// data-analytics module created, not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/environments/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/environments/gw-library-test","state":"RUNNING","labels":{"purpose":"terratest"},"config":{"softwareConfig":{"imageVersion":"composer-3-airflow-2","envVariables":{"PURPOSE":"terratest"}},"environmentSize":"ENVIRONMENT_SIZE_SMALL"}}`))
	})

	environment, err := gcp.GetComposerEnvironmentAttrsWithClient(context.Background(), newFakeComposerService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "RUNNING", environment.State)
	require.NotNil(t, environment.Config)
	require.NotNil(t, environment.Config.SoftwareConfig)
	assert.Equal(t, "terratest", environment.Config.SoftwareConfig.EnvVariables["PURPOSE"])
	assert.Equal(t, "ENVIRONMENT_SIZE_SMALL", environment.Config.EnvironmentSize)
}

func TestGetComposerEnvironmentAttrsWithClientMissingEnvironment(t *testing.T) {
	t.Parallel()

	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	// The error names the environment and everything that identifies it, and each of those is a
	// value no other part of the message contains, or its check could not fail.
	_, err := gcp.GetComposerEnvironmentAttrsWithClient(context.Background(), newFakeComposerService(t, handler), "gw-library-test-project", "us-central1", "gone")
	require.ErrorContains(t, err, "does not exist")
	require.ErrorContains(t, err, "gone")
	require.ErrorContains(t, err, "us-central1")
	require.ErrorContains(t, err, "gw-library-test-project")
}
