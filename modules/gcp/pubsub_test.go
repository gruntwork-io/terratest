package gcp_test

import (
	"context"
	"testing"
	"time"

	"cloud.google.com/go/pubsub/v2"
	"cloud.google.com/go/pubsub/v2/apiv1/pubsubpb"
	"cloud.google.com/go/pubsub/v2/pstest"
	"github.com/gruntwork-io/terratest/modules/gcp/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/durationpb"
)

// newFakePubSubClient wires a *pubsub.Client to the pstest in-memory fake Pub/Sub server —
// credential-free, gRPC-conformant.
func newFakePubSubClient(t *testing.T) *pubsub.Client {
	t.Helper()

	srv := pstest.NewServer()

	t.Cleanup(func() { _ = srv.Close() })

	conn, err := grpc.NewClient(srv.Addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	t.Cleanup(func() { _ = conn.Close() })

	client, err := pubsub.NewClient(context.Background(), "test-project", option.WithGRPCConn(conn))
	require.NoError(t, err)
	t.Cleanup(func() { _ = client.Close() })

	return client
}

func TestPubSubTopicLifecycleWithClient(t *testing.T) {
	t.Parallel()

	client := newFakePubSubClient(t)
	ctx := context.Background()

	// Missing topic: AssertTopicExists fails, Delete fails.
	require.ErrorContains(t, gcp.AssertTopicExistsWithClient(ctx, client, "missing"), "does not exist")
	require.Error(t, gcp.DeleteTopicWithClient(ctx, client, "missing"))

	// Create / assert / delete roundtrip.
	require.NoError(t, gcp.CreateTopicWithClient(ctx, client, "t"))
	require.NoError(t, gcp.AssertTopicExistsWithClient(ctx, client, "t"))
	require.NoError(t, gcp.DeleteTopicWithClient(ctx, client, "t"))
	require.ErrorContains(t, gcp.AssertTopicExistsWithClient(ctx, client, "t"), "does not exist")

	// Re-creating the same topic name within the same run is an error on pstest (matches prod).
	require.NoError(t, gcp.CreateTopicWithClient(ctx, client, "dup"))
	require.ErrorContains(t, gcp.CreateTopicWithClient(ctx, client, "dup"), "failed to create")
}

func TestGetTopicAttrsWithClient(t *testing.T) {
	t.Parallel()

	client := newFakePubSubClient(t)
	ctx := context.Background()

	// The error names the topic and the project as well as saying it is absent, so all three
	// are asserted rather than only the phrase.
	_, err := gcp.GetTopicAttrsWithClient(ctx, client, "missing")
	require.ErrorContains(t, err, "does not exist")
	require.ErrorContains(t, err, "missing")
	require.ErrorContains(t, err, "test-project")

	// The values are the ones the terraform-google-messaging topic module sets, because the point
	// of reading settings back is asserting a module configured the topic it was asked for.
	_, err = client.TopicAdminClient.CreateTopic(ctx, &pubsubpb.Topic{
		Name:                     "projects/test-project/topics/configured",
		Labels:                   map[string]string{"purpose": "terratest"},
		MessageRetentionDuration: durationpb.New(600 * time.Second),
		MessageStoragePolicy: &pubsubpb.MessageStoragePolicy{
			AllowedPersistenceRegions: []string{"us-central1"},
		},
	})
	require.NoError(t, err)

	topic, err := gcp.GetTopicAttrsWithClient(ctx, client, "configured")
	require.NoError(t, err)

	assert.Equal(t, "projects/test-project/topics/configured", topic.GetName())
	assert.Equal(t, map[string]string{"purpose": "terratest"}, topic.GetLabels())
	assert.Equal(t, 600*time.Second, topic.GetMessageRetentionDuration().AsDuration())
	assert.Equal(t, []string{"us-central1"}, topic.GetMessageStoragePolicy().GetAllowedPersistenceRegions())
}

func TestPubSubSubscriptionLifecycleWithClient(t *testing.T) {
	t.Parallel()

	client := newFakePubSubClient(t)
	ctx := context.Background()

	// Subscription on a missing topic errors.
	require.ErrorContains(t, gcp.CreateSubscriptionWithClient(ctx, client, "s", "missing-topic"), "failed to create")

	// Full roundtrip on a real topic.
	require.NoError(t, gcp.CreateTopicWithClient(ctx, client, "t"))
	require.NoError(t, gcp.CreateSubscriptionWithClient(ctx, client, "s", "t"))
	require.NoError(t, gcp.AssertSubscriptionExistsWithClient(ctx, client, "s"))
	require.NoError(t, gcp.DeleteSubscriptionWithClient(ctx, client, "s"))
	require.ErrorContains(t, gcp.AssertSubscriptionExistsWithClient(ctx, client, "s"), "does not exist")
}
