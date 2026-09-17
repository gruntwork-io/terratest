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
	"google.golang.org/api/cloudbuild/v1"
	"google.golang.org/api/option"
)

// newFakeCloudBuildRESTService points a real *cloudbuild.Service at a local test server, so the
// Google transport is exercised rather than a hand-written stand-in for a type we do not own.
func newFakeCloudBuildRESTService(t *testing.T, handler http.Handler) *cloudbuild.Service {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	service, err := cloudbuild.NewService(context.Background(),
		option.WithEndpoint(server.URL), option.WithoutAuthentication())
	require.NoError(t, err)

	return service
}

func TestGetWorkerPoolAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The values are the ones the terraform-google-devtools worker pool module sets, because the
	// point of reading settings back is asserting a module configured the pool it was asked for.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/workerPools/gw-library-test"))
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"name":"projects/gw-library-test-project/locations/us-central1/workerPools/gw-library-test",
			"displayName":"terratest pool",
			"state":"RUNNING",
			"privatePoolV1Config":{
				"workerConfig":{"machineType":"e2-medium","diskSizeGb":"100"},
				"networkConfig":{"egressOption":"NO_PUBLIC_EGRESS"}
			}
		}`))
	})

	pool, err := gcp.GetWorkerPoolAttrsWithClient(context.Background(), newFakeCloudBuildRESTService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "terratest pool", pool.DisplayName)
	assert.Equal(t, "RUNNING", pool.State)
	require.NotNil(t, pool.PrivatePoolV1Config)
	require.NotNil(t, pool.PrivatePoolV1Config.WorkerConfig)
	assert.Equal(t, "e2-medium", pool.PrivatePoolV1Config.WorkerConfig.MachineType)
	// The API sends the disk size as a JSON string and the Go client decodes it to an int64.
	assert.Equal(t, int64(100), pool.PrivatePoolV1Config.WorkerConfig.DiskSizeGb)
	require.NotNil(t, pool.PrivatePoolV1Config.NetworkConfig)
	assert.Equal(t, "NO_PUBLIC_EGRESS", pool.PrivatePoolV1Config.NetworkConfig.EgressOption)
}

func TestGetWorkerPoolAttrsWithClientMissingPool(t *testing.T) {
	t.Parallel()

	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	// The error names the pool and the project as well as saying it is absent, so all three are
	// asserted rather than only the phrase.
	_, err := gcp.GetWorkerPoolAttrsWithClient(context.Background(), newFakeCloudBuildRESTService(t, handler), "gw-library-test-project", "us-central1", "gone")
	require.ErrorContains(t, err, "does not exist")
	require.ErrorContains(t, err, "gone")
	require.ErrorContains(t, err, "gw-library-test-project")
}
