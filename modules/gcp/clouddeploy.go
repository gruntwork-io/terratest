package gcp

import (
	"context"
	"errors"
	"fmt"

	"github.com/gruntwork-io/terratest/modules/core/v2/logger"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/clouddeploy/v1"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
)

// GetDeliveryPipelineAttrs returns the settings Google Cloud holds for the given Cloud Deploy
// delivery pipeline, so a test can assert on what was actually created rather than only that it
// exists. A pipeline is regional, so the location it was created in has to be given.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetDeliveryPipelineAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, pipelineName string) *clouddeploy.DeliveryPipeline {
	pipeline, err := GetDeliveryPipelineAttrsE(t, ctx, projectID, location, pipelineName)
	require.NoError(t, err)

	return pipeline
}

// GetDeliveryPipelineAttrsE returns the settings Google Cloud holds for the given Cloud Deploy
// delivery pipeline.
// The ctx parameter supports cancellation and timeouts.
func GetDeliveryPipelineAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, pipelineName string) (*clouddeploy.DeliveryPipeline, error) {
	logger.Default.Logf(t, "Getting settings for delivery pipeline %s in project %s", pipelineName, projectID)

	service, err := NewCloudDeployServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetDeliveryPipelineAttrsWithClient(ctx, service, projectID, location, pipelineName)
}

// GetDeliveryPipelineAttrsWithClient returns the settings Google Cloud holds for the given Cloud
// Deploy delivery pipeline using the supplied *clouddeploy.Service. Prefer this variant in unit
// tests where the service is backed by an httptest fake server (see clouddeploy_test.go for the
// pattern).
// The ctx parameter supports cancellation and timeouts.
func GetDeliveryPipelineAttrsWithClient(ctx context.Context, service *clouddeploy.Service, projectID string, location string, pipelineName string) (*clouddeploy.DeliveryPipeline, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/deliveryPipelines/%s", projectID, location, pipelineName)

	pipeline, err := service.Projects.Locations.DeliveryPipelines.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("delivery pipeline %s does not exist in project %s", pipelineName, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for delivery pipeline %s in project %s: %w", pipelineName, projectID, err)
	}

	return pipeline, nil
}

// GetDeployTargetAttrs returns the settings Google Cloud holds for the given Cloud Deploy target, so
// a test can assert on what was actually created rather than only that it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetDeployTargetAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, targetName string) *clouddeploy.Target {
	target, err := GetDeployTargetAttrsE(t, ctx, projectID, location, targetName)
	require.NoError(t, err)

	return target
}

// GetDeployTargetAttrsE returns the settings Google Cloud holds for the given Cloud Deploy target.
// The ctx parameter supports cancellation and timeouts.
func GetDeployTargetAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, targetName string) (*clouddeploy.Target, error) {
	logger.Default.Logf(t, "Getting settings for Cloud Deploy target %s in location %s in project %s", targetName, location, projectID)

	service, err := NewCloudDeployServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetDeployTargetAttrsWithClient(ctx, service, projectID, location, targetName)
}

// GetDeployTargetAttrsWithClient returns the settings Google Cloud holds for the given Cloud Deploy
// target using the supplied *clouddeploy.Service. Prefer this variant in unit tests where the
// service is backed by an httptest fake server (see clouddeploy_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetDeployTargetAttrsWithClient(ctx context.Context, service *clouddeploy.Service, projectID string, location string, targetName string) (*clouddeploy.Target, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/targets/%s", projectID, location, targetName)

	target, err := service.Projects.Locations.Targets.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Cloud Deploy target %s does not exist in location %s in project %s", targetName, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Cloud Deploy target %s in location %s in project %s: %w", targetName, location, projectID, err)
	}

	return target, nil
}

// NewCloudDeployServiceE creates a Cloud Deploy service authenticated the same way every other
// client in this module is.
// The ctx parameter supports cancellation and timeouts.
func NewCloudDeployServiceE(t testing.TestingT, ctx context.Context) (*clouddeploy.Service, error) {
	return clouddeploy.NewService(ctx, append(withOptions(), option.WithScopes(clouddeploy.CloudPlatformScope))...)
}

// GetCloudDeployAutomationAttrs returns the settings Google Cloud holds for the given delivery pipeline automation, so a test can assert on
// what was actually created rather than only that it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetCloudDeployAutomationAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, pipelineID string, automationID string) *clouddeploy.Automation {
	automation, err := GetCloudDeployAutomationAttrsE(t, ctx, projectID, location, pipelineID, automationID)
	require.NoError(t, err)

	return automation
}

// GetCloudDeployAutomationAttrsE returns the settings Google Cloud holds for the given delivery pipeline automation.
// The ctx parameter supports cancellation and timeouts.
func GetCloudDeployAutomationAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, pipelineID string, automationID string) (*clouddeploy.Automation, error) {
	logger.Default.Logf(t, "Getting settings for automation %s on pipeline %s in %s in project %s", automationID, pipelineID, location, projectID)

	service, err := NewCloudDeployServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetCloudDeployAutomationAttrsWithClient(ctx, service, projectID, location, pipelineID, automationID)
}

// GetCloudDeployAutomationAttrsWithClient returns the settings Google Cloud holds for the given delivery pipeline automation using the supplied
// *clouddeploy.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see clouddeploy_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetCloudDeployAutomationAttrsWithClient(ctx context.Context, service *clouddeploy.Service, projectID string, location string, pipelineID string, automationID string) (*clouddeploy.Automation, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/deliveryPipelines/%s/automations/%s", projectID, location, pipelineID, automationID)

	automation, err := service.Projects.Locations.DeliveryPipelines.Automations.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the automation %s does not exist on pipeline %s in %s in project %s", automationID, pipelineID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for automation %s on pipeline %s in %s in project %s: %w", automationID, pipelineID, location, projectID, err)
	}

	return automation, nil
}

// GetCloudDeployCustomTargetTypeAttrs returns the settings Google Cloud holds for the given custom target type, so a test can assert on
// what was actually created rather than only that it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetCloudDeployCustomTargetTypeAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, targetTypeID string) *clouddeploy.CustomTargetType {
	targetType, err := GetCloudDeployCustomTargetTypeAttrsE(t, ctx, projectID, location, targetTypeID)
	require.NoError(t, err)

	return targetType
}

// GetCloudDeployCustomTargetTypeAttrsE returns the settings Google Cloud holds for the given custom target type.
// The ctx parameter supports cancellation and timeouts.
func GetCloudDeployCustomTargetTypeAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, targetTypeID string) (*clouddeploy.CustomTargetType, error) {
	logger.Default.Logf(t, "Getting settings for custom target type %s in %s in project %s", targetTypeID, location, projectID)

	service, err := NewCloudDeployServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetCloudDeployCustomTargetTypeAttrsWithClient(ctx, service, projectID, location, targetTypeID)
}

// GetCloudDeployCustomTargetTypeAttrsWithClient returns the settings Google Cloud holds for the given custom target type using the supplied
// *clouddeploy.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see clouddeploy_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetCloudDeployCustomTargetTypeAttrsWithClient(ctx context.Context, service *clouddeploy.Service, projectID string, location string, targetTypeID string) (*clouddeploy.CustomTargetType, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/customTargetTypes/%s", projectID, location, targetTypeID)

	targetType, err := service.Projects.Locations.CustomTargetTypes.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the custom target type %s does not exist in %s in project %s", targetTypeID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for custom target type %s in %s in project %s: %w", targetTypeID, location, projectID, err)
	}

	return targetType, nil
}

// GetCloudDeployPolicyAttrs returns the settings Google Cloud holds for the given deploy policy, so a test can assert on
// what was actually created rather than only that it exists. This is the policy that says when a rollout may happen,
// not the IAM policy that says who may start one.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetCloudDeployPolicyAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, policyID string) *clouddeploy.DeployPolicy {
	policy, err := GetCloudDeployPolicyAttrsE(t, ctx, projectID, location, policyID)
	require.NoError(t, err)

	return policy
}

// GetCloudDeployPolicyAttrsE returns the settings Google Cloud holds for the given deploy policy.
// The ctx parameter supports cancellation and timeouts.
func GetCloudDeployPolicyAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, policyID string) (*clouddeploy.DeployPolicy, error) {
	logger.Default.Logf(t, "Getting settings for deploy policy %s in %s in project %s", policyID, location, projectID)

	service, err := NewCloudDeployServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetCloudDeployPolicyAttrsWithClient(ctx, service, projectID, location, policyID)
}

// GetCloudDeployPolicyAttrsWithClient returns the settings Google Cloud holds for the given deploy policy using the supplied
// *clouddeploy.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see clouddeploy_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetCloudDeployPolicyAttrsWithClient(ctx context.Context, service *clouddeploy.Service, projectID string, location string, policyID string) (*clouddeploy.DeployPolicy, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/deployPolicies/%s", projectID, location, policyID)

	policy, err := service.Projects.Locations.DeployPolicies.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the deploy policy %s does not exist in %s in project %s", policyID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for deploy policy %s in %s in project %s: %w", policyID, location, projectID, err)
	}

	return policy, nil
}

// GetCloudDeployDeliveryPipelineIamPolicyAttrs returns the IAM policy Google Cloud holds for the given delivery pipeline, so a test
// can assert on who was actually granted access to it.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetCloudDeployDeliveryPipelineIamPolicyAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, pipelineID string) *clouddeploy.Policy {
	policy, err := GetCloudDeployDeliveryPipelineIamPolicyAttrsE(t, ctx, projectID, location, pipelineID)
	require.NoError(t, err)

	return policy
}

// GetCloudDeployDeliveryPipelineIamPolicyAttrsE returns the IAM policy Google Cloud holds for the given delivery pipeline.
// The ctx parameter supports cancellation and timeouts.
func GetCloudDeployDeliveryPipelineIamPolicyAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, pipelineID string) (*clouddeploy.Policy, error) {
	logger.Default.Logf(t, "Getting the IAM policy for delivery pipeline %s in %s in project %s", pipelineID, location, projectID)

	service, err := NewCloudDeployServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetCloudDeployDeliveryPipelineIamPolicyAttrsWithClient(ctx, service, projectID, location, pipelineID)
}

// GetCloudDeployDeliveryPipelineIamPolicyAttrsWithClient returns the IAM policy Google Cloud holds for the given delivery pipeline
// using the supplied *clouddeploy.Service. Prefer this variant in unit tests where the service is
// backed by an httptest fake server (see clouddeploy_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetCloudDeployDeliveryPipelineIamPolicyAttrsWithClient(ctx context.Context, service *clouddeploy.Service, projectID string, location string, pipelineID string) (*clouddeploy.Policy, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/deliveryPipelines/%s", projectID, location, pipelineID)

	// A policy carrying a conditional binding is only returned in full at version 3, so that is what
	// is asked for: at a lower version Google drops the condition or refuses the call outright.
	policy, err := service.Projects.Locations.DeliveryPipelines.GetIamPolicy(name).
		OptionsRequestedPolicyVersion(iamPolicyVersionWithConditions).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the delivery pipeline %s does not exist in %s in project %s", pipelineID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get the IAM policy for delivery pipeline %s in %s in project %s: %w", pipelineID, location, projectID, err)
	}

	return policy, nil
}

// GetCloudDeployTargetIamPolicyAttrs returns the IAM policy Google Cloud holds for the given target, so a test
// can assert on who was actually granted access to it.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetCloudDeployTargetIamPolicyAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, targetID string) *clouddeploy.Policy {
	policy, err := GetCloudDeployTargetIamPolicyAttrsE(t, ctx, projectID, location, targetID)
	require.NoError(t, err)

	return policy
}

// GetCloudDeployTargetIamPolicyAttrsE returns the IAM policy Google Cloud holds for the given target.
// The ctx parameter supports cancellation and timeouts.
func GetCloudDeployTargetIamPolicyAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, targetID string) (*clouddeploy.Policy, error) {
	logger.Default.Logf(t, "Getting the IAM policy for target %s in %s in project %s", targetID, location, projectID)

	service, err := NewCloudDeployServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetCloudDeployTargetIamPolicyAttrsWithClient(ctx, service, projectID, location, targetID)
}

// GetCloudDeployTargetIamPolicyAttrsWithClient returns the IAM policy Google Cloud holds for the given target
// using the supplied *clouddeploy.Service. Prefer this variant in unit tests where the service is
// backed by an httptest fake server (see clouddeploy_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetCloudDeployTargetIamPolicyAttrsWithClient(ctx context.Context, service *clouddeploy.Service, projectID string, location string, targetID string) (*clouddeploy.Policy, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/targets/%s", projectID, location, targetID)

	// A policy carrying a conditional binding is only returned in full at version 3, so that is what
	// is asked for: at a lower version Google drops the condition or refuses the call outright.
	policy, err := service.Projects.Locations.Targets.GetIamPolicy(name).
		OptionsRequestedPolicyVersion(iamPolicyVersionWithConditions).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the target %s does not exist in %s in project %s", targetID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get the IAM policy for target %s in %s in project %s: %w", targetID, location, projectID, err)
	}

	return policy, nil
}
