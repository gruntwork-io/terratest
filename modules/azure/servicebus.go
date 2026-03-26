package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/services/servicebus/mgmt/2017-04-01/servicebus"
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

func serviceBusNamespaceClientE(subscriptionID string) (*servicebus.NamespacesClient, error) {
	authorizer, err := NewAuthorizer()
	if err != nil {
		return nil, err
	}

	nsClient := servicebus.NewNamespacesClient(subscriptionID)
	nsClient.Authorizer = *authorizer

	return &nsClient, nil
}

func serviceBusTopicClientE(subscriptionID string) (*servicebus.TopicsClient, error) {
	authorizer, err := NewAuthorizer()
	if err != nil {
		return nil, err
	}

	tClient := servicebus.NewTopicsClient(subscriptionID)
	tClient.Authorizer = *authorizer

	return &tClient, nil
}

func serviceBusSubscriptionsClientE(subscriptionID string) (*servicebus.SubscriptionsClient, error) {
	authorizer, err := NewAuthorizer()
	if err != nil {
		return nil, err
	}

	sClient := servicebus.NewSubscriptionsClient(subscriptionID)
	sClient.Authorizer = *authorizer

	return &sClient, nil
}

// ListServiceBusNamespaceContextE lists all SB namespaces in all resource groups in the given subscription ID.
// The ctx parameter supports cancellation and timeouts.
func ListServiceBusNamespaceContextE(ctx context.Context, subscriptionID string) ([]servicebus.SBNamespace, error) {
	nsClient, err := serviceBusNamespaceClientE(subscriptionID)
	if err != nil {
		return nil, err
	}

	iteratorSBNamespace, err := nsClient.ListComplete(ctx)
	if err != nil {
		return nil, err
	}

	results := make([]servicebus.SBNamespace, 0)

	for iteratorSBNamespace.NotDone() {
		results = append(results, iteratorSBNamespace.Value())

		if err := iteratorSBNamespace.Next(); err != nil {
			return nil, err
		}
	}

	return results, nil
}

// ListServiceBusNamespaceE lists all SB namespaces in all resource groups in the given subscription ID.
//
// Deprecated: Use [ListServiceBusNamespaceContextE] instead.
func ListServiceBusNamespaceE(subscriptionID string) ([]servicebus.SBNamespace, error) {
	return ListServiceBusNamespaceContextE(context.Background(), subscriptionID)
}

// ListServiceBusNamespaceContext lists all SB namespaces in all resource groups in the given subscription ID.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func ListServiceBusNamespaceContext(t testing.TestingT, ctx context.Context, subscriptionID string) []servicebus.SBNamespace {
	t.Helper()

	results, err := ListServiceBusNamespaceContextE(ctx, subscriptionID)
	require.NoError(t, err)

	return results
}

// ListServiceBusNamespace lists all SB namespaces in all resource groups in the given subscription ID.
// This function would fail the test if there is an error.
//
// Deprecated: Use [ListServiceBusNamespaceContext] instead.
func ListServiceBusNamespace(t testing.TestingT, subscriptionID string) []servicebus.SBNamespace {
	t.Helper()

	return ListServiceBusNamespaceContext(t, context.Background(), subscriptionID)
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

// ListServiceBusNamespaceNamesE lists names of all SB namespaces in all resource groups in the given subscription ID.
//
// Deprecated: Use [ListServiceBusNamespaceNamesContextE] instead.
func ListServiceBusNamespaceNamesE(subscriptionID string) ([]string, error) {
	return ListServiceBusNamespaceNamesContextE(context.Background(), subscriptionID)
}

// BuildNamespaceNamesList is a helper method to build a namespace name list.
func BuildNamespaceNamesList(sbNamespace []servicebus.SBNamespace) []string {
	results := make([]string, 0, len(sbNamespace))

	for _, namespace := range sbNamespace {
		results = append(results, *namespace.Name)
	}

	return results
}

// BuildNamespaceIdsList is a helper method to build a namespace id list.
func BuildNamespaceIdsList(sbNamespace []servicebus.SBNamespace) []string {
	results := make([]string, 0, len(sbNamespace))

	for _, namespace := range sbNamespace {
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

// ListServiceBusNamespaceNames lists names of all SB namespaces in all resource groups in the given subscription ID.
// This function would fail the test if there is an error.
//
// Deprecated: Use [ListServiceBusNamespaceNamesContext] instead.
func ListServiceBusNamespaceNames(t testing.TestingT, subscriptionID string) []string {
	t.Helper()

	return ListServiceBusNamespaceNamesContext(t, context.Background(), subscriptionID)
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

// ListServiceBusNamespaceIDsE lists IDs of all SB namespaces in all resource groups in the given subscription ID.
//
// Deprecated: Use [ListServiceBusNamespaceIDsContextE] instead.
func ListServiceBusNamespaceIDsE(subscriptionID string) ([]string, error) {
	return ListServiceBusNamespaceIDsContextE(context.Background(), subscriptionID)
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

// ListServiceBusNamespaceIDs lists IDs of all SB namespaces in all resource groups in the given subscription ID.
// This function would fail the test if there is an error.
//
// Deprecated: Use [ListServiceBusNamespaceIDsContext] instead.
func ListServiceBusNamespaceIDs(t testing.TestingT, subscriptionID string) []string {
	t.Helper()

	return ListServiceBusNamespaceIDsContext(t, context.Background(), subscriptionID)
}

// ListServiceBusNamespaceByResourceGroupContextE lists all SB namespaces in the given resource group.
// The ctx parameter supports cancellation and timeouts.
func ListServiceBusNamespaceByResourceGroupContextE(ctx context.Context, subscriptionID string, resourceGroup string) ([]servicebus.SBNamespace, error) {
	nsClient, err := serviceBusNamespaceClientE(subscriptionID)
	if err != nil {
		return nil, err
	}

	iteratorSBNamespace, err := nsClient.ListByResourceGroupComplete(ctx, resourceGroup)
	if err != nil {
		return nil, err
	}

	results := make([]servicebus.SBNamespace, 0)

	for iteratorSBNamespace.NotDone() {
		results = append(results, iteratorSBNamespace.Value())

		if err := iteratorSBNamespace.Next(); err != nil {
			return nil, err
		}
	}

	return results, nil
}

// ListServiceBusNamespaceByResourceGroupE lists all SB namespaces in the given resource group.
//
// Deprecated: Use [ListServiceBusNamespaceByResourceGroupContextE] instead.
func ListServiceBusNamespaceByResourceGroupE(subscriptionID string, resourceGroup string) ([]servicebus.SBNamespace, error) {
	return ListServiceBusNamespaceByResourceGroupContextE(context.Background(), subscriptionID, resourceGroup)
}

// ListServiceBusNamespaceByResourceGroupContext lists all SB namespaces in the given resource group.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func ListServiceBusNamespaceByResourceGroupContext(t testing.TestingT, ctx context.Context, subscriptionID string, resourceGroup string) []servicebus.SBNamespace {
	t.Helper()

	results, err := ListServiceBusNamespaceByResourceGroupContextE(ctx, subscriptionID, resourceGroup)
	require.NoError(t, err)

	return results
}

// ListServiceBusNamespaceByResourceGroup lists all SB namespaces in the given resource group.
// This function would fail the test if there is an error.
//
// Deprecated: Use [ListServiceBusNamespaceByResourceGroupContext] instead.
func ListServiceBusNamespaceByResourceGroup(t testing.TestingT, subscriptionID string, resourceGroup string) []servicebus.SBNamespace {
	t.Helper()

	return ListServiceBusNamespaceByResourceGroupContext(t, context.Background(), subscriptionID, resourceGroup)
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

// ListServiceBusNamespaceNamesByResourceGroupE lists names of all SB namespaces in the given resource group.
//
// Deprecated: Use [ListServiceBusNamespaceNamesByResourceGroupContextE] instead.
func ListServiceBusNamespaceNamesByResourceGroupE(subscriptionID string, resourceGroup string) ([]string, error) {
	return ListServiceBusNamespaceNamesByResourceGroupContextE(context.Background(), subscriptionID, resourceGroup)
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

// ListServiceBusNamespaceNamesByResourceGroup lists names of all SB namespaces in the given resource group.
// This function would fail the test if there is an error.
//
// Deprecated: Use [ListServiceBusNamespaceNamesByResourceGroupContext] instead.
func ListServiceBusNamespaceNamesByResourceGroup(t testing.TestingT, subscriptionID string, resourceGroup string) []string {
	t.Helper()

	return ListServiceBusNamespaceNamesByResourceGroupContext(t, context.Background(), subscriptionID, resourceGroup)
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

// ListServiceBusNamespaceIDsByResourceGroupE lists IDs of all SB namespaces in the given resource group.
//
// Deprecated: Use [ListServiceBusNamespaceIDsByResourceGroupContextE] instead.
func ListServiceBusNamespaceIDsByResourceGroupE(subscriptionID string, resourceGroup string) ([]string, error) {
	return ListServiceBusNamespaceIDsByResourceGroupContextE(context.Background(), subscriptionID, resourceGroup)
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

// ListServiceBusNamespaceIDsByResourceGroup lists IDs of all SB namespaces in the given resource group.
// This function would fail the test if there is an error.
//
// Deprecated: Use [ListServiceBusNamespaceIDsByResourceGroupContext] instead.
func ListServiceBusNamespaceIDsByResourceGroup(t testing.TestingT, subscriptionID string, resourceGroup string) []string {
	t.Helper()

	return ListServiceBusNamespaceIDsByResourceGroupContext(t, context.Background(), subscriptionID, resourceGroup)
}

// ListNamespaceAuthRulesContextE authenticates the namespace client and enumerates all values to get a list
// of authorization rules for the given namespace name, automatically crossing page boundaries as required.
// The ctx parameter supports cancellation and timeouts.
func ListNamespaceAuthRulesContextE(ctx context.Context, subscriptionID string, namespace string, resourceGroup string) ([]string, error) {
	nsClient, err := serviceBusNamespaceClientE(subscriptionID)
	if err != nil {
		return nil, err
	}

	iteratorNamespaceRules, err := nsClient.ListAuthorizationRulesComplete(ctx, resourceGroup, namespace)
	if err != nil {
		return nil, err
	}

	results := []string{}

	for iteratorNamespaceRules.NotDone() {
		results = append(results, *(iteratorNamespaceRules.Value()).Name)

		if err := iteratorNamespaceRules.Next(); err != nil {
			return nil, err
		}
	}

	return results, nil
}

// ListNamespaceAuthRulesE authenticates the namespace client and enumerates all values to get a list
// of authorization rules for the given namespace name, automatically crossing page boundaries as required.
//
// Deprecated: Use [ListNamespaceAuthRulesContextE] instead.
func ListNamespaceAuthRulesE(subscriptionID string, namespace string, resourceGroup string) ([]string, error) {
	return ListNamespaceAuthRulesContextE(context.Background(), subscriptionID, namespace, resourceGroup)
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

// ListNamespaceAuthRules authenticates the namespace client and enumerates all values to get a list
// of authorization rules for the given namespace name, automatically crossing page boundaries as required.
// This function would fail the test if there is an error.
//
// Deprecated: Use [ListNamespaceAuthRulesContext] instead.
func ListNamespaceAuthRules(t testing.TestingT, subscriptionID string, namespace string, resourceGroup string) []string {
	t.Helper()

	return ListNamespaceAuthRulesContext(t, context.Background(), subscriptionID, namespace, resourceGroup)
}

// ListNamespaceTopicsContextE authenticates the topic client and enumerates all values,
// automatically crossing page boundaries as required.
// The ctx parameter supports cancellation and timeouts.
func ListNamespaceTopicsContextE(ctx context.Context, subscriptionID string, namespace string, resourceGroup string) ([]servicebus.SBTopic, error) {
	tClient, err := serviceBusTopicClientE(subscriptionID)
	if err != nil {
		return nil, err
	}

	iteratorTopics, err := tClient.ListByNamespaceComplete(ctx, resourceGroup, namespace, nil, nil)
	if err != nil {
		return nil, err
	}

	results := make([]servicebus.SBTopic, 0)

	for iteratorTopics.NotDone() {
		results = append(results, iteratorTopics.Value())

		if err := iteratorTopics.Next(); err != nil {
			return nil, err
		}
	}

	return results, nil
}

// ListNamespaceTopicsE authenticates the topic client and enumerates all values,
// automatically crossing page boundaries as required.
//
// Deprecated: Use [ListNamespaceTopicsContextE] instead.
func ListNamespaceTopicsE(subscriptionID string, namespace string, resourceGroup string) ([]servicebus.SBTopic, error) {
	return ListNamespaceTopicsContextE(context.Background(), subscriptionID, namespace, resourceGroup)
}

// ListNamespaceTopicsContext authenticates the topic client and enumerates all values,
// automatically crossing page boundaries as required.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func ListNamespaceTopicsContext(t testing.TestingT, ctx context.Context, subscriptionID string, namespace string, resourceGroup string) []servicebus.SBTopic {
	t.Helper()

	results, err := ListNamespaceTopicsContextE(ctx, subscriptionID, namespace, resourceGroup)
	require.NoError(t, err)

	return results
}

// ListNamespaceTopics authenticates the topic client and enumerates all values,
// automatically crossing page boundaries as required.
// This function would fail the test if there is an error.
//
// Deprecated: Use [ListNamespaceTopicsContext] instead.
func ListNamespaceTopics(t testing.TestingT, subscriptionID string, namespace string, resourceGroup string) []servicebus.SBTopic {
	t.Helper()

	return ListNamespaceTopicsContext(t, context.Background(), subscriptionID, namespace, resourceGroup)
}

// ListTopicSubscriptionsContextE authenticates the subscriptions client and enumerates all values,
// automatically crossing page boundaries as required.
// The ctx parameter supports cancellation and timeouts.
func ListTopicSubscriptionsContextE(ctx context.Context, subscriptionID string, namespace string, resourceGroup string, topicName string) ([]servicebus.SBSubscription, error) {
	sClient, err := serviceBusSubscriptionsClientE(subscriptionID)
	if err != nil {
		return nil, err
	}

	iteratorSubscription, err := sClient.ListByTopicComplete(ctx, resourceGroup, namespace, topicName, nil, nil)
	if err != nil {
		return nil, err
	}

	results := make([]servicebus.SBSubscription, 0)

	for iteratorSubscription.NotDone() {
		results = append(results, iteratorSubscription.Value())

		if err := iteratorSubscription.Next(); err != nil {
			return nil, err
		}
	}

	return results, nil
}

// ListTopicSubscriptionsE authenticates the subscriptions client and enumerates all values,
// automatically crossing page boundaries as required.
//
// Deprecated: Use [ListTopicSubscriptionsContextE] instead.
func ListTopicSubscriptionsE(subscriptionID string, namespace string, resourceGroup string, topicName string) ([]servicebus.SBSubscription, error) {
	return ListTopicSubscriptionsContextE(context.Background(), subscriptionID, namespace, resourceGroup, topicName)
}

// ListTopicSubscriptionsContext authenticates the subscriptions client and enumerates all values,
// automatically crossing page boundaries as required.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func ListTopicSubscriptionsContext(t testing.TestingT, ctx context.Context, subscriptionID string, namespace string, resourceGroup string, topicName string) []servicebus.SBSubscription {
	t.Helper()

	results, err := ListTopicSubscriptionsContextE(ctx, subscriptionID, namespace, resourceGroup, topicName)
	require.NoError(t, err)

	return results
}

// ListTopicSubscriptions authenticates the subscriptions client and enumerates all values,
// automatically crossing page boundaries as required.
// This function would fail the test if there is an error.
//
// Deprecated: Use [ListTopicSubscriptionsContext] instead.
func ListTopicSubscriptions(t testing.TestingT, subscriptionID string, namespace string, resourceGroup string, topicName string) []servicebus.SBSubscription {
	t.Helper()

	return ListTopicSubscriptionsContext(t, context.Background(), subscriptionID, namespace, resourceGroup, topicName)
}

// ListTopicSubscriptionsNameContextE authenticates the subscriptions client and enumerates all values to get
// a list of subscriptions for the given topic name, automatically crossing page boundaries as required.
// The ctx parameter supports cancellation and timeouts.
func ListTopicSubscriptionsNameContextE(ctx context.Context, subscriptionID string, namespace string, resourceGroup string, topicName string) ([]string, error) {
	sClient, err := serviceBusSubscriptionsClientE(subscriptionID)
	if err != nil {
		return nil, err
	}

	iteratorSubscription, err := sClient.ListByTopicComplete(ctx, resourceGroup, namespace, topicName, nil, nil)
	if err != nil {
		return nil, err
	}

	results := []string{}

	for iteratorSubscription.NotDone() {
		results = append(results, *(iteratorSubscription.Value()).Name)

		if err := iteratorSubscription.Next(); err != nil {
			return nil, err
		}
	}

	return results, nil
}

// ListTopicSubscriptionsNameE authenticates the subscriptions client and enumerates all values to get
// a list of subscriptions for the given topic name, automatically crossing page boundaries as required.
//
// Deprecated: Use [ListTopicSubscriptionsNameContextE] instead.
func ListTopicSubscriptionsNameE(subscriptionID string, namespace string, resourceGroup string, topicName string) ([]string, error) {
	return ListTopicSubscriptionsNameContextE(context.Background(), subscriptionID, namespace, resourceGroup, topicName)
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

// ListTopicSubscriptionsName authenticates the subscriptions client and enumerates all values to get
// a list of subscriptions for the given topic name, automatically crossing page boundaries as required.
// This function would fail the test if there is an error.
//
// Deprecated: Use [ListTopicSubscriptionsNameContext] instead.
func ListTopicSubscriptionsName(t testing.TestingT, subscriptionID string, namespace string, resourceGroup string, topicName string) []string {
	t.Helper()

	return ListTopicSubscriptionsNameContext(t, context.Background(), subscriptionID, namespace, resourceGroup, topicName)
}

// ListTopicAuthRulesContextE authenticates the topic client and enumerates all values to get a list
// of authorization rules for the given topic name, automatically crossing page boundaries as required.
// The ctx parameter supports cancellation and timeouts.
func ListTopicAuthRulesContextE(ctx context.Context, subscriptionID string, namespace string, resourceGroup string, topicName string) ([]string, error) {
	tClient, err := serviceBusTopicClientE(subscriptionID)
	if err != nil {
		return nil, err
	}

	iteratorTopicsRules, err := tClient.ListAuthorizationRulesComplete(ctx, resourceGroup, namespace, topicName)
	if err != nil {
		return nil, err
	}

	results := []string{}

	for iteratorTopicsRules.NotDone() {
		results = append(results, *(iteratorTopicsRules.Value()).Name)

		if err := iteratorTopicsRules.Next(); err != nil {
			return nil, err
		}
	}

	return results, nil
}

// ListTopicAuthRulesE authenticates the topic client and enumerates all values to get a list
// of authorization rules for the given topic name, automatically crossing page boundaries as required.
//
// Deprecated: Use [ListTopicAuthRulesContextE] instead.
func ListTopicAuthRulesE(subscriptionID string, namespace string, resourceGroup string, topicName string) ([]string, error) {
	return ListTopicAuthRulesContextE(context.Background(), subscriptionID, namespace, resourceGroup, topicName)
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

// ListTopicAuthRules authenticates the topic client and enumerates all values to get a list
// of authorization rules for the given topic name, automatically crossing page boundaries as required.
// This function would fail the test if there is an error.
//
// Deprecated: Use [ListTopicAuthRulesContext] instead.
func ListTopicAuthRules(t testing.TestingT, subscriptionID string, namespace string, resourceGroup string, topicName string) []string {
	t.Helper()

	return ListTopicAuthRulesContext(t, context.Background(), subscriptionID, namespace, resourceGroup, topicName)
}
