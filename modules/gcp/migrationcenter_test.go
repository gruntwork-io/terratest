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
	"google.golang.org/api/migrationcenter/v1"
	"google.golang.org/api/option"
)

// newFakeMigrationCenterService points a real Migration Center client at a local test server, so the Google transport
// is exercised rather than a hand-written stand-in for a type we do not own.
func newFakeMigrationCenterService(t *testing.T, handler http.Handler) *migrationcenter.Service {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	service, err := migrationcenter.NewService(context.Background(),
		option.WithEndpoint(server.URL), option.WithoutAuthentication())
	require.NoError(t, err)

	return service
}

func TestGetMigrationCenterGroupAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a resource the terraform-google-migration groups module created, not a
	// copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/groups/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/groups/gw-library-test","displayName":"terratest group","description":"created by terratest","labels":{"purpose":"terratest"}}`))
	})

	result, err := gcp.GetMigrationCenterGroupAttrsWithClient(context.Background(), newFakeMigrationCenterService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "terratest group", result.DisplayName)
	assert.Equal(t, map[string]string{"purpose": "terratest"}, result.Labels)
}

func TestGetMigrationCenterSourceAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a resource the terraform-google-migration sources module created, not a
	// copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/sources/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/sources/gw-library-test","displayName":"terratest source","type":"SOURCE_TYPE_MANUAL_UPLOAD","priority":3}`))
	})

	result, err := gcp.GetMigrationCenterSourceAttrsWithClient(context.Background(), newFakeMigrationCenterService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "SOURCE_TYPE_MANUAL_UPLOAD", result.Type)
	assert.Equal(t, int64(3), result.Priority)
}

func TestGetMigrationCenterPreferenceSetAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a resource the terraform-google-migration preferenceSets module created, not a
	// copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/preferenceSets/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/preferenceSets/gw-library-test","displayName":"terratest preference set","virtualMachinePreferences":{"targetProduct":"COMPUTE_MIGRATION_TARGET_PRODUCT_COMPUTE_ENGINE"}}`))
	})

	result, err := gcp.GetMigrationCenterPreferenceSetAttrsWithClient(context.Background(), newFakeMigrationCenterService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	require.NotNil(t, result.VirtualMachinePreferences)
	assert.Equal(t, "COMPUTE_MIGRATION_TARGET_PRODUCT_COMPUTE_ENGINE", result.VirtualMachinePreferences.TargetProduct)
}

func TestGetMigrationCenterReportConfigAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a resource the terraform-google-migration reportConfigs module created, not a
	// copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/reportConfigs/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/reportConfigs/gw-library-test","displayName":"terratest report config","description":"created by terratest"}`))
	})

	result, err := gcp.GetMigrationCenterReportConfigAttrsWithClient(context.Background(), newFakeMigrationCenterService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "terratest report config", result.DisplayName)
	assert.Equal(t, "created by terratest", result.Description)
}

func TestGetMigrationCenterImportJobAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a resource the terraform-google-migration importJobs module created, not a
	// copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/importJobs/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/importJobs/gw-library-test","displayName":"terratest import job","assetSource":"projects/gw-library-test-project/locations/us-central1/sources/gw-library-test"}`))
	})

	result, err := gcp.GetMigrationCenterImportJobAttrsWithClient(context.Background(), newFakeMigrationCenterService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "terratest import job", result.DisplayName)
	assert.Contains(t, result.AssetSource, "sources/gw-library-test")
}

func TestGetMigrationCenterDiscoveryClientAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a resource the terraform-google-migration discoveryClients module created, not a
	// copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/discoveryClients/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/discoveryClients/gw-library-test","displayName":"terratest discovery client","ttl":"3600s","serviceAccount":"terratest@gw-library-test-project.iam.gserviceaccount.com"}`))
	})

	result, err := gcp.GetMigrationCenterDiscoveryClientAttrsWithClient(context.Background(), newFakeMigrationCenterService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "3600s", result.Ttl)
	assert.Contains(t, result.ServiceAccount, "terratest@")
}

func TestGetMigrationCenterReportAttrsWithClient(t *testing.T) {
	t.Parallel()

	// A report lives under the config that produced it rather than under the project, so the path is
	// two collections deep.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/reportConfigs/gw-library-config/reports/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/reportConfigs/gw-library-config/reports/gw-library-test","displayName":"terratest report","type":"TOTAL_COST_OF_OWNERSHIP","state":"SUCCEEDED"}`))
	})

	report, err := gcp.GetMigrationCenterReportAttrsWithClient(context.Background(), newFakeMigrationCenterService(t, handler), "gw-library-test-project", "us-central1", "gw-library-config", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "TOTAL_COST_OF_OWNERSHIP", report.Type)
	assert.Equal(t, "terratest report", report.DisplayName)
}

func TestGetMigrationCenterGroupAttrsWithClientReportsAMissingGroup(t *testing.T) {
	t.Parallel()

	// A caller who asks for a group that is not there should be told that, rather than be handed the
	// transport's own wording for a 404.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error":{"code":404,"message":"Resource not found."}}`))
	})

	_, err := gcp.GetMigrationCenterGroupAttrsWithClient(context.Background(), newFakeMigrationCenterService(t, handler), "gw-library-test-project", "us-central1", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "does not exist")
}
