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

func TestFetchDiskWithClient(t *testing.T) {
	t.Parallel()

	// The values are the ones the terraform-google-compute disk module sets, because the point of
	// reading settings back is asserting a module configured the disk it was asked for.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/zones/us-central1-a/disks/gw-library-test"))
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"kind":"compute#disk",
			"name":"gw-library-test",
			"description":"created by terratest",
			"sizeGb":"10",
			"type":"https://www.googleapis.com/compute/v1/projects/gw-library-test-project/zones/us-central1-a/diskTypes/pd-standard",
			"zone":"https://www.googleapis.com/compute/v1/projects/gw-library-test-project/zones/us-central1-a",
			"status":"READY",
			"labels":{"managed-by":"terratest"}
		}`))
	})

	disk, err := gcp.FetchDiskWithClient(context.Background(), newFakeComputeService(t, handler), "gw-library-test-project", "us-central1-a", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "gw-library-test", disk.Name)
	assert.Equal(t, "created by terratest", disk.Description)
	// The API sends the size as a JSON string and the Go client decodes it to an int64.
	assert.Equal(t, int64(10), disk.SizeGb)
	assert.True(t, strings.HasSuffix(disk.Type, "/diskTypes/pd-standard"))
	assert.Equal(t, "READY", disk.Status)
	assert.Equal(t, "terratest", disk.Labels["managed-by"])
}

func TestFetchDiskWithClientMissingDisk(t *testing.T) {
	t.Parallel()

	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	// The error names the call, the project, the zone and the disk, in the shape the other compute
	// reads here use, so all of them are asserted rather than only that it failed.
	_, err := gcp.FetchDiskWithClient(context.Background(), newFakeComputeService(t, handler), "gw-library-test-project", "us-central1-a", "gone")
	require.ErrorContains(t, err, "Disks.Get(gw-library-test-project, us-central1-a, gone)")
}

func TestFetchRegionDiskWithClient(t *testing.T) {
	t.Parallel()

	// The values are the ones the terraform-google-compute region disk module sets, because the point of reading settings back
	// is asserting a module configured the disk it was asked for.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/regions/us-central1/disks/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"gw-library-test","sizeGb":"200","status":"READY","replicaZones":["https://www.googleapis.com/compute/v1/projects/gw-library-test-project/zones/us-central1-a","https://www.googleapis.com/compute/v1/projects/gw-library-test-project/zones/us-central1-b"],"labels":{"purpose":"terratest"}}`))
	})

	disk, err := gcp.FetchRegionDiskWithClient(context.Background(), newFakeComputeService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "gw-library-test", disk.Name)
	// The API sends the size as a JSON string and the Go client decodes it to an int64.
	assert.Equal(t, int64(200), disk.SizeGb)
	assert.Equal(t, "READY", disk.Status)
	require.Len(t, disk.ReplicaZones, 2, "a regional disk is replicated across exactly two zones")
	assert.Equal(t, "terratest", disk.Labels["purpose"])
}

func TestFetchSnapshotWithClient(t *testing.T) {
	t.Parallel()

	// The values are the ones the terraform-google-compute snapshot module sets, because the point of reading settings back
	// is asserting a module configured the snapshot it was asked for.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/global/snapshots/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"gw-library-test","description":"created by terratest","status":"READY","diskSizeGb":"10","storageLocations":["us-central1"],"labels":{"purpose":"terratest"}}`))
	})

	snapshot, err := gcp.FetchSnapshotWithClient(context.Background(), newFakeComputeService(t, handler), "gw-library-test-project", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "gw-library-test", snapshot.Name)
	assert.Equal(t, "created by terratest", snapshot.Description)
	assert.Equal(t, "READY", snapshot.Status)
	assert.Equal(t, int64(10), snapshot.DiskSizeGb)
	require.Len(t, snapshot.StorageLocations, 1)
	assert.Equal(t, "us-central1", snapshot.StorageLocations[0])
}
