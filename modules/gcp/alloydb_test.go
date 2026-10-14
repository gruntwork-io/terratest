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

func TestGetAlloyDBBackupAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a AlloyDB backup a Gruntwork module created, not a
	// copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/backups/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"clusterName":"projects/gw-library-test-project/locations/us-central1/clusters/gw-library-test","type":"ON_DEMAND","state":"READY","sizeBytes":"1048576","databaseVersion":"POSTGRES_15"}`))
	})

	result, err := gcp.GetAlloyDBBackupAttrsWithClient(context.Background(), newFakeAlloyDBService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Contains(t, result.ClusterName, "/clusters/gw-library-test")
	assert.Equal(t, "ON_DEMAND", result.Type)
	assert.Equal(t, "READY", result.State)
	assert.Equal(t, int64(1048576), result.SizeBytes)
}

func TestGetAlloyDBUserAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a AlloyDB user a Gruntwork module created, not a
	// copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/clusters/gw-library-test/users/gw_library_user"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/clusters/gw-library-test/users/gw_library_user","userType":"ALLOYDB_BUILT_IN","databaseRoles":["alloydbsuperuser"]}`))
	})

	result, err := gcp.GetAlloyDBUserAttrsWithClient(context.Background(), newFakeAlloyDBService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test", "gw_library_user")
	require.NoError(t, err)

	assert.Equal(t, "ALLOYDB_BUILT_IN", result.UserType)
	assert.Equal(t, []string{"alloydbsuperuser"}, result.DatabaseRoles)
}

func TestGetAlloyDBBackupAttrsWithClientReportsAMissingBackup(t *testing.T) {
	t.Parallel()

	// A caller who names a AlloyDB backup that is not there should read a sentence about that AlloyDB backup, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetAlloyDBBackupAttrsWithClient(context.Background(), newFakeAlloyDBService(t, handler), "gw-library-test-project", "us-central1", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}
