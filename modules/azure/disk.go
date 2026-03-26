package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2019-07-01/compute"
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

// DiskExists indicates whether the specified Azure Managed Disk exists.
// This function would fail the test if there is an error.
//
// Deprecated: Use [DiskExistsContext] instead.
func DiskExists(t testing.TestingT, diskName string, resGroupName string, subscriptionID string) bool {
	t.Helper()

	return DiskExistsContext(t, context.Background(), diskName, resGroupName, subscriptionID)
}

// DiskExistsContextE indicates whether the specified Azure Managed Disk exists in the specified Azure Resource Group.
// The ctx parameter supports cancellation and timeouts.
func DiskExistsContextE(ctx context.Context, diskName string, resGroupName string, subscriptionID string) (bool, error) {
	// Get the Disk object
	_, err := GetDiskContextE(ctx, diskName, resGroupName, subscriptionID)
	if err != nil {
		if ResourceNotFoundErrorExists(err) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

// DiskExistsE indicates whether the specified Azure Managed Disk exists in the specified Azure Resource Group.
//
// Deprecated: Use [DiskExistsContextE] instead.
func DiskExistsE(diskName string, resGroupName string, subscriptionID string) (bool, error) {
	return DiskExistsContextE(context.Background(), diskName, resGroupName, subscriptionID)
}

// GetDiskContext returns a Disk in the specified Azure Resource Group.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetDiskContext(t testing.TestingT, ctx context.Context, diskName string, resGroupName string, subscriptionID string) *compute.Disk {
	t.Helper()

	disk, err := GetDiskContextE(ctx, diskName, resGroupName, subscriptionID)
	require.NoError(t, err)

	return disk
}

// GetDisk returns a Disk in the specified Azure Resource Group.
// This function would fail the test if there is an error.
//
// Deprecated: Use [GetDiskContext] instead.
func GetDisk(t testing.TestingT, diskName string, resGroupName string, subscriptionID string) *compute.Disk {
	t.Helper()

	return GetDiskContext(t, context.Background(), diskName, resGroupName, subscriptionID)
}

// GetDiskContextE returns a Disk in the specified Azure Resource Group.
// The ctx parameter supports cancellation and timeouts.
func GetDiskContextE(ctx context.Context, diskName string, resGroupName string, subscriptionID string) (*compute.Disk, error) {
	// Validate resource group name and subscription ID
	resGroupName, err := getTargetAzureResourceGroupName(resGroupName)
	if err != nil {
		return nil, err
	}

	// Get the client reference
	client, err := CreateDisksClientE(subscriptionID)
	if err != nil {
		return nil, err
	}

	// Get the Disk
	disk, err := client.Get(ctx, resGroupName, diskName)
	if err != nil {
		return nil, err
	}

	return &disk, nil
}

// GetDiskE returns a Disk in the specified Azure Resource Group.
//
// Deprecated: Use [GetDiskContextE] instead.
func GetDiskE(diskName string, resGroupName string, subscriptionID string) (*compute.Disk, error) {
	return GetDiskContextE(context.Background(), diskName, resGroupName, subscriptionID)
}

// GetDiskClientE returns a new Disk client in the specified Azure Subscription.
// TODO: remove in next major/minor version
func GetDiskClientE(subscriptionID string) (*compute.DisksClient, error) {
	// Validate Azure subscription ID
	subscriptionID, err := getTargetAzureSubscription(subscriptionID)
	if err != nil {
		return nil, err
	}

	// Get the Disk client
	client := compute.NewDisksClient(subscriptionID)

	// Create an authorizer
	authorizer, err := NewAuthorizer()
	if err != nil {
		return nil, err
	}

	client.Authorizer = *authorizer

	return &client, nil
}
