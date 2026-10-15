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
	notebooksv2 "google.golang.org/api/notebooks/v2"
	"google.golang.org/api/option"
)

// newFakeWorkbenchService points a real Workbench client at a local test server, so the Google
// transport is exercised rather than a hand-written stand-in for a type we do not own.
func newFakeWorkbenchService(t *testing.T, handler http.Handler) *notebooksv2.Service {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	service, err := notebooksv2.NewService(context.Background(),
		option.WithEndpoint(server.URL), option.WithoutAuthentication())
	require.NoError(t, err)

	return service
}

func TestGetWorkbenchInstanceAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for an instance the Workbench instance module
	// created, not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1-a/instances/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1-a/instances/gw-library-test","state":"ACTIVE","healthState":"HEALTHY","proxyUri":"gw-library-test-dot-us-central1.notebooks.googleusercontent.com","disableProxyAccess":false,"instanceOwners":["terratest@example.com"],"gceSetup":{"machineType":"e2-standard-4","bootDisk":{"diskSizeGb":"150","diskType":"PD_BALANCED"}}}`))
	})

	instance, err := gcp.GetWorkbenchInstanceAttrsWithClient(context.Background(), newFakeWorkbenchService(t, handler), "gw-library-test-project", "us-central1-a", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "ACTIVE", instance.State)
	assert.Equal(t, []string{"terratest@example.com"}, instance.InstanceOwners)
	require.NotNil(t, instance.GceSetup)
	assert.Equal(t, "e2-standard-4", instance.GceSetup.MachineType)
}

func TestGetWorkbenchInstanceAttrsWithClientReportsAMissingInstance(t *testing.T) {
	t.Parallel()

	// A caller who asks for an instance that is not there should be told that, rather than be handed
	// the transport's own wording for a 404.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error":{"code":404,"message":"Not found."}}`))
	})

	_, err := gcp.GetWorkbenchInstanceAttrsWithClient(context.Background(), newFakeWorkbenchService(t, handler), "gw-library-test-project", "us-central1-a", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "does not exist")
}
