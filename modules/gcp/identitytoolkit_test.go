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
	identitytoolkit "google.golang.org/api/identitytoolkit/v2"
	"google.golang.org/api/option"
)

// newFakeIdentityPlatformService points a real Identity Toolkit client at an httptest server, so a read
// can be exercised against a response we control without reaching Google.
func newFakeIdentityPlatformService(t *testing.T, handler http.Handler) *identitytoolkit.Service {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	service, err := identitytoolkit.NewService(context.Background(),
		option.WithoutAuthentication(), option.WithEndpoint(server.URL), option.WithHTTPClient(server.Client()))
	require.NoError(t, err)

	return service
}

func TestGetIdentityPlatformConfigAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a Identity Platform config the terraform-google-identity Identity Platform config module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/config"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/config","signIn":{"allowDuplicateEmails":false,"anonymous":{"enabled":true}},"autodeleteAnonymousUsers":true}`))
	})

	attrs, err := gcp.GetIdentityPlatformConfigAttrsWithClient(context.Background(), newFakeIdentityPlatformService(t, handler), "gw-library-test-project")
	require.NoError(t, err)

	assert.True(t, attrs.AutodeleteAnonymousUsers)
}

func TestGetIdentityPlatformConfigAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a Identity Platform config that is not there should read a sentence about that Identity Platform config, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetIdentityPlatformConfigAttrsWithClient(context.Background(), newFakeIdentityPlatformService(t, handler), "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetIdentityPlatformTenantAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a Identity Platform tenant the terraform-google-identity Identity Platform tenant module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/tenants/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/tenants/gw-library-test","displayName":"terratest tenant","allowPasswordSignup":true,"enableEmailLinkSignin":true,"disableAuth":false}`))
	})

	attrs, err := gcp.GetIdentityPlatformTenantAttrsWithClient(context.Background(), newFakeIdentityPlatformService(t, handler), "gw-library-test-project", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "terratest tenant", attrs.DisplayName)
	assert.True(t, attrs.AllowPasswordSignup)
	assert.True(t, attrs.EnableEmailLinkSignin)
}

func TestGetIdentityPlatformTenantAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a Identity Platform tenant that is not there should read a sentence about that Identity Platform tenant, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetIdentityPlatformTenantAttrsWithClient(context.Background(), newFakeIdentityPlatformService(t, handler), "gw-library-test-project", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetIdentityPlatformDefaultSupportedIdpConfigAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a Identity Platform default IdP config the terraform-google-identity Identity Platform default IdP config module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/defaultSupportedIdpConfigs/facebook.com"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/defaultSupportedIdpConfigs/facebook.com","enabled":true,"clientId":"terratest-client"}`))
	})

	attrs, err := gcp.GetIdentityPlatformDefaultSupportedIdpConfigAttrsWithClient(context.Background(), newFakeIdentityPlatformService(t, handler), "gw-library-test-project", "facebook.com")
	require.NoError(t, err)

	assert.True(t, attrs.Enabled)
	assert.Equal(t, "terratest-client", attrs.ClientId)
}

func TestGetIdentityPlatformDefaultSupportedIdpConfigAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a Identity Platform default IdP config that is not there should read a sentence about that Identity Platform default IdP config, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetIdentityPlatformDefaultSupportedIdpConfigAttrsWithClient(context.Background(), newFakeIdentityPlatformService(t, handler), "gw-library-test-project", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetIdentityPlatformInboundSamlConfigAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a Identity Platform inbound SAML config the terraform-google-identity Identity Platform inbound SAML config module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/inboundSamlConfigs/saml.gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/inboundSamlConfigs/saml.gw-library-test","displayName":"terratest saml","enabled":true,"idpConfig":{"idpEntityId":"https://terratest.example.com/saml","ssoUrl":"https://terratest.example.com/sso"}}`))
	})

	attrs, err := gcp.GetIdentityPlatformInboundSamlConfigAttrsWithClient(context.Background(), newFakeIdentityPlatformService(t, handler), "gw-library-test-project", "saml.gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "terratest saml", attrs.DisplayName)
	assert.True(t, attrs.Enabled)
	assert.Equal(t, "https://terratest.example.com/saml", attrs.IdpConfig.IdpEntityId)
}

func TestGetIdentityPlatformInboundSamlConfigAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a Identity Platform inbound SAML config that is not there should read a sentence about that Identity Platform inbound SAML config, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetIdentityPlatformInboundSamlConfigAttrsWithClient(context.Background(), newFakeIdentityPlatformService(t, handler), "gw-library-test-project", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetIdentityPlatformOAuthIdpConfigAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a Identity Platform OAuth IdP config the terraform-google-identity Identity Platform OAuth IdP config module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/oauthIdpConfigs/oidc.gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/oauthIdpConfigs/oidc.gw-library-test","displayName":"terratest oidc","enabled":true,"issuer":"https://terratest.example.com","clientId":"terratest-client"}`))
	})

	attrs, err := gcp.GetIdentityPlatformOAuthIdpConfigAttrsWithClient(context.Background(), newFakeIdentityPlatformService(t, handler), "gw-library-test-project", "oidc.gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "terratest oidc", attrs.DisplayName)
	assert.Equal(t, "https://terratest.example.com", attrs.Issuer)
	assert.Equal(t, "terratest-client", attrs.ClientId)
}

func TestGetIdentityPlatformOAuthIdpConfigAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a Identity Platform OAuth IdP config that is not there should read a sentence about that Identity Platform OAuth IdP config, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetIdentityPlatformOAuthIdpConfigAttrsWithClient(context.Background(), newFakeIdentityPlatformService(t, handler), "gw-library-test-project", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetIdentityPlatformTenantDefaultSupportedIdpConfigAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a Identity Platform tenant default IdP config the terraform-google-identity Identity Platform tenant default IdP config module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/tenants/gw-library-parent/defaultSupportedIdpConfigs/facebook.com"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/tenants/gw-library-parent/defaultSupportedIdpConfigs/facebook.com","enabled":true,"clientId":"terratest-client"}`))
	})

	attrs, err := gcp.GetIdentityPlatformTenantDefaultSupportedIdpConfigAttrsWithClient(context.Background(), newFakeIdentityPlatformService(t, handler), "gw-library-test-project", "gw-library-parent", "facebook.com")
	require.NoError(t, err)

	assert.True(t, attrs.Enabled)
	assert.Equal(t, "terratest-client", attrs.ClientId)
}

func TestGetIdentityPlatformTenantDefaultSupportedIdpConfigAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a Identity Platform tenant default IdP config that is not there should read a sentence about that Identity Platform tenant default IdP config, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetIdentityPlatformTenantDefaultSupportedIdpConfigAttrsWithClient(context.Background(), newFakeIdentityPlatformService(t, handler), "gw-library-test-project", "gw-library-parent", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetIdentityPlatformTenantInboundSamlConfigAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a Identity Platform tenant inbound SAML config the terraform-google-identity Identity Platform tenant inbound SAML config module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/tenants/gw-library-parent/inboundSamlConfigs/saml.gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/tenants/gw-library-parent/inboundSamlConfigs/saml.gw-library-test","displayName":"terratest saml","enabled":true,"idpConfig":{"idpEntityId":"https://terratest.example.com/saml","ssoUrl":"https://terratest.example.com/sso"}}`))
	})

	attrs, err := gcp.GetIdentityPlatformTenantInboundSamlConfigAttrsWithClient(context.Background(), newFakeIdentityPlatformService(t, handler), "gw-library-test-project", "gw-library-parent", "saml.gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "terratest saml", attrs.DisplayName)
	assert.True(t, attrs.Enabled)
}

func TestGetIdentityPlatformTenantInboundSamlConfigAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a Identity Platform tenant inbound SAML config that is not there should read a sentence about that Identity Platform tenant inbound SAML config, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetIdentityPlatformTenantInboundSamlConfigAttrsWithClient(context.Background(), newFakeIdentityPlatformService(t, handler), "gw-library-test-project", "gw-library-parent", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetIdentityPlatformTenantOAuthIdpConfigAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a Identity Platform tenant OAuth IdP config the terraform-google-identity Identity Platform tenant OAuth IdP config module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/tenants/gw-library-parent/oauthIdpConfigs/oidc.gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/tenants/gw-library-parent/oauthIdpConfigs/oidc.gw-library-test","displayName":"terratest oidc","enabled":true,"issuer":"https://terratest.example.com","clientId":"terratest-client"}`))
	})

	attrs, err := gcp.GetIdentityPlatformTenantOAuthIdpConfigAttrsWithClient(context.Background(), newFakeIdentityPlatformService(t, handler), "gw-library-test-project", "gw-library-parent", "oidc.gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "terratest oidc", attrs.DisplayName)
	assert.Equal(t, "https://terratest.example.com", attrs.Issuer)
}

func TestGetIdentityPlatformTenantOAuthIdpConfigAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a Identity Platform tenant OAuth IdP config that is not there should read a sentence about that Identity Platform tenant OAuth IdP config, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetIdentityPlatformTenantOAuthIdpConfigAttrsWithClient(context.Background(), newFakeIdentityPlatformService(t, handler), "gw-library-test-project", "gw-library-parent", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}
