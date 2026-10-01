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
	"google.golang.org/api/developerconnect/v1"
	"google.golang.org/api/option"
)

// newFakeDeveloperConnectService points a real Developer Connect client at an httptest server, so a read can be exercised
// against a response we control without reaching Google.
func newFakeDeveloperConnectService(t *testing.T, handler http.Handler) *developerconnect.Service {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	service, err := developerconnect.NewService(context.Background(),
		option.WithoutAuthentication(), option.WithEndpoint(server.URL), option.WithHTTPClient(server.Client()))
	require.NoError(t, err)

	return service
}

func TestGetDeveloperConnectConnectionAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a Developer Connect connection the terraform-google-devtools Developer Connect connection module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/connections/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/connections/gw-library-test","disabled":false,"installationState":{"stage":"COMPLETE"},"labels":{"purpose":"terratest"}}`))
	})

	attrs, err := gcp.GetDeveloperConnectConnectionAttrsWithClient(context.Background(), newFakeDeveloperConnectService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "COMPLETE", attrs.InstallationState.Stage)
}

func TestGetDeveloperConnectConnectionAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a Developer Connect connection that is not there should read a sentence about that Developer Connect connection, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetDeveloperConnectConnectionAttrsWithClient(context.Background(), newFakeDeveloperConnectService(t, handler), "gw-library-test-project", "us-central1", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetDeveloperConnectGitRepositoryLinkAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a Developer Connect git repository link the terraform-google-devtools Developer Connect git repository link module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/connections/gw-library-parent/gitRepositoryLinks/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/connections/gw-library-parent/gitRepositoryLinks/gw-library-test","cloneUri":"https://github.com/terratest/terratest.git","labels":{"purpose":"terratest"}}`))
	})

	attrs, err := gcp.GetDeveloperConnectGitRepositoryLinkAttrsWithClient(context.Background(), newFakeDeveloperConnectService(t, handler), "gw-library-test-project", "us-central1", "gw-library-parent", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "https://github.com/terratest/terratest.git", attrs.CloneUri)
}

func TestGetDeveloperConnectGitRepositoryLinkAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a Developer Connect git repository link that is not there should read a sentence about that Developer Connect git repository link, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetDeveloperConnectGitRepositoryLinkAttrsWithClient(context.Background(), newFakeDeveloperConnectService(t, handler), "gw-library-test-project", "us-central1", "gw-library-parent", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetDeveloperConnectAccountConnectorAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a Developer Connect account connector the terraform-google-devtools Developer Connect account connector module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/accountConnectors/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/accountConnectors/gw-library-test","labels":{"purpose":"terratest"},"providerOauthConfig":{"systemProviderId":"GITHUB"}}`))
	})

	attrs, err := gcp.GetDeveloperConnectAccountConnectorAttrsWithClient(context.Background(), newFakeDeveloperConnectService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "GITHUB", attrs.ProviderOauthConfig.SystemProviderId)
}

func TestGetDeveloperConnectAccountConnectorAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a Developer Connect account connector that is not there should read a sentence about that Developer Connect account connector, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetDeveloperConnectAccountConnectorAttrsWithClient(context.Background(), newFakeDeveloperConnectService(t, handler), "gw-library-test-project", "us-central1", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetDeveloperConnectInsightsConfigAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a Developer Connect insights config the terraform-google-devtools Developer Connect insights config module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/insightsConfigs/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/insightsConfigs/gw-library-test","appHubApplication":"projects/gw-library-test-project/locations/us-central1/applications/gw-library-parent","state":"ACTIVE","labels":{"purpose":"terratest"}}`))
	})

	attrs, err := gcp.GetDeveloperConnectInsightsConfigAttrsWithClient(context.Background(), newFakeDeveloperConnectService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "ACTIVE", attrs.State)
}

func TestGetDeveloperConnectInsightsConfigAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a Developer Connect insights config that is not there should read a sentence about that Developer Connect insights config, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetDeveloperConnectInsightsConfigAttrsWithClient(context.Background(), newFakeDeveloperConnectService(t, handler), "gw-library-test-project", "us-central1", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}
