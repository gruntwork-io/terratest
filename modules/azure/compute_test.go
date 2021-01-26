// +build azure

// NOTE: We use build tags to differentiate azure testing because we currently do not have azure access setup for
// CircleCI.

package azure

import (
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

	_, err := GetVirtualMachineE(vmName, rgName, subID)

	require.Error(t, err)
}

func TestListVirtualMachinesForResourceGroupE(t *testing.T) {
	t.Parallel()

	rgName := ""
	subID := ""

	_, err := ListVirtualMachinesForResourceGroupE(rgName, subID)

	require.Error(t, err)
}

func TestGetVirtualMachinesForResourceGroupE(t *testing.T) {
	t.Parallel()

	rgName := ""
	subID := ""

	_, err := GetVirtualMachinesForResourceGroupE(rgName, subID)

	require.Error(t, err)
}

func TestGetVirtualMachineTagsE(t *testing.T) {
	t.Parallel()

	vmName := ""
	rgName := ""
	subID := ""

	_, err := GetVirtualMachineTagsE(vmName, rgName, subID)

	require.Error(t, err)
}

func TestGetSizeOfVirtualMachineE(t *testing.T) {
	t.Parallel()

	vmName := ""
	rgName := ""
	subID := ""

	_, err := GetSizeOfVirtualMachineE(vmName, rgName, subID)

	require.Error(t, err)
}

func TestGetVirtualMachineImageE(t *testing.T) {
	t.Parallel()

	vmName := ""
	rgName := ""
	subID := ""

	_, err := GetVirtualMachineImageE(vmName, rgName, subID)

	require.Error(t, err)
}

func TestGetVirtualMachineAvailabilitySetIDE(t *testing.T) {
	t.Parallel()

	vmName := ""
	rgName := ""
	subID := ""

	_, err := GetVirtualMachineAvailabilitySetIDE(vmName, rgName, subID)

	require.Error(t, err)
}

func TestGetVirtualMachineOSDiskNameE(t *testing.T) {
	t.Parallel()

	vmName := ""
	rgName := ""
	subID := ""

	_, err := GetVirtualMachineOSDiskNameE(vmName, rgName, subID)

	require.Error(t, err)
}

func TestGetVirtualMachineManagedDisksE(t *testing.T) {
	t.Parallel()

	vmName := ""
	rgName := ""
	subID := ""

	_, err := GetVirtualMachineManagedDisksE(vmName, rgName, subID)

	require.Error(t, err)
}

func TestGetVirtualMachineNicsE(t *testing.T) {
	t.Parallel()

	vmName := ""
	rgName := ""
	subID := ""

	_, err := GetVirtualMachineNicsE(vmName, rgName, subID)

	require.Error(t, err)
}

func TestVirtualMachineExistsE(t *testing.T) {
	t.Parallel()

	vmName := ""
	rgName := ""
	subID := ""

	_, err := VirtualMachineExistsE(vmName, rgName, subID)

	require.Error(t, err)
}
