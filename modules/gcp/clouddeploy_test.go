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
	"google.golang.org/api/clouddeploy/v1"
	"google.golang.org/api/option"
)

// newFakeCloudDeployService points a real *clouddeploy.Service at a local test server, so the
// Google transport is exercised rather than a hand-written stand-in for a type we do not own.
func newFakeCloudDeployService(t *testing.T, handler http.Handler) *clouddeploy.Service {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	service, err := clouddeploy.NewService(context.Background(),
		option.WithEndpoint(server.URL), option.WithoutAuthentication())
	require.NoError(t, err)

	return service
}

func TestGetDeliveryPipelineAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The values are the ones the terraform-google-devtools pipeline module sets, because the point
	// of reading settings back is asserting a module configured the pipeline it was asked for.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/locations/us-central1/deliveryPipelines/gw-library-test"))
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"name":"projects/gw-library-test-project/locations/us-central1/deliveryPipelines/gw-library-test",
			"description":"created by terratest",
			"labels":{"purpose":"terratest"},
			"annotations":{"owner":"terratest"},
			"suspended":true
		}`))
	})

	pipeline, err := gcp.GetDeliveryPipelineAttrsWithClient(context.Background(), newFakeCloudDeployService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "created by terratest", pipeline.Description)
	assert.Equal(t, map[string]string{"purpose": "terratest"}, pipeline.Labels)
	assert.Equal(t, map[string]string{"owner": "terratest"}, pipeline.Annotations)
	assert.True(t, pipeline.Suspended)
}

func TestGetDeliveryPipelineAttrsWithClientMissingPipeline(t *testing.T) {
	t.Parallel()

	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	// The error names the pipeline and the project as well as saying it is absent, so all three are
	// asserted rather than only the phrase.
	_, err := gcp.GetDeliveryPipelineAttrsWithClient(context.Background(), newFakeCloudDeployService(t, handler), "gw-library-test-project", "us-central1", "gone")
	require.ErrorContains(t, err, "does not exist")
	require.ErrorContains(t, err, "gone")
	require.ErrorContains(t, err, "gw-library-test-project")
}
