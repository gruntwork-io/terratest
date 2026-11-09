package gcp

import (
	"context"
	"errors"
	"fmt"

	"github.com/gruntwork-io/terratest/modules/core/v2/logger"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/appengine/v1"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
)

// GetAppEngineDomainMappingAttrs returns the settings Google Cloud holds for the App Engine domain mapping, so a test can assert on
// what was actually created rather than only that it exists.
// The records Google wants the domain's owner to publish are reported on the mapping, and a create
// does not list them.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetAppEngineDomainMappingAttrs(t testing.TestingT, ctx context.Context, projectID string, domain string) *appengine.DomainMapping {
	result, err := GetAppEngineDomainMappingAttrsE(t, ctx, projectID, domain)
	require.NoError(t, err)

	return result
}

// GetAppEngineDomainMappingAttrsE returns the settings Google Cloud holds for the App Engine domain mapping.
// The ctx parameter supports cancellation and timeouts.
func GetAppEngineDomainMappingAttrsE(t testing.TestingT, ctx context.Context, projectID string, domain string) (*appengine.DomainMapping, error) {
	logger.Default.Logf(t, "Getting settings for App Engine domain mapping %s in project %s", domain, projectID)

	service, err := NewAppEngineServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetAppEngineDomainMappingAttrsWithClient(ctx, service, projectID, domain)
}

// GetAppEngineDomainMappingAttrsWithClient returns the settings Google Cloud holds for the App Engine domain mapping using the supplied
// *appengine.APIService. Prefer this variant in unit tests where the service is backed by an httptest fake server
// (see appengine_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetAppEngineDomainMappingAttrsWithClient(ctx context.Context, service *appengine.APIService, projectID string, domain string) (*appengine.DomainMapping, error) {
	result, err := service.Apps.DomainMappings.Get(projectID, domain).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the App Engine domain mapping %s does not exist in project %s", domain, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for App Engine domain mapping %s in project %s: %w", domain, projectID, err)
	}

	return result, nil
}

// GetAppEngineFirewallRuleAttrs returns the settings Google Cloud holds for the App Engine firewall rule, so a test can assert on
// what was actually created rather than only that it exists.
// A rule is identified by its priority rather than by a name, so the priority is passed as the id
// the API asks for: a decimal string such as 1000.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetAppEngineFirewallRuleAttrs(t testing.TestingT, ctx context.Context, projectID string, priority string) *appengine.FirewallRule {
	result, err := GetAppEngineFirewallRuleAttrsE(t, ctx, projectID, priority)
	require.NoError(t, err)

	return result
}

// GetAppEngineFirewallRuleAttrsE returns the settings Google Cloud holds for the App Engine firewall rule.
// The ctx parameter supports cancellation and timeouts.
func GetAppEngineFirewallRuleAttrsE(t testing.TestingT, ctx context.Context, projectID string, priority string) (*appengine.FirewallRule, error) {
	logger.Default.Logf(t, "Getting settings for App Engine firewall rule %s in project %s", priority, projectID)

	service, err := NewAppEngineServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetAppEngineFirewallRuleAttrsWithClient(ctx, service, projectID, priority)
}

// GetAppEngineFirewallRuleAttrsWithClient returns the settings Google Cloud holds for the App Engine firewall rule using the supplied
// *appengine.APIService. Prefer this variant in unit tests where the service is backed by an httptest fake server
// (see appengine_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetAppEngineFirewallRuleAttrsWithClient(ctx context.Context, service *appengine.APIService, projectID string, priority string) (*appengine.FirewallRule, error) {
	result, err := service.Apps.Firewall.IngressRules.Get(projectID, priority).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the App Engine firewall rule %s does not exist in project %s", priority, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for App Engine firewall rule %s in project %s: %w", priority, projectID, err)
	}

	return result, nil
}

// GetAppEngineServiceAttrs returns the settings Google Cloud holds for the App Engine service, so a test can assert on
// what was actually created rather than only that it exists.
// One read answers for two modules: the traffic split across a service's versions and the service's
// network settings are both fields of the service rather than resources of their own.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetAppEngineServiceAttrs(t testing.TestingT, ctx context.Context, projectID string, serviceID string) *appengine.Service {
	result, err := GetAppEngineServiceAttrsE(t, ctx, projectID, serviceID)
	require.NoError(t, err)

	return result
}

// GetAppEngineServiceAttrsE returns the settings Google Cloud holds for the App Engine service.
// The ctx parameter supports cancellation and timeouts.
func GetAppEngineServiceAttrsE(t testing.TestingT, ctx context.Context, projectID string, serviceID string) (*appengine.Service, error) {
	logger.Default.Logf(t, "Getting settings for App Engine service %s in project %s", serviceID, projectID)

	service, err := NewAppEngineServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetAppEngineServiceAttrsWithClient(ctx, service, projectID, serviceID)
}

// GetAppEngineServiceAttrsWithClient returns the settings Google Cloud holds for the App Engine service using the supplied
// *appengine.APIService. Prefer this variant in unit tests where the service is backed by an httptest fake server
// (see appengine_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetAppEngineServiceAttrsWithClient(ctx context.Context, service *appengine.APIService, projectID string, serviceID string) (*appengine.Service, error) {
	result, err := service.Apps.Services.Get(projectID, serviceID).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the App Engine service %s does not exist in project %s", serviceID, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for App Engine service %s in project %s: %w", serviceID, projectID, err)
	}

	return result, nil
}

// GetAppEngineVersionAttrs returns the settings Google Cloud holds for the App Engine version, so a test can assert on
// what was actually created rather than only that it exists.
// The full view is asked for, because the basic view leaves out the handlers, the entrypoint and the
// scaling settings, which is most of what a version module sets.
// One read answers for both the standard and the flexible version modules.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetAppEngineVersionAttrs(t testing.TestingT, ctx context.Context, projectID string, serviceID string, versionID string) *appengine.Version {
	result, err := GetAppEngineVersionAttrsE(t, ctx, projectID, serviceID, versionID)
	require.NoError(t, err)

	return result
}

// GetAppEngineVersionAttrsE returns the settings Google Cloud holds for the App Engine version.
// The ctx parameter supports cancellation and timeouts.
func GetAppEngineVersionAttrsE(t testing.TestingT, ctx context.Context, projectID string, serviceID string, versionID string) (*appengine.Version, error) {
	logger.Default.Logf(t, "Getting settings for App Engine version %s of service %s in project %s", versionID, serviceID, projectID)

	service, err := NewAppEngineServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetAppEngineVersionAttrsWithClient(ctx, service, projectID, serviceID, versionID)
}

// GetAppEngineVersionAttrsWithClient returns the settings Google Cloud holds for the App Engine version using the supplied
// *appengine.APIService. Prefer this variant in unit tests where the service is backed by an httptest fake server
// (see appengine_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetAppEngineVersionAttrsWithClient(ctx context.Context, service *appengine.APIService, projectID string, serviceID string, versionID string) (*appengine.Version, error) {
	result, err := service.Apps.Services.Versions.Get(projectID, serviceID, versionID).View("FULL").Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the App Engine version %s does not exist on service %s in project %s", versionID, serviceID, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for App Engine version %s of service %s in project %s: %w", versionID, serviceID, projectID, err)
	}

	return result, nil
}

// NewAppEngineServiceE creates an App Engine Admin service authenticated the same way every other
// client in this module is.
// Every id this API takes is positional rather than part of a resource name, and the application's
// id is the project id.
// The ctx parameter supports cancellation and timeouts.
func NewAppEngineServiceE(t testing.TestingT, ctx context.Context) (*appengine.APIService, error) {
	return appengine.NewService(ctx, append(withOptions(), option.WithScopes(appengine.CloudPlatformScope))...)
}
