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
	"google.golang.org/api/privateca/v1"
)

// newFakePrivateCAService points a real Certificate Authority Service client at a local test server, so the Google transport
// is exercised rather than a hand-written stand-in for a type we do not own.
func newFakePrivateCAService(t *testing.T, handler http.Handler) *privateca.Service {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	service, err := privateca.NewService(context.Background(),
		option.WithEndpoint(server.URL), option.WithoutAuthentication())
	require.NoError(t, err)

	return service
}

func TestGetCertificateTemplateAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a template the terraform-google-security certificate template module created, not a
	// copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/certificateTemplates/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/certificateTemplates/gw-library-test","description":"created by terratest","labels":{"purpose":"terratest"},"identityConstraints":{"allowSubjectPassthrough":false,"allowSubjectAltNamesPassthrough":true},"maximumLifetime":"2592000s"}`))
	})

	template, err := gcp.GetCertificateTemplateAttrsWithClient(context.Background(), newFakePrivateCAService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "created by terratest", template.Description)
	assert.Equal(t, "2592000s", template.MaximumLifetime)
	require.NotNil(t, template.IdentityConstraints)
	assert.True(t, template.IdentityConstraints.AllowSubjectAltNamesPassthrough)
}

func TestGetCaPoolAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a pool the terraform-google-security CA pool module created, not a
	// copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/caPools/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/caPools/gw-library-test","tier":"DEVOPS","labels":{"purpose":"terratest"},"publishingOptions":{"publishCaCert":true,"publishCrl":false}}`))
	})

	pool, err := gcp.GetCaPoolAttrsWithClient(context.Background(), newFakePrivateCAService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "DEVOPS", pool.Tier)
	assert.Equal(t, map[string]string{"purpose": "terratest"}, pool.Labels)
	require.NotNil(t, pool.PublishingOptions)
	assert.True(t, pool.PublishingOptions.PublishCaCert)
}

func TestGetPrivateCACertificateAttrsWithClient(t *testing.T) {
	t.Parallel()

	// A certificate lives under the pool that issued it, so the path is two collections deep.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/caPools/gw-library-pool/certificates/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/caPools/gw-library-pool/certificates/gw-library-test","lifetime":"2592000s","labels":{"purpose":"terratest"}}`))
	})

	certificate, err := gcp.GetPrivateCACertificateAttrsWithClient(context.Background(), newFakePrivateCAService(t, handler), "gw-library-test-project", "us-central1", "gw-library-pool", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "2592000s", certificate.Lifetime)
	assert.Equal(t, map[string]string{"purpose": "terratest"}, certificate.Labels)
}
