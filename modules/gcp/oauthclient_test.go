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

func TestGetOAuthClientAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The values are the ones the terraform-google-identity OAuth client module sets, because the
	// point of reading settings back is asserting a module configured the client it was asked for.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/global/oauthClients/gw-library-test"))
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"name":"projects/gw-library-test-project/locations/global/oauthClients/gw-library-test",
			"clientId":"gw-library-test",
			"displayName":"terratest client",
			"description":"created by terratest",
			"clientType":"CONFIDENTIAL_CLIENT",
			"allowedGrantTypes":["AUTHORIZATION_CODE_GRANT"],
			"allowedRedirectUris":["https://example.com/callback"],
			"allowedScopes":["openid"],
			"disabled":true,
			"state":"ACTIVE"
		}`))
	})

	client, err := gcp.GetOAuthClientAttrsWithClient(context.Background(), newFakeIAMService(t, handler), "gw-library-test-project", "global", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "terratest client", client.DisplayName)
	assert.Equal(t, "created by terratest", client.Description)
	assert.Equal(t, "CONFIDENTIAL_CLIENT", client.ClientType)
	assert.Equal(t, []string{"AUTHORIZATION_CODE_GRANT"}, client.AllowedGrantTypes)
	assert.Equal(t, []string{"https://example.com/callback"}, client.AllowedRedirectUris)
	assert.Equal(t, []string{"openid"}, client.AllowedScopes)
	assert.True(t, client.Disabled)
}

func TestGetOAuthClientAttrsWithClientMissingClient(t *testing.T) {
	t.Parallel()

	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	// The error names the client and the project as well as saying it is absent, so all three are
	// asserted rather than only the phrase.
	_, err := gcp.GetOAuthClientAttrsWithClient(context.Background(), newFakeIAMService(t, handler), "gw-library-test-project", "global", "gone")
	require.ErrorContains(t, err, "does not exist")
	require.ErrorContains(t, err, "gone")
	require.ErrorContains(t, err, "gw-library-test-project")
}
