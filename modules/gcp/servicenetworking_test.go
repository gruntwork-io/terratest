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
	"google.golang.org/api/servicenetworking/v1"
)

// newFakeServiceNetworkingService points a real Service Networking client at a local test server, so
// the Google transport is exercised rather than a hand-written stand-in for a type we do not own.
func newFakeServiceNetworkingService(t *testing.T, handler http.Handler) *servicenetworking.APIService {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	service, err := servicenetworking.NewService(context.Background(),
		option.WithEndpoint(server.URL), option.WithoutAuthentication())
	require.NoError(t, err)

	return service
}

// connectionList answers the list call with the given connections, and asserts that the network the
// read filtered on reached the request, since that filter is what makes a single answer meaningful.
func connectionList(t *testing.T, connections string) http.Handler {
	t.Helper()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Contains(t, r.URL.Path, "/services/servicenetworking.googleapis.com/connections", "unexpected path")
		assert.Equal(t, "projects/gw-library-test-project/global/networks/gw-library-test", r.URL.Query().Get("network"))
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"connections":[` + connections + `]}`))
	})
}

func TestGetServiceNetworkingConnectionAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a connection the
	// terraform-google-networking module created, not a copy of any one fixture's values.
	handler := connectionList(t, `{
		"peering":"servicenetworking-googleapis-com",
		"network":"projects/gw-library-test-project/global/networks/gw-library-test",
		"reservedPeeringRanges":["gw-library-test"]
	}`)

	connection, err := gcp.GetServiceNetworkingConnectionAttrsWithClient(context.Background(), newFakeServiceNetworkingService(t, handler), "gw-library-test-project", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "servicenetworking-googleapis-com", connection.Peering)
	assert.Equal(t, []string{"gw-library-test"}, connection.ReservedPeeringRanges)
}

func TestGetServiceNetworkingConnectionAttrsWithClientNoConnection(t *testing.T) {
	t.Parallel()

	// A network with no connection answers with an empty list rather than a 404, so the error says
	// none exists and names the network and the project, neither of which the other contains.
	_, err := gcp.GetServiceNetworkingConnectionAttrsWithClient(context.Background(), newFakeServiceNetworkingService(t, connectionList(t, "")), "gw-library-test-project", "gw-library-test")
	require.ErrorContains(t, err, "no private connection exists")
	require.ErrorContains(t, err, "network gw-library-test ")
	require.ErrorContains(t, err, "gw-library-test-project")
}

func TestGetServiceNetworkingConnectionAttrsWithClientSeveralConnections(t *testing.T) {
	t.Parallel()

	// More than one connection on the filtered network is the test's answer rather than a choice
	// this read should make, so the error says how many there are and names them.
	handler := connectionList(t, `{"peering":"first"},{"peering":"second"}`)

	_, err := gcp.GetServiceNetworkingConnectionAttrsWithClient(context.Background(), newFakeServiceNetworkingService(t, handler), "gw-library-test-project", "gw-library-test")
	require.ErrorContains(t, err, "has 2 private connections")
	require.ErrorContains(t, err, "first, second")
}

func TestGetServiceNetworkingPeeredDNSDomainAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The API offers no get for one domain, so the helper lists the network's domains and picks one out. The list here
	// holds two, and the one asked for is not first.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/services/servicenetworking.googleapis.com/projects/370000000001/global/networks/gw-library-network/peeredDnsDomains"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"peeredDnsDomains":[{"name":"gw-library-other","dnsSuffix":"other.example.com."},{"name":"gw-library-test","dnsSuffix":"terratest.example.com."}]}`))
	})

	domain, err := gcp.GetServiceNetworkingPeeredDNSDomainAttrsWithClient(context.Background(), newFakeServiceNetworkingService(t, handler), "370000000001", "gw-library-network", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "gw-library-test", domain.Name)
	assert.Equal(t, "terratest.example.com.", domain.DnsSuffix)
}

func TestGetServiceNetworkingPeeredDNSDomainAttrsWithClientReportsAMissingDomain(t *testing.T) {
	t.Parallel()

	// A domain that is not in the list is reported as missing rather than as an empty result, because the list call
	// itself succeeded.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"peeredDnsDomains":[{"name":"gw-library-other","dnsSuffix":"other.example.com."}]}`))
	})

	_, err := gcp.GetServiceNetworkingPeeredDNSDomainAttrsWithClient(context.Background(), newFakeServiceNetworkingService(t, handler), "370000000001", "gw-library-network", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetServiceNetworkingVPCServiceControlsAttrsWithClient(t *testing.T) {
	t.Parallel()

	// Whether the reference architecture is enforced is an output only field on the consumer config, so a test can
	// read it but a module cannot claim it.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/services/servicenetworking.googleapis.com/projects/370000000001/global/networks/gw-library-network"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"vpcScReferenceArchitectureEnabled":true,"reservedRanges":[{"name":"gw-library-range","ipPrefixLength":16}]}`))
	})

	config, err := gcp.GetServiceNetworkingVPCServiceControlsAttrsWithClient(context.Background(), newFakeServiceNetworkingService(t, handler), "370000000001", "gw-library-network")
	require.NoError(t, err)

	assert.True(t, config.VpcScReferenceArchitectureEnabled)
}

func TestGetServiceNetworkingVPCServiceControlsAttrsWithClientReportsAMissingConnection(t *testing.T) {
	t.Parallel()

	// A network with no private connection answers 404, and the sentence has to name the connection rather than imply
	// the network is gone.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetServiceNetworkingVPCServiceControlsAttrsWithClient(context.Background(), newFakeServiceNetworkingService(t, handler), "370000000001", "gw-library-network")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-network")
	assert.Contains(t, err.Error(), "no service networking connection")
}
