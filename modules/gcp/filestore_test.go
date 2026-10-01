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
	"google.golang.org/api/file/v1"
	"google.golang.org/api/option"
)

// newFakeFilestoreService points a real Filestore client at a local test server, so the Google transport is
// exercised rather than a hand-written stand-in for a type we do not own.
func newFakeFilestoreService(t *testing.T, handler http.Handler) *file.Service {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	service, err := file.NewService(context.Background(),
		option.WithEndpoint(server.URL), option.WithoutAuthentication())
	require.NoError(t, err)

	return service
}

func TestGetFilestoreInstanceAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a instance the terraform-google-data-
	// storage module created, not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1-a/instances/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1-a/instances/gw-library-test","description":"created by terratest","tier":"BASIC_HDD","state":"READY","labels":{"purpose":"terratest"},"fileShares":[{"name":"share","capacityGb":"1024"}],"networks":[{"network":"gw-library-test","modes":["MODE_IPV4"]}]}`))
	})

	// Google sends a share's size as a JSON string and the Go client decodes it to an int64.
	instance, err := gcp.GetFilestoreInstanceAttrsWithClient(context.Background(), newFakeFilestoreService(t, handler), "gw-library-test-project", "us-central1-a", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "BASIC_HDD", instance.Tier)
	assert.Equal(t, "READY", instance.State)
	require.Len(t, instance.FileShares, 1)
	assert.Equal(t, int64(1024), instance.FileShares[0].CapacityGb)
	require.Len(t, instance.Networks, 1)
	assert.Equal(t, "gw-library-test", instance.Networks[0].Network)
}

func TestGetFilestoreInstanceAttrsWithClientMissingInstance(t *testing.T) {
	t.Parallel()

	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	// The error names the instance and everything that identifies it, and each of those is a value
	// no other part of the message contains, or its check could not fail.
	_, err := gcp.GetFilestoreInstanceAttrsWithClient(context.Background(), newFakeFilestoreService(t, handler), "gw-library-test-project", "us-central1-a", "gone")
	require.ErrorContains(t, err, "does not exist")
	require.ErrorContains(t, err, "gone")
	require.ErrorContains(t, err, "us-central1-a")
	require.ErrorContains(t, err, "gw-library-test-project")
}

func TestGetFilestoreSnapshotAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a snapshot the terraform-google-data-storage snapshot module created, not a copy of
	// any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1-a/instances/gw-library-test/snapshots/gw-library-snapshot"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1-a/instances/gw-library-test/snapshots/gw-library-snapshot","description":"created by terratest","state":"READY","labels":{"purpose":"terratest"}}`))
	})

	snapshot, err := gcp.GetFilestoreSnapshotAttrsWithClient(context.Background(), newFakeFilestoreService(t, handler), "gw-library-test-project", "us-central1-a", "gw-library-test", "gw-library-snapshot")
	require.NoError(t, err)

	assert.Equal(t, "created by terratest", snapshot.Description)
	assert.Equal(t, "READY", snapshot.State)
	assert.Equal(t, "terratest", snapshot.Labels["purpose"])
}

func TestGetFilestoreSnapshotAttrsWithClientMissingSnapshot(t *testing.T) {
	t.Parallel()

	// A caller who names a snapshot that is not there should read a sentence about that snapshot, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetFilestoreSnapshotAttrsWithClient(context.Background(), newFakeFilestoreService(t, handler), "gw-library-test-project", "us-central1-a", "gw-library-test", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetFilestoreBackupAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a backup the terraform-google-data-storage backup module created, not a copy of any
	// one fixture's values. A backup names a location rather than a zone, and carries the share it was
	// taken from.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/backups/gw-library-backup"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/backups/gw-library-backup","description":"created by terratest","sourceFileShare":"terratest_share","state":"READY","labels":{"purpose":"terratest"}}`))
	})

	backup, err := gcp.GetFilestoreBackupAttrsWithClient(context.Background(), newFakeFilestoreService(t, handler), "gw-library-test-project", "us-central1", "gw-library-backup")
	require.NoError(t, err)

	assert.Equal(t, "created by terratest", backup.Description)
	assert.Equal(t, "terratest_share", backup.SourceFileShare)
	assert.Equal(t, "terratest", backup.Labels["purpose"])
}

func TestGetFilestoreBackupAttrsWithClientMissingBackup(t *testing.T) {
	t.Parallel()

	// A caller who names a backup that is not there should read a sentence about that backup, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetFilestoreBackupAttrsWithClient(context.Background(), newFakeFilestoreService(t, handler), "gw-library-test-project", "us-central1", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}
