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

func TestGetDeployTargetAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a target the terraform-google-devtools
	// module created, not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/targets/gw-library-test"))
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"name":"projects/gw-library-test-project/locations/us-central1/targets/gw-library-test",
			"description":"created by terratest",
			"requireApproval":true,
			"labels":{"purpose":"terratest"},
			"run":{"location":"projects/gw-library-test-project/locations/us-central1"}
		}`))
	})

	target, err := gcp.GetDeployTargetAttrsWithClient(context.Background(), newFakeCloudDeployService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "created by terratest", target.Description)
	assert.True(t, target.RequireApproval)
	assert.Equal(t, map[string]string{"purpose": "terratest"}, target.Labels)
	require.NotNil(t, target.Run)
	assert.Equal(t, "projects/gw-library-test-project/locations/us-central1", target.Run.Location)
}

func TestGetDeployTargetAttrsWithClientMissingTarget(t *testing.T) {
	t.Parallel()

	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	// The error names the target, its location and the project as well as saying it is absent, so
	// all four are asserted rather than only the phrase.
	_, err := gcp.GetDeployTargetAttrsWithClient(context.Background(), newFakeCloudDeployService(t, handler), "gw-library-test-project", "us-central1", "gone")
	require.ErrorContains(t, err, "does not exist")
	require.ErrorContains(t, err, "gone")
	require.ErrorContains(t, err, "us-central1")
	require.ErrorContains(t, err, "gw-library-test-project")
}

func TestGetCloudDeployAutomationAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for an automation the terraform-google-devtools automation module created, not a
	// copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/deliveryPipelines/gw-library-test/automations/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/deliveryPipelines/gw-library-test/automations/gw-library-test","description":"created by terratest","serviceAccount":"gw-library-test@gw-library-test-project.iam.gserviceaccount.com","suspended":true,"selector":{"targets":[{"id":"gw-library-test"}]},"rules":[{"promoteReleaseRule":{"id":"promote","wait":"3600s"}}]}`))
	})

	automation, err := gcp.GetCloudDeployAutomationAttrsWithClient(context.Background(), newFakeCloudDeployService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "created by terratest", automation.Description)
	assert.True(t, automation.Suspended)
	require.Len(t, automation.Rules, 1)
	require.NotNil(t, automation.Rules[0].PromoteReleaseRule)
	assert.Equal(t, "3600s", automation.Rules[0].PromoteReleaseRule.Wait)
}

func TestGetCloudDeployCustomTargetTypeAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a target type the terraform-google-devtools custom target type module created, not a
	// copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/customTargetTypes/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/customTargetTypes/gw-library-test","description":"created by terratest","customActions":{"renderAction":"render","deployAction":"deploy"},"labels":{"purpose":"terratest"}}`))
	})

	targetType, err := gcp.GetCloudDeployCustomTargetTypeAttrsWithClient(context.Background(), newFakeCloudDeployService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "created by terratest", targetType.Description)
	require.NotNil(t, targetType.CustomActions)
	assert.Equal(t, "deploy", targetType.CustomActions.DeployAction)
}

func TestGetCloudDeployPolicyAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a policy the terraform-google-devtools deploy policy module created, not a
	// copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/deployPolicies/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/deployPolicies/gw-library-test","description":"created by terratest","suspended":true,"selectors":[{"deliveryPipeline":{"id":"gw-library-test"}}],"rules":[{"rolloutRestriction":{"id":"weekend","actions":["ADVANCE"]}}]}`))
	})

	policy, err := gcp.GetCloudDeployPolicyAttrsWithClient(context.Background(), newFakeCloudDeployService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "created by terratest", policy.Description)
	assert.True(t, policy.Suspended)
	require.Len(t, policy.Rules, 1)
	require.NotNil(t, policy.Rules[0].RolloutRestriction)
	assert.Equal(t, []string{"ADVANCE"}, policy.Rules[0].RolloutRestriction.Actions)
}

func TestGetCloudDeployDeliveryPipelineIamPolicyAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a pipeline the terraform-google-devtools delivery pipeline IAM policy module granted
	// access on, not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/deliveryPipelines/gw-library-test:getIamPolicy"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"version":1,"etag":"BwXhqw==","bindings":[{"role":"roles/clouddeploy.releaser","members":["serviceAccount:gw-library-test@gw-library-test-project.iam.gserviceaccount.com"]}]}`))
	})

	policy, err := gcp.GetCloudDeployDeliveryPipelineIamPolicyAttrsWithClient(context.Background(), newFakeCloudDeployService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	require.Len(t, policy.Bindings, 1)
	assert.Equal(t, "roles/clouddeploy.releaser", policy.Bindings[0].Role)
}

func TestGetCloudDeployTargetIamPolicyAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a target the terraform-google-devtools target IAM policy module granted
	// access on, not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/targets/gw-library-test:getIamPolicy"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"version":1,"etag":"BwXhqw==","bindings":[{"role":"roles/clouddeploy.viewer","members":["serviceAccount:gw-library-test@gw-library-test-project.iam.gserviceaccount.com"]}]}`))
	})

	policy, err := gcp.GetCloudDeployTargetIamPolicyAttrsWithClient(context.Background(), newFakeCloudDeployService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	require.Len(t, policy.Bindings, 1)
	assert.Equal(t, "roles/clouddeploy.viewer", policy.Bindings[0].Role)
}
