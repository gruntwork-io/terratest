//go:build azure
// +build azure

// NOTE: We use build tags to differentiate azure testing because we currently do not have azure access setup for
// CircleCI.
package azure

import (
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

/*
The below tests are currently stubbed out, with the expectation that they will throw errors.
If/when methods to create and delete resource groups are added, these tests can be extended.
*/

func TestResourceGroupExistsV2(t *testing.T) {
	t.Parallel()

	resourceGroupName := "fakeResourceGroupName"
	_, err := ResourceGroupExistsV2E(resourceGroupName, "")
	errAzure := &azcore.ResponseError{}
	require.ErrorAs(t, err, &errAzure)
	assert.Equal(t, errAzure.StatusCode, 404)
}

func TestGetAResourceGroupV2(t *testing.T) {
	t.Parallel()

	resourceGroupName := "fakeResourceGroupName"

	_, err := GetAResourceGroupV2E(resourceGroupName, "")
	errAzure := &azcore.ResponseError{}
	require.ErrorAs(t, err, &errAzure)
	assert.Equal(t, errAzure.StatusCode, 404)
}
