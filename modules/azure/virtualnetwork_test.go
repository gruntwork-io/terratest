//go:build azure || (azureslim && network)
// +build azure azureslim,network

// NOTE: We use build tags to differentiate azure testing because we currently do not have azure access setup for
// CircleCI.

package azure_test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/azure"
	"github.com/stretchr/testify/require"
)

/*
The below tests are currently stubbed out, with the expectation that they will throw errors.
If/when methods can be mocked or Create/Delete APIs are added, these tests can be extended.
*/

func TestGetVirtualNetworkE(t *testing.T) {
	t.Parallel()

	_, err := azure.GetVirtualNetworkContextE(t.Context(), "", "", "")

	require.Error(t, err)
}

func TestGetSubnetE(t *testing.T) {
	t.Parallel()

	_, err := azure.GetSubnetContextE(t.Context(), "", "", "", "")

	require.Error(t, err)
}

func TestGetVirtualNetworkDNSServerIPsE(t *testing.T) {
	t.Parallel()

	_, err := azure.GetVirtualNetworkDNSServerIPsContextE(t.Context(), "", "", "")

	require.Error(t, err)
}

func TestGetVirtualNetworkSubnetsE(t *testing.T) {
	t.Parallel()

	_, err := azure.GetVirtualNetworkSubnetsContextE(t.Context(), "", "", "")

	require.Error(t, err)
}

func TestCheckSubnetContainsIPE(t *testing.T) {
	t.Parallel()

	_, err := azure.CheckSubnetContainsIPContextE(t.Context(), "", "", "", "", "")

	require.Error(t, err)
}

func TestSubnetExistsE(t *testing.T) {
	t.Parallel()

	_, err := azure.SubnetExistsContextE(t.Context(), "", "", "", "")

	require.Error(t, err)
}

func TestVirtualNetworkExistsE(t *testing.T) {
	t.Parallel()

	_, err := azure.VirtualNetworkExistsContextE(t.Context(), "", "", "")

	require.Error(t, err)
}
