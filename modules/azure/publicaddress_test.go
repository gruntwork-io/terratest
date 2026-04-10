//go:build azure || (azureslim && network)
// +build azure azureslim,network

// NOTE: We use build tags to differentiate azure testing because we currently do not have azure access setup for
// CircleCI.

package azure

import (
	"testing"

	"github.com/stretchr/testify/require"
)

/*
The below tests are currently stubbed out, with the expectation that they will throw errors.
If/when methods can be mocked or Create/Delete APIs are added, these tests can be extended.
*/

func TestGetPublicIPAddressE(t *testing.T) {
	t.Parallel()

	paName := ""
	rgName := ""
	subID := ""

	_, err := GetPublicIPAddressContextE(t.Context(), paName, rgName, subID)

	require.Error(t, err)
}

func TestCheckPublicDNSNameAvailabilityE(t *testing.T) {
	t.Parallel()

	location := ""
	domain := ""
	subID := ""

	_, err := CheckPublicDNSNameAvailabilityContextE(t.Context(), location, domain, subID)

	require.Error(t, err)
}

func TestGetIPOfPublicIPAddressByNameE(t *testing.T) {
	t.Parallel()

	paName := ""
	rgName := ""
	subID := ""

	_, err := GetIPOfPublicIPAddressByNameContextE(t.Context(), paName, rgName, subID)

	require.Error(t, err)
}

func TestPublicAddressExistsE(t *testing.T) {
	t.Parallel()

	paName := ""
	rgName := ""
	subID := ""

	_, err := PublicAddressExistsContextE(t.Context(), paName, rgName, subID)

	require.Error(t, err)
}
