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
	"google.golang.org/api/memcache/v1"
	"google.golang.org/api/option"
)

// newFakeMemcacheService points a real Memorystore for Memcached client at a local test server, so the Google transport is
// exercised rather than a hand-written stand-in for a type we do not own.
func newFakeMemcacheService(t *testing.T, handler http.Handler) *memcache.Service {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	service, err := memcache.NewService(context.Background(),
		option.WithEndpoint(server.URL), option.WithoutAuthentication())
	require.NoError(t, err)

	return service
}

func TestGetMemcacheInstanceAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a Memcached instance a Gruntwork module created, not a
	// copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/instances/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"displayName":"terratest memcached","nodeCount":1,"state":"READY","memcacheVersion":"MEMCACHE_1_6_15","nodeConfig":{"cpuCount":1,"memorySizeMb":1024},"zones":["us-central1-a"],"authorizedNetwork":"projects/gw-library-test-project/global/networks/gw-library-test"}`))
	})

	result, err := gcp.GetMemcacheInstanceAttrsWithClient(context.Background(), newFakeMemcacheService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "terratest memcached", result.DisplayName)
	assert.Equal(t, int64(1), result.NodeCount)
	assert.Equal(t, "READY", result.State)
	require.NotNil(t, result.NodeConfig)
	assert.Equal(t, int64(1024), result.NodeConfig.MemorySizeMb)
	assert.Equal(t, []string{"us-central1-a"}, result.Zones)
}

func TestGetMemcacheInstanceAttrsWithClientReportsAMissingInstance(t *testing.T) {
	t.Parallel()

	// A caller who names a Memcached instance that is not there should read a sentence about that Memcached instance, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetMemcacheInstanceAttrsWithClient(context.Background(), newFakeMemcacheService(t, handler), "gw-library-test-project", "us-central1", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}
