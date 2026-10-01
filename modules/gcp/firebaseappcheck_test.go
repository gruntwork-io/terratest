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
	"google.golang.org/api/firebaseappcheck/v1"
	"google.golang.org/api/option"
)

// newFakeFirebaseAppCheckService points a real Firebase App Check client at an httptest server, so a read can be exercised
// against a response we control without reaching Google.
func newFakeFirebaseAppCheckService(t *testing.T, handler http.Handler) *firebaseappcheck.Service {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	service, err := firebaseappcheck.NewService(context.Background(),
		option.WithoutAuthentication(), option.WithEndpoint(server.URL), option.WithHTTPClient(server.Client()))
	require.NoError(t, err)

	return service
}

func TestGetFirebaseAppCheckAppAttestConfigAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a Firebase App Check App Attest config the terraform-google-firebase Firebase App Check App Attest config module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/apps/gw-library-app/appAttestConfig"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/apps/gw-library-app/appAttestConfig","tokenTtl":"3600s"}`))
	})

	attrs, err := gcp.GetFirebaseAppCheckAppAttestConfigAttrsWithClient(context.Background(), newFakeFirebaseAppCheckService(t, handler), "gw-library-test-project", "gw-library-app")
	require.NoError(t, err)

	assert.Equal(t, "3600s", attrs.TokenTtl)
}

func TestGetFirebaseAppCheckAppAttestConfigAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a Firebase App Check App Attest config that is not there should read a sentence about that Firebase App Check App Attest config, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetFirebaseAppCheckAppAttestConfigAttrsWithClient(context.Background(), newFakeFirebaseAppCheckService(t, handler), "gw-library-test-project", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetFirebaseAppCheckDeviceCheckConfigAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a Firebase App Check DeviceCheck config the terraform-google-firebase Firebase App Check DeviceCheck config module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/apps/gw-library-app/deviceCheckConfig"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/apps/gw-library-app/deviceCheckConfig","tokenTtl":"3600s","keyId":"terratest-key","privateKeySet":true}`))
	})

	attrs, err := gcp.GetFirebaseAppCheckDeviceCheckConfigAttrsWithClient(context.Background(), newFakeFirebaseAppCheckService(t, handler), "gw-library-test-project", "gw-library-app")
	require.NoError(t, err)

	assert.Equal(t, "3600s", attrs.TokenTtl)
	assert.Equal(t, "terratest-key", attrs.KeyId)
}

func TestGetFirebaseAppCheckDeviceCheckConfigAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a Firebase App Check DeviceCheck config that is not there should read a sentence about that Firebase App Check DeviceCheck config, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetFirebaseAppCheckDeviceCheckConfigAttrsWithClient(context.Background(), newFakeFirebaseAppCheckService(t, handler), "gw-library-test-project", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetFirebaseAppCheckPlayIntegrityConfigAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a Firebase App Check Play Integrity config the terraform-google-firebase Firebase App Check Play Integrity config module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/apps/gw-library-app/playIntegrityConfig"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/apps/gw-library-app/playIntegrityConfig","tokenTtl":"3600s"}`))
	})

	attrs, err := gcp.GetFirebaseAppCheckPlayIntegrityConfigAttrsWithClient(context.Background(), newFakeFirebaseAppCheckService(t, handler), "gw-library-test-project", "gw-library-app")
	require.NoError(t, err)

	assert.Equal(t, "3600s", attrs.TokenTtl)
}

func TestGetFirebaseAppCheckPlayIntegrityConfigAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a Firebase App Check Play Integrity config that is not there should read a sentence about that Firebase App Check Play Integrity config, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetFirebaseAppCheckPlayIntegrityConfigAttrsWithClient(context.Background(), newFakeFirebaseAppCheckService(t, handler), "gw-library-test-project", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetFirebaseAppCheckRecaptchaEnterpriseConfigAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a Firebase App Check reCAPTCHA Enterprise config the terraform-google-firebase Firebase App Check reCAPTCHA Enterprise config module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/apps/gw-library-app/recaptchaEnterpriseConfig"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/apps/gw-library-app/recaptchaEnterpriseConfig","tokenTtl":"3600s","siteKey":"terratest-site-key"}`))
	})

	attrs, err := gcp.GetFirebaseAppCheckRecaptchaEnterpriseConfigAttrsWithClient(context.Background(), newFakeFirebaseAppCheckService(t, handler), "gw-library-test-project", "gw-library-app")
	require.NoError(t, err)

	assert.Equal(t, "3600s", attrs.TokenTtl)
	assert.Equal(t, "terratest-site-key", attrs.SiteKey)
}

func TestGetFirebaseAppCheckRecaptchaEnterpriseConfigAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a Firebase App Check reCAPTCHA Enterprise config that is not there should read a sentence about that Firebase App Check reCAPTCHA Enterprise config, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetFirebaseAppCheckRecaptchaEnterpriseConfigAttrsWithClient(context.Background(), newFakeFirebaseAppCheckService(t, handler), "gw-library-test-project", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetFirebaseAppCheckRecaptchaV3ConfigAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a Firebase App Check reCAPTCHA v3 config the terraform-google-firebase Firebase App Check reCAPTCHA v3 config module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/apps/gw-library-app/recaptchaV3Config"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/apps/gw-library-app/recaptchaV3Config","tokenTtl":"3600s","siteSecretSet":true}`))
	})

	attrs, err := gcp.GetFirebaseAppCheckRecaptchaV3ConfigAttrsWithClient(context.Background(), newFakeFirebaseAppCheckService(t, handler), "gw-library-test-project", "gw-library-app")
	require.NoError(t, err)

	assert.Equal(t, "3600s", attrs.TokenTtl)
	assert.True(t, attrs.SiteSecretSet)
}

func TestGetFirebaseAppCheckRecaptchaV3ConfigAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a Firebase App Check reCAPTCHA v3 config that is not there should read a sentence about that Firebase App Check reCAPTCHA v3 config, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetFirebaseAppCheckRecaptchaV3ConfigAttrsWithClient(context.Background(), newFakeFirebaseAppCheckService(t, handler), "gw-library-test-project", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetFirebaseAppCheckDebugTokenAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a Firebase App Check debug token the terraform-google-firebase Firebase App Check debug token module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/apps/gw-library-app/debugTokens/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/apps/gw-library-app/debugTokens/gw-library-test","displayName":"terratest token"}`))
	})

	attrs, err := gcp.GetFirebaseAppCheckDebugTokenAttrsWithClient(context.Background(), newFakeFirebaseAppCheckService(t, handler), "gw-library-test-project", "gw-library-app")
	require.NoError(t, err)

	assert.Equal(t, "terratest token", attrs.DisplayName)
}

func TestGetFirebaseAppCheckDebugTokenAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a Firebase App Check debug token that is not there should read a sentence about that Firebase App Check debug token, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetFirebaseAppCheckDebugTokenAttrsWithClient(context.Background(), newFakeFirebaseAppCheckService(t, handler), "gw-library-test-project", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetFirebaseAppCheckServiceConfigAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a Firebase App Check service config the terraform-google-firebase App Check service config module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/services/firestore.googleapis.com"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/services/firestore.googleapis.com","enforcementMode":"ENFORCED"}`))
	})

	attrs, err := gcp.GetFirebaseAppCheckServiceConfigAttrsWithClient(context.Background(), newFakeFirebaseAppCheckService(t, handler), "gw-library-test-project", "firestore.googleapis.com")
	require.NoError(t, err)

	assert.Equal(t, "ENFORCED", attrs.EnforcementMode)
}

func TestGetFirebaseAppCheckServiceConfigAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a Firebase App Check service config that is not there should read a sentence about that Firebase App Check service config, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetFirebaseAppCheckServiceConfigAttrsWithClient(context.Background(), newFakeFirebaseAppCheckService(t, handler), "gw-library-test-project", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetFirebaseAppCheckResourcePolicyAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a Firebase App Check resource policy the terraform-google-firebase App Check resource policy module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/services/firestore.googleapis.com/resourcePolicies/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/services/firestore.googleapis.com/resourcePolicies/gw-library-test","targetResource":"//firestore.googleapis.com/projects/gw-library-test-project/databases/gw-library-db","enforcementMode":"UNENFORCED"}`))
	})

	attrs, err := gcp.GetFirebaseAppCheckResourcePolicyAttrsWithClient(context.Background(), newFakeFirebaseAppCheckService(t, handler), "gw-library-test-project", "firestore.googleapis.com", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "UNENFORCED", attrs.EnforcementMode)
	assert.Equal(t, "//firestore.googleapis.com/projects/gw-library-test-project/databases/gw-library-db", attrs.TargetResource)
}

func TestGetFirebaseAppCheckResourcePolicyAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a Firebase App Check resource policy that is not there should read a sentence about that Firebase App Check resource policy, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetFirebaseAppCheckResourcePolicyAttrsWithClient(context.Background(), newFakeFirebaseAppCheckService(t, handler), "gw-library-test-project", "firestore.googleapis.com", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}
