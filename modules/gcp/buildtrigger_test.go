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

func TestGetBuildTriggerAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The values are the ones the terraform-google-devtools trigger module sets, because the point of
	// reading settings back is asserting a module configured the trigger it was asked for.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/triggers/gw-library-test"))
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"id":"0f1e2d3c",
			"name":"gw-library-test",
			"description":"created by terratest",
			"disabled":true,
			"tags":["terratest"],
			"substitutions":{"_PURPOSE":"terratest"},
			"pubsubConfig":{"topic":"projects/gw-library-test-project/topics/gw-library-test"},
			"build":{"steps":[{"name":"ubuntu","args":["echo","hello"]}],"timeout":"600s"}
		}`))
	})

	trigger, err := gcp.GetBuildTriggerAttrsWithClient(context.Background(), newFakeCloudBuildRESTService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "gw-library-test", trigger.Name)
	assert.Equal(t, "created by terratest", trigger.Description)
	assert.True(t, trigger.Disabled)
	assert.Equal(t, []string{"terratest"}, trigger.Tags)
	assert.Equal(t, map[string]string{"_PURPOSE": "terratest"}, trigger.Substitutions)
	require.NotNil(t, trigger.PubsubConfig)
	assert.Equal(t, "projects/gw-library-test-project/topics/gw-library-test", trigger.PubsubConfig.Topic)
	require.NotNil(t, trigger.Build)
	require.Len(t, trigger.Build.Steps, 1)
	assert.Equal(t, "ubuntu", trigger.Build.Steps[0].Name)
	assert.Equal(t, []string{"echo", "hello"}, trigger.Build.Steps[0].Args)
}

func TestGetBuildTriggerAttrsWithClientMissingTrigger(t *testing.T) {
	t.Parallel()

	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	// The error names the trigger, its location and the project as well as saying it is absent, so
	// all four are asserted rather than only the phrase.
	_, err := gcp.GetBuildTriggerAttrsWithClient(context.Background(), newFakeCloudBuildRESTService(t, handler), "gw-library-test-project", "us-central1", "gone")
	require.ErrorContains(t, err, "does not exist")
	require.ErrorContains(t, err, "gone")
	require.ErrorContains(t, err, "us-central1")
	require.ErrorContains(t, err, "gw-library-test-project")
}
