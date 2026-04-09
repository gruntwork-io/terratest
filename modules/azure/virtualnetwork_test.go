//go:build azure || (azureslim && network)
// +build azure azureslim,network

// NOTE: We use build tags to differentiate azure testing because we currently do not have azure access setup for
// CircleCI.

package azure

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

/*
The below tests are currently stubbed out, with the expectation that they will throw errors.
If/when methods can be mocked or Create/Delete APIs are added, these tests can be extended.
*/

func TestGetVirtualNetworkE(t *testing.T) {
	t.Parallel()

	_, err := GetVirtualNetworkContextE(context.Background(), "", "", "")

	require.Error(t, err)
}

func TestGetSubnetE(t *testing.T) {
	t.Parallel()

	_, err := GetSubnetContextE(context.Background(), "", "", "", "")

	require.Error(t, err)
}

func TestGetVirtualNetworkDNSServerIPsE(t *testing.T) {
	t.Parallel()

	_, err := GetVirtualNetworkDNSServerIPsContextE(context.Background(), "", "", "")

	require.Error(t, err)
}

func TestGetVirtualNetworkSubnetsE(t *testing.T) {
	t.Parallel()

	_, err := GetVirtualNetworkSubnetsContextE(context.Background(), "", "", "")

	require.Error(t, err)
}

func TestCheckSubnetContainsIPE(t *testing.T) {
	t.Parallel()

	_, err := CheckSubnetContainsIPContextE(context.Background(), "", "", "", "", "")

	require.Error(t, err)
}

func TestSubnetExistsE(t *testing.T) {
	t.Parallel()

	_, err := SubnetExistsContextE(context.Background(), "", "", "", "")

	require.Error(t, err)
}

func TestVirtualNetworkExistsE(t *testing.T) {
	t.Parallel()

	_, err := VirtualNetworkExistsContextE(context.Background(), "", "", "")

	require.Error(t, err)
}
