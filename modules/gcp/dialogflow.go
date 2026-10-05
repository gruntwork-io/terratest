package gcp

import (
	"context"
	"errors"
	"fmt"
	"regexp"

	"github.com/gruntwork-io/terratest/modules/core/v2/logger"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/dialogflow/v3"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
)

// GetDialogflowAgentAttrs returns the settings Google Cloud holds for the given Dialogflow CX
// agent, so a test can assert on what was actually created rather than only that it exists. An
// agent lives in a location, which has to be given, and the id is the one Google assigns.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetDialogflowAgentAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, agentID string) *dialogflow.GoogleCloudDialogflowCxV3Agent {
	agent, err := GetDialogflowAgentAttrsE(t, ctx, projectID, location, agentID)
	require.NoError(t, err)

	return agent
}

// GetDialogflowAgentAttrsE returns the settings Google Cloud holds for the given Dialogflow CX
// agent.
// The ctx parameter supports cancellation and timeouts.
func GetDialogflowAgentAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, agentID string) (*dialogflow.GoogleCloudDialogflowCxV3Agent, error) {
	logger.Default.Logf(t, "Getting settings for Dialogflow agent %s in project %s", agentID, projectID)

	service, err := NewDialogflowServiceE(t, ctx, location)
	if err != nil {
		return nil, err
	}

	return GetDialogflowAgentAttrsWithClient(ctx, service, projectID, location, agentID)
}

// GetDialogflowAgentAttrsWithClient returns the settings Google Cloud holds for the given
// Dialogflow CX agent using the supplied *dialogflow.Service. Prefer this variant in unit tests
// where the service is backed by an httptest fake server (see dialogflow_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetDialogflowAgentAttrsWithClient(ctx context.Context, service *dialogflow.Service, projectID string, location string, agentID string) (*dialogflow.GoogleCloudDialogflowCxV3Agent, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/agents/%s", projectID, location, agentID)

	agent, err := service.Projects.Locations.Agents.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("Dialogflow agent %s does not exist in project %s", agentID, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Dialogflow agent %s in project %s: %w", agentID, projectID, err)
	}

	return agent, nil
}

// dialogflowLocationPattern is what a Google Cloud location identifier may contain. It is checked
// before a location reaches an endpoint, because a value carrying a slash, a colon or an at sign
// would build a URL pointing at a host of the caller's choosing rather than at Google.
var dialogflowLocationPattern = regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+)*$`)

// GetDialogflowCXGenerativeSettingsAttrs returns the settings Google Cloud holds for the given Dialogflow CX generative settings, so a test can assert on what was
// actually created rather than only that it exists.
// These settings hold the prompts and the banned phrases a generative agent works under, so the language they apply to is part of their name.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetDialogflowCXGenerativeSettingsAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, agentID string, languageCode string) *dialogflow.GoogleCloudDialogflowCxV3GenerativeSettings {
	attrs, err := GetDialogflowCXGenerativeSettingsAttrsE(t, ctx, projectID, location, agentID, languageCode)
	require.NoError(t, err)

	return attrs
}

// GetDialogflowCXGenerativeSettingsAttrsE returns the settings Google Cloud holds for the given Dialogflow CX generative settings.
// The ctx parameter supports cancellation and timeouts.
func GetDialogflowCXGenerativeSettingsAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, agentID string, languageCode string) (*dialogflow.GoogleCloudDialogflowCxV3GenerativeSettings, error) {
	logger.Default.Logf(t, "Getting settings for Dialogflow CX generative settings %s for agent %s in %s in project %s", languageCode, agentID, location, projectID)

	service, err := NewDialogflowServiceE(t, ctx, location)
	if err != nil {
		return nil, err
	}

	return GetDialogflowCXGenerativeSettingsAttrsWithClient(ctx, service, projectID, location, agentID, languageCode)
}

// GetDialogflowCXGenerativeSettingsAttrsWithClient returns the settings Google Cloud holds for the given Dialogflow CX generative settings using the supplied
// *dialogflow.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see dialogflow_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetDialogflowCXGenerativeSettingsAttrsWithClient(ctx context.Context, service *dialogflow.Service, projectID string, location string, agentID string, languageCode string) (*dialogflow.GoogleCloudDialogflowCxV3GenerativeSettings, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/agents/%s/generativeSettings", projectID, location, agentID)

	attrs, err := service.Projects.Locations.Agents.GetGenerativeSettings(name).LanguageCode(languageCode).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Dialogflow CX generative settings %s for agent %s in %s in project %s does not exist", languageCode, agentID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Dialogflow CX generative settings %s for agent %s in %s in project %s: %w", languageCode, agentID, location, projectID, err)
	}

	return attrs, nil
}

// GetDialogflowCXToolVersionAttrs returns the settings Google Cloud holds for the given Dialogflow CX tool version, so a test can assert on what was
// actually created rather than only that it exists.
// A tool version is a frozen copy of a tool a flow can call, so its display name is how a caller tells one from another.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetDialogflowCXToolVersionAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, agentID string, toolID string, versionID string) *dialogflow.GoogleCloudDialogflowCxV3ToolVersion {
	attrs, err := GetDialogflowCXToolVersionAttrsE(t, ctx, projectID, location, agentID, toolID, versionID)
	require.NoError(t, err)

	return attrs
}

// GetDialogflowCXToolVersionAttrsE returns the settings Google Cloud holds for the given Dialogflow CX tool version.
// The ctx parameter supports cancellation and timeouts.
func GetDialogflowCXToolVersionAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, agentID string, toolID string, versionID string) (*dialogflow.GoogleCloudDialogflowCxV3ToolVersion, error) {
	logger.Default.Logf(t, "Getting settings for Dialogflow CX tool version %s of tool %s for agent %s in %s in project %s", versionID, toolID, agentID, location, projectID)

	service, err := NewDialogflowServiceE(t, ctx, location)
	if err != nil {
		return nil, err
	}

	return GetDialogflowCXToolVersionAttrsWithClient(ctx, service, projectID, location, agentID, toolID, versionID)
}

// GetDialogflowCXToolVersionAttrsWithClient returns the settings Google Cloud holds for the given Dialogflow CX tool version using the supplied
// *dialogflow.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see dialogflow_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetDialogflowCXToolVersionAttrsWithClient(ctx context.Context, service *dialogflow.Service, projectID string, location string, agentID string, toolID string, versionID string) (*dialogflow.GoogleCloudDialogflowCxV3ToolVersion, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/agents/%s/tools/%s/versions/%s", projectID, location, agentID, toolID, versionID)

	attrs, err := service.Projects.Locations.Agents.Tools.Versions.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Dialogflow CX tool version %s of tool %s for agent %s in %s in project %s does not exist", versionID, toolID, agentID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Dialogflow CX tool version %s of tool %s for agent %s in %s in project %s: %w", versionID, toolID, agentID, location, projectID, err)
	}

	return attrs, nil
}

// NewDialogflowServiceE creates a Dialogflow service authenticated the same way every other client
// in this module is. Dialogflow answers on a per-location host for every location but `global`, so
// the location decides which endpoint the service talks to, and a location that is not a plain
// identifier is refused.
// The ctx parameter supports cancellation and timeouts.
func NewDialogflowServiceE(t testing.TestingT, ctx context.Context, location string) (*dialogflow.Service, error) {
	if !dialogflowLocationPattern.MatchString(location) {
		return nil, fmt.Errorf("%q is not a valid location: a location may hold only lowercase letters, digits and hyphens", location)
	}

	opts := append(withOptions(), option.WithScopes(dialogflow.CloudPlatformScope))
	if location != "global" {
		opts = append(opts, option.WithEndpoint(fmt.Sprintf("https://%s-dialogflow.googleapis.com/", location)))
	}

	return dialogflow.NewService(ctx, opts...)
}
