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
	"google.golang.org/api/bigquery/v2"
	"google.golang.org/api/option"
)

// newFakeBigQueryService points a real *bigquery.Service at a local test server, so the Google
// transport is exercised rather than a hand-written stand-in for a type we do not own.
func newFakeBigQueryService(t *testing.T, handler http.Handler) *bigquery.Service {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	service, err := bigquery.NewService(context.Background(),
		option.WithEndpoint(server.URL), option.WithoutAuthentication())
	require.NoError(t, err)

	return service
}

func TestGetBigQueryDatasetAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The values are the ones the terraform-google-data-analytics dataset module sets, because the
	// point of reading settings back is asserting a module configured the dataset it was asked for.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/datasets/gw_library_test"))
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"kind":"bigquery#dataset",
			"id":"gw-library-test-project:gw_library_test",
			"datasetReference":{"datasetId":"gw_library_test","projectId":"gw-library-test-project"},
			"friendlyName":"terratest dataset",
			"description":"created by terratest",
			"location":"US",
			"defaultTableExpirationMs":"7200000",
			"labels":{"purpose":"terratest"}
		}`))
	})

	dataset, err := gcp.GetBigQueryDatasetAttrsWithClient(context.Background(), newFakeBigQueryService(t, handler), "gw-library-test-project", "gw_library_test")
	require.NoError(t, err)

	assert.Equal(t, "gw_library_test", dataset.DatasetReference.DatasetId)
	assert.Equal(t, "terratest dataset", dataset.FriendlyName)
	assert.Equal(t, "created by terratest", dataset.Description)
	assert.Equal(t, "US", dataset.Location)
	assert.Equal(t, int64(7200000), dataset.DefaultTableExpirationMs)
	assert.Equal(t, map[string]string{"purpose": "terratest"}, dataset.Labels)
}

func TestGetBigQueryDatasetAttrsWithClientMissingDataset(t *testing.T) {
	t.Parallel()

	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	// The error names the dataset and the project as well as saying it is absent, so all three are
	// asserted rather than only the phrase.
	_, err := gcp.GetBigQueryDatasetAttrsWithClient(context.Background(), newFakeBigQueryService(t, handler), "gw-library-test-project", "gone")
	require.ErrorContains(t, err, "does not exist")
	require.ErrorContains(t, err, "gone")
	require.ErrorContains(t, err, "gw-library-test-project")
}
