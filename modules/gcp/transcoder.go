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
	"google.golang.org/api/transcoder/v1"
)

// GetTranscoderJobAttrs returns the settings Google Cloud holds for the Transcoder job, so a test can assert on what was
// actually created rather than only that it exists.
// A job carries the state Google put it in and the error it failed with, which is the only way a test
// can tell a job that was accepted from one that ran.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetTranscoderJobAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, jobID string) *transcoder.Job {
	result, err := GetTranscoderJobAttrsE(t, ctx, projectID, location, jobID)
	require.NoError(t, err)

	return result
}

// GetTranscoderJobAttrsE returns the settings Google Cloud holds for the Transcoder job.
// The ctx parameter supports cancellation and timeouts.
func GetTranscoderJobAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, jobID string) (*transcoder.Job, error) {
	logger.Default.Logf(t, "Getting settings for Transcoder job %s in %s in project %s", jobID, location, projectID)

	service, err := NewTranscoderServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetTranscoderJobAttrsWithClient(ctx, service, projectID, location, jobID)
}

// GetTranscoderJobAttrsWithClient returns the settings Google Cloud holds for the Transcoder job using the supplied
// *transcoder.Service. Prefer this variant in unit tests where the service is backed by an httptest
// fake server (see transcoder_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetTranscoderJobAttrsWithClient(ctx context.Context, service *transcoder.Service, projectID string, location string, jobID string) (*transcoder.Job, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/jobs/%s", projectID, location, jobID)

	result, err := service.Projects.Locations.Jobs.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Transcoder job %s in %s in project %s does not exist", jobID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Transcoder job %s in %s in project %s: %w", jobID, location, projectID, err)
	}

	return result, nil
}

// NewTranscoderServiceE creates a Transcoder service authenticated the same way every other client in this
// module is.
// The ctx parameter supports cancellation and timeouts.
func NewTranscoderServiceE(t testing.TestingT, ctx context.Context) (*transcoder.Service, error) {
	return transcoder.NewService(ctx, append(withOptions(), option.WithScopes(transcoder.CloudPlatformScope))...)
}
