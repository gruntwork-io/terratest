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

// GetVertexAIEndpointAttrs returns the settings Google Cloud holds for the Vertex AI endpoint, so a test can assert on what was
// actually created rather than only that it exists.
// A Vertex AI endpoint is where a model is served from, and its traffic split is the only place the
// share each deployed model receives is recorded.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetVertexAIEndpointAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, endpointID string) *aiplatform.GoogleCloudAiplatformV1Endpoint {
	result, err := GetVertexAIEndpointAttrsE(t, ctx, projectID, location, endpointID)
	require.NoError(t, err)

	return result
}

// GetVertexAIEndpointAttrsE returns the settings Google Cloud holds for the Vertex AI endpoint.
// The ctx parameter supports cancellation and timeouts.
func GetVertexAIEndpointAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, endpointID string) (*aiplatform.GoogleCloudAiplatformV1Endpoint, error) {
	logger.Default.Logf(t, "Getting settings for Vertex AI endpoint %s in %s in project %s", endpointID, location, projectID)

	service, err := NewVertexAIServiceE(t, ctx, location)
	if err != nil {
		return nil, err
	}

	return GetVertexAIEndpointAttrsWithClient(ctx, service, projectID, location, endpointID)
}

// GetVertexAIEndpointAttrsWithClient returns the settings Google Cloud holds for the Vertex AI endpoint using the supplied
// *aiplatform.Service. Prefer this variant in unit tests where the service is backed by an httptest
// fake server (see aiplatform_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetVertexAIEndpointAttrsWithClient(ctx context.Context, service *aiplatform.Service, projectID string, location string, endpointID string) (*aiplatform.GoogleCloudAiplatformV1Endpoint, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/endpoints/%s", projectID, location, endpointID)

	result, err := service.Projects.Locations.Endpoints.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Vertex AI endpoint %s in %s in project %s does not exist", endpointID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Vertex AI endpoint %s in %s in project %s: %w", endpointID, location, projectID, err)
	}

	return result, nil
}

// GetVertexAIIndexAttrs returns the settings Google Cloud holds for the Vertex AI index, so a test can assert on what was
// actually created rather than only that it exists.
// The index holds the vectors a search runs against, and Google reports the shape it settled on.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetVertexAIIndexAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, indexID string) *aiplatform.GoogleCloudAiplatformV1Index {
	result, err := GetVertexAIIndexAttrsE(t, ctx, projectID, location, indexID)
	require.NoError(t, err)

	return result
}

// GetVertexAIIndexAttrsE returns the settings Google Cloud holds for the Vertex AI index.
// The ctx parameter supports cancellation and timeouts.
func GetVertexAIIndexAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, indexID string) (*aiplatform.GoogleCloudAiplatformV1Index, error) {
	logger.Default.Logf(t, "Getting settings for Vertex AI index %s in %s in project %s", indexID, location, projectID)

	service, err := NewVertexAIServiceE(t, ctx, location)
	if err != nil {
		return nil, err
	}

	return GetVertexAIIndexAttrsWithClient(ctx, service, projectID, location, indexID)
}

// GetVertexAIIndexAttrsWithClient returns the settings Google Cloud holds for the Vertex AI index using the supplied
// *aiplatform.Service. Prefer this variant in unit tests where the service is backed by an httptest
// fake server (see aiplatform_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetVertexAIIndexAttrsWithClient(ctx context.Context, service *aiplatform.Service, projectID string, location string, indexID string) (*aiplatform.GoogleCloudAiplatformV1Index, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/indexes/%s", projectID, location, indexID)

	result, err := service.Projects.Locations.Indexes.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Vertex AI index %s in %s in project %s does not exist", indexID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Vertex AI index %s in %s in project %s: %w", indexID, location, projectID, err)
	}

	return result, nil
}

// GetVertexAIIndexEndpointAttrs returns the settings Google Cloud holds for the Vertex AI index endpoint, so a test can assert on what was
// actually created rather than only that it exists.
// The deployed indexes an index endpoint serves are reported on the endpoint rather than on
// themselves, so this read is also how a test asserts on a deployed index.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetVertexAIIndexEndpointAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, indexEndpointID string) *aiplatform.GoogleCloudAiplatformV1IndexEndpoint {
	result, err := GetVertexAIIndexEndpointAttrsE(t, ctx, projectID, location, indexEndpointID)
	require.NoError(t, err)

	return result
}

// GetVertexAIIndexEndpointAttrsE returns the settings Google Cloud holds for the Vertex AI index endpoint.
// The ctx parameter supports cancellation and timeouts.
func GetVertexAIIndexEndpointAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, indexEndpointID string) (*aiplatform.GoogleCloudAiplatformV1IndexEndpoint, error) {
	logger.Default.Logf(t, "Getting settings for Vertex AI index endpoint %s in %s in project %s", indexEndpointID, location, projectID)

	service, err := NewVertexAIServiceE(t, ctx, location)
	if err != nil {
		return nil, err
	}

	return GetVertexAIIndexEndpointAttrsWithClient(ctx, service, projectID, location, indexEndpointID)
}

// GetVertexAIIndexEndpointAttrsWithClient returns the settings Google Cloud holds for the Vertex AI index endpoint using the supplied
// *aiplatform.Service. Prefer this variant in unit tests where the service is backed by an httptest
// fake server (see aiplatform_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetVertexAIIndexEndpointAttrsWithClient(ctx context.Context, service *aiplatform.Service, projectID string, location string, indexEndpointID string) (*aiplatform.GoogleCloudAiplatformV1IndexEndpoint, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/indexEndpoints/%s", projectID, location, indexEndpointID)

	result, err := service.Projects.Locations.IndexEndpoints.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Vertex AI index endpoint %s in %s in project %s does not exist", indexEndpointID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Vertex AI index endpoint %s in %s in project %s: %w", indexEndpointID, location, projectID, err)
	}

	return result, nil
}

// GetVertexAIFeatureOnlineStoreAttrs returns the settings Google Cloud holds for the Vertex AI feature online store, so a test can assert on what was
// actually created rather than only that it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetVertexAIFeatureOnlineStoreAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, storeID string) *aiplatform.GoogleCloudAiplatformV1FeatureOnlineStore {
	result, err := GetVertexAIFeatureOnlineStoreAttrsE(t, ctx, projectID, location, storeID)
	require.NoError(t, err)

	return result
}

// GetVertexAIFeatureOnlineStoreAttrsE returns the settings Google Cloud holds for the Vertex AI feature online store.
// The ctx parameter supports cancellation and timeouts.
func GetVertexAIFeatureOnlineStoreAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, storeID string) (*aiplatform.GoogleCloudAiplatformV1FeatureOnlineStore, error) {
	logger.Default.Logf(t, "Getting settings for Vertex AI feature online store %s in %s in project %s", storeID, location, projectID)

	service, err := NewVertexAIServiceE(t, ctx, location)
	if err != nil {
		return nil, err
	}

	return GetVertexAIFeatureOnlineStoreAttrsWithClient(ctx, service, projectID, location, storeID)
}

// GetVertexAIFeatureOnlineStoreAttrsWithClient returns the settings Google Cloud holds for the Vertex AI feature online store using the supplied
// *aiplatform.Service. Prefer this variant in unit tests where the service is backed by an httptest
// fake server (see aiplatform_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetVertexAIFeatureOnlineStoreAttrsWithClient(ctx context.Context, service *aiplatform.Service, projectID string, location string, storeID string) (*aiplatform.GoogleCloudAiplatformV1FeatureOnlineStore, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/featureOnlineStores/%s", projectID, location, storeID)

	result, err := service.Projects.Locations.FeatureOnlineStores.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Vertex AI feature online store %s in %s in project %s does not exist", storeID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Vertex AI feature online store %s in %s in project %s: %w", storeID, location, projectID, err)
	}

	return result, nil
}

// GetVertexAIDeploymentResourcePoolAttrs returns the settings Google Cloud holds for the Vertex AI deployment resource pool, so a test can assert on what was
// actually created rather than only that it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetVertexAIDeploymentResourcePoolAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, poolID string) *aiplatform.GoogleCloudAiplatformV1DeploymentResourcePool {
	result, err := GetVertexAIDeploymentResourcePoolAttrsE(t, ctx, projectID, location, poolID)
	require.NoError(t, err)

	return result
}

// GetVertexAIDeploymentResourcePoolAttrsE returns the settings Google Cloud holds for the Vertex AI deployment resource pool.
// The ctx parameter supports cancellation and timeouts.
func GetVertexAIDeploymentResourcePoolAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, poolID string) (*aiplatform.GoogleCloudAiplatformV1DeploymentResourcePool, error) {
	logger.Default.Logf(t, "Getting settings for Vertex AI deployment resource pool %s in %s in project %s", poolID, location, projectID)

	service, err := NewVertexAIServiceE(t, ctx, location)
	if err != nil {
		return nil, err
	}

	return GetVertexAIDeploymentResourcePoolAttrsWithClient(ctx, service, projectID, location, poolID)
}

// GetVertexAIDeploymentResourcePoolAttrsWithClient returns the settings Google Cloud holds for the Vertex AI deployment resource pool using the supplied
// *aiplatform.Service. Prefer this variant in unit tests where the service is backed by an httptest
// fake server (see aiplatform_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetVertexAIDeploymentResourcePoolAttrsWithClient(ctx context.Context, service *aiplatform.Service, projectID string, location string, poolID string) (*aiplatform.GoogleCloudAiplatformV1DeploymentResourcePool, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/deploymentResourcePools/%s", projectID, location, poolID)

	result, err := service.Projects.Locations.DeploymentResourcePools.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Vertex AI deployment resource pool %s in %s in project %s does not exist", poolID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Vertex AI deployment resource pool %s in %s in project %s: %w", poolID, location, projectID, err)
	}

	return result, nil
}

// GetVertexAIPersistentResourceAttrs returns the settings Google Cloud holds for the Vertex AI persistent resource, so a test can assert on what was
// actually created rather than only that it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetVertexAIPersistentResourceAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, resourceID string) *aiplatform.GoogleCloudAiplatformV1PersistentResource {
	result, err := GetVertexAIPersistentResourceAttrsE(t, ctx, projectID, location, resourceID)
	require.NoError(t, err)

	return result
}

// GetVertexAIPersistentResourceAttrsE returns the settings Google Cloud holds for the Vertex AI persistent resource.
// The ctx parameter supports cancellation and timeouts.
func GetVertexAIPersistentResourceAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, resourceID string) (*aiplatform.GoogleCloudAiplatformV1PersistentResource, error) {
	logger.Default.Logf(t, "Getting settings for Vertex AI persistent resource %s in %s in project %s", resourceID, location, projectID)

	service, err := NewVertexAIServiceE(t, ctx, location)
	if err != nil {
		return nil, err
	}

	return GetVertexAIPersistentResourceAttrsWithClient(ctx, service, projectID, location, resourceID)
}

// GetVertexAIPersistentResourceAttrsWithClient returns the settings Google Cloud holds for the Vertex AI persistent resource using the supplied
// *aiplatform.Service. Prefer this variant in unit tests where the service is backed by an httptest
// fake server (see aiplatform_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetVertexAIPersistentResourceAttrsWithClient(ctx context.Context, service *aiplatform.Service, projectID string, location string, resourceID string) (*aiplatform.GoogleCloudAiplatformV1PersistentResource, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/persistentResources/%s", projectID, location, resourceID)

	result, err := service.Projects.Locations.PersistentResources.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Vertex AI persistent resource %s in %s in project %s does not exist", resourceID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Vertex AI persistent resource %s in %s in project %s: %w", resourceID, location, projectID, err)
	}

	return result, nil
}

// GetVertexAIScheduleAttrs returns the settings Google Cloud holds for the Vertex AI schedule, so a test can assert on what was
// actually created rather than only that it exists.
// The cron expression and the next run time are what say the schedule is actually armed.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetVertexAIScheduleAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, scheduleID string) *aiplatform.GoogleCloudAiplatformV1Schedule {
	result, err := GetVertexAIScheduleAttrsE(t, ctx, projectID, location, scheduleID)
	require.NoError(t, err)

	return result
}

// GetVertexAIScheduleAttrsE returns the settings Google Cloud holds for the Vertex AI schedule.
// The ctx parameter supports cancellation and timeouts.
func GetVertexAIScheduleAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, scheduleID string) (*aiplatform.GoogleCloudAiplatformV1Schedule, error) {
	logger.Default.Logf(t, "Getting settings for Vertex AI schedule %s in %s in project %s", scheduleID, location, projectID)

	service, err := NewVertexAIServiceE(t, ctx, location)
	if err != nil {
		return nil, err
	}

	return GetVertexAIScheduleAttrsWithClient(ctx, service, projectID, location, scheduleID)
}

// GetVertexAIScheduleAttrsWithClient returns the settings Google Cloud holds for the Vertex AI schedule using the supplied
// *aiplatform.Service. Prefer this variant in unit tests where the service is backed by an httptest
// fake server (see aiplatform_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetVertexAIScheduleAttrsWithClient(ctx context.Context, service *aiplatform.Service, projectID string, location string, scheduleID string) (*aiplatform.GoogleCloudAiplatformV1Schedule, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/schedules/%s", projectID, location, scheduleID)

	result, err := service.Projects.Locations.Schedules.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Vertex AI schedule %s in %s in project %s does not exist", scheduleID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Vertex AI schedule %s in %s in project %s: %w", scheduleID, location, projectID, err)
	}

	return result, nil
}

// GetVertexAIReasoningEngineAttrs returns the settings Google Cloud holds for the Vertex AI reasoning engine, so a test can assert on what was
// actually created rather than only that it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetVertexAIReasoningEngineAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, engineID string) *aiplatform.GoogleCloudAiplatformV1ReasoningEngine {
	result, err := GetVertexAIReasoningEngineAttrsE(t, ctx, projectID, location, engineID)
	require.NoError(t, err)

	return result
}

// GetVertexAIReasoningEngineAttrsE returns the settings Google Cloud holds for the Vertex AI reasoning engine.
// The ctx parameter supports cancellation and timeouts.
func GetVertexAIReasoningEngineAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, engineID string) (*aiplatform.GoogleCloudAiplatformV1ReasoningEngine, error) {
	logger.Default.Logf(t, "Getting settings for Vertex AI reasoning engine %s in %s in project %s", engineID, location, projectID)

	service, err := NewVertexAIServiceE(t, ctx, location)
	if err != nil {
		return nil, err
	}

	return GetVertexAIReasoningEngineAttrsWithClient(ctx, service, projectID, location, engineID)
}

// GetVertexAIReasoningEngineAttrsWithClient returns the settings Google Cloud holds for the Vertex AI reasoning engine using the supplied
// *aiplatform.Service. Prefer this variant in unit tests where the service is backed by an httptest
// fake server (see aiplatform_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetVertexAIReasoningEngineAttrsWithClient(ctx context.Context, service *aiplatform.Service, projectID string, location string, engineID string) (*aiplatform.GoogleCloudAiplatformV1ReasoningEngine, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/reasoningEngines/%s", projectID, location, engineID)

	result, err := service.Projects.Locations.ReasoningEngines.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Vertex AI reasoning engine %s in %s in project %s does not exist", engineID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Vertex AI reasoning engine %s in %s in project %s: %w", engineID, location, projectID, err)
	}

	return result, nil
}

// GetVertexAIFeatureViewAttrs returns the settings Google Cloud holds for the Vertex AI feature view, so a test can assert on what was
// actually created rather than only that it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetVertexAIFeatureViewAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, storeID string, viewID string) *aiplatform.GoogleCloudAiplatformV1FeatureView {
	result, err := GetVertexAIFeatureViewAttrsE(t, ctx, projectID, location, storeID, viewID)
	require.NoError(t, err)

	return result
}

// GetVertexAIFeatureViewAttrsE returns the settings Google Cloud holds for the Vertex AI feature view.
// The ctx parameter supports cancellation and timeouts.
func GetVertexAIFeatureViewAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, storeID string, viewID string) (*aiplatform.GoogleCloudAiplatformV1FeatureView, error) {
	logger.Default.Logf(t, "Getting settings for Vertex AI feature view %s of feature online store %s in %s in project %s", viewID, storeID, location, projectID)

	service, err := NewVertexAIServiceE(t, ctx, location)
	if err != nil {
		return nil, err
	}

	return GetVertexAIFeatureViewAttrsWithClient(ctx, service, projectID, location, storeID, viewID)
}

// GetVertexAIFeatureViewAttrsWithClient returns the settings Google Cloud holds for the Vertex AI feature view using the supplied
// *aiplatform.Service. Prefer this variant in unit tests where the service is backed by an httptest
// fake server (see aiplatform_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetVertexAIFeatureViewAttrsWithClient(ctx context.Context, service *aiplatform.Service, projectID string, location string, storeID string, viewID string) (*aiplatform.GoogleCloudAiplatformV1FeatureView, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/featureOnlineStores/%s/featureViews/%s", projectID, location, storeID, viewID)

	result, err := service.Projects.Locations.FeatureOnlineStores.FeatureViews.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Vertex AI feature view %s of feature online store %s in %s in project %s does not exist", viewID, storeID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Vertex AI feature view %s of feature online store %s in %s in project %s: %w", viewID, storeID, location, projectID, err)
	}

	return result, nil
}

// GetVertexAIFeatureGroupFeatureAttrs returns the settings Google Cloud holds for the Vertex AI feature group feature, so a test can assert on what was
// actually created rather than only that it exists.
// This is the feature a feature group reads out of a BigQuery column, which is a different resource
// from the feature a featurestore entity type holds: GetVertexAIFeatureAttrs reads that older one.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetVertexAIFeatureGroupFeatureAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, featureGroupID string, featureID string) *aiplatform.GoogleCloudAiplatformV1Feature {
	result, err := GetVertexAIFeatureGroupFeatureAttrsE(t, ctx, projectID, location, featureGroupID, featureID)
	require.NoError(t, err)

	return result
}

// GetVertexAIFeatureGroupFeatureAttrsE returns the settings Google Cloud holds for the Vertex AI feature group feature.
// The ctx parameter supports cancellation and timeouts.
func GetVertexAIFeatureGroupFeatureAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, featureGroupID string, featureID string) (*aiplatform.GoogleCloudAiplatformV1Feature, error) {
	logger.Default.Logf(t, "Getting settings for Vertex AI feature group feature %s of feature group %s in %s in project %s", featureID, featureGroupID, location, projectID)

	service, err := NewVertexAIServiceE(t, ctx, location)
	if err != nil {
		return nil, err
	}

	return GetVertexAIFeatureGroupFeatureAttrsWithClient(ctx, service, projectID, location, featureGroupID, featureID)
}

// GetVertexAIFeatureGroupFeatureAttrsWithClient returns the settings Google Cloud holds for the Vertex AI feature group feature using the supplied
// *aiplatform.Service. Prefer this variant in unit tests where the service is backed by an httptest
// fake server (see aiplatform_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetVertexAIFeatureGroupFeatureAttrsWithClient(ctx context.Context, service *aiplatform.Service, projectID string, location string, featureGroupID string, featureID string) (*aiplatform.GoogleCloudAiplatformV1Feature, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/featureGroups/%s/features/%s", projectID, location, featureGroupID, featureID)

	result, err := service.Projects.Locations.FeatureGroups.Features.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Vertex AI feature group feature %s of feature group %s in %s in project %s does not exist", featureID, featureGroupID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Vertex AI feature group feature %s of feature group %s in %s in project %s: %w", featureID, featureGroupID, location, projectID, err)
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
