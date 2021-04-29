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

func TestCreateAvailabilitySetClientE(t *testing.T) {
	t.Parallel()

	subscriptionID := ""

	client, err := CreateAvailabilitySetClientE(subscriptionID)

	require.NoError(t, err)
	assert.NotEmpty(t, *client)
}

func TestGetAvailabilitySetE(t *testing.T) {
	t.Parallel()

	avsName := ""
	rgName := ""
	subscriptionID := ""

	_, err := GetAvailabilitySetE(t, avsName, rgName, subscriptionID)

	require.Error(t, err)
}

func TestCheckAvailabilitySetContainsVME(t *testing.T) {
	t.Parallel()

	vmName := ""
	avsName := ""
	rgName := ""
	subscriptionID := ""

	_, err := CheckAvailabilitySetContainsVME(t, vmName, avsName, rgName, subscriptionID)

	require.Error(t, err)
}

func TestGetAvailabilitySetVMNamesInCapsE(t *testing.T) {
	t.Parallel()

	avsName := ""
	rgName := ""
	subscriptionID := ""

	_, err := GetAvailabilitySetVMNamesInCapsE(t, avsName, rgName, subscriptionID)

	require.Error(t, err)
}

func TestGetAvailabilitySetFaultDomainCountE(t *testing.T) {
	t.Parallel()

	avsName := ""
	rgName := ""
	subscriptionID := ""

	_, err := GetAvailabilitySetFaultDomainCountE(t, avsName, rgName, subscriptionID)

	require.Error(t, err)
}

func TestAvailabilitySetExistsE(t *testing.T) {
	t.Parallel()

	avsName := ""
	rgName := ""
	subscriptionID := ""

	_, err := AvailabilitySetExistsE(t, avsName, rgName, subscriptionID)

	require.Error(t, err)
}
