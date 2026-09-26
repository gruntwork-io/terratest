package gcp

import (
	"context"
	"fmt"

	"github.com/gruntwork-io/terratest/modules/core/v2/logger"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/compute/v1"
)

// FetchDiskContext queries GCP to return the settings it holds for the given zonal persistent
// disk, so a test can assert on what was actually created rather than only that it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func FetchDiskContext(t testing.TestingT, ctx context.Context, projectID string, zone string, name string) *compute.Disk {
	disk, err := FetchDiskContextE(t, ctx, projectID, zone, name)
	require.NoError(t, err)

	return disk
}

// FetchDiskContextE queries GCP to return the settings it holds for the given zonal persistent
// disk.
// The ctx parameter supports cancellation and timeouts.
func FetchDiskContextE(t testing.TestingT, ctx context.Context, projectID string, zone string, name string) (*compute.Disk, error) {
	logger.Default.Logf(t, "Getting disk %s in zone %s", name, zone)

	service, err := NewComputeServiceContextE(t, ctx)
	if err != nil {
		return nil, err
	}

	return FetchDiskWithClient(ctx, service, projectID, zone, name)
}

// FetchDiskWithClient queries GCP to return the settings it holds for the given zonal persistent
// disk using the supplied *compute.Service. Prefer this variant in unit tests where the service is
// backed by an httptest fake server (see disk_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func FetchDiskWithClient(ctx context.Context, service *compute.Service, projectID string, zone string, name string) (*compute.Disk, error) {
	disk, err := service.Disks.Get(projectID, zone, name).Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("Disks.Get(%s, %s, %s) got error: %w", projectID, zone, name, err)
	}

	return disk, nil
}

// FetchRegionDisk queries GCP to return the settings it holds for the given regional persistent disk, so a test can
// assert on what was actually created rather than only that it exists. A regional disk is replicated across two zones of one region, so it is read by region rather than by zone.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func FetchRegionDisk(t testing.TestingT, ctx context.Context, projectID string, region string, name string) *compute.Disk {
	disk, err := FetchRegionDiskE(t, ctx, projectID, region, name)
	require.NoError(t, err)

	return disk
}

// FetchRegionDiskE queries GCP to return the settings it holds for the given regional persistent disk.
// The ctx parameter supports cancellation and timeouts.
func FetchRegionDiskE(t testing.TestingT, ctx context.Context, projectID string, region string, name string) (*compute.Disk, error) {
	logger.Default.Logf(t, "Getting regional persistent disk %s in region %s", name, region)

	service, err := NewComputeServiceContextE(t, ctx)
	if err != nil {
		return nil, err
	}

	return FetchRegionDiskWithClient(ctx, service, projectID, region, name)
}

// FetchRegionDiskWithClient queries GCP to return the settings it holds for the given regional persistent disk using the
// supplied *compute.Service. Prefer this variant in unit tests where the service is backed by an
// httptest fake server (see disk_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func FetchRegionDiskWithClient(ctx context.Context, service *compute.Service, projectID string, region string, name string) (*compute.Disk, error) {
	disk, err := service.RegionDisks.Get(projectID, region, name).Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("RegionDisks.Get(%s, %s, %s) got error: %w", projectID, region, name, err)
	}

	return disk, nil
}

// FetchSnapshot queries GCP to return the settings it holds for the given snapshot, so a test can
// assert on what was actually created rather than only that it exists. A snapshot is global, so it is read by name alone.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func FetchSnapshot(t testing.TestingT, ctx context.Context, projectID string, name string) *compute.Snapshot {
	snapshot, err := FetchSnapshotE(t, ctx, projectID, name)
	require.NoError(t, err)

	return snapshot
}

// FetchSnapshotE queries GCP to return the settings it holds for the given snapshot.
// The ctx parameter supports cancellation and timeouts.
func FetchSnapshotE(t testing.TestingT, ctx context.Context, projectID string, name string) (*compute.Snapshot, error) {
	logger.Default.Logf(t, "Getting snapshot %s", name)

	service, err := NewComputeServiceContextE(t, ctx)
	if err != nil {
		return nil, err
	}

	return FetchSnapshotWithClient(ctx, service, projectID, name)
}

// FetchSnapshotWithClient queries GCP to return the settings it holds for the given snapshot using the
// supplied *compute.Service. Prefer this variant in unit tests where the service is backed by an
// httptest fake server (see disk_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func FetchSnapshotWithClient(ctx context.Context, service *compute.Service, projectID string, name string) (*compute.Snapshot, error) {
	snapshot, err := service.Snapshots.Get(projectID, name).Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("Snapshots.Get(%s, %s) got error: %w", projectID, name, err)
	}

	return snapshot, nil
}
