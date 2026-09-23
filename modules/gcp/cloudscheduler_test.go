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
	"google.golang.org/api/cloudscheduler/v1"
	"google.golang.org/api/option"
)

// newFakeCloudSchedulerService points a real Cloud Scheduler client at a local test server, so the Google transport is
// exercised rather than a hand-written stand-in for a type we do not own.
func newFakeCloudSchedulerService(t *testing.T, handler http.Handler) *cloudscheduler.Service {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	service, err := cloudscheduler.NewService(context.Background(),
		option.WithEndpoint(server.URL), option.WithoutAuthentication())
	require.NoError(t, err)

	return service
}

func TestSchedulerJobAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a job the terraform-google-messaging
	// module created, not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Contains(t, r.URL.Path, "/locations/us-central1/jobs/gw-library-test", "unexpected path")
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/jobs/gw-library-test","description":"created by terratest","schedule":"0 4 * * *","timeZone":"Etc/UTC","state":"PAUSED","pubsubTarget":{"topicName":"projects/gw-library-test-project/topics/gw-library-test","data":"aGVsbG8="}}`))
	})

	job, err := gcp.GetCloudSchedulerJobAttrsWithClient(context.Background(), newFakeCloudSchedulerService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "0 4 * * *", job.Schedule)
	assert.Equal(t, "Etc/UTC", job.TimeZone)
	assert.Equal(t, "PAUSED", job.State)
	require.NotNil(t, job.PubsubTarget)
	assert.True(t, strings.HasSuffix(job.PubsubTarget.TopicName, "/topics/gw-library-test"))
}

func TestSchedulerJobAttrsWithClientMissingJob(t *testing.T) {
	t.Parallel()

	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	// The error names the job and everything that identifies it, and each of those is a value no
	// other part of the message contains, or its check could not fail.
	_, err := gcp.GetCloudSchedulerJobAttrsWithClient(context.Background(), newFakeCloudSchedulerService(t, handler), "gw-library-test-project", "us-central1", "gone")
	require.ErrorContains(t, err, "does not exist")
	require.ErrorContains(t, err, "gone")
	require.ErrorContains(t, err, "us-central1")
	require.ErrorContains(t, err, "gw-library-test-project")
}
