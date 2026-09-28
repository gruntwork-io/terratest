package gcp

import (
	"context"
	"errors"
	"fmt"

	"github.com/gruntwork-io/terratest/modules/core/v2/logger"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
	"google.golang.org/api/parametermanager/v1"
)

// GetParameterAttrs returns the settings Google Cloud holds for the given parameter, so a test can assert on
// what was actually created rather than only that it exists. Global parameters live in the location called global, and a
// regional one in its region.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetParameterAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, parameterID string) *parametermanager.Parameter {
	parameter, err := GetParameterAttrsE(t, ctx, projectID, location, parameterID)
	require.NoError(t, err)

	return parameter
}

// GetParameterAttrsE returns the settings Google Cloud holds for the given parameter.
// The ctx parameter supports cancellation and timeouts.
func GetParameterAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, parameterID string) (*parametermanager.Parameter, error) {
	logger.Default.Logf(t, "Getting settings for parameter %s in %s in project %s", parameterID, location, projectID)

	service, err := NewParameterManagerServiceE(t, ctx, location)
	if err != nil {
		return nil, err
	}

	return GetParameterAttrsWithClient(ctx, service, projectID, location, parameterID)
}

// GetParameterAttrsWithClient returns the settings Google Cloud holds for the given parameter using the supplied
// *parametermanager.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see parametermanager_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetParameterAttrsWithClient(ctx context.Context, service *parametermanager.Service, projectID string, location string, parameterID string) (*parametermanager.Parameter, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/parameters/%s", projectID, location, parameterID)

	parameter, err := service.Projects.Locations.Parameters.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the parameter %s does not exist in %s in project %s", parameterID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for parameter %s in %s in project %s: %w", parameterID, location, projectID, err)
	}

	return parameter, nil
}

// NewParameterManagerServiceE creates a Parameter Manager service authenticated the same way every
// other client in this module is. A regional parameter is only served by its own region's endpoint,
// and the global one only by the global endpoint, so the location decides which is used.
// The ctx parameter supports cancellation and timeouts.
func NewParameterManagerServiceE(t testing.TestingT, ctx context.Context, location string) (*parametermanager.Service, error) {
	opts := append(withOptions(), option.WithScopes(parametermanager.CloudPlatformScope))
	if location != "global" {
		opts = append(opts, option.WithEndpoint(fmt.Sprintf("https://parametermanager.%s.rep.googleapis.com/", location)))
	}

	return parametermanager.NewService(ctx, opts...)
}
