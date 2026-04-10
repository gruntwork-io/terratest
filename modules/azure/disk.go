package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v6"
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

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
	_, err := GetDiskContextE(ctx, diskName, resGroupName, subscriptionID)
	if err != nil {
		if ResourceNotFoundErrorExists(err) {
			return false, nil
		}

		return false, err
	}

	return true, nil
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

	client, err := CreateDisksClientE(subscriptionID)
	if err != nil {
		return nil, err
	}

	resp, err := client.Get(ctx, resGroupName, diskName, nil)
	if err != nil {
		return nil, err
	}

	return &resp.Disk, nil
}
