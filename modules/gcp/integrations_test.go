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
	"google.golang.org/api/integrations/v1"
	"google.golang.org/api/option"
)

// newFakeIntegrationsService points a real Application Integration client at a local test server, so
// the Google transport is exercised rather than a hand-written stand-in for a type we do not own.
func newFakeIntegrationsService(t *testing.T, handler http.Handler) *integrations.Service {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	service, err := integrations.NewService(context.Background(),
		option.WithEndpoint(server.URL), option.WithoutAuthentication())
	require.NoError(t, err)

	return service
}

func TestGetIntegrationsAuthConfigAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for an auth config the auth config module
	// created, not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/authConfigs/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/authConfigs/gw-library-test","displayName":"terratest auth config","description":"created by terratest","credentialType":"USERNAME_AND_PASSWORD","state":"VALID","visibility":"PRIVATE"}`))
	})

	config, err := gcp.GetIntegrationsAuthConfigAttrsWithClient(context.Background(), newFakeIntegrationsService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "terratest auth config", config.DisplayName)
	assert.Equal(t, "USERNAME_AND_PASSWORD", config.CredentialType)
	assert.Equal(t, "VALID", config.State)
}

func TestGetIntegrationsAuthConfigAttrsWithClientReportsAMissingConfig(t *testing.T) {
	t.Parallel()

	// A caller who asks for an auth config that is not there should be told that, rather than be
	// handed the transport's own wording for a 404.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error":{"code":404,"message":"Not found."}}`))
	})

	_, err := gcp.GetIntegrationsAuthConfigAttrsWithClient(context.Background(), newFakeIntegrationsService(t, handler), "gw-library-test-project", "us-central1", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "does not exist")
}
