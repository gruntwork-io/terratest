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
	"google.golang.org/api/recaptchaenterprise/v1"
)

// newFakeRecaptchaEnterpriseService points a real reCAPTCHA Enterprise client at a local test server, so the Google transport
// is exercised rather than a hand-written stand-in for a type we do not own.
func newFakeRecaptchaEnterpriseService(t *testing.T, handler http.Handler) *recaptchaenterprise.Service {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	service, err := recaptchaenterprise.NewService(context.Background(),
		option.WithEndpoint(server.URL), option.WithoutAuthentication())
	require.NoError(t, err)

	return service
}

func TestGetRecaptchaKeyAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a key a Gruntwork module created, not a
	// copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/keys/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/keys/gw-library-test","displayName":"terratest key","labels":{"purpose":"terratest"},"webSettings":{"integrationType":"CHECKBOX","allowAllDomains":false,"allowedDomains":["terratest.example.com"]}}`))
	})

	result, err := gcp.GetRecaptchaKeyAttrsWithClient(context.Background(), newFakeRecaptchaEnterpriseService(t, handler), "gw-library-test-project", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "terratest key", result.DisplayName)
	require.NotNil(t, result.WebSettings)
	assert.Equal(t, "CHECKBOX", result.WebSettings.IntegrationType)
	assert.Equal(t, []string{"terratest.example.com"}, result.WebSettings.AllowedDomains)
}

func TestGetRecaptchaKeyAttrsWithClientReportsAMissingOne(t *testing.T) {
	t.Parallel()

	// A caller who asks for something that is not there should be told that, rather than be handed
	// the transport's own wording for a 404.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error":{"code":404,"message":"Resource not found."}}`))
	})

	_, err := gcp.GetRecaptchaKeyAttrsWithClient(context.Background(), newFakeRecaptchaEnterpriseService(t, handler), "gw-library-test-project", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "does not exist")
}
