package gcp_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
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

func TestGetCloudTasksQueueIamPolicyAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a queue the terraform-google-messaging
	// queue IAM policy module granted access on, not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/queues/gw-library-test:getIamPolicy"), "unexpected path %s", r.URL.Path)

		// A conditional binding only comes back at version 3, so the read has to ask for it, which for
		// this call means in the request body rather than in the query.
		var request struct {
			Options struct {
				RequestedPolicyVersion int64 `json:"requestedPolicyVersion"`
			} `json:"options"`
		}

		// assert rather than require: a failed require inside a handler stops the wrong goroutine.
		if assert.NoError(t, json.NewDecoder(r.Body).Decode(&request)) {
			assert.Equal(t, int64(3), request.Options.RequestedPolicyVersion)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"version":3,"etag":"BwXhqw==","bindings":[{"role":"roles/cloudtasks.enqueuer","members":["serviceAccount:gw-library-test@gw-library-test-project.iam.gserviceaccount.com"],"condition":{"title":"until 2030","expression":"request.time < timestamp(\"2030-01-01T00:00:00Z\")"}}]}`))
	})

	policy, err := gcp.GetCloudTasksQueueIamPolicyAttrsWithClient(context.Background(), newFakeCloudTasksService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	require.Len(t, policy.Bindings, 1)
	assert.Equal(t, int64(3), policy.Version)
	assert.Equal(t, "roles/cloudtasks.enqueuer", policy.Bindings[0].Role)
	require.NotNil(t, policy.Bindings[0].Condition, "a conditional binding should keep its condition")
	assert.Equal(t, `request.time < timestamp("2030-01-01T00:00:00Z")`, policy.Bindings[0].Condition.Expression)
}
