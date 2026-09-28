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

// GetDNSAuthorizationAttrs returns the settings Google Cloud holds for the given DNS authorization, so a test can assert on
// what was actually created rather than only that it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetDNSAuthorizationAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, authorizationID string) *certificatemanager.DnsAuthorization {
	authorization, err := GetDNSAuthorizationAttrsE(t, ctx, projectID, location, authorizationID)
	require.NoError(t, err)

	return authorization
}

// GetDNSAuthorizationAttrsE returns the settings Google Cloud holds for the given DNS authorization.
// The ctx parameter supports cancellation and timeouts.
func GetDNSAuthorizationAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, authorizationID string) (*certificatemanager.DnsAuthorization, error) {
	logger.Default.Logf(t, "Getting settings for DNS authorization %s in %s in project %s", authorizationID, location, projectID)

	service, err := NewCertificateManagerServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetDNSAuthorizationAttrsWithClient(ctx, service, projectID, location, authorizationID)
}

// GetDNSAuthorizationAttrsWithClient returns the settings Google Cloud holds for the given DNS authorization using the supplied
// *certificatemanager.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see certificatemanager_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetDNSAuthorizationAttrsWithClient(ctx context.Context, service *certificatemanager.Service, projectID string, location string, authorizationID string) (*certificatemanager.DnsAuthorization, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/dnsAuthorizations/%s", projectID, location, authorizationID)

	authorization, err := service.Projects.Locations.DnsAuthorizations.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the DNS authorization %s does not exist in %s in project %s", authorizationID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for DNS authorization %s in %s in project %s: %w", authorizationID, location, projectID, err)
	}

	return authorization, nil
}

// GetCertificateMapAttrs returns the settings Google Cloud holds for the given certificate map, so a test can assert on
// what was actually created rather than only that it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetCertificateMapAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, mapID string) *certificatemanager.CertificateMap {
	certificateMap, err := GetCertificateMapAttrsE(t, ctx, projectID, location, mapID)
	require.NoError(t, err)

	return certificateMap
}

// GetCertificateMapAttrsE returns the settings Google Cloud holds for the given certificate map.
// The ctx parameter supports cancellation and timeouts.
func GetCertificateMapAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, mapID string) (*certificatemanager.CertificateMap, error) {
	logger.Default.Logf(t, "Getting settings for certificate map %s in %s in project %s", mapID, location, projectID)

	service, err := NewCertificateManagerServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetCertificateMapAttrsWithClient(ctx, service, projectID, location, mapID)
}

// GetCertificateMapAttrsWithClient returns the settings Google Cloud holds for the given certificate map using the supplied
// *certificatemanager.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see certificatemanager_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetCertificateMapAttrsWithClient(ctx context.Context, service *certificatemanager.Service, projectID string, location string, mapID string) (*certificatemanager.CertificateMap, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/certificateMaps/%s", projectID, location, mapID)

	certificateMap, err := service.Projects.Locations.CertificateMaps.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the certificate map %s does not exist in %s in project %s", mapID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for certificate map %s in %s in project %s: %w", mapID, location, projectID, err)
	}

	return certificateMap, nil
}

// GetTrustConfigAttrs returns the settings Google Cloud holds for the given trust configuration, so a test can assert on
// what was actually created rather than only that it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetTrustConfigAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, configID string) *certificatemanager.TrustConfig {
	trustConfig, err := GetTrustConfigAttrsE(t, ctx, projectID, location, configID)
	require.NoError(t, err)

	return trustConfig
}

// GetTrustConfigAttrsE returns the settings Google Cloud holds for the given trust configuration.
// The ctx parameter supports cancellation and timeouts.
func GetTrustConfigAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, configID string) (*certificatemanager.TrustConfig, error) {
	logger.Default.Logf(t, "Getting settings for trust configuration %s in %s in project %s", configID, location, projectID)

	service, err := NewCertificateManagerServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetTrustConfigAttrsWithClient(ctx, service, projectID, location, configID)
}

// GetTrustConfigAttrsWithClient returns the settings Google Cloud holds for the given trust configuration using the supplied
// *certificatemanager.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see certificatemanager_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetTrustConfigAttrsWithClient(ctx context.Context, service *certificatemanager.Service, projectID string, location string, configID string) (*certificatemanager.TrustConfig, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/trustConfigs/%s", projectID, location, configID)

	trustConfig, err := service.Projects.Locations.TrustConfigs.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the trust configuration %s does not exist in %s in project %s", configID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for trust configuration %s in %s in project %s: %w", configID, location, projectID, err)
	}

	return trustConfig, nil
}
