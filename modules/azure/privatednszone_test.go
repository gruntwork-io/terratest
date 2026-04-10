//go:build azure
// +build azure

package azure_test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/azure"
	"github.com/stretchr/testify/require"
)

func TestPrivateDNSZoneExistsContextE(t *testing.T) {
	t.Parallel()

	zoneName := ""
	resourceGroupName := ""
	subscriptionID := ""

	exists, err := azure.PrivateDNSZoneExistsContextE(t.Context(), zoneName, resourceGroupName, subscriptionID)

	require.False(t, exists)
	require.Error(t, err)
}

func TestGetPrivateDNSZoneContextE(t *testing.T) {
	t.Parallel()

	zoneName := ""
	resGroupName := ""
	subscriptionID := ""

	_, err := azure.GetPrivateDNSZoneContextE(t.Context(), zoneName, resGroupName, subscriptionID)
	require.Error(t, err)
}
