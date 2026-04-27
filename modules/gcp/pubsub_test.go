//go:build gcp
// +build gcp

// NOTE: We use build tags to differentiate GCP testing for better isolation and parallelism when executing our tests.

//nolint:testpackage // uses unexported newPubSubClient
package gcp

import (
	"context"
	"testing"

	"cloud.google.com/go/pubsub/v2/apiv1/pubsubpb"
	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAssertTopicExistsNoFalseNegative(t *testing.T) {
	t.Parallel()

	projectID := GetGoogleProjectIDFromEnvVar(t)
	topicName := "pubsub-topic-" + random.UniqueID()
	logger.Default.Logf(t, "Creating Pub/Sub topic %s to verify existence check works", topicName)

	CreateTopicContext(t, context.Background(), projectID, topicName)
	defer DeleteTopicContext(t, context.Background(), projectID, topicName)

	AssertTopicExistsContext(t, context.Background(), projectID, topicName)
}

func TestAssertTopicExistsNoFalsePositive(t *testing.T) {
	t.Parallel()

	projectID := GetGoogleProjectIDFromEnvVar(t)
	topicName := "pubsub-topic-" + random.UniqueID()
	logger.Default.Logf(t, "Checking that non-existent Pub/Sub topic %s returns an error", topicName)

	err := AssertTopicExistsContextE(t, context.Background(), projectID, topicName)
	require.Error(t, err, "Expected an error for non-existent Pub/Sub topic, but got none")
}

func TestAssertSubscriptionExistsNoFalseNegative(t *testing.T) {
	t.Parallel()

	projectID := GetGoogleProjectIDFromEnvVar(t)
	topicName := "pubsub-topic-" + random.UniqueID()
	subscriptionName := "pubsub-sub-" + random.UniqueID()
	logger.Default.Logf(t, "Creating Pub/Sub topic %s and subscription %s to verify existence check works", topicName, subscriptionName)

	CreateTopicContext(t, context.Background(), projectID, topicName)
	defer DeleteTopicContext(t, context.Background(), projectID, topicName)

	CreateSubscriptionContext(t, context.Background(), projectID, subscriptionName, topicName)
	defer DeleteSubscriptionContext(t, context.Background(), projectID, subscriptionName)

	AssertSubscriptionExistsContext(t, context.Background(), projectID, subscriptionName)
}

func TestAssertSubscriptionExistsNoFalsePositive(t *testing.T) {
	t.Parallel()

	projectID := GetGoogleProjectIDFromEnvVar(t)
	subscriptionName := "pubsub-sub-" + random.UniqueID()
	logger.Default.Logf(t, "Checking that non-existent Pub/Sub subscription %s returns an error", subscriptionName)

	err := AssertSubscriptionExistsContextE(t, context.Background(), projectID, subscriptionName)
	require.Error(t, err, "Expected an error for non-existent Pub/Sub subscription, but got none")
}

func TestAssertTopicAndSubscriptionExist(t *testing.T) {
	t.Parallel()

	projectID := GetGoogleProjectIDFromEnvVar(t)
	topicName := "pubsub-topic-" + random.UniqueID()
	subscriptionName := "pubsub-sub-" + random.UniqueID()
	logger.Default.Logf(t, "Creating Pub/Sub topic %s and subscription %s", topicName, subscriptionName)

	CreateTopicContext(t, context.Background(), projectID, topicName)
	defer DeleteTopicContext(t, context.Background(), projectID, topicName)

	CreateSubscriptionContext(t, context.Background(), projectID, subscriptionName, topicName)
	defer DeleteSubscriptionContext(t, context.Background(), projectID, subscriptionName)

	AssertTopicExistsContext(t, context.Background(), projectID, topicName)
	AssertSubscriptionExistsContext(t, context.Background(), projectID, subscriptionName)

	// Verify subscription is linked to the correct topic
	client, err := newPubSubClient(context.Background(), projectID)
	require.NoError(t, err)

	defer func() { _ = client.Close() }()

	sub, err := client.SubscriptionAdminClient.GetSubscription(context.Background(), &pubsubpb.GetSubscriptionRequest{
		Subscription: subscriptionResource(projectID, subscriptionName),
	})
	require.NoError(t, err)
	assert.Equal(t, topicResource(projectID, topicName), sub.GetTopic())
}
