package gcp

import (
	"context"
	"errors"
	"fmt"
	"regexp"

	"github.com/gruntwork-io/terratest/modules/core/v2/logger"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/dataproc/v1"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
)

// GetDataprocClusterAttrs returns the settings Google Cloud holds for the given Dataproc cluster, so
// a test can assert on what was actually created rather than only that it exists. A cluster is read
// from the region it was created in, which is part of its address rather than a filter.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetDataprocClusterAttrs(t testing.TestingT, ctx context.Context, projectID string, region string, clusterID string) *dataproc.Cluster {
	cluster, err := GetDataprocClusterAttrsE(t, ctx, projectID, region, clusterID)
	require.NoError(t, err)

	return cluster
}

// GetDataprocClusterAttrsE returns the settings Google Cloud holds for the given Dataproc cluster.
// The ctx parameter supports cancellation and timeouts.
func GetDataprocClusterAttrsE(t testing.TestingT, ctx context.Context, projectID string, region string, clusterID string) (*dataproc.Cluster, error) {
	logger.Default.Logf(t, "Getting settings for Dataproc cluster %s in region %s in project %s", clusterID, region, projectID)

	service, err := NewDataprocServiceE(t, ctx, region)
	if err != nil {
		return nil, err
	}

	return GetDataprocClusterAttrsWithClient(ctx, service, projectID, region, clusterID)
}

// GetDataprocClusterAttrsWithClient returns the settings Google Cloud holds for the given Dataproc
// cluster using the supplied *dataproc.Service, which has to point at that region's host. Prefer
// this variant in unit tests where the service is backed by an httptest fake server (see
// dataproc_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetDataprocClusterAttrsWithClient(ctx context.Context, service *dataproc.Service, projectID string, region string, clusterID string) (*dataproc.Cluster, error) {
	cluster, err := service.Projects.Regions.Clusters.Get(projectID, region, clusterID).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Dataproc cluster %s does not exist in region %s in project %s", clusterID, region, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Dataproc cluster %s in region %s in project %s: %w", clusterID, region, projectID, err)
	}

	return cluster, nil
}

// dataprocRegionPattern is what a Google Cloud region identifier may contain. It is checked before
// a region reaches an endpoint, because a value carrying a slash, a colon or an at sign would build
// a URL pointing at a host of the caller's choosing rather than at Google.
var dataprocRegionPattern = regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+)*$`)

// GetDataprocBatchAttrs returns the settings Google Cloud holds for the given Dataproc batch, so a test can assert on what was
// actually created rather than only that it exists.
// A batch is one serverless Spark run, so what it ran and how it ended are the only things it has to say.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetDataprocBatchAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, batchID string) *dataproc.Batch {
	attrs, err := GetDataprocBatchAttrsE(t, ctx, projectID, location, batchID)
	require.NoError(t, err)

	return attrs
}

// GetDataprocBatchAttrsE returns the settings Google Cloud holds for the given Dataproc batch.
// The ctx parameter supports cancellation and timeouts.
func GetDataprocBatchAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, batchID string) (*dataproc.Batch, error) {
	logger.Default.Logf(t, "Getting settings for Dataproc batch %s in %s in project %s", batchID, location, projectID)

	service, err := NewDataprocServiceE(t, ctx, location)
	if err != nil {
		return nil, err
	}

	return GetDataprocBatchAttrsWithClient(ctx, service, projectID, location, batchID)
}

// GetDataprocBatchAttrsWithClient returns the settings Google Cloud holds for the given Dataproc batch using the supplied
// *dataproc.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see dataproc_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetDataprocBatchAttrsWithClient(ctx context.Context, service *dataproc.Service, projectID string, location string, batchID string) (*dataproc.Batch, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/batches/%s", projectID, location, batchID)

	attrs, err := service.Projects.Locations.Batches.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Dataproc batch %s in %s in project %s does not exist", batchID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Dataproc batch %s in %s in project %s: %w", batchID, location, projectID, err)
	}

	return attrs, nil
}

// GetDataprocJobAttrs returns the settings Google Cloud holds for the given Dataproc job, so a test can assert on what was
// actually created rather than only that it exists.
// A job is one submission to a cluster, so which cluster ran it and how it ended are the point.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetDataprocJobAttrs(t testing.TestingT, ctx context.Context, projectID string, region string, jobID string) *dataproc.Job {
	attrs, err := GetDataprocJobAttrsE(t, ctx, projectID, region, jobID)
	require.NoError(t, err)

	return attrs
}

// GetDataprocJobAttrsE returns the settings Google Cloud holds for the given Dataproc job.
// The ctx parameter supports cancellation and timeouts.
func GetDataprocJobAttrsE(t testing.TestingT, ctx context.Context, projectID string, region string, jobID string) (*dataproc.Job, error) {
	logger.Default.Logf(t, "Getting settings for Dataproc job %s in %s in project %s", jobID, region, projectID)

	service, err := NewDataprocServiceE(t, ctx, region)
	if err != nil {
		return nil, err
	}

	return GetDataprocJobAttrsWithClient(ctx, service, projectID, region, jobID)
}

// GetDataprocJobAttrsWithClient returns the settings Google Cloud holds for the given Dataproc job using the supplied
// *dataproc.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see dataproc_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetDataprocJobAttrsWithClient(ctx context.Context, service *dataproc.Service, projectID string, region string, jobID string) (*dataproc.Job, error) {
	attrs, err := service.Projects.Regions.Jobs.Get(projectID, region, jobID).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Dataproc job %s in %s in project %s does not exist", jobID, region, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Dataproc job %s in %s in project %s: %w", jobID, region, projectID, err)
	}

	return attrs, nil
}

// GetDataprocAutoscalingPolicyAttrs returns the settings Google Cloud holds for the given Dataproc autoscaling policy, so a test can assert on
// what was actually created rather than only that it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetDataprocAutoscalingPolicyAttrs(t testing.TestingT, ctx context.Context, projectID string, region string, policyID string) *dataproc.AutoscalingPolicy {
	policy, err := GetDataprocAutoscalingPolicyAttrsE(t, ctx, projectID, region, policyID)
	require.NoError(t, err)

	return policy
}

// GetDataprocAutoscalingPolicyAttrsE returns the settings Google Cloud holds for the given Dataproc autoscaling policy.
// The ctx parameter supports cancellation and timeouts.
func GetDataprocAutoscalingPolicyAttrsE(t testing.TestingT, ctx context.Context, projectID string, region string, policyID string) (*dataproc.AutoscalingPolicy, error) {
	logger.Default.Logf(t, "Getting settings for Dataproc autoscaling policy %s in %s in project %s", policyID, region, projectID)

	service, err := NewDataprocServiceE(t, ctx, region)
	if err != nil {
		return nil, err
	}

	return GetDataprocAutoscalingPolicyAttrsWithClient(ctx, service, projectID, region, policyID)
}

// GetDataprocAutoscalingPolicyAttrsWithClient returns the settings Google Cloud holds for the given Dataproc autoscaling policy using the supplied
// *dataproc.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see dataproc_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetDataprocAutoscalingPolicyAttrsWithClient(ctx context.Context, service *dataproc.Service, projectID string, region string, policyID string) (*dataproc.AutoscalingPolicy, error) {
	name := fmt.Sprintf("projects/%s/regions/%s/autoscalingPolicies/%s", projectID, region, policyID)

	policy, err := service.Projects.Regions.AutoscalingPolicies.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Dataproc autoscaling policy %s does not exist in %s in project %s", policyID, region, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Dataproc autoscaling policy %s in %s in project %s: %w", policyID, region, projectID, err)
	}

	return policy, nil
}

// GetDataprocWorkflowTemplateAttrs returns the settings Google Cloud holds for the given Dataproc workflow template, so a test can assert on
// what was actually created rather than only that it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetDataprocWorkflowTemplateAttrs(t testing.TestingT, ctx context.Context, projectID string, region string, templateID string) *dataproc.WorkflowTemplate {
	template, err := GetDataprocWorkflowTemplateAttrsE(t, ctx, projectID, region, templateID)
	require.NoError(t, err)

	return template
}

// GetDataprocWorkflowTemplateAttrsE returns the settings Google Cloud holds for the given Dataproc workflow template.
// The ctx parameter supports cancellation and timeouts.
func GetDataprocWorkflowTemplateAttrsE(t testing.TestingT, ctx context.Context, projectID string, region string, templateID string) (*dataproc.WorkflowTemplate, error) {
	logger.Default.Logf(t, "Getting settings for Dataproc workflow template %s in %s in project %s", templateID, region, projectID)

	service, err := NewDataprocServiceE(t, ctx, region)
	if err != nil {
		return nil, err
	}

	return GetDataprocWorkflowTemplateAttrsWithClient(ctx, service, projectID, region, templateID)
}

// GetDataprocWorkflowTemplateAttrsWithClient returns the settings Google Cloud holds for the given Dataproc workflow template using the supplied
// *dataproc.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see dataproc_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetDataprocWorkflowTemplateAttrsWithClient(ctx context.Context, service *dataproc.Service, projectID string, region string, templateID string) (*dataproc.WorkflowTemplate, error) {
	name := fmt.Sprintf("projects/%s/regions/%s/workflowTemplates/%s", projectID, region, templateID)

	template, err := service.Projects.Regions.WorkflowTemplates.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Dataproc workflow template %s does not exist in %s in project %s", templateID, region, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Dataproc workflow template %s in %s in project %s: %w", templateID, region, projectID, err)
	}

	return template, nil
}

// GetDataprocSessionTemplateAttrs returns the settings Google Cloud holds for the given Dataproc session template, so a test can assert on
// what was actually created rather than only that it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetDataprocSessionTemplateAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, templateID string) *dataproc.SessionTemplate {
	template, err := GetDataprocSessionTemplateAttrsE(t, ctx, projectID, location, templateID)
	require.NoError(t, err)

	return template
}

// GetDataprocSessionTemplateAttrsE returns the settings Google Cloud holds for the given Dataproc session template.
// The ctx parameter supports cancellation and timeouts.
func GetDataprocSessionTemplateAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, templateID string) (*dataproc.SessionTemplate, error) {
	logger.Default.Logf(t, "Getting settings for Dataproc session template %s in %s in project %s", templateID, location, projectID)

	service, err := NewDataprocServiceE(t, ctx, location)
	if err != nil {
		return nil, err
	}

	return GetDataprocSessionTemplateAttrsWithClient(ctx, service, projectID, location, templateID)
}

// GetDataprocSessionTemplateAttrsWithClient returns the settings Google Cloud holds for the given Dataproc session template using the supplied
// *dataproc.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see dataproc_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetDataprocSessionTemplateAttrsWithClient(ctx context.Context, service *dataproc.Service, projectID string, location string, templateID string) (*dataproc.SessionTemplate, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/sessionTemplates/%s", projectID, location, templateID)

	template, err := service.Projects.Locations.SessionTemplates.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Dataproc session template %s does not exist in %s in project %s", templateID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Dataproc session template %s in %s in project %s: %w", templateID, location, projectID, err)
	}

	return template, nil
}

// GetDataprocJobIamPolicyAttrs returns the IAM policy Google Cloud holds for the given Dataproc job, so a test can assert on what was
// actually created rather than only that it exists.
// Who may read or cancel the job is what the policy decides, and a job Google created on its own carries no binding at all.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetDataprocJobIamPolicyAttrs(t testing.TestingT, ctx context.Context, projectID string, region string, jobID string) *dataproc.Policy {
	policy, err := GetDataprocJobIamPolicyAttrsE(t, ctx, projectID, region, jobID)
	require.NoError(t, err)

	return policy
}

// GetDataprocJobIamPolicyAttrsE returns the IAM policy Google Cloud holds for the given Dataproc job.
// The ctx parameter supports cancellation and timeouts.
func GetDataprocJobIamPolicyAttrsE(t testing.TestingT, ctx context.Context, projectID string, region string, jobID string) (*dataproc.Policy, error) {
	logger.Default.Logf(t, "Getting the IAM policy for Dataproc job %s in %s in project %s", jobID, region, projectID)

	service, err := NewDataprocServiceE(t, ctx, region)
	if err != nil {
		return nil, err
	}

	return GetDataprocJobIamPolicyAttrsWithClient(ctx, service, projectID, region, jobID)
}

// GetDataprocJobIamPolicyAttrsWithClient returns the IAM policy Google Cloud holds for the given Dataproc job using the supplied
// *dataproc.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see dataproc_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetDataprocJobIamPolicyAttrsWithClient(ctx context.Context, service *dataproc.Service, projectID string, region string, jobID string) (*dataproc.Policy, error) {
	resource := fmt.Sprintf("projects/%s/regions/%s/jobs/%s", projectID, region, jobID)

	policy, err := service.Projects.Regions.Jobs.GetIamPolicy(resource, &dataproc.GetIamPolicyRequest{
		Options: &dataproc.GetPolicyOptions{RequestedPolicyVersion: iamPolicyVersionWithConditions},
	}).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Dataproc job %s in %s in project %s does not exist", jobID, region, projectID)
		}

		return nil, fmt.Errorf("failed to get the IAM policy for Dataproc job %s in %s in project %s: %w", jobID, region, projectID, err)
	}

	return policy, nil
}

// NewDataprocServiceE creates a Dataproc service authenticated the same way every other client in
// this module is. Dataproc answers on a per-region host for every region but `global`, so the
// region the caller is asking about decides which one this talks to, and a region that is not a
// plain identifier is refused.
// The ctx parameter supports cancellation and timeouts.
func NewDataprocServiceE(t testing.TestingT, ctx context.Context, region string) (*dataproc.Service, error) {
	if !dataprocRegionPattern.MatchString(region) {
		return nil, fmt.Errorf("%q is not a valid region: a region may hold only lowercase letters, digits and hyphens", region)
	}

	opts := append(withOptions(), option.WithScopes(dataproc.CloudPlatformScope))
	if region != "global" {
		opts = append(opts, option.WithEndpoint(fmt.Sprintf("https://%s-dataproc.googleapis.com/", region)))
	}

	return dataproc.NewService(ctx, opts...)
}
