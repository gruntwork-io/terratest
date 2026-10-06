package gcp

import (
	"context"
	"errors"
	"fmt"

	"github.com/gruntwork-io/terratest/modules/core/v2/logger"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/discoveryengine/v1"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
)

// Discovery Engine keeps most resources twice: once under a collection and once directly under a
// location. Terraform writes to the collection form, where the collection is almost always
// `default_collection`, so these reads take the collection as a parameter rather than assuming it.

// GetDiscoveryEngineACLConfigAttrs returns the settings Google Cloud holds for the given Discovery Engine ACL config, so a test can assert on what was
// actually created rather than only that it exists.
// The config names the identity provider a data store checks documents against, so without it nothing is access controlled.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetDiscoveryEngineACLConfigAttrs(t testing.TestingT, ctx context.Context, projectID string, location string) *discoveryengine.GoogleCloudDiscoveryengineV1AclConfig {
	attrs, err := GetDiscoveryEngineACLConfigAttrsE(t, ctx, projectID, location)
	require.NoError(t, err)

	return attrs
}

// GetDiscoveryEngineACLConfigAttrsE returns the settings Google Cloud holds for the given Discovery Engine ACL config.
// The ctx parameter supports cancellation and timeouts.
func GetDiscoveryEngineACLConfigAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string) (*discoveryengine.GoogleCloudDiscoveryengineV1AclConfig, error) {
	logger.Default.Logf(t, "Getting settings for Discovery Engine ACL config in %s in project %s", location, projectID)

	service, err := NewDiscoveryEngineServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetDiscoveryEngineACLConfigAttrsWithClient(ctx, service, projectID, location)
}

// GetDiscoveryEngineACLConfigAttrsWithClient returns the settings Google Cloud holds for the given Discovery Engine ACL config using the supplied
// *discoveryengine.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see discoveryengine_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetDiscoveryEngineACLConfigAttrsWithClient(ctx context.Context, service *discoveryengine.Service, projectID string, location string) (*discoveryengine.GoogleCloudDiscoveryengineV1AclConfig, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/aclConfig", projectID, location)

	attrs, err := service.Projects.Locations.GetAclConfig(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Discovery Engine ACL config in %s in project %s does not exist", location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Discovery Engine ACL config in %s in project %s: %w", location, projectID, err)
	}

	return attrs, nil
}

// GetDiscoveryEngineCmekConfigAttrs returns the settings Google Cloud holds for the given Discovery Engine CMEK config, so a test can assert on what was
// actually created rather than only that it exists.
// The config names the key a data store's indexes are encrypted with, so which key it is and whether it is the default are the point.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetDiscoveryEngineCmekConfigAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, configID string) *discoveryengine.GoogleCloudDiscoveryengineV1CmekConfig {
	attrs, err := GetDiscoveryEngineCmekConfigAttrsE(t, ctx, projectID, location, configID)
	require.NoError(t, err)

	return attrs
}

// GetDiscoveryEngineCmekConfigAttrsE returns the settings Google Cloud holds for the given Discovery Engine CMEK config.
// The ctx parameter supports cancellation and timeouts.
func GetDiscoveryEngineCmekConfigAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, configID string) (*discoveryengine.GoogleCloudDiscoveryengineV1CmekConfig, error) {
	logger.Default.Logf(t, "Getting settings for Discovery Engine CMEK config %s in %s in project %s", configID, location, projectID)

	service, err := NewDiscoveryEngineServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetDiscoveryEngineCmekConfigAttrsWithClient(ctx, service, projectID, location, configID)
}

// GetDiscoveryEngineCmekConfigAttrsWithClient returns the settings Google Cloud holds for the given Discovery Engine CMEK config using the supplied
// *discoveryengine.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see discoveryengine_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetDiscoveryEngineCmekConfigAttrsWithClient(ctx context.Context, service *discoveryengine.Service, projectID string, location string, configID string) (*discoveryengine.GoogleCloudDiscoveryengineV1CmekConfig, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/cmekConfigs/%s", projectID, location, configID)

	attrs, err := service.Projects.Locations.CmekConfigs.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Discovery Engine CMEK config %s in %s in project %s does not exist", configID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Discovery Engine CMEK config %s in %s in project %s: %w", configID, location, projectID, err)
	}

	return attrs, nil
}

// GetDiscoveryEngineDataStoreAttrs returns the settings Google Cloud holds for the given Discovery Engine data store, so a test can assert on what was
// actually created rather than only that it exists.
// A data store is what documents are indexed into, so its industry vertical, its content config and the solutions it serves decide what can search it.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetDiscoveryEngineDataStoreAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, collectionID string, dataStoreID string) *discoveryengine.GoogleCloudDiscoveryengineV1DataStore {
	attrs, err := GetDiscoveryEngineDataStoreAttrsE(t, ctx, projectID, location, collectionID, dataStoreID)
	require.NoError(t, err)

	return attrs
}

// GetDiscoveryEngineDataStoreAttrsE returns the settings Google Cloud holds for the given Discovery Engine data store.
// The ctx parameter supports cancellation and timeouts.
func GetDiscoveryEngineDataStoreAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, collectionID string, dataStoreID string) (*discoveryengine.GoogleCloudDiscoveryengineV1DataStore, error) {
	logger.Default.Logf(t, "Getting settings for Discovery Engine data store %s in collection %s in %s in project %s", dataStoreID, collectionID, location, projectID)

	service, err := NewDiscoveryEngineServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetDiscoveryEngineDataStoreAttrsWithClient(ctx, service, projectID, location, collectionID, dataStoreID)
}

// GetDiscoveryEngineDataStoreAttrsWithClient returns the settings Google Cloud holds for the given Discovery Engine data store using the supplied
// *discoveryengine.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see discoveryengine_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetDiscoveryEngineDataStoreAttrsWithClient(ctx context.Context, service *discoveryengine.Service, projectID string, location string, collectionID string, dataStoreID string) (*discoveryengine.GoogleCloudDiscoveryengineV1DataStore, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/collections/%s/dataStores/%s", projectID, location, collectionID, dataStoreID)

	attrs, err := service.Projects.Locations.Collections.DataStores.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Discovery Engine data store %s in collection %s in %s in project %s does not exist", dataStoreID, collectionID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Discovery Engine data store %s in collection %s in %s in project %s: %w", dataStoreID, collectionID, location, projectID, err)
	}

	return attrs, nil
}

// GetDiscoveryEngineEngineAttrs returns the settings Google Cloud holds for the given Discovery Engine engine, so a test can assert on what was
// actually created rather than only that it exists.
// One read serves the search, chat and recommendation engine modules, because Google stores all three as an engine and the solution type says which it is.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetDiscoveryEngineEngineAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, collectionID string, engineID string) *discoveryengine.GoogleCloudDiscoveryengineV1Engine {
	attrs, err := GetDiscoveryEngineEngineAttrsE(t, ctx, projectID, location, collectionID, engineID)
	require.NoError(t, err)

	return attrs
}

// GetDiscoveryEngineEngineAttrsE returns the settings Google Cloud holds for the given Discovery Engine engine.
// The ctx parameter supports cancellation and timeouts.
func GetDiscoveryEngineEngineAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, collectionID string, engineID string) (*discoveryengine.GoogleCloudDiscoveryengineV1Engine, error) {
	logger.Default.Logf(t, "Getting settings for Discovery Engine engine %s in collection %s in %s in project %s", engineID, collectionID, location, projectID)

	service, err := NewDiscoveryEngineServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetDiscoveryEngineEngineAttrsWithClient(ctx, service, projectID, location, collectionID, engineID)
}

// GetDiscoveryEngineEngineAttrsWithClient returns the settings Google Cloud holds for the given Discovery Engine engine using the supplied
// *discoveryengine.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see discoveryengine_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetDiscoveryEngineEngineAttrsWithClient(ctx context.Context, service *discoveryengine.Service, projectID string, location string, collectionID string, engineID string) (*discoveryengine.GoogleCloudDiscoveryengineV1Engine, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/collections/%s/engines/%s", projectID, location, collectionID, engineID)

	attrs, err := service.Projects.Locations.Collections.Engines.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Discovery Engine engine %s in collection %s in %s in project %s does not exist", engineID, collectionID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Discovery Engine engine %s in collection %s in %s in project %s: %w", engineID, collectionID, location, projectID, err)
	}

	return attrs, nil
}

// GetDiscoveryEngineEngineIamPolicyAttrs returns the IAM policy Google Cloud holds for the given Discovery Engine engine, so a test can assert on what was
// actually created rather than only that it exists.
// Who may query the engine is what the policy decides, and an engine Google created on its own carries no binding at all.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetDiscoveryEngineEngineIamPolicyAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, collectionID string, engineID string) *discoveryengine.GoogleIamV1Policy {
	policy, err := GetDiscoveryEngineEngineIamPolicyAttrsE(t, ctx, projectID, location, collectionID, engineID)
	require.NoError(t, err)

	return policy
}

// GetDiscoveryEngineEngineIamPolicyAttrsE returns the IAM policy Google Cloud holds for the given Discovery Engine engine.
// The ctx parameter supports cancellation and timeouts.
func GetDiscoveryEngineEngineIamPolicyAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, collectionID string, engineID string) (*discoveryengine.GoogleIamV1Policy, error) {
	logger.Default.Logf(t, "Getting the IAM policy for Discovery Engine engine %s in collection %s in %s in project %s", engineID, collectionID, location, projectID)

	service, err := NewDiscoveryEngineServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetDiscoveryEngineEngineIamPolicyAttrsWithClient(ctx, service, projectID, location, collectionID, engineID)
}

// GetDiscoveryEngineEngineIamPolicyAttrsWithClient returns the IAM policy Google Cloud holds for the given Discovery Engine engine using the supplied
// *discoveryengine.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see discoveryengine_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetDiscoveryEngineEngineIamPolicyAttrsWithClient(ctx context.Context, service *discoveryengine.Service, projectID string, location string, collectionID string, engineID string) (*discoveryengine.GoogleIamV1Policy, error) {
	resource := fmt.Sprintf("projects/%s/locations/%s/collections/%s/engines/%s", projectID, location, collectionID, engineID)

	policy, err := service.Projects.Locations.Collections.Engines.GetIamPolicy(resource).OptionsRequestedPolicyVersion(iamPolicyVersionWithConditions).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Discovery Engine engine %s in collection %s in %s in project %s does not exist", engineID, collectionID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get the IAM policy for Discovery Engine engine %s in collection %s in %s in project %s: %w", engineID, collectionID, location, projectID, err)
	}

	return policy, nil
}

// GetDiscoveryEngineAssistantAttrs returns the settings Google Cloud holds for the given Discovery Engine assistant, so a test can assert on what was
// actually created rather than only that it exists.
// An assistant is the agent an engine answers through, so its display name and the tools it is allowed are what it can do.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetDiscoveryEngineAssistantAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, collectionID string, engineID string, assistantID string) *discoveryengine.GoogleCloudDiscoveryengineV1Assistant {
	attrs, err := GetDiscoveryEngineAssistantAttrsE(t, ctx, projectID, location, collectionID, engineID, assistantID)
	require.NoError(t, err)

	return attrs
}

// GetDiscoveryEngineAssistantAttrsE returns the settings Google Cloud holds for the given Discovery Engine assistant.
// The ctx parameter supports cancellation and timeouts.
func GetDiscoveryEngineAssistantAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, collectionID string, engineID string, assistantID string) (*discoveryengine.GoogleCloudDiscoveryengineV1Assistant, error) {
	logger.Default.Logf(t, "Getting settings for Discovery Engine assistant %s on engine %s in collection %s in %s in project %s", assistantID, engineID, collectionID, location, projectID)

	service, err := NewDiscoveryEngineServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetDiscoveryEngineAssistantAttrsWithClient(ctx, service, projectID, location, collectionID, engineID, assistantID)
}

// GetDiscoveryEngineAssistantAttrsWithClient returns the settings Google Cloud holds for the given Discovery Engine assistant using the supplied
// *discoveryengine.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see discoveryengine_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetDiscoveryEngineAssistantAttrsWithClient(ctx context.Context, service *discoveryengine.Service, projectID string, location string, collectionID string, engineID string, assistantID string) (*discoveryengine.GoogleCloudDiscoveryengineV1Assistant, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/collections/%s/engines/%s/assistants/%s", projectID, location, collectionID, engineID, assistantID)

	attrs, err := service.Projects.Locations.Collections.Engines.Assistants.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Discovery Engine assistant %s on engine %s in collection %s in %s in project %s does not exist", assistantID, engineID, collectionID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Discovery Engine assistant %s on engine %s in collection %s in %s in project %s: %w", assistantID, engineID, collectionID, location, projectID, err)
	}

	return attrs, nil
}

// GetDiscoveryEngineDataConnectorAttrs returns the settings Google Cloud holds for the given Discovery Engine data connector, so a test can assert on what was
// actually created rather than only that it exists.
// The connector is what pulls documents in from another system, so the source it names and its refresh interval decide what gets indexed and how often.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetDiscoveryEngineDataConnectorAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, collectionID string) *discoveryengine.GoogleCloudDiscoveryengineV1DataConnector {
	attrs, err := GetDiscoveryEngineDataConnectorAttrsE(t, ctx, projectID, location, collectionID)
	require.NoError(t, err)

	return attrs
}

// GetDiscoveryEngineDataConnectorAttrsE returns the settings Google Cloud holds for the given Discovery Engine data connector.
// The ctx parameter supports cancellation and timeouts.
func GetDiscoveryEngineDataConnectorAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, collectionID string) (*discoveryengine.GoogleCloudDiscoveryengineV1DataConnector, error) {
	logger.Default.Logf(t, "Getting settings for Discovery Engine data connector in collection %s in %s in project %s", collectionID, location, projectID)

	service, err := NewDiscoveryEngineServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetDiscoveryEngineDataConnectorAttrsWithClient(ctx, service, projectID, location, collectionID)
}

// GetDiscoveryEngineDataConnectorAttrsWithClient returns the settings Google Cloud holds for the given Discovery Engine data connector using the supplied
// *discoveryengine.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see discoveryengine_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetDiscoveryEngineDataConnectorAttrsWithClient(ctx context.Context, service *discoveryengine.Service, projectID string, location string, collectionID string) (*discoveryengine.GoogleCloudDiscoveryengineV1DataConnector, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/collections/%s/dataConnector", projectID, location, collectionID)

	attrs, err := service.Projects.Locations.Collections.GetDataConnector(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Discovery Engine data connector in collection %s in %s in project %s does not exist", collectionID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Discovery Engine data connector in collection %s in %s in project %s: %w", collectionID, location, projectID, err)
	}

	return attrs, nil
}

// GetDiscoveryEngineSchemaAttrs returns the settings Google Cloud holds for the given Discovery Engine schema, so a test can assert on what was
// actually created rather than only that it exists.
// A schema says which document fields are searchable, so what it declares decides what a query can filter on.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetDiscoveryEngineSchemaAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, collectionID string, dataStoreID string, schemaID string) *discoveryengine.GoogleCloudDiscoveryengineV1Schema {
	attrs, err := GetDiscoveryEngineSchemaAttrsE(t, ctx, projectID, location, collectionID, dataStoreID, schemaID)
	require.NoError(t, err)

	return attrs
}

// GetDiscoveryEngineSchemaAttrsE returns the settings Google Cloud holds for the given Discovery Engine schema.
// The ctx parameter supports cancellation and timeouts.
func GetDiscoveryEngineSchemaAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, collectionID string, dataStoreID string, schemaID string) (*discoveryengine.GoogleCloudDiscoveryengineV1Schema, error) {
	logger.Default.Logf(t, "Getting settings for Discovery Engine schema %s in data store %s in collection %s in %s in project %s", schemaID, dataStoreID, collectionID, location, projectID)

	service, err := NewDiscoveryEngineServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetDiscoveryEngineSchemaAttrsWithClient(ctx, service, projectID, location, collectionID, dataStoreID, schemaID)
}

// GetDiscoveryEngineSchemaAttrsWithClient returns the settings Google Cloud holds for the given Discovery Engine schema using the supplied
// *discoveryengine.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see discoveryengine_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetDiscoveryEngineSchemaAttrsWithClient(ctx context.Context, service *discoveryengine.Service, projectID string, location string, collectionID string, dataStoreID string, schemaID string) (*discoveryengine.GoogleCloudDiscoveryengineV1Schema, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/collections/%s/dataStores/%s/schemas/%s", projectID, location, collectionID, dataStoreID, schemaID)

	attrs, err := service.Projects.Locations.Collections.DataStores.Schemas.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Discovery Engine schema %s in data store %s in collection %s in %s in project %s does not exist", schemaID, dataStoreID, collectionID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Discovery Engine schema %s in data store %s in collection %s in %s in project %s: %w", schemaID, dataStoreID, collectionID, location, projectID, err)
	}

	return attrs, nil
}

// GetDiscoveryEngineServingConfigAttrs returns the settings Google Cloud holds for the given Discovery Engine serving config, so a test can assert on what was
// actually created rather than only that it exists.
// A serving config is the set of knobs one search surface uses, so which controls it applies decides how results are ordered.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetDiscoveryEngineServingConfigAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, collectionID string, dataStoreID string, configID string) *discoveryengine.GoogleCloudDiscoveryengineV1ServingConfig {
	attrs, err := GetDiscoveryEngineServingConfigAttrsE(t, ctx, projectID, location, collectionID, dataStoreID, configID)
	require.NoError(t, err)

	return attrs
}

// GetDiscoveryEngineServingConfigAttrsE returns the settings Google Cloud holds for the given Discovery Engine serving config.
// The ctx parameter supports cancellation and timeouts.
func GetDiscoveryEngineServingConfigAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, collectionID string, dataStoreID string, configID string) (*discoveryengine.GoogleCloudDiscoveryengineV1ServingConfig, error) {
	logger.Default.Logf(t, "Getting settings for Discovery Engine serving config %s in data store %s in collection %s in %s in project %s", configID, dataStoreID, collectionID, location, projectID)

	service, err := NewDiscoveryEngineServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetDiscoveryEngineServingConfigAttrsWithClient(ctx, service, projectID, location, collectionID, dataStoreID, configID)
}

// GetDiscoveryEngineServingConfigAttrsWithClient returns the settings Google Cloud holds for the given Discovery Engine serving config using the supplied
// *discoveryengine.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see discoveryengine_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetDiscoveryEngineServingConfigAttrsWithClient(ctx context.Context, service *discoveryengine.Service, projectID string, location string, collectionID string, dataStoreID string, configID string) (*discoveryengine.GoogleCloudDiscoveryengineV1ServingConfig, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/collections/%s/dataStores/%s/servingConfigs/%s", projectID, location, collectionID, dataStoreID, configID)

	attrs, err := service.Projects.Locations.Collections.DataStores.ServingConfigs.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Discovery Engine serving config %s in data store %s in collection %s in %s in project %s does not exist", configID, dataStoreID, collectionID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Discovery Engine serving config %s in data store %s in collection %s in %s in project %s: %w", configID, dataStoreID, collectionID, location, projectID, err)
	}

	return attrs, nil
}

// GetDiscoveryEngineControlAttrs returns the settings Google Cloud holds for the given Discovery Engine control, so a test can assert on what was
// actually created rather than only that it exists.
// A control is a rule that boosts, filters or redirects results, so its condition and its action are the whole point.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetDiscoveryEngineControlAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, collectionID string, dataStoreID string, controlID string) *discoveryengine.GoogleCloudDiscoveryengineV1Control {
	attrs, err := GetDiscoveryEngineControlAttrsE(t, ctx, projectID, location, collectionID, dataStoreID, controlID)
	require.NoError(t, err)

	return attrs
}

// GetDiscoveryEngineControlAttrsE returns the settings Google Cloud holds for the given Discovery Engine control.
// The ctx parameter supports cancellation and timeouts.
func GetDiscoveryEngineControlAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, collectionID string, dataStoreID string, controlID string) (*discoveryengine.GoogleCloudDiscoveryengineV1Control, error) {
	logger.Default.Logf(t, "Getting settings for Discovery Engine control %s in data store %s in collection %s in %s in project %s", controlID, dataStoreID, collectionID, location, projectID)

	service, err := NewDiscoveryEngineServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetDiscoveryEngineControlAttrsWithClient(ctx, service, projectID, location, collectionID, dataStoreID, controlID)
}

// GetDiscoveryEngineControlAttrsWithClient returns the settings Google Cloud holds for the given Discovery Engine control using the supplied
// *discoveryengine.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see discoveryengine_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetDiscoveryEngineControlAttrsWithClient(ctx context.Context, service *discoveryengine.Service, projectID string, location string, collectionID string, dataStoreID string, controlID string) (*discoveryengine.GoogleCloudDiscoveryengineV1Control, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/collections/%s/dataStores/%s/controls/%s", projectID, location, collectionID, dataStoreID, controlID)

	attrs, err := service.Projects.Locations.Collections.DataStores.Controls.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Discovery Engine control %s in data store %s in collection %s in %s in project %s does not exist", controlID, dataStoreID, collectionID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Discovery Engine control %s in data store %s in collection %s in %s in project %s: %w", controlID, dataStoreID, collectionID, location, projectID, err)
	}

	return attrs, nil
}

// GetDiscoveryEngineTargetSiteAttrs returns the settings Google Cloud holds for the given Discovery Engine target site, so a test can assert on what was
// actually created rather than only that it exists.
// A target site is a pattern of URLs a website data store may crawl, so the pattern and whether it includes or excludes are the point.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetDiscoveryEngineTargetSiteAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, collectionID string, dataStoreID string, siteID string) *discoveryengine.GoogleCloudDiscoveryengineV1TargetSite {
	attrs, err := GetDiscoveryEngineTargetSiteAttrsE(t, ctx, projectID, location, collectionID, dataStoreID, siteID)
	require.NoError(t, err)

	return attrs
}

// GetDiscoveryEngineTargetSiteAttrsE returns the settings Google Cloud holds for the given Discovery Engine target site.
// The ctx parameter supports cancellation and timeouts.
func GetDiscoveryEngineTargetSiteAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, collectionID string, dataStoreID string, siteID string) (*discoveryengine.GoogleCloudDiscoveryengineV1TargetSite, error) {
	logger.Default.Logf(t, "Getting settings for Discovery Engine target site %s in data store %s in collection %s in %s in project %s", siteID, dataStoreID, collectionID, location, projectID)

	service, err := NewDiscoveryEngineServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetDiscoveryEngineTargetSiteAttrsWithClient(ctx, service, projectID, location, collectionID, dataStoreID, siteID)
}

// GetDiscoveryEngineTargetSiteAttrsWithClient returns the settings Google Cloud holds for the given Discovery Engine target site using the supplied
// *discoveryengine.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see discoveryengine_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetDiscoveryEngineTargetSiteAttrsWithClient(ctx context.Context, service *discoveryengine.Service, projectID string, location string, collectionID string, dataStoreID string, siteID string) (*discoveryengine.GoogleCloudDiscoveryengineV1TargetSite, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/collections/%s/dataStores/%s/siteSearchEngine/targetSites/%s", projectID, location, collectionID, dataStoreID, siteID)

	attrs, err := service.Projects.Locations.Collections.DataStores.SiteSearchEngine.TargetSites.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Discovery Engine target site %s in data store %s in collection %s in %s in project %s does not exist", siteID, dataStoreID, collectionID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Discovery Engine target site %s in data store %s in collection %s in %s in project %s: %w", siteID, dataStoreID, collectionID, location, projectID, err)
	}

	return attrs, nil
}

// GetDiscoveryEngineWidgetConfigAttrs returns the settings Google Cloud holds for the given Discovery Engine widget config, so a test can assert on what was
// actually created rather than only that it exists.
// A widget config is the embeddable search box, so whether results are public and which features are on decide what a visitor sees.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetDiscoveryEngineWidgetConfigAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, collectionID string, dataStoreID string, configID string) *discoveryengine.GoogleCloudDiscoveryengineV1WidgetConfig {
	attrs, err := GetDiscoveryEngineWidgetConfigAttrsE(t, ctx, projectID, location, collectionID, dataStoreID, configID)
	require.NoError(t, err)

	return attrs
}

// GetDiscoveryEngineWidgetConfigAttrsE returns the settings Google Cloud holds for the given Discovery Engine widget config.
// The ctx parameter supports cancellation and timeouts.
func GetDiscoveryEngineWidgetConfigAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, collectionID string, dataStoreID string, configID string) (*discoveryengine.GoogleCloudDiscoveryengineV1WidgetConfig, error) {
	logger.Default.Logf(t, "Getting settings for Discovery Engine widget config %s in data store %s in collection %s in %s in project %s", configID, dataStoreID, collectionID, location, projectID)

	service, err := NewDiscoveryEngineServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetDiscoveryEngineWidgetConfigAttrsWithClient(ctx, service, projectID, location, collectionID, dataStoreID, configID)
}

// GetDiscoveryEngineWidgetConfigAttrsWithClient returns the settings Google Cloud holds for the given Discovery Engine widget config using the supplied
// *discoveryengine.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see discoveryengine_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetDiscoveryEngineWidgetConfigAttrsWithClient(ctx context.Context, service *discoveryengine.Service, projectID string, location string, collectionID string, dataStoreID string, configID string) (*discoveryengine.GoogleCloudDiscoveryengineV1WidgetConfig, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/collections/%s/dataStores/%s/widgetConfigs/%s", projectID, location, collectionID, dataStoreID, configID)

	attrs, err := service.Projects.Locations.Collections.DataStores.WidgetConfigs.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Discovery Engine widget config %s in data store %s in collection %s in %s in project %s does not exist", configID, dataStoreID, collectionID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Discovery Engine widget config %s in data store %s in collection %s in %s in project %s: %w", configID, dataStoreID, collectionID, location, projectID, err)
	}

	return attrs, nil
}

// GetDiscoveryEngineLicenseConfigAttrs returns the settings Google Cloud holds for the given Discovery Engine license config, so a test can assert on what was
// actually created rather than only that it exists.
// A licence config is how many seats of a paid tier are bought, so the count and its state are what is billed.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetDiscoveryEngineLicenseConfigAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, configID string) *discoveryengine.GoogleCloudDiscoveryengineV1LicenseConfig {
	attrs, err := GetDiscoveryEngineLicenseConfigAttrsE(t, ctx, projectID, location, configID)
	require.NoError(t, err)

	return attrs
}

// GetDiscoveryEngineLicenseConfigAttrsE returns the settings Google Cloud holds for the given Discovery Engine license config.
// The ctx parameter supports cancellation and timeouts.
func GetDiscoveryEngineLicenseConfigAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, configID string) (*discoveryengine.GoogleCloudDiscoveryengineV1LicenseConfig, error) {
	logger.Default.Logf(t, "Getting settings for Discovery Engine license config %s in %s in project %s", configID, location, projectID)

	service, err := NewDiscoveryEngineServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetDiscoveryEngineLicenseConfigAttrsWithClient(ctx, service, projectID, location, configID)
}

// GetDiscoveryEngineLicenseConfigAttrsWithClient returns the settings Google Cloud holds for the given Discovery Engine license config using the supplied
// *discoveryengine.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see discoveryengine_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetDiscoveryEngineLicenseConfigAttrsWithClient(ctx context.Context, service *discoveryengine.Service, projectID string, location string, configID string) (*discoveryengine.GoogleCloudDiscoveryengineV1LicenseConfig, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/licenseConfigs/%s", projectID, location, configID)

	attrs, err := service.Projects.Locations.LicenseConfigs.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Discovery Engine license config %s in %s in project %s does not exist", configID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Discovery Engine license config %s in %s in project %s: %w", configID, location, projectID, err)
	}

	return attrs, nil
}

// GetDiscoveryEngineUserStoreAttrs returns the settings Google Cloud holds for the given Discovery Engine user store, so a test can assert on what was
// actually created rather than only that it exists.
// A user store is where licensed users are recorded, so the licence config it points at is what it draws seats from.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetDiscoveryEngineUserStoreAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, storeID string) *discoveryengine.GoogleCloudDiscoveryengineV1UserStore {
	attrs, err := GetDiscoveryEngineUserStoreAttrsE(t, ctx, projectID, location, storeID)
	require.NoError(t, err)

	return attrs
}

// GetDiscoveryEngineUserStoreAttrsE returns the settings Google Cloud holds for the given Discovery Engine user store.
// The ctx parameter supports cancellation and timeouts.
func GetDiscoveryEngineUserStoreAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, storeID string) (*discoveryengine.GoogleCloudDiscoveryengineV1UserStore, error) {
	logger.Default.Logf(t, "Getting settings for Discovery Engine user store %s in %s in project %s", storeID, location, projectID)

	service, err := NewDiscoveryEngineServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetDiscoveryEngineUserStoreAttrsWithClient(ctx, service, projectID, location, storeID)
}

// GetDiscoveryEngineUserStoreAttrsWithClient returns the settings Google Cloud holds for the given Discovery Engine user store using the supplied
// *discoveryengine.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see discoveryengine_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetDiscoveryEngineUserStoreAttrsWithClient(ctx context.Context, service *discoveryengine.Service, projectID string, location string, storeID string) (*discoveryengine.GoogleCloudDiscoveryengineV1UserStore, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/userStores/%s", projectID, location, storeID)

	attrs, err := service.Projects.Locations.UserStores.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Discovery Engine user store %s in %s in project %s does not exist", storeID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Discovery Engine user store %s in %s in project %s: %w", storeID, location, projectID, err)
	}

	return attrs, nil
}

// NewDiscoveryEngineServiceE creates a Discovery Engine service authenticated the same way every other
// client in this module is.
// The ctx parameter supports cancellation and timeouts.
func NewDiscoveryEngineServiceE(t testing.TestingT, ctx context.Context) (*discoveryengine.Service, error) {
	return discoveryengine.NewService(ctx, append(withOptions(), option.WithScopes(discoveryengine.CloudPlatformScope))...)
}
