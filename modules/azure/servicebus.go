package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/servicebus/armservicebus/v2"
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

func serviceBusNamespaceClientE(subscriptionID string) (*armservicebus.NamespacesClient, error) {
	return CreateServiceBusNamespacesClientE(subscriptionID) //nolint:contextcheck
}

func serviceBusTopicClientE(subscriptionID string) (*armservicebus.TopicsClient, error) {
	return CreateServiceBusTopicsClientE(subscriptionID) //nolint:contextcheck
}

func serviceBusSubscriptionsClientE(subscriptionID string) (*armservicebus.SubscriptionsClient, error) {
	return CreateServiceBusSubscriptionsClientE(subscriptionID) //nolint:contextcheck
}

// ListServiceBusNamespaceContextE lists all SB namespaces in all resource groups in the given subscription ID.
// The ctx parameter supports cancellation and timeouts.
func ListServiceBusNamespaceContextE(ctx context.Context, subscriptionID string) ([]*armservicebus.SBNamespace, error) {
	nsClient, err := serviceBusNamespaceClientE(subscriptionID) //nolint:contextcheck
	if err != nil {
		return nil, err
	}

	return ListServiceBusNamespaceWithClient(ctx, nsClient)
}

// ListServiceBusNamespaceWithClient lists all SB namespaces using the provided NamespacesClient.
func ListServiceBusNamespaceWithClient(ctx context.Context, client *armservicebus.NamespacesClient) ([]*armservicebus.SBNamespace, error) {
	pager := client.NewListPager(nil)

	var results []*armservicebus.SBNamespace

	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		results = append(results, page.Value...)
	}

	return results, nil
}

// ListServiceBusNamespaceContext lists all SB namespaces in all resource groups in the given subscription ID.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func ListServiceBusNamespaceContext(t testing.TestingT, ctx context.Context, subscriptionID string) []*armservicebus.SBNamespace {
	t.Helper()

	results, err := ListServiceBusNamespaceContextE(ctx, subscriptionID)
	require.NoError(t, err)

	return results
}

// ListServiceBusNamespaceNamesContextE lists names of all SB namespaces in all resource groups in the given subscription ID.
// The ctx parameter supports cancellation and timeouts.
func ListServiceBusNamespaceNamesContextE(ctx context.Context, subscriptionID string) ([]string, error) {
	sbNamespace, err := ListServiceBusNamespaceContextE(ctx, subscriptionID)
	if err != nil {
		return nil, err
	}

	results := BuildNamespaceNamesList(sbNamespace)

	return results, nil
}

// BuildNamespaceNamesList is a helper method to build a namespace name list.
func BuildNamespaceNamesList(sbNamespace []*armservicebus.SBNamespace) []string {
	results := make([]string, 0, len(sbNamespace))

	for _, namespace := range sbNamespace {
		if namespace == nil || namespace.Name == nil {
			continue
		}

		results = append(results, *namespace.Name)
	}

	return results
}

// BuildNamespaceIdsList is a helper method to build a namespace id list.
func BuildNamespaceIdsList(sbNamespace []*armservicebus.SBNamespace) []string {
	results := make([]string, 0, len(sbNamespace))

	for _, namespace := range sbNamespace {
		if namespace == nil || namespace.ID == nil {
			continue
		}

		results = append(results, *namespace.ID)
	}

	return results
}

// ListServiceBusNamespaceNamesContext lists names of all SB namespaces in all resource groups in the given subscription ID.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func ListServiceBusNamespaceNamesContext(t testing.TestingT, ctx context.Context, subscriptionID string) []string {
	t.Helper()

	results, err := ListServiceBusNamespaceNamesContextE(ctx, subscriptionID)
	require.NoError(t, err)

	return results
}

// ListServiceBusNamespaceIDsContextE lists IDs of all SB namespaces in all resource groups in the given subscription ID.
// The ctx parameter supports cancellation and timeouts.
func ListServiceBusNamespaceIDsContextE(ctx context.Context, subscriptionID string) ([]string, error) {
	sbNamespace, err := ListServiceBusNamespaceContextE(ctx, subscriptionID)
	if err != nil {
		return nil, err
	}

	results := BuildNamespaceIdsList(sbNamespace)

	return results, nil
}

// ListServiceBusNamespaceIDsContext lists IDs of all SB namespaces in all resource groups in the given subscription ID.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func ListServiceBusNamespaceIDsContext(t testing.TestingT, ctx context.Context, subscriptionID string) []string {
	t.Helper()

	results, err := ListServiceBusNamespaceIDsContextE(ctx, subscriptionID)
	require.NoError(t, err)

	return results
}

// ListServiceBusNamespaceByResourceGroupContextE lists all SB namespaces in the given resource group.
// The ctx parameter supports cancellation and timeouts.
func ListServiceBusNamespaceByResourceGroupContextE(ctx context.Context, subscriptionID string, resourceGroup string) ([]*armservicebus.SBNamespace, error) {
	nsClient, err := serviceBusNamespaceClientE(subscriptionID) //nolint:contextcheck
	if err != nil {
		return nil, err
	}

	pager := nsClient.NewListByResourceGroupPager(resourceGroup, nil)

	var results []*armservicebus.SBNamespace

	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		results = append(results, page.Value...)
	}

	return results, nil
}

// ListServiceBusNamespaceByResourceGroupContext lists all SB namespaces in the given resource group.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func ListServiceBusNamespaceByResourceGroupContext(t testing.TestingT, ctx context.Context, subscriptionID string, resourceGroup string) []*armservicebus.SBNamespace {
	t.Helper()

	results, err := ListServiceBusNamespaceByResourceGroupContextE(ctx, subscriptionID, resourceGroup)
	require.NoError(t, err)

	return results
}

// ListServiceBusNamespaceNamesByResourceGroupContextE lists names of all SB namespaces in the given resource group.
// The ctx parameter supports cancellation and timeouts.
func ListServiceBusNamespaceNamesByResourceGroupContextE(ctx context.Context, subscriptionID string, resourceGroup string) ([]string, error) {
	sbNamespace, err := ListServiceBusNamespaceByResourceGroupContextE(ctx, subscriptionID, resourceGroup)
	if err != nil {
		return nil, err
	}

	results := BuildNamespaceNamesList(sbNamespace)

	return results, nil
}

// ListServiceBusNamespaceNamesByResourceGroupContext lists names of all SB namespaces in the given resource group.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func ListServiceBusNamespaceNamesByResourceGroupContext(t testing.TestingT, ctx context.Context, subscriptionID string, resourceGroup string) []string {
	t.Helper()

	results, err := ListServiceBusNamespaceNamesByResourceGroupContextE(ctx, subscriptionID, resourceGroup)
	require.NoError(t, err)

	return results
}

// ListServiceBusNamespaceIDsByResourceGroupContextE lists IDs of all SB namespaces in the given resource group.
// The ctx parameter supports cancellation and timeouts.
func ListServiceBusNamespaceIDsByResourceGroupContextE(ctx context.Context, subscriptionID string, resourceGroup string) ([]string, error) {
	sbNamespace, err := ListServiceBusNamespaceByResourceGroupContextE(ctx, subscriptionID, resourceGroup)
	if err != nil {
		return nil, err
	}

	results := BuildNamespaceIdsList(sbNamespace)

	return results, nil
}

// ListServiceBusNamespaceIDsByResourceGroupContext lists IDs of all SB namespaces in the given resource group.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func ListServiceBusNamespaceIDsByResourceGroupContext(t testing.TestingT, ctx context.Context, subscriptionID string, resourceGroup string) []string {
	t.Helper()

	results, err := ListServiceBusNamespaceIDsByResourceGroupContextE(ctx, subscriptionID, resourceGroup)
	require.NoError(t, err)

	return results
}

// ListNamespaceAuthRulesContextE authenticates the namespace client and enumerates all values to get a list
// of authorization rules for the given namespace name, automatically crossing page boundaries as required.
// The ctx parameter supports cancellation and timeouts.
func ListNamespaceAuthRulesContextE(ctx context.Context, subscriptionID string, namespace string, resourceGroup string) ([]string, error) {
	nsClient, err := serviceBusNamespaceClientE(subscriptionID) //nolint:contextcheck
	if err != nil {
		return nil, err
	}

	pager := nsClient.NewListAuthorizationRulesPager(resourceGroup, namespace, nil)

	var results []string

	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, rule := range page.Value {
			if rule == nil || rule.Name == nil {
				continue
			}

			results = append(results, *rule.Name)
		}
	}

	return results, nil
}

// ListNamespaceAuthRulesContext authenticates the namespace client and enumerates all values to get a list
// of authorization rules for the given namespace name, automatically crossing page boundaries as required.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func ListNamespaceAuthRulesContext(t testing.TestingT, ctx context.Context, subscriptionID string, namespace string, resourceGroup string) []string {
	t.Helper()

	results, err := ListNamespaceAuthRulesContextE(ctx, subscriptionID, namespace, resourceGroup)
	require.NoError(t, err)

	return results
}

// ListNamespaceTopicsContextE authenticates the topic client and enumerates all values,
// automatically crossing page boundaries as required.
// The ctx parameter supports cancellation and timeouts.
func ListNamespaceTopicsContextE(ctx context.Context, subscriptionID string, namespace string, resourceGroup string) ([]*armservicebus.SBTopic, error) {
	tClient, err := serviceBusTopicClientE(subscriptionID) //nolint:contextcheck
	if err != nil {
		return nil, err
	}

	pager := tClient.NewListByNamespacePager(resourceGroup, namespace, nil)

	var results []*armservicebus.SBTopic

	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		results = append(results, page.Value...)
	}

	return results, nil
}

// ListNamespaceTopicsContext authenticates the topic client and enumerates all values,
// automatically crossing page boundaries as required.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func ListNamespaceTopicsContext(t testing.TestingT, ctx context.Context, subscriptionID string, namespace string, resourceGroup string) []*armservicebus.SBTopic {
	t.Helper()

	results, err := ListNamespaceTopicsContextE(ctx, subscriptionID, namespace, resourceGroup)
	require.NoError(t, err)

	return results
}

// ListTopicSubscriptionsContextE authenticates the subscriptions client and enumerates all values,
// automatically crossing page boundaries as required.
// The ctx parameter supports cancellation and timeouts.
func ListTopicSubscriptionsContextE(ctx context.Context, subscriptionID string, namespace string, resourceGroup string, topicName string) ([]*armservicebus.SBSubscription, error) {
	sClient, err := serviceBusSubscriptionsClientE(subscriptionID) //nolint:contextcheck
	if err != nil {
		return nil, err
	}

	pager := sClient.NewListByTopicPager(resourceGroup, namespace, topicName, nil)

	var results []*armservicebus.SBSubscription

	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		results = append(results, page.Value...)
	}

	return results, nil
}

// ListTopicSubscriptionsContext authenticates the subscriptions client and enumerates all values,
// automatically crossing page boundaries as required.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func ListTopicSubscriptionsContext(t testing.TestingT, ctx context.Context, subscriptionID string, namespace string, resourceGroup string, topicName string) []*armservicebus.SBSubscription {
	t.Helper()

	results, err := ListTopicSubscriptionsContextE(ctx, subscriptionID, namespace, resourceGroup, topicName)
	require.NoError(t, err)

	return results
}

// ListTopicSubscriptionsNameContextE authenticates the subscriptions client and enumerates all values to get
// a list of subscriptions for the given topic name, automatically crossing page boundaries as required.
// The ctx parameter supports cancellation and timeouts.
func ListTopicSubscriptionsNameContextE(ctx context.Context, subscriptionID string, namespace string, resourceGroup string, topicName string) ([]string, error) {
	sClient, err := serviceBusSubscriptionsClientE(subscriptionID) //nolint:contextcheck
	if err != nil {
		return nil, err
	}

	pager := sClient.NewListByTopicPager(resourceGroup, namespace, topicName, nil)

	var results []string

	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, sub := range page.Value {
			if sub == nil || sub.Name == nil {
				continue
			}

			results = append(results, *sub.Name)
		}
	}

	return results, nil
}

// ListTopicSubscriptionsNameContext authenticates the subscriptions client and enumerates all values to get
// a list of subscriptions for the given topic name, automatically crossing page boundaries as required.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func ListTopicSubscriptionsNameContext(t testing.TestingT, ctx context.Context, subscriptionID string, namespace string, resourceGroup string, topicName string) []string {
	t.Helper()

	results, err := ListTopicSubscriptionsNameContextE(ctx, subscriptionID, namespace, resourceGroup, topicName)
	require.NoError(t, err)

	return results
}

// ListTopicAuthRulesContextE authenticates the topic client and enumerates all values to get a list
// of authorization rules for the given topic name, automatically crossing page boundaries as required.
// The ctx parameter supports cancellation and timeouts.
func ListTopicAuthRulesContextE(ctx context.Context, subscriptionID string, namespace string, resourceGroup string, topicName string) ([]string, error) {
	tClient, err := serviceBusTopicClientE(subscriptionID) //nolint:contextcheck
	if err != nil {
		return nil, err
	}

	pager := tClient.NewListAuthorizationRulesPager(resourceGroup, namespace, topicName, nil)

	var results []string

	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, rule := range page.Value {
			if rule == nil || rule.Name == nil {
				continue
			}

			results = append(results, *rule.Name)
		}
	}

	return results, nil
}

// ListTopicAuthRulesContext authenticates the topic client and enumerates all values to get a list
// of authorization rules for the given topic name, automatically crossing page boundaries as required.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func ListTopicAuthRulesContext(t testing.TestingT, ctx context.Context, subscriptionID string, namespace string, resourceGroup string, topicName string) []string {
	t.Helper()

	results, err := ListTopicAuthRulesContextE(ctx, subscriptionID, namespace, resourceGroup, topicName)
	require.NoError(t, err)

	return results
}
