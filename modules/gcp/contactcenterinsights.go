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

// GetContactCenterInsightsAssessmentRuleAttrs returns the settings Google Cloud holds for the given
// assessment rule, so a test can assert on which conversations it samples rather than only that it
// exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetContactCenterInsightsAssessmentRuleAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, ruleID string) *contactcenterinsights.GoogleCloudContactcenterinsightsV1AssessmentRule {
	rule, err := GetContactCenterInsightsAssessmentRuleAttrsE(t, ctx, projectID, location, ruleID)
	require.NoError(t, err)

	return rule
}

// GetContactCenterInsightsAssessmentRuleAttrsE returns the settings Google Cloud holds for the given
// assessment rule.
// The ctx parameter supports cancellation and timeouts.
func GetContactCenterInsightsAssessmentRuleAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, ruleID string) (*contactcenterinsights.GoogleCloudContactcenterinsightsV1AssessmentRule, error) {
	logger.Default.Logf(t, "Getting settings for assessment rule %s in %s in project %s", ruleID, location, projectID)

	service, err := NewContactCenterInsightsServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetContactCenterInsightsAssessmentRuleAttrsWithClient(ctx, service, projectID, location, ruleID)
}

// GetContactCenterInsightsAssessmentRuleAttrsWithClient returns the settings Google Cloud holds for
// the given assessment rule using the supplied *contactcenterinsights.Service. Prefer this variant
// in unit tests where the service is backed by an httptest fake server (see
// contactcenterinsights_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetContactCenterInsightsAssessmentRuleAttrsWithClient(ctx context.Context, service *contactcenterinsights.Service, projectID string, location string, ruleID string) (*contactcenterinsights.GoogleCloudContactcenterinsightsV1AssessmentRule, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/assessmentRules/%s", projectID, location, ruleID)

	rule, err := service.Projects.Locations.AssessmentRules.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the assessment rule %s in %s in project %s does not exist", ruleID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for assessment rule %s in %s in project %s: %w", ruleID, location, projectID, err)
	}

	return rule, nil
}

// GetContactCenterInsightsAutoLabelingRuleAttrs returns the settings Google Cloud holds for the given
// auto labeling rule, so a test can assert on which label it applies and to what.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetContactCenterInsightsAutoLabelingRuleAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, ruleID string) *contactcenterinsights.GoogleCloudContactcenterinsightsV1AutoLabelingRule {
	rule, err := GetContactCenterInsightsAutoLabelingRuleAttrsE(t, ctx, projectID, location, ruleID)
	require.NoError(t, err)

	return rule
}

// GetContactCenterInsightsAutoLabelingRuleAttrsE returns the settings Google Cloud holds for the
// given auto labeling rule.
// The ctx parameter supports cancellation and timeouts.
func GetContactCenterInsightsAutoLabelingRuleAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, ruleID string) (*contactcenterinsights.GoogleCloudContactcenterinsightsV1AutoLabelingRule, error) {
	logger.Default.Logf(t, "Getting settings for auto labeling rule %s in %s in project %s", ruleID, location, projectID)

	service, err := NewContactCenterInsightsServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetContactCenterInsightsAutoLabelingRuleAttrsWithClient(ctx, service, projectID, location, ruleID)
}

// GetContactCenterInsightsAutoLabelingRuleAttrsWithClient returns the settings Google Cloud holds for
// the given auto labeling rule using the supplied *contactcenterinsights.Service. Prefer this variant
// in unit tests where the service is backed by an httptest fake server (see
// contactcenterinsights_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetContactCenterInsightsAutoLabelingRuleAttrsWithClient(ctx context.Context, service *contactcenterinsights.Service, projectID string, location string, ruleID string) (*contactcenterinsights.GoogleCloudContactcenterinsightsV1AutoLabelingRule, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/autoLabelingRules/%s", projectID, location, ruleID)

	rule, err := service.Projects.Locations.AutoLabelingRules.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the auto labeling rule %s in %s in project %s does not exist", ruleID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for auto labeling rule %s in %s in project %s: %w", ruleID, location, projectID, err)
	}

	return rule, nil
}

// GetContactCenterInsightsQaScorecardAttrs returns the settings Google Cloud holds for the given QA
// scorecard, so a test can assert on what was actually created rather than only that it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetContactCenterInsightsQaScorecardAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, scorecardID string) *contactcenterinsights.GoogleCloudContactcenterinsightsV1QaScorecard {
	scorecard, err := GetContactCenterInsightsQaScorecardAttrsE(t, ctx, projectID, location, scorecardID)
	require.NoError(t, err)

	return scorecard
}

// GetContactCenterInsightsQaScorecardAttrsE returns the settings Google Cloud holds for the given QA
// scorecard.
// The ctx parameter supports cancellation and timeouts.
func GetContactCenterInsightsQaScorecardAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, scorecardID string) (*contactcenterinsights.GoogleCloudContactcenterinsightsV1QaScorecard, error) {
	logger.Default.Logf(t, "Getting settings for QA scorecard %s in %s in project %s", scorecardID, location, projectID)

	service, err := NewContactCenterInsightsServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetContactCenterInsightsQaScorecardAttrsWithClient(ctx, service, projectID, location, scorecardID)
}

// GetContactCenterInsightsQaScorecardAttrsWithClient returns the settings Google Cloud holds for the
// given QA scorecard using the supplied *contactcenterinsights.Service. Prefer this variant in unit
// tests where the service is backed by an httptest fake server (see contactcenterinsights_test.go
// for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetContactCenterInsightsQaScorecardAttrsWithClient(ctx context.Context, service *contactcenterinsights.Service, projectID string, location string, scorecardID string) (*contactcenterinsights.GoogleCloudContactcenterinsightsV1QaScorecard, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/qaScorecards/%s", projectID, location, scorecardID)

	scorecard, err := service.Projects.Locations.QaScorecards.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the QA scorecard %s in %s in project %s does not exist", scorecardID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for QA scorecard %s in %s in project %s: %w", scorecardID, location, projectID, err)
	}

	return scorecard, nil
}
