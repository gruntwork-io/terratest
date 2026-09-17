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

func TestFetchSecurityPolicyWithClient(t *testing.T) {
	t.Parallel()

	// The values are the ones the terraform-google-security policy module sets. A rule whose action
	// or source range came out wrong is the case most worth catching here.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/global/securityPolicies/gw-library-test"))
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"kind":"compute#securityPolicy",
			"name":"gw-library-test",
			"description":"created by terratest",
			"type":"CLOUD_ARMOR",
			"rules":[
				{"action":"deny(403)","priority":1000,"preview":true,"description":"block a private range",
				 "match":{"versionedExpr":"SRC_IPS_V1","config":{"srcIpRanges":["10.0.0.0/8"]}}},
				{"action":"allow","priority":2147483647,"description":"default rule",
				 "match":{"versionedExpr":"SRC_IPS_V1","config":{"srcIpRanges":["*"]}}}
			]
		}`))
	})

	policy, err := gcp.FetchSecurityPolicyWithClient(context.Background(), newFakeComputeService(t, handler), "gw-library-test-project", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "gw-library-test", policy.Name)
	assert.Equal(t, "created by terratest", policy.Description)
	assert.Equal(t, "CLOUD_ARMOR", policy.Type)
	require.Len(t, policy.Rules, 2)
	assert.Equal(t, "deny(403)", policy.Rules[0].Action)
	assert.Equal(t, int64(1000), policy.Rules[0].Priority)
	assert.True(t, policy.Rules[0].Preview)
	require.NotNil(t, policy.Rules[0].Match)
	require.NotNil(t, policy.Rules[0].Match.Config)
	assert.Equal(t, []string{"10.0.0.0/8"}, policy.Rules[0].Match.Config.SrcIpRanges)
}

func TestFetchSecurityPolicyWithClientMissingPolicy(t *testing.T) {
	t.Parallel()

	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	// The error names the call, the project and the policy, in the shape the other compute reads
	// here use, so all three are asserted rather than only that it failed.
	_, err := gcp.FetchSecurityPolicyWithClient(context.Background(), newFakeComputeService(t, handler), "gw-library-test-project", "gone")
	require.ErrorContains(t, err, "SecurityPolicies.Get(gw-library-test-project, gone)")
}
