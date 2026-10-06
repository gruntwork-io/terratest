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

func TestGetGKEBackupChannelAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a GKE backup channel the the library suite GKE backup channel module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/backupChannels/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/backupChannels/gw-library-test","destinationProject":"projects/gw-library-other-project","description":"created by terratest","labels":{"purpose":"terratest"}}`))
	})

	attrs, err := gcp.GetGKEBackupChannelAttrsWithClient(context.Background(), newFakeGKEBackupService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "projects/gw-library-other-project", attrs.DestinationProject)
	assert.Equal(t, "created by terratest", attrs.Description)
}

func TestGetGKEBackupChannelAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a GKE backup channel that is not there should read a sentence about that GKE backup channel, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetGKEBackupChannelAttrsWithClient(context.Background(), newFakeGKEBackupService(t, handler), "gw-library-test-project", "us-central1", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetGKERestoreChannelAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a GKE restore channel the the library suite GKE restore channel module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/restoreChannels/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/restoreChannels/gw-library-test","destinationProject":"projects/gw-library-other-project","description":"created by terratest","labels":{"purpose":"terratest"}}`))
	})

	attrs, err := gcp.GetGKERestoreChannelAttrsWithClient(context.Background(), newFakeGKEBackupService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "projects/gw-library-other-project", attrs.DestinationProject)
	assert.Equal(t, "created by terratest", attrs.Description)
}

func TestGetGKERestoreChannelAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a GKE restore channel that is not there should read a sentence about that GKE restore channel, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetGKERestoreChannelAttrsWithClient(context.Background(), newFakeGKEBackupService(t, handler), "gw-library-test-project", "us-central1", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetGKEBackupPlanIamPolicyAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a GKE backup plan the the library suite GKE backup plan module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/backupPlans/gw-library-test:getIamPolicy"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"version":3,"etag":"BwXhqw==","bindings":[{"role":"roles/gkebackup.backupAdmin","members":["serviceAccount:gw-library-test@gw-library-test-project.iam.gserviceaccount.com"],"condition":{"title":"until 2030","expression":"request.time < timestamp(\"2030-01-01T00:00:00Z\")"}}]}`))
	})

	policy, err := gcp.GetGKEBackupPlanIamPolicyAttrsWithClient(context.Background(), newFakeGKEBackupService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, int64(3), policy.Version)
	assert.Equal(t, "roles/gkebackup.backupAdmin", policy.Bindings[0].Role)
	assert.Equal(t, `request.time < timestamp("2030-01-01T00:00:00Z")`, policy.Bindings[0].Condition.Expression)
}

func TestGetGKEBackupPlanIamPolicyAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a GKE backup plan that is not there should read a sentence about that GKE backup plan, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetGKEBackupPlanIamPolicyAttrsWithClient(context.Background(), newFakeGKEBackupService(t, handler), "gw-library-test-project", "us-central1", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetGKERestorePlanIamPolicyAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a GKE restore plan the the library suite GKE restore plan module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/restorePlans/gw-library-test:getIamPolicy"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"version":3,"etag":"BwXhqw==","bindings":[{"role":"roles/gkebackup.restoreAdmin","members":["serviceAccount:gw-library-test@gw-library-test-project.iam.gserviceaccount.com"],"condition":{"title":"until 2030","expression":"request.time < timestamp(\"2030-01-01T00:00:00Z\")"}}]}`))
	})

	policy, err := gcp.GetGKERestorePlanIamPolicyAttrsWithClient(context.Background(), newFakeGKEBackupService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, int64(3), policy.Version)
	assert.Equal(t, "roles/gkebackup.restoreAdmin", policy.Bindings[0].Role)
	assert.Equal(t, `request.time < timestamp("2030-01-01T00:00:00Z")`, policy.Bindings[0].Condition.Expression)
}

func TestGetGKERestorePlanIamPolicyAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a GKE restore plan that is not there should read a sentence about that GKE restore plan, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetGKERestorePlanIamPolicyAttrsWithClient(context.Background(), newFakeGKEBackupService(t, handler), "gw-library-test-project", "us-central1", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}
