package gcp_test

import (
	"context"
	"net/http"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/gcp/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFetchNetworkFirewallPolicyWithClient(t *testing.T) {
	t.Parallel()

	// The values are the ones the terraform-google-networking network firewall policy module sets, because the point
	// of reading settings back is asserting a module configured the policy it was asked for.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/global/firewallPolicies/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"gw-library-test","description":"created by terratest","shortName":"gw-library-test","ruleTupleCount":2}`))
	})

	policy, err := gcp.FetchNetworkFirewallPolicyWithClient(context.Background(), newFakeComputeService(t, handler), "gw-library-test-project", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "gw-library-test", policy.Name)
	assert.Equal(t, "created by terratest", policy.Description)
	assert.Equal(t, "gw-library-test", policy.ShortName)
}

func TestFetchRegionNetworkFirewallPolicyWithClient(t *testing.T) {
	t.Parallel()

	// The values are the ones the terraform-google-networking region network firewall policy module sets, because the point
	// of reading settings back is asserting a module configured the policy it was asked for.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/regions/us-central1/firewallPolicies/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"gw-library-test","description":"created by terratest","shortName":"gw-library-test","region":"https://www.googleapis.com/compute/v1/projects/gw-library-test-project/regions/us-central1"}`))
	})

	policy, err := gcp.FetchRegionNetworkFirewallPolicyWithClient(context.Background(), newFakeComputeService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "gw-library-test", policy.Name)
	assert.Equal(t, "created by terratest", policy.Description)
	assert.True(t, strings.HasSuffix(policy.Region, "/regions/us-central1"))
}

func TestFetchNetworkFirewallPolicyIamPolicyWithClient(t *testing.T) {
	t.Parallel()

	// The values are the ones the terraform-google-networking network firewall policy IAM module sets, because the point
	// of reading settings back is asserting a module configured the binding it was asked for.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/global/firewallPolicies/gw-library-test/getIamPolicy"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"version":1,"etag":"BwXhqw==","bindings":[{"role":"roles/compute.networkViewer","members":["serviceAccount:gw-library-test@gw-library-test-project.iam.gserviceaccount.com"]}]}`))
	})

	policy, err := gcp.FetchNetworkFirewallPolicyIamPolicyWithClient(context.Background(), newFakeComputeService(t, handler), "gw-library-test-project", "gw-library-test")
	require.NoError(t, err)

	require.Len(t, policy.Bindings, 1, "the policy should carry exactly the binding the fixture gave it")
	assert.Equal(t, "roles/compute.networkViewer", policy.Bindings[0].Role)
	require.Len(t, policy.Bindings[0].Members, 1)
	assert.True(t, strings.HasPrefix(policy.Bindings[0].Members[0], "serviceAccount:gw-library-test@"))
}

func TestFetchRegionNetworkFirewallPolicyIamPolicyWithClient(t *testing.T) {
	t.Parallel()

	// The values are the ones the terraform-google-networking region network firewall policy IAM module sets, because the point
	// of reading settings back is asserting a module configured the binding it was asked for.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/regions/us-central1/firewallPolicies/gw-library-test/getIamPolicy"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"version":1,"etag":"BwXhqw==","bindings":[{"role":"roles/compute.networkViewer","members":["serviceAccount:gw-library-test@gw-library-test-project.iam.gserviceaccount.com"]}]}`))
	})

	policy, err := gcp.FetchRegionNetworkFirewallPolicyIamPolicyWithClient(context.Background(), newFakeComputeService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	require.Len(t, policy.Bindings, 1, "the policy should carry exactly the binding the fixture gave it")
	assert.Equal(t, "roles/compute.networkViewer", policy.Bindings[0].Role)
}
