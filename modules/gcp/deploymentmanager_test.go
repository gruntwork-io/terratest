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
	"google.golang.org/api/deploymentmanager/v2"
	"google.golang.org/api/option"
)

// newFakeDeploymentManagerService points a real Deployment Manager client at a local test server, so
// the Google transport is exercised rather than a hand-written stand-in for a type we do not own.
func newFakeDeploymentManagerService(t *testing.T, handler http.Handler) *deploymentmanager.Service {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	service, err := deploymentmanager.NewService(context.Background(),
		option.WithEndpoint(server.URL), option.WithoutAuthentication())
	require.NoError(t, err)

	return service
}

func TestGetDeploymentManagerDeploymentAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a deployment the Deployment Manager
	// module created, not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/gw-library-test-project/global/deployments/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"gw-library-test","id":"4242424242","description":"created by terratest","fingerprint":"gw-library-fingerprint","manifest":"https://www.googleapis.com/deploymentmanager/v2/projects/gw-library-test-project/global/deployments/gw-library-test/manifests/manifest-1"}`))
	})

	deployment, err := gcp.GetDeploymentManagerDeploymentAttrsWithClient(context.Background(), newFakeDeploymentManagerService(t, handler), "gw-library-test-project", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, uint64(4242424242), deployment.Id)
	assert.Equal(t, "gw-library-fingerprint", deployment.Fingerprint)
	assert.Contains(t, deployment.Manifest, "/manifests/manifest-1")
}

func TestGetDeploymentManagerDeploymentAttrsWithClientReportsAMissingDeployment(t *testing.T) {
	t.Parallel()

	// A caller who asks for a deployment that is not there should be told that, rather than be handed
	// the transport's own wording for a 404.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error":{"code":404,"message":"Not found."}}`))
	})

	_, err := gcp.GetDeploymentManagerDeploymentAttrsWithClient(context.Background(), newFakeDeploymentManagerService(t, handler), "gw-library-test-project", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "does not exist")
}
