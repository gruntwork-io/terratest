package gcp

import (
	"context"
	"fmt"

	"github.com/gruntwork-io/terratest/modules/core/v2/logger"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/compute/v1"
)

// FetchResourcePolicy queries GCP to return the settings it holds for the given resource policy, so a test can
// assert on what was actually created rather than only that it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func FetchResourcePolicy(t testing.TestingT, ctx context.Context, projectID string, region string, name string) *compute.ResourcePolicy {
	policy, err := FetchResourcePolicyE(t, ctx, projectID, region, name)
	require.NoError(t, err)

	return policy
}

// FetchResourcePolicyE queries GCP to return the settings it holds for the given resource policy.
// The ctx parameter supports cancellation and timeouts.
func FetchResourcePolicyE(t testing.TestingT, ctx context.Context, projectID string, region string, name string) (*compute.ResourcePolicy, error) {
	logger.Default.Logf(t, "Getting resource policy %s in region %s", name, region)

	service, err := NewComputeServiceContextE(t, ctx)
	if err != nil {
		return nil, err
	}

	return FetchResourcePolicyWithClient(ctx, service, projectID, region, name)
}

// FetchResourcePolicyWithClient queries GCP to return the settings it holds for the given resource policy using the
// supplied *compute.Service. Prefer this variant in unit tests where the service is backed by an
// httptest fake server (see resourcepolicy_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func FetchResourcePolicyWithClient(ctx context.Context, service *compute.Service, projectID string, region string, name string) (*compute.ResourcePolicy, error) {
	policy, err := service.ResourcePolicies.Get(projectID, region, name).Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("ResourcePolicies.Get(%s, %s, %s) got error: %w", projectID, region, name, err)
	}

	return policy, nil
}
