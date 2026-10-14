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
	iamv2 "google.golang.org/api/iam/v2"
	"google.golang.org/api/option"
)

// newFakeIAMV2Service points a real IAM v2 client at a local test server, so the Google transport is
// exercised rather than a hand-written stand-in for a type we do not own.
func newFakeIAMV2Service(t *testing.T, handler http.Handler) *iamv2.Service {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	service, err := iamv2.NewService(context.Background(),
		option.WithEndpoint(server.URL), option.WithoutAuthentication())
	require.NoError(t, err)

	return service
}

func TestGetIAMDenyPolicyAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a deny policy a Gruntwork module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		// The percent encoded attachment point survives the round trip rather than being unescaped into
		// extra path segments, which is the whole reason this helper builds the name itself.
		assert.True(t, strings.HasSuffix(r.URL.Path, "/policies/cloudresourcemanager.googleapis.com%2Fprojects%2Fgw-library-test-project/denypolicies/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"policies/cloudresourcemanager.googleapis.com%2Fprojects%2Fgw-library-test-project/denypolicies/gw-library-test","displayName":"terratest deny policy","kind":"DenyPolicy","etag":"BwXhqw==","rules":[{"description":"no storage deletes","denyRule":{"deniedPrincipals":["principalSet://goog/public:all"],"deniedPermissions":["storage.googleapis.com/buckets.delete"],"exceptionPrincipals":["principal://iam.googleapis.com/projects/-/serviceAccounts/gw-library-test@gw-library-test-project.iam.gserviceaccount.com"]}}]}`))
	})

	result, err := gcp.GetIAMDenyPolicyAttrsWithClient(context.Background(), newFakeIAMV2Service(t, handler), "gw-library-test-project", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "terratest deny policy", result.DisplayName)
	require.Len(t, result.Rules, 1)
	require.NotNil(t, result.Rules[0].DenyRule)
	assert.Equal(t, []string{"storage.googleapis.com/buckets.delete"}, result.Rules[0].DenyRule.DeniedPermissions)
	assert.Equal(t, []string{"principalSet://goog/public:all"}, result.Rules[0].DenyRule.DeniedPrincipals)
}

func TestGetIAMDenyPolicyAttrsWithClientReportsAMissingPolicy(t *testing.T) {
	t.Parallel()

	// A caller who names a deny policy that is not there should read a sentence about that policy, not
	// a status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetIAMDenyPolicyAttrsWithClient(context.Background(), newFakeIAMV2Service(t, handler), "gw-library-test-project", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}
