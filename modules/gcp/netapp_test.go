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
	"google.golang.org/api/netapp/v1"
	"google.golang.org/api/option"
)

// newFakeNetAppService points a real NetApp Volumes client at a local test server, so the Google transport
// is exercised rather than a hand-written stand-in for a type we do not own.
func newFakeNetAppService(t *testing.T, handler http.Handler) *netapp.Service {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	service, err := netapp.NewService(context.Background(),
		option.WithEndpoint(server.URL), option.WithoutAuthentication())
	require.NoError(t, err)

	return service
}

func TestGetNetAppBackupVaultAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a backup vault a Gruntwork module created, not a
	// copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/backupVaults/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/backupVaults/gw-library-test","description":"created by terratest","backupVaultType":"IN_REGION","labels":{"purpose":"terratest"},"state":"READY"}`))
	})

	result, err := gcp.GetNetAppBackupVaultAttrsWithClient(context.Background(), newFakeNetAppService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "created by terratest", result.Description)
	assert.Equal(t, "IN_REGION", result.BackupVaultType)
	assert.Equal(t, map[string]string{"purpose": "terratest"}, result.Labels)
}

func TestGetNetAppBackupPolicyAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a backup policy a Gruntwork module created, not a
	// copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/backupPolicies/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/backupPolicies/gw-library-test","description":"created by terratest","enabled":false,"dailyBackupLimit":3,"weeklyBackupLimit":2,"monthlyBackupLimit":1}`))
	})

	result, err := gcp.GetNetAppBackupPolicyAttrsWithClient(context.Background(), newFakeNetAppService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, int64(3), result.DailyBackupLimit)
	assert.Equal(t, int64(2), result.WeeklyBackupLimit)
	assert.False(t, result.Enabled)
}

func TestGetNetAppBackupVaultAttrsWithClientReportsAMissingOne(t *testing.T) {
	t.Parallel()

	// A caller who asks for something that is not there should be told that, rather than be handed the
	// transport's own wording for a 404.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error":{"code":404,"message":"Resource not found."}}`))
	})

	_, err := gcp.GetNetAppBackupVaultAttrsWithClient(context.Background(), newFakeNetAppService(t, handler), "gw-library-test-project", "us-central1", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "does not exist")
}
