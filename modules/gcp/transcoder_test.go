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
	"google.golang.org/api/option"
	"google.golang.org/api/transcoder/v1"
)

// newFakeTranscoderService points a real Transcoder client at a local test server, so the Google transport is
// exercised rather than a hand-written stand-in for a type we do not own.
func newFakeTranscoderService(t *testing.T, handler http.Handler) *transcoder.Service {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	service, err := transcoder.NewService(context.Background(),
		option.WithEndpoint(server.URL), option.WithoutAuthentication())
	require.NoError(t, err)

	return service
}

func TestGetTranscoderJobAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a Transcoder job a Gruntwork module created, not a
	// copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/jobs/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"inputUri":"gs://gw-library-test/input.mp4","outputUri":"gs://gw-library-test/output/","state":"SUCCEEDED","templateId":"preset/web-hd","mode":"PROCESSING_MODE_INTERACTIVE"}`))
	})

	result, err := gcp.GetTranscoderJobAttrsWithClient(context.Background(), newFakeTranscoderService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "gs://gw-library-test/input.mp4", result.InputUri)
	assert.Equal(t, "gs://gw-library-test/output/", result.OutputUri)
	assert.Equal(t, "SUCCEEDED", result.State)
	assert.Equal(t, "preset/web-hd", result.TemplateId)
}

func TestGetTranscoderJobAttrsWithClientReportsAMissingJob(t *testing.T) {
	t.Parallel()

	// A caller who names a Transcoder job that is not there should read a sentence about that Transcoder job, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetTranscoderJobAttrsWithClient(context.Background(), newFakeTranscoderService(t, handler), "gw-library-test-project", "us-central1", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}
