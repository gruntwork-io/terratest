package gcp

import (
	"context"
	"errors"
	"fmt"

	"github.com/gruntwork-io/terratest/modules/core/v2/logger"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/bigqueryreservation/v1"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
)

// GetBigQueryBiReservationAttrs returns the settings Google Cloud holds for the BI Engine reservation in the given location, so a test can assert on
// what was actually created rather than only that it exists. A project holds one per location, which exists with a size of
// zero whether anyone asked for it or not, so what is read back is its size.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetBigQueryBiReservationAttrs(t testing.TestingT, ctx context.Context, projectID string, location string) *bigqueryreservation.BiReservation {
	reservation, err := GetBigQueryBiReservationAttrsE(t, ctx, projectID, location)
	require.NoError(t, err)

	return reservation
}

// GetBigQueryBiReservationAttrsE returns the settings Google Cloud holds for the BI Engine reservation in the given location.
// The ctx parameter supports cancellation and timeouts.
func GetBigQueryBiReservationAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string) (*bigqueryreservation.BiReservation, error) {
	logger.Default.Logf(t, "Getting settings for the BI Engine reservation in %s in project %s", location, projectID)

	service, err := NewBigQueryReservationServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetBigQueryBiReservationAttrsWithClient(ctx, service, projectID, location)
}

// GetBigQueryBiReservationAttrsWithClient returns the settings Google Cloud holds for the BI Engine reservation in the given location using the supplied
// *bigqueryreservation.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see bigqueryreservation_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetBigQueryBiReservationAttrsWithClient(ctx context.Context, service *bigqueryreservation.Service, projectID string, location string) (*bigqueryreservation.BiReservation, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/biReservation", projectID, location)

	reservation, err := service.Projects.Locations.GetBiReservation(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("there is no BI Engine reservation in %s in project %s", location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for the BI Engine reservation in %s in project %s: %w", location, projectID, err)
	}

	return reservation, nil
}

// NewBigQueryReservationServiceE creates a BigQuery Reservation service authenticated the same way every other client in this
// module is.
// The ctx parameter supports cancellation and timeouts.
func NewBigQueryReservationServiceE(t testing.TestingT, ctx context.Context) (*bigqueryreservation.Service, error) {
	return bigqueryreservation.NewService(ctx, append(withOptions(), option.WithScopes(bigqueryreservation.CloudPlatformScope))...)
}
