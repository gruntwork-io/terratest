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
	"google.golang.org/api/privateca/v1"
)

// GetCertificateTemplateAttrs returns the settings Google Cloud holds for the given certificate template, so a test can assert on
// what was actually created rather than only that it exists. A template shapes the certificates a pool will issue, and
// exists whether or not any pool references it.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetCertificateTemplateAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, templateID string) *privateca.CertificateTemplate {
	template, err := GetCertificateTemplateAttrsE(t, ctx, projectID, location, templateID)
	require.NoError(t, err)

	return template
}

// GetCertificateTemplateAttrsE returns the settings Google Cloud holds for the given certificate template.
// The ctx parameter supports cancellation and timeouts.
func GetCertificateTemplateAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, templateID string) (*privateca.CertificateTemplate, error) {
	logger.Default.Logf(t, "Getting settings for certificate template %s in %s in project %s", templateID, location, projectID)

	service, err := NewPrivateCAServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetCertificateTemplateAttrsWithClient(ctx, service, projectID, location, templateID)
}

// GetCertificateTemplateAttrsWithClient returns the settings Google Cloud holds for the given certificate template using the supplied
// *privateca.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see privateca_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetCertificateTemplateAttrsWithClient(ctx context.Context, service *privateca.Service, projectID string, location string, templateID string) (*privateca.CertificateTemplate, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/certificateTemplates/%s", projectID, location, templateID)

	template, err := service.Projects.Locations.CertificateTemplates.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the certificate template %s does not exist in %s in project %s", templateID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for certificate template %s in %s in project %s: %w", templateID, location, projectID, err)
	}

	return template, nil
}

// NewPrivateCAServiceE creates a Certificate Authority Service service authenticated the same way every other client in this
// module is.
// The ctx parameter supports cancellation and timeouts.
func NewPrivateCAServiceE(t testing.TestingT, ctx context.Context) (*privateca.Service, error) {
	return privateca.NewService(ctx, append(withOptions(), option.WithScopes(privateca.CloudPlatformScope))...)
}

// GetCaPoolAttrs returns the settings Google Cloud holds for the given CA pool, so a test can assert
// on the tier and issuance policy it was given rather than only that it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetCaPoolAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, poolID string) *privateca.CaPool {
	pool, err := GetCaPoolAttrsE(t, ctx, projectID, location, poolID)
	require.NoError(t, err)

	return pool
}

// GetCaPoolAttrsE returns the settings Google Cloud holds for the given CA pool.
// The ctx parameter supports cancellation and timeouts.
func GetCaPoolAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, poolID string) (*privateca.CaPool, error) {
	logger.Default.Logf(t, "Getting settings for CA pool %s in %s in project %s", poolID, location, projectID)

	service, err := NewPrivateCAServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetCaPoolAttrsWithClient(ctx, service, projectID, location, poolID)
}

// GetCaPoolAttrsWithClient returns the settings Google Cloud holds for the given CA pool using the
// supplied *privateca.Service. Prefer this variant in unit tests where the service is backed by an
// httptest fake server (see privateca_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetCaPoolAttrsWithClient(ctx context.Context, service *privateca.Service, projectID string, location string, poolID string) (*privateca.CaPool, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/caPools/%s", projectID, location, poolID)

	pool, err := service.Projects.Locations.CaPools.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the CA pool %s in %s in project %s does not exist", poolID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for CA pool %s in %s in project %s: %w", poolID, location, projectID, err)
	}

	return pool, nil
}

// GetPrivateCACertificateAttrs returns the settings Google Cloud holds for the given certificate that
// a CA pool issued, so a test can assert on its lifetime and subject rather than only that it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetPrivateCACertificateAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, poolID string, certificateID string) *privateca.Certificate {
	certificate, err := GetPrivateCACertificateAttrsE(t, ctx, projectID, location, poolID, certificateID)
	require.NoError(t, err)

	return certificate
}

// GetPrivateCACertificateAttrsE returns the settings Google Cloud holds for the given certificate
// that a CA pool issued.
// The ctx parameter supports cancellation and timeouts.
func GetPrivateCACertificateAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, poolID string, certificateID string) (*privateca.Certificate, error) {
	logger.Default.Logf(t, "Getting settings for certificate %s in CA pool %s in %s in project %s", certificateID, poolID, location, projectID)

	service, err := NewPrivateCAServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetPrivateCACertificateAttrsWithClient(ctx, service, projectID, location, poolID, certificateID)
}

// GetPrivateCACertificateAttrsWithClient returns the settings Google Cloud holds for the given
// certificate using the supplied *privateca.Service. Prefer this variant in unit tests where the
// service is backed by an httptest fake server (see privateca_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetPrivateCACertificateAttrsWithClient(ctx context.Context, service *privateca.Service, projectID string, location string, poolID string, certificateID string) (*privateca.Certificate, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/caPools/%s/certificates/%s", projectID, location, poolID, certificateID)

	certificate, err := service.Projects.Locations.CaPools.Certificates.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the certificate %s in CA pool %s in %s in project %s does not exist", certificateID, poolID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for certificate %s in CA pool %s in %s in project %s: %w", certificateID, poolID, location, projectID, err)
	}

	return certificate, nil
}
