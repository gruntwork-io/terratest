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
	"google.golang.org/api/firestore/v1"
	"google.golang.org/api/option"
)

// newFakeFirestoreAdminService points a real Firestore client at a local test server, so the Google
// transport is exercised rather than a hand-written stand-in for a type we do not own.
func newFakeFirestoreAdminService(t *testing.T, handler http.Handler) *firestore.Service {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	service, err := firestore.NewService(context.Background(),
		option.WithEndpoint(server.URL), option.WithoutAuthentication())
	require.NoError(t, err)

	return service
}

func TestGetFirestoreIndexAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a index the terraform-google-data-storage index module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/databases/gw-library-test/collectionGroups/terratest/indexes/CICAgJjK"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/databases/gw-library-test/collectionGroups/terratest/indexes/CICAgJjK","queryScope":"COLLECTION","state":"READY","fields":[{"fieldPath":"created","order":"DESCENDING"},{"fieldPath":"__name__","order":"DESCENDING"}]}`))
	})

	index, err := gcp.GetFirestoreIndexAttrsWithClient(context.Background(), newFakeFirestoreAdminService(t, handler), "gw-library-test-project", "gw-library-test", "terratest", "CICAgJjK")
	require.NoError(t, err)

	assert.Equal(t, "COLLECTION", index.QueryScope)
	assert.Equal(t, "READY", index.State)
	require.Len(t, index.Fields, 2)
	assert.Equal(t, "created", index.Fields[0].FieldPath)
	assert.Equal(t, "DESCENDING", index.Fields[0].Order)
}

func TestGetFirestoreFieldAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a field the terraform-google-data-storage field module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/databases/gw-library-test/collectionGroups/terratest/fields/created"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/databases/gw-library-test/collectionGroups/terratest/fields/created","indexConfig":{"indexes":[{"queryScope":"COLLECTION","fields":[{"fieldPath":"created","order":"ASCENDING"}]}]}}`))
	})

	field, err := gcp.GetFirestoreFieldAttrsWithClient(context.Background(), newFakeFirestoreAdminService(t, handler), "gw-library-test-project", "gw-library-test", "terratest", "created")
	require.NoError(t, err)

	require.NotNil(t, field.IndexConfig, "the field should carry the index configuration the fixture gave it")
	require.Len(t, field.IndexConfig.Indexes, 1)
	assert.Equal(t, "COLLECTION", field.IndexConfig.Indexes[0].QueryScope)
}

func TestGetFirestoreDocumentAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a document the terraform-google-data-storage document module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/databases/gw-library-test/documents/terratest/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/databases/gw-library-test/documents/terratest/gw-library-test","fields":{"purpose":{"stringValue":"terratest"}}}`))
	})

	document, err := gcp.GetFirestoreDocumentAttrsWithClient(context.Background(), newFakeFirestoreAdminService(t, handler), "gw-library-test-project", "gw-library-test", "terratest", "gw-library-test")
	require.NoError(t, err)

	require.Contains(t, document.Fields, "purpose", "the document should carry the field the fixture wrote")
	assert.Equal(t, "terratest", document.Fields["purpose"].StringValue)
}
