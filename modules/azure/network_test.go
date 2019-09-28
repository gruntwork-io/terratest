// +build azure

// NOTE: We use build tags to differentiate azure testing because we currently do not have azure access setup for
// CircleCI.

package azure

import (
	"fmt"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	// subscriptionID is overridden by the environment variable "ARM_SUBSCRIPTION_ID"
	subscriptionID = ""
	rgName         = "terratest-rg"
)

func TestGetVirtualNetworkClient(t *testing.T) {
	t.Parallel()

	vnetClient, err := GetVirtualNetworkClient(subscriptionID)

	require.NoError(t, err)
	assert.NotEmpty(t, *vnetClient)
}

func TestGetPublicIPClient(t *testing.T) {
	t.Parallel()

	ipClient, err := GetPublicIPClient(subscriptionID)

	require.NoError(t, err)
	assert.NotEmpty(t, *ipClient)
}

func TestCheckPublicDNSNameAvailability(t *testing.T) {
	t.Parallel()

	randomsuffix := strings.ToLower(fmt.Sprintf("%s%s", random.UniqueId(), random.UniqueId()))
	nonExistentDomainNameLabel := fmt.Sprintf("nonexistent-%s", randomsuffix)
	location := GetRandomStableRegion(t, []string{}, []string{"australiacentral2"}, subscriptionID)

	available := CheckPublicDNSNameAvailability(t, location, nonExistentDomainNameLabel, subscriptionID)

	assert.True(t, available)
}

/*
The below tests are currently stubbed out, with the expectation that they will throw errors.
If/when methods to create and delete network resources are added, these tests can be extended.
(see AWS S3 tests for reference).
*/

func TestGetVirtualNetworkE(t *testing.T) {
	t.Parallel()

	vnetName := ""

	_, err := GetVirtualNetworkE(t, vnetName, rgName, subscriptionID)

	require.Error(t, err)
}

func TestGetSubnetsForVirtualNetworkE(t *testing.T) {
	t.Parallel()

	vnetName := ""

	_, err := GetSubnetsForVirtualNetworkE(t, vnetName, rgName, subscriptionID)

	require.Error(t, err)
}

func TestGetTagsForVirtualNetworkE(t *testing.T) {
	t.Parallel()

	vnetName := ""

	_, err := GetTagsForVirtualNetworkE(t, vnetName, rgName, subscriptionID)

	require.Error(t, err)
}

func TestGetNetworkSecurityGroupForSubnetE(t *testing.T) {
	t.Parallel()

	vnetName := ""
	subnetName := ""

	_, err := GetNetworkSecurityGroupForSubnetE(t, subnetName, vnetName, rgName, subscriptionID)

	require.Error(t, err)
}

func TestGetPublicIPE(t *testing.T) {
	t.Parallel()

	publicIPName := ""

	_, err := GetPublicIPE(t, rgName, publicIPName, subscriptionID)

	require.Error(t, err)
}
