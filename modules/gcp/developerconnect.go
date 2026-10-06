package gcp

import (
	"context"
	"errors"
	"fmt"

	"github.com/gruntwork-io/terratest/modules/core/v2/logger"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/developerconnect/v1"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
)

// GetDeveloperConnectConnectionAttrs returns the settings Google Cloud holds for the given Developer Connect connection, so a test can assert on what was
// actually created rather than only that it exists.
// A connection is the link to a git host, so which host it names and whether it is installed decide whether anything can be linked.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetDeveloperConnectConnectionAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, connectionID string) *developerconnect.Connection {
	attrs, err := GetDeveloperConnectConnectionAttrsE(t, ctx, projectID, location, connectionID)
	require.NoError(t, err)

	return attrs
}

// GetDeveloperConnectConnectionAttrsE returns the settings Google Cloud holds for the given Developer Connect connection.
// The ctx parameter supports cancellation and timeouts.
func GetDeveloperConnectConnectionAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, connectionID string) (*developerconnect.Connection, error) {
	logger.Default.Logf(t, "Getting settings for Developer Connect connection %s in %s in project %s", connectionID, location, projectID)

	service, err := NewDeveloperConnectServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetDeveloperConnectConnectionAttrsWithClient(ctx, service, projectID, location, connectionID)
}

// GetDeveloperConnectConnectionAttrsWithClient returns the settings Google Cloud holds for the given Developer Connect connection using the supplied
// *developerconnect.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see developerconnect_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetDeveloperConnectConnectionAttrsWithClient(ctx context.Context, service *developerconnect.Service, projectID string, location string, connectionID string) (*developerconnect.Connection, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/connections/%s", projectID, location, connectionID)

	attrs, err := service.Projects.Locations.Connections.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Developer Connect connection %s in %s in project %s does not exist", connectionID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Developer Connect connection %s in %s in project %s: %w", connectionID, location, projectID, err)
	}

	return attrs, nil
}

// GetDeveloperConnectGitRepositoryLinkAttrs returns the settings Google Cloud holds for the given Developer Connect git repository link, so a test can assert on what was
// actually created rather than only that it exists.
// A link is one repository made available through a connection, so the clone URI it names is the whole point.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetDeveloperConnectGitRepositoryLinkAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, connectionID string, linkID string) *developerconnect.GitRepositoryLink {
	attrs, err := GetDeveloperConnectGitRepositoryLinkAttrsE(t, ctx, projectID, location, connectionID, linkID)
	require.NoError(t, err)

	return attrs
}

// GetDeveloperConnectGitRepositoryLinkAttrsE returns the settings Google Cloud holds for the given Developer Connect git repository link.
// The ctx parameter supports cancellation and timeouts.
func GetDeveloperConnectGitRepositoryLinkAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, connectionID string, linkID string) (*developerconnect.GitRepositoryLink, error) {
	logger.Default.Logf(t, "Getting settings for Developer Connect git repository link %s %s in %s in project %s", linkID, connectionID, location, projectID)

	service, err := NewDeveloperConnectServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetDeveloperConnectGitRepositoryLinkAttrsWithClient(ctx, service, projectID, location, connectionID, linkID)
}

// GetDeveloperConnectGitRepositoryLinkAttrsWithClient returns the settings Google Cloud holds for the given Developer Connect git repository link using the supplied
// *developerconnect.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see developerconnect_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetDeveloperConnectGitRepositoryLinkAttrsWithClient(ctx context.Context, service *developerconnect.Service, projectID string, location string, connectionID string, linkID string) (*developerconnect.GitRepositoryLink, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/connections/%s/gitRepositoryLinks/%s", projectID, location, connectionID, linkID)

	attrs, err := service.Projects.Locations.Connections.GitRepositoryLinks.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Developer Connect git repository link %s %s in %s in project %s does not exist", linkID, connectionID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Developer Connect git repository link %s %s in %s in project %s: %w", linkID, connectionID, location, projectID, err)
	}

	return attrs, nil
}

// GetDeveloperConnectAccountConnectorAttrs returns the settings Google Cloud holds for the given Developer Connect account connector, so a test can assert on what was
// actually created rather than only that it exists.
// An account connector holds a user's own credentials for a third party, so the provider it names decides what it can act on.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetDeveloperConnectAccountConnectorAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, connectorID string) *developerconnect.AccountConnector {
	attrs, err := GetDeveloperConnectAccountConnectorAttrsE(t, ctx, projectID, location, connectorID)
	require.NoError(t, err)

	return attrs
}

// GetDeveloperConnectAccountConnectorAttrsE returns the settings Google Cloud holds for the given Developer Connect account connector.
// The ctx parameter supports cancellation and timeouts.
func GetDeveloperConnectAccountConnectorAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, connectorID string) (*developerconnect.AccountConnector, error) {
	logger.Default.Logf(t, "Getting settings for Developer Connect account connector %s in %s in project %s", connectorID, location, projectID)

	service, err := NewDeveloperConnectServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetDeveloperConnectAccountConnectorAttrsWithClient(ctx, service, projectID, location, connectorID)
}

// GetDeveloperConnectAccountConnectorAttrsWithClient returns the settings Google Cloud holds for the given Developer Connect account connector using the supplied
// *developerconnect.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see developerconnect_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetDeveloperConnectAccountConnectorAttrsWithClient(ctx context.Context, service *developerconnect.Service, projectID string, location string, connectorID string) (*developerconnect.AccountConnector, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/accountConnectors/%s", projectID, location, connectorID)

	attrs, err := service.Projects.Locations.AccountConnectors.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Developer Connect account connector %s in %s in project %s does not exist", connectorID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Developer Connect account connector %s in %s in project %s: %w", connectorID, location, projectID, err)
	}

	return attrs, nil
}

// GetDeveloperConnectInsightsConfigAttrs returns the settings Google Cloud holds for the given Developer Connect insights config, so a test can assert on what was
// actually created rather than only that it exists.
// An insights config ties a deployed app back to the code that built it, so the app it names is the point.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetDeveloperConnectInsightsConfigAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, configID string) *developerconnect.InsightsConfig {
	attrs, err := GetDeveloperConnectInsightsConfigAttrsE(t, ctx, projectID, location, configID)
	require.NoError(t, err)

	return attrs
}

// GetDeveloperConnectInsightsConfigAttrsE returns the settings Google Cloud holds for the given Developer Connect insights config.
// The ctx parameter supports cancellation and timeouts.
func GetDeveloperConnectInsightsConfigAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, configID string) (*developerconnect.InsightsConfig, error) {
	logger.Default.Logf(t, "Getting settings for Developer Connect insights config %s in %s in project %s", configID, location, projectID)

	service, err := NewDeveloperConnectServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetDeveloperConnectInsightsConfigAttrsWithClient(ctx, service, projectID, location, configID)
}

// GetDeveloperConnectInsightsConfigAttrsWithClient returns the settings Google Cloud holds for the given Developer Connect insights config using the supplied
// *developerconnect.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see developerconnect_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetDeveloperConnectInsightsConfigAttrsWithClient(ctx context.Context, service *developerconnect.Service, projectID string, location string, configID string) (*developerconnect.InsightsConfig, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/insightsConfigs/%s", projectID, location, configID)

	attrs, err := service.Projects.Locations.InsightsConfigs.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Developer Connect insights config %s in %s in project %s does not exist", configID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Developer Connect insights config %s in %s in project %s: %w", configID, location, projectID, err)
	}

	return attrs, nil
}

// NewDeveloperConnectServiceE creates a Developer Connect service authenticated the same way every other client in this
// module is.
// The ctx parameter supports cancellation and timeouts.
func NewDeveloperConnectServiceE(t testing.TestingT, ctx context.Context) (*developerconnect.Service, error) {
	return developerconnect.NewService(ctx, append(withOptions(), option.WithScopes(developerconnect.CloudPlatformScope))...)
}
