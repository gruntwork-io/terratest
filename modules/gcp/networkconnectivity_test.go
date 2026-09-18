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
	"google.golang.org/api/networkconnectivity/v1"
	"google.golang.org/api/option"
)

// newFakeNetworkConnectivityService points a real *networkconnectivity.Service at a local test
// server, so the Google transport is exercised rather than a hand-written stand-in for a type we
// do not own.
func newFakeNetworkConnectivityService(t *testing.T, handler http.Handler) *networkconnectivity.Service {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	service, err := networkconnectivity.NewService(context.Background(),
		option.WithEndpoint(server.URL), option.WithoutAuthentication())
	require.NoError(t, err)

	return service
}

func TestGetNetworkConnectivityHubAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The values are the ones the terraform-google-networking hub module sets, because the point of
	// reading settings back is asserting a module configured the hub it was asked for.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/global/hubs/gw-library-test"))
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"name":"projects/gw-library-test-project/locations/global/hubs/gw-library-test",
			"description":"created by terratest",
			"state":"ACTIVE",
			"policyMode":"PRESET",
			"presetTopology":"MESH",
			"exportPsc":true,
			"labels":{"managed-by":"terratest"}
		}`))
	})

	hub, err := gcp.GetNetworkConnectivityHubAttrsWithClient(context.Background(), newFakeNetworkConnectivityService(t, handler), "gw-library-test-project", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "created by terratest", hub.Description)
	assert.Equal(t, "ACTIVE", hub.State)
	assert.Equal(t, "MESH", hub.PresetTopology)
	assert.True(t, hub.ExportPsc)
	assert.Equal(t, "terratest", hub.Labels["managed-by"])
}

func TestGetNetworkConnectivityHubAttrsWithClientMissingHub(t *testing.T) {
	t.Parallel()

	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	// The error names the hub and the project as well as saying it is absent, so all three are
	// asserted rather than only the phrase.
	_, err := gcp.GetNetworkConnectivityHubAttrsWithClient(context.Background(), newFakeNetworkConnectivityService(t, handler), "gw-library-test-project", "gone")
	require.ErrorContains(t, err, "does not exist")
	require.ErrorContains(t, err, "gone")
	require.ErrorContains(t, err, "gw-library-test-project")
}
