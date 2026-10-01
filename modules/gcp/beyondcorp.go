package gcp

import (
	"context"
	"errors"
	"fmt"

	"github.com/gruntwork-io/terratest/modules/core/v2/logger"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/beyondcorp/v1"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
)

// GetBeyondCorpAppConnectionAttrs returns the settings Google Cloud holds for the given BeyondCorp app connection, so a test can assert on what was
// actually created rather than only that it exists.
// A connection is what publishes one internal application, so the endpoint it points at and the connectors that serve it decide what is reachable.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetBeyondCorpAppConnectionAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, connectionID string) *beyondcorp.GoogleCloudBeyondcorpAppconnectionsV1AppConnection {
	attrs, err := GetBeyondCorpAppConnectionAttrsE(t, ctx, projectID, location, connectionID)
	require.NoError(t, err)

	return attrs
}

// GetBeyondCorpAppConnectionAttrsE returns the settings Google Cloud holds for the given BeyondCorp app connection.
// The ctx parameter supports cancellation and timeouts.
func GetBeyondCorpAppConnectionAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, connectionID string) (*beyondcorp.GoogleCloudBeyondcorpAppconnectionsV1AppConnection, error) {
	logger.Default.Logf(t, "Getting settings for BeyondCorp app connection %s in %s in project %s", connectionID, location, projectID)

	service, err := NewBeyondCorpServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetBeyondCorpAppConnectionAttrsWithClient(ctx, service, projectID, location, connectionID)
}

// GetBeyondCorpAppConnectionAttrsWithClient returns the settings Google Cloud holds for the given BeyondCorp app connection using the supplied
// *beyondcorp.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see beyondcorp_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetBeyondCorpAppConnectionAttrsWithClient(ctx context.Context, service *beyondcorp.Service, projectID string, location string, connectionID string) (*beyondcorp.GoogleCloudBeyondcorpAppconnectionsV1AppConnection, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/appConnections/%s", projectID, location, connectionID)

	attrs, err := service.Projects.Locations.AppConnections.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the BeyondCorp app connection %s in %s in project %s does not exist", connectionID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for BeyondCorp app connection %s in %s in project %s: %w", connectionID, location, projectID, err)
	}

	return attrs, nil
}

// GetBeyondCorpAppConnectorAttrs returns the settings Google Cloud holds for the given BeyondCorp app connector, so a test can assert on what was
// actually created rather than only that it exists.
// A connector is the agent inside the network that dials out, so the service account it runs as and its state decide whether anything connects.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetBeyondCorpAppConnectorAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, connectorID string) *beyondcorp.GoogleCloudBeyondcorpAppconnectorsV1AppConnector {
	attrs, err := GetBeyondCorpAppConnectorAttrsE(t, ctx, projectID, location, connectorID)
	require.NoError(t, err)

	return attrs
}

// GetBeyondCorpAppConnectorAttrsE returns the settings Google Cloud holds for the given BeyondCorp app connector.
// The ctx parameter supports cancellation and timeouts.
func GetBeyondCorpAppConnectorAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, connectorID string) (*beyondcorp.GoogleCloudBeyondcorpAppconnectorsV1AppConnector, error) {
	logger.Default.Logf(t, "Getting settings for BeyondCorp app connector %s in %s in project %s", connectorID, location, projectID)

	service, err := NewBeyondCorpServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetBeyondCorpAppConnectorAttrsWithClient(ctx, service, projectID, location, connectorID)
}

// GetBeyondCorpAppConnectorAttrsWithClient returns the settings Google Cloud holds for the given BeyondCorp app connector using the supplied
// *beyondcorp.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see beyondcorp_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetBeyondCorpAppConnectorAttrsWithClient(ctx context.Context, service *beyondcorp.Service, projectID string, location string, connectorID string) (*beyondcorp.GoogleCloudBeyondcorpAppconnectorsV1AppConnector, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/appConnectors/%s", projectID, location, connectorID)

	attrs, err := service.Projects.Locations.AppConnectors.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the BeyondCorp app connector %s in %s in project %s does not exist", connectorID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for BeyondCorp app connector %s in %s in project %s: %w", connectorID, location, projectID, err)
	}

	return attrs, nil
}

// GetBeyondCorpAppGatewayAttrs returns the settings Google Cloud holds for the given BeyondCorp app gateway, so a test can assert on what was
// actually created rather than only that it exists.
// A gateway is where clients arrive, so its type, its host type and the URI it answers on are the point.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetBeyondCorpAppGatewayAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, gatewayID string) *beyondcorp.AppGateway {
	attrs, err := GetBeyondCorpAppGatewayAttrsE(t, ctx, projectID, location, gatewayID)
	require.NoError(t, err)

	return attrs
}

// GetBeyondCorpAppGatewayAttrsE returns the settings Google Cloud holds for the given BeyondCorp app gateway.
// The ctx parameter supports cancellation and timeouts.
func GetBeyondCorpAppGatewayAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, gatewayID string) (*beyondcorp.AppGateway, error) {
	logger.Default.Logf(t, "Getting settings for BeyondCorp app gateway %s in %s in project %s", gatewayID, location, projectID)

	service, err := NewBeyondCorpServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetBeyondCorpAppGatewayAttrsWithClient(ctx, service, projectID, location, gatewayID)
}

// GetBeyondCorpAppGatewayAttrsWithClient returns the settings Google Cloud holds for the given BeyondCorp app gateway using the supplied
// *beyondcorp.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see beyondcorp_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetBeyondCorpAppGatewayAttrsWithClient(ctx context.Context, service *beyondcorp.Service, projectID string, location string, gatewayID string) (*beyondcorp.AppGateway, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/appGateways/%s", projectID, location, gatewayID)

	attrs, err := service.Projects.Locations.AppGateways.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the BeyondCorp app gateway %s in %s in project %s does not exist", gatewayID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for BeyondCorp app gateway %s in %s in project %s: %w", gatewayID, location, projectID, err)
	}

	return attrs, nil
}

// GetBeyondCorpSecurityGatewayAttrs returns the settings Google Cloud holds for the given BeyondCorp security gateway, so a test can assert on what was
// actually created rather than only that it exists.
// A security gateway is the newer way to publish applications, so the hubs it runs in decide where clients can arrive.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetBeyondCorpSecurityGatewayAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, gatewayID string) *beyondcorp.GoogleCloudBeyondcorpSecuritygatewaysV1SecurityGateway {
	attrs, err := GetBeyondCorpSecurityGatewayAttrsE(t, ctx, projectID, location, gatewayID)
	require.NoError(t, err)

	return attrs
}

// GetBeyondCorpSecurityGatewayAttrsE returns the settings Google Cloud holds for the given BeyondCorp security gateway.
// The ctx parameter supports cancellation and timeouts.
func GetBeyondCorpSecurityGatewayAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, gatewayID string) (*beyondcorp.GoogleCloudBeyondcorpSecuritygatewaysV1SecurityGateway, error) {
	logger.Default.Logf(t, "Getting settings for BeyondCorp security gateway %s in %s in project %s", gatewayID, location, projectID)

	service, err := NewBeyondCorpServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetBeyondCorpSecurityGatewayAttrsWithClient(ctx, service, projectID, location, gatewayID)
}

// GetBeyondCorpSecurityGatewayAttrsWithClient returns the settings Google Cloud holds for the given BeyondCorp security gateway using the supplied
// *beyondcorp.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see beyondcorp_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetBeyondCorpSecurityGatewayAttrsWithClient(ctx context.Context, service *beyondcorp.Service, projectID string, location string, gatewayID string) (*beyondcorp.GoogleCloudBeyondcorpSecuritygatewaysV1SecurityGateway, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/securityGateways/%s", projectID, location, gatewayID)

	attrs, err := service.Projects.Locations.SecurityGateways.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the BeyondCorp security gateway %s in %s in project %s does not exist", gatewayID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for BeyondCorp security gateway %s in %s in project %s: %w", gatewayID, location, projectID, err)
	}

	return attrs, nil
}

// GetBeyondCorpSecurityGatewayApplicationAttrs returns the settings Google Cloud holds for the given BeyondCorp security gateway application, so a test can assert on what was
// actually created rather than only that it exists.
// An application is one thing published through a gateway, so the endpoint matchers it carries decide which hostnames and ports reach it.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetBeyondCorpSecurityGatewayApplicationAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, gatewayID string, applicationID string) *beyondcorp.GoogleCloudBeyondcorpSecuritygatewaysV1Application {
	attrs, err := GetBeyondCorpSecurityGatewayApplicationAttrsE(t, ctx, projectID, location, gatewayID, applicationID)
	require.NoError(t, err)

	return attrs
}

// GetBeyondCorpSecurityGatewayApplicationAttrsE returns the settings Google Cloud holds for the given BeyondCorp security gateway application.
// The ctx parameter supports cancellation and timeouts.
func GetBeyondCorpSecurityGatewayApplicationAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, gatewayID string, applicationID string) (*beyondcorp.GoogleCloudBeyondcorpSecuritygatewaysV1Application, error) {
	logger.Default.Logf(t, "Getting settings for BeyondCorp security gateway application %s in gateway %s in %s in project %s", applicationID, gatewayID, location, projectID)

	service, err := NewBeyondCorpServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetBeyondCorpSecurityGatewayApplicationAttrsWithClient(ctx, service, projectID, location, gatewayID, applicationID)
}

// GetBeyondCorpSecurityGatewayApplicationAttrsWithClient returns the settings Google Cloud holds for the given BeyondCorp security gateway application using the supplied
// *beyondcorp.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see beyondcorp_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetBeyondCorpSecurityGatewayApplicationAttrsWithClient(ctx context.Context, service *beyondcorp.Service, projectID string, location string, gatewayID string, applicationID string) (*beyondcorp.GoogleCloudBeyondcorpSecuritygatewaysV1Application, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/securityGateways/%s/applications/%s", projectID, location, gatewayID, applicationID)

	attrs, err := service.Projects.Locations.SecurityGateways.Applications.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the BeyondCorp security gateway application %s in gateway %s in %s in project %s does not exist", applicationID, gatewayID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for BeyondCorp security gateway application %s in gateway %s in %s in project %s: %w", applicationID, gatewayID, location, projectID, err)
	}

	return attrs, nil
}

// GetBeyondCorpSecurityGatewayIamPolicyAttrs returns the IAM policy Google Cloud holds for the given BeyondCorp security gateway, so a test can assert on what was
// actually created rather than only that it exists.
// Who may use the gateway is what the policy decides, and a gateway Google created on its own carries no binding at all.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetBeyondCorpSecurityGatewayIamPolicyAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, gatewayID string) *beyondcorp.GoogleIamV1Policy {
	policy, err := GetBeyondCorpSecurityGatewayIamPolicyAttrsE(t, ctx, projectID, location, gatewayID)
	require.NoError(t, err)

	return policy
}

// GetBeyondCorpSecurityGatewayIamPolicyAttrsE returns the IAM policy Google Cloud holds for the given BeyondCorp security gateway.
// The ctx parameter supports cancellation and timeouts.
func GetBeyondCorpSecurityGatewayIamPolicyAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, gatewayID string) (*beyondcorp.GoogleIamV1Policy, error) {
	logger.Default.Logf(t, "Getting settings for BeyondCorp security gateway %s in %s in project %s", gatewayID, location, projectID)

	service, err := NewBeyondCorpServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetBeyondCorpSecurityGatewayIamPolicyAttrsWithClient(ctx, service, projectID, location, gatewayID)
}

// GetBeyondCorpSecurityGatewayIamPolicyAttrsWithClient returns the IAM policy Google Cloud holds for the given BeyondCorp security gateway using the supplied
// *beyondcorp.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see beyondcorp_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetBeyondCorpSecurityGatewayIamPolicyAttrsWithClient(ctx context.Context, service *beyondcorp.Service, projectID string, location string, gatewayID string) (*beyondcorp.GoogleIamV1Policy, error) {
	resource := fmt.Sprintf("projects/%s/locations/%s/securityGateways/%s", projectID, location, gatewayID)

	policy, err := service.Projects.Locations.SecurityGateways.GetIamPolicy(resource).OptionsRequestedPolicyVersion(iamPolicyVersionWithConditions).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the BeyondCorp security gateway %s in %s in project %s does not exist", gatewayID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for BeyondCorp security gateway %s in %s in project %s: %w", gatewayID, location, projectID, err)
	}

	return policy, nil
}

// GetBeyondCorpSecurityGatewayApplicationIamPolicyAttrs returns the IAM policy Google Cloud holds for the given BeyondCorp security gateway application, so a test can assert on what was
// actually created rather than only that it exists.
// Who may reach this one application is what the policy decides, which is finer than the policy on the gateway itself.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetBeyondCorpSecurityGatewayApplicationIamPolicyAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, gatewayID string, applicationID string) *beyondcorp.GoogleIamV1Policy {
	policy, err := GetBeyondCorpSecurityGatewayApplicationIamPolicyAttrsE(t, ctx, projectID, location, gatewayID, applicationID)
	require.NoError(t, err)

	return policy
}

// GetBeyondCorpSecurityGatewayApplicationIamPolicyAttrsE returns the IAM policy Google Cloud holds for the given BeyondCorp security gateway application.
// The ctx parameter supports cancellation and timeouts.
func GetBeyondCorpSecurityGatewayApplicationIamPolicyAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, gatewayID string, applicationID string) (*beyondcorp.GoogleIamV1Policy, error) {
	logger.Default.Logf(t, "Getting settings for BeyondCorp security gateway application %s in gateway %s in %s in project %s", applicationID, gatewayID, location, projectID)

	service, err := NewBeyondCorpServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetBeyondCorpSecurityGatewayApplicationIamPolicyAttrsWithClient(ctx, service, projectID, location, gatewayID, applicationID)
}

// GetBeyondCorpSecurityGatewayApplicationIamPolicyAttrsWithClient returns the IAM policy Google Cloud holds for the given BeyondCorp security gateway application using the supplied
// *beyondcorp.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see beyondcorp_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetBeyondCorpSecurityGatewayApplicationIamPolicyAttrsWithClient(ctx context.Context, service *beyondcorp.Service, projectID string, location string, gatewayID string, applicationID string) (*beyondcorp.GoogleIamV1Policy, error) {
	resource := fmt.Sprintf("projects/%s/locations/%s/securityGateways/%s/applications/%s", projectID, location, gatewayID, applicationID)

	policy, err := service.Projects.Locations.SecurityGateways.Applications.GetIamPolicy(resource).OptionsRequestedPolicyVersion(iamPolicyVersionWithConditions).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the BeyondCorp security gateway application %s in gateway %s in %s in project %s does not exist", applicationID, gatewayID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for BeyondCorp security gateway application %s in gateway %s in %s in project %s: %w", applicationID, gatewayID, location, projectID, err)
	}

	return policy, nil
}

// NewBeyondCorpServiceE creates a BeyondCorp service authenticated the same way every other client in this
// module is.
// The ctx parameter supports cancellation and timeouts.
func NewBeyondCorpServiceE(t testing.TestingT, ctx context.Context) (*beyondcorp.Service, error) {
	return beyondcorp.NewService(ctx, append(withOptions(), option.WithScopes(beyondcorp.CloudPlatformScope))...)
}
