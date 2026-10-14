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
	"google.golang.org/api/ids/v1"
	"google.golang.org/api/option"
)

// newFakeCloudIDSService points a real Cloud IDS client at a local test server, so the Google transport is
// exercised rather than a hand-written stand-in for a type we do not own.
func newFakeCloudIDSService(t *testing.T, handler http.Handler) *ids.Service {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	service, err := ids.NewService(context.Background(),
		option.WithEndpoint(server.URL), option.WithoutAuthentication())
	require.NoError(t, err)

	return service
}

func TestGetCloudIDSEndpointAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a Cloud IDS endpoint a Gruntwork module created, not a
	// copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/endpoints/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"network":"projects/gw-library-test-project/global/networks/gw-library-test","severity":"INFORMATIONAL","state":"READY","trafficLogs":true,"endpointIp":"10.0.0.5","threatExceptions":["12345"]}`))
	})

	result, err := gcp.GetCloudIDSEndpointAttrsWithClient(context.Background(), newFakeCloudIDSService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Contains(t, result.Network, "/networks/gw-library-test")
	assert.Equal(t, "INFORMATIONAL", result.Severity)
	assert.Equal(t, "READY", result.State)
	assert.True(t, result.TrafficLogs)
	assert.Equal(t, []string{"12345"}, result.ThreatExceptions)
}

func TestGetCloudIDSEndpointAttrsWithClientReportsAMissingEndpoint(t *testing.T) {
	t.Parallel()

	// A caller who names a Cloud IDS endpoint that is not there should read a sentence about that Cloud IDS endpoint, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetCloudIDSEndpointAttrsWithClient(context.Background(), newFakeCloudIDSService(t, handler), "gw-library-test-project", "us-central1", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}
