package gcp

import (
	"context"
	"errors"
	"fmt"

	"github.com/gruntwork-io/terratest/modules/core/v2/logger"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/binaryauthorization/v1"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
)

// GetBinaryAuthorizationAttestorAttrs returns the settings Google Cloud holds for the Binary Authorization attestor, so a test can assert
// on what was actually created rather than only that it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetBinaryAuthorizationAttestorAttrs(t testing.TestingT, ctx context.Context, projectID string, attestorID string) *binaryauthorization.Attestor {
	result, err := GetBinaryAuthorizationAttestorAttrsE(t, ctx, projectID, attestorID)
	require.NoError(t, err)

	return result
}

// GetBinaryAuthorizationAttestorAttrsE returns the settings Google Cloud holds for the Binary Authorization attestor.
// The ctx parameter supports cancellation and timeouts.
func GetBinaryAuthorizationAttestorAttrsE(t testing.TestingT, ctx context.Context, projectID string, attestorID string) (*binaryauthorization.Attestor, error) {
	logger.Default.Logf(t, "Getting settings for Binary Authorization attestor %s in project %s", attestorID, projectID)

	service, err := NewBinaryAuthorizationServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetBinaryAuthorizationAttestorAttrsWithClient(ctx, service, projectID, attestorID)
}

// GetBinaryAuthorizationAttestorAttrsWithClient returns the settings Google Cloud holds for the Binary Authorization attestor using the
// supplied *binaryauthorization.Service. Prefer this variant in unit tests where the service is backed by an
// httptest fake server (see binaryauthorization_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetBinaryAuthorizationAttestorAttrsWithClient(ctx context.Context, service *binaryauthorization.Service, projectID string, attestorID string) (*binaryauthorization.Attestor, error) {
	name := fmt.Sprintf("projects/%s/attestors/%s", projectID, attestorID)

	result, err := service.Projects.Attestors.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Binary Authorization attestor %s in project %s does not exist", attestorID, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Binary Authorization attestor %s in project %s: %w", attestorID, projectID, err)
	}

	return result, nil
}

// NewBinaryAuthorizationServiceE creates a Binary Authorization service authenticated the same way every other client in
// this module is.
// The ctx parameter supports cancellation and timeouts.
func NewBinaryAuthorizationServiceE(t testing.TestingT, ctx context.Context) (*binaryauthorization.Service, error) {
	return binaryauthorization.NewService(ctx, append(withOptions(), option.WithScopes(binaryauthorization.CloudPlatformScope))...)
}
