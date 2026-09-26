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
	"google.golang.org/api/gkebackup/v1"
	"google.golang.org/api/option"
)

// newFakeGKEBackupService points a real Backup for GKE client at a local test server, so the Google transport
// is exercised rather than a hand-written stand-in for a type we do not own.
func newFakeGKEBackupService(t *testing.T, handler http.Handler) *gkebackup.Service {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	service, err := gkebackup.NewService(context.Background(),
		option.WithEndpoint(server.URL), option.WithoutAuthentication())
	require.NoError(t, err)

	return service
}

func TestGetGKEBackupPlanAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a plan the terraform-google-containers backup plan module created, not a
	// copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/backupPlans/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/backupPlans/gw-library-test","cluster":"projects/gw-library-test-project/locations/us-central1/clusters/gw-library-test","description":"created by terratest","deactivated":true,"labels":{"purpose":"terratest"},"retentionPolicy":{"backupDeleteLockDays":1,"backupRetainDays":7},"backupConfig":{"allNamespaces":true,"includeVolumeData":false}}`))
	})

	plan, err := gcp.GetGKEBackupPlanAttrsWithClient(context.Background(), newFakeGKEBackupService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "created by terratest", plan.Description)
	assert.True(t, plan.Deactivated)
	require.NotNil(t, plan.RetentionPolicy)
	assert.Equal(t, int64(7), plan.RetentionPolicy.BackupRetainDays)
	require.NotNil(t, plan.BackupConfig)
	assert.True(t, plan.BackupConfig.AllNamespaces)
}

func TestGetGKERestorePlanAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a plan the terraform-google-containers restore plan module created, not a
	// copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/restorePlans/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/restorePlans/gw-library-test","backupPlan":"projects/gw-library-test-project/locations/us-central1/backupPlans/gw-library-test","cluster":"projects/gw-library-test-project/locations/us-central1/clusters/gw-library-test","description":"created by terratest","restoreConfig":{"allNamespaces":true,"namespacedResourceRestoreMode":"FAIL_ON_CONFLICT","volumeDataRestorePolicy":"NO_VOLUME_DATA_RESTORATION"}}`))
	})

	plan, err := gcp.GetGKERestorePlanAttrsWithClient(context.Background(), newFakeGKEBackupService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "created by terratest", plan.Description)
	require.NotNil(t, plan.RestoreConfig)
	assert.Equal(t, "FAIL_ON_CONFLICT", plan.RestoreConfig.NamespacedResourceRestoreMode)
	assert.True(t, plan.RestoreConfig.AllNamespaces)
}
