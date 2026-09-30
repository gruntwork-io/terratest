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
