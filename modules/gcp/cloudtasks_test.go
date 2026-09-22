package gcp_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gruntwork-io/terratest/modules/gcp/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/cloudtasks/v2"
	"google.golang.org/api/option"
)

// newFakeCloudTasksService points a real Cloud Tasks client at a local test server, so the Google transport is
// exercised rather than a hand-written stand-in for a type we do not own.
func newFakeCloudTasksService(t *testing.T, handler http.Handler) *cloudtasks.Service {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	service, err := cloudtasks.NewService(context.Background(),
		option.WithEndpoint(server.URL), option.WithoutAuthentication())
	require.NoError(t, err)

	return service
}

func TestTasksQueueAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a queue the terraform-google-messaging
	// module created, not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Contains(t, r.URL.Path, "/locations/us-central1/queues/gw-library-test", "unexpected path")
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/queues/gw-library-test","state":"PAUSED","rateLimits":{"maxDispatchesPerSecond":2.5,"maxConcurrentDispatches":3},"retryConfig":{"maxAttempts":7,"maxRetryDuration":"600s"}}`))
	})

	queue, err := gcp.GetTasksQueueAttrsWithClient(context.Background(), newFakeCloudTasksService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "PAUSED", queue.State)
	require.NotNil(t, queue.RateLimits)
	assert.InDelta(t, 2.5, queue.RateLimits.MaxDispatchesPerSecond, 1e-9)
	require.NotNil(t, queue.RetryConfig)
	assert.Equal(t, int64(7), queue.RetryConfig.MaxAttempts)
}

func TestTasksQueueAttrsWithClientMissingQueue(t *testing.T) {
	t.Parallel()

	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	// The error names the queue and everything that identifies it, and each of those is a value no
	// other part of the message contains, or its check could not fail.
	_, err := gcp.GetTasksQueueAttrsWithClient(context.Background(), newFakeCloudTasksService(t, handler), "gw-library-test-project", "us-central1", "gone")
	require.ErrorContains(t, err, "does not exist")
	require.ErrorContains(t, err, "gone")
	require.ErrorContains(t, err, "us-central1")
	require.ErrorContains(t, err, "gw-library-test-project")
}
