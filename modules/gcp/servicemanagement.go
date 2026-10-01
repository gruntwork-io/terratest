package gcp

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/gruntwork-io/terratest/modules/core/v2/logger"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
	"google.golang.org/api/servicemanagement/v1"
)

// GetEndpointsServiceAttrs returns the settings Google Cloud holds for the given Cloud Endpoints
// service, so a test can assert which project produces it. An Endpoints service is named by its own
// DNS name rather than by a path under the project, so the caller passes that name.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetEndpointsServiceAttrs(t testing.TestingT, ctx context.Context, serviceName string) *servicemanagement.ManagedService {
	service, err := GetEndpointsServiceAttrsE(t, ctx, serviceName)
	require.NoError(t, err)

	return service
}

// GetEndpointsServiceAttrsE returns the settings Google Cloud holds for the given Cloud Endpoints
// service.
// The ctx parameter supports cancellation and timeouts.
func GetEndpointsServiceAttrsE(t testing.TestingT, ctx context.Context, serviceName string) (*servicemanagement.ManagedService, error) {
	logger.Default.Logf(t, "Getting settings for Endpoints service %s", serviceName)

	client, err := NewServiceManagementServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetEndpointsServiceAttrsWithClient(ctx, client, serviceName)
}

// GetEndpointsServiceAttrsWithClient returns the settings Google Cloud holds for the given Cloud
// Endpoints service using the supplied *servicemanagement.APIService. Prefer this variant in unit tests
// where the service is backed by an httptest fake server (see servicemanagement_test.go for the
// pattern).
// The ctx parameter supports cancellation and timeouts.
func GetEndpointsServiceAttrsWithClient(ctx context.Context, client *servicemanagement.APIService, serviceName string) (*servicemanagement.ManagedService, error) {
	managed, err := client.Services.Get(serviceName).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Endpoints service %s does not exist", serviceName)
		}

		return nil, fmt.Errorf("failed to get settings for Endpoints service %s: %w", serviceName, err)
	}

	return managed, nil
}

// GetEndpointsServiceIamPolicyAttrs returns the IAM policy Google Cloud holds for the given Cloud
// Endpoints service, so a test can assert on who may manage it.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetEndpointsServiceIamPolicyAttrs(t testing.TestingT, ctx context.Context, serviceName string) *servicemanagement.Policy {
	policy, err := GetEndpointsServiceIamPolicyAttrsE(t, ctx, serviceName)
	require.NoError(t, err)

	return policy
}

// GetEndpointsServiceIamPolicyAttrsE returns the IAM policy Google Cloud holds for the given Cloud
// Endpoints service.
// The ctx parameter supports cancellation and timeouts.
func GetEndpointsServiceIamPolicyAttrsE(t testing.TestingT, ctx context.Context, serviceName string) (*servicemanagement.Policy, error) {
	logger.Default.Logf(t, "Getting the IAM policy for Endpoints service %s", serviceName)

	client, err := NewServiceManagementServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetEndpointsServiceIamPolicyAttrsWithClient(ctx, client, serviceName)
}

// GetEndpointsServiceIamPolicyAttrsWithClient returns the IAM policy Google Cloud holds for the given
// Cloud Endpoints service using the supplied *servicemanagement.APIService. Prefer this variant in unit
// tests where the service is backed by an httptest fake server (see servicemanagement_test.go for the
// pattern).
// The ctx parameter supports cancellation and timeouts.
func GetEndpointsServiceIamPolicyAttrsWithClient(ctx context.Context, client *servicemanagement.APIService, serviceName string) (*servicemanagement.Policy, error) {
	// The policy hangs off the service's own resource path rather than off its bare DNS name, so the
	// prefix is added here. A caller who already holds the whole path passes it through unchanged.
	resource := "services/" + serviceName
	if strings.HasPrefix(serviceName, "services/") {
		resource = serviceName
	}

	// This call takes a request body rather than a plain resource name. A policy carrying a conditional
	// binding is only returned in full at version 3, so that is what is asked for.
	policy, err := client.Services.GetIamPolicy(resource, &servicemanagement.GetIamPolicyRequest{
		Options: &servicemanagement.GetPolicyOptions{RequestedPolicyVersion: iamPolicyVersionWithConditions},
	}).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Endpoints service %s does not exist", serviceName)
		}

		return nil, fmt.Errorf("failed to get the IAM policy for Endpoints service %s: %w", serviceName, err)
	}

	return policy, nil
}

// GetEndpointsServiceConsumerIamPolicyAttrs returns the IAM policy Google Cloud holds for one consumer
// of the given Cloud Endpoints service, so a test can assert on who in that project may call it. The
// consumer policy is separate from the service's own, and the two answer different questions.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetEndpointsServiceConsumerIamPolicyAttrs(t testing.TestingT, ctx context.Context, serviceName string, consumerProject string) *servicemanagement.Policy {
	policy, err := GetEndpointsServiceConsumerIamPolicyAttrsE(t, ctx, serviceName, consumerProject)
	require.NoError(t, err)

	return policy
}

// GetEndpointsServiceConsumerIamPolicyAttrsE returns the IAM policy Google Cloud holds for one consumer
// of the given Cloud Endpoints service.
// The ctx parameter supports cancellation and timeouts.
func GetEndpointsServiceConsumerIamPolicyAttrsE(t testing.TestingT, ctx context.Context, serviceName string, consumerProject string) (*servicemanagement.Policy, error) {
	logger.Default.Logf(t, "Getting the IAM policy for consumer %s of Endpoints service %s", consumerProject, serviceName)

	client, err := NewServiceManagementServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetEndpointsServiceConsumerIamPolicyAttrsWithClient(ctx, client, serviceName, consumerProject)
}

// GetEndpointsServiceConsumerIamPolicyAttrsWithClient returns the IAM policy Google Cloud holds for one
// consumer of the given Cloud Endpoints service using the supplied *servicemanagement.APIService.
// Prefer this variant in unit tests where the service is backed by an httptest fake server (see
// servicemanagement_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetEndpointsServiceConsumerIamPolicyAttrsWithClient(ctx context.Context, client *servicemanagement.APIService, serviceName string, consumerProject string) (*servicemanagement.Policy, error) {
	// Google's reference spells a consumer as project:<id>, but the API answers 404 for that spelling
	// here and serves the policy under the bare project id, which is what the provider writes to.
	resource := fmt.Sprintf("services/%s/consumers/%s", serviceName, consumerProject)

	policy, err := client.Services.Consumers.GetIamPolicy(resource, &servicemanagement.GetIamPolicyRequest{
		Options: &servicemanagement.GetPolicyOptions{RequestedPolicyVersion: iamPolicyVersionWithConditions},
	}).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the consumer %s of Endpoints service %s does not exist", consumerProject, serviceName)
		}

		return nil, fmt.Errorf("failed to get the IAM policy for consumer %s of Endpoints service %s: %w", consumerProject, serviceName, err)
	}

	return policy, nil
}

// NewServiceManagementServiceE creates a Service Management service authenticated the same way every
// other client in this module is. The Go client calls this type APIService, because Service Management
// already spends the name Service on a resource.
// The ctx parameter supports cancellation and timeouts.
func NewServiceManagementServiceE(t testing.TestingT, ctx context.Context) (*servicemanagement.APIService, error) {
	return servicemanagement.NewService(ctx, append(withOptions(), option.WithScopes(servicemanagement.CloudPlatformScope))...)
}
