package gcp

import (
	"context"
	"errors"
	"fmt"

	pubsubapi "cloud.google.com/go/pubsub/v2/apiv1"
	"cloud.google.com/go/pubsub/v2/apiv1/pubsubpb"
	"github.com/gruntwork-io/terratest/modules/core/v2/logger"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GetSchemaAttrs returns the settings Google Cloud holds for the given Pub/Sub schema, so a test
// can assert on what was actually created rather than only that it exists. The definition comes
// back with it, so a caller can check the schema text Google accepted.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetSchemaAttrs(t testing.TestingT, ctx context.Context, projectID string, schemaName string) *pubsubpb.Schema {
	schema, err := GetSchemaAttrsE(t, ctx, projectID, schemaName)
	require.NoError(t, err)

	return schema
}

// GetSchemaAttrsE returns the settings Google Cloud holds for the given Pub/Sub schema.
// The ctx parameter supports cancellation and timeouts.
func GetSchemaAttrsE(t testing.TestingT, ctx context.Context, projectID string, schemaName string) (schema *pubsubpb.Schema, err error) {
	logger.Default.Logf(t, "Getting settings for Pub/Sub schema %s in project %s", schemaName, projectID)

	client, err := pubsubapi.NewSchemaClient(ctx, withOptions()...)
	if err != nil {
		return nil, err
	}

	defer func() { err = errors.Join(err, client.Close()) }()

	return GetSchemaAttrsWithClient(ctx, client, projectID, schemaName)
}

// GetSchemaAttrsWithClient returns the settings Google Cloud holds for the given Pub/Sub schema
// using the supplied *pubsubapi.SchemaClient. Prefer this variant in unit tests where the client is
// backed by a pstest in-memory fake server (see pubsubschema_test.go for the pattern). The schema
// client is separate from *pubsub.Client, which has no schema surface.
// The ctx parameter supports cancellation and timeouts.
func GetSchemaAttrsWithClient(ctx context.Context, client *pubsubapi.SchemaClient, projectID string, schemaName string) (*pubsubpb.Schema, error) {
	schema, err := client.GetSchema(ctx, &pubsubpb.GetSchemaRequest{
		Name: fmt.Sprintf("projects/%s/schemas/%s", projectID, schemaName),
		View: pubsubpb.SchemaView_FULL,
	})
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, fmt.Errorf("Pub/Sub schema %s does not exist in project %s", schemaName, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Pub/Sub schema %s in project %s: %w", schemaName, projectID, err)
	}

	return schema, nil
}
