package gcp_test

import (
	"context"
	"net/http"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/gcp/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetCloudBuildBitbucketServerConfigAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a config the Bitbucket Server config
	// module created, not a copy of any one fixture's values. The secrets the module sends are write
	// only, so they are absent from it.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/bitbucketServerConfigs/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/bitbucketServerConfigs/gw-library-test","hostUri":"https://bitbucket.example.com","username":"terratest","webhookKey":"gw-library-webhook","peeredNetwork":"projects/gw-library-test-project/global/networks/gw-library-test"}`))
	})

	config, err := gcp.GetCloudBuildBitbucketServerConfigAttrsWithClient(context.Background(), newFakeCloudBuildRESTService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "https://bitbucket.example.com", config.HostUri)
	assert.Equal(t, "terratest", config.Username)
	assert.Equal(t, "gw-library-webhook", config.WebhookKey)
}

func TestGetCloudBuildBitbucketServerConfigAttrsWithClientReportsAMissingConfig(t *testing.T) {
	t.Parallel()

	// A caller who asks for a config that is not there should be told that, rather than be handed the
	// transport's own wording for a 404.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error":{"code":404,"message":"Not found."}}`))
	})

	_, err := gcp.GetCloudBuildBitbucketServerConfigAttrsWithClient(context.Background(), newFakeCloudBuildRESTService(t, handler), "gw-library-test-project", "us-central1", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "does not exist")
}
