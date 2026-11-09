package gcp

import (
	"context"
	"errors"
	"fmt"

	"github.com/gruntwork-io/terratest/modules/core/v2/logger"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/datafusion/v1"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
)

// GetDataFusionInstanceAttrs returns the settings Google Cloud holds for the given Cloud Data Fusion instance, so a test can assert on
// what was actually created rather than only that it exists.
// The tenant project and the service endpoint Google allocates for an instance are reported here and
// nowhere in the create request.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetDataFusionInstanceAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, instanceID string) *datafusion.Instance {
	result, err := GetDataFusionInstanceAttrsE(t, ctx, projectID, location, instanceID)
	require.NoError(t, err)

	return result
}

// GetDataFusionInstanceAttrsE returns the settings Google Cloud holds for the given Cloud Data Fusion instance.
// The ctx parameter supports cancellation and timeouts.
func GetDataFusionInstanceAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, instanceID string) (*datafusion.Instance, error) {
	logger.Default.Logf(t, "Getting settings for Data Fusion instance %s in %s of project %s", instanceID, location, projectID)

	service, err := NewDataFusionServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetDataFusionInstanceAttrsWithClient(ctx, service, projectID, location, instanceID)
}

// GetDataFusionInstanceAttrsWithClient returns the settings Google Cloud holds for the given Cloud Data Fusion instance using the supplied
// *datafusion.Service. Prefer this variant in unit tests where the service is backed by an httptest fake server
// (see datafusion_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetDataFusionInstanceAttrsWithClient(ctx context.Context, service *datafusion.Service, projectID string, location string, instanceID string) (*datafusion.Instance, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/instances/%s", projectID, location, instanceID)

	result, err := service.Projects.Locations.Instances.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Data Fusion instance %s does not exist in %s of project %s", instanceID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Data Fusion instance %s in %s of project %s: %w", instanceID, location, projectID, err)
	}

	return result, nil
}

// NewDataFusionServiceE creates a Cloud Data Fusion service authenticated the same way every other client in
// this module is.
// The ctx parameter supports cancellation and timeouts.
func NewDataFusionServiceE(t testing.TestingT, ctx context.Context) (*datafusion.Service, error) {
	return datafusion.NewService(ctx, append(withOptions(), option.WithScopes(datafusion.CloudPlatformScope))...)
}
