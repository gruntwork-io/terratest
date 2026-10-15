package gcp

import (
	"context"
	"errors"
	"fmt"

	"github.com/gruntwork-io/terratest/modules/core/v2/logger"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/iap/v1"
	"google.golang.org/api/option"
)

// GetIAPIamPolicyAttrs returns the IAM policy Google Cloud holds for the given Identity Aware Proxy
// resource, so a test can assert on who may reach it through IAP rather than only that a policy was
// applied. Every IAP resource answers on one method, and they differ only in the resource name, so
// the caller passes that name in full: for example
// projects/PROJECT/iap_web/compute/services/BACKEND_SERVICE.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetIAPIamPolicyAttrs(t testing.TestingT, ctx context.Context, resource string) *iap.Policy {
	policy, err := GetIAPIamPolicyAttrsE(t, ctx, resource)
	require.NoError(t, err)

	return policy
}

// GetIAPIamPolicyAttrsE returns the IAM policy Google Cloud holds for the given Identity Aware Proxy
// resource.
// The ctx parameter supports cancellation and timeouts.
func GetIAPIamPolicyAttrsE(t testing.TestingT, ctx context.Context, resource string) (*iap.Policy, error) {
	logger.Default.Logf(t, "Getting the IAP IAM policy for %s", resource)

	service, err := NewIAPServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetIAPIamPolicyAttrsWithClient(ctx, service, resource)
}

// GetIAPIamPolicyAttrsWithClient returns the IAM policy Google Cloud holds for the given Identity
// Aware Proxy resource using the supplied *iap.Service. Prefer this variant in unit tests where the
// service is backed by an httptest fake server (see iap_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetIAPIamPolicyAttrsWithClient(ctx context.Context, service *iap.Service, resource string) (*iap.Policy, error) {
	// This call takes a request body rather than a plain resource name. A policy carrying a
	// conditional binding is only returned in full at version 3, so that is what is asked for.
	policy, err := service.V1.GetIamPolicy(resource, &iap.GetIamPolicyRequest{
		Options: &iap.GetPolicyOptions{RequestedPolicyVersion: iamPolicyVersionWithConditions},
	}).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the IAP resource %s does not exist", resource)
		}

		return nil, fmt.Errorf("failed to get the IAP IAM policy for %s: %w", resource, err)
	}

	return policy, nil
}

// GetIAPBrandAttrs returns the settings Google Cloud holds for the given Identity-Aware Proxy brand, so a test can assert on
// what was actually created rather than only that it exists.
// Whether Google treats the brand as internal only is decided at create time from the support email's
// domain, so a read is the only way to see it.
// The brand id is the project number, which is also what a create answers with.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetIAPBrandAttrs(t testing.TestingT, ctx context.Context, projectID string, brandID string) *iap.Brand {
	result, err := GetIAPBrandAttrsE(t, ctx, projectID, brandID)
	require.NoError(t, err)

	return result
}

// GetIAPBrandAttrsE returns the settings Google Cloud holds for the given Identity-Aware Proxy brand.
// The ctx parameter supports cancellation and timeouts.
func GetIAPBrandAttrsE(t testing.TestingT, ctx context.Context, projectID string, brandID string) (*iap.Brand, error) {
	logger.Default.Logf(t, "Getting settings for Identity-Aware Proxy brand %s in project %s", brandID, projectID)

	service, err := NewIAPServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetIAPBrandAttrsWithClient(ctx, service, projectID, brandID)
}

// GetIAPBrandAttrsWithClient returns the settings Google Cloud holds for the given Identity-Aware Proxy brand using the supplied
// *iap.Service. Prefer this variant in unit tests where the service is backed by an httptest fake server
// (see iap_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetIAPBrandAttrsWithClient(ctx context.Context, service *iap.Service, projectID string, brandID string) (*iap.Brand, error) {
	name := fmt.Sprintf("projects/%s/brands/%s", projectID, brandID)

	result, err := service.Projects.Brands.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Identity-Aware Proxy brand %s does not exist in project %s", brandID, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Identity-Aware Proxy brand %s in project %s: %w", brandID, projectID, err)
	}

	return result, nil
}

// GetIAPTunnelDestGroupAttrs returns the settings Google Cloud holds for the given Identity-Aware Proxy tunnel destination group, so a test can assert on
// what was actually created rather than only that it exists.
// The tunnel segment of the name is spelled with an underscore, which is how the API serves it and not
// a typo.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetIAPTunnelDestGroupAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, groupID string) *iap.TunnelDestGroup {
	result, err := GetIAPTunnelDestGroupAttrsE(t, ctx, projectID, location, groupID)
	require.NoError(t, err)

	return result
}

// GetIAPTunnelDestGroupAttrsE returns the settings Google Cloud holds for the given Identity-Aware Proxy tunnel destination group.
// The ctx parameter supports cancellation and timeouts.
func GetIAPTunnelDestGroupAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, groupID string) (*iap.TunnelDestGroup, error) {
	logger.Default.Logf(t, "Getting settings for Identity-Aware Proxy tunnel destination group %s in %s of project %s", groupID, location, projectID)

	service, err := NewIAPServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetIAPTunnelDestGroupAttrsWithClient(ctx, service, projectID, location, groupID)
}

// GetIAPTunnelDestGroupAttrsWithClient returns the settings Google Cloud holds for the given Identity-Aware Proxy tunnel destination group using the supplied
// *iap.Service. Prefer this variant in unit tests where the service is backed by an httptest fake server
// (see iap_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetIAPTunnelDestGroupAttrsWithClient(ctx context.Context, service *iap.Service, projectID string, location string, groupID string) (*iap.TunnelDestGroup, error) {
	name := fmt.Sprintf("projects/%s/iap_tunnel/locations/%s/destGroups/%s", projectID, location, groupID)

	result, err := service.Projects.IapTunnel.Locations.DestGroups.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Identity-Aware Proxy tunnel destination group %s does not exist in %s of project %s", groupID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Identity-Aware Proxy tunnel destination group %s in %s of project %s: %w", groupID, location, projectID, err)
	}

	return result, nil
}

// NewIAPServiceE creates an Identity Aware Proxy service authenticated the same way every other
// client in this module is.
// The ctx parameter supports cancellation and timeouts.
func NewIAPServiceE(t testing.TestingT, ctx context.Context) (*iap.Service, error) {
	return iap.NewService(ctx, append(withOptions(), option.WithScopes(iap.CloudPlatformScope))...)
}
