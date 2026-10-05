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
	"google.golang.org/api/biglake/v1"
	"google.golang.org/api/option"
)

// newFakeBigLakeDatabaseService points a real BigLake client at a local test server, so the Google transport
// is exercised rather than a hand-written stand-in for a type we do not own.
func newFakeBigLakeDatabaseService(t *testing.T, handler http.Handler) *biglake.Service {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	service, err := biglake.NewService(context.Background(),
		option.WithEndpoint(server.URL), option.WithoutAuthentication())
	require.NoError(t, err)

	return service
}

func TestGetBigLakeDatabaseAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a database a Gruntwork module created, not a
	// copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/catalogs/gw_library_catalog/databases/gw_library_test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"...","type":"HIVE","hiveOptions":{"locationUri":"gs://gw-library-test/db","parameters":{"purpose":"terratest"}}}`))
	})

	result, err := gcp.GetBigLakeDatabaseAttrsWithClient(context.Background(), newFakeBigLakeDatabaseService(t, handler), "gw-library-test-project", "us-central1", "gw_library_catalog", "gw_library_test")
	require.NoError(t, err)

	assert.Equal(t, "HIVE", result.Type)
	require.NotNil(t, result.HiveOptions)
	assert.Equal(t, "gs://gw-library-test/db", result.HiveOptions.LocationUri)
}

func TestGetBigLakeTableAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a table a Gruntwork module created, not a
	// copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/catalogs/gw_library_catalog/databases/gw_library_db/tables/gw_library_test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"...","type":"HIVE","hiveOptions":{"tableType":"EXTERNAL_TABLE","parameters":{"purpose":"terratest"}}}`))
	})

	result, err := gcp.GetBigLakeTableAttrsWithClient(context.Background(), newFakeBigLakeDatabaseService(t, handler), "gw-library-test-project", "us-central1", "gw_library_catalog", "gw_library_db", "gw_library_test")
	require.NoError(t, err)

	assert.Equal(t, "HIVE", result.Type)
	require.NotNil(t, result.HiveOptions)
	assert.Equal(t, "EXTERNAL_TABLE", result.HiveOptions.TableType)
}

func TestGetBigLakeDatabaseAttrsWithClientReportsAMissingOne(t *testing.T) {
	t.Parallel()

	// A caller who asks for something that is not there should be told that, rather than be handed the
	// transport's own wording for a 404.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error":{"code":404,"message":"Resource not found."}}`))
	})

	_, err := gcp.GetBigLakeDatabaseAttrsWithClient(context.Background(), newFakeBigLakeDatabaseService(t, handler), "gw-library-test-project", "us-central1", "gw_library_catalog", "gw_library_missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "does not exist")
}
