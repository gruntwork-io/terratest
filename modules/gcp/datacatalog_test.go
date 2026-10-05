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
	"google.golang.org/api/datacatalog/v1"
	"google.golang.org/api/option"
)

// newFakeDataCatalogService points a real Data Catalog client at a local test server, so the Google transport
// is exercised rather than a hand-written stand-in for a type we do not own.
func newFakeDataCatalogService(t *testing.T, handler http.Handler) *datacatalog.Service {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	service, err := datacatalog.NewService(context.Background(),
		option.WithEndpoint(server.URL), option.WithoutAuthentication())
	require.NoError(t, err)

	return service
}

func TestGetDataCatalogTaxonomyAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a taxonomy the terraform-google-data-analytics taxonomy module created, not a
	// copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/taxonomies/1234567890"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/taxonomies/1234567890","displayName":"terratest taxonomy","description":"created by terratest","activatedPolicyTypes":["FINE_GRAINED_ACCESS_CONTROL"]}`))
	})

	taxonomy, err := gcp.GetDataCatalogTaxonomyAttrsWithClient(context.Background(), newFakeDataCatalogService(t, handler), "gw-library-test-project", "us-central1", "1234567890")
	require.NoError(t, err)

	assert.Equal(t, "terratest taxonomy", taxonomy.DisplayName)
	assert.Equal(t, []string{"FINE_GRAINED_ACCESS_CONTROL"}, taxonomy.ActivatedPolicyTypes)
}

func TestGetDataCatalogPolicyTagAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a policy tag the terraform-google-data-analytics policy tag module created, not a
	// copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/taxonomies/1234567890/policyTags/987654321"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/taxonomies/1234567890/policyTags/987654321","displayName":"terratest policy tag","description":"created by terratest"}`))
	})

	policyTag, err := gcp.GetDataCatalogPolicyTagAttrsWithClient(context.Background(), newFakeDataCatalogService(t, handler), "gw-library-test-project", "us-central1", "1234567890", "987654321")
	require.NoError(t, err)

	assert.Equal(t, "terratest policy tag", policyTag.DisplayName)
	assert.Equal(t, "created by terratest", policyTag.Description)
}
