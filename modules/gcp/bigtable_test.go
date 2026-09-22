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
	"google.golang.org/api/bigtableadmin/v2"
	"google.golang.org/api/option"
)

// newFakeBigtableAdminService points a real Bigtable client at a local test server, so the Google transport is
// exercised rather than a hand-written stand-in for a type we do not own.
func newFakeBigtableAdminService(t *testing.T, handler http.Handler) *bigtableadmin.Service {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	service, err := bigtableadmin.NewService(context.Background(),
		option.WithEndpoint(server.URL), option.WithoutAuthentication())
	require.NoError(t, err)

	return service
}

func TestGetBigtableInstanceAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a instance the terraform-google-data-
	// storage module created, not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/instances/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/instances/gw-library-test","displayName":"terratest instance","type":"DEVELOPMENT","state":"READY","labels":{"purpose":"terratest"}}`))
	})

	instance, err := gcp.GetBigtableInstanceAttrsWithClient(context.Background(), newFakeBigtableAdminService(t, handler), "gw-library-test-project", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "terratest instance", instance.DisplayName)
	assert.Equal(t, "DEVELOPMENT", instance.Type)
	assert.Equal(t, "READY", instance.State)
	assert.Equal(t, "terratest", instance.Labels["purpose"])
}

func TestGetBigtableInstanceAttrsWithClientMissingInstance(t *testing.T) {
	t.Parallel()

	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	// The error names the instance and everything that identifies it, and each of those is a value
	// no other part of the message contains, or its check could not fail.
	_, err := gcp.GetBigtableInstanceAttrsWithClient(context.Background(), newFakeBigtableAdminService(t, handler), "gw-library-test-project", "gone")
	require.ErrorContains(t, err, "does not exist")
	require.ErrorContains(t, err, "gone")
	require.ErrorContains(t, err, "gw-library-test-project")
}

func TestGetBigtableClusterAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a cluster the terraform-google-data-
	// storage module created, not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/instances/gw-library-test/clusters/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/instances/gw-library-test/clusters/gw-library-test","location":"projects/gw-library-test-project/locations/us-central1-b","serveNodes":1,"defaultStorageType":"SSD","state":"READY"}`))
	})

	cluster, err := gcp.GetBigtableClusterAttrsWithClient(context.Background(), newFakeBigtableAdminService(t, handler),
		"gw-library-test-project", "gw-library-test", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "projects/gw-library-test-project/locations/us-central1-b", cluster.Location)
	assert.Equal(t, int64(1), cluster.ServeNodes)
	assert.Equal(t, "SSD", cluster.DefaultStorageType)
}

func TestGetBigtableClusterAttrsWithClientMissingCluster(t *testing.T) {
	t.Parallel()

	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error":{"code":404,"message":"not found"}}`))
	})

	_, err := gcp.GetBigtableClusterAttrsWithClient(context.Background(), newFakeBigtableAdminService(t, handler),
		"gw-library-test-project", "gw-library-test", "gone")
	require.Error(t, err)
	require.ErrorContains(t, err, "does not exist")
	require.ErrorContains(t, err, "gone")
	require.ErrorContains(t, err, "gw-library-test-project")
}
