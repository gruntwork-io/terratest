package gcp

import (
	"context"
	"errors"
	"fmt"

	"github.com/gruntwork-io/terratest/modules/core/v2/logger"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/googleapi"
	iamv2 "google.golang.org/api/iam/v2"
	"google.golang.org/api/option"
)

// GetIAMDenyPolicyAttrs returns the settings Google Cloud holds for the IAM deny policy, so a test can assert on what was
// actually created rather than only that it exists.
// A deny policy hangs off an attachment point rather than a parent resource, and the forward slashes
// in that attachment point have to arrive percent encoded, which is what this helper builds.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetIAMDenyPolicyAttrs(t testing.TestingT, ctx context.Context, projectID string, policyID string) *iamv2.GoogleIamV2Policy {
	result, err := GetIAMDenyPolicyAttrsE(t, ctx, projectID, policyID)
	require.NoError(t, err)

	return result
}

// GetIAMDenyPolicyAttrsE returns the settings Google Cloud holds for the IAM deny policy.
// The ctx parameter supports cancellation and timeouts.
func GetIAMDenyPolicyAttrsE(t testing.TestingT, ctx context.Context, projectID string, policyID string) (*iamv2.GoogleIamV2Policy, error) {
	logger.Default.Logf(t, "Getting settings for IAM deny policy %s on project %s", policyID, projectID)

	service, err := NewIAMV2ServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetIAMDenyPolicyAttrsWithClient(ctx, service, projectID, policyID)
}

// GetIAMDenyPolicyAttrsWithClient returns the settings Google Cloud holds for the IAM deny policy using the supplied
// *iamv2.Service. Prefer this variant in unit tests where the service is backed by an httptest
// fake server (see iamdenypolicy_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetIAMDenyPolicyAttrsWithClient(ctx context.Context, service *iamv2.Service, projectID string, policyID string) (*iamv2.GoogleIamV2Policy, error) {
	name := fmt.Sprintf("policies/cloudresourcemanager.googleapis.com%%2Fprojects%%2F%s/denypolicies/%s", projectID, policyID)

	result, err := service.Policies.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the IAM deny policy %s on project %s does not exist", policyID, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for IAM deny policy %s on project %s: %w", policyID, projectID, err)
	}

	return result, nil
}

// NewIAMV2ServiceE creates an IAM v2 service authenticated the same way every other client in this
// module is. Deny policies live only on the v2 surface, which is a separate client from the v1 one
// the service account and role reads in iam.go use.
// The ctx parameter supports cancellation and timeouts.
func NewIAMV2ServiceE(t testing.TestingT, ctx context.Context) (*iamv2.Service, error) {
	return iamv2.NewService(ctx, append(withOptions(), option.WithScopes(iamv2.CloudPlatformScope))...)
}
