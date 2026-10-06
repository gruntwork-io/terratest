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
	"google.golang.org/api/eventarc/v1"
	"google.golang.org/api/option"
)

// newFakeEventarcService points a real Eventarc client at a local test server, so the Google transport is
// exercised rather than a hand-written stand-in for a type we do not own.
func newFakeEventarcService(t *testing.T, handler http.Handler) *eventarc.Service {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	service, err := eventarc.NewService(context.Background(),
		option.WithEndpoint(server.URL), option.WithoutAuthentication())
	require.NoError(t, err)

	return service
}

func TestGetEventarcTriggerAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a trigger the terraform-google-
	// messaging module created, not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Contains(t, r.URL.Path, "/locations/us-central1/triggers/gw-library-test", "unexpected path")
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/triggers/gw-library-test","labels":{"purpose":"terratest"},"eventFilters":[{"attribute":"type","value":"google.cloud.pubsub.topic.v1.messagePublished"}],"destination":{"cloudRun":{"service":"gw-library-test","region":"us-central1"}},"serviceAccount":"gw-library-test@gw-library-test-project.iam.gserviceaccount.com"}`))
	})

	trigger, err := gcp.GetEventarcTriggerAttrsWithClient(context.Background(), newFakeEventarcService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	require.Len(t, trigger.EventFilters, 1)
	assert.Equal(t, "type", trigger.EventFilters[0].Attribute)
	assert.Equal(t, "google.cloud.pubsub.topic.v1.messagePublished", trigger.EventFilters[0].Value)
	require.NotNil(t, trigger.Destination)
	require.NotNil(t, trigger.Destination.CloudRun)
	assert.Equal(t, "gw-library-test", trigger.Destination.CloudRun.Service)
}

func TestGetEventarcTriggerAttrsWithClientMissingTrigger(t *testing.T) {
	t.Parallel()

	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	// The error names the trigger and everything that identifies it, and each of those is a value no
	// other part of the message contains, or its check could not fail.
	_, err := gcp.GetEventarcTriggerAttrsWithClient(context.Background(), newFakeEventarcService(t, handler), "gw-library-test-project", "us-central1", "gone")
	require.ErrorContains(t, err, "does not exist")
	require.ErrorContains(t, err, "gone")
	require.ErrorContains(t, err, "us-central1")
	require.ErrorContains(t, err, "gw-library-test-project")
}

func TestGetEventarcChannelAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a channel the terraform-google-messaging channel module created, not a
	// copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/channels/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/channels/gw-library-test","state":"PENDING","provider":"projects/gw-library-test-project/locations/us-central1/providers/datadog","activationToken":"terratest"}`))
	})

	channel, err := gcp.GetEventarcChannelAttrsWithClient(context.Background(), newFakeEventarcService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "PENDING", channel.State)
	assert.Contains(t, channel.Provider, "datadog")
}

func TestGetEventarcPipelineAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a pipeline the terraform-google-messaging pipeline module created, not a
	// copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/pipelines/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/pipelines/gw-library-test","labels":{"purpose":"terratest"},"destinations":[{"topic":"projects/gw-library-test-project/topics/gw-library-test"}],"retryPolicy":{"maxAttempts":3,"minRetryDelay":"5s","maxRetryDelay":"60s"}}`))
	})

	pipeline, err := gcp.GetEventarcPipelineAttrsWithClient(context.Background(), newFakeEventarcService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	require.Len(t, pipeline.Destinations, 1)
	require.NotNil(t, pipeline.RetryPolicy)
	assert.Equal(t, int64(3), pipeline.RetryPolicy.MaxAttempts)
	assert.Equal(t, "5s", pipeline.RetryPolicy.MinRetryDelay)
}

func TestGetEventarcEnrollmentAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a Eventarc enrollment the the library suite Eventarc enrollment module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/enrollments/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/enrollments/gw-library-test","celMatch":"message.type == \"terratest\"","displayName":"terratest enrollment","labels":{"purpose":"terratest"}}`))
	})

	attrs, err := gcp.GetEventarcEnrollmentAttrsWithClient(context.Background(), newFakeEventarcService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, `message.type == "terratest"`, attrs.CelMatch)
	assert.Equal(t, "terratest enrollment", attrs.DisplayName)
}

func TestGetEventarcEnrollmentAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a Eventarc enrollment that is not there should read a sentence about that Eventarc enrollment, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetEventarcEnrollmentAttrsWithClient(context.Background(), newFakeEventarcService(t, handler), "gw-library-test-project", "us-central1", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetEventarcGoogleAPISourceAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a Eventarc Google API source the the library suite Eventarc Google API source module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/googleApiSources/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/googleApiSources/gw-library-test","destination":"projects/gw-library-test-project/locations/us-central1/messageBuses/gw-library-parent","displayName":"terratest source","labels":{"purpose":"terratest"}}`))
	})

	attrs, err := gcp.GetEventarcGoogleAPISourceAttrsWithClient(context.Background(), newFakeEventarcService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "projects/gw-library-test-project/locations/us-central1/messageBuses/gw-library-parent", attrs.Destination)
	assert.Equal(t, "terratest source", attrs.DisplayName)
}

func TestGetEventarcGoogleAPISourceAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a Eventarc Google API source that is not there should read a sentence about that Eventarc Google API source, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetEventarcGoogleAPISourceAttrsWithClient(context.Background(), newFakeEventarcService(t, handler), "gw-library-test-project", "us-central1", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetEventarcMessageBusAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a Eventarc message bus the the library suite Eventarc message bus module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/messageBuses/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/messageBuses/gw-library-test","displayName":"terratest bus","labels":{"purpose":"terratest"}}`))
	})

	attrs, err := gcp.GetEventarcMessageBusAttrsWithClient(context.Background(), newFakeEventarcService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "terratest bus", attrs.DisplayName)
}

func TestGetEventarcMessageBusAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a Eventarc message bus that is not there should read a sentence about that Eventarc message bus, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetEventarcMessageBusAttrsWithClient(context.Background(), newFakeEventarcService(t, handler), "gw-library-test-project", "us-central1", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}
