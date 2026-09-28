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
	"google.golang.org/api/cloudasset/v1"
	"google.golang.org/api/option"
)

// newFakeCloudAssetService points a real Cloud Asset client at a local test server, so the Google transport
// is exercised rather than a hand-written stand-in for a type we do not own.
func newFakeCloudAssetService(t *testing.T, handler http.Handler) *cloudasset.Service {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	service, err := cloudasset.NewService(context.Background(),
		option.WithEndpoint(server.URL), option.WithoutAuthentication())
	require.NoError(t, err)

	return service
}

func TestGetAssetFeedAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a feed the terraform-google-management project feed module created, not a
	// copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/37950160017/feeds/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/37950160017/feeds/gw-library-test","assetTypes":["storage.googleapis.com/Bucket"],"contentType":"RESOURCE","feedOutputConfig":{"pubsubDestination":{"topic":"projects/gw-library-test-project/topics/gw-library-test"}}}`))
	})

	feed, err := gcp.GetAssetFeedAttrsWithClient(context.Background(), newFakeCloudAssetService(t, handler), "37950160017", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, []string{"storage.googleapis.com/Bucket"}, feed.AssetTypes)
	assert.Equal(t, "RESOURCE", feed.ContentType)
	require.NotNil(t, feed.FeedOutputConfig)
	require.NotNil(t, feed.FeedOutputConfig.PubsubDestination)
	assert.Equal(t, "projects/gw-library-test-project/topics/gw-library-test", feed.FeedOutputConfig.PubsubDestination.Topic)
}

func TestGetAssetFeedAttrsWithClientReportsAMissingFeed(t *testing.T) {
	t.Parallel()

	// A caller who asks for a feed that is not there should be told that, rather than be handed the
	// transport's own wording for a 404.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error":{"code":404,"message":"Resource not found."}}`))
	})

	_, err := gcp.GetAssetFeedAttrsWithClient(context.Background(), newFakeCloudAssetService(t, handler), "37950160017", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "does not exist")
}
