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

func TestGetNetworkSecurityAuthzPolicyAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a authorization policy the terraform-google-networking authorization policy module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/authzPolicies/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/authzPolicies/gw-library-test","description":"created by terratest","action":"DENY","labels":{"purpose":"terratest"}}`))
	})

	attrs, err := gcp.GetNetworkSecurityAuthzPolicyAttrsWithClient(context.Background(), newFakeNetworkSecurityService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "created by terratest", attrs.Description)
	assert.Equal(t, "DENY", attrs.Action)
}

func TestGetNetworkSecurityAuthzPolicyAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a authorization policy that is not there should read a sentence about that authorization policy, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetNetworkSecurityAuthzPolicyAttrsWithClient(context.Background(), newFakeNetworkSecurityService(t, handler), "gw-library-test-project", "us-central1", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetNetworkSecurityBackendAuthenticationConfigAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a backend authentication config the terraform-google-networking backend authentication config module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/backendAuthenticationConfigs/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/backendAuthenticationConfigs/gw-library-test","description":"created by terratest","wellKnownRoots":"PUBLIC_ROOTS","labels":{"purpose":"terratest"}}`))
	})

	attrs, err := gcp.GetNetworkSecurityBackendAuthenticationConfigAttrsWithClient(context.Background(), newFakeNetworkSecurityService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "created by terratest", attrs.Description)
	assert.Equal(t, "PUBLIC_ROOTS", attrs.WellKnownRoots)
}

func TestGetNetworkSecurityBackendAuthenticationConfigAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a backend authentication config that is not there should read a sentence about that backend authentication config, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetNetworkSecurityBackendAuthenticationConfigAttrsWithClient(context.Background(), newFakeNetworkSecurityService(t, handler), "gw-library-test-project", "us-central1", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetNetworkSecurityDNSThreatDetectorAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a DNS threat detector the terraform-google-networking DNS threat detector module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/dnsThreatDetectors/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/dnsThreatDetectors/gw-library-test","provider":"INFOBLOX","labels":{"purpose":"terratest"}}`))
	})

	attrs, err := gcp.GetNetworkSecurityDNSThreatDetectorAttrsWithClient(context.Background(), newFakeNetworkSecurityService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "INFOBLOX", attrs.Provider)
}

func TestGetNetworkSecurityDNSThreatDetectorAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a DNS threat detector that is not there should read a sentence about that DNS threat detector, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetNetworkSecurityDNSThreatDetectorAttrsWithClient(context.Background(), newFakeNetworkSecurityService(t, handler), "gw-library-test-project", "us-central1", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetNetworkSecurityFirewallEndpointAssociationAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a firewall endpoint association the terraform-google-networking firewall endpoint association module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1-a/firewallEndpointAssociations/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1-a/firewallEndpointAssociations/gw-library-test","network":"projects/gw-library-test-project/global/networks/gw-library-test","state":"ACTIVE","labels":{"purpose":"terratest"}}`))
	})

	attrs, err := gcp.GetNetworkSecurityFirewallEndpointAssociationAttrsWithClient(context.Background(), newFakeNetworkSecurityService(t, handler), "gw-library-test-project", "us-central1-a", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "ACTIVE", attrs.State)
	assert.Equal(t, "projects/gw-library-test-project/global/networks/gw-library-test", attrs.Network)
}

func TestGetNetworkSecurityFirewallEndpointAssociationAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a firewall endpoint association that is not there should read a sentence about that firewall endpoint association, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetNetworkSecurityFirewallEndpointAssociationAttrsWithClient(context.Background(), newFakeNetworkSecurityService(t, handler), "gw-library-test-project", "us-central1-a", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetNetworkSecurityInterceptDeploymentAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a intercept deployment the terraform-google-networking intercept deployment module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1-a/interceptDeployments/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1-a/interceptDeployments/gw-library-test","forwardingRule":"projects/gw-library-test-project/regions/us-central1/forwardingRules/gw-library-test","state":"ACTIVE","labels":{"purpose":"terratest"}}`))
	})

	attrs, err := gcp.GetNetworkSecurityInterceptDeploymentAttrsWithClient(context.Background(), newFakeNetworkSecurityService(t, handler), "gw-library-test-project", "us-central1-a", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "ACTIVE", attrs.State)
}

func TestGetNetworkSecurityInterceptDeploymentAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a intercept deployment that is not there should read a sentence about that intercept deployment, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetNetworkSecurityInterceptDeploymentAttrsWithClient(context.Background(), newFakeNetworkSecurityService(t, handler), "gw-library-test-project", "us-central1-a", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetNetworkSecurityInterceptDeploymentGroupAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a intercept deployment group the terraform-google-networking intercept deployment group module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/global/interceptDeploymentGroups/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/global/interceptDeploymentGroups/gw-library-test","network":"projects/gw-library-test-project/global/networks/gw-library-test","state":"ACTIVE","labels":{"purpose":"terratest"}}`))
	})

	attrs, err := gcp.GetNetworkSecurityInterceptDeploymentGroupAttrsWithClient(context.Background(), newFakeNetworkSecurityService(t, handler), "gw-library-test-project", "global", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "ACTIVE", attrs.State)
}

func TestGetNetworkSecurityInterceptDeploymentGroupAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a intercept deployment group that is not there should read a sentence about that intercept deployment group, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetNetworkSecurityInterceptDeploymentGroupAttrsWithClient(context.Background(), newFakeNetworkSecurityService(t, handler), "gw-library-test-project", "global", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetNetworkSecurityInterceptEndpointGroupAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a intercept endpoint group the terraform-google-networking intercept endpoint group module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/global/interceptEndpointGroups/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/global/interceptEndpointGroups/gw-library-test","interceptDeploymentGroup":"projects/gw-library-test-project/locations/global/interceptDeploymentGroups/gw-library-test","state":"ACTIVE","labels":{"purpose":"terratest"}}`))
	})

	attrs, err := gcp.GetNetworkSecurityInterceptEndpointGroupAttrsWithClient(context.Background(), newFakeNetworkSecurityService(t, handler), "gw-library-test-project", "global", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "ACTIVE", attrs.State)
}

func TestGetNetworkSecurityInterceptEndpointGroupAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a intercept endpoint group that is not there should read a sentence about that intercept endpoint group, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetNetworkSecurityInterceptEndpointGroupAttrsWithClient(context.Background(), newFakeNetworkSecurityService(t, handler), "gw-library-test-project", "global", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetNetworkSecurityInterceptEndpointGroupAssociationAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a intercept endpoint group association the terraform-google-networking intercept endpoint group association module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/global/interceptEndpointGroupAssociations/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/global/interceptEndpointGroupAssociations/gw-library-test","network":"projects/gw-library-test-project/global/networks/gw-library-test","state":"ACTIVE","labels":{"purpose":"terratest"}}`))
	})

	attrs, err := gcp.GetNetworkSecurityInterceptEndpointGroupAssociationAttrsWithClient(context.Background(), newFakeNetworkSecurityService(t, handler), "gw-library-test-project", "global", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "ACTIVE", attrs.State)
}

func TestGetNetworkSecurityInterceptEndpointGroupAssociationAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a intercept endpoint group association that is not there should read a sentence about that intercept endpoint group association, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetNetworkSecurityInterceptEndpointGroupAssociationAttrsWithClient(context.Background(), newFakeNetworkSecurityService(t, handler), "gw-library-test-project", "global", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetNetworkSecurityMirroringDeploymentAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a mirroring deployment the terraform-google-networking mirroring deployment module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1-a/mirroringDeployments/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1-a/mirroringDeployments/gw-library-test","forwardingRule":"projects/gw-library-test-project/regions/us-central1/forwardingRules/gw-library-test","state":"ACTIVE","labels":{"purpose":"terratest"}}`))
	})

	attrs, err := gcp.GetNetworkSecurityMirroringDeploymentAttrsWithClient(context.Background(), newFakeNetworkSecurityService(t, handler), "gw-library-test-project", "us-central1-a", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "ACTIVE", attrs.State)
}

func TestGetNetworkSecurityMirroringDeploymentAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a mirroring deployment that is not there should read a sentence about that mirroring deployment, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetNetworkSecurityMirroringDeploymentAttrsWithClient(context.Background(), newFakeNetworkSecurityService(t, handler), "gw-library-test-project", "us-central1-a", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetNetworkSecurityMirroringDeploymentGroupAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a mirroring deployment group the terraform-google-networking mirroring deployment group module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/global/mirroringDeploymentGroups/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/global/mirroringDeploymentGroups/gw-library-test","network":"projects/gw-library-test-project/global/networks/gw-library-test","state":"ACTIVE","labels":{"purpose":"terratest"}}`))
	})

	attrs, err := gcp.GetNetworkSecurityMirroringDeploymentGroupAttrsWithClient(context.Background(), newFakeNetworkSecurityService(t, handler), "gw-library-test-project", "global", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "ACTIVE", attrs.State)
}

func TestGetNetworkSecurityMirroringDeploymentGroupAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a mirroring deployment group that is not there should read a sentence about that mirroring deployment group, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetNetworkSecurityMirroringDeploymentGroupAttrsWithClient(context.Background(), newFakeNetworkSecurityService(t, handler), "gw-library-test-project", "global", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetNetworkSecurityMirroringEndpointGroupAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a mirroring endpoint group the terraform-google-networking mirroring endpoint group module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/global/mirroringEndpointGroups/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/global/mirroringEndpointGroups/gw-library-test","mirroringDeploymentGroup":"projects/gw-library-test-project/locations/global/mirroringDeploymentGroups/gw-library-test","state":"ACTIVE","labels":{"purpose":"terratest"}}`))
	})

	attrs, err := gcp.GetNetworkSecurityMirroringEndpointGroupAttrsWithClient(context.Background(), newFakeNetworkSecurityService(t, handler), "gw-library-test-project", "global", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "ACTIVE", attrs.State)
}

func TestGetNetworkSecurityMirroringEndpointGroupAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a mirroring endpoint group that is not there should read a sentence about that mirroring endpoint group, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetNetworkSecurityMirroringEndpointGroupAttrsWithClient(context.Background(), newFakeNetworkSecurityService(t, handler), "gw-library-test-project", "global", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetNetworkSecurityMirroringEndpointGroupAssociationAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a mirroring endpoint group association the terraform-google-networking mirroring endpoint group association module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/global/mirroringEndpointGroupAssociations/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/global/mirroringEndpointGroupAssociations/gw-library-test","network":"projects/gw-library-test-project/global/networks/gw-library-test","state":"ACTIVE","labels":{"purpose":"terratest"}}`))
	})

	attrs, err := gcp.GetNetworkSecurityMirroringEndpointGroupAssociationAttrsWithClient(context.Background(), newFakeNetworkSecurityService(t, handler), "gw-library-test-project", "global", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "ACTIVE", attrs.State)
}

func TestGetNetworkSecurityMirroringEndpointGroupAssociationAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a mirroring endpoint group association that is not there should read a sentence about that mirroring endpoint group association, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetNetworkSecurityMirroringEndpointGroupAssociationAttrsWithClient(context.Background(), newFakeNetworkSecurityService(t, handler), "gw-library-test-project", "global", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetNetworkSecurityTLSInspectionPolicyAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a TLS inspection policy the terraform-google-networking TLS inspection policy module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/tlsInspectionPolicies/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/tlsInspectionPolicies/gw-library-test","caPool":"projects/gw-library-test-project/locations/us-central1/caPools/gw-library-test","excludePublicCaSet":true}`))
	})

	attrs, err := gcp.GetNetworkSecurityTLSInspectionPolicyAttrsWithClient(context.Background(), newFakeNetworkSecurityService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "projects/gw-library-test-project/locations/us-central1/caPools/gw-library-test", attrs.CaPool)
	assert.True(t, attrs.ExcludePublicCaSet)
}

func TestGetNetworkSecurityTLSInspectionPolicyAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a TLS inspection policy that is not there should read a sentence about that TLS inspection policy, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetNetworkSecurityTLSInspectionPolicyAttrsWithClient(context.Background(), newFakeNetworkSecurityService(t, handler), "gw-library-test-project", "us-central1", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}
