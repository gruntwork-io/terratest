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
	"google.golang.org/api/contentwarehouse/v1"
	"google.golang.org/api/option"
)

// newFakeContentWarehouseService points a real Document AI Warehouse client at a local test server, so the Google transport is
// exercised rather than a hand-written stand-in for a type we do not own.
func newFakeContentWarehouseService(t *testing.T, handler http.Handler) *contentwarehouse.Service {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	service, err := contentwarehouse.NewService(context.Background(),
		option.WithEndpoint(server.URL), option.WithoutAuthentication())
	require.NoError(t, err)

	return service
}

func TestGetDocumentAIWarehouseDocumentSchemaAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a document schema a Gruntwork module created, not a
	// copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/documentSchemas/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"displayName":"terratest schema","description":"created by terratest","documentIsFolder":false,"propertyDefinitions":[{"name":"invoice_total","displayName":"Invoice total","isFilterable":true,"floatTypeOptions":{}}]}`))
	})

	result, err := gcp.GetDocumentAIWarehouseDocumentSchemaAttrsWithClient(context.Background(), newFakeContentWarehouseService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "terratest schema", result.DisplayName)
	assert.False(t, result.DocumentIsFolder)
	require.Len(t, result.PropertyDefinitions, 1)
	assert.Equal(t, "invoice_total", result.PropertyDefinitions[0].Name)
	assert.True(t, result.PropertyDefinitions[0].IsFilterable)
}

func TestGetDocumentAIWarehouseDocumentSchemaAttrsWithClientReportsAMissingSchema(t *testing.T) {
	t.Parallel()

	// A caller who names a document schema that is not there should read a sentence about that document schema, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetDocumentAIWarehouseDocumentSchemaAttrsWithClient(context.Background(), newFakeContentWarehouseService(t, handler), "gw-library-test-project", "us-central1", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}
