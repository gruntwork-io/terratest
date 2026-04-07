package azure_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	azure "github.com/gruntwork-io/terratest/modules/azure"
)

/*
The below tests are currently stubbed out, with the expectation that they will throw errors.
If/when CRUD methods are introduced for Azure Synapse, these tests can be extended
*/
func TestPrivateDNSZoneExists(t *testing.T) {
	t.Parallel()

	zoneName := ""
	resourceGroupName := ""
	subscriptionID := ""

	exists, err := azure.PrivateDNSZoneExistsContextE(context.Background(), zoneName, resourceGroupName, subscriptionID)

	require.False(t, exists)
	require.Error(t, err)
}

func TestPrivateDNSZoneExistsE(t *testing.T) {
	t.Parallel()

	resGroupName := ""
	subscriptionID := ""
	zoneName := ""

	_, err := azure.GetPrivateDNSZoneContextE(context.Background(), subscriptionID, resGroupName, zoneName)
	require.Error(t, err)
}
