//go:build azure
// +build azure

// NOTE: We use build tags to differentiate azure testing because we currently do not have azure access setup for
// CircleCI.

package azure

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

/*
The below tests are currently stubbed out, with the expectation that they will throw errors.
If/when CRUD methods are introduced for Azure Virtual Machines, these tests can be extended.
*/

func TestGetVirtualMachineE(t *testing.T) {
	t.Parallel()

	vmName := ""
	rgName := ""
	subID := ""

	_, err := GetVirtualMachineContextE(context.Background(), vmName, rgName, subID)

	require.Error(t, err)
}

func TestListVirtualMachinesForResourceGroupE(t *testing.T) {
	t.Parallel()

	rgName := ""
	subID := ""

	_, err := ListVirtualMachinesForResourceGroupContextE(context.Background(), rgName, subID)

	require.Error(t, err)
}

func TestGetVirtualMachinesForResourceGroupE(t *testing.T) {
	t.Parallel()

	rgName := ""
	subID := ""

	_, err := GetVirtualMachinesForResourceGroupContextE(context.Background(), rgName, subID)

	require.Error(t, err)
}

func TestGetVirtualMachineTagsE(t *testing.T) {
	t.Parallel()

	vmName := ""
	rgName := ""
	subID := ""

	_, err := GetVirtualMachineTagsContextE(context.Background(), vmName, rgName, subID)

	require.Error(t, err)
}

func TestGetSizeOfVirtualMachineE(t *testing.T) {
	t.Parallel()

	vmName := ""
	rgName := ""
	subID := ""

	_, err := GetSizeOfVirtualMachineContextE(context.Background(), vmName, rgName, subID)

	require.Error(t, err)
}

func TestGetVirtualMachineImageE(t *testing.T) {
	t.Parallel()

	vmName := ""
	rgName := ""
	subID := ""

	_, err := GetVirtualMachineImageContextE(context.Background(), vmName, rgName, subID)

	require.Error(t, err)
}

func TestGetVirtualMachineAvailabilitySetIDE(t *testing.T) {
	t.Parallel()

	vmName := ""
	rgName := ""
	subID := ""

	_, err := GetVirtualMachineAvailabilitySetIDContextE(context.Background(), vmName, rgName, subID)

	require.Error(t, err)
}

func TestGetVirtualMachineOSDiskNameE(t *testing.T) {
	t.Parallel()

	vmName := ""
	rgName := ""
	subID := ""

	_, err := GetVirtualMachineOSDiskNameContextE(context.Background(), vmName, rgName, subID)

	require.Error(t, err)
}

func TestGetVirtualMachineManagedDisksE(t *testing.T) {
	t.Parallel()

	vmName := ""
	rgName := ""
	subID := ""

	_, err := GetVirtualMachineManagedDisksContextE(context.Background(), vmName, rgName, subID)

	require.Error(t, err)
}

func TestGetVirtualMachineNicsE(t *testing.T) {
	t.Parallel()

	vmName := ""
	rgName := ""
	subID := ""

	_, err := GetVirtualMachineNicsContextE(context.Background(), vmName, rgName, subID)

	require.Error(t, err)
}

func TestVirtualMachineExistsE(t *testing.T) {
	t.Parallel()

	vmName := ""
	rgName := ""
	subID := ""

	_, err := VirtualMachineExistsContextE(context.Background(), vmName, rgName, subID)

	require.Error(t, err)
}
