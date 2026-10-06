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
	"google.golang.org/api/parallelstore/v1"
)

// newFakeParallelstoreService points a real Parallelstore client at an httptest server, so a read can be exercised
// against a response we control without reaching Google.
func newFakeParallelstoreService(t *testing.T, handler http.Handler) *parallelstore.Service {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	service, err := parallelstore.NewService(context.Background(),
		option.WithoutAuthentication(), option.WithEndpoint(server.URL), option.WithHTTPClient(server.Client()))
	require.NoError(t, err)

	return service
}

func TestGetParallelstoreInstanceAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a Parallelstore instance the terraform-google-data-storage Parallelstore instance module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/instances/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1-a/instances/gw-library-test","description":"created by terratest","capacityGib":"12000","deploymentType":"SCRATCH","state":"ACTIVE","labels":{"purpose":"terratest"}}`))
	})

	attrs, err := gcp.GetParallelstoreInstanceAttrsWithClient(context.Background(), newFakeParallelstoreService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "created by terratest", attrs.Description)
	assert.Equal(t, int64(12000), attrs.CapacityGib)
	assert.Equal(t, "SCRATCH", attrs.DeploymentType)
}

func TestGetParallelstoreInstanceAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a Parallelstore instance that is not there should read a sentence about that Parallelstore instance, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetParallelstoreInstanceAttrsWithClient(context.Background(), newFakeParallelstoreService(t, handler), "gw-library-test-project", "us-central1", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}
