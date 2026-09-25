package gcp

import (
	"context"
	"errors"
	"fmt"

	"github.com/gruntwork-io/terratest/modules/core/v2/logger"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/contactcenterinsights/v1"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
)

// GetInsightsViewAttrs returns the settings Google Cloud holds for the given Contact Center
// Insights view, so a test can assert on what was actually created rather than only that it
// exists. A view lives in a location, which has to be given, and the id is the one Google assigns.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetInsightsViewAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, viewID string) *contactcenterinsights.GoogleCloudContactcenterinsightsV1View {
	view, err := GetInsightsViewAttrsE(t, ctx, projectID, location, viewID)
	require.NoError(t, err)

	return view
}

// GetInsightsViewAttrsE returns the settings Google Cloud holds for the given Contact Center
// Insights view.
// The ctx parameter supports cancellation and timeouts.
func GetInsightsViewAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, viewID string) (*contactcenterinsights.GoogleCloudContactcenterinsightsV1View, error) {
	logger.Default.Logf(t, "Getting settings for Contact Center Insights view %s in project %s", viewID, projectID)

	service, err := NewContactCenterInsightsServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetInsightsViewAttrsWithClient(ctx, service, projectID, location, viewID)
}

// GetInsightsViewAttrsWithClient returns the settings Google Cloud holds for the given Contact
// Center Insights view using the supplied *contactcenterinsights.Service. Prefer this variant in
// unit tests where the service is backed by an httptest fake server (see
// contactcenterinsights_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetInsightsViewAttrsWithClient(ctx context.Context, service *contactcenterinsights.Service, projectID string, location string, viewID string) (*contactcenterinsights.GoogleCloudContactcenterinsightsV1View, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/views/%s", projectID, location, viewID)

	view, err := service.Projects.Locations.Views.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("Contact Center Insights view %s does not exist in project %s", viewID, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Contact Center Insights view %s in project %s: %w", viewID, projectID, err)
	}

	return view, nil
}

// NewContactCenterInsightsServiceE creates a Contact Center Insights service authenticated the
// same way every other client in this module is. Every location answers on the one host, so
// unlike Dialogflow the location goes in the resource name rather than the endpoint.
// The ctx parameter supports cancellation and timeouts.
func NewContactCenterInsightsServiceE(t testing.TestingT, ctx context.Context) (*contactcenterinsights.Service, error) {
	return contactcenterinsights.NewService(ctx, append(withOptions(), option.WithScopes(contactcenterinsights.CloudPlatformScope))...)
}

// GetContactCenterInsightsAnalysisRuleAttrs returns the settings Google Cloud holds for the given analysis rule, so a test can assert on
// what was actually created rather than only that it exists. A rule decides which conversations get analysed and how
// much of them.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetContactCenterInsightsAnalysisRuleAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, ruleID string) *contactcenterinsights.GoogleCloudContactcenterinsightsV1AnalysisRule {
	rule, err := GetContactCenterInsightsAnalysisRuleAttrsE(t, ctx, projectID, location, ruleID)
	require.NoError(t, err)

	return rule
}

// GetContactCenterInsightsAnalysisRuleAttrsE returns the settings Google Cloud holds for the given analysis rule.
// The ctx parameter supports cancellation and timeouts.
func GetContactCenterInsightsAnalysisRuleAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, ruleID string) (*contactcenterinsights.GoogleCloudContactcenterinsightsV1AnalysisRule, error) {
	logger.Default.Logf(t, "Getting settings for analysis rule %s in %s in project %s", ruleID, location, projectID)

	service, err := NewContactCenterInsightsServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetContactCenterInsightsAnalysisRuleAttrsWithClient(ctx, service, projectID, location, ruleID)
}

// GetContactCenterInsightsAnalysisRuleAttrsWithClient returns the settings Google Cloud holds for the given analysis rule using the supplied
// *contactcenterinsights.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see contactcenterinsights_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetContactCenterInsightsAnalysisRuleAttrsWithClient(ctx context.Context, service *contactcenterinsights.Service, projectID string, location string, ruleID string) (*contactcenterinsights.GoogleCloudContactcenterinsightsV1AnalysisRule, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/analysisRules/%s", projectID, location, ruleID)

	rule, err := service.Projects.Locations.AnalysisRules.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the analysis rule %s does not exist in %s in project %s", ruleID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for analysis rule %s in %s in project %s: %w", ruleID, location, projectID, err)
	}

	return rule, nil
}
