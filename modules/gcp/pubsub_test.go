package gcp_test

import (
	"context"
	"testing"

	"cloud.google.com/go/pubsub/v2"
	"cloud.google.com/go/pubsub/v2/pstest"
	"github.com/gruntwork-io/terratest/modules/gcp"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
