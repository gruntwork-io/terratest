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
	"google.golang.org/api/securesourcemanager/v1"
)

// GetSecureSourceManagerInstanceAttrs returns the settings Google Cloud holds for the given Secure Source Manager instance, so a test can assert on what was
// actually created rather than only that it exists.
// The instance is the git host itself, so its state and the key it encrypts repositories with are what it is.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetSecureSourceManagerInstanceAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, instanceID string) *securesourcemanager.Instance {
	attrs, err := GetSecureSourceManagerInstanceAttrsE(t, ctx, projectID, location, instanceID)
	require.NoError(t, err)

	return attrs
}

// GetSecureSourceManagerInstanceAttrsE returns the settings Google Cloud holds for the given Secure Source Manager instance.
// The ctx parameter supports cancellation and timeouts.
func GetSecureSourceManagerInstanceAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, instanceID string) (*securesourcemanager.Instance, error) {
	logger.Default.Logf(t, "Getting settings for Secure Source Manager instance %s in %s in project %s", instanceID, location, projectID)

	service, err := NewSecureSourceManagerServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetSecureSourceManagerInstanceAttrsWithClient(ctx, service, projectID, location, instanceID)
}

// GetSecureSourceManagerInstanceAttrsWithClient returns the settings Google Cloud holds for the given Secure Source Manager instance using the supplied
// *securesourcemanager.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see securesourcemanager_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetSecureSourceManagerInstanceAttrsWithClient(ctx context.Context, service *securesourcemanager.Service, projectID string, location string, instanceID string) (*securesourcemanager.Instance, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/instances/%s", projectID, location, instanceID)

	attrs, err := service.Projects.Locations.Instances.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Secure Source Manager instance %s in %s in project %s does not exist", instanceID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Secure Source Manager instance %s in %s in project %s: %w", instanceID, location, projectID, err)
	}

	return attrs, nil
}

// GetSecureSourceManagerRepositoryAttrs returns the settings Google Cloud holds for the given Secure Source Manager repository, so a test can assert on what was
// actually created rather than only that it exists.
// A repository is where code lives, so the instance that hosts it and its default branch are what a clone depends on.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetSecureSourceManagerRepositoryAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, repositoryID string) *securesourcemanager.Repository {
	attrs, err := GetSecureSourceManagerRepositoryAttrsE(t, ctx, projectID, location, repositoryID)
	require.NoError(t, err)

	return attrs
}

// GetSecureSourceManagerRepositoryAttrsE returns the settings Google Cloud holds for the given Secure Source Manager repository.
// The ctx parameter supports cancellation and timeouts.
func GetSecureSourceManagerRepositoryAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, repositoryID string) (*securesourcemanager.Repository, error) {
	logger.Default.Logf(t, "Getting settings for Secure Source Manager repository %s in %s in project %s", repositoryID, location, projectID)

	service, err := NewSecureSourceManagerServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetSecureSourceManagerRepositoryAttrsWithClient(ctx, service, projectID, location, repositoryID)
}

// GetSecureSourceManagerRepositoryAttrsWithClient returns the settings Google Cloud holds for the given Secure Source Manager repository using the supplied
// *securesourcemanager.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see securesourcemanager_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetSecureSourceManagerRepositoryAttrsWithClient(ctx context.Context, service *securesourcemanager.Service, projectID string, location string, repositoryID string) (*securesourcemanager.Repository, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/repositories/%s", projectID, location, repositoryID)

	attrs, err := service.Projects.Locations.Repositories.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Secure Source Manager repository %s in %s in project %s does not exist", repositoryID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Secure Source Manager repository %s in %s in project %s: %w", repositoryID, location, projectID, err)
	}

	return attrs, nil
}

// GetSecureSourceManagerBranchRuleAttrs returns the settings Google Cloud holds for the given Secure Source Manager branch rule, so a test can assert on what was
// actually created rather than only that it exists.
// A branch rule is what stops a push, so the pattern it matches and whether it demands review decide what gets in.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetSecureSourceManagerBranchRuleAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, repositoryID string, ruleID string) *securesourcemanager.BranchRule {
	attrs, err := GetSecureSourceManagerBranchRuleAttrsE(t, ctx, projectID, location, repositoryID, ruleID)
	require.NoError(t, err)

	return attrs
}

// GetSecureSourceManagerBranchRuleAttrsE returns the settings Google Cloud holds for the given Secure Source Manager branch rule.
// The ctx parameter supports cancellation and timeouts.
func GetSecureSourceManagerBranchRuleAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, repositoryID string, ruleID string) (*securesourcemanager.BranchRule, error) {
	logger.Default.Logf(t, "Getting settings for Secure Source Manager branch rule %s in repository %s in %s in project %s", ruleID, repositoryID, location, projectID)

	service, err := NewSecureSourceManagerServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetSecureSourceManagerBranchRuleAttrsWithClient(ctx, service, projectID, location, repositoryID, ruleID)
}

// GetSecureSourceManagerBranchRuleAttrsWithClient returns the settings Google Cloud holds for the given Secure Source Manager branch rule using the supplied
// *securesourcemanager.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see securesourcemanager_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetSecureSourceManagerBranchRuleAttrsWithClient(ctx context.Context, service *securesourcemanager.Service, projectID string, location string, repositoryID string, ruleID string) (*securesourcemanager.BranchRule, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/repositories/%s/branchRules/%s", projectID, location, repositoryID, ruleID)

	attrs, err := service.Projects.Locations.Repositories.BranchRules.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Secure Source Manager branch rule %s in repository %s in %s in project %s does not exist", ruleID, repositoryID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Secure Source Manager branch rule %s in repository %s in %s in project %s: %w", ruleID, repositoryID, location, projectID, err)
	}

	return attrs, nil
}

// GetSecureSourceManagerHookAttrs returns the settings Google Cloud holds for the given Secure Source Manager hook, so a test can assert on what was
// actually created rather than only that it exists.
// A hook is what a push notifies, so the URL it posts to and which events trigger it are the point.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetSecureSourceManagerHookAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, repositoryID string, hookID string) *securesourcemanager.Hook {
	attrs, err := GetSecureSourceManagerHookAttrsE(t, ctx, projectID, location, repositoryID, hookID)
	require.NoError(t, err)

	return attrs
}

// GetSecureSourceManagerHookAttrsE returns the settings Google Cloud holds for the given Secure Source Manager hook.
// The ctx parameter supports cancellation and timeouts.
func GetSecureSourceManagerHookAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, repositoryID string, hookID string) (*securesourcemanager.Hook, error) {
	logger.Default.Logf(t, "Getting settings for Secure Source Manager hook %s in repository %s in %s in project %s", hookID, repositoryID, location, projectID)

	service, err := NewSecureSourceManagerServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetSecureSourceManagerHookAttrsWithClient(ctx, service, projectID, location, repositoryID, hookID)
}

// GetSecureSourceManagerHookAttrsWithClient returns the settings Google Cloud holds for the given Secure Source Manager hook using the supplied
// *securesourcemanager.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see securesourcemanager_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetSecureSourceManagerHookAttrsWithClient(ctx context.Context, service *securesourcemanager.Service, projectID string, location string, repositoryID string, hookID string) (*securesourcemanager.Hook, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/repositories/%s/hooks/%s", projectID, location, repositoryID, hookID)

	attrs, err := service.Projects.Locations.Repositories.Hooks.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Secure Source Manager hook %s in repository %s in %s in project %s does not exist", hookID, repositoryID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Secure Source Manager hook %s in repository %s in %s in project %s: %w", hookID, repositoryID, location, projectID, err)
	}

	return attrs, nil
}

// GetSecureSourceManagerInstanceIamPolicyAttrs returns the IAM policy Google Cloud holds for the given Secure Source Manager instance, so a test can assert on what was
// actually created rather than only that it exists.
// Who may read or push is what the policy decides, and a repository Google created on its own carries no binding at all.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetSecureSourceManagerInstanceIamPolicyAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, id string) *securesourcemanager.Policy {
	policy, err := GetSecureSourceManagerInstanceIamPolicyAttrsE(t, ctx, projectID, location, id)
	require.NoError(t, err)

	return policy
}

// GetSecureSourceManagerInstanceIamPolicyAttrsE returns the IAM policy Google Cloud holds for the given Secure Source Manager instance.
// The ctx parameter supports cancellation and timeouts.
func GetSecureSourceManagerInstanceIamPolicyAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, id string) (*securesourcemanager.Policy, error) {
	logger.Default.Logf(t, "Getting the IAM policy for Secure Source Manager instance %s in %s in project %s", id, location, projectID)

	service, err := NewSecureSourceManagerServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetSecureSourceManagerInstanceIamPolicyAttrsWithClient(ctx, service, projectID, location, id)
}

// GetSecureSourceManagerInstanceIamPolicyAttrsWithClient returns the IAM policy Google Cloud holds for the given Secure Source Manager instance using the supplied
// *securesourcemanager.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see securesourcemanager_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetSecureSourceManagerInstanceIamPolicyAttrsWithClient(ctx context.Context, service *securesourcemanager.Service, projectID string, location string, id string) (*securesourcemanager.Policy, error) {
	resource := fmt.Sprintf("projects/%s/locations/%s/instances/%s", projectID, location, id)

	policy, err := service.Projects.Locations.Instances.GetIamPolicy(resource).OptionsRequestedPolicyVersion(iamPolicyVersionWithConditions).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Secure Source Manager instance %s in %s in project %s does not exist", id, location, projectID)
		}

		return nil, fmt.Errorf("failed to get the IAM policy for Secure Source Manager instance %s in %s in project %s: %w", id, location, projectID, err)
	}

	return policy, nil
}

// GetSecureSourceManagerRepositoryIamPolicyAttrs returns the IAM policy Google Cloud holds for the given Secure Source Manager repository, so a test can assert on what was
// actually created rather than only that it exists.
// Who may read or push is what the policy decides, and a repository Google created on its own carries no binding at all.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetSecureSourceManagerRepositoryIamPolicyAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, id string) *securesourcemanager.Policy {
	policy, err := GetSecureSourceManagerRepositoryIamPolicyAttrsE(t, ctx, projectID, location, id)
	require.NoError(t, err)

	return policy
}

// GetSecureSourceManagerRepositoryIamPolicyAttrsE returns the IAM policy Google Cloud holds for the given Secure Source Manager repository.
// The ctx parameter supports cancellation and timeouts.
func GetSecureSourceManagerRepositoryIamPolicyAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, id string) (*securesourcemanager.Policy, error) {
	logger.Default.Logf(t, "Getting the IAM policy for Secure Source Manager repository %s in %s in project %s", id, location, projectID)

	service, err := NewSecureSourceManagerServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetSecureSourceManagerRepositoryIamPolicyAttrsWithClient(ctx, service, projectID, location, id)
}

// GetSecureSourceManagerRepositoryIamPolicyAttrsWithClient returns the IAM policy Google Cloud holds for the given Secure Source Manager repository using the supplied
// *securesourcemanager.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see securesourcemanager_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetSecureSourceManagerRepositoryIamPolicyAttrsWithClient(ctx context.Context, service *securesourcemanager.Service, projectID string, location string, id string) (*securesourcemanager.Policy, error) {
	resource := fmt.Sprintf("projects/%s/locations/%s/repositories/%s", projectID, location, id)

	policy, err := service.Projects.Locations.Repositories.GetIamPolicy(resource).OptionsRequestedPolicyVersion(iamPolicyVersionWithConditions).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Secure Source Manager repository %s in %s in project %s does not exist", id, location, projectID)
		}

		return nil, fmt.Errorf("failed to get the IAM policy for Secure Source Manager repository %s in %s in project %s: %w", id, location, projectID, err)
	}

	return policy, nil
}

// NewSecureSourceManagerServiceE creates a Secure Source Manager service authenticated the same way every other client in this
// module is.
// The ctx parameter supports cancellation and timeouts.
func NewSecureSourceManagerServiceE(t testing.TestingT, ctx context.Context) (*securesourcemanager.Service, error) {
	return securesourcemanager.NewService(ctx, append(withOptions(), option.WithScopes(securesourcemanager.CloudPlatformScope))...)
}
