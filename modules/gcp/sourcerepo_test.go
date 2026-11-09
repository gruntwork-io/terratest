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
	"google.golang.org/api/sourcerepo/v1"
)

// newFakeSourceRepoService points a real Cloud Source Repositories client at a local test server, so
// the Google transport is exercised rather than a hand-written stand-in for a type we do not own.
func newFakeSourceRepoService(t *testing.T, handler http.Handler) *sourcerepo.Service {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	service, err := sourcerepo.NewService(context.Background(),
		option.WithEndpoint(server.URL), option.WithoutAuthentication())
	require.NoError(t, err)

	return service
}

func TestGetSourceRepoRepositoryAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a repository the source repository module
	// created, not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/repos/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/repos/gw-library-test","size":"2048","url":"https://source.developers.google.com/p/gw-library-test-project/r/gw-library-test","pubsubConfigs":{"projects/gw-library-test-project/topics/gw-library-test":{"messageFormat":"JSON","serviceAccountEmail":"terratest@gw-library-test-project.iam.gserviceaccount.com"}}}`))
	})

	repo, err := gcp.GetSourceRepoRepositoryAttrsWithClient(context.Background(), newFakeSourceRepoService(t, handler), "gw-library-test-project", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, int64(2048), repo.Size)
	assert.Equal(t, "https://source.developers.google.com/p/gw-library-test-project/r/gw-library-test", repo.Url)
	assert.Equal(t, "JSON", repo.PubsubConfigs["projects/gw-library-test-project/topics/gw-library-test"].MessageFormat)
}

func TestGetSourceRepoRepositoryAttrsWithClientReportsAMissingRepository(t *testing.T) {
	t.Parallel()

	// A caller who asks for a repository that is not there should be told that, rather than be handed
	// the transport's own wording for a 404.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error":{"code":404,"message":"Not found."}}`))
	})

	_, err := gcp.GetSourceRepoRepositoryAttrsWithClient(context.Background(), newFakeSourceRepoService(t, handler), "gw-library-test-project", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "does not exist")
}
