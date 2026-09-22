package gcp

import (
	"context"
	"errors"
	"fmt"

	"github.com/gruntwork-io/terratest/modules/core/v2/logger"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/certificatemanager/v1"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
)

// GetCertificateManagerCertificateAttrs returns the settings Google Cloud holds for the given
// Certificate Manager certificate, so a test can assert on what was actually created rather than
// only that it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetCertificateManagerCertificateAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, certificateID string) *certificatemanager.Certificate {
	certificate, err := GetCertificateManagerCertificateAttrsE(t, ctx, projectID, location, certificateID)
	require.NoError(t, err)

	return certificate
}

// GetCertificateManagerCertificateAttrsE returns the settings Google Cloud holds for the given
// Certificate Manager certificate.
// The ctx parameter supports cancellation and timeouts.
func GetCertificateManagerCertificateAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, certificateID string) (*certificatemanager.Certificate, error) {
	logger.Default.Logf(t, "Getting settings for Certificate Manager certificate %s in location %s in project %s", certificateID, location, projectID)

	service, err := NewCertificateManagerServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetCertificateManagerCertificateAttrsWithClient(ctx, service, projectID, location, certificateID)
}

// GetCertificateManagerCertificateAttrsWithClient returns the settings Google Cloud holds for the
// given Certificate Manager certificate using the supplied *certificatemanager.Service. Prefer this
// variant in unit tests where the service is backed by an httptest fake server (see
// certificatemanager_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetCertificateManagerCertificateAttrsWithClient(ctx context.Context, service *certificatemanager.Service, projectID string, location string, certificateID string) (*certificatemanager.Certificate, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/certificates/%s", projectID, location, certificateID)

	certificate, err := service.Projects.Locations.Certificates.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Certificate Manager certificate %s does not exist in location %s in project %s", certificateID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Certificate Manager certificate %s in location %s in project %s: %w", certificateID, location, projectID, err)
	}

	return certificate, nil
}

// NewCertificateManagerServiceE creates a Certificate Manager service authenticated the same way
// every other client in this module is.
// The ctx parameter supports cancellation and timeouts.
func NewCertificateManagerServiceE(t testing.TestingT, ctx context.Context) (*certificatemanager.Service, error) {
	return certificatemanager.NewService(ctx, append(withOptions(), option.WithScopes(certificatemanager.CloudPlatformScope))...)
}
