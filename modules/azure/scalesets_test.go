// +build azure

// NOTE: We use build tags to differentiate azure testing because we currently do not have azure access setup for
// CircleCI.

package azure

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

/*
The below tests are currently stubbed out, with the expectation that they will throw errors.
If/when methods to create and delete network resources are added, these tests can be extended.
*/

func TestGetVirtualMachineScaleSetsClientE(t *testing.T) {
	t.Parallel()

	subscriptionID := ""

	client, err := GetVirtualMachineScaleSetsClientE(subscriptionID)

	require.NoError(t, err)
	assert.NotEmpty(t, *client)
}

func TestVirtualMachineScaleSetExistsE(t *testing.T) {
	t.Parallel()

	vmssName := ""
	rgName := ""
	subID := ""

	_, err := VirtualMachineScaleSetExistsE(vmssName, rgName, subID)

	require.Error(t, err)
}

func TestGetVirtualMachineScaleSetE(t *testing.T) {
	t.Parallel()

	vmssName := ""
	rgName := ""
	subID := ""

	_, err := GetVirtualMachineScaleSetE(vmssName, rgName, subID)

	require.Error(t, err)
}
