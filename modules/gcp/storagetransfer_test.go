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
	"google.golang.org/api/storagetransfer/v1"
)

// newFakeStorageTransferService points a real *storagetransfer.Service at a local test server, so
// the Google transport is exercised rather than a hand-written stand-in for a type we do not own.
func newFakeStorageTransferService(t *testing.T, handler http.Handler) *storagetransfer.Service {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	service, err := storagetransfer.NewService(context.Background(),
		option.WithEndpoint(server.URL), option.WithoutAuthentication())
	require.NoError(t, err)

	return service
}

func TestGetAgentPoolAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The values are the ones the terraform-google-migration agent pool module sets, because the
	// point of reading settings back is asserting a module configured the pool it was asked for.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/agentPools/gw-library-test"))
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"name":"projects/gw-library-test-project/agentPools/gw-library-test",
			"displayName":"terratest pool",
			"state":"CREATED",
			"bandwidthLimit":{"limitMbps":"50"}
		}`))
	})

	pool, err := gcp.GetAgentPoolAttrsWithClient(context.Background(), newFakeStorageTransferService(t, handler), "gw-library-test-project", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "terratest pool", pool.DisplayName)
	assert.Equal(t, "CREATED", pool.State)
	require.NotNil(t, pool.BandwidthLimit)
	assert.Equal(t, int64(50), pool.BandwidthLimit.LimitMbps)
}

func TestGetAgentPoolAttrsWithClientMissingPool(t *testing.T) {
	t.Parallel()

	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	// The error names the pool and the project as well as saying it is absent, so all three are
	// asserted rather than only the phrase.
	_, err := gcp.GetAgentPoolAttrsWithClient(context.Background(), newFakeStorageTransferService(t, handler), "gw-library-test-project", "gone")
	require.ErrorContains(t, err, "does not exist")
	require.ErrorContains(t, err, "gone")
	require.ErrorContains(t, err, "gw-library-test-project")
}

func TestGetStorageTransferJobAttrsWithClient(t *testing.T) {
	t.Parallel()

	// A transfer job's name is unique across Google rather than within a project, so the project is a
	// query parameter rather than part of the path.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/transferJobs/1234567890"), "unexpected path %s", r.URL.Path)
		assert.Equal(t, "gw-library-test-project", r.URL.Query().Get("projectId"))
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"transferJobs/1234567890","description":"created by terratest","status":"ENABLED","projectId":"gw-library-test-project","transferSpec":{"gcsDataSource":{"bucketName":"gw-library-test-source"},"gcsDataSink":{"bucketName":"gw-library-test-sink"}}}`))
	})

	job, err := gcp.GetStorageTransferJobAttrsWithClient(context.Background(), newFakeStorageTransferService(t, handler), "gw-library-test-project", "transferJobs/1234567890")
	require.NoError(t, err)

	assert.Equal(t, "created by terratest", job.Description)
	assert.Equal(t, "ENABLED", job.Status)
	require.NotNil(t, job.TransferSpec)
	require.NotNil(t, job.TransferSpec.GcsDataSource)
	assert.Equal(t, "gw-library-test-source", job.TransferSpec.GcsDataSource.BucketName)
	assert.Equal(t, "gw-library-test-sink", job.TransferSpec.GcsDataSink.BucketName)
}
