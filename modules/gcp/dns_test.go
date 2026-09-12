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
	"google.golang.org/api/dns/v1"
	"google.golang.org/api/option"
)

// newFakeDNSService points a real *dns.Service at a local test server, so the Google transport is
// exercised rather than a hand-written stand-in for a type we do not own.
func newFakeDNSService(t *testing.T, handler http.Handler) *dns.Service {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	service, err := dns.NewService(context.Background(),
		option.WithEndpoint(server.URL), option.WithoutAuthentication())
	require.NoError(t, err)

	return service
}

func TestGetManagedZoneAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The values are the ones the terraform-google-networking managed zone module sets, because the
	// point of reading settings back is asserting a module configured the zone it was asked for.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/managedZones/gw-library-test"))
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"kind":"dns#managedZone",
			"name":"gw-library-test",
			"dnsName":"gw-library-test.example.com.",
			"description":"created by terratest",
			"visibility":"private",
			"labels":{"purpose":"terratest"}
		}`))
	})

	zone, err := gcp.GetManagedZoneAttrsWithClient(context.Background(), newFakeDNSService(t, handler), "gw-library-test-project", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "gw-library-test", zone.Name)
	assert.Equal(t, "gw-library-test.example.com.", zone.DnsName)
	assert.Equal(t, "created by terratest", zone.Description)
	assert.Equal(t, "private", zone.Visibility)
	assert.Equal(t, map[string]string{"purpose": "terratest"}, zone.Labels)
}

func TestGetManagedZoneAttrsWithClientMissingZone(t *testing.T) {
	t.Parallel()

	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	// The error names the zone and the project as well as saying it is absent, so all three are
	// asserted rather than only the phrase.
	_, err := gcp.GetManagedZoneAttrsWithClient(context.Background(), newFakeDNSService(t, handler), "gw-library-test-project", "gone")
	require.ErrorContains(t, err, "does not exist")
	require.ErrorContains(t, err, "gone")
	require.ErrorContains(t, err, "gw-library-test-project")
}
