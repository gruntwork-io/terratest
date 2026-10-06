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
	"google.golang.org/api/vmwareengine/v1"
)

// newFakeVMwareEngineService points a real VMware Engine client at an httptest server, so a read can be exercised
// against a response we control without reaching Google.
func newFakeVMwareEngineService(t *testing.T, handler http.Handler) *vmwareengine.Service {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	service, err := vmwareengine.NewService(context.Background(),
		option.WithoutAuthentication(), option.WithEndpoint(server.URL), option.WithHTTPClient(server.Client()))
	require.NoError(t, err)

	return service
}

func TestGetVMwareEnginePrivateCloudAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a VMware Engine private cloud the terraform-google-compute VMware Engine private cloud module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/privateClouds/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/privateClouds/gw-library-test","description":"created by terratest","type":"TIME_LIMITED","state":"ACTIVE"}`))
	})

	attrs, err := gcp.GetVMwareEnginePrivateCloudAttrsWithClient(context.Background(), newFakeVMwareEngineService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "created by terratest", attrs.Description)
	assert.Equal(t, "TIME_LIMITED", attrs.Type)
	assert.Equal(t, "ACTIVE", attrs.State)
}

func TestGetVMwareEnginePrivateCloudAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a VMware Engine private cloud that is not there should read a sentence about that VMware Engine private cloud, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetVMwareEnginePrivateCloudAttrsWithClient(context.Background(), newFakeVMwareEngineService(t, handler), "gw-library-test-project", "us-central1", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetVMwareEngineClusterAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a VMware Engine cluster the terraform-google-compute VMware Engine cluster module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/privateClouds/gw-library-parent/clusters/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/privateClouds/gw-library-parent/clusters/gw-library-test","state":"ACTIVE","management":false}`))
	})

	attrs, err := gcp.GetVMwareEngineClusterAttrsWithClient(context.Background(), newFakeVMwareEngineService(t, handler), "gw-library-test-project", "us-central1", "gw-library-parent", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "ACTIVE", attrs.State)
}

func TestGetVMwareEngineClusterAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a VMware Engine cluster that is not there should read a sentence about that VMware Engine cluster, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetVMwareEngineClusterAttrsWithClient(context.Background(), newFakeVMwareEngineService(t, handler), "gw-library-test-project", "us-central1", "gw-library-parent", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetVMwareEngineSubnetAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a VMware Engine subnet the terraform-google-compute VMware Engine subnet module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/privateClouds/gw-library-parent/subnets/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/privateClouds/gw-library-parent/subnets/gw-library-test","ipCidrRange":"10.98.0.0/24","type":"USER_DEFINED","state":"ACTIVE"}`))
	})

	attrs, err := gcp.GetVMwareEngineSubnetAttrsWithClient(context.Background(), newFakeVMwareEngineService(t, handler), "gw-library-test-project", "us-central1", "gw-library-parent", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "10.98.0.0/24", attrs.IpCidrRange)
	assert.Equal(t, "USER_DEFINED", attrs.Type)
}

func TestGetVMwareEngineSubnetAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a VMware Engine subnet that is not there should read a sentence about that VMware Engine subnet, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetVMwareEngineSubnetAttrsWithClient(context.Background(), newFakeVMwareEngineService(t, handler), "gw-library-test-project", "us-central1", "gw-library-parent", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetVMwareEngineExternalAddressAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a VMware Engine external address the terraform-google-compute VMware Engine external address module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/privateClouds/gw-library-parent/externalAddresses/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/privateClouds/gw-library-parent/externalAddresses/gw-library-test","internalIp":"192.168.0.10","description":"created by terratest","state":"ACTIVE"}`))
	})

	attrs, err := gcp.GetVMwareEngineExternalAddressAttrsWithClient(context.Background(), newFakeVMwareEngineService(t, handler), "gw-library-test-project", "us-central1", "gw-library-parent", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "192.168.0.10", attrs.InternalIp)
	assert.Equal(t, "created by terratest", attrs.Description)
}

func TestGetVMwareEngineExternalAddressAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a VMware Engine external address that is not there should read a sentence about that VMware Engine external address, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetVMwareEngineExternalAddressAttrsWithClient(context.Background(), newFakeVMwareEngineService(t, handler), "gw-library-test-project", "us-central1", "gw-library-parent", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetVMwareEngineNetworkAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a VMware Engine network the terraform-google-compute VMware Engine network module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/vmwareEngineNetworks/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/vmwareEngineNetworks/gw-library-test","description":"created by terratest","type":"STANDARD","state":"ACTIVE"}`))
	})

	attrs, err := gcp.GetVMwareEngineNetworkAttrsWithClient(context.Background(), newFakeVMwareEngineService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "created by terratest", attrs.Description)
	assert.Equal(t, "STANDARD", attrs.Type)
}

func TestGetVMwareEngineNetworkAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a VMware Engine network that is not there should read a sentence about that VMware Engine network, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetVMwareEngineNetworkAttrsWithClient(context.Background(), newFakeVMwareEngineService(t, handler), "gw-library-test-project", "us-central1", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetVMwareEngineNetworkPeeringAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a VMware Engine network peering the terraform-google-compute VMware Engine network peering module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/networkPeerings/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/networkPeerings/gw-library-test","description":"created by terratest","peerNetworkType":"STANDARD","state":"ACTIVE"}`))
	})

	attrs, err := gcp.GetVMwareEngineNetworkPeeringAttrsWithClient(context.Background(), newFakeVMwareEngineService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "created by terratest", attrs.Description)
	assert.Equal(t, "STANDARD", attrs.PeerNetworkType)
}

func TestGetVMwareEngineNetworkPeeringAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a VMware Engine network peering that is not there should read a sentence about that VMware Engine network peering, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetVMwareEngineNetworkPeeringAttrsWithClient(context.Background(), newFakeVMwareEngineService(t, handler), "gw-library-test-project", "us-central1", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetVMwareEngineNetworkPolicyAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a VMware Engine network policy the terraform-google-compute VMware Engine network policy module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/networkPolicies/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/networkPolicies/gw-library-test","description":"created by terratest","edgeServicesCidr":"192.168.30.0/26","internetAccess":{"enabled":true}}`))
	})

	attrs, err := gcp.GetVMwareEngineNetworkPolicyAttrsWithClient(context.Background(), newFakeVMwareEngineService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "created by terratest", attrs.Description)
	assert.Equal(t, "192.168.30.0/26", attrs.EdgeServicesCidr)
}

func TestGetVMwareEngineNetworkPolicyAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a VMware Engine network policy that is not there should read a sentence about that VMware Engine network policy, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetVMwareEngineNetworkPolicyAttrsWithClient(context.Background(), newFakeVMwareEngineService(t, handler), "gw-library-test-project", "us-central1", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetVMwareEngineExternalAccessRuleAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a VMware Engine external access rule the terraform-google-compute VMware Engine external access rule module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/networkPolicies/gw-library-parent/externalAccessRules/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/networkPolicies/gw-library-parent/externalAccessRules/gw-library-test","description":"created by terratest","priority":1000,"action":"ALLOW","ipProtocol":"TCP","state":"ACTIVE"}`))
	})

	attrs, err := gcp.GetVMwareEngineExternalAccessRuleAttrsWithClient(context.Background(), newFakeVMwareEngineService(t, handler), "gw-library-test-project", "us-central1", "gw-library-parent", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "created by terratest", attrs.Description)
	assert.Equal(t, int64(1000), attrs.Priority)
	assert.Equal(t, "ALLOW", attrs.Action)
}

func TestGetVMwareEngineExternalAccessRuleAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a VMware Engine external access rule that is not there should read a sentence about that VMware Engine external access rule, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetVMwareEngineExternalAccessRuleAttrsWithClient(context.Background(), newFakeVMwareEngineService(t, handler), "gw-library-test-project", "us-central1", "gw-library-parent", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetVMwareEngineDatastoreAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a VMware Engine datastore the terraform-google-compute VMware Engine datastore module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/datastores/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/datastores/gw-library-test","state":"ACTIVE"}`))
	})

	attrs, err := gcp.GetVMwareEngineDatastoreAttrsWithClient(context.Background(), newFakeVMwareEngineService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "ACTIVE", attrs.State)
}

func TestGetVMwareEngineDatastoreAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a VMware Engine datastore that is not there should read a sentence about that VMware Engine datastore, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetVMwareEngineDatastoreAttrsWithClient(context.Background(), newFakeVMwareEngineService(t, handler), "gw-library-test-project", "us-central1", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}
