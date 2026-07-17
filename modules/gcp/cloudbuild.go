package gcp

import (
	"context"
	"errors"
	"fmt"

	cloudbuild "cloud.google.com/go/cloudbuild/apiv1/v2"
	cloudbuildpb "cloud.google.com/go/cloudbuild/apiv1/v2/cloudbuildpb"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/iterator"
)

// CreateBuildContext creates a new build blocking until the operation is complete.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func CreateBuildContext(t testing.TestingT, ctx context.Context, projectID string, build *cloudbuildpb.Build) *cloudbuildpb.Build {
	out, err := CreateBuildContextE(t, ctx, projectID, build)
	require.NoError(t, err)

	return out
}

// CreateBuildContextE creates a new build blocking until the operation is complete.
// The ctx parameter supports cancellation and timeouts.
func CreateBuildContextE(t testing.TestingT, ctx context.Context, projectID string, build *cloudbuildpb.Build) (*cloudbuildpb.Build, error) {
	service, err := NewCloudBuildServiceContextE(t, ctx)
	if err != nil {
		return nil, err
	}

	defer func() { _ = service.Close() }()

	return CreateBuildWithClient(ctx, service, projectID, build)
}

// CreateBuildWithClient creates a new build blocking until the operation is complete using the
// supplied *cloudbuild.Client. Prefer this variant in unit tests where the client is backed by a
// mock gRPC server.
// The ctx parameter supports cancellation and timeouts.
func CreateBuildWithClient(ctx context.Context, service *cloudbuild.Client, projectID string, build *cloudbuildpb.Build) (*cloudbuildpb.Build, error) {
	req := &cloudbuildpb.CreateBuildRequest{
		ProjectId: projectID,
		Build:     build,
	}

	op, err := service.CreateBuild(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("CreateBuildContextE.CreateBuild(%s) got error: %w", projectID, err)
	}

	resp, err := op.Wait(ctx)
	if err != nil {
		return nil, fmt.Errorf("CreateBuildContextE.Wait(%s) got error: %w", projectID, err)
	}

	return resp, nil
}

// GetBuildContext gets the given build.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetBuildContext(t testing.TestingT, ctx context.Context, projectID string, buildID string) *cloudbuildpb.Build {
	out, err := GetBuildContextE(t, ctx, projectID, buildID)
	require.NoError(t, err)

	return out
}

// GetBuildContextE gets the given build.
// The ctx parameter supports cancellation and timeouts.
func GetBuildContextE(t testing.TestingT, ctx context.Context, projectID string, buildID string) (*cloudbuildpb.Build, error) {
	service, err := NewCloudBuildServiceContextE(t, ctx)
	if err != nil {
		return nil, err
	}

	defer func() { _ = service.Close() }()

	return GetBuildWithClient(ctx, service, projectID, buildID)
}

// GetBuildWithClient gets the given build using the supplied *cloudbuild.Client. Prefer this variant
// in unit tests where the client is backed by a mock gRPC server.
// The ctx parameter supports cancellation and timeouts.
func GetBuildWithClient(ctx context.Context, service *cloudbuild.Client, projectID string, buildID string) (*cloudbuildpb.Build, error) {
	req := &cloudbuildpb.GetBuildRequest{
		ProjectId: projectID,
		Id:        buildID,
	}

	resp, err := service.GetBuild(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("GetBuildContextE.GetBuild(%s, %s) got error: %w", projectID, buildID, err)
	}

	return resp, nil
}

// GetBuildsContext gets the list of builds for a given project.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetBuildsContext(t testing.TestingT, ctx context.Context, projectID string) []*cloudbuildpb.Build {
	out, err := GetBuildsContextE(t, ctx, projectID)
	require.NoError(t, err)

	return out
}

// GetBuildsContextE gets the list of builds for a given project.
// The ctx parameter supports cancellation and timeouts.
func GetBuildsContextE(t testing.TestingT, ctx context.Context, projectID string) ([]*cloudbuildpb.Build, error) {
	service, err := NewCloudBuildServiceContextE(t, ctx)
	if err != nil {
		return nil, err
	}

	defer func() { _ = service.Close() }()

	return GetBuildsWithClient(ctx, service, projectID)
}

// GetBuildsWithClient gets the list of builds for a given project using the supplied
// *cloudbuild.Client. Prefer this variant in unit tests where the client is backed by a mock gRPC
// server.
// The ctx parameter supports cancellation and timeouts.
func GetBuildsWithClient(ctx context.Context, service *cloudbuild.Client, projectID string) ([]*cloudbuildpb.Build, error) {
	req := &cloudbuildpb.ListBuildsRequest{
		ProjectId: projectID,
	}

	it := service.ListBuilds(ctx, req)
	builds := []*cloudbuildpb.Build{}

	for {
		resp, err := it.Next()

		if errors.Is(err, iterator.Done) {
			break
		}

		if err != nil {
			return nil, fmt.Errorf("GetBuildsContextE.ListBuilds(%s) got error: %w", projectID, err)
		}

		builds = append(builds, resp)
	}

	return builds, nil
}

// GetBuildsForTriggerContext gets a list of builds for a specific cloud build trigger.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetBuildsForTriggerContext(t testing.TestingT, ctx context.Context, projectID string, triggerID string) []*cloudbuildpb.Build {
	out, err := GetBuildsForTriggerContextE(t, ctx, projectID, triggerID)
	require.NoError(t, err)

	return out
}

// GetBuildsForTriggerContextE gets a list of builds for a specific cloud build trigger.
// The ctx parameter supports cancellation and timeouts.
func GetBuildsForTriggerContextE(t testing.TestingT, ctx context.Context, projectID string, triggerID string) ([]*cloudbuildpb.Build, error) {
	builds, err := GetBuildsContextE(t, ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("GetBuildsForTriggerContextE.ListBuilds(%s) got error: %w", projectID, err)
	}

	return filterBuildsByTrigger(builds, triggerID), nil
}

// GetBuildsForTriggerWithClient gets a list of builds for a specific cloud build trigger using the
// supplied *cloudbuild.Client. Prefer this variant in unit tests where the client is backed by a
// mock gRPC server.
// The ctx parameter supports cancellation and timeouts.
func GetBuildsForTriggerWithClient(ctx context.Context, service *cloudbuild.Client, projectID string, triggerID string) ([]*cloudbuildpb.Build, error) {
	builds, err := GetBuildsWithClient(ctx, service, projectID)
	if err != nil {
		return nil, fmt.Errorf("GetBuildsForTriggerContextE.ListBuilds(%s) got error: %w", projectID, err)
	}

	return filterBuildsByTrigger(builds, triggerID), nil
}

// filterBuildsByTrigger returns the subset of the given builds that were produced by the trigger
// with the given ID.
func filterBuildsByTrigger(builds []*cloudbuildpb.Build, triggerID string) []*cloudbuildpb.Build {
	filtered := []*cloudbuildpb.Build{}

	for _, build := range builds {
		if build.GetBuildTriggerId() == triggerID {
			filtered = append(filtered, build)
		}
	}

	return filtered
}

// NewCloudBuildServiceContext creates a new Cloud Build service, which is used to make Cloud Build API calls.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func NewCloudBuildServiceContext(t testing.TestingT, ctx context.Context) *cloudbuild.Client {
	service, err := NewCloudBuildServiceContextE(t, ctx)
	require.NoError(t, err)

	return service
}

// NewCloudBuildServiceContextE creates a new Cloud Build service, which is used to make Cloud Build API calls.
// The ctx parameter supports cancellation and timeouts.
func NewCloudBuildServiceContextE(t testing.TestingT, ctx context.Context) (*cloudbuild.Client, error) {
	service, err := cloudbuild.NewClient(ctx, withOptions()...)
	if err != nil {
		return nil, err
	}

	return service, nil
}
