package gcp

import (
	"context"
	"errors"
	"fmt"

	"github.com/gruntwork-io/terratest/modules/core/v2/logger"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/cloudbuild/v1"
	"google.golang.org/api/googleapi"
)

// GetBuildTriggerAttrs returns the settings Google Cloud holds for the given Cloud Build trigger,
// so a test can assert on what was actually created rather than only that it exists. A trigger
// lives in a location, which has to be given.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetBuildTriggerAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, trigger string) *cloudbuild.BuildTrigger {
	buildTrigger, err := GetBuildTriggerAttrsE(t, ctx, projectID, location, trigger)
	require.NoError(t, err)

	return buildTrigger
}

// GetBuildTriggerAttrsE returns the settings Google Cloud holds for the given Cloud Build trigger.
// The ctx parameter supports cancellation and timeouts.
func GetBuildTriggerAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, trigger string) (*cloudbuild.BuildTrigger, error) {
	logger.Default.Logf(t, "Getting settings for build trigger %s in location %s in project %s", trigger, location, projectID)

	service, err := NewCloudBuildRESTServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetBuildTriggerAttrsWithClient(ctx, service, projectID, location, trigger)
}

// GetBuildTriggerAttrsWithClient returns the settings Google Cloud holds for the given Cloud Build
// trigger using the supplied *cloudbuild.Service. Prefer this variant in unit tests where the
// service is backed by an httptest fake server (see workerpool_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetBuildTriggerAttrsWithClient(ctx context.Context, service *cloudbuild.Service, projectID string, location string, trigger string) (*cloudbuild.BuildTrigger, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/triggers/%s", projectID, location, trigger)

	buildTrigger, err := service.Projects.Locations.Triggers.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("build trigger %s does not exist in location %s in project %s", trigger, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for build trigger %s in location %s in project %s: %w", trigger, location, projectID, err)
	}

	return buildTrigger, nil
}
