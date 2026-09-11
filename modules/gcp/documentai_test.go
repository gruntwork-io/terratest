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
	"google.golang.org/api/documentai/v1"
	"google.golang.org/api/option"
)

// newFakeDocumentAIService points a real *documentai.Service at a local test server, so the Google
// transport is exercised rather than a hand-written stand-in for a type we do not own.
func newFakeDocumentAIService(t *testing.T, handler http.Handler) *documentai.Service {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	service, err := documentai.NewService(context.Background(),
		option.WithEndpoint(server.URL), option.WithoutAuthentication())
	require.NoError(t, err)

	return service
}

func TestGetDocumentAIProcessorAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The values are the ones the terraform-google-ml processor module sets, because the point of
	// reading settings back is asserting a module configured the processor it was asked for.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/locations/us/processors/abc123"))
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"name":"projects/gw-library-test-project/locations/us/processors/abc123",
			"type":"OCR_PROCESSOR",
			"displayName":"terratest processor",
			"state":"ENABLED"
		}`))
	})

	processor, err := gcp.GetDocumentAIProcessorAttrsWithClient(context.Background(), newFakeDocumentAIService(t, handler), "gw-library-test-project", "us", "abc123")
	require.NoError(t, err)

	assert.Equal(t, "OCR_PROCESSOR", processor.Type)
	assert.Equal(t, "terratest processor", processor.DisplayName)
	assert.Equal(t, "ENABLED", processor.State)
}

func TestGetDocumentAIProcessorAttrsWithClientMissingProcessor(t *testing.T) {
	t.Parallel()

	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	// The error names the processor and the project as well as saying it is absent, so all three
	// are asserted rather than only the phrase.
	_, err := gcp.GetDocumentAIProcessorAttrsWithClient(context.Background(), newFakeDocumentAIService(t, handler), "gw-library-test-project", "us", "gone")
	require.ErrorContains(t, err, "does not exist")
	require.ErrorContains(t, err, "gone")
	require.ErrorContains(t, err, "gw-library-test-project")
}
