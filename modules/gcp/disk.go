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
