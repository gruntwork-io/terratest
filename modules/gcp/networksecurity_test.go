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
	"google.golang.org/api/networksecurity/v1"
	"google.golang.org/api/option"
)

// newFakeNetworkSecurityService points a real Network Security client at a local test server, so the Google transport
// is exercised rather than a hand-written stand-in for a type we do not own.
func newFakeNetworkSecurityService(t *testing.T, handler http.Handler) *networksecurity.Service {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	service, err := networksecurity.NewService(context.Background(),
		option.WithEndpoint(server.URL), option.WithoutAuthentication())
	require.NoError(t, err)

	return service
}

func TestGetNetworkSecurityAddressGroupAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a group the terraform-google-networking address group module created, not a
	// copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/addressGroups/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/addressGroups/gw-library-test","description":"created by terratest","type":"IPV4","capacity":100,"items":["10.10.0.0/16"]}`))
	})

	group, err := gcp.GetNetworkSecurityAddressGroupAttrsWithClient(context.Background(), newFakeNetworkSecurityService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "IPV4", group.Type)
	assert.Equal(t, int64(100), group.Capacity)
	assert.Equal(t, []string{"10.10.0.0/16"}, group.Items)
}

func TestGetNetworkSecurityClientTLSPolicyAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a policy the terraform-google-networking client TLS policy module created, not a
	// copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/clientTlsPolicies/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/clientTlsPolicies/gw-library-test","description":"created by terratest","sni":"terratest.example.com","labels":{"purpose":"terratest"}}`))
	})

	policy, err := gcp.GetNetworkSecurityClientTLSPolicyAttrsWithClient(context.Background(), newFakeNetworkSecurityService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "terratest.example.com", policy.Sni)
	assert.Equal(t, "created by terratest", policy.Description)
}

func TestGetNetworkSecurityServerTLSPolicyAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a policy the terraform-google-networking server TLS policy module created, not a
	// copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/serverTlsPolicies/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/serverTlsPolicies/gw-library-test","description":"created by terratest","allowOpen":true,"labels":{"purpose":"terratest"}}`))
	})

	policy, err := gcp.GetNetworkSecurityServerTLSPolicyAttrsWithClient(context.Background(), newFakeNetworkSecurityService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.True(t, policy.AllowOpen)
	assert.Equal(t, "created by terratest", policy.Description)
}

func TestGetNetworkSecurityURLListAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a list the terraform-google-networking URL list module created, not a
	// copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/urlLists/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/urlLists/gw-library-test","description":"created by terratest","values":["terratest.example.com/*"]}`))
	})

	list, err := gcp.GetNetworkSecurityURLListAttrsWithClient(context.Background(), newFakeNetworkSecurityService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, []string{"terratest.example.com/*"}, list.Values)
}

func TestGetNetworkSecurityGatewaySecurityPolicyAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a policy the terraform-google-networking gateway security policy module created, not a
	// copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/gatewaySecurityPolicies/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/gatewaySecurityPolicies/gw-library-test","description":"created by terratest"}`))
	})

	policy, err := gcp.GetNetworkSecurityGatewaySecurityPolicyAttrsWithClient(context.Background(), newFakeNetworkSecurityService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "created by terratest", policy.Description)
}

func TestGetNetworkSecurityAddressGroupIamPolicyAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a group the terraform-google-networking
	// address group IAM policy module granted access on, not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/addressGroups/gw-library-test:getIamPolicy"), "unexpected path %s", r.URL.Path)

		// A conditional binding only comes back at version 3, so the read has to ask for it.
		assert.Equal(t, "3", r.URL.Query().Get("options.requestedPolicyVersion"))
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"version":3,"etag":"BwXhqw==","bindings":[{"role":"roles/compute.networkViewer","members":["serviceAccount:gw-library-test@gw-library-test-project.iam.gserviceaccount.com"],"condition":{"title":"until 2030","expression":"request.time < timestamp(\"2030-01-01T00:00:00Z\")"}}]}`))
	})

	policy, err := gcp.GetNetworkSecurityAddressGroupIamPolicyAttrsWithClient(context.Background(), newFakeNetworkSecurityService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	require.Len(t, policy.Bindings, 1)
	assert.Equal(t, int64(3), policy.Version)
	require.NotNil(t, policy.Bindings[0].Condition, "a conditional binding should keep its condition")
	assert.Equal(t, `request.time < timestamp("2030-01-01T00:00:00Z")`, policy.Bindings[0].Condition.Expression)
}
