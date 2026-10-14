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
	"google.golang.org/api/networkmanagement/v1"
	"google.golang.org/api/option"
)

// newFakeNetworkManagementService points a real Network Management client at a local test server, so the Google transport
// is exercised rather than a hand-written stand-in for a type we do not own.
func newFakeNetworkManagementService(t *testing.T, handler http.Handler) *networkmanagement.Service {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	service, err := networkmanagement.NewService(context.Background(),
		option.WithEndpoint(server.URL), option.WithoutAuthentication())
	require.NoError(t, err)

	return service
}

func TestGetConnectivityTestAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a test the terraform-google-networking connectivity test module created, not a
	// copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/v1/projects/gw-library-test-project/locations/global/connectivityTests/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/global/connectivityTests/gw-library-test","description":"created by terratest","protocol":"TCP","source":{"ipAddress":"10.0.0.1","network":"projects/gw-library-test-project/global/networks/gw-library-test"},"destination":{"ipAddress":"10.0.0.2","port":8080},"labels":{"purpose":"terratest"}}`))
	})

	connectivityTest, err := gcp.GetConnectivityTestAttrsWithClient(context.Background(), newFakeNetworkManagementService(t, handler), "gw-library-test-project", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "TCP", connectivityTest.Protocol)
	require.NotNil(t, connectivityTest.Source)
	assert.Equal(t, "10.0.0.1", connectivityTest.Source.IpAddress)
	require.NotNil(t, connectivityTest.Destination)
	assert.Equal(t, int64(8080), connectivityTest.Destination.Port)
}

func TestGetNetworkManagementVpcFlowLogsConfigAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a VPC flow logs config a Gruntwork module created, not a
	// copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/vpcFlowLogsConfigs/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"state":"ENABLED","aggregationInterval":"INTERVAL_5_SEC","flowSampling":0.5,"metadata":"INCLUDE_ALL_METADATA","interconnectAttachment":"projects/gw-library-test-project/regions/us-central1/interconnectAttachments/gw-library-test"}`))
	})

	result, err := gcp.GetNetworkManagementVpcFlowLogsConfigAttrsWithClient(context.Background(), newFakeNetworkManagementService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "ENABLED", result.State)
	assert.Equal(t, "INTERVAL_5_SEC", result.AggregationInterval)
	assert.InDelta(t, 0.5, result.FlowSampling, 0.0001)
	assert.Contains(t, result.InterconnectAttachment, "/interconnectAttachments/gw-library-test")
}

func TestGetNetworkManagementVpcFlowLogsConfigAttrsWithClientReportsAMissingConfig(t *testing.T) {
	t.Parallel()

	// A caller who names a VPC flow logs config that is not there should read a sentence about that VPC flow logs config, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetNetworkManagementVpcFlowLogsConfigAttrsWithClient(context.Background(), newFakeNetworkManagementService(t, handler), "gw-library-test-project", "us-central1", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}
