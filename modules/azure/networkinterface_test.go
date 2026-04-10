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

func TestGetNetworkInterfaceE(t *testing.T) {
	t.Parallel()

	_, err := azure.GetNetworkInterfaceContextE(t.Context(), "", "", "")

	require.Error(t, err)
}

func TestGetNetworkInterfacePrivateIPsE(t *testing.T) {
	t.Parallel()

	_, err := azure.GetNetworkInterfacePrivateIPsContextE(t.Context(), "", "", "")

	require.Error(t, err)
}

func TestGetNetworkInterfacePublicIPsE(t *testing.T) {
	t.Parallel()

	_, err := azure.GetNetworkInterfacePublicIPsContextE(t.Context(), "", "", "")

	require.Error(t, err)
}

func TestNetworkInterfaceExistsE(t *testing.T) {
	t.Parallel()

	_, err := azure.NetworkInterfaceExistsContextE(t.Context(), "", "", "")

	require.Error(t, err)
}
