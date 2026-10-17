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
	"google.golang.org/api/artifactregistry/v1"
	"google.golang.org/api/option"
)

// newFakeArtifactRegistryService points a real *artifactregistry.Service at a local test server, so
// the Google transport is exercised rather than a hand-written stand-in for a type we do not own.
func newFakeArtifactRegistryService(t *testing.T, handler http.Handler) *artifactregistry.Service {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	service, err := artifactregistry.NewService(context.Background(),
		option.WithEndpoint(server.URL), option.WithoutAuthentication())
	require.NoError(t, err)

	return service
}

func TestGetArtifactRegistryRepositoryAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The values are the ones the terraform-google-containers repository module sets, because the
	// point of reading settings back is asserting a module configured the repository it was asked
	// for.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/repositories/gw-library-test"))
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"name":"projects/gw-library-test-project/locations/us-central1/repositories/gw-library-test",
			"format":"DOCKER",
			"mode":"STANDARD_REPOSITORY",
			"description":"created by terratest",
			"labels":{"purpose":"terratest"}
		}`))
	})

	repository, err := gcp.GetArtifactRegistryRepositoryAttrsWithClient(context.Background(), newFakeArtifactRegistryService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "DOCKER", repository.Format)
	assert.Equal(t, "STANDARD_REPOSITORY", repository.Mode)
	assert.Equal(t, "created by terratest", repository.Description)
	assert.Equal(t, map[string]string{"purpose": "terratest"}, repository.Labels)
}

func TestGetArtifactRegistryRepositoryAttrsWithClientMissingRepository(t *testing.T) {
	t.Parallel()

	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	// The error names the repository and the project as well as saying it is absent, so all three
	// are asserted rather than only the phrase.
	_, err := gcp.GetArtifactRegistryRepositoryAttrsWithClient(context.Background(), newFakeArtifactRegistryService(t, handler), "gw-library-test-project", "us-central1", "gone")
	require.ErrorContains(t, err, "does not exist")
	require.ErrorContains(t, err, "gone")
	require.ErrorContains(t, err, "gw-library-test-project")
}

func TestGetArtifactRegistryProjectConfigAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The config is one object per project and location, and what it holds is whether Artifact
	// Registry writes platform logs there and at which severity.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/projectConfig"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/projectConfig","platformLogsConfig":{"loggingState":"LOGGING_STATE_ENABLED","severityLevel":"VERBOSE"}}`))
	})

	config, err := gcp.GetArtifactRegistryProjectConfigAttrsWithClient(context.Background(), newFakeArtifactRegistryService(t, handler), "gw-library-test-project", "us-central1")
	require.NoError(t, err)

	require.NotNil(t, config.PlatformLogsConfig)
	assert.Equal(t, "LOGGING_STATE_ENABLED", config.PlatformLogsConfig.LoggingState)
	assert.Equal(t, "VERBOSE", config.PlatformLogsConfig.SeverityLevel)
}

func TestGetArtifactRegistryProjectConfigAttrsWithClientReportsAMissingConfig(t *testing.T) {
	t.Parallel()

	// A caller who names a location Artifact Registry does not answer for should read a sentence
	// about that project and location, not a status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetArtifactRegistryProjectConfigAttrsWithClient(context.Background(), newFakeArtifactRegistryService(t, handler), "gw-library-missing", "us-central1")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "no Artifact Registry config")
}
