package gcp

import (
	"context"
	"errors"
	"fmt"

	"github.com/gruntwork-io/terratest/modules/core/v2/logger"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/ml/v1"
	"google.golang.org/api/option"
)

// GetMLEngineModelAttrs returns the settings Google Cloud holds for the given AI Platform model, so a test can assert on
// what was actually created rather than only that it exists.
// The default version Google promoted is reported on the model rather than set by the caller, so a
// read is how a test sees which version serves traffic.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetMLEngineModelAttrs(t testing.TestingT, ctx context.Context, projectID string, modelID string) *ml.GoogleCloudMlV1__Model {
	result, err := GetMLEngineModelAttrsE(t, ctx, projectID, modelID)
	require.NoError(t, err)

	return result
}

// GetMLEngineModelAttrsE returns the settings Google Cloud holds for the given AI Platform model.
// The ctx parameter supports cancellation and timeouts.
func GetMLEngineModelAttrsE(t testing.TestingT, ctx context.Context, projectID string, modelID string) (*ml.GoogleCloudMlV1__Model, error) {
	logger.Default.Logf(t, "Getting settings for AI Platform model %s in project %s", modelID, projectID)

	service, err := NewMLEngineServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetMLEngineModelAttrsWithClient(ctx, service, projectID, modelID)
}

// GetMLEngineModelAttrsWithClient returns the settings Google Cloud holds for the given AI Platform model using the supplied
// *ml.Service. Prefer this variant in unit tests where the service is backed by an httptest fake server
// (see mlengine_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetMLEngineModelAttrsWithClient(ctx context.Context, service *ml.Service, projectID string, modelID string) (*ml.GoogleCloudMlV1__Model, error) {
	name := fmt.Sprintf("projects/%s/models/%s", projectID, modelID)

	result, err := service.Projects.Models.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the AI Platform model %s does not exist in project %s", modelID, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for AI Platform model %s in project %s: %w", modelID, projectID, err)
	}

	return result, nil
}

// NewMLEngineServiceE creates a AI Platform Training and Prediction service authenticated the same way every other client in
// this module is.
// The ctx parameter supports cancellation and timeouts.
func NewMLEngineServiceE(t testing.TestingT, ctx context.Context) (*ml.Service, error) {
	return ml.NewService(ctx, append(withOptions(), option.WithScopes(ml.CloudPlatformScope))...)
}
