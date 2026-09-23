package gcp

import (
	"context"
	"fmt"

	"github.com/gruntwork-io/terratest/modules/core/v2/logger"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/compute/v1"
)

// FetchSecurityPolicy queries GCP to return the settings it holds for the given Cloud Armor
// security policy, so a test can assert on what was actually created rather than only that it
// exists. The rules come back with the policy, including the default rule Google adds.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func FetchSecurityPolicy(t testing.TestingT, ctx context.Context, projectID string, name string) *compute.SecurityPolicy {
	policy, err := FetchSecurityPolicyE(t, ctx, projectID, name)
	require.NoError(t, err)

	return policy
}

// FetchSecurityPolicyE queries GCP to return the settings it holds for the given Cloud Armor
// security policy.
// The ctx parameter supports cancellation and timeouts.
func FetchSecurityPolicyE(t testing.TestingT, ctx context.Context, projectID string, name string) (*compute.SecurityPolicy, error) {
	logger.Default.Logf(t, "Getting security policy %s", name)

	service, err := NewComputeServiceContextE(t, ctx)
	if err != nil {
		return nil, err
	}

	return FetchSecurityPolicyWithClient(ctx, service, projectID, name)
}

// FetchSecurityPolicyWithClient queries GCP to return the settings it holds for the given Cloud
// Armor security policy using the supplied *compute.Service. Prefer this variant in unit tests
// where the service is backed by an httptest fake server (see cloudarmor_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func FetchSecurityPolicyWithClient(ctx context.Context, service *compute.Service, projectID string, name string) (*compute.SecurityPolicy, error) {
	policy, err := service.SecurityPolicies.Get(projectID, name).Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("SecurityPolicies.Get(%s, %s) got error: %w", projectID, name, err)
	}

	return policy, nil
}
