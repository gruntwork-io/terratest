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
	"google.golang.org/api/observability/v1"
	"google.golang.org/api/option"
)

// newFakeObservabilityService points a real Observability client at a local test server, so the Google transport
// is exercised rather than a hand-written stand-in for a type we do not own.
func newFakeObservabilityService(t *testing.T, handler http.Handler) *observability.Service {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	service, err := observability.NewService(context.Background(),
		option.WithEndpoint(server.URL), option.WithoutAuthentication())
	require.NoError(t, err)

	return service
}

func TestGetTraceScopeAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a scope the terraform-google-observability trace scope module created, not a
	// copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/global/traceScopes/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/global/traceScopes/gw-library-test","description":"created by terratest","resourceNames":["projects/gw-library-test-project"]}`))
	})

	scope, err := gcp.GetTraceScopeAttrsWithClient(context.Background(), newFakeObservabilityService(t, handler), "gw-library-test-project", "global", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "created by terratest", scope.Description)
	assert.Equal(t, []string{"projects/gw-library-test-project"}, scope.ResourceNames)
}

func TestGetTraceScopeAttrsWithClientReportsAMissingScope(t *testing.T) {
	t.Parallel()

	// A caller who asks for a scope that is not there should be told that, rather than be handed the
	// transport's own wording for a 404.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error":{"code":404,"message":"Resource not found."}}`))
	})

	_, err := gcp.GetTraceScopeAttrsWithClient(context.Background(), newFakeObservabilityService(t, handler), "gw-library-test-project", "global", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "does not exist")
}
