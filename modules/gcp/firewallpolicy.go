package gcp

import (
	"context"
	"fmt"

	"github.com/gruntwork-io/terratest/modules/core/v2/logger"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/compute/v1"
)

// FetchNetworkFirewallPolicy queries GCP to return the settings it holds for the given global network firewall policy, so a test can
// assert on what was actually created rather than only that it exists. A network firewall policy belongs to a project, unlike the hierarchical policy that belongs to an organization or a folder.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func FetchNetworkFirewallPolicy(t testing.TestingT, ctx context.Context, projectID string, name string) *compute.FirewallPolicy {
	policy, err := FetchNetworkFirewallPolicyE(t, ctx, projectID, name)
	require.NoError(t, err)

	return policy
}

// FetchNetworkFirewallPolicyE queries GCP to return the settings it holds for the given global network firewall policy.
// The ctx parameter supports cancellation and timeouts.
func FetchNetworkFirewallPolicyE(t testing.TestingT, ctx context.Context, projectID string, name string) (*compute.FirewallPolicy, error) {
	logger.Default.Logf(t, "Getting global network firewall policy %s", name)

	service, err := NewComputeServiceContextE(t, ctx)
	if err != nil {
		return nil, err
	}

	return FetchNetworkFirewallPolicyWithClient(ctx, service, projectID, name)
}

// FetchNetworkFirewallPolicyWithClient queries GCP to return the settings it holds for the given global network firewall policy using the
// supplied *compute.Service. Prefer this variant in unit tests where the service is backed by an
// httptest fake server (see firewallpolicy_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func FetchNetworkFirewallPolicyWithClient(ctx context.Context, service *compute.Service, projectID string, name string) (*compute.FirewallPolicy, error) {
	policy, err := service.NetworkFirewallPolicies.Get(projectID, name).Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("NetworkFirewallPolicies.Get(%s, %s) got error: %w", projectID, name, err)
	}

	return policy, nil
}

// FetchRegionNetworkFirewallPolicy queries GCP to return the settings it holds for the given regional network firewall policy, so a test can
// assert on what was actually created rather than only that it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func FetchRegionNetworkFirewallPolicy(t testing.TestingT, ctx context.Context, projectID string, region string, name string) *compute.FirewallPolicy {
	policy, err := FetchRegionNetworkFirewallPolicyE(t, ctx, projectID, region, name)
	require.NoError(t, err)

	return policy
}

// FetchRegionNetworkFirewallPolicyE queries GCP to return the settings it holds for the given regional network firewall policy.
// The ctx parameter supports cancellation and timeouts.
func FetchRegionNetworkFirewallPolicyE(t testing.TestingT, ctx context.Context, projectID string, region string, name string) (*compute.FirewallPolicy, error) {
	logger.Default.Logf(t, "Getting regional network firewall policy %s in region %s", name, region)

	service, err := NewComputeServiceContextE(t, ctx)
	if err != nil {
		return nil, err
	}

	return FetchRegionNetworkFirewallPolicyWithClient(ctx, service, projectID, region, name)
}

// FetchRegionNetworkFirewallPolicyWithClient queries GCP to return the settings it holds for the given regional network firewall policy using the
// supplied *compute.Service. Prefer this variant in unit tests where the service is backed by an
// httptest fake server (see firewallpolicy_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func FetchRegionNetworkFirewallPolicyWithClient(ctx context.Context, service *compute.Service, projectID string, region string, name string) (*compute.FirewallPolicy, error) {
	policy, err := service.RegionNetworkFirewallPolicies.Get(projectID, region, name).Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("RegionNetworkFirewallPolicies.Get(%s, %s, %s) got error: %w", projectID, region, name, err)
	}

	return policy, nil
}

// iamPolicyVersionWithConditions is the policy version that carries conditional bindings. Google
// returns a policy at version 1 unless asked otherwise, and a version 1 answer has no room for a
// condition.
const iamPolicyVersionWithConditions = 3

// FetchNetworkFirewallPolicyIamPolicy queries GCP to return the settings it holds for the given IAM policy of a global network firewall policy, so a test can
// assert on what was actually created rather than only that it exists. This reads who may act on the policy, which is a different question from what the policy allows through.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func FetchNetworkFirewallPolicyIamPolicy(t testing.TestingT, ctx context.Context, projectID string, name string) *compute.Policy {
	policy, err := FetchNetworkFirewallPolicyIamPolicyE(t, ctx, projectID, name)
	require.NoError(t, err)

	return policy
}

// FetchNetworkFirewallPolicyIamPolicyE queries GCP to return the settings it holds for the given IAM policy of a global network firewall policy.
// The ctx parameter supports cancellation and timeouts.
func FetchNetworkFirewallPolicyIamPolicyE(t testing.TestingT, ctx context.Context, projectID string, name string) (*compute.Policy, error) {
	logger.Default.Logf(t, "Getting IAM policy of a global network firewall policy %s", name)

	service, err := NewComputeServiceContextE(t, ctx)
	if err != nil {
		return nil, err
	}

	return FetchNetworkFirewallPolicyIamPolicyWithClient(ctx, service, projectID, name)
}

// FetchNetworkFirewallPolicyIamPolicyWithClient queries GCP to return the settings it holds for the given IAM policy of a global network firewall policy using the
// supplied *compute.Service. Prefer this variant in unit tests where the service is backed by an
// httptest fake server (see firewallpolicy_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func FetchNetworkFirewallPolicyIamPolicyWithClient(ctx context.Context, service *compute.Service, projectID string, name string) (*compute.Policy, error) {
	// A policy carrying a conditional binding is only returned in full at version 3, so that is what
	// is asked for: at a lower version Google drops the condition or refuses the call outright.
	policy, err := service.NetworkFirewallPolicies.GetIamPolicy(projectID, name).
		OptionsRequestedPolicyVersion(iamPolicyVersionWithConditions).Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("NetworkFirewallPolicies.GetIamPolicy(%s, %s) got error: %w", projectID, name, err)
	}

	return policy, nil
}

// FetchRegionNetworkFirewallPolicyIamPolicy queries GCP to return the settings it holds for the given IAM policy of a regional network firewall policy, so a test can
// assert on what was actually created rather than only that it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func FetchRegionNetworkFirewallPolicyIamPolicy(t testing.TestingT, ctx context.Context, projectID string, region string, name string) *compute.Policy {
	policy, err := FetchRegionNetworkFirewallPolicyIamPolicyE(t, ctx, projectID, region, name)
	require.NoError(t, err)

	return policy
}

// FetchRegionNetworkFirewallPolicyIamPolicyE queries GCP to return the settings it holds for the given IAM policy of a regional network firewall policy.
// The ctx parameter supports cancellation and timeouts.
func FetchRegionNetworkFirewallPolicyIamPolicyE(t testing.TestingT, ctx context.Context, projectID string, region string, name string) (*compute.Policy, error) {
	logger.Default.Logf(t, "Getting IAM policy of a regional network firewall policy %s in region %s", name, region)

	service, err := NewComputeServiceContextE(t, ctx)
	if err != nil {
		return nil, err
	}

	return FetchRegionNetworkFirewallPolicyIamPolicyWithClient(ctx, service, projectID, region, name)
}

// FetchRegionNetworkFirewallPolicyIamPolicyWithClient queries GCP to return the settings it holds for the given IAM policy of a regional network firewall policy using the
// supplied *compute.Service. Prefer this variant in unit tests where the service is backed by an
// httptest fake server (see firewallpolicy_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func FetchRegionNetworkFirewallPolicyIamPolicyWithClient(ctx context.Context, service *compute.Service, projectID string, region string, name string) (*compute.Policy, error) {
	// A policy carrying a conditional binding is only returned in full at version 3, so that is what
	// is asked for: at a lower version Google drops the condition or refuses the call outright.
	policy, err := service.RegionNetworkFirewallPolicies.GetIamPolicy(projectID, region, name).
		OptionsRequestedPolicyVersion(iamPolicyVersionWithConditions).Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("RegionNetworkFirewallPolicies.GetIamPolicy(%s, %s, %s) got error: %w", projectID, region, name, err)
	}

	return policy, nil
}
