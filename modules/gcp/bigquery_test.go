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

func TestGetBigQueryTableAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The values are the ones the terraform-google-data-analytics table module sets, because the
	// point of reading settings back is asserting a module configured the table it was asked for.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/datasets/gw_library_test/tables/events"))
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"tableReference":{"projectId":"gw-library-test-project","datasetId":"gw_library_test","tableId":"events"},
			"friendlyName":"terratest table",
			"description":"created by terratest",
			"type":"TABLE",
			"labels":{"managed-by":"terratest"},
			"schema":{"fields":[
				{"name":"id","type":"STRING","mode":"REQUIRED"},
				{"name":"occurred_at","type":"TIMESTAMP","mode":"NULLABLE"}
			]},
			"timePartitioning":{"type":"DAY","field":"occurred_at"},
			"clustering":{"fields":["id"]}
		}`))
	})

	table, err := gcp.GetBigQueryTableAttrsWithClient(context.Background(), newFakeBigQueryService(t, handler), "gw-library-test-project", "gw_library_test", "events")
	require.NoError(t, err)

	assert.Equal(t, "terratest table", table.FriendlyName)
	assert.Equal(t, "created by terratest", table.Description)
	assert.Equal(t, "TABLE", table.Type)
	assert.Equal(t, map[string]string{"managed-by": "terratest"}, table.Labels)
	require.NotNil(t, table.Schema)
	require.Len(t, table.Schema.Fields, 2)
	assert.Equal(t, "id", table.Schema.Fields[0].Name)
	assert.Equal(t, "REQUIRED", table.Schema.Fields[0].Mode)
	require.NotNil(t, table.TimePartitioning)
	assert.Equal(t, "DAY", table.TimePartitioning.Type)
	assert.Equal(t, "occurred_at", table.TimePartitioning.Field)
	require.NotNil(t, table.Clustering)
	assert.Equal(t, []string{"id"}, table.Clustering.Fields)
}

func TestGetBigQueryTableAttrsWithClientMissingTable(t *testing.T) {
	t.Parallel()

	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	// The error names the dataset, the table and the project as well as saying it is absent, so all
	// of them are asserted rather than only the phrase.
	_, err := gcp.GetBigQueryTableAttrsWithClient(context.Background(), newFakeBigQueryService(t, handler), "gw-library-test-project", "gw_library_test", "gone")
	require.ErrorContains(t, err, "does not exist")
	require.ErrorContains(t, err, "gw_library_test.gone")
	require.ErrorContains(t, err, "gw-library-test-project")
}

func TestGetBigQueryJobAttrsWithClient(t *testing.T) {
	t.Parallel()

	// A job is a record of work rather than a resource that persists, so the response carries the
	// configuration it ran with and the state it finished in.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/jobs/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"kind":"bigquery#job","jobReference":{"projectId":"gw-library-test-project","jobId":"gw-library-test"},"configuration":{"jobType":"QUERY","labels":{"purpose":"terratest"},"query":{"query":"SELECT 1","useLegacySql":false,"priority":"BATCH"}},"status":{"state":"DONE"}}`))
	})

	job, err := gcp.GetBigQueryJobAttrsWithClient(context.Background(), newFakeBigQueryService(t, handler), "gw-library-test-project", "gw-library-test")
	require.NoError(t, err)

	require.NotNil(t, job.Configuration)
	require.NotNil(t, job.Configuration.Query)
	assert.Equal(t, "QUERY", job.Configuration.JobType)
	assert.Equal(t, "SELECT 1", job.Configuration.Query.Query)
	assert.Equal(t, "BATCH", job.Configuration.Query.Priority)
	require.NotNil(t, job.Status)
	assert.Equal(t, "DONE", job.Status.State)
}
