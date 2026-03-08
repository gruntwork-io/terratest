package gcp

import (
	"context"
	"fmt"

	"cloud.google.com/go/pubsub"
	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/testing"
)

// AssertTopicExistsContext checks if the given Pub/Sub topic exists and fails the test if it does not.
// The ctx parameter supports cancellation and timeouts.
func AssertTopicExistsContext(t testing.TestingT, ctx context.Context, projectID string, topicName string) {
	err := AssertTopicExistsContextE(t, ctx, projectID, topicName)
	if err != nil {
		t.Fatal(err)
	}
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

	exists, err := client.Topic(topicName).Exists(ctx)
	if err != nil {
		return fmt.Errorf("failed to check if Pub/Sub topic %s exists in project %s: %w", topicName, projectID, err)
	}

	if !exists {
		return fmt.Errorf("Pub/Sub topic %s does not exist in project %s", topicName, projectID)
	}

	return nil
}

// AssertSubscriptionExistsContext checks if the given Pub/Sub subscription exists and fails the test if it does not.
// The ctx parameter supports cancellation and timeouts.
func AssertSubscriptionExistsContext(t testing.TestingT, ctx context.Context, projectID string, subscriptionName string) {
	err := AssertSubscriptionExistsContextE(t, ctx, projectID, subscriptionName)
	if err != nil {
		t.Fatal(err)
	}
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

	exists, err := client.Subscription(subscriptionName).Exists(ctx)
	if err != nil {
		return fmt.Errorf("failed to check if Pub/Sub subscription %s exists in project %s: %w", subscriptionName, projectID, err)
	}

	if !exists {
		return fmt.Errorf("Pub/Sub subscription %s does not exist in project %s", subscriptionName, projectID)
	}

	return nil
}

// CreateTopicContext creates a new Pub/Sub topic and fails the test if it cannot.
// The ctx parameter supports cancellation and timeouts.
func CreateTopicContext(t testing.TestingT, ctx context.Context, projectID string, topicName string) {
	err := CreateTopicContextE(t, ctx, projectID, topicName)
	if err != nil {
		t.Fatal(err)
	}
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

	_, err = client.CreateTopic(ctx, topicName)
	if err != nil {
		return fmt.Errorf("failed to create Pub/Sub topic %s in project %s: %w", topicName, projectID, err)
	}

	return nil
}

// DeleteTopicContext deletes the given Pub/Sub topic and fails the test if it cannot.
// The ctx parameter supports cancellation and timeouts.
func DeleteTopicContext(t testing.TestingT, ctx context.Context, projectID string, topicName string) {
	err := DeleteTopicContextE(t, ctx, projectID, topicName)
	if err != nil {
		t.Fatal(err)
	}
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

	if err := client.Topic(topicName).Delete(ctx); err != nil {
		return fmt.Errorf("failed to delete Pub/Sub topic %s in project %s: %w", topicName, projectID, err)
	}

	return nil
}

// CreateSubscriptionContext creates a new Pub/Sub subscription on the given topic and fails the test if it cannot.
// The ctx parameter supports cancellation and timeouts.
func CreateSubscriptionContext(t testing.TestingT, ctx context.Context, projectID string, subscriptionName string, topicName string) {
	err := CreateSubscriptionContextE(t, ctx, projectID, subscriptionName, topicName)
	if err != nil {
		t.Fatal(err)
	}
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

	_, err = client.CreateSubscription(ctx, subscriptionName, pubsub.SubscriptionConfig{
		Topic: client.Topic(topicName),
	})
	if err != nil {
		return fmt.Errorf("failed to create Pub/Sub subscription %s in project %s: %w", subscriptionName, projectID, err)
	}

	return nil
}

// DeleteSubscriptionContext deletes the given Pub/Sub subscription and fails the test if it cannot.
// The ctx parameter supports cancellation and timeouts.
func DeleteSubscriptionContext(t testing.TestingT, ctx context.Context, projectID string, subscriptionName string) {
	err := DeleteSubscriptionContextE(t, ctx, projectID, subscriptionName)
	if err != nil {
		t.Fatal(err)
	}
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

	if err := client.Subscription(subscriptionName).Delete(ctx); err != nil {
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
