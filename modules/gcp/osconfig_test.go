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
	"google.golang.org/api/osconfig/v1"
)

// newFakeOSConfigService points a real OS Config client at a local test server, so the Google transport
// is exercised rather than a hand-written stand-in for a type we do not own.
func newFakeOSConfigService(t *testing.T, handler http.Handler) *osconfig.Service {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	service, err := osconfig.NewService(context.Background(),
		option.WithEndpoint(server.URL), option.WithoutAuthentication())
	require.NoError(t, err)

	return service
}

func TestGetPatchDeploymentAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a deployment the terraform-google-management patch deployment module created, not a
	// copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/patchDeployments/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/patchDeployments/gw-library-test","description":"created by terratest","duration":"3600s","instanceFilter":{"groupLabels":[{"labels":{"purpose":"terratest"}}]},"patchConfig":{"rebootConfig":"NEVER"}}`))
	})

	deployment, err := gcp.GetPatchDeploymentAttrsWithClient(context.Background(), newFakeOSConfigService(t, handler), "gw-library-test-project", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "3600s", deployment.Duration)
	require.NotNil(t, deployment.PatchConfig)
	assert.Equal(t, "NEVER", deployment.PatchConfig.RebootConfig)
	require.NotNil(t, deployment.InstanceFilter)
	require.Len(t, deployment.InstanceFilter.GroupLabels, 1)
	assert.Equal(t, map[string]string{"purpose": "terratest"}, deployment.InstanceFilter.GroupLabels[0].Labels)
}

func TestGetOSPolicyAssignmentAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for an assignment the terraform-google-management OS policy assignment module created, not a
	// copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1-a/osPolicyAssignments/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1-a/osPolicyAssignments/gw-library-test","description":"created by terratest","osPolicies":[{"id":"terratest-policy","mode":"VALIDATION"}],"rollout":{"minWaitDuration":"60s"}}`))
	})

	assignment, err := gcp.GetOSPolicyAssignmentAttrsWithClient(context.Background(), newFakeOSConfigService(t, handler), "gw-library-test-project", "us-central1-a", "gw-library-test")
	require.NoError(t, err)

	require.Len(t, assignment.OsPolicies, 1)
	assert.Equal(t, "terratest-policy", assignment.OsPolicies[0].Id)
	assert.Equal(t, "VALIDATION", assignment.OsPolicies[0].Mode)
	require.NotNil(t, assignment.Rollout)
	assert.Equal(t, "60s", assignment.Rollout.MinWaitDuration)
}

func TestGetPatchDeploymentAttrsWithClientReportsAMissingDeployment(t *testing.T) {
	t.Parallel()

	// A caller who asks for a deployment that is not there should be told that, rather than be handed
	// the transport's own wording for a 404.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error":{"code":404,"message":"Resource not found."}}`))
	})

	_, err := gcp.GetPatchDeploymentAttrsWithClient(context.Background(), newFakeOSConfigService(t, handler), "gw-library-test-project", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetOSPolicyAssignmentAttrsWithClientReportsAMissingAssignment(t *testing.T) {
	t.Parallel()

	// A caller who asks for an assignment that is not there should be told that, rather than be
	// handed the transport's own wording for a 404.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error":{"code":404,"message":"Resource not found."}}`))
	})

	_, err := gcp.GetOSPolicyAssignmentAttrsWithClient(context.Background(), newFakeOSConfigService(t, handler), "gw-library-test-project", "us-central1-a", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "does not exist")
}
