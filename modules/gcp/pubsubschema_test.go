package gcp_test

import (
	"context"
	"testing"

	pubsubapi "cloud.google.com/go/pubsub/v2/apiv1"
	"cloud.google.com/go/pubsub/v2/apiv1/pubsubpb"
	"cloud.google.com/go/pubsub/v2/pstest"
	"github.com/gruntwork-io/terratest/modules/gcp/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// newFakeSchemaClient wires a *pubsubapi.SchemaClient to the pstest in-memory fake Pub/Sub server,
// credential-free and gRPC-conformant, the way newFakePubSubClient does for topics.
func newFakeSchemaClient(t *testing.T) *pubsubapi.SchemaClient {
	t.Helper()

	srv := pstest.NewServer()
	t.Cleanup(func() { _ = srv.Close() })

	conn, err := grpc.NewClient(srv.Addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	t.Cleanup(func() { _ = conn.Close() })

	client, err := pubsubapi.NewSchemaClient(context.Background(), option.WithGRPCConn(conn))
	require.NoError(t, err)
	t.Cleanup(func() { _ = client.Close() })

	return client
}

func TestGetSchemaAttrsWithClient(t *testing.T) {
	t.Parallel()

	client := newFakeSchemaClient(t)
	ctx := context.Background()

	// The values are the ones the terraform-google-messaging schema module sets, because the point
	// of reading settings back is asserting a module configured the schema it was asked for.
	definition := `{"type":"record","name":"Event","fields":[{"name":"id","type":"string"}]}`
	_, err := client.CreateSchema(ctx, &pubsubpb.CreateSchemaRequest{
		Parent:   "projects/test-project",
		SchemaId: "events",
		Schema:   &pubsubpb.Schema{Type: pubsubpb.Schema_AVRO, Definition: definition},
	})
	require.NoError(t, err)

	schema, err := gcp.GetSchemaAttrsWithClient(ctx, client, "test-project", "events")
	require.NoError(t, err)

	assert.Equal(t, "projects/test-project/schemas/events", schema.Name)
	assert.Equal(t, pubsubpb.Schema_AVRO, schema.Type)
	assert.Equal(t, definition, schema.Definition)
}

func TestGetSchemaAttrsWithClientMissingSchema(t *testing.T) {
	t.Parallel()

	// The error names the schema and the project as well as saying it is absent, so all three are
	// asserted rather than only the phrase.
	_, err := gcp.GetSchemaAttrsWithClient(context.Background(), newFakeSchemaClient(t), "test-project", "gone")
	require.ErrorContains(t, err, "does not exist")
	require.ErrorContains(t, err, "gone")
	require.ErrorContains(t, err, "test-project")
}
