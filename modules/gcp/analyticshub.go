package gcp

import (
	"context"
	"errors"
	"fmt"

	"github.com/gruntwork-io/terratest/modules/core/v2/logger"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/analyticshub/v1"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
)

// GetAnalyticsHubDataExchangeAttrs returns the settings Google Cloud holds for the given Analytics Hub data exchange, so a test can assert on
// what was actually created rather than only that it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetAnalyticsHubDataExchangeAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, exchangeID string) *analyticshub.DataExchange {
	exchange, err := GetAnalyticsHubDataExchangeAttrsE(t, ctx, projectID, location, exchangeID)
	require.NoError(t, err)

	return exchange
}

// GetAnalyticsHubDataExchangeAttrsE returns the settings Google Cloud holds for the given Analytics Hub data exchange.
// The ctx parameter supports cancellation and timeouts.
func GetAnalyticsHubDataExchangeAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, exchangeID string) (*analyticshub.DataExchange, error) {
	logger.Default.Logf(t, "Getting settings for data exchange %s in %s in project %s", exchangeID, location, projectID)

	service, err := NewAnalyticsHubServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetAnalyticsHubDataExchangeAttrsWithClient(ctx, service, projectID, location, exchangeID)
}

// GetAnalyticsHubDataExchangeAttrsWithClient returns the settings Google Cloud holds for the given Analytics Hub data exchange using the supplied
// *analyticshub.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see analyticshub_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetAnalyticsHubDataExchangeAttrsWithClient(ctx context.Context, service *analyticshub.Service, projectID string, location string, exchangeID string) (*analyticshub.DataExchange, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/dataExchanges/%s", projectID, location, exchangeID)

	exchange, err := service.Projects.Locations.DataExchanges.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the data exchange %s does not exist in %s in project %s", exchangeID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for data exchange %s in %s in project %s: %w", exchangeID, location, projectID, err)
	}

	return exchange, nil
}

// GetAnalyticsHubListingAttrs returns the settings Google Cloud holds for the given Analytics Hub listing, so a test can assert on
// what was actually created rather than only that it exists. A listing belongs to a data exchange, so it is named by the
// exchange's id as well as its own.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetAnalyticsHubListingAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, exchangeID string, listingID string) *analyticshub.Listing {
	listing, err := GetAnalyticsHubListingAttrsE(t, ctx, projectID, location, exchangeID, listingID)
	require.NoError(t, err)

	return listing
}

// GetAnalyticsHubListingAttrsE returns the settings Google Cloud holds for the given Analytics Hub listing.
// The ctx parameter supports cancellation and timeouts.
func GetAnalyticsHubListingAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, exchangeID string, listingID string) (*analyticshub.Listing, error) {
	logger.Default.Logf(t, "Getting settings for listing %s in data exchange %s in %s in project %s", listingID, exchangeID, location, projectID)

	service, err := NewAnalyticsHubServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetAnalyticsHubListingAttrsWithClient(ctx, service, projectID, location, exchangeID, listingID)
}

// GetAnalyticsHubListingAttrsWithClient returns the settings Google Cloud holds for the given Analytics Hub listing using the supplied
// *analyticshub.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see analyticshub_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetAnalyticsHubListingAttrsWithClient(ctx context.Context, service *analyticshub.Service, projectID string, location string, exchangeID string, listingID string) (*analyticshub.Listing, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/dataExchanges/%s/listings/%s", projectID, location, exchangeID, listingID)

	listing, err := service.Projects.Locations.DataExchanges.Listings.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the listing %s does not exist in data exchange %s in %s in project %s", listingID, exchangeID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for listing %s in data exchange %s in %s in project %s: %w", listingID, exchangeID, location, projectID, err)
	}

	return listing, nil
}

// NewAnalyticsHubServiceE creates a Analytics Hub service authenticated the same way every other client in this
// module is.
// The ctx parameter supports cancellation and timeouts.
func NewAnalyticsHubServiceE(t testing.TestingT, ctx context.Context) (*analyticshub.Service, error) {
	return analyticshub.NewService(ctx, append(withOptions(), option.WithScopes(analyticshub.CloudPlatformScope))...)
}
