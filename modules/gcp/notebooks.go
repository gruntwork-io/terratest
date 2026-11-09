package gcp

import (
	"context"
	"errors"
	"fmt"

	"github.com/gruntwork-io/terratest/modules/core/v2/logger"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/notebooks/v1"
	"google.golang.org/api/option"
)

// GetNotebooksEnvironmentAttrs returns the settings Google Cloud holds for the given Notebooks
// environment, so a test can assert on the image an instance would boot from rather than only that
// the environment exists. An environment describes a machine; it does not run one.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetNotebooksEnvironmentAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, environmentID string) *notebooks.Environment {
	environment, err := GetNotebooksEnvironmentAttrsE(t, ctx, projectID, location, environmentID)
	require.NoError(t, err)

	return environment
}

// GetNotebooksEnvironmentAttrsE returns the settings Google Cloud holds for the given Notebooks
// environment.
// The ctx parameter supports cancellation and timeouts.
func GetNotebooksEnvironmentAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, environmentID string) (*notebooks.Environment, error) {
	logger.Default.Logf(t, "Getting settings for Notebooks environment %s in %s in project %s", environmentID, location, projectID)

	service, err := NewNotebooksServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetNotebooksEnvironmentAttrsWithClient(ctx, service, projectID, location, environmentID)
}

// GetNotebooksEnvironmentAttrsWithClient returns the settings Google Cloud holds for the given
// Notebooks environment using the supplied *notebooks.Service. Prefer this variant in unit tests
// where the service is backed by an httptest fake server (see notebooks_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetNotebooksEnvironmentAttrsWithClient(ctx context.Context, service *notebooks.Service, projectID string, location string, environmentID string) (*notebooks.Environment, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/environments/%s", projectID, location, environmentID)

	environment, err := service.Projects.Locations.Environments.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Notebooks environment %s in %s in project %s does not exist", environmentID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Notebooks environment %s in %s in project %s: %w", environmentID, location, projectID, err)
	}

	return environment, nil
}

// GetNotebooksRuntimeAttrs returns the settings Google Cloud holds for the given managed notebooks runtime, so a test can assert on
// what was actually created rather than only that it exists.
// The proxy URI a user reaches the runtime on is assigned by Google, so it can only be asserted on
// from a read.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetNotebooksRuntimeAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, runtimeID string) *notebooks.Runtime {
	result, err := GetNotebooksRuntimeAttrsE(t, ctx, projectID, location, runtimeID)
	require.NoError(t, err)

	return result
}

// GetNotebooksRuntimeAttrsE returns the settings Google Cloud holds for the given managed notebooks runtime.
// The ctx parameter supports cancellation and timeouts.
func GetNotebooksRuntimeAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, runtimeID string) (*notebooks.Runtime, error) {
	logger.Default.Logf(t, "Getting settings for managed notebooks runtime %s in %s of project %s", runtimeID, location, projectID)

	service, err := NewNotebooksServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetNotebooksRuntimeAttrsWithClient(ctx, service, projectID, location, runtimeID)
}

// GetNotebooksRuntimeAttrsWithClient returns the settings Google Cloud holds for the given managed notebooks runtime using the supplied
// *notebooks.Service. Prefer this variant in unit tests where the service is backed by an httptest fake server
// (see notebooks_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetNotebooksRuntimeAttrsWithClient(ctx context.Context, service *notebooks.Service, projectID string, location string, runtimeID string) (*notebooks.Runtime, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/runtimes/%s", projectID, location, runtimeID)

	result, err := service.Projects.Locations.Runtimes.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the managed notebooks runtime %s does not exist in %s of project %s", runtimeID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for managed notebooks runtime %s in %s of project %s: %w", runtimeID, location, projectID, err)
	}

	return result, nil
}

// GetNotebooksRuntimeIamPolicyAttrs returns the IAM policy Google Cloud holds for the given managed
// notebooks runtime, so a test can assert on the bindings that were applied rather than only that
// the apply finished.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetNotebooksRuntimeIamPolicyAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, runtimeID string) *notebooks.Policy {
	policy, err := GetNotebooksRuntimeIamPolicyAttrsE(t, ctx, projectID, location, runtimeID)
	require.NoError(t, err)

	return policy
}

// GetNotebooksRuntimeIamPolicyAttrsE returns the IAM policy Google Cloud holds for the given managed
// notebooks runtime.
// The ctx parameter supports cancellation and timeouts.
func GetNotebooksRuntimeIamPolicyAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, runtimeID string) (*notebooks.Policy, error) {
	logger.Default.Logf(t, "Getting the IAM policy for managed notebooks runtime %s in %s of project %s", runtimeID, location, projectID)

	service, err := NewNotebooksServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetNotebooksRuntimeIamPolicyAttrsWithClient(ctx, service, projectID, location, runtimeID)
}

// GetNotebooksRuntimeIamPolicyAttrsWithClient returns the IAM policy Google Cloud holds for the given
// managed notebooks runtime using the supplied *notebooks.Service. Prefer this variant in unit tests
// where the service is backed by an httptest fake server (see notebooks_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetNotebooksRuntimeIamPolicyAttrsWithClient(ctx context.Context, service *notebooks.Service, projectID string, location string, runtimeID string) (*notebooks.Policy, error) {
	resource := fmt.Sprintf("projects/%s/locations/%s/runtimes/%s", projectID, location, runtimeID)

	policy, err := service.Projects.Locations.Runtimes.GetIamPolicy(resource).OptionsRequestedPolicyVersion(iamPolicyVersionWithConditions).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the managed notebooks runtime %s does not exist in %s of project %s", runtimeID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get the IAM policy for managed notebooks runtime %s in %s of project %s: %w", runtimeID, location, projectID, err)
	}

	return policy, nil
}

// NewNotebooksServiceE creates a Notebooks service authenticated the same way every other client in
// this module is.
// The ctx parameter supports cancellation and timeouts.
func NewNotebooksServiceE(t testing.TestingT, ctx context.Context) (*notebooks.Service, error) {
	return notebooks.NewService(ctx, append(withOptions(), option.WithScopes(notebooks.CloudPlatformScope))...)
}
