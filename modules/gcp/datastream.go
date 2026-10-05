package gcp

import (
	"context"
	"errors"
	"fmt"

	"github.com/gruntwork-io/terratest/modules/core/v2/logger"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/datastream/v1"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
)

// GetDatastreamConnectionProfileAttrs returns the settings Google Cloud holds for the Datastream connection profile, so a test can assert
// on what was actually created rather than only that it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetDatastreamConnectionProfileAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, profileID string) *datastream.ConnectionProfile {
	result, err := GetDatastreamConnectionProfileAttrsE(t, ctx, projectID, location, profileID)
	require.NoError(t, err)

	return result
}

// GetDatastreamConnectionProfileAttrsE returns the settings Google Cloud holds for the Datastream connection profile.
// The ctx parameter supports cancellation and timeouts.
func GetDatastreamConnectionProfileAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, profileID string) (*datastream.ConnectionProfile, error) {
	logger.Default.Logf(t, "Getting settings for Datastream connection profile %s in project %s", profileID, projectID)

	service, err := NewDatastreamServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetDatastreamConnectionProfileAttrsWithClient(ctx, service, projectID, location, profileID)
}

// GetDatastreamConnectionProfileAttrsWithClient returns the settings Google Cloud holds for the Datastream connection profile using the
// supplied *datastream.Service. Prefer this variant in unit tests where the service is backed by an
// httptest fake server (see datastream_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetDatastreamConnectionProfileAttrsWithClient(ctx context.Context, service *datastream.Service, projectID string, location string, profileID string) (*datastream.ConnectionProfile, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/connectionProfiles/%s", projectID, location, profileID)

	result, err := service.Projects.Locations.ConnectionProfiles.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Datastream connection profile %s in project %s does not exist", profileID, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Datastream connection profile %s in project %s: %w", profileID, projectID, err)
	}

	return result, nil
}

// NewDatastreamServiceE creates a Datastream service authenticated the same way every other client in
// this module is.
// The ctx parameter supports cancellation and timeouts.
func NewDatastreamServiceE(t testing.TestingT, ctx context.Context) (*datastream.Service, error) {
	return datastream.NewService(ctx, append(withOptions(), option.WithScopes(datastream.CloudPlatformScope))...)
}

// GetDatastreamPrivateConnectionAttrs returns the settings Google Cloud holds for the given private
// connection, so a test can assert which network Datastream would reach a source through. A private
// connection is a peering Datastream holds open; a stream names it rather than the network.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetDatastreamPrivateConnectionAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, connectionID string) *datastream.PrivateConnection {
	connection, err := GetDatastreamPrivateConnectionAttrsE(t, ctx, projectID, location, connectionID)
	require.NoError(t, err)

	return connection
}

// GetDatastreamPrivateConnectionAttrsE returns the settings Google Cloud holds for the given private
// connection.
// The ctx parameter supports cancellation and timeouts.
func GetDatastreamPrivateConnectionAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, connectionID string) (*datastream.PrivateConnection, error) {
	logger.Default.Logf(t, "Getting settings for Datastream private connection %s in %s in project %s", connectionID, location, projectID)

	service, err := NewDatastreamServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetDatastreamPrivateConnectionAttrsWithClient(ctx, service, projectID, location, connectionID)
}

// GetDatastreamPrivateConnectionAttrsWithClient returns the settings Google Cloud holds for the given
// private connection using the supplied *datastream.Service. Prefer this variant in unit tests where the
// service is backed by an httptest fake server (see datastream_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetDatastreamPrivateConnectionAttrsWithClient(ctx context.Context, service *datastream.Service, projectID string, location string, connectionID string) (*datastream.PrivateConnection, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/privateConnections/%s", projectID, location, connectionID)

	connection, err := service.Projects.Locations.PrivateConnections.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Datastream private connection %s in %s in project %s does not exist", connectionID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Datastream private connection %s in %s in project %s: %w", connectionID, location, projectID, err)
	}

	return connection, nil
}
