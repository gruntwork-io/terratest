package gcp

import (
	"context"
	"errors"
	"fmt"

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

// NewDialogflowServiceE creates a Dialogflow service authenticated the same way every other client
// in this module is. Dialogflow answers on a per-location host for every location but `global`, so
// the location decides which endpoint the service talks to.
// The ctx parameter supports cancellation and timeouts.
func NewDialogflowServiceE(t testing.TestingT, ctx context.Context, location string) (*dialogflow.Service, error) {
	opts := append(withOptions(), option.WithScopes(dialogflow.CloudPlatformScope))
	if location != "global" {
		opts = append(opts, option.WithEndpoint(fmt.Sprintf("https://%s-dialogflow.googleapis.com/", location)))
	}

	return dialogflow.NewService(ctx, opts...)
}
