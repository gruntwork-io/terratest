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
If/when methods to create and delete resource groups are added, these tests can be extended.
*/

func TestResourceGroupExists(t *testing.T) {
	t.Parallel()

	resourceGroupName := "fakeResourceGroupName"
	_, err := ResourceGroupExistsE(resourceGroupName, "")
	require.Error(t, err)
}

func TestGetAResourceGroup(t *testing.T) {
	t.Parallel()

	resourceGroupName := "fakeResourceGroupName"

	_, err := GetAResourceGroupE(resourceGroupName, "")
	require.Error(t, err)
}
