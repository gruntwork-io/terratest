package gcp

import (
	"context"
	"errors"
	"fmt"

	"github.com/gruntwork-io/terratest/modules/core/v2/logger"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/cloudidentity/v1"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
)

// GetCloudIdentityGroupAttrs returns the settings Google Cloud holds for the given Cloud Identity group, so a test can assert on
// what was actually created rather than only that it exists.
// A group is named groups/{id} with an id Google assigns, which is what Terraform records, so the
// caller passes the whole name rather than a project and an id.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetCloudIdentityGroupAttrs(t testing.TestingT, ctx context.Context, groupName string) *cloudidentity.Group {
	result, err := GetCloudIdentityGroupAttrsE(t, ctx, groupName)
	require.NoError(t, err)

	return result
}

// GetCloudIdentityGroupAttrsE returns the settings Google Cloud holds for the given Cloud Identity group.
// The ctx parameter supports cancellation and timeouts.
func GetCloudIdentityGroupAttrsE(t testing.TestingT, ctx context.Context, groupName string) (*cloudidentity.Group, error) {
	logger.Default.Logf(t, "Getting settings for Cloud Identity group %s", groupName)

	service, err := NewCloudIdentityServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetCloudIdentityGroupAttrsWithClient(ctx, service, groupName)
}

// GetCloudIdentityGroupAttrsWithClient returns the settings Google Cloud holds for the given Cloud Identity group using the supplied
// *cloudidentity.Service. Prefer this variant in unit tests where the service is backed by an httptest fake server
// (see cloudidentity_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetCloudIdentityGroupAttrsWithClient(ctx context.Context, service *cloudidentity.Service, groupName string) (*cloudidentity.Group, error) {
	result, err := service.Groups.Get(groupName).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Cloud Identity group %s does not exist", groupName)
		}

		return nil, fmt.Errorf("failed to get settings for Cloud Identity group %s: %w", groupName, err)
	}

	return result, nil
}

// NewCloudIdentityServiceE creates a Cloud Identity service authenticated the same way every other client in
// this module is.
// Reading a group needs the Cloud Identity groups scope rather than cloud-platform, which the API
// rejects.
// The ctx parameter supports cancellation and timeouts.
func NewCloudIdentityServiceE(t testing.TestingT, ctx context.Context) (*cloudidentity.Service, error) {
	return cloudidentity.NewService(ctx, append(withOptions(), option.WithScopes(cloudidentity.CloudIdentityGroupsScope))...)
}
