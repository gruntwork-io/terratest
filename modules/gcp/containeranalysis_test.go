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
	"google.golang.org/api/containeranalysis/v1"
	"google.golang.org/api/option"
)

// newFakeContainerAnalysisService points a real Container Analysis client at a local test server, so the Google transport
// is exercised rather than a hand-written stand-in for a type we do not own.
func newFakeContainerAnalysisService(t *testing.T, handler http.Handler) *containeranalysis.Service {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	service, err := containeranalysis.NewService(context.Background(),
		option.WithEndpoint(server.URL), option.WithoutAuthentication())
	require.NoError(t, err)

	return service
}

func TestGetContainerAnalysisNoteAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a note the terraform-google-containers analysis note module created, not a
	// copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/v1/projects/gw-library-test-project/notes/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/notes/gw-library-test","shortDescription":"terratest note","longDescription":"created by terratest","attestation":{"hint":{"humanReadableName":"terratest attestation"}},"relatedUrl":[{"url":"https://gruntwork.io/terratest","label":"terratest"}]}`))
	})

	note, err := gcp.GetContainerAnalysisNoteAttrsWithClient(context.Background(), newFakeContainerAnalysisService(t, handler), "gw-library-test-project", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "terratest note", note.ShortDescription)
	require.NotNil(t, note.Attestation)
	require.NotNil(t, note.Attestation.Hint)
	assert.Equal(t, "terratest attestation", note.Attestation.Hint.HumanReadableName)
}

func TestGetContainerAnalysisOccurrenceAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for an occurrence the terraform-google-containers analysis occurrence module created, not a
	// copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/v1/projects/gw-library-test-project/occurrences/8f3c1f47-0000-2c0a-0000-000000000000"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/occurrences/8f3c1f47-0000-2c0a-0000-000000000000","resourceUri":"https://gcr.io/gw-library-test-project/terratest@sha256:0000000000000000000000000000000000000000000000000000000000000000","noteName":"projects/gw-library-test-project/notes/gw-library-test","kind":"ATTESTATION","attestation":{"serializedPayload":"dGVycmF0ZXN0"}}`))
	})

	occurrence, err := gcp.GetContainerAnalysisOccurrenceAttrsWithClient(context.Background(), newFakeContainerAnalysisService(t, handler), "gw-library-test-project", "8f3c1f47-0000-2c0a-0000-000000000000")
	require.NoError(t, err)

	assert.Equal(t, "ATTESTATION", occurrence.Kind)
	assert.Equal(t, "projects/gw-library-test-project/notes/gw-library-test", occurrence.NoteName)
	require.NotNil(t, occurrence.Attestation)
}
