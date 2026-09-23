package gcp_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gruntwork-io/terratest/modules/gcp/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/certificatemanager/v1"
	"google.golang.org/api/option"
)

// newFakeCertificateManagerService points a real Certificate Manager client at a local test server, so the Google transport is
// exercised rather than a hand-written stand-in for a type we do not own.
func newFakeCertificateManagerService(t *testing.T, handler http.Handler) *certificatemanager.Service {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	service, err := certificatemanager.NewService(context.Background(),
		option.WithEndpoint(server.URL), option.WithoutAuthentication())
	require.NoError(t, err)

	return service
}

func TestGetCertificateManagerCertificateAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a certificate the terraform-google-
	// security module created, not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Contains(t, r.URL.Path, "/locations/us-central1/certificates/gw-library-test", "unexpected path")
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/certificates/gw-library-test","description":"created by terratest","scope":"DEFAULT","labels":{"purpose":"terratest"},"managed":{"domains":["app.example.com"],"state":"PROVISIONING"}}`))
	})

	certificate, err := gcp.GetCertificateManagerCertificateAttrsWithClient(context.Background(), newFakeCertificateManagerService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "created by terratest", certificate.Description)
	assert.Equal(t, "DEFAULT", certificate.Scope)
	require.NotNil(t, certificate.Managed)
	assert.Equal(t, []string{"app.example.com"}, certificate.Managed.Domains)
}

func TestGetCertificateManagerCertificateAttrsWithClientMissingCertificate(t *testing.T) {
	t.Parallel()

	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	// The error names the certificate and everything that identifies it, and each of those is a
	// value no other part of the message contains, or its check could not fail.
	_, err := gcp.GetCertificateManagerCertificateAttrsWithClient(context.Background(), newFakeCertificateManagerService(t, handler), "gw-library-test-project", "us-central1", "gone")
	require.ErrorContains(t, err, "does not exist")
	require.ErrorContains(t, err, "gone")
	require.ErrorContains(t, err, "us-central1")
	require.ErrorContains(t, err, "gw-library-test-project")
}
