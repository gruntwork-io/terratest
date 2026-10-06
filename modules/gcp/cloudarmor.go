package gcp

import (
	"context"
	"errors"
	"fmt"

	"github.com/gruntwork-io/terratest/modules/core/v2/logger"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/googleapi"
)

// FetchSecurityPolicyContext queries GCP to return the settings it holds for the given Cloud Armor
// security policy, so a test can assert on what was actually created rather than only that it
// exists. The rules come back with the policy, including the default rule Google adds.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func FetchSecurityPolicyContext(t testing.TestingT, ctx context.Context, projectID string, name string) *compute.SecurityPolicy {
	policy, err := FetchSecurityPolicyContextE(t, ctx, projectID, name)
	require.NoError(t, err)

	return policy
}

// FetchSecurityPolicyContextE queries GCP to return the settings it holds for the given Cloud Armor
// security policy.
// The ctx parameter supports cancellation and timeouts.
func FetchSecurityPolicyContextE(t testing.TestingT, ctx context.Context, projectID string, name string) (*compute.SecurityPolicy, error) {
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

// FetchRegionSecurityPolicy returns the settings Google Cloud holds for the given regional security
// policy, so a test can assert on what was actually created rather than only that it exists. A
// regional policy protects a regional load balancer; the global one protects a global load balancer,
// and they are separate resources with separate reads.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func FetchRegionSecurityPolicy(t testing.TestingT, ctx context.Context, projectID string, region string, name string) *compute.SecurityPolicy {
	policy, err := FetchRegionSecurityPolicyE(t, ctx, projectID, region, name)
	require.NoError(t, err)

	return policy
}

// FetchRegionSecurityPolicyE returns the settings Google Cloud holds for the given regional security
// policy.
// The ctx parameter supports cancellation and timeouts.
func FetchRegionSecurityPolicyE(t testing.TestingT, ctx context.Context, projectID string, region string, name string) (*compute.SecurityPolicy, error) {
	logger.Default.Logf(t, "Getting settings for regional security policy %s in %s in project %s", name, region, projectID)

	service, err := NewComputeServiceContextE(t, ctx)
	if err != nil {
		return nil, err
	}

	return FetchRegionSecurityPolicyWithClient(ctx, service, projectID, region, name)
}

// FetchRegionSecurityPolicyWithClient returns the settings Google Cloud holds for the given regional
// security policy using the supplied *compute.Service. Prefer this variant in unit tests where the
// service is backed by an httptest fake server (see cloudarmor_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func FetchRegionSecurityPolicyWithClient(ctx context.Context, service *compute.Service, projectID string, region string, name string) (*compute.SecurityPolicy, error) {
	policy, err := service.RegionSecurityPolicies.Get(projectID, region, name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the regional security policy %s does not exist in %s in project %s", name, region, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for regional security policy %s in %s in project %s: %w", name, region, projectID, err)
	}

	return policy, nil
}
