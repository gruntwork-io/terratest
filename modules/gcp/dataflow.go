package gcp

import (
	"context"
	"errors"
	"fmt"

	"github.com/gruntwork-io/terratest/modules/core/v2/logger"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
	dataflow "google.golang.org/api/dataflow/v1b3"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
)

// GetDataflowJobAttrs returns the settings Google Cloud holds for the given Dataflow job, so a test
// can assert on what was actually created rather than only that it exists. Google assigns a job's
// id when it is launched, so the caller passes the id it got back rather than the name it asked for.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetDataflowJobAttrs(t testing.TestingT, ctx context.Context, projectID string, region string, jobID string) *dataflow.Job {
	job, err := GetDataflowJobAttrsE(t, ctx, projectID, region, jobID)
	require.NoError(t, err)

	return job
}

// GetDataflowJobAttrsE returns the settings Google Cloud holds for the given Dataflow job.
// The ctx parameter supports cancellation and timeouts.
func GetDataflowJobAttrsE(t testing.TestingT, ctx context.Context, projectID string, region string, jobID string) (*dataflow.Job, error) {
	logger.Default.Logf(t, "Getting settings for Dataflow job %s in region %s in project %s", jobID, region, projectID)

	service, err := NewDataflowServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetDataflowJobAttrsWithClient(ctx, service, projectID, region, jobID)
}

// GetDataflowJobAttrsWithClient returns the settings Google Cloud holds for the given Dataflow job
// using the supplied *dataflow.Service. Prefer this variant in unit tests where the service is
// backed by an httptest fake server (see dataflow_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetDataflowJobAttrsWithClient(ctx context.Context, service *dataflow.Service, projectID string, region string, jobID string) (*dataflow.Job, error) {
	// Without a view Google answers with the summary only, which carries no labels and no
	// environment, so a caller asserting on either would read an empty map or a nil pointer.
	job, err := service.Projects.Locations.Jobs.Get(projectID, region, jobID).View("JOB_VIEW_ALL").Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Dataflow job %s does not exist in region %s in project %s", jobID, region, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Dataflow job %s in region %s in project %s: %w", jobID, region, projectID, err)
	}

	return job, nil
}

// NewDataflowServiceE creates a Dataflow service authenticated the same way every other client in
// this module is.
// The ctx parameter supports cancellation and timeouts.
func NewDataflowServiceE(t testing.TestingT, ctx context.Context) (*dataflow.Service, error) {
	return dataflow.NewService(ctx, append(withOptions(), option.WithScopes(dataflow.CloudPlatformScope))...)
}
