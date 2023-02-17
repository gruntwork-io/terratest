//go:build azure
// +build azure

// NOTE: We use build tags to differentiate azure testing because we currently do not have azure access setup for
// CircleCI.
package azure

import (
	"testing"

	"github.com/stretchr/testify/require"
)

/*
The below tests are currently stubbed out, with the expectation that they will throw errors. These tests can be extended.
*/
func TestListServiceBusNamespaceNamesE(t *testing.T) {
	t.Parallel()

	subscriptionID := ""

	_, err := ListServiceBusNamespaceNamesE(subscriptionID)
	require.Error(t, err)
}

func TestListServiceBusNamespaceIDsByResourceGroupE(t *testing.T) {
	t.Parallel()

	subscriptionID := ""
	resourceGroup := ""

	_, err := ListServiceBusNamespaceIDsByResourceGroupE(subscriptionID, resourceGroup)
	require.Error(t, err)
}

func TestListNamespaceAuthRulesE(t *testing.T) {
	t.Parallel()

	subscriptionID := ""
	namespace := ""
	resourceGroup := ""

	_, err := ListNamespaceAuthRulesE(subscriptionID, namespace, resourceGroup)
	require.Error(t, err)
}

func TestListNamespaceTopicsE(t *testing.T) {
	t.Parallel()

	subscriptionID := ""
	namespace := ""
	resourceGroup := ""

	_, err := ListNamespaceTopicsE(subscriptionID, namespace, resourceGroup)
	require.Error(t, err)
}

func TestListTopicAuthRulesE(t *testing.T) {
	t.Parallel()

	subscriptionID := ""
	namespace := ""
	resourceGroup := ""
	topicName := ""

	_, err := ListTopicAuthRulesE(subscriptionID, namespace, resourceGroup, topicName)
	require.Error(t, err)
}

func TestListTopicSubscriptionsNameE(t *testing.T) {
	t.Parallel()

	subscriptionID := ""
	namespace := ""
	resourceGroup := ""
	topicName := ""

	_, err := ListTopicSubscriptionsNameE(subscriptionID, namespace, resourceGroup, topicName)
	require.Error(t, err)
}
