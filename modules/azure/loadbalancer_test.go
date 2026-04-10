//go:build azure
// +build azure

// NOTE: We use build tags to differentiate azure testing because we currently do not have azure access setup for
// CircleCI.

package azure_test

import (
	"context"
	"testing"

	"github.com/gruntwork-io/terratest/modules/azure"
	"github.com/stretchr/testify/require"
)

/*
The below tests are currently stubbed out, with the expectation that they will throw errors.
If/when methods can be mocked or Create/Delete APIs are added, these tests can be extended.
*/

func TestLoadBalancerExistsE(t *testing.T) {
	t.Parallel()

	loadBalancerName := ""
	resourceGroupName := ""
	subscriptionID := ""

	_, err := azure.LoadBalancerExistsContextE(context.Background(), loadBalancerName, resourceGroupName, subscriptionID)

	require.Error(t, err)
}

func TestGetLoadBalancerE(t *testing.T) {
	t.Parallel()

	loadBalancerName := ""
	resourceGroupName := ""
	subscriptionID := ""

	_, err := azure.GetLoadBalancerContextE(context.Background(), loadBalancerName, resourceGroupName, subscriptionID)

	require.Error(t, err)
}
