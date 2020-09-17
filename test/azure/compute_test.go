// +build azure

// NOTE: We use build tags to differentiate azure testing because we currently do not have azure access setup for
// CircleCI.

package test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/azure"

	"github.com/stretchr/testify/require"
)

/*
The below tests are currently stubbed out, with the expectation that they will throw errors.
If/when CRUD methods are introduced for Azure Virtual Machines, these tests can be extended
(see AWS S3 tests for reference).
*/

func TestGetTagsForVirtualMachineE(t *testing.T) {
	t.Parallel()

	vmName := ""
	rgName := ""
	subID := ""

	_, err := azure.GetTagsForVirtualMachineE(t, vmName, rgName, subID)

	require.Error(t, err)
}

func TestGetSizeOfVirtualMachineE(t *testing.T) {
	t.Parallel()

	vmName := ""
	rgName := ""
	subID := ""

	_, err := azure.GetSizeOfVirtualMachineE(t, vmName, rgName, subID)

	require.Error(t, err)
}
