package gcp

import (
	"context"
	"fmt"

	"cloud.google.com/go/pubsub/v2"
	"cloud.google.com/go/pubsub/v2/apiv1/pubsubpb"
	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// AssertTopicExists checks if the given Pub/Sub topic exists and fails the test if it does not.
//
// Deprecated: Use [AssertTopicExistsContext] instead.
func AssertTopicExists(t testing.TestingT, projectID string, topicName string) {
	AssertTopicExistsContext(t, context.Background(), projectID, topicName)
}

// AssertTopicExistsContext checks if the given Pub/Sub topic exists and fails the test if it does not.
// The ctx parameter supports cancellation and timeouts.
func AssertTopicExistsContext(t testing.TestingT, ctx context.Context, projectID string, topicName string) {
	err := AssertTopicExistsContextE(t, ctx, projectID, topicName)
	require.NoError(t, err)
}

// AssertTopicExistsE checks if the given Pub/Sub topic exists and returns an error if it does not.
//
// Deprecated: Use [AssertTopicExistsContextE] instead.
func AssertTopicExistsE(t testing.TestingT, projectID string, topicName string) error {
	return AssertTopicExistsContextE(t, context.Background(), projectID, topicName)
}

// AssertTopicExistsContextE checks if the given Pub/Sub topic exists and returns an error if it does not.
// The ctx parameter supports cancellation and timeouts.
func AssertTopicExistsContextE(t testing.TestingT, ctx context.Context, projectID string, topicName string) error {
	logger.Default.Logf(t, "Verifying Pub/Sub topic %s exists in project %s", topicName, projectID)

	client, err := newPubSubClient(ctx, projectID)
	if err != nil {
		return err
	}

	defer func() { _ = client.Close() }()

	return AssertTopicExistsWithClient(ctx, client, topicName)
}

// AssertTopicExistsWithClient checks if the given Pub/Sub topic exists using the supplied *pubsub.Client.
// Prefer this variant in unit tests where the client is backed by a pstest in-memory fake server
// (see pubsub_unit_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func AssertTopicExistsWithClient(ctx context.Context, client *pubsub.Client, topicName string) error {
	projectID := client.Project()

	_, err := client.TopicAdminClient.GetTopic(ctx, &pubsubpb.GetTopicRequest{
		Topic: topicResource(projectID, topicName),
	})
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return fmt.Errorf("Pub/Sub topic %s does not exist in project %s", topicName, projectID)
		}

		return fmt.Errorf("failed to check if Pub/Sub topic %s exists in project %s: %w", topicName, projectID, err)
	}

	return nil
}

// AssertSubscriptionExists checks if the given Pub/Sub subscription exists and fails the test if it does not.
//
// Deprecated: Use [AssertSubscriptionExistsContext] instead.
func AssertSubscriptionExists(t testing.TestingT, projectID string, subscriptionName string) {
	AssertSubscriptionExistsContext(t, context.Background(), projectID, subscriptionName)
}

// AssertSubscriptionExistsContext checks if the given Pub/Sub subscription exists and fails the test if it does not.
// The ctx parameter supports cancellation and timeouts.
func AssertSubscriptionExistsContext(t testing.TestingT, ctx context.Context, projectID string, subscriptionName string) {
	err := AssertSubscriptionExistsContextE(t, ctx, projectID, subscriptionName)
	require.NoError(t, err)
}

// AssertSubscriptionExistsE checks if the given Pub/Sub subscription exists and returns an error if it does not.
//
// Deprecated: Use [AssertSubscriptionExistsContextE] instead.
func AssertSubscriptionExistsE(t testing.TestingT, projectID string, subscriptionName string) error {
	return AssertSubscriptionExistsContextE(t, context.Background(), projectID, subscriptionName)
}

// AssertSubscriptionExistsContextE checks if the given Pub/Sub subscription exists and returns an error if it does not.
// The ctx parameter supports cancellation and timeouts.
func AssertSubscriptionExistsContextE(t testing.TestingT, ctx context.Context, projectID string, subscriptionName string) error {
	logger.Default.Logf(t, "Verifying Pub/Sub subscription %s exists in project %s", subscriptionName, projectID)

	client, err := newPubSubClient(ctx, projectID)
	if err != nil {
		return err
	}

	defer func() { _ = client.Close() }()

	return AssertSubscriptionExistsWithClient(ctx, client, subscriptionName)
}

// AssertSubscriptionExistsWithClient checks if the given Pub/Sub subscription exists using the supplied *pubsub.Client.
// Prefer this variant in unit tests where the client is backed by a pstest in-memory fake server
// (see pubsub_unit_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func AssertSubscriptionExistsWithClient(ctx context.Context, client *pubsub.Client, subscriptionName string) error {
	projectID := client.Project()

	_, err := client.SubscriptionAdminClient.GetSubscription(ctx, &pubsubpb.GetSubscriptionRequest{
		Subscription: subscriptionResource(projectID, subscriptionName),
	})
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return fmt.Errorf("Pub/Sub subscription %s does not exist in project %s", subscriptionName, projectID)
		}

		return fmt.Errorf("failed to check if Pub/Sub subscription %s exists in project %s: %w", subscriptionName, projectID, err)
	}

	return nil
}

// CreateTopic creates a new Pub/Sub topic and fails the test if it cannot.
//
// Deprecated: Use [CreateTopicContext] instead.
func CreateTopic(t testing.TestingT, projectID string, topicName string) {
	CreateTopicContext(t, context.Background(), projectID, topicName)
}

// CreateTopicContext creates a new Pub/Sub topic and fails the test if it cannot.
// The ctx parameter supports cancellation and timeouts.
func CreateTopicContext(t testing.TestingT, ctx context.Context, projectID string, topicName string) {
	err := CreateTopicContextE(t, ctx, projectID, topicName)
	require.NoError(t, err)
}

// CreateTopicE creates a new Pub/Sub topic and returns an error if it fails.
//
// Deprecated: Use [CreateTopicContextE] instead.
func CreateTopicE(t testing.TestingT, projectID string, topicName string) error {
	return CreateTopicContextE(t, context.Background(), projectID, topicName)
}

// CreateTopicContextE creates a new Pub/Sub topic and returns an error if it fails.
// The ctx parameter supports cancellation and timeouts.
func CreateTopicContextE(t testing.TestingT, ctx context.Context, projectID string, topicName string) error {
	logger.Default.Logf(t, "Creating Pub/Sub topic %s in project %s", topicName, projectID)

	client, err := newPubSubClient(ctx, projectID)
	if err != nil {
		return err
	}

	defer func() { _ = client.Close() }()

	return CreateTopicWithClient(ctx, client, topicName)
}

// CreateTopicWithClient creates a new Pub/Sub topic using the supplied *pubsub.Client.
// Prefer this variant in unit tests where the client is backed by a pstest in-memory fake server
// (see pubsub_unit_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func CreateTopicWithClient(ctx context.Context, client *pubsub.Client, topicName string) error {
	projectID := client.Project()
	if _, err := client.TopicAdminClient.CreateTopic(ctx, &pubsubpb.Topic{
		Name: topicResource(projectID, topicName),
	}); err != nil {
		return fmt.Errorf("failed to create Pub/Sub topic %s in project %s: %w", topicName, projectID, err)
	}

	return nil
}

// DeleteTopic deletes the given Pub/Sub topic and fails the test if it cannot.
//
// Deprecated: Use [DeleteTopicContext] instead.
func DeleteTopic(t testing.TestingT, projectID string, topicName string) {
	DeleteTopicContext(t, context.Background(), projectID, topicName)
}

// DeleteTopicContext deletes the given Pub/Sub topic and fails the test if it cannot.
// The ctx parameter supports cancellation and timeouts.
func DeleteTopicContext(t testing.TestingT, ctx context.Context, projectID string, topicName string) {
	err := DeleteTopicContextE(t, ctx, projectID, topicName)
	require.NoError(t, err)
}

// DeleteTopicE deletes the given Pub/Sub topic and returns an error if it fails.
//
// Deprecated: Use [DeleteTopicContextE] instead.
func DeleteTopicE(t testing.TestingT, projectID string, topicName string) error {
	return DeleteTopicContextE(t, context.Background(), projectID, topicName)
}

// DeleteTopicContextE deletes the given Pub/Sub topic and returns an error if it fails.
// The ctx parameter supports cancellation and timeouts.
func DeleteTopicContextE(t testing.TestingT, ctx context.Context, projectID string, topicName string) error {
	logger.Default.Logf(t, "Deleting Pub/Sub topic %s in project %s", topicName, projectID)

	client, err := newPubSubClient(ctx, projectID)
	if err != nil {
		return err
	}

	defer func() { _ = client.Close() }()

	return DeleteTopicWithClient(ctx, client, topicName)
}

// DeleteTopicWithClient deletes the given Pub/Sub topic using the supplied *pubsub.Client.
// Prefer this variant in unit tests where the client is backed by a pstest in-memory fake server
// (see pubsub_unit_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func DeleteTopicWithClient(ctx context.Context, client *pubsub.Client, topicName string) error {
	projectID := client.Project()
	if err := client.TopicAdminClient.DeleteTopic(ctx, &pubsubpb.DeleteTopicRequest{
		Topic: topicResource(projectID, topicName),
	}); err != nil {
		return fmt.Errorf("failed to delete Pub/Sub topic %s in project %s: %w", topicName, projectID, err)
	}

	return nil
}

// CreateSubscription creates a new Pub/Sub subscription on the given topic and fails the test if it cannot.
//
// Deprecated: Use [CreateSubscriptionContext] instead.
func CreateSubscription(t testing.TestingT, projectID string, subscriptionName string, topicName string) {
	CreateSubscriptionContext(t, context.Background(), projectID, subscriptionName, topicName)
}

// CreateSubscriptionContext creates a new Pub/Sub subscription on the given topic and fails the test if it cannot.
// The ctx parameter supports cancellation and timeouts.
func CreateSubscriptionContext(t testing.TestingT, ctx context.Context, projectID string, subscriptionName string, topicName string) {
	err := CreateSubscriptionContextE(t, ctx, projectID, subscriptionName, topicName)
	require.NoError(t, err)
}

// CreateSubscriptionE creates a new Pub/Sub subscription on the given topic and returns an error if it fails.
//
// Deprecated: Use [CreateSubscriptionContextE] instead.
func CreateSubscriptionE(t testing.TestingT, projectID string, subscriptionName string, topicName string) error {
	return CreateSubscriptionContextE(t, context.Background(), projectID, subscriptionName, topicName)
}

// CreateSubscriptionContextE creates a new Pub/Sub subscription on the given topic and returns an error if it fails.
// The ctx parameter supports cancellation and timeouts.
func CreateSubscriptionContextE(t testing.TestingT, ctx context.Context, projectID string, subscriptionName string, topicName string) error {
	logger.Default.Logf(t, "Creating Pub/Sub subscription %s on topic %s in project %s", subscriptionName, topicName, projectID)

	client, err := newPubSubClient(ctx, projectID)
	if err != nil {
		return err
	}

	defer func() { _ = client.Close() }()

	return CreateSubscriptionWithClient(ctx, client, subscriptionName, topicName)
}

// CreateSubscriptionWithClient creates a new Pub/Sub subscription on the given topic using the supplied *pubsub.Client.
// Prefer this variant in unit tests where the client is backed by a pstest in-memory fake server
// (see pubsub_unit_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func CreateSubscriptionWithClient(ctx context.Context, client *pubsub.Client, subscriptionName string, topicName string) error {
	projectID := client.Project()
	if _, err := client.SubscriptionAdminClient.CreateSubscription(ctx, &pubsubpb.Subscription{
		Name:  subscriptionResource(projectID, subscriptionName),
		Topic: topicResource(projectID, topicName),
	}); err != nil {
		return fmt.Errorf("failed to create Pub/Sub subscription %s on topic %s in project %s: %w", subscriptionName, topicName, projectID, err)
	}

	return nil
}

// DeleteSubscription deletes the given Pub/Sub subscription and fails the test if it cannot.
//
// Deprecated: Use [DeleteSubscriptionContext] instead.
func DeleteSubscription(t testing.TestingT, projectID string, subscriptionName string) {
	DeleteSubscriptionContext(t, context.Background(), projectID, subscriptionName)
}

// DeleteSubscriptionContext deletes the given Pub/Sub subscription and fails the test if it cannot.
// The ctx parameter supports cancellation and timeouts.
func DeleteSubscriptionContext(t testing.TestingT, ctx context.Context, projectID string, subscriptionName string) {
	err := DeleteSubscriptionContextE(t, ctx, projectID, subscriptionName)
	require.NoError(t, err)
}

// DeleteSubscriptionE deletes the given Pub/Sub subscription and returns an error if it fails.
//
// Deprecated: Use [DeleteSubscriptionContextE] instead.
func DeleteSubscriptionE(t testing.TestingT, projectID string, subscriptionName string) error {
	return DeleteSubscriptionContextE(t, context.Background(), projectID, subscriptionName)
}

// DeleteSubscriptionContextE deletes the given Pub/Sub subscription and returns an error if it fails.
// The ctx parameter supports cancellation and timeouts.
func DeleteSubscriptionContextE(t testing.TestingT, ctx context.Context, projectID string, subscriptionName string) error {
	logger.Default.Logf(t, "Deleting Pub/Sub subscription %s in project %s", subscriptionName, projectID)

	client, err := newPubSubClient(ctx, projectID)
	if err != nil {
		return err
	}

	defer func() { _ = client.Close() }()

	return DeleteSubscriptionWithClient(ctx, client, subscriptionName)
}

// DeleteSubscriptionWithClient deletes the given Pub/Sub subscription using the supplied *pubsub.Client.
// Prefer this variant in unit tests where the client is backed by a pstest in-memory fake server
// (see pubsub_unit_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func DeleteSubscriptionWithClient(ctx context.Context, client *pubsub.Client, subscriptionName string) error {
	projectID := client.Project()
	if err := client.SubscriptionAdminClient.DeleteSubscription(ctx, &pubsubpb.DeleteSubscriptionRequest{
		Subscription: subscriptionResource(projectID, subscriptionName),
	}); err != nil {
		return fmt.Errorf("failed to delete Pub/Sub subscription %s in project %s: %w", subscriptionName, projectID, err)
	}

	return nil
}

// newPubSubClient creates a new Pub/Sub client using the provided project ID and global GCP auth options.
func newPubSubClient(ctx context.Context, projectID string) (*pubsub.Client, error) {
	client, err := pubsub.NewClient(ctx, projectID, withOptions()...)
	if err != nil {
		return nil, err
	}

	return client, nil
}

// topicResource returns the fully-qualified Pub/Sub topic resource name ("projects/<p>/topics/<t>"),
// required by the v2 admin client APIs.
func topicResource(projectID, topicName string) string {
	return fmt.Sprintf("projects/%s/topics/%s", projectID, topicName)
}

// subscriptionResource returns the fully-qualified Pub/Sub subscription resource name
// ("projects/<p>/subscriptions/<s>"), required by the v2 admin client APIs.
func subscriptionResource(projectID, subscriptionName string) string {
	return fmt.Sprintf("projects/%s/subscriptions/%s", projectID, subscriptionName)
}
