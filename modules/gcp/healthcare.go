package gcp

import (
	"context"
	"errors"
	"fmt"

	"github.com/gruntwork-io/terratest/modules/core/v2/logger"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/healthcare/v1"
	"google.golang.org/api/option"
)

// GetHealthcareDatasetAttrs returns the settings Google Cloud holds for the given Cloud Healthcare
// dataset, so a test can assert on what was actually created rather than only that it exists. A
// dataset is regional, so the location it was created in has to be given.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetHealthcareDatasetAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, datasetID string) *healthcare.Dataset {
	dataset, err := GetHealthcareDatasetAttrsE(t, ctx, projectID, location, datasetID)
	require.NoError(t, err)

	return dataset
}

// GetHealthcareDatasetAttrsE returns the settings Google Cloud holds for the given Cloud Healthcare
// dataset.
// The ctx parameter supports cancellation and timeouts.
func GetHealthcareDatasetAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, datasetID string) (*healthcare.Dataset, error) {
	logger.Default.Logf(t, "Getting settings for healthcare dataset %s in project %s", datasetID, projectID)

	service, err := NewHealthcareServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetHealthcareDatasetAttrsWithClient(ctx, service, projectID, location, datasetID)
}

// GetHealthcareDatasetAttrsWithClient returns the settings Google Cloud holds for the given Cloud
// Healthcare dataset using the supplied *healthcare.Service. Prefer this variant in unit tests
// where the service is backed by an httptest fake server (see healthcare_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetHealthcareDatasetAttrsWithClient(ctx context.Context, service *healthcare.Service, projectID string, location string, datasetID string) (*healthcare.Dataset, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/datasets/%s", projectID, location, datasetID)

	dataset, err := service.Projects.Locations.Datasets.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("healthcare dataset %s does not exist in project %s", datasetID, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for healthcare dataset %s in project %s: %w", datasetID, projectID, err)
	}

	return dataset, nil
}

// NewHealthcareServiceE creates a Cloud Healthcare service authenticated the same way every other
// client in this module is.
// The ctx parameter supports cancellation and timeouts.
func NewHealthcareServiceE(t testing.TestingT, ctx context.Context) (*healthcare.Service, error) {
	return healthcare.NewService(ctx, append(withOptions(), option.WithScopes(healthcare.CloudPlatformScope))...)
}
