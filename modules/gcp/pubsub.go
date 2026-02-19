package gcp

import (
    "context"

    "cloud.google.com/go/pubsub"

    "github.com/gruntwork-io/terratest/modules/logger"
    "github.com/gruntwork-io/terratest/modules/testing"
)

// AssertPubSubTopicExists checks if the given pubsub topic exists and fails the test if it does not.
func AssertPubSubTopicExists(t testing.TestingT, projectID, topicID string) {
    err := AssertPubSubTopicExistsE(t, projectID, topicID)
    if err != nil {
        t.Fatal(err)
    }
}

// AssertPubSubTopicExistsE checks if the given pubsub topic exists and returns an error if it does not.
func AssertPubSubTopicExistsE(t testing.TestingT, projectID, topicID string) error {
    logger.Logf(t, "Finding topic %s", topicID)
    
    ctx := context.Background()
    
    // Creates a client.
    client, err := pubsub.NewClient(ctx, projectID)
    if err != nil {
        return err
    }
    
    topic := client.Topic(topicID)
    
    ok, err := topic.Exists(ctx)
    if err != nil {
        return err
    }
    if !ok {
        return err
    }
    return nil
}

// AssertPubSubSubscriptionExists checks if the given pubsub subscription exists and fails the test if it does not.
func AssertPubSubSubscriptionExists(t testing.TestingT, projectID, subID string) {
    err := AssertPubSubSubscriptionExistsE(t, projectID, subID)

    if err != nil {
        t.Fatal(err)
    }
}

// AssertPubSubSubscriptionExistsE checks if the given pubsub subscription exists and returns an error if it does not.
func AssertPubSubSubscriptionExistsE(t testing.TestingT, projectID, subID string) error {
    logger.Logf(t, "Finding subscription %s", subID)
    
    ctx := context.Background()
    
    // Creates a client.
    client, err := pubsub.NewClient(ctx, projectID)
    if err != nil {
        return err
    }
    
    sub := client.Subscription(subID)
    ok, err := sub.Exists(ctx)
    if err != nil {
        return err
    }
    if !ok {
        return err
    }
    return nil
}