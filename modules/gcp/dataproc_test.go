package gcp_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gruntwork-io/terratest/modules/gcp/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/dataproc/v1"
	"google.golang.org/api/option"
)

// newFakeDataprocService points a real Dataproc client at a local test server, so the Google transport is
// exercised rather than a hand-written stand-in for a type we do not own.
func newFakeDataprocService(t *testing.T, handler http.Handler) *dataproc.Service {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	service, err := dataproc.NewService(context.Background(),
		option.WithEndpoint(server.URL), option.WithoutAuthentication())
	require.NoError(t, err)

	return service
}

func TestGetDataprocClusterAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a cluster the terraform-google-data-
	// analytics module created, not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Contains(t, r.URL.Path, "/projects/gw-library-test-project/regions/us-central1/clusters/gw-library-test", "unexpected path")
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"clusterName":"gw-library-test","status":{"state":"RUNNING"},"labels":{"purpose":"terratest"},"config":{"masterConfig":{"numInstances":1,"machineTypeUri":"n2-standard-2"},"softwareConfig":{"imageVersion":"2.2-debian12"}}}`))
	})

	cluster, err := gcp.GetDataprocClusterAttrsWithClient(context.Background(), newFakeDataprocService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "gw-library-test", cluster.ClusterName)
	require.NotNil(t, cluster.Status)
	assert.Equal(t, "RUNNING", cluster.Status.State)
	require.NotNil(t, cluster.Config)
	require.NotNil(t, cluster.Config.MasterConfig)
	assert.Equal(t, int64(1), cluster.Config.MasterConfig.NumInstances)
}

func TestGetDataprocClusterAttrsWithClientMissingCluster(t *testing.T) {
	t.Parallel()

	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	// The error names the cluster and everything that identifies it, and each of those is a value no
	// other part of the message contains, or its check could not fail.
	_, err := gcp.GetDataprocClusterAttrsWithClient(context.Background(), newFakeDataprocService(t, handler), "gw-library-test-project", "us-central1", "gone")
	require.ErrorContains(t, err, "does not exist")
	require.ErrorContains(t, err, "gone")
	require.ErrorContains(t, err, "us-central1")
	require.ErrorContains(t, err, "gw-library-test-project")
}

func TestNewDataprocServiceERefusesABadRegion(t *testing.T) {
	t.Parallel()

	// Each of these would build a URL pointing somewhere other than Google, so the constructor has
	// to refuse them before the endpoint is built.
	for _, region := range []string{"us/../evil.com", "evil.com", "us:8080", "user@evil.com", "US", ""} {
		_, err := gcp.NewDataprocServiceE(t, context.Background(), region)
		require.ErrorContains(t, err, "not a valid region", "region %q should be refused", region)
	}

	_, err := gcp.NewDataprocServiceE(t, context.Background(), "us-central1")
	require.NoError(t, err)
}
