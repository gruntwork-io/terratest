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
	"google.golang.org/api/managedidentities/v1"
	"google.golang.org/api/option"
)

// newFakeManagedIdentitiesService points a real Managed Microsoft AD client at a local test server,
// so the Google transport is exercised rather than a hand-written stand-in for a type we do not own.
func newFakeManagedIdentitiesService(t *testing.T, handler http.Handler) *managedidentities.Service {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	service, err := managedidentities.NewService(context.Background(),
		option.WithEndpoint(server.URL), option.WithoutAuthentication())
	require.NoError(t, err)

	return service
}

func TestGetActiveDirectoryDomainAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a domain the Active Directory domain
	// module created, with the trust the domain trust module added, not a copy of any one fixture's
	// values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/global/domains/terratest.example.com"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/global/domains/terratest.example.com","fqdn":"terratest.example.com","state":"READY","locations":["us-central1"],"reservedIpRange":"10.1.0.0/24","authorizedNetworks":["projects/gw-library-test-project/global/networks/gw-library-test"],"admin":"setupadmin","trusts":[{"targetDomainName":"peer.example.com","trustType":"FOREST","trustDirection":"BIDIRECTIONAL","state":"CONNECTED"}]}`))
	})

	domain, err := gcp.GetActiveDirectoryDomainAttrsWithClient(context.Background(), newFakeManagedIdentitiesService(t, handler), "gw-library-test-project", "terratest.example.com")
	require.NoError(t, err)

	assert.Equal(t, "READY", domain.State)
	assert.Equal(t, "10.1.0.0/24", domain.ReservedIpRange)
	require.Len(t, domain.Trusts, 1)
	assert.Equal(t, "peer.example.com", domain.Trusts[0].TargetDomainName)
	assert.Equal(t, "BIDIRECTIONAL", domain.Trusts[0].TrustDirection)
}

func TestGetActiveDirectoryDomainAttrsWithClientReportsAMissingDomain(t *testing.T) {
	t.Parallel()

	// A caller who asks for a domain that is not there should be told that, rather than be handed the
	// transport's own wording for a 404.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error":{"code":404,"message":"Not found."}}`))
	})

	_, err := gcp.GetActiveDirectoryDomainAttrsWithClient(context.Background(), newFakeManagedIdentitiesService(t, handler), "gw-library-test-project", "missing.example.com")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "does not exist")
}
