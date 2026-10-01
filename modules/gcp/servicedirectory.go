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
	"google.golang.org/api/servicedirectory/v1"
)

// GetServiceDirectoryNamespaceAttrs returns the settings Google Cloud holds for the given namespace, so a test can assert on
// what was actually created rather than only that it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetServiceDirectoryNamespaceAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, namespaceID string) *servicedirectory.Namespace {
	namespace, err := GetServiceDirectoryNamespaceAttrsE(t, ctx, projectID, location, namespaceID)
	require.NoError(t, err)

	return namespace
}

// GetServiceDirectoryNamespaceAttrsE returns the settings Google Cloud holds for the given namespace.
// The ctx parameter supports cancellation and timeouts.
func GetServiceDirectoryNamespaceAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, namespaceID string) (*servicedirectory.Namespace, error) {
	logger.Default.Logf(t, "Getting settings for namespace %s in %s in project %s", namespaceID, location, projectID)

	service, err := NewServiceDirectoryServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetServiceDirectoryNamespaceAttrsWithClient(ctx, service, projectID, location, namespaceID)
}

// GetServiceDirectoryNamespaceAttrsWithClient returns the settings Google Cloud holds for the given namespace using the supplied
// *servicedirectory.APIService. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see servicedirectory_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetServiceDirectoryNamespaceAttrsWithClient(ctx context.Context, service *servicedirectory.APIService, projectID string, location string, namespaceID string) (*servicedirectory.Namespace, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/namespaces/%s", projectID, location, namespaceID)

	namespace, err := service.Projects.Locations.Namespaces.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the namespace %s does not exist in %s in project %s", namespaceID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for namespace %s in %s in project %s: %w", namespaceID, location, projectID, err)
	}

	return namespace, nil
}

// GetServiceDirectoryServiceAttrs returns the settings Google Cloud holds for the given registered service, so a test can assert on
// what was actually created rather than only that it exists. The variable is called registered because Service is the
// name of the client itself.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetServiceDirectoryServiceAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, namespaceID string, serviceID string) *servicedirectory.Service {
	registered, err := GetServiceDirectoryServiceAttrsE(t, ctx, projectID, location, namespaceID, serviceID)
	require.NoError(t, err)

	return registered
}

// GetServiceDirectoryServiceAttrsE returns the settings Google Cloud holds for the given registered service.
// The ctx parameter supports cancellation and timeouts.
func GetServiceDirectoryServiceAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, namespaceID string, serviceID string) (*servicedirectory.Service, error) {
	logger.Default.Logf(t, "Getting settings for service %s in namespace %s in %s in project %s", serviceID, namespaceID, location, projectID)

	service, err := NewServiceDirectoryServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetServiceDirectoryServiceAttrsWithClient(ctx, service, projectID, location, namespaceID, serviceID)
}

// GetServiceDirectoryServiceAttrsWithClient returns the settings Google Cloud holds for the given registered service using the supplied
// *servicedirectory.APIService. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see servicedirectory_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetServiceDirectoryServiceAttrsWithClient(ctx context.Context, service *servicedirectory.APIService, projectID string, location string, namespaceID string, serviceID string) (*servicedirectory.Service, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/namespaces/%s/services/%s", projectID, location, namespaceID, serviceID)

	registered, err := service.Projects.Locations.Namespaces.Services.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the service %s does not exist in namespace %s in %s in project %s", serviceID, namespaceID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for service %s in namespace %s in %s in project %s: %w", serviceID, namespaceID, location, projectID, err)
	}

	return registered, nil
}

// GetServiceDirectoryEndpointAttrs returns the settings Google Cloud holds for the given endpoint, so a test can assert on
// what was actually created rather than only that it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetServiceDirectoryEndpointAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, namespaceID string, serviceID string, endpointID string) *servicedirectory.Endpoint {
	endpoint, err := GetServiceDirectoryEndpointAttrsE(t, ctx, projectID, location, namespaceID, serviceID, endpointID)
	require.NoError(t, err)

	return endpoint
}

// GetServiceDirectoryEndpointAttrsE returns the settings Google Cloud holds for the given endpoint.
// The ctx parameter supports cancellation and timeouts.
func GetServiceDirectoryEndpointAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, namespaceID string, serviceID string, endpointID string) (*servicedirectory.Endpoint, error) {
	logger.Default.Logf(t, "Getting settings for endpoint %s on service %s in namespace %s in %s in project %s", endpointID, serviceID, namespaceID, location, projectID)

	service, err := NewServiceDirectoryServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetServiceDirectoryEndpointAttrsWithClient(ctx, service, projectID, location, namespaceID, serviceID, endpointID)
}

// GetServiceDirectoryEndpointAttrsWithClient returns the settings Google Cloud holds for the given endpoint using the supplied
// *servicedirectory.APIService. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see servicedirectory_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetServiceDirectoryEndpointAttrsWithClient(ctx context.Context, service *servicedirectory.APIService, projectID string, location string, namespaceID string, serviceID string, endpointID string) (*servicedirectory.Endpoint, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/namespaces/%s/services/%s/endpoints/%s", projectID, location, namespaceID, serviceID, endpointID)

	endpoint, err := service.Projects.Locations.Namespaces.Services.Endpoints.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the endpoint %s does not exist on service %s in namespace %s in %s in project %s", endpointID, serviceID, namespaceID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for endpoint %s on service %s in namespace %s in %s in project %s: %w", endpointID, serviceID, namespaceID, location, projectID, err)
	}

	return endpoint, nil
}

// NewServiceDirectoryServiceE creates a Service Directory client authenticated the same way every
// other client in this module is. The generated client is called APIService, because Service is the
// name of a resource in this API.
// The ctx parameter supports cancellation and timeouts.
func NewServiceDirectoryServiceE(t testing.TestingT, ctx context.Context) (*servicedirectory.APIService, error) {
	return servicedirectory.NewService(ctx, append(withOptions(), option.WithScopes(servicedirectory.CloudPlatformScope))...)
}

// GetServiceDirectoryNamespaceIamPolicyAttrs returns the IAM policy Google Cloud holds for the given
// namespace, so a test can assert on who was actually granted access to it.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetServiceDirectoryNamespaceIamPolicyAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, namespaceID string) *servicedirectory.Policy {
	policy, err := GetServiceDirectoryNamespaceIamPolicyAttrsE(t, ctx, projectID, location, namespaceID)
	require.NoError(t, err)

	return policy
}

// GetServiceDirectoryNamespaceIamPolicyAttrsE returns the IAM policy Google Cloud holds for the given namespace.
// The ctx parameter supports cancellation and timeouts.
func GetServiceDirectoryNamespaceIamPolicyAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, namespaceID string) (*servicedirectory.Policy, error) {
	logger.Default.Logf(t, "Getting the IAM policy for namespace %s in %s in project %s", namespaceID, location, projectID)

	service, err := NewServiceDirectoryServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetServiceDirectoryNamespaceIamPolicyAttrsWithClient(ctx, service, projectID, location, namespaceID)
}

// GetServiceDirectoryNamespaceIamPolicyAttrsWithClient returns the IAM policy Google Cloud holds for the
// given namespace using the supplied *servicedirectory.APIService. Prefer this variant in unit tests
// where the service is backed by an httptest fake server (see servicedirectory_test.go).
// The ctx parameter supports cancellation and timeouts.
func GetServiceDirectoryNamespaceIamPolicyAttrsWithClient(ctx context.Context, service *servicedirectory.APIService, projectID string, location string, namespaceID string) (*servicedirectory.Policy, error) {
	resource := fmt.Sprintf("projects/%s/locations/%s/namespaces/%s", projectID, location, namespaceID)

	// This call takes a request body rather than a plain resource name. A policy carrying a conditional
	// binding is only returned in full at version 3, so that is what is asked for.
	policy, err := service.Projects.Locations.Namespaces.GetIamPolicy(resource, &servicedirectory.GetIamPolicyRequest{
		Options: &servicedirectory.GetPolicyOptions{RequestedPolicyVersion: iamPolicyVersionWithConditions},
	}).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the namespace %s in %s in project %s does not exist", namespaceID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get the IAM policy for namespace %s in %s in project %s: %w", namespaceID, location, projectID, err)
	}

	return policy, nil
}

// GetServiceDirectoryServiceIamPolicyAttrs returns the IAM policy Google Cloud holds for the given
// service, so a test can assert on who was actually granted access to it.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetServiceDirectoryServiceIamPolicyAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, namespaceID string, serviceID string) *servicedirectory.Policy {
	policy, err := GetServiceDirectoryServiceIamPolicyAttrsE(t, ctx, projectID, location, namespaceID, serviceID)
	require.NoError(t, err)

	return policy
}

// GetServiceDirectoryServiceIamPolicyAttrsE returns the IAM policy Google Cloud holds for the given service.
// The ctx parameter supports cancellation and timeouts.
func GetServiceDirectoryServiceIamPolicyAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, namespaceID string, serviceID string) (*servicedirectory.Policy, error) {
	logger.Default.Logf(t, "Getting the IAM policy for service %s in namespace %s in %s in project %s", serviceID, namespaceID, location, projectID)

	service, err := NewServiceDirectoryServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetServiceDirectoryServiceIamPolicyAttrsWithClient(ctx, service, projectID, location, namespaceID, serviceID)
}

// GetServiceDirectoryServiceIamPolicyAttrsWithClient returns the IAM policy Google Cloud holds for the
// given service using the supplied *servicedirectory.APIService. Prefer this variant in unit tests
// where the service is backed by an httptest fake server (see servicedirectory_test.go).
// The ctx parameter supports cancellation and timeouts.
func GetServiceDirectoryServiceIamPolicyAttrsWithClient(ctx context.Context, service *servicedirectory.APIService, projectID string, location string, namespaceID string, serviceID string) (*servicedirectory.Policy, error) {
	resource := fmt.Sprintf("projects/%s/locations/%s/namespaces/%s/services/%s", projectID, location, namespaceID, serviceID)

	// This call takes a request body rather than a plain resource name. A policy carrying a conditional
	// binding is only returned in full at version 3, so that is what is asked for.
	policy, err := service.Projects.Locations.Namespaces.Services.GetIamPolicy(resource, &servicedirectory.GetIamPolicyRequest{
		Options: &servicedirectory.GetPolicyOptions{RequestedPolicyVersion: iamPolicyVersionWithConditions},
	}).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the service %s in namespace %s in %s in project %s does not exist", serviceID, namespaceID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get the IAM policy for service %s in namespace %s in %s in project %s: %w", serviceID, namespaceID, location, projectID, err)
	}

	return policy, nil
}
