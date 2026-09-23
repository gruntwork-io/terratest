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
	"google.golang.org/api/alloydb/v1"
	"google.golang.org/api/option"
)

// newFakeAlloyDBService points a real AlloyDB client at a local test server, so the Google transport is
// exercised rather than a hand-written stand-in for a type we do not own.
func newFakeAlloyDBService(t *testing.T, handler http.Handler) *alloydb.Service {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	service, err := alloydb.NewService(context.Background(),
		option.WithEndpoint(server.URL), option.WithoutAuthentication())
	require.NoError(t, err)

	return service
}

func TestGetAlloyDBClusterAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a cluster the terraform-google-data-
	// storage module created, not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/clusters/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/clusters/gw-library-test","databaseVersion":"POSTGRES_15","state":"READY","labels":{"purpose":"terratest"},"networkConfig":{"network":"projects/gw-library-test-project/global/networks/gw-library-test"},"initialUser":{"user":"postgres"}}`))
	})

	cluster, err := gcp.GetAlloyDBClusterAttrsWithClient(context.Background(), newFakeAlloyDBService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "POSTGRES_15", cluster.DatabaseVersion)
	assert.Equal(t, "READY", cluster.State)
	require.NotNil(t, cluster.NetworkConfig)
	assert.True(t, strings.HasSuffix(cluster.NetworkConfig.Network, "/global/networks/gw-library-test"))
}

func TestGetAlloyDBClusterAttrsWithClientMissingCluster(t *testing.T) {
	t.Parallel()

	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	// The error names the cluster and everything that identifies it, and each of those is a value no
	// other part of the message contains, or its check could not fail.
	_, err := gcp.GetAlloyDBClusterAttrsWithClient(context.Background(), newFakeAlloyDBService(t, handler), "gw-library-test-project", "us-central1", "gone")
	require.ErrorContains(t, err, "does not exist")
	require.ErrorContains(t, err, "gone")
	require.ErrorContains(t, err, "us-central1")
	require.ErrorContains(t, err, "gw-library-test-project")
}
