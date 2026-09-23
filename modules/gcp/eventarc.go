package gcp

import (
	"context"
	"errors"
	"fmt"

	"github.com/gruntwork-io/terratest/modules/core/v2/logger"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/eventarc/v1"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
)

// GetEventarcTriggerAttrs returns the settings Google Cloud holds for the given Eventarc trigger,
// so a test can assert on what was actually created rather than only that it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetEventarcTriggerAttrs(t testing.TestingT, ctx context.Context, projectID string, region string, triggerID string) *eventarc.Trigger {
	trigger, err := GetEventarcTriggerAttrsE(t, ctx, projectID, region, triggerID)
	require.NoError(t, err)

	return trigger
}

// GetEventarcTriggerAttrsE returns the settings Google Cloud holds for the given Eventarc trigger.
// The ctx parameter supports cancellation and timeouts.
func GetEventarcTriggerAttrsE(t testing.TestingT, ctx context.Context, projectID string, region string, triggerID string) (*eventarc.Trigger, error) {
	logger.Default.Logf(t, "Getting settings for Eventarc trigger %s in region %s in project %s", triggerID, region, projectID)

	service, err := NewEventarcServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetEventarcTriggerAttrsWithClient(ctx, service, projectID, region, triggerID)
}

// GetEventarcTriggerAttrsWithClient returns the settings Google Cloud holds for the given Eventarc
// trigger using the supplied *eventarc.Service. Prefer this variant in unit tests where the service
// is backed by an httptest fake server (see eventarc_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetEventarcTriggerAttrsWithClient(ctx context.Context, service *eventarc.Service, projectID string, region string, triggerID string) (*eventarc.Trigger, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/triggers/%s", projectID, region, triggerID)

	trigger, err := service.Projects.Locations.Triggers.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Eventarc trigger %s does not exist in region %s in project %s", triggerID, region, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Eventarc trigger %s in region %s in project %s: %w", triggerID, region, projectID, err)
	}

	return trigger, nil
}

// NewEventarcServiceE creates a Eventarc service authenticated the same way every other client in
// this module is.
// The ctx parameter supports cancellation and timeouts.
func NewEventarcServiceE(t testing.TestingT, ctx context.Context) (*eventarc.Service, error) {
	return eventarc.NewService(ctx, append(withOptions(), option.WithScopes(eventarc.CloudPlatformScope))...)
}
