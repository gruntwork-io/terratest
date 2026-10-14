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
	"google.golang.org/api/redis/v1"
)

// newFakeRedisService points a real Memorystore for Redis client at a local test server, so the Google transport is
// exercised rather than a hand-written stand-in for a type we do not own.
func newFakeRedisService(t *testing.T, handler http.Handler) *redis.Service {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	service, err := redis.NewService(context.Background(),
		option.WithEndpoint(server.URL), option.WithoutAuthentication())
	require.NoError(t, err)

	return service
}

func TestGetRedisInstanceAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a instance the terraform-google-data-
	// storage module created, not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/instances/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/instances/gw-library-test","displayName":"terratest instance","tier":"BASIC","memorySizeGb":1,"redisVersion":"REDIS_7_0","state":"READY","labels":{"purpose":"terratest"}}`))
	})

	instance, err := gcp.GetRedisInstanceAttrsWithClient(context.Background(), newFakeRedisService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "BASIC", instance.Tier)
	assert.Equal(t, int64(1), instance.MemorySizeGb)
	assert.Equal(t, "REDIS_7_0", instance.RedisVersion)
	assert.Equal(t, "READY", instance.State)
}

func TestGetRedisInstanceAttrsWithClientMissingInstance(t *testing.T) {
	t.Parallel()

	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	// The error names the instance and everything that identifies it, and each of those is a value
	// no other part of the message contains, or its check could not fail.
	_, err := gcp.GetRedisInstanceAttrsWithClient(context.Background(), newFakeRedisService(t, handler), "gw-library-test-project", "us-central1", "gone")
	require.ErrorContains(t, err, "does not exist")
	require.ErrorContains(t, err, "gone")
	require.ErrorContains(t, err, "us-central1")
	require.ErrorContains(t, err, "gw-library-test-project")
}

func TestGetRedisClusterAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a Redis cluster a Gruntwork module created, not a
	// copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/clusters/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"state":"ACTIVE","shardCount":3,"replicaCount":1,"authorizationMode":"AUTH_MODE_DISABLED","discoveryEndpoints":[{"address":"10.0.0.3","port":6379}]}`))
	})

	result, err := gcp.GetRedisClusterAttrsWithClient(context.Background(), newFakeRedisService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "ACTIVE", result.State)
	assert.Equal(t, int64(3), result.ShardCount)
	assert.Equal(t, int64(1), result.ReplicaCount)
	require.Len(t, result.DiscoveryEndpoints, 1)
	assert.Equal(t, int64(6379), result.DiscoveryEndpoints[0].Port)
}

func TestGetRedisClusterAttrsWithClientReportsAMissingCluster(t *testing.T) {
	t.Parallel()

	// A caller who names a Redis cluster that is not there should read a sentence about that Redis cluster, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetRedisClusterAttrsWithClient(context.Background(), newFakeRedisService(t, handler), "gw-library-test-project", "us-central1", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}
