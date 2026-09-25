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
	"google.golang.org/api/analyticshub/v1"
	"google.golang.org/api/option"
)

// newFakeAnalyticsHubService points a real Analytics Hub client at a local test server, so the Google transport
// is exercised rather than a hand-written stand-in for a type we do not own.
func newFakeAnalyticsHubService(t *testing.T, handler http.Handler) *analyticshub.Service {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	service, err := analyticshub.NewService(context.Background(),
		option.WithEndpoint(server.URL), option.WithoutAuthentication())
	require.NoError(t, err)

	return service
}

func TestGetAnalyticsHubDataExchangeAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for an exchange the terraform-google-data-analytics data exchange module created, not a
	// copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/dataExchanges/gw_library_test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/dataExchanges/gw_library_test","displayName":"terratest exchange","description":"created by terratest","primaryContact":"terratest@gruntwork.io","documentation":"https://example.com/terratest","discoveryType":"DISCOVERY_TYPE_PRIVATE"}`))
	})

	exchange, err := gcp.GetAnalyticsHubDataExchangeAttrsWithClient(context.Background(), newFakeAnalyticsHubService(t, handler), "gw-library-test-project", "us-central1", "gw_library_test")
	require.NoError(t, err)

	assert.Equal(t, "terratest exchange", exchange.DisplayName)
	assert.Equal(t, "terratest@gruntwork.io", exchange.PrimaryContact)
	assert.Equal(t, "DISCOVERY_TYPE_PRIVATE", exchange.DiscoveryType)
}

func TestGetAnalyticsHubListingAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a listing the terraform-google-data-analytics listing module created, not a
	// copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/dataExchanges/gw_library_test/listings/gw_library_test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/dataExchanges/gw_library_test/listings/gw_library_test","displayName":"terratest listing","description":"created by terratest","state":"ACTIVE","categories":["CATEGORY_OTHERS"],"requestAccess":"terratest@gruntwork.io","bigqueryDataset":{"dataset":"projects/gw-library-test-project/datasets/gw_library_test"}}`))
	})

	listing, err := gcp.GetAnalyticsHubListingAttrsWithClient(context.Background(), newFakeAnalyticsHubService(t, handler), "gw-library-test-project", "us-central1", "gw_library_test", "gw_library_test")
	require.NoError(t, err)

	assert.Equal(t, "terratest listing", listing.DisplayName)
	assert.Equal(t, []string{"CATEGORY_OTHERS"}, listing.Categories)
	require.NotNil(t, listing.BigqueryDataset)
	assert.Equal(t, "projects/gw-library-test-project/datasets/gw_library_test", listing.BigqueryDataset.Dataset)
}
