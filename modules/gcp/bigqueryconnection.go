package gcp

import (
	"context"
	"errors"
	"fmt"

	"github.com/gruntwork-io/terratest/modules/core/v2/logger"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/bigqueryconnection/v1"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
)

// GetBigQueryConnectionAttrs returns the settings Google Cloud holds for the given BigQuery connection, so a test can assert on
// what was actually created rather than only that it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetBigQueryConnectionAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, connectionID string) *bigqueryconnection.Connection {
	connection, err := GetBigQueryConnectionAttrsE(t, ctx, projectID, location, connectionID)
	require.NoError(t, err)

	return connection
}

// GetBigQueryConnectionAttrsE returns the settings Google Cloud holds for the given BigQuery connection.
// The ctx parameter supports cancellation and timeouts.
func GetBigQueryConnectionAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, connectionID string) (*bigqueryconnection.Connection, error) {
	logger.Default.Logf(t, "Getting settings for BigQuery connection %s in %s in project %s", connectionID, location, projectID)

	service, err := NewBigQueryConnectionServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetBigQueryConnectionAttrsWithClient(ctx, service, projectID, location, connectionID)
}

// GetBigQueryConnectionAttrsWithClient returns the settings Google Cloud holds for the given BigQuery connection using the supplied
// *bigqueryconnection.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see bigqueryconnection_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetBigQueryConnectionAttrsWithClient(ctx context.Context, service *bigqueryconnection.Service, projectID string, location string, connectionID string) (*bigqueryconnection.Connection, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/connections/%s", projectID, location, connectionID)

	connection, err := service.Projects.Locations.Connections.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the BigQuery connection %s does not exist in %s in project %s", connectionID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for BigQuery connection %s in %s in project %s: %w", connectionID, location, projectID, err)
	}

	return connection, nil
}

// NewBigQueryConnectionServiceE creates a BigQuery Connection service authenticated the same way every other client in this
// module is.
// The ctx parameter supports cancellation and timeouts.
func NewBigQueryConnectionServiceE(t testing.TestingT, ctx context.Context) (*bigqueryconnection.Service, error) {
	return bigqueryconnection.NewService(ctx, append(withOptions(), option.WithScopes(bigqueryconnection.CloudPlatformScope))...)
}
