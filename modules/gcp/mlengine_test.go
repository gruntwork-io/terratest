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
	"google.golang.org/api/ml/v1"
	"google.golang.org/api/option"
)

// newFakeMLEngineService points a real AI Platform client at a local test server, so the Google
// transport is exercised rather than a hand-written stand-in for a type we do not own.
func newFakeMLEngineService(t *testing.T, handler http.Handler) *ml.Service {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	service, err := ml.NewService(context.Background(),
		option.WithEndpoint(server.URL), option.WithoutAuthentication())
	require.NoError(t, err)

	return service
}

func TestGetMLEngineModelAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a model the AI Platform model module
	// created, not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/models/gw_library_test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/models/gw_library_test","description":"created by terratest","regions":["us-central1"],"onlinePredictionLogging":true,"labels":{"terratest":"true"}}`))
	})

	model, err := gcp.GetMLEngineModelAttrsWithClient(context.Background(), newFakeMLEngineService(t, handler), "gw-library-test-project", "gw_library_test")
	require.NoError(t, err)

	assert.Equal(t, []string{"us-central1"}, model.Regions)
	assert.True(t, model.OnlinePredictionLogging)
	assert.Equal(t, map[string]string{"terratest": "true"}, model.Labels)
}

func TestGetMLEngineModelAttrsWithClientReportsAMissingModel(t *testing.T) {
	t.Parallel()

	// A caller who asks for a model that is not there should be told that, rather than be handed the
	// transport's own wording for a 404.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error":{"code":404,"message":"Not found."}}`))
	})

	_, err := gcp.GetMLEngineModelAttrsWithClient(context.Background(), newFakeMLEngineService(t, handler), "gw-library-test-project", "gw_library_missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "does not exist")
}
