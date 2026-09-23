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
	"google.golang.org/api/container/v1"
	"google.golang.org/api/option"
)

// newFakeGKEService points a real GKE client at a local test server, so the Google transport is
// exercised rather than a hand-written stand-in for a type we do not own.
func newFakeGKEService(t *testing.T, handler http.Handler) *container.Service {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	service, err := container.NewService(context.Background(),
		option.WithEndpoint(server.URL), option.WithoutAuthentication())
	require.NoError(t, err)

	return service
}

func TestGetGKEClusterAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a cluster the terraform-google-
	// containers module created, not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1-b/clusters/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"gw-library-test","description":"created by terratest","status":"RUNNING","location":"us-central1-b","initialClusterVersion":"1.31.1","releaseChannel":{"channel":"REGULAR"},"network":"gw-library-test","resourceLabels":{"purpose":"terratest"}}`))
	})

	cluster, err := gcp.GetGKEClusterAttrsWithClient(context.Background(), newFakeGKEService(t, handler), "gw-library-test-project", "us-central1-b", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "created by terratest", cluster.Description)
	assert.Equal(t, "RUNNING", cluster.Status)
	require.NotNil(t, cluster.ReleaseChannel)
	assert.Equal(t, "REGULAR", cluster.ReleaseChannel.Channel)
	assert.Equal(t, "terratest", cluster.ResourceLabels["purpose"])
}

func TestGetGKEClusterAttrsWithClientMissingCluster(t *testing.T) {
	t.Parallel()

	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	// The error names the cluster and everything that identifies it, and each of those is a value no
	// other part of the message contains, or its check could not fail.
	_, err := gcp.GetGKEClusterAttrsWithClient(context.Background(), newFakeGKEService(t, handler), "gw-library-test-project", "us-central1-b", "gone")
	require.ErrorContains(t, err, "does not exist")
	require.ErrorContains(t, err, "gone")
	require.ErrorContains(t, err, "us-central1-b")
	require.ErrorContains(t, err, "gw-library-test-project")
}

func TestGetGKENodePoolAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a node pool the terraform-google-
	// containers module created, not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1-b/clusters/gw-cluster/nodePools/gw-pool"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"gw-pool","initialNodeCount":1,"status":"RUNNING","config":{"machineType":"e2-small","diskSizeGb":32,"spot":true},"management":{"autoRepair":true,"autoUpgrade":false}}`))
	})

	nodePool, err := gcp.GetGKENodePoolAttrsWithClient(context.Background(), newFakeGKEService(t, handler), "gw-library-test-project", "us-central1-b", "gw-cluster", "gw-pool")
	require.NoError(t, err)

	assert.Equal(t, int64(1), nodePool.InitialNodeCount)
	require.NotNil(t, nodePool.Config)
	assert.Equal(t, "e2-small", nodePool.Config.MachineType)
	assert.True(t, nodePool.Config.Spot)
	require.NotNil(t, nodePool.Management)
	assert.False(t, nodePool.Management.AutoUpgrade)
}

func TestGetGKENodePoolAttrsWithClientMissingNodePool(t *testing.T) {
	t.Parallel()

	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	// The error names the node pool and everything that identifies it, and each of those is a value
	// no other part of the message contains, or its check could not fail.
	_, err := gcp.GetGKENodePoolAttrsWithClient(context.Background(), newFakeGKEService(t, handler), "gw-library-test-project", "us-central1-b", "gw-cluster", "gone")
	require.ErrorContains(t, err, "does not exist")
	require.ErrorContains(t, err, "gone")
	require.ErrorContains(t, err, "gw-cluster")
	require.ErrorContains(t, err, "gw-library-test-project")
}
