package gcp

import (
	"context"
	"fmt"

	"cloud.google.com/go/pubsub"
	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/testing"
)

// AssertTopicExists checks if the given Pub/Sub topic exists and fails the test if it does not.
func AssertTopicExists(t testing.TestingT, projectID string, topicName string) {
	err := AssertTopicExistsE(t, projectID, topicName)
	if err != nil {
		t.Fatal(err)
	}
}

// AssertTopicExistsE checks if the given Pub/Sub topic exists and returns an error if it does not.
func AssertTopicExistsE(t testing.TestingT, projectID string, topicName string) error {
	logger.Default.Logf(t, "Verifying Pub/Sub topic %s exists in project %s", topicName, projectID)

	ctx := context.Background()

	client, err := newPubSubClient(projectID)
	if err != nil {
		return err
	}
	defer func() { _ = client.Close() }()

	exists, err := client.Topic(topicName).Exists(ctx)
	if err != nil {
		return err
	}

	if !exists {
		return fmt.Errorf("Pub/Sub Topic %s does not exist in project %s", topicName, projectID)
	}

	return nil
}

// AssertSubscriptionExists checks if the given Pub/Sub subscription exists and fails the test if it does not.
func AssertSubscriptionExists(t testing.TestingT, projectID string, subscriptionName string) {
	err := AssertSubscriptionExistsE(t, projectID, subscriptionName)
	if err != nil {
		t.Fatal(err)
	}
}

// AssertSubscriptionExistsE checks if the given Pub/Sub subscription exists and returns an error if it does not.
func AssertSubscriptionExistsE(t testing.TestingT, projectID string, subscriptionName string) error {
	logger.Default.Logf(t, "Verifying Pub/Sub subscription %s exists in project %s", subscriptionName, projectID)

	ctx := context.Background()

	client, err := newPubSubClient(projectID)
	if err != nil {
		return err
	}
	defer func() { _ = client.Close() }()

	exists, err := client.Subscription(subscriptionName).Exists(ctx)
	if err != nil {
		return err
	}

	if !exists {
		return fmt.Errorf("Pub/Sub Subscription %s does not exist in project %s", subscriptionName, projectID)
	}

	return nil
}

// newPubSubClient creates a new Pub/Sub client using the provided project ID and global GCP auth options.
func newPubSubClient(projectID string) (*pubsub.Client, error) {
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID, withOptions()...)
	if err != nil {
		return nil, err
	}
	return client, nil
}
