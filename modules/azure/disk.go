package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v6"
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// DiskExists indicates whether the specified Azure Managed Disk exists.
// This function would fail the test if there is an error.
//
// Deprecated: Use [DiskExistsContext] instead.
func DiskExists(t testing.TestingT, diskName string, resGroupName string, subscriptionID string) bool {
	t.Helper()

	return DiskExistsContext(t, context.Background(), diskName, resGroupName, subscriptionID)
}

// DiskExistsE indicates whether the specified Azure Managed Disk exists in the specified Azure Resource Group.
//
// Deprecated: Use [DiskExistsContextE] instead.
func DiskExistsE(diskName string, resGroupName string, subscriptionID string) (bool, error) {
	return DiskExistsContextE(context.Background(), diskName, resGroupName, subscriptionID)
}

// DiskExistsContext indicates whether the specified Azure Managed Disk exists.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func DiskExistsContext(t testing.TestingT, ctx context.Context, diskName string, resGroupName string, subscriptionID string) bool {
	t.Helper()

	exists, err := DiskExistsContextE(ctx, diskName, resGroupName, subscriptionID)
	require.NoError(t, err)

	return exists
}

// DiskExistsContextE indicates whether the specified Azure Managed Disk exists in the specified Azure Resource Group.
// The ctx parameter supports cancellation and timeouts.
func DiskExistsContextE(ctx context.Context, diskName string, resGroupName string, subscriptionID string) (bool, error) {
	resGroupName, err := getTargetAzureResourceGroupName(resGroupName)
	if err != nil {
		return false, err
	}

	client, err := CreateDisksClientContextE(ctx, subscriptionID)
	if err != nil {
		return false, err
	}

	return DiskExistsWithClient(ctx, client, resGroupName, diskName)
}

// GetDisk returns a Disk in the specified Azure Resource Group.
// This function would fail the test if there is an error.
//
// Deprecated: Use [GetDiskContext] instead.
func GetDisk(t testing.TestingT, diskName string, resGroupName string, subscriptionID string) *armcompute.Disk {
	t.Helper()

	return GetDiskContext(t, context.Background(), diskName, resGroupName, subscriptionID)
}

// GetDiskE returns a Disk in the specified Azure Resource Group.
//
// Deprecated: Use [GetDiskContextE] instead.
func GetDiskE(diskName string, resGroupName string, subscriptionID string) (*armcompute.Disk, error) {
	return GetDiskContextE(context.Background(), diskName, resGroupName, subscriptionID)
}

// GetDiskContext returns a Disk in the specified Azure Resource Group.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetDiskContext(t testing.TestingT, ctx context.Context, diskName string, resGroupName string, subscriptionID string) *armcompute.Disk {
	t.Helper()

	disk, err := GetDiskContextE(ctx, diskName, resGroupName, subscriptionID)
	require.NoError(t, err)

	return disk
}

// GetDiskContextE returns a Disk in the specified Azure Resource Group.
// The ctx parameter supports cancellation and timeouts.
func GetDiskContextE(ctx context.Context, diskName string, resGroupName string, subscriptionID string) (*armcompute.Disk, error) {
	resGroupName, err := getTargetAzureResourceGroupName(resGroupName)
	if err != nil {
		return nil, err
	}

	client, err := CreateDisksClientContextE(ctx, subscriptionID)
	if err != nil {
		return nil, err
	}

	return GetDiskWithClient(ctx, client, resGroupName, diskName)
}

// GetDiskWithClient returns a Disk using the provided DisksClient.
func GetDiskWithClient(ctx context.Context, client *armcompute.DisksClient, resGroupName string, diskName string) (*armcompute.Disk, error) {
	resp, err := client.Get(ctx, resGroupName, diskName, nil)
	if err != nil {
		return nil, err
	}

	return &resp.Disk, nil
}

// DiskExistsWithClient checks if a Disk exists using the provided DisksClient.
func DiskExistsWithClient(ctx context.Context, client *armcompute.DisksClient, resGroupName string, diskName string) (bool, error) {
	_, err := GetDiskWithClient(ctx, client, resGroupName, diskName)
	if err != nil {
		if ResourceNotFoundErrorExists(err) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}
