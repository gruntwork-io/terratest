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
	osconfigv2 "google.golang.org/api/osconfig/v2"
)

// newFakeOSConfigV2Service points a real OS Config v2 client at a local test server, so the Google
// transport is exercised rather than a hand-written stand-in for a type we do not own.
func newFakeOSConfigV2Service(t *testing.T, handler http.Handler) *osconfigv2.Service {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	service, err := osconfigv2.NewService(context.Background(),
		option.WithEndpoint(server.URL), option.WithoutAuthentication())
	require.NoError(t, err)

	return service
}

func TestGetOSConfigPolicyOrchestratorAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for an orchestrator the policy orchestrator
	// module created, not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/global/policyOrchestrators/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/global/policyOrchestrators/gw-library-test","action":"UPSERT","state":"ACTIVE","description":"created by terratest","orchestrationScope":{"selectors":[{"locationSelector":{"includedLocations":["us-central1-a"]}}]},"orchestrationState":{"currentIterationState":{"state":"SUCCEEDED"}}}`))
	})

	orchestrator, err := gcp.GetOSConfigPolicyOrchestratorAttrsWithClient(context.Background(), newFakeOSConfigV2Service(t, handler), "gw-library-test-project", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "UPSERT", orchestrator.Action)
	assert.Equal(t, "ACTIVE", orchestrator.State)
	require.NotNil(t, orchestrator.OrchestrationState)
	require.NotNil(t, orchestrator.OrchestrationState.CurrentIterationState)
	assert.Equal(t, "SUCCEEDED", orchestrator.OrchestrationState.CurrentIterationState.State)
}

func TestGetOSConfigPolicyOrchestratorAttrsWithClientReportsAMissingOrchestrator(t *testing.T) {
	t.Parallel()

	// A caller who asks for an orchestrator that is not there should be told that, rather than be
	// handed the transport's own wording for a 404.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error":{"code":404,"message":"Not found."}}`))
	})

	_, err := gcp.GetOSConfigPolicyOrchestratorAttrsWithClient(context.Background(), newFakeOSConfigV2Service(t, handler), "gw-library-test-project", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "does not exist")
}
