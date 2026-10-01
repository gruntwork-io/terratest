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
	"google.golang.org/api/dataplex/v1"
	"google.golang.org/api/option"
)

// newFakeDataplexService points a real *dataplex.Service at a local test server, so the Google
// transport is exercised rather than a hand-written stand-in for a type we do not own.
func newFakeDataplexService(t *testing.T, handler http.Handler) *dataplex.Service {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	service, err := dataplex.NewService(context.Background(),
		option.WithEndpoint(server.URL), option.WithoutAuthentication())
	require.NoError(t, err)

	return service
}

func TestGetDataplexLakeAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The values are the ones the terraform-google-data-analytics lake module sets, because the
	// point of reading settings back is asserting a module configured the lake it was asked for.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/lakes/gw-library-test"))
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"name":"projects/gw-library-test-project/locations/us-central1/lakes/gw-library-test",
			"displayName":"terratest lake",
			"description":"created by terratest",
			"state":"ACTIVE",
			"labels":{"managed-by":"terratest"}
		}`))
	})

	lake, err := gcp.GetDataplexLakeAttrsWithClient(context.Background(), newFakeDataplexService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "terratest lake", lake.DisplayName)
	assert.Equal(t, "created by terratest", lake.Description)
	assert.Equal(t, "ACTIVE", lake.State)
	assert.Equal(t, "terratest", lake.Labels["managed-by"])
}

func TestGetDataplexLakeAttrsWithClientMissingLake(t *testing.T) {
	t.Parallel()

	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	// The error names the lake and the project as well as saying it is absent, so all three are
	// asserted rather than only the phrase.
	_, err := gcp.GetDataplexLakeAttrsWithClient(context.Background(), newFakeDataplexService(t, handler), "gw-library-test-project", "us-central1", "gone")
	require.ErrorContains(t, err, "does not exist")
	require.ErrorContains(t, err, "gone")
	require.ErrorContains(t, err, "gw-library-test-project")
}

func TestGetDataplexZoneAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a zone the terraform-google-data-analytics zone module created, not a
	// copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/lakes/gw-library-test/zones/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/lakes/gw-library-test/zones/gw-library-test","displayName":"terratest zone","type":"RAW","state":"ACTIVE","discoverySpec":{"enabled":true,"schedule":"0 6 * * *"},"resourceSpec":{"locationType":"SINGLE_REGION"}}`))
	})

	zone, err := gcp.GetDataplexZoneAttrsWithClient(context.Background(), newFakeDataplexService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "RAW", zone.Type)
	assert.Equal(t, "terratest zone", zone.DisplayName)
	require.NotNil(t, zone.DiscoverySpec)
	assert.Equal(t, "0 6 * * *", zone.DiscoverySpec.Schedule)
	require.NotNil(t, zone.ResourceSpec)
	assert.Equal(t, "SINGLE_REGION", zone.ResourceSpec.LocationType)
}

func TestGetDataplexAssetAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for an asset the terraform-google-data-analytics asset module created, not a
	// copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/lakes/gw-library-test/zones/gw-library-test/assets/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/lakes/gw-library-test/zones/gw-library-test/assets/gw-library-test","displayName":"terratest asset","state":"ACTIVE","resourceSpec":{"type":"STORAGE_BUCKET","name":"projects/gw-library-test-project/buckets/gw-library-test"},"discoverySpec":{"enabled":false}}`))
	})

	asset, err := gcp.GetDataplexAssetAttrsWithClient(context.Background(), newFakeDataplexService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test", "gw-library-test", "gw-library-test")
	require.NoError(t, err)

	require.NotNil(t, asset.ResourceSpec)
	assert.Equal(t, "STORAGE_BUCKET", asset.ResourceSpec.Type)
	assert.Equal(t, "terratest asset", asset.DisplayName)
}

func TestGetDataplexEntryGroupAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for an entry group the terraform-google-data-analytics entry group module created, not a
	// copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/entryGroups/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/entryGroups/gw-library-test","displayName":"terratest entry group","description":"created by terratest","labels":{"purpose":"terratest"}}`))
	})

	entryGroup, err := gcp.GetDataplexEntryGroupAttrsWithClient(context.Background(), newFakeDataplexService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "terratest entry group", entryGroup.DisplayName)
	assert.Equal(t, "created by terratest", entryGroup.Description)
	assert.Equal(t, "terratest", entryGroup.Labels["purpose"])
}

func TestGetDataplexEntryTypeAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for an entry type the terraform-google-data-analytics entry type module created, not a
	// copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/entryTypes/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/entryTypes/gw-library-test","displayName":"terratest entry type","description":"created by terratest","platform":"terratest","system":"terratest","typeAliases":["TABLE"]}`))
	})

	entryType, err := gcp.GetDataplexEntryTypeAttrsWithClient(context.Background(), newFakeDataplexService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "terratest entry type", entryType.DisplayName)
	assert.Equal(t, "terratest", entryType.Platform)
	assert.Equal(t, []string{"TABLE"}, entryType.TypeAliases)
}

func TestGetDataplexAspectTypeAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for an aspect type the terraform-google-data-analytics aspect type module created, not a
	// copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/aspectTypes/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/aspectTypes/gw-library-test","displayName":"terratest aspect type","description":"created by terratest","metadataTemplate":{"name":"terratest","type":"record","recordFields":[{"name":"owner","type":"string","index":1}]}}`))
	})

	aspectType, err := gcp.GetDataplexAspectTypeAttrsWithClient(context.Background(), newFakeDataplexService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "terratest aspect type", aspectType.DisplayName)
	require.NotNil(t, aspectType.MetadataTemplate)
	require.Len(t, aspectType.MetadataTemplate.RecordFields, 1)
	assert.Equal(t, "owner", aspectType.MetadataTemplate.RecordFields[0].Name)
}

func TestGetDataplexEntryAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for an entry the terraform-google-data-analytics entry module created, not a
	// copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/entryGroups/gw-library-test/entries/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/entryGroups/gw-library-test/entries/gw-library-test","entryType":"projects/gw-library-test-project/locations/us-central1/entryTypes/gw-library-test","fullyQualifiedName":"custom:terratest","entrySource":{"resource":"terratest","system":"terratest","description":"created by terratest"}}`))
	})

	entry, err := gcp.GetDataplexEntryAttrsWithClient(context.Background(), newFakeDataplexService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "custom:terratest", entry.FullyQualifiedName)
	require.NotNil(t, entry.EntrySource)
	assert.Equal(t, "created by terratest", entry.EntrySource.Description)
}

func TestGetDataplexGlossaryAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a glossary the terraform-google-data-analytics glossary module created, not a
	// copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/global/glossaries/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/global/glossaries/gw-library-test","displayName":"terratest glossary","description":"created by terratest","termCount":1,"categoryCount":1}`))
	})

	glossary, err := gcp.GetDataplexGlossaryAttrsWithClient(context.Background(), newFakeDataplexService(t, handler), "gw-library-test-project", "global", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "terratest glossary", glossary.DisplayName)
	assert.Equal(t, "created by terratest", glossary.Description)
}

func TestGetDataplexGlossaryCategoryAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a category the terraform-google-data-analytics glossary category module created, not a
	// copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/global/glossaries/gw-library-test/categories/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/global/glossaries/gw-library-test/categories/gw-library-test","displayName":"terratest category","description":"created by terratest","parent":"projects/gw-library-test-project/locations/global/glossaries/gw-library-test"}`))
	})

	category, err := gcp.GetDataplexGlossaryCategoryAttrsWithClient(context.Background(), newFakeDataplexService(t, handler), "gw-library-test-project", "global", "gw-library-test", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "terratest category", category.DisplayName)
	assert.Equal(t, "created by terratest", category.Description)
}

func TestGetDataplexGlossaryTermAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a term the terraform-google-data-analytics glossary term module created, not a
	// copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/global/glossaries/gw-library-test/terms/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/global/glossaries/gw-library-test/terms/gw-library-test","displayName":"terratest term","description":"created by terratest","parent":"projects/gw-library-test-project/locations/global/glossaries/gw-library-test"}`))
	})

	term, err := gcp.GetDataplexGlossaryTermAttrsWithClient(context.Background(), newFakeDataplexService(t, handler), "gw-library-test-project", "global", "gw-library-test", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "terratest term", term.DisplayName)
	assert.Equal(t, "created by terratest", term.Description)
}

func TestGetDataplexDataScanAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a scan the terraform-google-data-analytics datascan module created, not a
	// copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/dataScans/gw-library-test"), "unexpected path %s", r.URL.Path)

		// The default answer carries no scan specification, so the read has to ask for the full view
		// and this is what fails if it stops doing so.
		assert.Equal(t, "FULL", r.URL.Query().Get("view"))
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/dataScans/gw-library-test","displayName":"terratest scan","state":"ACTIVE","data":{"resource":"//bigquery.googleapis.com/projects/gw-library-test-project/datasets/gw_library_test/tables/gw_library_test"},"executionSpec":{"trigger":{"onDemand":{}}},"dataQualitySpec":{"samplingPercent":10,"rules":[{"column":"purpose","nonNullExpectation":{},"dimension":"COMPLETENESS"}]}}`))
	})

	scan, err := gcp.GetDataplexDataScanAttrsWithClient(context.Background(), newFakeDataplexService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "terratest scan", scan.DisplayName)
	require.NotNil(t, scan.DataQualitySpec)
	assert.InDelta(t, 10, scan.DataQualitySpec.SamplingPercent, 0.001)
	require.Len(t, scan.DataQualitySpec.Rules, 1)
	assert.Equal(t, "COMPLETENESS", scan.DataQualitySpec.Rules[0].Dimension)
}
