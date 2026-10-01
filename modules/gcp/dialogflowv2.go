package gcp

import (
	"context"
	"errors"
	"fmt"

	"github.com/gruntwork-io/terratest/modules/core/v2/logger"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
	dialogflowv2 "google.golang.org/api/dialogflow/v2"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
)

// Dialogflow serves two generations of API from one host, and they are different products rather than
// two versions of one. ES agents, their intents and entity types live in v2; CX agents and their flows
// live in v3, which dialogflow.go covers. The reads here say V2 in their names so a caller cannot pick
// the wrong generation by accident.

// GetDialogflowV2AgentAttrs returns the settings Google Cloud holds for the given Dialogflow ES agent, so a test can assert on what was
// actually created rather than only that it exists.
// An agent is the whole conversational model, so its language, its time zone and its match mode decide how every request is understood.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetDialogflowV2AgentAttrs(t testing.TestingT, ctx context.Context, projectID string) *dialogflowv2.GoogleCloudDialogflowV2Agent {
	attrs, err := GetDialogflowV2AgentAttrsE(t, ctx, projectID)
	require.NoError(t, err)

	return attrs
}

// GetDialogflowV2AgentAttrsE returns the settings Google Cloud holds for the given Dialogflow ES agent.
// The ctx parameter supports cancellation and timeouts.
func GetDialogflowV2AgentAttrsE(t testing.TestingT, ctx context.Context, projectID string) (*dialogflowv2.GoogleCloudDialogflowV2Agent, error) {
	logger.Default.Logf(t, "Getting settings for Dialogflow ES agent in project %s", projectID)

	service, err := NewDialogflowV2ServiceE(t, ctx, "global")
	if err != nil {
		return nil, err
	}

	return GetDialogflowV2AgentAttrsWithClient(ctx, service, projectID)
}

// GetDialogflowV2AgentAttrsWithClient returns the settings Google Cloud holds for the given Dialogflow ES agent using the supplied
// *dialogflowv2.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see dialogflowv2_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetDialogflowV2AgentAttrsWithClient(ctx context.Context, service *dialogflowv2.Service, projectID string) (*dialogflowv2.GoogleCloudDialogflowV2Agent, error) {
	parent := "projects/" + projectID

	attrs, err := service.Projects.GetAgent(parent).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Dialogflow ES agent in project %s does not exist", projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Dialogflow ES agent in project %s: %w", projectID, err)
	}

	return attrs, nil
}

// GetDialogflowV2ConversationProfileAttrs returns the settings Google Cloud holds for the given Dialogflow ES conversation profile, so a test can assert on what was
// actually created rather than only that it exists.
// A profile is the set of suggestion features a conversation uses, so its language and its display name are what a caller picks it by.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetDialogflowV2ConversationProfileAttrs(t testing.TestingT, ctx context.Context, projectID string, profileID string) *dialogflowv2.GoogleCloudDialogflowV2ConversationProfile {
	attrs, err := GetDialogflowV2ConversationProfileAttrsE(t, ctx, projectID, profileID)
	require.NoError(t, err)

	return attrs
}

// GetDialogflowV2ConversationProfileAttrsE returns the settings Google Cloud holds for the given Dialogflow ES conversation profile.
// The ctx parameter supports cancellation and timeouts.
func GetDialogflowV2ConversationProfileAttrsE(t testing.TestingT, ctx context.Context, projectID string, profileID string) (*dialogflowv2.GoogleCloudDialogflowV2ConversationProfile, error) {
	logger.Default.Logf(t, "Getting settings for Dialogflow ES conversation profile %s in project %s", profileID, projectID)

	service, err := NewDialogflowV2ServiceE(t, ctx, "global")
	if err != nil {
		return nil, err
	}

	return GetDialogflowV2ConversationProfileAttrsWithClient(ctx, service, projectID, profileID)
}

// GetDialogflowV2ConversationProfileAttrsWithClient returns the settings Google Cloud holds for the given Dialogflow ES conversation profile using the supplied
// *dialogflowv2.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see dialogflowv2_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetDialogflowV2ConversationProfileAttrsWithClient(ctx context.Context, service *dialogflowv2.Service, projectID string, profileID string) (*dialogflowv2.GoogleCloudDialogflowV2ConversationProfile, error) {
	name := fmt.Sprintf("projects/%s/conversationProfiles/%s", projectID, profileID)

	attrs, err := service.Projects.ConversationProfiles.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Dialogflow ES conversation profile %s in project %s does not exist", profileID, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Dialogflow ES conversation profile %s in project %s: %w", profileID, projectID, err)
	}

	return attrs, nil
}

// GetDialogflowV2EncryptionSpecAttrs returns the settings Google Cloud holds for the given Dialogflow ES encryption spec, so a test can assert on what was
// actually created rather than only that it exists.
// The spec names the key a region's conversation data is encrypted with, so which key it is is the only thing it says.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetDialogflowV2EncryptionSpecAttrs(t testing.TestingT, ctx context.Context, projectID string, location string) *dialogflowv2.GoogleCloudDialogflowV2EncryptionSpec {
	attrs, err := GetDialogflowV2EncryptionSpecAttrsE(t, ctx, projectID, location)
	require.NoError(t, err)

	return attrs
}

// GetDialogflowV2EncryptionSpecAttrsE returns the settings Google Cloud holds for the given Dialogflow ES encryption spec.
// The ctx parameter supports cancellation and timeouts.
func GetDialogflowV2EncryptionSpecAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string) (*dialogflowv2.GoogleCloudDialogflowV2EncryptionSpec, error) {
	logger.Default.Logf(t, "Getting settings for Dialogflow ES encryption spec in %s in project %s", location, projectID)

	service, err := NewDialogflowV2ServiceE(t, ctx, location)
	if err != nil {
		return nil, err
	}

	return GetDialogflowV2EncryptionSpecAttrsWithClient(ctx, service, projectID, location)
}

// GetDialogflowV2EncryptionSpecAttrsWithClient returns the settings Google Cloud holds for the given Dialogflow ES encryption spec using the supplied
// *dialogflowv2.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see dialogflowv2_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetDialogflowV2EncryptionSpecAttrsWithClient(ctx context.Context, service *dialogflowv2.Service, projectID string, location string) (*dialogflowv2.GoogleCloudDialogflowV2EncryptionSpec, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/encryptionSpec", projectID, location)

	attrs, err := service.Projects.Locations.GetEncryptionSpec(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Dialogflow ES encryption spec in %s in project %s does not exist", location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Dialogflow ES encryption spec in %s in project %s: %w", location, projectID, err)
	}

	return attrs, nil
}

// GetDialogflowV2EntityTypeAttrs returns the settings Google Cloud holds for the given Dialogflow ES entity type, so a test can assert on what was
// actually created rather than only that it exists.
// An entity type is what the agent extracts from what a user says, so its kind and whether it auto-expands decide what it will match.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetDialogflowV2EntityTypeAttrs(t testing.TestingT, ctx context.Context, projectID string, entityTypeID string) *dialogflowv2.GoogleCloudDialogflowV2EntityType {
	attrs, err := GetDialogflowV2EntityTypeAttrsE(t, ctx, projectID, entityTypeID)
	require.NoError(t, err)

	return attrs
}

// GetDialogflowV2EntityTypeAttrsE returns the settings Google Cloud holds for the given Dialogflow ES entity type.
// The ctx parameter supports cancellation and timeouts.
func GetDialogflowV2EntityTypeAttrsE(t testing.TestingT, ctx context.Context, projectID string, entityTypeID string) (*dialogflowv2.GoogleCloudDialogflowV2EntityType, error) {
	logger.Default.Logf(t, "Getting settings for Dialogflow ES entity type %s in project %s", entityTypeID, projectID)

	service, err := NewDialogflowV2ServiceE(t, ctx, "global")
	if err != nil {
		return nil, err
	}

	return GetDialogflowV2EntityTypeAttrsWithClient(ctx, service, projectID, entityTypeID)
}

// GetDialogflowV2EntityTypeAttrsWithClient returns the settings Google Cloud holds for the given Dialogflow ES entity type using the supplied
// *dialogflowv2.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see dialogflowv2_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetDialogflowV2EntityTypeAttrsWithClient(ctx context.Context, service *dialogflowv2.Service, projectID string, entityTypeID string) (*dialogflowv2.GoogleCloudDialogflowV2EntityType, error) {
	name := fmt.Sprintf("projects/%s/agent/entityTypes/%s", projectID, entityTypeID)

	attrs, err := service.Projects.Agent.EntityTypes.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Dialogflow ES entity type %s in project %s does not exist", entityTypeID, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Dialogflow ES entity type %s in project %s: %w", entityTypeID, projectID, err)
	}

	return attrs, nil
}

// GetDialogflowV2EnvironmentAttrs returns the settings Google Cloud holds for the given Dialogflow ES environment, so a test can assert on what was
// actually created rather than only that it exists.
// An environment pins one version of the agent for a set of users, so which version it serves is the whole point.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetDialogflowV2EnvironmentAttrs(t testing.TestingT, ctx context.Context, projectID string, environmentID string) *dialogflowv2.GoogleCloudDialogflowV2Environment {
	attrs, err := GetDialogflowV2EnvironmentAttrsE(t, ctx, projectID, environmentID)
	require.NoError(t, err)

	return attrs
}

// GetDialogflowV2EnvironmentAttrsE returns the settings Google Cloud holds for the given Dialogflow ES environment.
// The ctx parameter supports cancellation and timeouts.
func GetDialogflowV2EnvironmentAttrsE(t testing.TestingT, ctx context.Context, projectID string, environmentID string) (*dialogflowv2.GoogleCloudDialogflowV2Environment, error) {
	logger.Default.Logf(t, "Getting settings for Dialogflow ES environment %s in project %s", environmentID, projectID)

	service, err := NewDialogflowV2ServiceE(t, ctx, "global")
	if err != nil {
		return nil, err
	}

	return GetDialogflowV2EnvironmentAttrsWithClient(ctx, service, projectID, environmentID)
}

// GetDialogflowV2EnvironmentAttrsWithClient returns the settings Google Cloud holds for the given Dialogflow ES environment using the supplied
// *dialogflowv2.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see dialogflowv2_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetDialogflowV2EnvironmentAttrsWithClient(ctx context.Context, service *dialogflowv2.Service, projectID string, environmentID string) (*dialogflowv2.GoogleCloudDialogflowV2Environment, error) {
	name := fmt.Sprintf("projects/%s/agent/environments/%s", projectID, environmentID)

	attrs, err := service.Projects.Agent.Environments.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Dialogflow ES environment %s in project %s does not exist", environmentID, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Dialogflow ES environment %s in project %s: %w", environmentID, projectID, err)
	}

	return attrs, nil
}

// GetDialogflowV2FulfillmentAttrs returns the settings Google Cloud holds for the given Dialogflow ES fulfillment, so a test can assert on what was
// actually created rather than only that it exists.
// Fulfillment is the webhook the agent calls, so the URL it posts to and whether it is enabled decide whether anything happens at all.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetDialogflowV2FulfillmentAttrs(t testing.TestingT, ctx context.Context, projectID string) *dialogflowv2.GoogleCloudDialogflowV2Fulfillment {
	attrs, err := GetDialogflowV2FulfillmentAttrsE(t, ctx, projectID)
	require.NoError(t, err)

	return attrs
}

// GetDialogflowV2FulfillmentAttrsE returns the settings Google Cloud holds for the given Dialogflow ES fulfillment.
// The ctx parameter supports cancellation and timeouts.
func GetDialogflowV2FulfillmentAttrsE(t testing.TestingT, ctx context.Context, projectID string) (*dialogflowv2.GoogleCloudDialogflowV2Fulfillment, error) {
	logger.Default.Logf(t, "Getting settings for Dialogflow ES fulfillment in project %s", projectID)

	service, err := NewDialogflowV2ServiceE(t, ctx, "global")
	if err != nil {
		return nil, err
	}

	return GetDialogflowV2FulfillmentAttrsWithClient(ctx, service, projectID)
}

// GetDialogflowV2FulfillmentAttrsWithClient returns the settings Google Cloud holds for the given Dialogflow ES fulfillment using the supplied
// *dialogflowv2.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see dialogflowv2_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetDialogflowV2FulfillmentAttrsWithClient(ctx context.Context, service *dialogflowv2.Service, projectID string) (*dialogflowv2.GoogleCloudDialogflowV2Fulfillment, error) {
	name := fmt.Sprintf("projects/%s/agent/fulfillment", projectID)

	attrs, err := service.Projects.Agent.GetFulfillment(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Dialogflow ES fulfillment in project %s does not exist", projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Dialogflow ES fulfillment in project %s: %w", projectID, err)
	}

	return attrs, nil
}

// GetDialogflowV2GeneratorAttrs returns the settings Google Cloud holds for the given Dialogflow ES generator, so a test can assert on what was
// actually created rather than only that it exists.
// A generator is the prompt a summary or suggestion is produced from, so its trigger event and its description are what it is for.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetDialogflowV2GeneratorAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, generatorID string) *dialogflowv2.GoogleCloudDialogflowV2Generator {
	attrs, err := GetDialogflowV2GeneratorAttrsE(t, ctx, projectID, location, generatorID)
	require.NoError(t, err)

	return attrs
}

// GetDialogflowV2GeneratorAttrsE returns the settings Google Cloud holds for the given Dialogflow ES generator.
// The ctx parameter supports cancellation and timeouts.
func GetDialogflowV2GeneratorAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, generatorID string) (*dialogflowv2.GoogleCloudDialogflowV2Generator, error) {
	logger.Default.Logf(t, "Getting settings for Dialogflow ES generator %s in %s in project %s", generatorID, location, projectID)

	service, err := NewDialogflowV2ServiceE(t, ctx, location)
	if err != nil {
		return nil, err
	}

	return GetDialogflowV2GeneratorAttrsWithClient(ctx, service, projectID, location, generatorID)
}

// GetDialogflowV2GeneratorAttrsWithClient returns the settings Google Cloud holds for the given Dialogflow ES generator using the supplied
// *dialogflowv2.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see dialogflowv2_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetDialogflowV2GeneratorAttrsWithClient(ctx context.Context, service *dialogflowv2.Service, projectID string, location string, generatorID string) (*dialogflowv2.GoogleCloudDialogflowV2Generator, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/generators/%s", projectID, location, generatorID)

	attrs, err := service.Projects.Locations.Generators.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Dialogflow ES generator %s in %s in project %s does not exist", generatorID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Dialogflow ES generator %s in %s in project %s: %w", generatorID, location, projectID, err)
	}

	return attrs, nil
}

// GetDialogflowV2IntentAttrs returns the settings Google Cloud holds for the given Dialogflow ES intent, so a test can assert on what was
// actually created rather than only that it exists.
// An intent is what the agent does when it recognises something, so its priority and whether it ends the conversation decide the flow.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetDialogflowV2IntentAttrs(t testing.TestingT, ctx context.Context, projectID string, intentID string) *dialogflowv2.GoogleCloudDialogflowV2Intent {
	attrs, err := GetDialogflowV2IntentAttrsE(t, ctx, projectID, intentID)
	require.NoError(t, err)

	return attrs
}

// GetDialogflowV2IntentAttrsE returns the settings Google Cloud holds for the given Dialogflow ES intent.
// The ctx parameter supports cancellation and timeouts.
func GetDialogflowV2IntentAttrsE(t testing.TestingT, ctx context.Context, projectID string, intentID string) (*dialogflowv2.GoogleCloudDialogflowV2Intent, error) {
	logger.Default.Logf(t, "Getting settings for Dialogflow ES intent %s in project %s", intentID, projectID)

	service, err := NewDialogflowV2ServiceE(t, ctx, "global")
	if err != nil {
		return nil, err
	}

	return GetDialogflowV2IntentAttrsWithClient(ctx, service, projectID, intentID)
}

// GetDialogflowV2IntentAttrsWithClient returns the settings Google Cloud holds for the given Dialogflow ES intent using the supplied
// *dialogflowv2.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see dialogflowv2_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetDialogflowV2IntentAttrsWithClient(ctx context.Context, service *dialogflowv2.Service, projectID string, intentID string) (*dialogflowv2.GoogleCloudDialogflowV2Intent, error) {
	name := fmt.Sprintf("projects/%s/agent/intents/%s", projectID, intentID)

	attrs, err := service.Projects.Agent.Intents.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Dialogflow ES intent %s in project %s does not exist", intentID, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Dialogflow ES intent %s in project %s: %w", intentID, projectID, err)
	}

	return attrs, nil
}

// GetDialogflowV2SipTrunkAttrs returns the settings Google Cloud holds for the given Dialogflow ES SIP trunk, so a test can assert on what was
// actually created rather than only that it exists.
// A trunk is how telephone calls reach the agent, so the connections it names and its expected hostname decide which calls arrive.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetDialogflowV2SipTrunkAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, trunkID string) *dialogflowv2.GoogleCloudDialogflowV2SipTrunk {
	attrs, err := GetDialogflowV2SipTrunkAttrsE(t, ctx, projectID, location, trunkID)
	require.NoError(t, err)

	return attrs
}

// GetDialogflowV2SipTrunkAttrsE returns the settings Google Cloud holds for the given Dialogflow ES SIP trunk.
// The ctx parameter supports cancellation and timeouts.
func GetDialogflowV2SipTrunkAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, trunkID string) (*dialogflowv2.GoogleCloudDialogflowV2SipTrunk, error) {
	logger.Default.Logf(t, "Getting settings for Dialogflow ES SIP trunk %s in %s in project %s", trunkID, location, projectID)

	service, err := NewDialogflowV2ServiceE(t, ctx, location)
	if err != nil {
		return nil, err
	}

	return GetDialogflowV2SipTrunkAttrsWithClient(ctx, service, projectID, location, trunkID)
}

// GetDialogflowV2SipTrunkAttrsWithClient returns the settings Google Cloud holds for the given Dialogflow ES SIP trunk using the supplied
// *dialogflowv2.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see dialogflowv2_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetDialogflowV2SipTrunkAttrsWithClient(ctx context.Context, service *dialogflowv2.Service, projectID string, location string, trunkID string) (*dialogflowv2.GoogleCloudDialogflowV2SipTrunk, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/sipTrunks/%s", projectID, location, trunkID)

	attrs, err := service.Projects.Locations.SipTrunks.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Dialogflow ES SIP trunk %s in %s in project %s does not exist", trunkID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Dialogflow ES SIP trunk %s in %s in project %s: %w", trunkID, location, projectID, err)
	}

	return attrs, nil
}

// GetDialogflowV2VersionAttrs returns the settings Google Cloud holds for the given Dialogflow ES agent version, so a test can assert on what was
// actually created rather than only that it exists.
// A version is a frozen copy of the agent an environment can serve, so its description and its status are what say whether it is usable.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetDialogflowV2VersionAttrs(t testing.TestingT, ctx context.Context, projectID string, versionID string) *dialogflowv2.GoogleCloudDialogflowV2Version {
	attrs, err := GetDialogflowV2VersionAttrsE(t, ctx, projectID, versionID)
	require.NoError(t, err)

	return attrs
}

// GetDialogflowV2VersionAttrsE returns the settings Google Cloud holds for the given Dialogflow ES agent version.
// The ctx parameter supports cancellation and timeouts.
func GetDialogflowV2VersionAttrsE(t testing.TestingT, ctx context.Context, projectID string, versionID string) (*dialogflowv2.GoogleCloudDialogflowV2Version, error) {
	logger.Default.Logf(t, "Getting settings for Dialogflow ES agent version %s in project %s", versionID, projectID)

	service, err := NewDialogflowV2ServiceE(t, ctx, "global")
	if err != nil {
		return nil, err
	}

	return GetDialogflowV2VersionAttrsWithClient(ctx, service, projectID, versionID)
}

// GetDialogflowV2VersionAttrsWithClient returns the settings Google Cloud holds for the given Dialogflow ES agent version using the supplied
// *dialogflowv2.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see dialogflowv2_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetDialogflowV2VersionAttrsWithClient(ctx context.Context, service *dialogflowv2.Service, projectID string, versionID string) (*dialogflowv2.GoogleCloudDialogflowV2Version, error) {
	name := fmt.Sprintf("projects/%s/agent/versions/%s", projectID, versionID)

	attrs, err := service.Projects.Agent.Versions.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Dialogflow ES agent version %s in project %s does not exist", versionID, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Dialogflow ES agent version %s in project %s: %w", versionID, projectID, err)
	}

	return attrs, nil
}

// NewDialogflowV2ServiceE creates a Dialogflow ES service authenticated the same way every other client
// in this module is. Dialogflow answers on a per-location host for every location but `global`, so the
// location decides which endpoint the service talks to, and a location that is not a plain identifier is
// refused before it can send a request, and credentials, to a host of the caller's choosing.
// The ctx parameter supports cancellation and timeouts.
func NewDialogflowV2ServiceE(t testing.TestingT, ctx context.Context, location string) (*dialogflowv2.Service, error) {
	if !dialogflowLocationPattern.MatchString(location) {
		return nil, fmt.Errorf("%q is not a valid location: a location may hold only lowercase letters, digits and hyphens", location)
	}

	opts := append(withOptions(), option.WithScopes(dialogflowv2.CloudPlatformScope))
	if location != "global" {
		opts = append(opts, option.WithEndpoint(fmt.Sprintf("https://%s-dialogflow.googleapis.com/", location)))
	}

	return dialogflowv2.NewService(ctx, opts...)
}
