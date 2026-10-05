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
	"google.golang.org/api/backupdr/v1"
	"google.golang.org/api/option"
)

// newFakeBackupDRService points a real Backup and DR client at a local test server, so the Google transport
// is exercised rather than a hand-written stand-in for a type we do not own.
func newFakeBackupDRService(t *testing.T, handler http.Handler) *backupdr.Service {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	service, err := backupdr.NewService(context.Background(),
		option.WithEndpoint(server.URL), option.WithoutAuthentication())
	require.NoError(t, err)

	return service
}

func TestGetBackupDRBackupVaultAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a backup vault a Gruntwork module created, not a
	// copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/backupVaults/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/backupVaults/gw-library-test","description":"created by terratest","backupMinimumEnforcedRetentionDuration":"100000s","labels":{"purpose":"terratest"}}`))
	})

	result, err := gcp.GetBackupDRBackupVaultAttrsWithClient(context.Background(), newFakeBackupDRService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "created by terratest", result.Description)
	assert.Equal(t, "100000s", result.BackupMinimumEnforcedRetentionDuration)
}

func TestGetBackupDRBackupPlanAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a backup plan a Gruntwork module created, not a
	// copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/backupPlans/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/backupPlans/gw-library-test","description":"created by terratest","resourceType":"compute.googleapis.com/Instance","backupRules":[{"ruleId":"terratest-daily","backupRetentionDays":5}]}`))
	})

	result, err := gcp.GetBackupDRBackupPlanAttrsWithClient(context.Background(), newFakeBackupDRService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "compute.googleapis.com/Instance", result.ResourceType)
	require.Len(t, result.BackupRules, 1)
	assert.Equal(t, "terratest-daily", result.BackupRules[0].RuleId)
	assert.Equal(t, int64(5), result.BackupRules[0].BackupRetentionDays)
}

func TestGetBackupDRBackupVaultAttrsWithClientReportsAMissingOne(t *testing.T) {
	t.Parallel()

	// A caller who asks for something that is not there should be told that, rather than be handed the
	// transport's own wording for a 404.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error":{"code":404,"message":"Resource not found."}}`))
	})

	_, err := gcp.GetBackupDRBackupVaultAttrsWithClient(context.Background(), newFakeBackupDRService(t, handler), "gw-library-test-project", "us-central1", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "does not exist")
}
