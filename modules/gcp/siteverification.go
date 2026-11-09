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
	"google.golang.org/api/siteverification/v1"
)

// GetSiteVerificationWebResourceAttrs returns the settings Google Cloud holds for the given verified site or domain, so a test can assert on
// what was actually created rather than only that it exists.
// The id is the one the verification API uses, such as dns://example.com or
// https://www.example.com/, and not a project scoped resource name.
// One read answers for two modules: a web resource's owners are a field of the resource rather than a
// resource of their own.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetSiteVerificationWebResourceAttrs(t testing.TestingT, ctx context.Context, id string) *siteverification.SiteVerificationWebResourceResource {
	result, err := GetSiteVerificationWebResourceAttrsE(t, ctx, id)
	require.NoError(t, err)

	return result
}

// GetSiteVerificationWebResourceAttrsE returns the settings Google Cloud holds for the given verified site or domain.
// The ctx parameter supports cancellation and timeouts.
func GetSiteVerificationWebResourceAttrsE(t testing.TestingT, ctx context.Context, id string) (*siteverification.SiteVerificationWebResourceResource, error) {
	logger.Default.Logf(t, "Getting the verification state of %s", id)

	service, err := NewSiteVerificationServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetSiteVerificationWebResourceAttrsWithClient(ctx, service, id)
}

// GetSiteVerificationWebResourceAttrsWithClient returns the settings Google Cloud holds for the given verified site or domain using the supplied
// *siteverification.Service. Prefer this variant in unit tests where the service is backed by an httptest fake server
// (see siteverification_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetSiteVerificationWebResourceAttrsWithClient(ctx context.Context, service *siteverification.Service, id string) (*siteverification.SiteVerificationWebResourceResource, error) {
	result, err := service.WebResource.Get(id).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("%s is not a verified site or domain", id)
		}

		return nil, fmt.Errorf("failed to get the verification state of %s: %w", id, err)
	}

	return result, nil
}

// NewSiteVerificationServiceE creates a Site Verification service authenticated the same way every other client in
// this module is.
// The API offers no cloud-platform scope, so this client asks for the verification scope instead.
// The ctx parameter supports cancellation and timeouts.
func NewSiteVerificationServiceE(t testing.TestingT, ctx context.Context) (*siteverification.Service, error) {
	return siteverification.NewService(ctx, append(withOptions(), option.WithScopes(siteverification.SiteverificationScope))...)
}
