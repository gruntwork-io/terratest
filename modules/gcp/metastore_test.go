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
	"google.golang.org/api/metastore/v1"
	"google.golang.org/api/option"
)

// newFakeMetastoreService points a real Dataproc Metastore client at an httptest server, so a read can be
// exercised against a response we control without reaching Google.
func newFakeMetastoreService(t *testing.T, handler http.Handler) *metastore.APIService {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	service, err := metastore.NewService(context.Background(),
		option.WithoutAuthentication(), option.WithEndpoint(server.URL), option.WithHTTPClient(server.Client()))
	require.NoError(t, err)

	return service
}

func TestGetMetastoreServiceAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a Dataproc Metastore service the terraform-google-data-analytics Dataproc Metastore service module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/services/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/services/gw-library-test","tier":"DEVELOPER","databaseType":"MYSQL","state":"ACTIVE","hiveMetastoreConfig":{"version":"3.1.2"},"labels":{"purpose":"terratest"}}`))
	})

	attrs, err := gcp.GetMetastoreServiceAttrsWithClient(context.Background(), newFakeMetastoreService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "DEVELOPER", attrs.Tier)
	assert.Equal(t, "MYSQL", attrs.DatabaseType)
	assert.Equal(t, "ACTIVE", attrs.State)
}

func TestGetMetastoreServiceAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a Dataproc Metastore service that is not there should read a sentence about that Dataproc Metastore service, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetMetastoreServiceAttrsWithClient(context.Background(), newFakeMetastoreService(t, handler), "gw-library-test-project", "us-central1", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetMetastoreFederationAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a Dataproc Metastore federation the terraform-google-data-analytics Dataproc Metastore federation module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/federations/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/federations/gw-library-test","version":"3.1.2","state":"ACTIVE","labels":{"purpose":"terratest"}}`))
	})

	attrs, err := gcp.GetMetastoreFederationAttrsWithClient(context.Background(), newFakeMetastoreService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "3.1.2", attrs.Version)
	assert.Equal(t, "ACTIVE", attrs.State)
}

func TestGetMetastoreFederationAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a Dataproc Metastore federation that is not there should read a sentence about that Dataproc Metastore federation, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetMetastoreFederationAttrsWithClient(context.Background(), newFakeMetastoreService(t, handler), "gw-library-test-project", "us-central1", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetMetastoreServiceIamPolicyAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a Dataproc Metastore service the terraform-google-data-analytics Dataproc Metastore service module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/services/gw-library-test:getIamPolicy"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"version":3,"etag":"BwXhqw==","bindings":[{"role":"roles/metastore.editor","members":["serviceAccount:gw-library-test@gw-library-test-project.iam.gserviceaccount.com"],"condition":{"title":"until 2030","expression":"request.time < timestamp(\"2030-01-01T00:00:00Z\")"}}]}`))
	})

	policy, err := gcp.GetMetastoreServiceIamPolicyAttrsWithClient(context.Background(), newFakeMetastoreService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, int64(3), policy.Version)
	assert.Equal(t, "roles/metastore.editor", policy.Bindings[0].Role)
}

func TestGetMetastoreServiceIamPolicyAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a Dataproc Metastore service that is not there should read a sentence about that Dataproc Metastore service, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetMetastoreServiceIamPolicyAttrsWithClient(context.Background(), newFakeMetastoreService(t, handler), "gw-library-test-project", "us-central1", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetMetastoreFederationIamPolicyAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a Dataproc Metastore federation the terraform-google-data-analytics Dataproc Metastore federation module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/federations/gw-library-test:getIamPolicy"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"version":3,"etag":"BwXhqw==","bindings":[{"role":"roles/metastore.viewer","members":["serviceAccount:gw-library-test@gw-library-test-project.iam.gserviceaccount.com"],"condition":{"title":"until 2030","expression":"request.time < timestamp(\"2030-01-01T00:00:00Z\")"}}]}`))
	})

	policy, err := gcp.GetMetastoreFederationIamPolicyAttrsWithClient(context.Background(), newFakeMetastoreService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, int64(3), policy.Version)
	assert.Equal(t, "roles/metastore.viewer", policy.Bindings[0].Role)
}

func TestGetMetastoreFederationIamPolicyAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a Dataproc Metastore federation that is not there should read a sentence about that Dataproc Metastore federation, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetMetastoreFederationIamPolicyAttrsWithClient(context.Background(), newFakeMetastoreService(t, handler), "gw-library-test-project", "us-central1", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetMetastoreDatabaseIamPolicyAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a Dataproc Metastore database the terraform-google-data-analytics Dataproc Metastore database module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/services/gw-library-parent/databases/gw-library-test:getIamPolicy"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"version":3,"etag":"BwXhqw==","bindings":[{"role":"roles/metastore.metadataViewer","members":["serviceAccount:gw-library-test@gw-library-test-project.iam.gserviceaccount.com"],"condition":{"title":"until 2030","expression":"request.time < timestamp(\"2030-01-01T00:00:00Z\")"}}]}`))
	})

	policy, err := gcp.GetMetastoreDatabaseIamPolicyAttrsWithClient(context.Background(), newFakeMetastoreService(t, handler), "gw-library-test-project", "us-central1", "gw-library-parent", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, int64(3), policy.Version)
	assert.Equal(t, "roles/metastore.metadataViewer", policy.Bindings[0].Role)
}

func TestGetMetastoreDatabaseIamPolicyAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a Dataproc Metastore database that is not there should read a sentence about that Dataproc Metastore database, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetMetastoreDatabaseIamPolicyAttrsWithClient(context.Background(), newFakeMetastoreService(t, handler), "gw-library-test-project", "us-central1", "gw-library-parent", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetMetastoreTableIamPolicyAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a Dataproc Metastore table the terraform-google-data-analytics Dataproc Metastore table module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/services/gw-library-parent/databases/gw-library-db/tables/gw-library-test:getIamPolicy"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"version":3,"etag":"BwXhqw==","bindings":[{"role":"roles/metastore.metadataViewer","members":["serviceAccount:gw-library-test@gw-library-test-project.iam.gserviceaccount.com"],"condition":{"title":"until 2030","expression":"request.time < timestamp(\"2030-01-01T00:00:00Z\")"}}]}`))
	})

	policy, err := gcp.GetMetastoreTableIamPolicyAttrsWithClient(context.Background(), newFakeMetastoreService(t, handler), "gw-library-test-project", "us-central1", "gw-library-parent", "gw-library-db", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, int64(3), policy.Version)
	assert.Equal(t, "roles/metastore.metadataViewer", policy.Bindings[0].Role)
}

func TestGetMetastoreTableIamPolicyAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a Dataproc Metastore table that is not there should read a sentence about that Dataproc Metastore table, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetMetastoreTableIamPolicyAttrsWithClient(context.Background(), newFakeMetastoreService(t, handler), "gw-library-test-project", "us-central1", "gw-library-parent", "gw-library-db", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}
