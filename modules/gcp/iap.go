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

// NewIAPServiceE creates an Identity Aware Proxy service authenticated the same way every other
// client in this module is.
// The ctx parameter supports cancellation and timeouts.
func NewIAPServiceE(t testing.TestingT, ctx context.Context) (*iap.Service, error) {
	return iap.NewService(ctx, append(withOptions(), option.WithScopes(iap.CloudPlatformScope))...)
}
