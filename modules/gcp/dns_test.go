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

func TestGetDNSPolicyAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The values are the ones the terraform-google-networking policy module sets, because the point
	// of reading settings back is asserting a module configured the policy it was asked for.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/policies/gw-library-test"))
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"name":"gw-library-test",
			"description":"created by terratest",
			"enableInboundForwarding":true,
			"enableLogging":true,
			"networks":[{"networkUrl":"https://www.googleapis.com/compute/v1/projects/gw-library-test-project/global/networks/gw-library-test"}]
		}`))
	})

	policy, err := gcp.GetDNSPolicyAttrsWithClient(context.Background(), newFakeDNSService(t, handler), "gw-library-test-project", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "created by terratest", policy.Description)
	assert.True(t, policy.EnableInboundForwarding)
	assert.True(t, policy.EnableLogging)
	require.Len(t, policy.Networks, 1)
	assert.True(t, strings.HasSuffix(policy.Networks[0].NetworkUrl, "/global/networks/gw-library-test"))
}

func TestGetDNSPolicyAttrsWithClientMissingPolicy(t *testing.T) {
	t.Parallel()

	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	// The error names the policy and the project as well as saying it is absent, so all three are
	// asserted rather than only the phrase.
	_, err := gcp.GetDNSPolicyAttrsWithClient(context.Background(), newFakeDNSService(t, handler), "gw-library-test-project", "gone")
	require.ErrorContains(t, err, "does not exist")
	require.ErrorContains(t, err, "gone")
	require.ErrorContains(t, err, "gw-library-test-project")
}

func TestGetDNSRecordSetAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The values are the ones the terraform-google-networking record set module sets, because the
	// point of reading settings back is asserting a module configured the record it was asked for.
	// A record is named by its zone, its fully qualified name and its type together, so all three
	// have to reach the request.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/managedZones/gw-library-test/rrsets/www.gw-library-test.example.com./A"))
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"kind":"dns#resourceRecordSet",
			"name":"www.gw-library-test.example.com.",
			"type":"A",
			"ttl":300,
			"rrdatas":["10.0.0.10","10.0.0.11"]
		}`))
	})

	recordSet, err := gcp.GetDNSRecordSetAttrsWithClient(context.Background(), newFakeDNSService(t, handler), "gw-library-test-project", "gw-library-test", "www.gw-library-test.example.com.", "A")
	require.NoError(t, err)

	assert.Equal(t, "www.gw-library-test.example.com.", recordSet.Name)
	assert.Equal(t, "A", recordSet.Type)
	assert.Equal(t, int64(300), recordSet.Ttl)
	assert.Equal(t, []string{"10.0.0.10", "10.0.0.11"}, recordSet.Rrdatas)
}

func TestGetDNSRecordSetAttrsWithClientMissingRecordSet(t *testing.T) {
	t.Parallel()

	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	// The error names the record, its type, the zone and the project as well as saying it is absent,
	// since a record is only identified by all of them together.
	// The zone is named so that no other value in the error contains it, or its check could not fail.
	_, err := gcp.GetDNSRecordSetAttrsWithClient(context.Background(), newFakeDNSService(t, handler), "gw-library-test-project", "gw-zone", "gone.example.com.", "TXT")
	require.ErrorContains(t, err, "does not exist")
	require.ErrorContains(t, err, "gone.example.com.")
	require.ErrorContains(t, err, "TXT")
	require.ErrorContains(t, err, "DNS managed zone gw-zone ")
	require.ErrorContains(t, err, "gw-library-test-project")
}
