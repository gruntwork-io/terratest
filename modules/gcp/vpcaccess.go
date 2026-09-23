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
	"google.golang.org/api/vpcaccess/v1"
)

// GetVPCAccessConnectorAttrs returns the settings Google Cloud holds for the given Serverless VPC
// Access connector, so a test can assert on what was actually created rather than only that it
// exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetVPCAccessConnectorAttrs(t testing.TestingT, ctx context.Context, projectID string, region string, connectorID string) *vpcaccess.Connector {
	connector, err := GetVPCAccessConnectorAttrsE(t, ctx, projectID, region, connectorID)
	require.NoError(t, err)

	return connector
}

// GetVPCAccessConnectorAttrsE returns the settings Google Cloud holds for the given Serverless VPC
// Access connector.
// The ctx parameter supports cancellation and timeouts.
func GetVPCAccessConnectorAttrsE(t testing.TestingT, ctx context.Context, projectID string, region string, connectorID string) (*vpcaccess.Connector, error) {
	logger.Default.Logf(t, "Getting settings for Serverless VPC Access connector %s in region %s in project %s", connectorID, region, projectID)

	service, err := NewVPCAccessServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetVPCAccessConnectorAttrsWithClient(ctx, service, projectID, region, connectorID)
}

// GetVPCAccessConnectorAttrsWithClient returns the settings Google Cloud holds for the given
// Serverless VPC Access connector using the supplied *vpcaccess.Service. Prefer this variant in
// unit tests where the service is backed by an httptest fake server (see vpcaccess_test.go for the
// pattern).
// The ctx parameter supports cancellation and timeouts.
func GetVPCAccessConnectorAttrsWithClient(ctx context.Context, service *vpcaccess.Service, projectID string, region string, connectorID string) (*vpcaccess.Connector, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/connectors/%s", projectID, region, connectorID)

	connector, err := service.Projects.Locations.Connectors.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Serverless VPC Access connector %s does not exist in region %s in project %s", connectorID, region, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Serverless VPC Access connector %s in region %s in project %s: %w", connectorID, region, projectID, err)
	}

	return connector, nil
}

// NewVPCAccessServiceE creates a Serverless VPC Access service authenticated the same way every
// other client in this module is.
// The ctx parameter supports cancellation and timeouts.
func NewVPCAccessServiceE(t testing.TestingT, ctx context.Context) (*vpcaccess.Service, error) {
	return vpcaccess.NewService(ctx, append(withOptions(), option.WithScopes(vpcaccess.CloudPlatformScope))...)
}
