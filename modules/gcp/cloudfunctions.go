package gcp

import (
	"context"
	"errors"
	"fmt"

	"github.com/gruntwork-io/terratest/modules/core/v2/logger"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
	cloudfunctionsv1 "google.golang.org/api/cloudfunctions/v1"
	"google.golang.org/api/cloudfunctions/v2"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
)

// GetCloudFunctionV2Attrs returns the settings Google Cloud holds for the given Cloud Function, so a
// test can assert on what was actually created rather than only that it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetCloudFunctionV2Attrs(t testing.TestingT, ctx context.Context, projectID string, region string, functionID string) *cloudfunctions.Function {
	function, err := GetCloudFunctionAttrsE(t, ctx, projectID, region, functionID)
	require.NoError(t, err)

	return function
}

// GetCloudFunctionAttrsE returns the settings Google Cloud holds for the given Cloud Function.
// The ctx parameter supports cancellation and timeouts.
func GetCloudFunctionAttrsE(t testing.TestingT, ctx context.Context, projectID string, region string, functionID string) (*cloudfunctions.Function, error) {
	logger.Default.Logf(t, "Getting settings for Cloud Function %s in region %s in project %s", functionID, region, projectID)

	service, err := NewCloudFunctionsServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetCloudFunctionAttrsWithClient(ctx, service, projectID, region, functionID)
}

// GetCloudFunctionAttrsWithClient returns the settings Google Cloud holds for the given Cloud
// Function using the supplied *cloudfunctions.Service. Prefer this variant in unit tests where the
// service is backed by an httptest fake server (see cloudfunctions_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetCloudFunctionAttrsWithClient(ctx context.Context, service *cloudfunctions.Service, projectID string, region string, functionID string) (*cloudfunctions.Function, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/functions/%s", projectID, region, functionID)

	function, err := service.Projects.Locations.Functions.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Cloud Function %s does not exist in region %s in project %s", functionID, region, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Cloud Function %s in region %s in project %s: %w", functionID, region, projectID, err)
	}

	return function, nil
}

// NewCloudFunctionsServiceE creates a Cloud Functions service authenticated the same way every
// other client in this module is.
// The ctx parameter supports cancellation and timeouts.
func NewCloudFunctionsServiceE(t testing.TestingT, ctx context.Context) (*cloudfunctions.Service, error) {
	return cloudfunctions.NewService(ctx, append(withOptions(), option.WithScopes(cloudfunctions.CloudPlatformScope))...)
}

// GetCloudFunctionV1Attrs returns the settings Google Cloud holds for the given first generation
// Cloud Function, so a test can assert on what the older module created. The two generations are
// separate API surfaces, and a function created by one is not visible to the other.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetCloudFunctionV1Attrs(t testing.TestingT, ctx context.Context, projectID string, location string, functionID string) *cloudfunctionsv1.CloudFunction {
	function, err := GetCloudFunctionV1AttrsE(t, ctx, projectID, location, functionID)
	require.NoError(t, err)

	return function
}

// GetCloudFunctionV1AttrsE returns the settings Google Cloud holds for the given first generation
// Cloud Function.
// The ctx parameter supports cancellation and timeouts.
func GetCloudFunctionV1AttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, functionID string) (*cloudfunctionsv1.CloudFunction, error) {
	logger.Default.Logf(t, "Getting settings for first generation Cloud Function %s in %s in project %s", functionID, location, projectID)

	service, err := NewCloudFunctionsV1ServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetCloudFunctionV1AttrsWithClient(ctx, service, projectID, location, functionID)
}

// GetCloudFunctionV1AttrsWithClient returns the settings Google Cloud holds for the given first
// generation Cloud Function using the supplied *cloudfunctionsv1.Service. Prefer this variant in
// unit tests where the service is backed by an httptest fake server (see cloudfunctions_test.go for
// the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetCloudFunctionV1AttrsWithClient(ctx context.Context, service *cloudfunctionsv1.Service, projectID string, location string, functionID string) (*cloudfunctionsv1.CloudFunction, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/functions/%s", projectID, location, functionID)

	function, err := service.Projects.Locations.Functions.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the first generation Cloud Function %s in %s in project %s does not exist", functionID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for first generation Cloud Function %s in %s in project %s: %w", functionID, location, projectID, err)
	}

	return function, nil
}

// NewCloudFunctionsV1ServiceE creates a first generation Cloud Functions service authenticated the
// same way every other client in this module is.
// The ctx parameter supports cancellation and timeouts.
func NewCloudFunctionsV1ServiceE(t testing.TestingT, ctx context.Context) (*cloudfunctionsv1.Service, error) {
	return cloudfunctionsv1.NewService(ctx, append(withOptions(), option.WithScopes(cloudfunctionsv1.CloudPlatformScope))...)
}
