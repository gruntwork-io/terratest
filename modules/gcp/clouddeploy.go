package gcp

import (
	"context"
	"errors"
	"fmt"

	"github.com/gruntwork-io/terratest/modules/core/v2/logger"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/clouddeploy/v1"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
)

// GetDeliveryPipelineAttrs returns the settings Google Cloud holds for the given Cloud Deploy
// delivery pipeline, so a test can assert on what was actually created rather than only that it
// exists. A pipeline is regional, so the location it was created in has to be given.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetDeliveryPipelineAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, pipelineName string) *clouddeploy.DeliveryPipeline {
	pipeline, err := GetDeliveryPipelineAttrsE(t, ctx, projectID, location, pipelineName)
	require.NoError(t, err)

	return pipeline
}

// GetDeliveryPipelineAttrsE returns the settings Google Cloud holds for the given Cloud Deploy
// delivery pipeline.
// The ctx parameter supports cancellation and timeouts.
func GetDeliveryPipelineAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, pipelineName string) (*clouddeploy.DeliveryPipeline, error) {
	logger.Default.Logf(t, "Getting settings for delivery pipeline %s in project %s", pipelineName, projectID)

	service, err := NewCloudDeployServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetDeliveryPipelineAttrsWithClient(ctx, service, projectID, location, pipelineName)
}

// GetDeliveryPipelineAttrsWithClient returns the settings Google Cloud holds for the given Cloud
// Deploy delivery pipeline using the supplied *clouddeploy.Service. Prefer this variant in unit
// tests where the service is backed by an httptest fake server (see clouddeploy_test.go for the
// pattern).
// The ctx parameter supports cancellation and timeouts.
func GetDeliveryPipelineAttrsWithClient(ctx context.Context, service *clouddeploy.Service, projectID string, location string, pipelineName string) (*clouddeploy.DeliveryPipeline, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/deliveryPipelines/%s", projectID, location, pipelineName)

	pipeline, err := service.Projects.Locations.DeliveryPipelines.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("delivery pipeline %s does not exist in project %s", pipelineName, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for delivery pipeline %s in project %s: %w", pipelineName, projectID, err)
	}

	return pipeline, nil
}

// NewCloudDeployServiceE creates a Cloud Deploy service authenticated the same way every other
// client in this module is.
// The ctx parameter supports cancellation and timeouts.
func NewCloudDeployServiceE(t testing.TestingT, ctx context.Context) (*clouddeploy.Service, error) {
	return clouddeploy.NewService(ctx, append(withOptions(), option.WithScopes(clouddeploy.CloudPlatformScope))...)
}
