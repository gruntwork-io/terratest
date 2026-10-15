package gcp_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gruntwork-io/terratest/modules/gcp/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/option"
	"google.golang.org/api/siteverification/v1"
)

// newFakeSiteVerificationService points a real Site Verification client at a local test server, so
// the Google transport is exercised rather than a hand-written stand-in for a type we do not own.
func newFakeSiteVerificationService(t *testing.T, handler http.Handler) *siteverification.Service {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	service, err := siteverification.NewService(context.Background(),
		option.WithEndpoint(server.URL), option.WithoutAuthentication())
	require.NoError(t, err)

	return service
}

func TestGetSiteVerificationWebResourceAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a domain the web resource module
	// verified, not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"id":"dns://example.com/","site":{"type":"INET_DOMAIN","identifier":"example.com"},"owners":["terratest@example.com"]}`))
	})

	resource, err := gcp.GetSiteVerificationWebResourceAttrsWithClient(context.Background(), newFakeSiteVerificationService(t, handler), "dns://example.com/")
	require.NoError(t, err)

	require.NotNil(t, resource.Site)
	assert.Equal(t, "INET_DOMAIN", resource.Site.Type)
	assert.Equal(t, "example.com", resource.Site.Identifier)
	assert.Equal(t, []string{"terratest@example.com"}, resource.Owners)
}

func TestGetSiteVerificationWebResourceAttrsWithClientReportsAnUnverifiedSite(t *testing.T) {
	t.Parallel()

	// A site nobody has verified answers 404, and a caller should be told that in those words rather
	// than be handed the transport's own.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error":{"code":404,"message":"Not found."}}`))
	})

	_, err := gcp.GetSiteVerificationWebResourceAttrsWithClient(context.Background(), newFakeSiteVerificationService(t, handler), "dns://missing.example.com/")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "is not a verified site or domain")
}
