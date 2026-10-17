package gcp

import (
	"context"
	"errors"
	"fmt"

	"github.com/gruntwork-io/terratest/modules/core/v2/logger"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/contentwarehouse/v1"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
)

// GetDocumentAIWarehouseDocumentSchemaAttrs returns the settings Google Cloud holds for the Document AI Warehouse document schema, so a test can assert on what was
// actually created rather than only that it exists.
// The property definitions a schema declares are the whole resource, and they are nested, so a test
// asserts on the names and types rather than on the schema as a whole.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetDocumentAIWarehouseDocumentSchemaAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, schemaID string) *contentwarehouse.GoogleCloudContentwarehouseV1DocumentSchema {
	result, err := GetDocumentAIWarehouseDocumentSchemaAttrsE(t, ctx, projectID, location, schemaID)
	require.NoError(t, err)

	return result
}

// GetDocumentAIWarehouseDocumentSchemaAttrsE returns the settings Google Cloud holds for the Document AI Warehouse document schema.
// The ctx parameter supports cancellation and timeouts.
func GetDocumentAIWarehouseDocumentSchemaAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, schemaID string) (*contentwarehouse.GoogleCloudContentwarehouseV1DocumentSchema, error) {
	logger.Default.Logf(t, "Getting settings for Document AI Warehouse document schema %s in %s in project %s", schemaID, location, projectID)

	service, err := NewContentWarehouseServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetDocumentAIWarehouseDocumentSchemaAttrsWithClient(ctx, service, projectID, location, schemaID)
}

// GetDocumentAIWarehouseDocumentSchemaAttrsWithClient returns the settings Google Cloud holds for the Document AI Warehouse document schema using the supplied
// *contentwarehouse.Service. Prefer this variant in unit tests where the service is backed by an httptest
// fake server (see contentwarehouse_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetDocumentAIWarehouseDocumentSchemaAttrsWithClient(ctx context.Context, service *contentwarehouse.Service, projectID string, location string, schemaID string) (*contentwarehouse.GoogleCloudContentwarehouseV1DocumentSchema, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/documentSchemas/%s", projectID, location, schemaID)

	result, err := service.Projects.Locations.DocumentSchemas.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Document AI Warehouse document schema %s in %s in project %s does not exist", schemaID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Document AI Warehouse document schema %s in %s in project %s: %w", schemaID, location, projectID, err)
	}

	return result, nil
}

// NewContentWarehouseServiceE creates a Document AI Warehouse service authenticated the same way every other client in this
// module is.
// Google serves Document AI Warehouse from the contentwarehouse API.
// The ctx parameter supports cancellation and timeouts.
func NewContentWarehouseServiceE(t testing.TestingT, ctx context.Context) (*contentwarehouse.Service, error) {
	return contentwarehouse.NewService(ctx, append(withOptions(), option.WithScopes(contentwarehouse.CloudPlatformScope))...)
}
