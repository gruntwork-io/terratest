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

func TestGetDataCatalogEntryAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for an entry the Data Catalog entry module
	// created, not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/entryGroups/gw_library_parent/entries/gw_library_test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/entryGroups/gw_library_parent/entries/gw_library_test","displayName":"terratest entry","description":"created by terratest","userSpecifiedType":"my_type","userSpecifiedSystem":"my_system","linkedResource":"//bigquery.googleapis.com/projects/gw-library-test-project"}`))
	})

	entry, err := gcp.GetDataCatalogEntryAttrsWithClient(context.Background(), newFakeDataCatalogService(t, handler), "gw-library-test-project", "us-central1", "gw_library_parent", "gw_library_test")
	require.NoError(t, err)

	assert.Equal(t, "terratest entry", entry.DisplayName)
	assert.Equal(t, "my_type", entry.UserSpecifiedType)
	assert.Equal(t, "my_system", entry.UserSpecifiedSystem)
}

func TestGetDataCatalogEntryAttrsWithClientReportsAMissingEntry(t *testing.T) {
	t.Parallel()

	// A caller who asks for an entry that is not there should be told that, rather than be handed the
	// transport's own wording for a 404.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error":{"code":404,"message":"Not found."}}`))
	})

	_, err := gcp.GetDataCatalogEntryAttrsWithClient(context.Background(), newFakeDataCatalogService(t, handler), "gw-library-test-project", "us-central1", "gw_library_parent", "gw_library_missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetDataCatalogEntryTagsAttrsWithClientFollowsEveryPage(t *testing.T) {
	t.Parallel()

	// Google pages the tags on an entry, so the read has to follow the page token rather than return
	// only what the first call answered.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		w.Header().Set("Content-Type", "application/json")

		if r.URL.Query().Get("pageToken") == "" {
			_, _ = w.Write([]byte(`{"tags":[{"name":"projects/gw-library-test-project/locations/us-central1/entryGroups/gw_library_parent/entries/gw_library_test/tags/first","template":"projects/gw-library-test-project/locations/us-central1/tagTemplates/gw_library_template","fields":{"source":{"stringValue":"terratest"}}}],"nextPageToken":"second"}`))

			return
		}

		_, _ = w.Write([]byte(`{"tags":[{"name":"projects/gw-library-test-project/locations/us-central1/entryGroups/gw_library_parent/entries/gw_library_test/tags/second","template":"projects/gw-library-test-project/locations/us-central1/tagTemplates/gw_library_template"}]}`))
	})

	tags, err := gcp.GetDataCatalogEntryTagsAttrsWithClient(context.Background(), newFakeDataCatalogService(t, handler), "gw-library-test-project", "us-central1", "gw_library_parent", "gw_library_test")
	require.NoError(t, err)

	require.Len(t, tags, 2)
	assert.Equal(t, "terratest", tags[0].Fields["source"].StringValue)
}

func TestGetDataCatalogEntryTagsAttrsWithClientReportsAMissingEntry(t *testing.T) {
	t.Parallel()

	// Listing the tags of an entry that is not there answers 404, and a caller should be told that in
	// those words rather than be handed the transport's own.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error":{"code":404,"message":"Not found."}}`))
	})

	_, err := gcp.GetDataCatalogEntryTagsAttrsWithClient(context.Background(), newFakeDataCatalogService(t, handler), "gw-library-test-project", "us-central1", "gw_library_parent", "gw_library_missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "does not exist")
}
