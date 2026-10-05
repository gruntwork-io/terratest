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
	"google.golang.org/api/binaryauthorization/v1"
	"google.golang.org/api/option"
)

// newFakeBinaryAuthorizationService points a real Binary Authorization client at a local test server, so the Google transport
// is exercised rather than a hand-written stand-in for a type we do not own.
func newFakeBinaryAuthorizationService(t *testing.T, handler http.Handler) *binaryauthorization.Service {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	service, err := binaryauthorization.NewService(context.Background(),
		option.WithEndpoint(server.URL), option.WithoutAuthentication())
	require.NoError(t, err)

	return service
}

func TestGetBinaryAuthorizationAttestorAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for an attestor a Gruntwork module created, not a
	// copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/attestors/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/attestors/gw-library-test","description":"created by terratest","userOwnedGrafeasNote":{"noteReference":"projects/gw-library-test-project/notes/gw-library-test","delegationServiceAccountEmail":"service-37950160017@gcp-sa-binaryauthorization.iam.gserviceaccount.com"}}`))
	})

	result, err := gcp.GetBinaryAuthorizationAttestorAttrsWithClient(context.Background(), newFakeBinaryAuthorizationService(t, handler), "gw-library-test-project", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "created by terratest", result.Description)
	require.NotNil(t, result.UserOwnedGrafeasNote)
	assert.Equal(t, "projects/gw-library-test-project/notes/gw-library-test", result.UserOwnedGrafeasNote.NoteReference)
}

func TestGetBinaryAuthorizationAttestorAttrsWithClientReportsAMissingOne(t *testing.T) {
	t.Parallel()

	// A caller who asks for something that is not there should be told that, rather than be handed
	// the transport's own wording for a 404.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error":{"code":404,"message":"Resource not found."}}`))
	})

	_, err := gcp.GetBinaryAuthorizationAttestorAttrsWithClient(context.Background(), newFakeBinaryAuthorizationService(t, handler), "gw-library-test-project", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "does not exist")
}
