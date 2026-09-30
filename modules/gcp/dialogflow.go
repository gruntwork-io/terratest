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

// GetDialogflowCXIntentAttrs returns the settings Google Cloud holds for the Dialogflow CX intent, so a test can assert on what was
// actually created rather than only that it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetDialogflowCXIntentAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, agentID string, intentID string) *dialogflow.GoogleCloudDialogflowCxV3Intent {
	result, err := GetDialogflowCXIntentAttrsE(t, ctx, projectID, location, agentID, intentID)
	require.NoError(t, err)

	return result
}

// GetDialogflowCXIntentAttrsE returns the settings Google Cloud holds for the Dialogflow CX intent.
// The ctx parameter supports cancellation and timeouts.
func GetDialogflowCXIntentAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, agentID string, intentID string) (*dialogflow.GoogleCloudDialogflowCxV3Intent, error) {
	logger.Default.Logf(t, "Getting settings for Dialogflow CX intent %s in %s in project %s", intentID, location, projectID)

	service, err := NewDialogflowServiceE(t, ctx, location)
	if err != nil {
		return nil, err
	}

	return GetDialogflowCXIntentAttrsWithClient(ctx, service, projectID, location, agentID, intentID)
}

// GetDialogflowCXIntentAttrsWithClient returns the settings Google Cloud holds for the Dialogflow CX intent using the supplied
// *dialogflow.Service. Prefer this variant in unit tests where the service is backed by an httptest
// fake server (see dialogflow_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetDialogflowCXIntentAttrsWithClient(ctx context.Context, service *dialogflow.Service, projectID string, location string, agentID string, intentID string) (*dialogflow.GoogleCloudDialogflowCxV3Intent, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/agents/%s/intents/%s", projectID, location, agentID, intentID)

	result, err := service.Projects.Locations.Agents.Intents.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Dialogflow CX intent %s in %s in project %s does not exist", intentID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Dialogflow CX intent %s in %s in project %s: %w", intentID, location, projectID, err)
	}

	return result, nil
}

// GetDialogflowCXEntityTypeAttrs returns the settings Google Cloud holds for the Dialogflow CX entity type, so a test can assert on what was
// actually created rather than only that it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetDialogflowCXEntityTypeAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, agentID string, entityTypeID string) *dialogflow.GoogleCloudDialogflowCxV3EntityType {
	result, err := GetDialogflowCXEntityTypeAttrsE(t, ctx, projectID, location, agentID, entityTypeID)
	require.NoError(t, err)

	return result
}

// GetDialogflowCXEntityTypeAttrsE returns the settings Google Cloud holds for the Dialogflow CX entity type.
// The ctx parameter supports cancellation and timeouts.
func GetDialogflowCXEntityTypeAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, agentID string, entityTypeID string) (*dialogflow.GoogleCloudDialogflowCxV3EntityType, error) {
	logger.Default.Logf(t, "Getting settings for Dialogflow CX entity type %s in %s in project %s", entityTypeID, location, projectID)

	service, err := NewDialogflowServiceE(t, ctx, location)
	if err != nil {
		return nil, err
	}

	return GetDialogflowCXEntityTypeAttrsWithClient(ctx, service, projectID, location, agentID, entityTypeID)
}

// GetDialogflowCXEntityTypeAttrsWithClient returns the settings Google Cloud holds for the Dialogflow CX entity type using the supplied
// *dialogflow.Service. Prefer this variant in unit tests where the service is backed by an httptest
// fake server (see dialogflow_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetDialogflowCXEntityTypeAttrsWithClient(ctx context.Context, service *dialogflow.Service, projectID string, location string, agentID string, entityTypeID string) (*dialogflow.GoogleCloudDialogflowCxV3EntityType, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/agents/%s/entityTypes/%s", projectID, location, agentID, entityTypeID)

	result, err := service.Projects.Locations.Agents.EntityTypes.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Dialogflow CX entity type %s in %s in project %s does not exist", entityTypeID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Dialogflow CX entity type %s in %s in project %s: %w", entityTypeID, location, projectID, err)
	}

	return result, nil
}

// GetDialogflowCXFlowAttrs returns the settings Google Cloud holds for the Dialogflow CX flow, so a test can assert on what was
// actually created rather than only that it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetDialogflowCXFlowAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, agentID string, flowID string) *dialogflow.GoogleCloudDialogflowCxV3Flow {
	result, err := GetDialogflowCXFlowAttrsE(t, ctx, projectID, location, agentID, flowID)
	require.NoError(t, err)

	return result
}

// GetDialogflowCXFlowAttrsE returns the settings Google Cloud holds for the Dialogflow CX flow.
// The ctx parameter supports cancellation and timeouts.
func GetDialogflowCXFlowAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, agentID string, flowID string) (*dialogflow.GoogleCloudDialogflowCxV3Flow, error) {
	logger.Default.Logf(t, "Getting settings for Dialogflow CX flow %s in %s in project %s", flowID, location, projectID)

	service, err := NewDialogflowServiceE(t, ctx, location)
	if err != nil {
		return nil, err
	}

	return GetDialogflowCXFlowAttrsWithClient(ctx, service, projectID, location, agentID, flowID)
}

// GetDialogflowCXFlowAttrsWithClient returns the settings Google Cloud holds for the Dialogflow CX flow using the supplied
// *dialogflow.Service. Prefer this variant in unit tests where the service is backed by an httptest
// fake server (see dialogflow_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetDialogflowCXFlowAttrsWithClient(ctx context.Context, service *dialogflow.Service, projectID string, location string, agentID string, flowID string) (*dialogflow.GoogleCloudDialogflowCxV3Flow, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/agents/%s/flows/%s", projectID, location, agentID, flowID)

	result, err := service.Projects.Locations.Agents.Flows.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Dialogflow CX flow %s in %s in project %s does not exist", flowID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Dialogflow CX flow %s in %s in project %s: %w", flowID, location, projectID, err)
	}

	return result, nil
}

// GetDialogflowCXWebhookAttrs returns the settings Google Cloud holds for the Dialogflow CX webhook, so a test can assert on what was
// actually created rather than only that it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetDialogflowCXWebhookAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, agentID string, webhookID string) *dialogflow.GoogleCloudDialogflowCxV3Webhook {
	result, err := GetDialogflowCXWebhookAttrsE(t, ctx, projectID, location, agentID, webhookID)
	require.NoError(t, err)

	return result
}

// GetDialogflowCXWebhookAttrsE returns the settings Google Cloud holds for the Dialogflow CX webhook.
// The ctx parameter supports cancellation and timeouts.
func GetDialogflowCXWebhookAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, agentID string, webhookID string) (*dialogflow.GoogleCloudDialogflowCxV3Webhook, error) {
	logger.Default.Logf(t, "Getting settings for Dialogflow CX webhook %s in %s in project %s", webhookID, location, projectID)

	service, err := NewDialogflowServiceE(t, ctx, location)
	if err != nil {
		return nil, err
	}

	return GetDialogflowCXWebhookAttrsWithClient(ctx, service, projectID, location, agentID, webhookID)
}

// GetDialogflowCXWebhookAttrsWithClient returns the settings Google Cloud holds for the Dialogflow CX webhook using the supplied
// *dialogflow.Service. Prefer this variant in unit tests where the service is backed by an httptest
// fake server (see dialogflow_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetDialogflowCXWebhookAttrsWithClient(ctx context.Context, service *dialogflow.Service, projectID string, location string, agentID string, webhookID string) (*dialogflow.GoogleCloudDialogflowCxV3Webhook, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/agents/%s/webhooks/%s", projectID, location, agentID, webhookID)

	result, err := service.Projects.Locations.Agents.Webhooks.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Dialogflow CX webhook %s in %s in project %s does not exist", webhookID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Dialogflow CX webhook %s in %s in project %s: %w", webhookID, location, projectID, err)
	}

	return result, nil
}

// GetDialogflowCXEnvironmentAttrs returns the settings Google Cloud holds for the Dialogflow CX environment, so a test can assert on what was
// actually created rather than only that it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetDialogflowCXEnvironmentAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, agentID string, environmentID string) *dialogflow.GoogleCloudDialogflowCxV3Environment {
	result, err := GetDialogflowCXEnvironmentAttrsE(t, ctx, projectID, location, agentID, environmentID)
	require.NoError(t, err)

	return result
}

// GetDialogflowCXEnvironmentAttrsE returns the settings Google Cloud holds for the Dialogflow CX environment.
// The ctx parameter supports cancellation and timeouts.
func GetDialogflowCXEnvironmentAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, agentID string, environmentID string) (*dialogflow.GoogleCloudDialogflowCxV3Environment, error) {
	logger.Default.Logf(t, "Getting settings for Dialogflow CX environment %s in %s in project %s", environmentID, location, projectID)

	service, err := NewDialogflowServiceE(t, ctx, location)
	if err != nil {
		return nil, err
	}

	return GetDialogflowCXEnvironmentAttrsWithClient(ctx, service, projectID, location, agentID, environmentID)
}

// GetDialogflowCXEnvironmentAttrsWithClient returns the settings Google Cloud holds for the Dialogflow CX environment using the supplied
// *dialogflow.Service. Prefer this variant in unit tests where the service is backed by an httptest
// fake server (see dialogflow_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetDialogflowCXEnvironmentAttrsWithClient(ctx context.Context, service *dialogflow.Service, projectID string, location string, agentID string, environmentID string) (*dialogflow.GoogleCloudDialogflowCxV3Environment, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/agents/%s/environments/%s", projectID, location, agentID, environmentID)

	result, err := service.Projects.Locations.Agents.Environments.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Dialogflow CX environment %s in %s in project %s does not exist", environmentID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Dialogflow CX environment %s in %s in project %s: %w", environmentID, location, projectID, err)
	}

	return result, nil
}

// GetDialogflowCXTestCaseAttrs returns the settings Google Cloud holds for the Dialogflow CX test case, so a test can assert on what was
// actually created rather than only that it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetDialogflowCXTestCaseAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, agentID string, testCaseID string) *dialogflow.GoogleCloudDialogflowCxV3TestCase {
	result, err := GetDialogflowCXTestCaseAttrsE(t, ctx, projectID, location, agentID, testCaseID)
	require.NoError(t, err)

	return result
}

// GetDialogflowCXTestCaseAttrsE returns the settings Google Cloud holds for the Dialogflow CX test case.
// The ctx parameter supports cancellation and timeouts.
func GetDialogflowCXTestCaseAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, agentID string, testCaseID string) (*dialogflow.GoogleCloudDialogflowCxV3TestCase, error) {
	logger.Default.Logf(t, "Getting settings for Dialogflow CX test case %s in %s in project %s", testCaseID, location, projectID)

	service, err := NewDialogflowServiceE(t, ctx, location)
	if err != nil {
		return nil, err
	}

	return GetDialogflowCXTestCaseAttrsWithClient(ctx, service, projectID, location, agentID, testCaseID)
}

// GetDialogflowCXTestCaseAttrsWithClient returns the settings Google Cloud holds for the Dialogflow CX test case using the supplied
// *dialogflow.Service. Prefer this variant in unit tests where the service is backed by an httptest
// fake server (see dialogflow_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetDialogflowCXTestCaseAttrsWithClient(ctx context.Context, service *dialogflow.Service, projectID string, location string, agentID string, testCaseID string) (*dialogflow.GoogleCloudDialogflowCxV3TestCase, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/agents/%s/testCases/%s", projectID, location, agentID, testCaseID)

	result, err := service.Projects.Locations.Agents.TestCases.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Dialogflow CX test case %s in %s in project %s does not exist", testCaseID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Dialogflow CX test case %s in %s in project %s: %w", testCaseID, location, projectID, err)
	}

	return result, nil
}

// GetDialogflowCXGeneratorAttrs returns the settings Google Cloud holds for the Dialogflow CX generator, so a test can assert on what was
// actually created rather than only that it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetDialogflowCXGeneratorAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, agentID string, generatorID string) *dialogflow.GoogleCloudDialogflowCxV3Generator {
	result, err := GetDialogflowCXGeneratorAttrsE(t, ctx, projectID, location, agentID, generatorID)
	require.NoError(t, err)

	return result
}

// GetDialogflowCXGeneratorAttrsE returns the settings Google Cloud holds for the Dialogflow CX generator.
// The ctx parameter supports cancellation and timeouts.
func GetDialogflowCXGeneratorAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, agentID string, generatorID string) (*dialogflow.GoogleCloudDialogflowCxV3Generator, error) {
	logger.Default.Logf(t, "Getting settings for Dialogflow CX generator %s in %s in project %s", generatorID, location, projectID)

	service, err := NewDialogflowServiceE(t, ctx, location)
	if err != nil {
		return nil, err
	}

	return GetDialogflowCXGeneratorAttrsWithClient(ctx, service, projectID, location, agentID, generatorID)
}

// GetDialogflowCXGeneratorAttrsWithClient returns the settings Google Cloud holds for the Dialogflow CX generator using the supplied
// *dialogflow.Service. Prefer this variant in unit tests where the service is backed by an httptest
// fake server (see dialogflow_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetDialogflowCXGeneratorAttrsWithClient(ctx context.Context, service *dialogflow.Service, projectID string, location string, agentID string, generatorID string) (*dialogflow.GoogleCloudDialogflowCxV3Generator, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/agents/%s/generators/%s", projectID, location, agentID, generatorID)

	result, err := service.Projects.Locations.Agents.Generators.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Dialogflow CX generator %s in %s in project %s does not exist", generatorID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Dialogflow CX generator %s in %s in project %s: %w", generatorID, location, projectID, err)
	}

	return result, nil
}

// GetDialogflowCXToolAttrs returns the settings Google Cloud holds for the Dialogflow CX tool, so a test can assert on what was
// actually created rather than only that it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetDialogflowCXToolAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, agentID string, toolID string) *dialogflow.GoogleCloudDialogflowCxV3Tool {
	result, err := GetDialogflowCXToolAttrsE(t, ctx, projectID, location, agentID, toolID)
	require.NoError(t, err)

	return result
}

// GetDialogflowCXToolAttrsE returns the settings Google Cloud holds for the Dialogflow CX tool.
// The ctx parameter supports cancellation and timeouts.
func GetDialogflowCXToolAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, agentID string, toolID string) (*dialogflow.GoogleCloudDialogflowCxV3Tool, error) {
	logger.Default.Logf(t, "Getting settings for Dialogflow CX tool %s in %s in project %s", toolID, location, projectID)

	service, err := NewDialogflowServiceE(t, ctx, location)
	if err != nil {
		return nil, err
	}

	return GetDialogflowCXToolAttrsWithClient(ctx, service, projectID, location, agentID, toolID)
}

// GetDialogflowCXToolAttrsWithClient returns the settings Google Cloud holds for the Dialogflow CX tool using the supplied
// *dialogflow.Service. Prefer this variant in unit tests where the service is backed by an httptest
// fake server (see dialogflow_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetDialogflowCXToolAttrsWithClient(ctx context.Context, service *dialogflow.Service, projectID string, location string, agentID string, toolID string) (*dialogflow.GoogleCloudDialogflowCxV3Tool, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/agents/%s/tools/%s", projectID, location, agentID, toolID)

	result, err := service.Projects.Locations.Agents.Tools.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Dialogflow CX tool %s in %s in project %s does not exist", toolID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Dialogflow CX tool %s in %s in project %s: %w", toolID, location, projectID, err)
	}

	return result, nil
}

// GetDialogflowCXPlaybookAttrs returns the settings Google Cloud holds for the Dialogflow CX playbook, so a test can assert on what was
// actually created rather than only that it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetDialogflowCXPlaybookAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, agentID string, playbookID string) *dialogflow.GoogleCloudDialogflowCxV3Playbook {
	result, err := GetDialogflowCXPlaybookAttrsE(t, ctx, projectID, location, agentID, playbookID)
	require.NoError(t, err)

	return result
}

// GetDialogflowCXPlaybookAttrsE returns the settings Google Cloud holds for the Dialogflow CX playbook.
// The ctx parameter supports cancellation and timeouts.
func GetDialogflowCXPlaybookAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, agentID string, playbookID string) (*dialogflow.GoogleCloudDialogflowCxV3Playbook, error) {
	logger.Default.Logf(t, "Getting settings for Dialogflow CX playbook %s in %s in project %s", playbookID, location, projectID)

	service, err := NewDialogflowServiceE(t, ctx, location)
	if err != nil {
		return nil, err
	}

	return GetDialogflowCXPlaybookAttrsWithClient(ctx, service, projectID, location, agentID, playbookID)
}

// GetDialogflowCXPlaybookAttrsWithClient returns the settings Google Cloud holds for the Dialogflow CX playbook using the supplied
// *dialogflow.Service. Prefer this variant in unit tests where the service is backed by an httptest
// fake server (see dialogflow_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetDialogflowCXPlaybookAttrsWithClient(ctx context.Context, service *dialogflow.Service, projectID string, location string, agentID string, playbookID string) (*dialogflow.GoogleCloudDialogflowCxV3Playbook, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/agents/%s/playbooks/%s", projectID, location, agentID, playbookID)

	result, err := service.Projects.Locations.Agents.Playbooks.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Dialogflow CX playbook %s in %s in project %s does not exist", playbookID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Dialogflow CX playbook %s in %s in project %s: %w", playbookID, location, projectID, err)
	}

	return result, nil
}

// GetDialogflowCXPageAttrs returns the settings Google Cloud holds for the Dialogflow CX page, so a test can assert on what was
// actually created rather than only that it exists. A page belongs to a flow rather than
// straight to the agent, so the caller names both.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetDialogflowCXPageAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, agentID string, flowID string, pageID string) *dialogflow.GoogleCloudDialogflowCxV3Page {
	result, err := GetDialogflowCXPageAttrsE(t, ctx, projectID, location, agentID, flowID, pageID)
	require.NoError(t, err)

	return result
}

// GetDialogflowCXPageAttrsE returns the settings Google Cloud holds for the Dialogflow CX page.
// The ctx parameter supports cancellation and timeouts.
func GetDialogflowCXPageAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, agentID string, flowID string, pageID string) (*dialogflow.GoogleCloudDialogflowCxV3Page, error) {
	logger.Default.Logf(t, "Getting settings for Dialogflow CX page %s in %s in project %s", pageID, location, projectID)

	service, err := NewDialogflowServiceE(t, ctx, location)
	if err != nil {
		return nil, err
	}

	return GetDialogflowCXPageAttrsWithClient(ctx, service, projectID, location, agentID, flowID, pageID)
}

// GetDialogflowCXPageAttrsWithClient returns the settings Google Cloud holds for the Dialogflow CX page using the supplied
// *dialogflow.Service. Prefer this variant in unit tests where the service is backed by an httptest
// fake server (see dialogflow_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetDialogflowCXPageAttrsWithClient(ctx context.Context, service *dialogflow.Service, projectID string, location string, agentID string, flowID string, pageID string) (*dialogflow.GoogleCloudDialogflowCxV3Page, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/agents/%s/flows/%s/pages/%s", projectID, location, agentID, flowID, pageID)

	result, err := service.Projects.Locations.Agents.Flows.Pages.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Dialogflow CX page %s in %s in project %s does not exist", pageID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Dialogflow CX page %s in %s in project %s: %w", pageID, location, projectID, err)
	}

	return result, nil
}

// GetDialogflowCXVersionAttrs returns the settings Google Cloud holds for the Dialogflow CX version, so a test can assert on what was
// actually created rather than only that it exists. A version is a snapshot of one flow rather
// than of the whole agent, so the caller names both.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetDialogflowCXVersionAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, agentID string, flowID string, versionID string) *dialogflow.GoogleCloudDialogflowCxV3Version {
	result, err := GetDialogflowCXVersionAttrsE(t, ctx, projectID, location, agentID, flowID, versionID)
	require.NoError(t, err)

	return result
}

// GetDialogflowCXVersionAttrsE returns the settings Google Cloud holds for the Dialogflow CX version.
// The ctx parameter supports cancellation and timeouts.
func GetDialogflowCXVersionAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, agentID string, flowID string, versionID string) (*dialogflow.GoogleCloudDialogflowCxV3Version, error) {
	logger.Default.Logf(t, "Getting settings for Dialogflow CX version %s in %s in project %s", versionID, location, projectID)

	service, err := NewDialogflowServiceE(t, ctx, location)
	if err != nil {
		return nil, err
	}

	return GetDialogflowCXVersionAttrsWithClient(ctx, service, projectID, location, agentID, flowID, versionID)
}

// GetDialogflowCXVersionAttrsWithClient returns the settings Google Cloud holds for the Dialogflow CX version using the supplied
// *dialogflow.Service. Prefer this variant in unit tests where the service is backed by an httptest
// fake server (see dialogflow_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetDialogflowCXVersionAttrsWithClient(ctx context.Context, service *dialogflow.Service, projectID string, location string, agentID string, flowID string, versionID string) (*dialogflow.GoogleCloudDialogflowCxV3Version, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/agents/%s/flows/%s/versions/%s", projectID, location, agentID, flowID, versionID)

	result, err := service.Projects.Locations.Agents.Flows.Versions.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Dialogflow CX version %s in %s in project %s does not exist", versionID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Dialogflow CX version %s in %s in project %s: %w", versionID, location, projectID, err)
	}

	return result, nil
}

// GetDialogflowCXSecuritySettingsAttrs returns the settings Google Cloud holds for the Dialogflow CX security settings, so a test can assert on what was
// actually created rather than only that it exists. Security settings belong to a project and
// location rather than to an agent, so no agent is named.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetDialogflowCXSecuritySettingsAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, settingsID string) *dialogflow.GoogleCloudDialogflowCxV3SecuritySettings {
	result, err := GetDialogflowCXSecuritySettingsAttrsE(t, ctx, projectID, location, settingsID)
	require.NoError(t, err)

	return result
}

// GetDialogflowCXSecuritySettingsAttrsE returns the settings Google Cloud holds for the Dialogflow CX security settings.
// The ctx parameter supports cancellation and timeouts.
func GetDialogflowCXSecuritySettingsAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, settingsID string) (*dialogflow.GoogleCloudDialogflowCxV3SecuritySettings, error) {
	logger.Default.Logf(t, "Getting settings for Dialogflow CX security settings %s in %s in project %s", settingsID, location, projectID)

	service, err := NewDialogflowServiceE(t, ctx, location)
	if err != nil {
		return nil, err
	}

	return GetDialogflowCXSecuritySettingsAttrsWithClient(ctx, service, projectID, location, settingsID)
}

// GetDialogflowCXSecuritySettingsAttrsWithClient returns the settings Google Cloud holds for the Dialogflow CX security settings using the supplied
// *dialogflow.Service. Prefer this variant in unit tests where the service is backed by an httptest
// fake server (see dialogflow_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetDialogflowCXSecuritySettingsAttrsWithClient(ctx context.Context, service *dialogflow.Service, projectID string, location string, settingsID string) (*dialogflow.GoogleCloudDialogflowCxV3SecuritySettings, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/securitySettings/%s", projectID, location, settingsID)

	result, err := service.Projects.Locations.SecuritySettings.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Dialogflow CX security settings %s in %s in project %s does not exist", settingsID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Dialogflow CX security settings %s in %s in project %s: %w", settingsID, location, projectID, err)
	}

	return result, nil
}
