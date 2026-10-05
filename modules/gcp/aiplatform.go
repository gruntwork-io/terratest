package gcp

import (
	"context"
	"errors"
	"fmt"
	"regexp"

	"github.com/gruntwork-io/terratest/modules/core/v2/logger"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/aiplatform/v1"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
)

// GetVertexAIDatasetAttrs returns the settings Google Cloud holds for the Vertex AI dataset, so a test can assert on what was
// actually created rather than only that it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetVertexAIDatasetAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, datasetID string) *aiplatform.GoogleCloudAiplatformV1Dataset {
	result, err := GetVertexAIDatasetAttrsE(t, ctx, projectID, location, datasetID)
	require.NoError(t, err)

	return result
}

// GetVertexAIDatasetAttrsE returns the settings Google Cloud holds for the Vertex AI dataset.
// The ctx parameter supports cancellation and timeouts.
func GetVertexAIDatasetAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, datasetID string) (*aiplatform.GoogleCloudAiplatformV1Dataset, error) {
	logger.Default.Logf(t, "Getting settings for Vertex AI dataset %s in %s in project %s", datasetID, location, projectID)

	service, err := NewVertexAIServiceE(t, ctx, location)
	if err != nil {
		return nil, err
	}

	return GetVertexAIDatasetAttrsWithClient(ctx, service, projectID, location, datasetID)
}

// GetVertexAIDatasetAttrsWithClient returns the settings Google Cloud holds for the Vertex AI dataset using the supplied
// *aiplatform.Service. Prefer this variant in unit tests where the service is backed by an httptest
// fake server (see aiplatform_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetVertexAIDatasetAttrsWithClient(ctx context.Context, service *aiplatform.Service, projectID string, location string, datasetID string) (*aiplatform.GoogleCloudAiplatformV1Dataset, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/datasets/%s", projectID, location, datasetID)

	result, err := service.Projects.Locations.Datasets.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Vertex AI dataset %s in %s in project %s does not exist", datasetID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Vertex AI dataset %s in %s in project %s: %w", datasetID, location, projectID, err)
	}

	return result, nil
}

// GetVertexAITensorboardAttrs returns the settings Google Cloud holds for the Vertex AI tensorboard, so a test can assert on what was
// actually created rather than only that it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetVertexAITensorboardAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, tensorboardID string) *aiplatform.GoogleCloudAiplatformV1Tensorboard {
	result, err := GetVertexAITensorboardAttrsE(t, ctx, projectID, location, tensorboardID)
	require.NoError(t, err)

	return result
}

// GetVertexAITensorboardAttrsE returns the settings Google Cloud holds for the Vertex AI tensorboard.
// The ctx parameter supports cancellation and timeouts.
func GetVertexAITensorboardAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, tensorboardID string) (*aiplatform.GoogleCloudAiplatformV1Tensorboard, error) {
	logger.Default.Logf(t, "Getting settings for Vertex AI tensorboard %s in %s in project %s", tensorboardID, location, projectID)

	service, err := NewVertexAIServiceE(t, ctx, location)
	if err != nil {
		return nil, err
	}

	return GetVertexAITensorboardAttrsWithClient(ctx, service, projectID, location, tensorboardID)
}

// GetVertexAITensorboardAttrsWithClient returns the settings Google Cloud holds for the Vertex AI tensorboard using the supplied
// *aiplatform.Service. Prefer this variant in unit tests where the service is backed by an httptest
// fake server (see aiplatform_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetVertexAITensorboardAttrsWithClient(ctx context.Context, service *aiplatform.Service, projectID string, location string, tensorboardID string) (*aiplatform.GoogleCloudAiplatformV1Tensorboard, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/tensorboards/%s", projectID, location, tensorboardID)

	result, err := service.Projects.Locations.Tensorboards.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Vertex AI tensorboard %s in %s in project %s does not exist", tensorboardID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Vertex AI tensorboard %s in %s in project %s: %w", tensorboardID, location, projectID, err)
	}

	return result, nil
}

// GetVertexAIFeaturestoreAttrs returns the settings Google Cloud holds for the Vertex AI featurestore, so a test can assert on what was
// actually created rather than only that it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetVertexAIFeaturestoreAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, featurestoreID string) *aiplatform.GoogleCloudAiplatformV1Featurestore {
	result, err := GetVertexAIFeaturestoreAttrsE(t, ctx, projectID, location, featurestoreID)
	require.NoError(t, err)

	return result
}

// GetVertexAIFeaturestoreAttrsE returns the settings Google Cloud holds for the Vertex AI featurestore.
// The ctx parameter supports cancellation and timeouts.
func GetVertexAIFeaturestoreAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, featurestoreID string) (*aiplatform.GoogleCloudAiplatformV1Featurestore, error) {
	logger.Default.Logf(t, "Getting settings for Vertex AI featurestore %s in %s in project %s", featurestoreID, location, projectID)

	service, err := NewVertexAIServiceE(t, ctx, location)
	if err != nil {
		return nil, err
	}

	return GetVertexAIFeaturestoreAttrsWithClient(ctx, service, projectID, location, featurestoreID)
}

// GetVertexAIFeaturestoreAttrsWithClient returns the settings Google Cloud holds for the Vertex AI featurestore using the supplied
// *aiplatform.Service. Prefer this variant in unit tests where the service is backed by an httptest
// fake server (see aiplatform_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetVertexAIFeaturestoreAttrsWithClient(ctx context.Context, service *aiplatform.Service, projectID string, location string, featurestoreID string) (*aiplatform.GoogleCloudAiplatformV1Featurestore, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/featurestores/%s", projectID, location, featurestoreID)

	result, err := service.Projects.Locations.Featurestores.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Vertex AI featurestore %s in %s in project %s does not exist", featurestoreID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Vertex AI featurestore %s in %s in project %s: %w", featurestoreID, location, projectID, err)
	}

	return result, nil
}

// GetVertexAIFeatureGroupAttrs returns the settings Google Cloud holds for the Vertex AI feature group, so a test can assert on what was
// actually created rather than only that it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetVertexAIFeatureGroupAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, groupID string) *aiplatform.GoogleCloudAiplatformV1FeatureGroup {
	result, err := GetVertexAIFeatureGroupAttrsE(t, ctx, projectID, location, groupID)
	require.NoError(t, err)

	return result
}

// GetVertexAIFeatureGroupAttrsE returns the settings Google Cloud holds for the Vertex AI feature group.
// The ctx parameter supports cancellation and timeouts.
func GetVertexAIFeatureGroupAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, groupID string) (*aiplatform.GoogleCloudAiplatformV1FeatureGroup, error) {
	logger.Default.Logf(t, "Getting settings for Vertex AI feature group %s in %s in project %s", groupID, location, projectID)

	service, err := NewVertexAIServiceE(t, ctx, location)
	if err != nil {
		return nil, err
	}

	return GetVertexAIFeatureGroupAttrsWithClient(ctx, service, projectID, location, groupID)
}

// GetVertexAIFeatureGroupAttrsWithClient returns the settings Google Cloud holds for the Vertex AI feature group using the supplied
// *aiplatform.Service. Prefer this variant in unit tests where the service is backed by an httptest
// fake server (see aiplatform_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetVertexAIFeatureGroupAttrsWithClient(ctx context.Context, service *aiplatform.Service, projectID string, location string, groupID string) (*aiplatform.GoogleCloudAiplatformV1FeatureGroup, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/featureGroups/%s", projectID, location, groupID)

	result, err := service.Projects.Locations.FeatureGroups.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Vertex AI feature group %s in %s in project %s does not exist", groupID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Vertex AI feature group %s in %s in project %s: %w", groupID, location, projectID, err)
	}

	return result, nil
}

// GetVertexAINotebookRuntimeTemplateAttrs returns the settings Google Cloud holds for the Vertex AI notebook runtime template, so a test can assert on what was
// actually created rather than only that it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetVertexAINotebookRuntimeTemplateAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, templateID string) *aiplatform.GoogleCloudAiplatformV1NotebookRuntimeTemplate {
	result, err := GetVertexAINotebookRuntimeTemplateAttrsE(t, ctx, projectID, location, templateID)
	require.NoError(t, err)

	return result
}

// GetVertexAINotebookRuntimeTemplateAttrsE returns the settings Google Cloud holds for the Vertex AI notebook runtime template.
// The ctx parameter supports cancellation and timeouts.
func GetVertexAINotebookRuntimeTemplateAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, templateID string) (*aiplatform.GoogleCloudAiplatformV1NotebookRuntimeTemplate, error) {
	logger.Default.Logf(t, "Getting settings for Vertex AI notebook runtime template %s in %s in project %s", templateID, location, projectID)

	service, err := NewVertexAIServiceE(t, ctx, location)
	if err != nil {
		return nil, err
	}

	return GetVertexAINotebookRuntimeTemplateAttrsWithClient(ctx, service, projectID, location, templateID)
}

// GetVertexAINotebookRuntimeTemplateAttrsWithClient returns the settings Google Cloud holds for the Vertex AI notebook runtime template using the supplied
// *aiplatform.Service. Prefer this variant in unit tests where the service is backed by an httptest
// fake server (see aiplatform_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetVertexAINotebookRuntimeTemplateAttrsWithClient(ctx context.Context, service *aiplatform.Service, projectID string, location string, templateID string) (*aiplatform.GoogleCloudAiplatformV1NotebookRuntimeTemplate, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/notebookRuntimeTemplates/%s", projectID, location, templateID)

	result, err := service.Projects.Locations.NotebookRuntimeTemplates.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Vertex AI notebook runtime template %s in %s in project %s does not exist", templateID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Vertex AI notebook runtime template %s in %s in project %s: %w", templateID, location, projectID, err)
	}

	return result, nil
}

// GetVertexAITensorboardExperimentAttrs returns the settings Google Cloud holds for the Vertex AI tensorboard experiment, so a test can assert on what was
// actually created rather than only that it exists. An experiment belongs to a tensorboard, so the caller names both.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetVertexAITensorboardExperimentAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, tensorboardID string, experimentID string) *aiplatform.GoogleCloudAiplatformV1TensorboardExperiment {
	result, err := GetVertexAITensorboardExperimentAttrsE(t, ctx, projectID, location, tensorboardID, experimentID)
	require.NoError(t, err)

	return result
}

// GetVertexAITensorboardExperimentAttrsE returns the settings Google Cloud holds for the Vertex AI tensorboard experiment.
// The ctx parameter supports cancellation and timeouts.
func GetVertexAITensorboardExperimentAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, tensorboardID string, experimentID string) (*aiplatform.GoogleCloudAiplatformV1TensorboardExperiment, error) {
	logger.Default.Logf(t, "Getting settings for Vertex AI tensorboard experiment %s in %s in project %s", experimentID, location, projectID)

	service, err := NewVertexAIServiceE(t, ctx, location)
	if err != nil {
		return nil, err
	}

	return GetVertexAITensorboardExperimentAttrsWithClient(ctx, service, projectID, location, tensorboardID, experimentID)
}

// GetVertexAITensorboardExperimentAttrsWithClient returns the settings Google Cloud holds for the Vertex AI tensorboard experiment using the supplied
// *aiplatform.Service. Prefer this variant in unit tests where the service is backed by an httptest
// fake server (see aiplatform_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetVertexAITensorboardExperimentAttrsWithClient(ctx context.Context, service *aiplatform.Service, projectID string, location string, tensorboardID string, experimentID string) (*aiplatform.GoogleCloudAiplatformV1TensorboardExperiment, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/tensorboards/%s/experiments/%s", projectID, location, tensorboardID, experimentID)

	result, err := service.Projects.Locations.Tensorboards.Experiments.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Vertex AI tensorboard experiment %s in %s in project %s does not exist", experimentID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Vertex AI tensorboard experiment %s in %s in project %s: %w", experimentID, location, projectID, err)
	}

	return result, nil
}

// GetVertexAIEntityTypeAttrs returns the settings Google Cloud holds for the Vertex AI entity type, so a test can assert on what was
// actually created rather than only that it exists. An entity type belongs to a featurestore, so the caller names both.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetVertexAIEntityTypeAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, featurestoreID string, entityTypeID string) *aiplatform.GoogleCloudAiplatformV1EntityType {
	result, err := GetVertexAIEntityTypeAttrsE(t, ctx, projectID, location, featurestoreID, entityTypeID)
	require.NoError(t, err)

	return result
}

// GetVertexAIEntityTypeAttrsE returns the settings Google Cloud holds for the Vertex AI entity type.
// The ctx parameter supports cancellation and timeouts.
func GetVertexAIEntityTypeAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, featurestoreID string, entityTypeID string) (*aiplatform.GoogleCloudAiplatformV1EntityType, error) {
	logger.Default.Logf(t, "Getting settings for Vertex AI entity type %s in %s in project %s", entityTypeID, location, projectID)

	service, err := NewVertexAIServiceE(t, ctx, location)
	if err != nil {
		return nil, err
	}

	return GetVertexAIEntityTypeAttrsWithClient(ctx, service, projectID, location, featurestoreID, entityTypeID)
}

// GetVertexAIEntityTypeAttrsWithClient returns the settings Google Cloud holds for the Vertex AI entity type using the supplied
// *aiplatform.Service. Prefer this variant in unit tests where the service is backed by an httptest
// fake server (see aiplatform_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetVertexAIEntityTypeAttrsWithClient(ctx context.Context, service *aiplatform.Service, projectID string, location string, featurestoreID string, entityTypeID string) (*aiplatform.GoogleCloudAiplatformV1EntityType, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/featurestores/%s/entityTypes/%s", projectID, location, featurestoreID, entityTypeID)

	result, err := service.Projects.Locations.Featurestores.EntityTypes.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Vertex AI entity type %s in %s in project %s does not exist", entityTypeID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Vertex AI entity type %s in %s in project %s: %w", entityTypeID, location, projectID, err)
	}

	return result, nil
}

// GetVertexAIFeatureAttrs returns the settings Google Cloud holds for the Vertex AI feature, so a test can assert on what was
// actually created rather than only that it exists. A feature belongs to an entity type inside a featurestore, so the caller names all three.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetVertexAIFeatureAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, featurestoreID string, entityTypeID string, featureID string) *aiplatform.GoogleCloudAiplatformV1Feature {
	result, err := GetVertexAIFeatureAttrsE(t, ctx, projectID, location, featurestoreID, entityTypeID, featureID)
	require.NoError(t, err)

	return result
}

// GetVertexAIFeatureAttrsE returns the settings Google Cloud holds for the Vertex AI feature.
// The ctx parameter supports cancellation and timeouts.
func GetVertexAIFeatureAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, featurestoreID string, entityTypeID string, featureID string) (*aiplatform.GoogleCloudAiplatformV1Feature, error) {
	logger.Default.Logf(t, "Getting settings for Vertex AI feature %s in %s in project %s", featureID, location, projectID)

	service, err := NewVertexAIServiceE(t, ctx, location)
	if err != nil {
		return nil, err
	}

	return GetVertexAIFeatureAttrsWithClient(ctx, service, projectID, location, featurestoreID, entityTypeID, featureID)
}

// GetVertexAIFeatureAttrsWithClient returns the settings Google Cloud holds for the Vertex AI feature using the supplied
// *aiplatform.Service. Prefer this variant in unit tests where the service is backed by an httptest
// fake server (see aiplatform_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetVertexAIFeatureAttrsWithClient(ctx context.Context, service *aiplatform.Service, projectID string, location string, featurestoreID string, entityTypeID string, featureID string) (*aiplatform.GoogleCloudAiplatformV1Feature, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/featurestores/%s/entityTypes/%s/features/%s", projectID, location, featurestoreID, entityTypeID, featureID)

	result, err := service.Projects.Locations.Featurestores.EntityTypes.Features.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Vertex AI feature %s in %s in project %s does not exist", featureID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Vertex AI feature %s in %s in project %s: %w", featureID, location, projectID, err)
	}

	return result, nil
}

// NewVertexAIServiceE creates a Vertex AI service pointed at one region, authenticated the same way
// every other client in this module is. Vertex AI serves each region from its own endpoint, and a
// call to the global one finds nothing.
// The ctx parameter supports cancellation and timeouts.
func NewVertexAIServiceE(t testing.TestingT, ctx context.Context, location string) (*aiplatform.Service, error) {
	if !vertexAILocationPattern.MatchString(location) {
		return nil, fmt.Errorf("%q is not a valid location: a location may hold only lowercase letters, digits and hyphens", location)
	}

	endpoint := fmt.Sprintf("https://%s-aiplatform.googleapis.com/", location)

	return aiplatform.NewService(ctx, append(withOptions(),
		option.WithScopes(aiplatform.CloudPlatformScope), option.WithEndpoint(endpoint))...)
}

// vertexAILocationPattern is what a Vertex AI location may look like. The location goes into the
// endpoint host, so anything else would send the request somewhere the caller did not name.
var vertexAILocationPattern = regexp.MustCompile(`^[a-z0-9-]+$`)
