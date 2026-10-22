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
	"google.golang.org/api/aiplatform/v1"
	"google.golang.org/api/option"
)

// newFakeVertexAIService points a real Vertex AI client at a local test server, so the Google transport
// is exercised rather than a hand-written stand-in for a type we do not own.
func newFakeVertexAIService(t *testing.T, handler http.Handler) *aiplatform.Service {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	service, err := aiplatform.NewService(context.Background(),
		option.WithEndpoint(server.URL), option.WithoutAuthentication())
	require.NoError(t, err)

	return service
}

func TestGetVertexAIDatasetAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a dataset a Gruntwork module created, not a
	// copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/datasets/1234567890"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"...","displayName":"terratest dataset","description":"created by terratest","metadataSchemaUri":"gs://google-cloud-aiplatform/schema/dataset/metadata/text_1.0.0.yaml","labels":{"purpose":"terratest"}}`))
	})

	result, err := gcp.GetVertexAIDatasetAttrsWithClient(context.Background(), newFakeVertexAIService(t, handler), "gw-library-test-project", "us-central1", "1234567890")
	require.NoError(t, err)

	assert.Equal(t, "terratest dataset", result.DisplayName)
	assert.Contains(t, result.MetadataSchemaUri, "text_1.0.0.yaml")
	assert.Equal(t, map[string]string{"purpose": "terratest"}, result.Labels)
}

func TestGetVertexAITensorboardAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a tensorboard a Gruntwork module created, not a
	// copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/tensorboards/1234567890"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"...","displayName":"terratest tensorboard","description":"created by terratest","isDefault":false,"labels":{"purpose":"terratest"}}`))
	})

	result, err := gcp.GetVertexAITensorboardAttrsWithClient(context.Background(), newFakeVertexAIService(t, handler), "gw-library-test-project", "us-central1", "1234567890")
	require.NoError(t, err)

	assert.Equal(t, "terratest tensorboard", result.DisplayName)
	assert.False(t, result.IsDefault)
}

func TestGetVertexAIFeaturestoreAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a featurestore a Gruntwork module created, not a
	// copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/featurestores/gw_library_test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"...","labels":{"purpose":"terratest"},"onlineServingConfig":{"fixedNodeCount":0},"onlineStorageTtlDays":7}`))
	})

	result, err := gcp.GetVertexAIFeaturestoreAttrsWithClient(context.Background(), newFakeVertexAIService(t, handler), "gw-library-test-project", "us-central1", "gw_library_test")
	require.NoError(t, err)

	require.NotNil(t, result.OnlineServingConfig)
	assert.Equal(t, int64(0), result.OnlineServingConfig.FixedNodeCount)
	assert.Equal(t, int64(7), result.OnlineStorageTtlDays)
}

func TestGetVertexAIFeatureGroupAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a feature group a Gruntwork module created, not a
	// copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/featureGroups/gw_library_test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"...","description":"created by terratest","labels":{"purpose":"terratest"},"bigQuery":{"bigQuerySource":{"inputUri":"bq://gw-library-test-project.gw_library_test.features"},"entityIdColumns":["entity_id"]}}`))
	})

	result, err := gcp.GetVertexAIFeatureGroupAttrsWithClient(context.Background(), newFakeVertexAIService(t, handler), "gw-library-test-project", "us-central1", "gw_library_test")
	require.NoError(t, err)

	assert.Equal(t, "created by terratest", result.Description)
	require.NotNil(t, result.BigQuery)
	require.NotNil(t, result.BigQuery.BigQuerySource)
	assert.Contains(t, result.BigQuery.BigQuerySource.InputUri, "gw_library_test.features")
	assert.Equal(t, []string{"entity_id"}, result.BigQuery.EntityIdColumns)
}

func TestGetVertexAINotebookRuntimeTemplateAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a runtime template a Gruntwork module created, not a
	// copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/notebookRuntimeTemplates/1234567890"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"...","displayName":"terratest runtime template","description":"created by terratest","isDefault":false,"machineSpec":{"machineType":"e2-standard-2"},"idleShutdownConfig":{"idleTimeout":"3600s"}}`))
	})

	result, err := gcp.GetVertexAINotebookRuntimeTemplateAttrsWithClient(context.Background(), newFakeVertexAIService(t, handler), "gw-library-test-project", "us-central1", "1234567890")
	require.NoError(t, err)

	assert.Equal(t, "terratest runtime template", result.DisplayName)
	require.NotNil(t, result.MachineSpec)
	assert.Equal(t, "e2-standard-2", result.MachineSpec.MachineType)
	require.NotNil(t, result.IdleShutdownConfig)
	assert.Equal(t, "3600s", result.IdleShutdownConfig.IdleTimeout)
}

func TestGetVertexAITensorboardExperimentAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for an experiment a Gruntwork module created, not a
	// copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/tensorboards/1234567890/experiments/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"...","displayName":"terratest experiment","description":"created by terratest","labels":{"purpose":"terratest"}}`))
	})

	result, err := gcp.GetVertexAITensorboardExperimentAttrsWithClient(context.Background(), newFakeVertexAIService(t, handler), "gw-library-test-project", "us-central1", "1234567890", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "terratest experiment", result.DisplayName)
	assert.Equal(t, "created by terratest", result.Description)
}

func TestGetVertexAIEntityTypeAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for an entity type a Gruntwork module created, not a
	// copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/featurestores/gw_library_test/entityTypes/gw_library_entity"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"...","description":"created by terratest","labels":{"purpose":"terratest"},"offlineStorageTtlDays":5}`))
	})

	result, err := gcp.GetVertexAIEntityTypeAttrsWithClient(context.Background(), newFakeVertexAIService(t, handler), "gw-library-test-project", "us-central1", "gw_library_test", "gw_library_entity")
	require.NoError(t, err)

	assert.Equal(t, "created by terratest", result.Description)
	assert.Equal(t, int64(5), result.OfflineStorageTtlDays)
}

func TestGetVertexAIFeatureAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a feature a Gruntwork module created, not a
	// copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/featurestores/gw_library_test/entityTypes/gw_library_entity/features/gw_library_feature"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"...","description":"created by terratest","valueType":"STRING","labels":{"purpose":"terratest"}}`))
	})

	result, err := gcp.GetVertexAIFeatureAttrsWithClient(context.Background(), newFakeVertexAIService(t, handler), "gw-library-test-project", "us-central1", "gw_library_test", "gw_library_entity", "gw_library_feature")
	require.NoError(t, err)

	assert.Equal(t, "STRING", result.ValueType)
	assert.Equal(t, "created by terratest", result.Description)
}

func TestGetVertexAIDatasetAttrsWithClientReportsAMissingDataset(t *testing.T) {
	t.Parallel()

	// A caller who asks for a dataset that is not there should be told that, rather than be handed the
	// transport's own wording for a 404.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error":{"code":404,"message":"Resource not found."}}`))
	})

	_, err := gcp.GetVertexAIDatasetAttrsWithClient(context.Background(), newFakeVertexAIService(t, handler), "gw-library-test-project", "us-central1", "999")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "does not exist")
}

func TestNewVertexAIServiceERefusesAnInvalidLocation(t *testing.T) {
	t.Parallel()

	// The location goes into the endpoint host, so anything but a plain location name would send the
	// request somewhere the caller did not name.
	_, err := gcp.NewVertexAIServiceE(t, context.Background(), "us-central1/../evil.example.com")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "is not a valid location")
}

func TestGetVertexAIEndpointAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a Vertex AI endpoint a Gruntwork module created, not a
	// copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/endpoints/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"displayName":"terratest endpoint","description":"created by terratest","trafficSplit":{"1234567890":100},"dedicatedEndpointEnabled":true}`))
	})

	result, err := gcp.GetVertexAIEndpointAttrsWithClient(context.Background(), newFakeVertexAIService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "terratest endpoint", result.DisplayName)
	assert.Equal(t, map[string]int64{"1234567890": 100}, result.TrafficSplit)
	assert.True(t, result.DedicatedEndpointEnabled)
}

func TestGetVertexAIIndexAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a Vertex AI index a Gruntwork module created, not a
	// copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/indexes/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"displayName":"terratest index","metadataSchemaUri":"gs://google-cloud-aiplatform/schema/matchingengine/metadata/nearest_neighbor_search_1.0.0.yaml","indexUpdateMethod":"BATCH_UPDATE"}`))
	})

	result, err := gcp.GetVertexAIIndexAttrsWithClient(context.Background(), newFakeVertexAIService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "terratest index", result.DisplayName)
	assert.Equal(t, "BATCH_UPDATE", result.IndexUpdateMethod)
	assert.Contains(t, result.MetadataSchemaUri, "nearest_neighbor_search_1.0.0.yaml")
}

func TestGetVertexAIIndexEndpointAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a Vertex AI index endpoint a Gruntwork module created, not a
	// copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/indexEndpoints/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"displayName":"terratest index endpoint","publicEndpointEnabled":true,"deployedIndexes":[{"id":"gw_library_test","index":"projects/gw-library-test-project/locations/us-central1/indexes/1234567890"}]}`))
	})

	result, err := gcp.GetVertexAIIndexEndpointAttrsWithClient(context.Background(), newFakeVertexAIService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "terratest index endpoint", result.DisplayName)
	assert.True(t, result.PublicEndpointEnabled)
	require.Len(t, result.DeployedIndexes, 1)
	assert.Equal(t, "gw_library_test", result.DeployedIndexes[0].Id)
}

func TestGetVertexAIFeatureOnlineStoreAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a Vertex AI feature online store a Gruntwork module created, not a
	// copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/featureOnlineStores/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"state":"STABLE","labels":{"purpose":"terratest"},"bigtable":{"autoScaling":{"minNodeCount":1,"maxNodeCount":2,"cpuUtilizationTarget":50}}}`))
	})

	result, err := gcp.GetVertexAIFeatureOnlineStoreAttrsWithClient(context.Background(), newFakeVertexAIService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "STABLE", result.State)
	assert.Equal(t, map[string]string{"purpose": "terratest"}, result.Labels)
	require.NotNil(t, result.Bigtable)
	require.NotNil(t, result.Bigtable.AutoScaling)
	assert.Equal(t, int64(1), result.Bigtable.AutoScaling.MinNodeCount)
}

func TestGetVertexAIFeatureViewAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a Vertex AI feature view a Gruntwork module created, not a
	// copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/featureOnlineStores/gw_library_store/featureViews/gw_library_view"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"bigQuerySource":{"uri":"bq://gw-library-test-project.gw_library_test.features","entityIdColumns":["entity_id"]},"syncConfig":{"cron":"0 0 * * *"}}`))
	})

	result, err := gcp.GetVertexAIFeatureViewAttrsWithClient(context.Background(), newFakeVertexAIService(t, handler), "gw-library-test-project", "us-central1", "gw_library_store", "gw_library_view")
	require.NoError(t, err)

	require.NotNil(t, result.BigQuerySource)
	assert.Equal(t, []string{"entity_id"}, result.BigQuerySource.EntityIdColumns)
	require.NotNil(t, result.SyncConfig)
	assert.Equal(t, "0 0 * * *", result.SyncConfig.Cron)
}

func TestGetVertexAIDeploymentResourcePoolAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a Vertex AI deployment resource pool a Gruntwork module created, not a
	// copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/deploymentResourcePools/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"dedicatedResources":{"machineSpec":{"machineType":"n1-standard-2"},"minReplicaCount":1,"maxReplicaCount":1},"disableContainerLogging":true}`))
	})

	result, err := gcp.GetVertexAIDeploymentResourcePoolAttrsWithClient(context.Background(), newFakeVertexAIService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	require.NotNil(t, result.DedicatedResources)
	require.NotNil(t, result.DedicatedResources.MachineSpec)
	assert.Equal(t, "n1-standard-2", result.DedicatedResources.MachineSpec.MachineType)
	assert.Equal(t, int64(1), result.DedicatedResources.MinReplicaCount)
	assert.True(t, result.DisableContainerLogging)
}

func TestGetVertexAIPersistentResourceAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a Vertex AI persistent resource a Gruntwork module created, not a
	// copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/persistentResources/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"displayName":"terratest persistent resource","state":"RUNNING","resourcePools":[{"id":"gw-library-test","machineSpec":{"machineType":"n1-standard-4"},"replicaCount":"1"}]}`))
	})

	result, err := gcp.GetVertexAIPersistentResourceAttrsWithClient(context.Background(), newFakeVertexAIService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "terratest persistent resource", result.DisplayName)
	assert.Equal(t, "RUNNING", result.State)
	require.Len(t, result.ResourcePools, 1)
	assert.Equal(t, "n1-standard-4", result.ResourcePools[0].MachineSpec.MachineType)
}

func TestGetVertexAIScheduleAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a Vertex AI schedule a Gruntwork module created, not a
	// copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/schedules/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"displayName":"terratest schedule","cron":"0 9 * * MON","state":"ACTIVE","maxConcurrentRunCount":"1","allowQueueing":true}`))
	})

	result, err := gcp.GetVertexAIScheduleAttrsWithClient(context.Background(), newFakeVertexAIService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "terratest schedule", result.DisplayName)
	assert.Equal(t, "0 9 * * MON", result.Cron)
	assert.Equal(t, "ACTIVE", result.State)
	assert.Equal(t, int64(1), result.MaxConcurrentRunCount)
	assert.True(t, result.AllowQueueing)
}

func TestGetVertexAIReasoningEngineAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a Vertex AI reasoning engine a Gruntwork module created, not a
	// copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/reasoningEngines/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"displayName":"terratest reasoning engine","description":"created by terratest","spec":{"packageSpec":{"pythonVersion":"3.11","pickleObjectGcsUri":"gs://gw-library-test/agent.pkl"}}}`))
	})

	result, err := gcp.GetVertexAIReasoningEngineAttrsWithClient(context.Background(), newFakeVertexAIService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "terratest reasoning engine", result.DisplayName)
	require.NotNil(t, result.Spec)
	require.NotNil(t, result.Spec.PackageSpec)
	assert.Equal(t, "3.11", result.Spec.PackageSpec.PythonVersion)
}

func TestGetVertexAIFeatureGroupFeatureAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a Vertex AI feature group feature a Gruntwork module created, not a
	// copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/featureGroups/gw_library_group/features/gw_library_feature"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"description":"created by terratest","versionColumnName":"feature_value","labels":{"purpose":"terratest"}}`))
	})

	result, err := gcp.GetVertexAIFeatureGroupFeatureAttrsWithClient(context.Background(), newFakeVertexAIService(t, handler), "gw-library-test-project", "us-central1", "gw_library_group", "gw_library_feature")
	require.NoError(t, err)

	assert.Equal(t, "created by terratest", result.Description)
	assert.Equal(t, "feature_value", result.VersionColumnName)
	assert.Equal(t, map[string]string{"purpose": "terratest"}, result.Labels)
}

func TestGetVertexAIEndpointAttrsWithClientReportsAMissingVertexAiEndpoint(t *testing.T) {
	t.Parallel()

	// A caller who names a Vertex AI endpoint that is not there should read a sentence about that Vertex AI endpoint, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetVertexAIEndpointAttrsWithClient(context.Background(), newFakeVertexAIService(t, handler), "gw-library-test-project", "us-central1", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetVertexAICacheConfigAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns once the cache config module has turned
	// context caching off, not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/cacheConfig"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/cacheConfig","disableCache":true}`))
	})

	config, err := gcp.GetVertexAICacheConfigAttrsWithClient(context.Background(), newFakeVertexAIService(t, handler), "gw-library-test-project")
	require.NoError(t, err)

	assert.True(t, config.DisableCache)
}

func TestGetVertexAICacheConfigAttrsWithClientReportsAMissingConfig(t *testing.T) {
	t.Parallel()

	// A project Vertex AI has never been enabled on answers 404, and a caller should be told that in
	// those words rather than be handed the transport's own.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error":{"code":404,"message":"Not found."}}`))
	})

	_, err := gcp.GetVertexAICacheConfigAttrsWithClient(context.Background(), newFakeVertexAIService(t, handler), "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "has no Vertex AI cache config")
}

func TestGetVertexAIRagEngineConfigAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns once the RAG engine config module has chosen
	// a database tier, not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/ragEngineConfig"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/ragEngineConfig","ragManagedDbConfig":{"basic":{}}}`))
	})

	config, err := gcp.GetVertexAIRagEngineConfigAttrsWithClient(context.Background(), newFakeVertexAIService(t, handler), "gw-library-test-project", "us-central1")
	require.NoError(t, err)

	require.NotNil(t, config.RagManagedDbConfig)
	assert.NotNil(t, config.RagManagedDbConfig.Basic)
}

func TestGetVertexAIRagEngineConfigAttrsWithClientReportsAMissingConfig(t *testing.T) {
	t.Parallel()

	// A location the RAG engine is not available in answers 404, and a caller should be told that in
	// those words rather than be handed the transport's own.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error":{"code":404,"message":"Not found."}}`))
	})

	_, err := gcp.GetVertexAIRagEngineConfigAttrsWithClient(context.Background(), newFakeVertexAIService(t, handler), "gw-library-test-project", "us-central1")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "no Vertex AI RAG engine config exists")
}

func TestGetVertexAINotebookRuntimeAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a runtime the Colab Enterprise runtime
	// module created, not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/notebookRuntimes/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/notebookRuntimes/gw-library-test","displayName":"terratest runtime","runtimeState":"RUNNING","healthState":"HEALTHY","notebookRuntimeTemplateRef":{"notebookRuntimeTemplate":"projects/gw-library-test-project/locations/us-central1/notebookRuntimeTemplates/gw-library-template"},"machineSpec":{"machineType":"e2-standard-4"}}`))
	})

	runtime, err := gcp.GetVertexAINotebookRuntimeAttrsWithClient(context.Background(), newFakeVertexAIService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "terratest runtime", runtime.DisplayName)
	require.NotNil(t, runtime.NotebookRuntimeTemplateRef)
	assert.Contains(t, runtime.NotebookRuntimeTemplateRef.NotebookRuntimeTemplate, "/notebookRuntimeTemplates/gw-library-template")
}

func TestGetVertexAINotebookRuntimeAttrsWithClientReportsAMissingRuntime(t *testing.T) {
	t.Parallel()

	// A caller who asks for a runtime that is not there should be told that, rather than be handed the
	// transport's own wording for a 404.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error":{"code":404,"message":"Not found."}}`))
	})

	_, err := gcp.GetVertexAINotebookRuntimeAttrsWithClient(context.Background(), newFakeVertexAIService(t, handler), "gw-library-test-project", "us-central1", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetVertexAINotebookExecutionJobAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a job the notebook execution module
	// created, not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/notebookExecutionJobs/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/notebookExecutionJobs/gw-library-test","displayName":"terratest job","jobState":"JOB_STATE_SUCCEEDED","gcsOutputUri":"gs://gw-library-test/output","gcsNotebookSource":{"uri":"gs://gw-library-test/notebook.ipynb"},"executionTimeout":"600s"}`))
	})

	job, err := gcp.GetVertexAINotebookExecutionJobAttrsWithClient(context.Background(), newFakeVertexAIService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "JOB_STATE_SUCCEEDED", job.JobState)
	assert.Equal(t, "gs://gw-library-test/output", job.GcsOutputUri)
	require.NotNil(t, job.GcsNotebookSource)
	assert.Equal(t, "gs://gw-library-test/notebook.ipynb", job.GcsNotebookSource.Uri)
}

func TestGetVertexAINotebookExecutionJobAttrsWithClientReportsAMissingJob(t *testing.T) {
	t.Parallel()

	// A caller who asks for a job that is not there should be told that, rather than be handed the
	// transport's own wording for a 404.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error":{"code":404,"message":"Not found."}}`))
	})

	_, err := gcp.GetVertexAINotebookExecutionJobAttrsWithClient(context.Background(), newFakeVertexAIService(t, handler), "gw-library-test-project", "us-central1", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "does not exist")
}
