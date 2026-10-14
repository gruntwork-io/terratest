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

// newFakeBigLakeService points a real BigLake client at a local test server, so the Google transport
// is exercised rather than a hand-written stand-in for a type we do not own.
func newFakeBigLakeService(t *testing.T, handler http.Handler) *biglake.Service {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	service, err := biglake.NewService(context.Background(),
		option.WithEndpoint(server.URL), option.WithoutAuthentication())
	require.NoError(t, err)

	return service
}

func TestGetBigLakeCatalogAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a catalog a Gruntwork module created, not a
	// copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/catalogs/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/catalogs/gw-library-test","createTime":"2026-09-28T12:00:00Z"}`))
	})

	result, err := gcp.GetBigLakeCatalogAttrsWithClient(context.Background(), newFakeBigLakeService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "projects/gw-library-test-project/locations/us-central1/catalogs/gw-library-test", result.Name)
	assert.Equal(t, "2026-09-28T12:00:00Z", result.CreateTime)
}

func TestGetBigLakeCatalogAttrsWithClientReportsAMissingOne(t *testing.T) {
	t.Parallel()

	// A caller who asks for something that is not there should be told that, rather than be handed
	// the transport's own wording for a 404.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error":{"code":404,"message":"Resource not found."}}`))
	})

	_, err := gcp.GetBigLakeCatalogAttrsWithClient(context.Background(), newFakeBigLakeService(t, handler), "gw-library-test-project", "us-central1", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetBigLakeIcebergCatalogIamPolicyAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a BigLake Iceberg catalog a Gruntwork module created, not a
	// copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/catalogs/gw_library_test:getIamPolicy"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"version":3,"etag":"BwXhqw==","bindings":[{"role":"roles/biglake.admin","members":["serviceAccount:gw-library-test@gw-library-test-project.iam.gserviceaccount.com"],"condition":{"title":"until 2030","expression":"request.time < timestamp(\"2030-01-01T00:00:00Z\")"}}]}`))
	})

	result, err := gcp.GetBigLakeIcebergCatalogIamPolicyAttrsWithClient(context.Background(), newFakeBigLakeService(t, handler), "gw-library-test-project", "gw_library_test")
	require.NoError(t, err)

	assert.Equal(t, int64(3), result.Version)
	assert.Equal(t, "roles/biglake.admin", result.Bindings[0].Role)
	assert.Equal(t, `request.time < timestamp("2030-01-01T00:00:00Z")`, result.Bindings[0].Condition.Expression)
}

func TestGetBigLakeIcebergTableIamPolicyAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a BigLake Iceberg table a Gruntwork module created, not a
	// copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/catalogs/gw_library_test/namespaces/gw_library_ns/tables/gw_library_table:getIamPolicy"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"version":3,"etag":"BwXhqw==","bindings":[{"role":"roles/biglake.admin","members":["serviceAccount:gw-library-test@gw-library-test-project.iam.gserviceaccount.com"],"condition":{"title":"until 2030","expression":"request.time < timestamp(\"2030-01-01T00:00:00Z\")"}}]}`))
	})

	result, err := gcp.GetBigLakeIcebergTableIamPolicyAttrsWithClient(context.Background(), newFakeBigLakeService(t, handler), "gw-library-test-project", "gw_library_test", "gw_library_ns", "gw_library_table")
	require.NoError(t, err)

	assert.Equal(t, int64(3), result.Version)
	assert.Equal(t, "roles/biglake.admin", result.Bindings[0].Role)
	assert.Equal(t, `request.time < timestamp("2030-01-01T00:00:00Z")`, result.Bindings[0].Condition.Expression)
}

func TestGetBigLakeIcebergCatalogIamPolicyAttrsWithClientReportsAMissingBiglakeIcebergCatalog(t *testing.T) {
	t.Parallel()

	// A caller who names a BigLake Iceberg catalog that is not there should read a sentence about that BigLake Iceberg catalog, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetBigLakeIcebergCatalogIamPolicyAttrsWithClient(context.Background(), newFakeBigLakeService(t, handler), "gw-library-test-project", "gw_library_missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw_library_missing")
	assert.Contains(t, err.Error(), "does not exist")
}
