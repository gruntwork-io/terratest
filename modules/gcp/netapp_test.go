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

func TestGetNetAppActiveDirectoryAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a NetApp Active Directory connection the the library suite NetApp Active Directory connection module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/activeDirectories/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/activeDirectories/gw-library-test","domain":"terratest.example.com","dns":"198.51.100.1","netBiosPrefix":"terratest","state":"READY"}`))
	})

	attrs, err := gcp.GetNetAppActiveDirectoryAttrsWithClient(context.Background(), newFakeNetAppService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "terratest.example.com", attrs.Domain)
	assert.Equal(t, "198.51.100.1", attrs.Dns)
	assert.Equal(t, "terratest", attrs.NetBiosPrefix)
}

func TestGetNetAppActiveDirectoryAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a NetApp Active Directory connection that is not there should read a sentence about that NetApp Active Directory connection, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetNetAppActiveDirectoryAttrsWithClient(context.Background(), newFakeNetAppService(t, handler), "gw-library-test-project", "us-central1", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetNetAppStoragePoolAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a NetApp storage pool the the library suite NetApp storage pool module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/storagePools/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/storagePools/gw-library-test","serviceLevel":"PREMIUM","capacityGib":"2048","description":"created by terratest","labels":{"purpose":"terratest"}}`))
	})

	attrs, err := gcp.GetNetAppStoragePoolAttrsWithClient(context.Background(), newFakeNetAppService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "PREMIUM", attrs.ServiceLevel)
	assert.Equal(t, int64(2048), attrs.CapacityGib)
	assert.Equal(t, "created by terratest", attrs.Description)
}

func TestGetNetAppStoragePoolAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a NetApp storage pool that is not there should read a sentence about that NetApp storage pool, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetNetAppStoragePoolAttrsWithClient(context.Background(), newFakeNetAppService(t, handler), "gw-library-test-project", "us-central1", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetNetAppVolumeAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a NetApp volume the the library suite NetApp volume module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/volumes/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/volumes/gw-library-test","shareName":"terratest","capacityGib":"100","protocols":["NFSV3"],"description":"created by terratest","state":"READY"}`))
	})

	attrs, err := gcp.GetNetAppVolumeAttrsWithClient(context.Background(), newFakeNetAppService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "terratest", attrs.ShareName)
	assert.Equal(t, int64(100), attrs.CapacityGib)
	assert.Equal(t, "created by terratest", attrs.Description)
}

func TestGetNetAppVolumeAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a NetApp volume that is not there should read a sentence about that NetApp volume, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetNetAppVolumeAttrsWithClient(context.Background(), newFakeNetAppService(t, handler), "gw-library-test-project", "us-central1", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetNetAppVolumeSnapshotAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a NetApp volume snapshot the the library suite NetApp volume snapshot module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/volumes/gw-library-parent/snapshots/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/volumes/gw-library-parent/snapshots/gw-library-test","description":"created by terratest","state":"READY","labels":{"purpose":"terratest"}}`))
	})

	attrs, err := gcp.GetNetAppVolumeSnapshotAttrsWithClient(context.Background(), newFakeNetAppService(t, handler), "gw-library-test-project", "us-central1", "gw-library-parent", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "created by terratest", attrs.Description)
	assert.Equal(t, "READY", attrs.State)
}

func TestGetNetAppVolumeSnapshotAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a NetApp volume snapshot that is not there should read a sentence about that NetApp volume snapshot, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetNetAppVolumeSnapshotAttrsWithClient(context.Background(), newFakeNetAppService(t, handler), "gw-library-test-project", "us-central1", "gw-library-parent", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetNetAppVolumeReplicationAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a NetApp volume replication the the library suite NetApp volume replication module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/volumes/gw-library-parent/replications/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/volumes/gw-library-parent/replications/gw-library-test","replicationSchedule":"HOURLY","description":"created by terratest","mirrorState":"MIRRORED"}`))
	})

	attrs, err := gcp.GetNetAppVolumeReplicationAttrsWithClient(context.Background(), newFakeNetAppService(t, handler), "gw-library-test-project", "us-central1", "gw-library-parent", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "HOURLY", attrs.ReplicationSchedule)
	assert.Equal(t, "created by terratest", attrs.Description)
	assert.Equal(t, "MIRRORED", attrs.MirrorState)
}

func TestGetNetAppVolumeReplicationAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a NetApp volume replication that is not there should read a sentence about that NetApp volume replication, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetNetAppVolumeReplicationAttrsWithClient(context.Background(), newFakeNetAppService(t, handler), "gw-library-test-project", "us-central1", "gw-library-parent", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetNetAppVolumeQuotaRuleAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a NetApp volume quota rule the the library suite NetApp volume quota rule module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/volumes/gw-library-parent/quotaRules/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/volumes/gw-library-parent/quotaRules/gw-library-test","type":"DEFAULT_USER_QUOTA","diskLimitMib":1024,"description":"created by terratest","state":"READY"}`))
	})

	attrs, err := gcp.GetNetAppVolumeQuotaRuleAttrsWithClient(context.Background(), newFakeNetAppService(t, handler), "gw-library-test-project", "us-central1", "gw-library-parent", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "DEFAULT_USER_QUOTA", attrs.Type)
	assert.Equal(t, int64(1024), attrs.DiskLimitMib)
	assert.Equal(t, "created by terratest", attrs.Description)
}

func TestGetNetAppVolumeQuotaRuleAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a NetApp volume quota rule that is not there should read a sentence about that NetApp volume quota rule, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetNetAppVolumeQuotaRuleAttrsWithClient(context.Background(), newFakeNetAppService(t, handler), "gw-library-test-project", "us-central1", "gw-library-parent", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetNetAppBackupAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a NetApp backup the the library suite NetApp backup module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/backupVaults/gw-library-parent/backups/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/backupVaults/gw-library-parent/backups/gw-library-test","description":"created by terratest","sourceVolume":"projects/gw-library-test-project/locations/us-central1/volumes/gw-library-parent","state":"READY"}`))
	})

	attrs, err := gcp.GetNetAppBackupAttrsWithClient(context.Background(), newFakeNetAppService(t, handler), "gw-library-test-project", "us-central1", "gw-library-parent", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "created by terratest", attrs.Description)
	assert.Equal(t, "READY", attrs.State)
}

func TestGetNetAppBackupAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a NetApp backup that is not there should read a sentence about that NetApp backup, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetNetAppBackupAttrsWithClient(context.Background(), newFakeNetAppService(t, handler), "gw-library-test-project", "us-central1", "gw-library-parent", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetNetAppHostGroupAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a host group a Gruntwork module created, not a
	// copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/hostGroups/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"...","description":"created by terratest","type":"ISCSI_INITIATOR","osType":"LINUX","hosts":["iqn.1993-08.org.debian:01:terratest"],"labels":{"purpose":"terratest"}}`))
	})

	result, err := gcp.GetNetAppHostGroupAttrsWithClient(context.Background(), newFakeNetAppService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "ISCSI_INITIATOR", result.Type)
	assert.Equal(t, "LINUX", result.OsType)
	require.Len(t, result.Hosts, 1)
}
func TestGetNetAppKmsConfigAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a KMS config a Gruntwork module created, not a
	// copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/kmsConfigs/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"...","description":"created by terratest","cryptoKeyName":"projects/gw-library-test-project/locations/us-central1/keyRings/gw-library-test/cryptoKeys/gw-library-test","labels":{"purpose":"terratest"}}`))
	})

	result, err := gcp.GetNetAppKmsConfigAttrsWithClient(context.Background(), newFakeNetAppService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "created by terratest", result.Description)
	assert.Contains(t, result.CryptoKeyName, "cryptoKeys/gw-library-test")
}
func TestGetNetAppHostGroupAttrsWithClientReportsAMissingOne(t *testing.T) {
	t.Parallel()

	// A caller who asks for something that is not there should be told that, rather than be handed the
	// transport's own wording for a 404.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error":{"code":404,"message":"Resource not found."}}`))
	})

	_, err := gcp.GetNetAppHostGroupAttrsWithClient(context.Background(), newFakeNetAppService(t, handler), "gw-library-test-project", "us-central1", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "does not exist")
}
