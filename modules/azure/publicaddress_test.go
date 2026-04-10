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

func TestGetPublicIPAddressE(t *testing.T) {
	t.Parallel()

	_, err := GetPublicIPAddressContextE(context.Background(), "", "", "")

	require.Error(t, err)
}

func TestCheckPublicDNSNameAvailabilityE(t *testing.T) {
	t.Parallel()

	_, err := CheckPublicDNSNameAvailabilityContextE(context.Background(), "", "", "")

	require.Error(t, err)
}

func TestGetIPOfPublicIPAddressByNameE(t *testing.T) {
	t.Parallel()

	_, err := GetIPOfPublicIPAddressByNameContextE(context.Background(), "", "", "")

	require.Error(t, err)
}

func TestPublicAddressExistsE(t *testing.T) {
	t.Parallel()

	_, err := PublicAddressExistsContextE(context.Background(), "", "", "")

	require.Error(t, err)
}
