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
	"google.golang.org/api/datastream/v1"
	"google.golang.org/api/option"
)

// newFakeDatastreamService points a real Datastream client at a local test server, so the Google transport
// is exercised rather than a hand-written stand-in for a type we do not own.
func newFakeDatastreamService(t *testing.T, handler http.Handler) *datastream.Service {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	service, err := datastream.NewService(context.Background(),
		option.WithEndpoint(server.URL), option.WithoutAuthentication())
	require.NoError(t, err)

	return service
}

func TestGetDatastreamConnectionProfileAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a connection profile a Gruntwork module created, not a
	// copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/connectionProfiles/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/connectionProfiles/gw-library-test","displayName":"terratest profile","labels":{"purpose":"terratest"},"mysqlProfile":{"hostname":"10.70.0.1","port":3306,"username":"terratest"}}`))
	})

	result, err := gcp.GetDatastreamConnectionProfileAttrsWithClient(context.Background(), newFakeDatastreamService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "terratest profile", result.DisplayName)
	require.NotNil(t, result.MysqlProfile)
	assert.Equal(t, "10.70.0.1", result.MysqlProfile.Hostname)
	assert.Equal(t, int64(3306), result.MysqlProfile.Port)
}

func TestGetDatastreamConnectionProfileAttrsWithClientReportsAMissingOne(t *testing.T) {
	t.Parallel()

	// A caller who asks for something that is not there should be told that, rather than be handed
	// the transport's own wording for a 404.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error":{"code":404,"message":"Resource not found."}}`))
	})

	_, err := gcp.GetDatastreamConnectionProfileAttrsWithClient(context.Background(), newFakeDatastreamService(t, handler), "gw-library-test-project", "us-central1", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "does not exist")
}
