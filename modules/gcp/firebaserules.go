package gcp

import (
	"context"
	"errors"
	"fmt"

	"github.com/gruntwork-io/terratest/modules/core/v2/logger"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/firebaserules/v1"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
)

// GetFirebaseRulesetAttrs returns the settings Google Cloud holds for the given Firebase Security
// Rules ruleset, so a test can assert on what was actually created rather than only that it
// exists. A ruleset has no caller-chosen name, so the id is the one Google assigned.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetFirebaseRulesetAttrs(t testing.TestingT, ctx context.Context, projectID string, rulesetID string) *firebaserules.Ruleset {
	ruleset, err := GetFirebaseRulesetAttrsE(t, ctx, projectID, rulesetID)
	require.NoError(t, err)

	return ruleset
}

// GetFirebaseRulesetAttrsE returns the settings Google Cloud holds for the given Firebase Security
// Rules ruleset.
// The ctx parameter supports cancellation and timeouts.
func GetFirebaseRulesetAttrsE(t testing.TestingT, ctx context.Context, projectID string, rulesetID string) (*firebaserules.Ruleset, error) {
	logger.Default.Logf(t, "Getting settings for Firebase ruleset %s in project %s", rulesetID, projectID)

	service, err := NewFirebaseRulesServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetFirebaseRulesetAttrsWithClient(ctx, service, projectID, rulesetID)
}

// GetFirebaseRulesetAttrsWithClient returns the settings Google Cloud holds for the given Firebase
// Security Rules ruleset using the supplied *firebaserules.Service. Prefer this variant in unit
// tests where the service is backed by an httptest fake server (see firebaserules_test.go for the
// pattern).
// The ctx parameter supports cancellation and timeouts.
func GetFirebaseRulesetAttrsWithClient(ctx context.Context, service *firebaserules.Service, projectID string, rulesetID string) (*firebaserules.Ruleset, error) {
	name := fmt.Sprintf("projects/%s/rulesets/%s", projectID, rulesetID)

	ruleset, err := service.Projects.Rulesets.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("Firebase ruleset %s does not exist in project %s", rulesetID, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Firebase ruleset %s in project %s: %w", rulesetID, projectID, err)
	}

	return ruleset, nil
}

// NewFirebaseRulesServiceE creates a Firebase Rules service authenticated the same way every other
// client in this module is.
// The ctx parameter supports cancellation and timeouts.
func NewFirebaseRulesServiceE(t testing.TestingT, ctx context.Context) (*firebaserules.Service, error) {
	return firebaserules.NewService(ctx, append(withOptions(), option.WithScopes(firebaserules.CloudPlatformScope))...)
}
