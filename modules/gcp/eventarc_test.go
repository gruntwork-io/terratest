package gcp_test

import (
	"context"
	"net/http"
	"net/http/httptest"
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
