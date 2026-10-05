package gcp

import (
	"context"
	"errors"
	"fmt"

	"github.com/gruntwork-io/terratest/modules/core/v2/logger"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
	"google.golang.org/api/workstations/v1"
)

// GetWorkstationClusterAttrs returns the settings Google Cloud holds for the given workstation cluster, so a test can assert on what was
// actually created rather than only that it exists.
// A cluster is the network a team's workstations sit on, so the network and subnetwork it names decide what they can reach.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetWorkstationClusterAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, clusterID string) *workstations.WorkstationCluster {
	attrs, err := GetWorkstationClusterAttrsE(t, ctx, projectID, location, clusterID)
	require.NoError(t, err)

	return attrs
}

// GetWorkstationClusterAttrsE returns the settings Google Cloud holds for the given workstation cluster.
// The ctx parameter supports cancellation and timeouts.
func GetWorkstationClusterAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, clusterID string) (*workstations.WorkstationCluster, error) {
	logger.Default.Logf(t, "Getting settings for workstation cluster %s in %s in project %s", clusterID, location, projectID)

	service, err := NewWorkstationsServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetWorkstationClusterAttrsWithClient(ctx, service, projectID, location, clusterID)
}

// GetWorkstationClusterAttrsWithClient returns the settings Google Cloud holds for the given workstation cluster using the supplied
// *workstations.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see workstations_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetWorkstationClusterAttrsWithClient(ctx context.Context, service *workstations.Service, projectID string, location string, clusterID string) (*workstations.WorkstationCluster, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/workstationClusters/%s", projectID, location, clusterID)

	attrs, err := service.Projects.Locations.WorkstationClusters.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the workstation cluster %s in %s in project %s does not exist", clusterID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for workstation cluster %s in %s in project %s: %w", clusterID, location, projectID, err)
	}

	return attrs, nil
}

// GetWorkstationConfigAttrs returns the settings Google Cloud holds for the given workstation config, so a test can assert on what was
// actually created rather than only that it exists.
// A config is the machine every workstation made from it gets, so its machine type, its idle timeout and its disks are what a developer actually works on.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetWorkstationConfigAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, clusterID string, configID string) *workstations.WorkstationConfig {
	attrs, err := GetWorkstationConfigAttrsE(t, ctx, projectID, location, clusterID, configID)
	require.NoError(t, err)

	return attrs
}

// GetWorkstationConfigAttrsE returns the settings Google Cloud holds for the given workstation config.
// The ctx parameter supports cancellation and timeouts.
func GetWorkstationConfigAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, clusterID string, configID string) (*workstations.WorkstationConfig, error) {
	logger.Default.Logf(t, "Getting settings for workstation config %s in cluster %s in %s in project %s", configID, clusterID, location, projectID)

	service, err := NewWorkstationsServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetWorkstationConfigAttrsWithClient(ctx, service, projectID, location, clusterID, configID)
}

// GetWorkstationConfigAttrsWithClient returns the settings Google Cloud holds for the given workstation config using the supplied
// *workstations.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see workstations_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetWorkstationConfigAttrsWithClient(ctx context.Context, service *workstations.Service, projectID string, location string, clusterID string, configID string) (*workstations.WorkstationConfig, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/workstationClusters/%s/workstationConfigs/%s", projectID, location, clusterID, configID)

	attrs, err := service.Projects.Locations.WorkstationClusters.WorkstationConfigs.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the workstation config %s in cluster %s in %s in project %s does not exist", configID, clusterID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for workstation config %s in cluster %s in %s in project %s: %w", configID, clusterID, location, projectID, err)
	}

	return attrs, nil
}

// GetWorkstationAttrs returns the settings Google Cloud holds for the given workstation, so a test can assert on what was
// actually created rather than only that it exists.
// A workstation is one developer's machine, so its state and the host it is reachable at are what say whether it can be used.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetWorkstationAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, clusterID string, configID string, workstationID string) *workstations.Workstation {
	attrs, err := GetWorkstationAttrsE(t, ctx, projectID, location, clusterID, configID, workstationID)
	require.NoError(t, err)

	return attrs
}

// GetWorkstationAttrsE returns the settings Google Cloud holds for the given workstation.
// The ctx parameter supports cancellation and timeouts.
func GetWorkstationAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, clusterID string, configID string, workstationID string) (*workstations.Workstation, error) {
	logger.Default.Logf(t, "Getting settings for workstation %s in config %s in cluster %s in %s in project %s", workstationID, configID, clusterID, location, projectID)

	service, err := NewWorkstationsServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetWorkstationAttrsWithClient(ctx, service, projectID, location, clusterID, configID, workstationID)
}

// GetWorkstationAttrsWithClient returns the settings Google Cloud holds for the given workstation using the supplied
// *workstations.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see workstations_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetWorkstationAttrsWithClient(ctx context.Context, service *workstations.Service, projectID string, location string, clusterID string, configID string, workstationID string) (*workstations.Workstation, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/workstationClusters/%s/workstationConfigs/%s/workstations/%s", projectID, location, clusterID, configID, workstationID)

	attrs, err := service.Projects.Locations.WorkstationClusters.WorkstationConfigs.Workstations.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the workstation %s in config %s in cluster %s in %s in project %s does not exist", workstationID, configID, clusterID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for workstation %s in config %s in cluster %s in %s in project %s: %w", workstationID, configID, clusterID, location, projectID, err)
	}

	return attrs, nil
}

// GetWorkstationConfigIamPolicyAttrs returns the IAM policy Google Cloud holds for the given workstation config, so a test can assert on what was
// actually created rather than only that it exists.
// Who may create a workstation from the config is what the policy decides, and it is separate from the policy on any workstation already made from it.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetWorkstationConfigIamPolicyAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, clusterID string, configID string) *workstations.Policy {
	policy, err := GetWorkstationConfigIamPolicyAttrsE(t, ctx, projectID, location, clusterID, configID)
	require.NoError(t, err)

	return policy
}

// GetWorkstationConfigIamPolicyAttrsE returns the IAM policy Google Cloud holds for the given workstation config.
// The ctx parameter supports cancellation and timeouts.
func GetWorkstationConfigIamPolicyAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, clusterID string, configID string) (*workstations.Policy, error) {
	logger.Default.Logf(t, "Getting the IAM policy for workstation config %s in cluster %s in %s in project %s", configID, clusterID, location, projectID)

	service, err := NewWorkstationsServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetWorkstationConfigIamPolicyAttrsWithClient(ctx, service, projectID, location, clusterID, configID)
}

// GetWorkstationConfigIamPolicyAttrsWithClient returns the IAM policy Google Cloud holds for the given workstation config using the supplied
// *workstations.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see workstations_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetWorkstationConfigIamPolicyAttrsWithClient(ctx context.Context, service *workstations.Service, projectID string, location string, clusterID string, configID string) (*workstations.Policy, error) {
	resource := fmt.Sprintf("projects/%s/locations/%s/workstationClusters/%s/workstationConfigs/%s", projectID, location, clusterID, configID)

	policy, err := service.Projects.Locations.WorkstationClusters.WorkstationConfigs.GetIamPolicy(resource).OptionsRequestedPolicyVersion(iamPolicyVersionWithConditions).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the workstation config %s in cluster %s in %s in project %s does not exist", configID, clusterID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get the IAM policy for workstation config %s in cluster %s in %s in project %s: %w", configID, clusterID, location, projectID, err)
	}

	return policy, nil
}

// GetWorkstationIamPolicyAttrs returns the IAM policy Google Cloud holds for the given workstation, so a test can assert on what was
// actually created rather than only that it exists.
// Who may use the running machine is what the policy decides, and that is a different question from who may create one.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetWorkstationIamPolicyAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, clusterID string, configID string, workstationID string) *workstations.Policy {
	policy, err := GetWorkstationIamPolicyAttrsE(t, ctx, projectID, location, clusterID, configID, workstationID)
	require.NoError(t, err)

	return policy
}

// GetWorkstationIamPolicyAttrsE returns the IAM policy Google Cloud holds for the given workstation.
// The ctx parameter supports cancellation and timeouts.
func GetWorkstationIamPolicyAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, clusterID string, configID string, workstationID string) (*workstations.Policy, error) {
	logger.Default.Logf(t, "Getting the IAM policy for workstation %s in config %s in cluster %s in %s in project %s", workstationID, configID, clusterID, location, projectID)

	service, err := NewWorkstationsServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetWorkstationIamPolicyAttrsWithClient(ctx, service, projectID, location, clusterID, configID, workstationID)
}

// GetWorkstationIamPolicyAttrsWithClient returns the IAM policy Google Cloud holds for the given workstation using the supplied
// *workstations.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see workstations_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetWorkstationIamPolicyAttrsWithClient(ctx context.Context, service *workstations.Service, projectID string, location string, clusterID string, configID string, workstationID string) (*workstations.Policy, error) {
	resource := fmt.Sprintf("projects/%s/locations/%s/workstationClusters/%s/workstationConfigs/%s/workstations/%s", projectID, location, clusterID, configID, workstationID)

	policy, err := service.Projects.Locations.WorkstationClusters.WorkstationConfigs.Workstations.GetIamPolicy(resource).OptionsRequestedPolicyVersion(iamPolicyVersionWithConditions).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the workstation %s in config %s in cluster %s in %s in project %s does not exist", workstationID, configID, clusterID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get the IAM policy for workstation %s in config %s in cluster %s in %s in project %s: %w", workstationID, configID, clusterID, location, projectID, err)
	}

	return policy, nil
}

// NewWorkstationsServiceE creates a Cloud Workstations service authenticated the same way every other client in this
// module is.
// The ctx parameter supports cancellation and timeouts.
func NewWorkstationsServiceE(t testing.TestingT, ctx context.Context) (*workstations.Service, error) {
	return workstations.NewService(ctx, append(withOptions(), option.WithScopes(workstations.CloudPlatformScope))...)
}
