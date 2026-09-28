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

func TestFetchInstanceTemplateWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a template the terraform-google-compute
	// module created, not a copy of any one fixture's values. Everything a template sets lives under
	// its properties.
	handler := respond(t, http.MethodGet, "/projects/gw-library-test-project/global/instanceTemplates/gw-library-test", http.StatusOK, `{
		"kind":"compute#instanceTemplate",
		"name":"gw-library-test",
		"description":"created by terratest",
		"properties":{
			"machineType":"e2-micro",
			"labels":{"purpose":"terratest"},
			"disks":[{"boot":true,"autoDelete":true,"initializeParams":{"diskSizeGb":"10","sourceImage":"projects/debian-cloud/global/images/family/debian-12"}}],
			"networkInterfaces":[{"network":"https://www.googleapis.com/compute/v1/projects/gw-library-test-project/global/networks/gw-library-test"}],
			"scheduling":{"preemptible":true,"automaticRestart":false}
		}
	}`)

	template, err := gcp.FetchInstanceTemplateWithClient(context.Background(), newFakeComputeService(t, handler), "gw-library-test-project", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "created by terratest", template.Description)
	require.NotNil(t, template.Properties)
	assert.Equal(t, "e2-micro", template.Properties.MachineType)
	assert.Equal(t, "terratest", template.Properties.Labels["purpose"])
	require.Len(t, template.Properties.Disks, 1)
	require.NotNil(t, template.Properties.Disks[0].InitializeParams)
	// The API sends the size as a JSON string and the Go client decodes it to an int64.
	assert.Equal(t, int64(10), template.Properties.Disks[0].InitializeParams.DiskSizeGb)
	require.NotNil(t, template.Properties.Scheduling)
	assert.True(t, template.Properties.Scheduling.Preemptible)
}

func TestFetchInstanceTemplateWithClientMissingTemplate(t *testing.T) {
	t.Parallel()

	_, err := gcp.FetchInstanceTemplateWithClient(context.Background(), newFakeComputeService(t, respond(t, "", "", http.StatusNotFound, "")), "gw-library-test-project", "gone")
	require.ErrorContains(t, err, "InstanceTemplates.Get(gw-library-test-project, gone)")
}

func TestFetchRegionalInstanceGroupManagerWithClient(t *testing.T) {
	t.Parallel()

	handler := respond(t, http.MethodGet, "/projects/gw-library-test-project/regions/us-central1/instanceGroupManagers/gw-library-test", http.StatusOK, `{
		"kind":"compute#instanceGroupManager",
		"name":"gw-library-test",
		"baseInstanceName":"gw-library-test",
		"instanceTemplate":"https://www.googleapis.com/compute/v1/projects/gw-library-test-project/global/instanceTemplates/gw-library-test",
		"targetSize":2,
		"distributionPolicy":{"targetShape":"EVEN"},
		"updatePolicy":{"type":"OPPORTUNISTIC","minimalAction":"REPLACE"}
	}`)

	manager, err := gcp.FetchRegionalInstanceGroupManagerWithClient(context.Background(), newFakeComputeService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "gw-library-test", manager.BaseInstanceName)
	assert.True(t, strings.HasSuffix(manager.InstanceTemplate, "/global/instanceTemplates/gw-library-test"))
	assert.Equal(t, int64(2), manager.TargetSize)
	require.NotNil(t, manager.DistributionPolicy)
	assert.Equal(t, "EVEN", manager.DistributionPolicy.TargetShape)
	require.NotNil(t, manager.UpdatePolicy)
	assert.Equal(t, "OPPORTUNISTIC", manager.UpdatePolicy.Type)
}

func TestFetchRegionalInstanceGroupManagerWithClientMissingManager(t *testing.T) {
	t.Parallel()

	_, err := gcp.FetchRegionalInstanceGroupManagerWithClient(context.Background(), newFakeComputeService(t, respond(t, "", "", http.StatusNotFound, "")), "gw-library-test-project", "us-central1", "gone")
	require.ErrorContains(t, err, "RegionInstanceGroupManagers.Get(gw-library-test-project, us-central1, gone)")
}

func TestFetchRegionInstanceTemplateWithClient(t *testing.T) {
	t.Parallel()

	// The values are the ones the terraform-google-compute region instance template module sets, because the point of reading settings back
	// is asserting a module configured the template it was asked for.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/regions/us-central1/instanceTemplates/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"gw-library-test","description":"created by terratest","properties":{"machineType":"e2-small","labels":{"purpose":"terratest"}}}`))
	})

	template, err := gcp.FetchRegionInstanceTemplateWithClient(context.Background(), newFakeComputeService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "gw-library-test", template.Name)
	assert.Equal(t, "created by terratest", template.Description)
	require.NotNil(t, template.Properties)
	assert.Equal(t, "e2-small", template.Properties.MachineType)
	assert.Equal(t, "terratest", template.Properties.Labels["purpose"])
}

func TestFetchInstanceGroupManagerWithClient(t *testing.T) {
	t.Parallel()

	// The values are the ones the terraform-google-compute instance group manager module sets, because the point of reading settings back
	// is asserting a module configured the group it was asked for.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/zones/us-central1-a/instanceGroupManagers/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"gw-library-test","description":"created by terratest","baseInstanceName":"gw-lib","targetSize":2,"updatePolicy":{"type":"OPPORTUNISTIC","minimalAction":"REPLACE"}}`))
	})

	manager, err := gcp.FetchInstanceGroupManagerWithClient(context.Background(), newFakeComputeService(t, handler), "gw-library-test-project", "us-central1-a", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "gw-library-test", manager.Name)
	assert.Equal(t, "created by terratest", manager.Description)
	assert.Equal(t, "gw-lib", manager.BaseInstanceName)
	assert.Equal(t, int64(2), manager.TargetSize)
	require.NotNil(t, manager.UpdatePolicy)
	assert.Equal(t, "OPPORTUNISTIC", manager.UpdatePolicy.Type)
}

func TestFetchAutoscalerWithClient(t *testing.T) {
	t.Parallel()

	// The values are the ones the terraform-google-compute autoscaler module sets, because the point of reading settings back
	// is asserting a module configured the autoscaler it was asked for.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/zones/us-central1-a/autoscalers/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"gw-library-test","description":"created by terratest","status":"ACTIVE","autoscalingPolicy":{"minNumReplicas":1,"maxNumReplicas":3,"coolDownPeriodSec":90,"cpuUtilization":{"utilizationTarget":0.75}}}`))
	})

	autoscaler, err := gcp.FetchAutoscalerWithClient(context.Background(), newFakeComputeService(t, handler), "gw-library-test-project", "us-central1-a", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "gw-library-test", autoscaler.Name)
	assert.Equal(t, "ACTIVE", autoscaler.Status)
	require.NotNil(t, autoscaler.AutoscalingPolicy)
	assert.Equal(t, int64(1), autoscaler.AutoscalingPolicy.MinNumReplicas)
	assert.Equal(t, int64(3), autoscaler.AutoscalingPolicy.MaxNumReplicas)
	assert.Equal(t, int64(90), autoscaler.AutoscalingPolicy.CoolDownPeriodSec)
	require.NotNil(t, autoscaler.AutoscalingPolicy.CpuUtilization)
	assert.InDelta(t, 0.75, autoscaler.AutoscalingPolicy.CpuUtilization.UtilizationTarget, 0.001)
}
