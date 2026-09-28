package gcp

import (
	"context"
	"errors"
	"fmt"

	"github.com/gruntwork-io/terratest/modules/core/v2/logger"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/logging/v2"
	"google.golang.org/api/option"
)

// GetLogMetricAttrs returns the settings Google Cloud holds for the given log-based metric, so a
// test can assert on what was actually created rather than only that it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetLogMetricAttrs(t testing.TestingT, ctx context.Context, projectID string, metricName string) *logging.LogMetric {
	metric, err := GetLogMetricAttrsE(t, ctx, projectID, metricName)
	require.NoError(t, err)

	return metric
}

// GetLogMetricAttrsE returns the settings Google Cloud holds for the given log-based metric.
// The ctx parameter supports cancellation and timeouts.
func GetLogMetricAttrsE(t testing.TestingT, ctx context.Context, projectID string, metricName string) (*logging.LogMetric, error) {
	logger.Default.Logf(t, "Getting settings for log-based metric %s in project %s", metricName, projectID)

	service, err := NewLoggingServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetLogMetricAttrsWithClient(ctx, service, projectID, metricName)
}

// GetLogMetricAttrsWithClient returns the settings Google Cloud holds for the given log-based
// metric using the supplied *logging.Service. Prefer this variant in unit tests where the service
// is backed by an httptest fake server (see logging_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetLogMetricAttrsWithClient(ctx context.Context, service *logging.Service, projectID string, metricName string) (*logging.LogMetric, error) {
	name := fmt.Sprintf("projects/%s/metrics/%s", projectID, metricName)

	metric, err := service.Projects.Metrics.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("log-based metric %s does not exist in project %s", metricName, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for log-based metric %s in project %s: %w", metricName, projectID, err)
	}

	return metric, nil
}

// GetLogSinkAttrs returns the settings Google Cloud holds for the given project log sink, so a
// test can assert on what was actually created rather than only that it exists. The writer
// identity comes back with it, which is the service account a destination has to grant.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetLogSinkAttrs(t testing.TestingT, ctx context.Context, projectID string, sinkName string) *logging.LogSink {
	sink, err := GetLogSinkAttrsE(t, ctx, projectID, sinkName)
	require.NoError(t, err)

	return sink
}

// GetLogSinkAttrsE returns the settings Google Cloud holds for the given project log sink.
// The ctx parameter supports cancellation and timeouts.
func GetLogSinkAttrsE(t testing.TestingT, ctx context.Context, projectID string, sinkName string) (*logging.LogSink, error) {
	logger.Default.Logf(t, "Getting settings for log sink %s in project %s", sinkName, projectID)

	service, err := NewLoggingServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetLogSinkAttrsWithClient(ctx, service, projectID, sinkName)
}

// GetLogSinkAttrsWithClient returns the settings Google Cloud holds for the given project log sink
// using the supplied *logging.Service. Prefer this variant in unit tests where the service is
// backed by an httptest fake server (see logging_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetLogSinkAttrsWithClient(ctx context.Context, service *logging.Service, projectID string, sinkName string) (*logging.LogSink, error) {
	name := fmt.Sprintf("projects/%s/sinks/%s", projectID, sinkName)

	sink, err := service.Projects.Sinks.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("log sink %s does not exist in project %s", sinkName, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for log sink %s in project %s: %w", sinkName, projectID, err)
	}

	return sink, nil
}

// GetLogBucketAttrs returns the settings Google Cloud holds for the given log bucket, so a test can
// assert on what was actually created rather than only that it exists. A bucket is named by its
// location as well as its id, and a deleted one is still returned for seven days with a
// lifecycleState of DELETE_REQUESTED rather than as absent.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetLogBucketAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, bucketID string) *logging.LogBucket {
	bucket, err := GetLogBucketAttrsE(t, ctx, projectID, location, bucketID)
	require.NoError(t, err)

	return bucket
}

// GetLogBucketAttrsE returns the settings Google Cloud holds for the given log bucket.
// The ctx parameter supports cancellation and timeouts.
func GetLogBucketAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, bucketID string) (*logging.LogBucket, error) {
	logger.Default.Logf(t, "Getting settings for log bucket %s in location %s in project %s", bucketID, location, projectID)

	service, err := NewLoggingServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetLogBucketAttrsWithClient(ctx, service, projectID, location, bucketID)
}

// GetLogBucketAttrsWithClient returns the settings Google Cloud holds for the given log bucket using
// the supplied *logging.Service. Prefer this variant in unit tests where the service is backed by an
// httptest fake server (see logging_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetLogBucketAttrsWithClient(ctx context.Context, service *logging.Service, projectID string, location string, bucketID string) (*logging.LogBucket, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/buckets/%s", projectID, location, bucketID)

	bucket, err := service.Projects.Locations.Buckets.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("log bucket %s does not exist in location %s in project %s", bucketID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for log bucket %s in location %s in project %s: %w", bucketID, location, projectID, err)
	}

	return bucket, nil
}

// NewLoggingServiceE creates a Cloud Logging service authenticated the same way every other client
// in this module is.
// The ctx parameter supports cancellation and timeouts.
func NewLoggingServiceE(t testing.TestingT, ctx context.Context) (*logging.Service, error) {
	return logging.NewService(ctx, append(withOptions(), option.WithScopes(logging.CloudPlatformScope))...)
}

// GetLogExclusionAttrs returns the settings Google Cloud holds for the given log exclusion, so a test can assert on
// what was actually created rather than only that it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetLogExclusionAttrs(t testing.TestingT, ctx context.Context, projectID string, exclusionName string) *logging.LogExclusion {
	exclusion, err := GetLogExclusionAttrsE(t, ctx, projectID, exclusionName)
	require.NoError(t, err)

	return exclusion
}

// GetLogExclusionAttrsE returns the settings Google Cloud holds for the given log exclusion.
// The ctx parameter supports cancellation and timeouts.
func GetLogExclusionAttrsE(t testing.TestingT, ctx context.Context, projectID string, exclusionName string) (*logging.LogExclusion, error) {
	logger.Default.Logf(t, "Getting settings for log exclusion %s in project %s", exclusionName, projectID)

	service, err := NewLoggingServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetLogExclusionAttrsWithClient(ctx, service, projectID, exclusionName)
}

// GetLogExclusionAttrsWithClient returns the settings Google Cloud holds for the given log exclusion using the supplied
// *logging.Service. Prefer this variant in unit tests where the service is backed by an httptest
// fake server (see logging_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetLogExclusionAttrsWithClient(ctx context.Context, service *logging.Service, projectID string, exclusionName string) (*logging.LogExclusion, error) {
	name := fmt.Sprintf("projects/%s/exclusions/%s", projectID, exclusionName)

	exclusion, err := service.Projects.Exclusions.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the log exclusion %s does not exist in project %s", exclusionName, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for log exclusion %s in project %s: %w", exclusionName, projectID, err)
	}

	return exclusion, nil
}

// GetLogViewAttrs returns the settings Google Cloud holds for the given log view, so a test can assert on
// what was actually created rather than only that it exists. A view belongs to a log bucket, so it is named by the
// bucket's location and id as well as its own.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetLogViewAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, bucketID string, viewID string) *logging.LogView {
	view, err := GetLogViewAttrsE(t, ctx, projectID, location, bucketID, viewID)
	require.NoError(t, err)

	return view
}

// GetLogViewAttrsE returns the settings Google Cloud holds for the given log view.
// The ctx parameter supports cancellation and timeouts.
func GetLogViewAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, bucketID string, viewID string) (*logging.LogView, error) {
	logger.Default.Logf(t, "Getting settings for log view %s on bucket %s in %s in project %s", viewID, bucketID, location, projectID)

	service, err := NewLoggingServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetLogViewAttrsWithClient(ctx, service, projectID, location, bucketID, viewID)
}

// GetLogViewAttrsWithClient returns the settings Google Cloud holds for the given log view using the supplied
// *logging.Service. Prefer this variant in unit tests where the service is backed by an httptest
// fake server (see logging_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetLogViewAttrsWithClient(ctx context.Context, service *logging.Service, projectID string, location string, bucketID string, viewID string) (*logging.LogView, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/buckets/%s/views/%s", projectID, location, bucketID, viewID)

	view, err := service.Projects.Locations.Buckets.Views.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the log view %s does not exist on bucket %s in %s in project %s", viewID, bucketID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for log view %s on bucket %s in %s in project %s: %w", viewID, bucketID, location, projectID, err)
	}

	return view, nil
}

// GetLogLinkAttrs returns the settings Google Cloud holds for the given link between a log bucket and a BigQuery dataset, so a test can assert on
// what was actually created rather than only that it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetLogLinkAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, bucketID string, linkID string) *logging.Link {
	link, err := GetLogLinkAttrsE(t, ctx, projectID, location, bucketID, linkID)
	require.NoError(t, err)

	return link
}

// GetLogLinkAttrsE returns the settings Google Cloud holds for the given link between a log bucket and a BigQuery dataset.
// The ctx parameter supports cancellation and timeouts.
func GetLogLinkAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, bucketID string, linkID string) (*logging.Link, error) {
	logger.Default.Logf(t, "Getting settings for link %s on bucket %s in %s in project %s", linkID, bucketID, location, projectID)

	service, err := NewLoggingServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetLogLinkAttrsWithClient(ctx, service, projectID, location, bucketID, linkID)
}

// GetLogLinkAttrsWithClient returns the settings Google Cloud holds for the given link between a log bucket and a BigQuery dataset using the supplied
// *logging.Service. Prefer this variant in unit tests where the service is backed by an httptest
// fake server (see logging_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetLogLinkAttrsWithClient(ctx context.Context, service *logging.Service, projectID string, location string, bucketID string, linkID string) (*logging.Link, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/buckets/%s/links/%s", projectID, location, bucketID, linkID)

	link, err := service.Projects.Locations.Buckets.Links.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the link %s does not exist on bucket %s in %s in project %s", linkID, bucketID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for link %s on bucket %s in %s in project %s: %w", linkID, bucketID, location, projectID, err)
	}

	return link, nil
}

// GetLogScopeAttrs returns the settings Google Cloud holds for the given log scope, so a test can assert on
// what was actually created rather than only that it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetLogScopeAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, scopeID string) *logging.LogScope {
	scope, err := GetLogScopeAttrsE(t, ctx, projectID, location, scopeID)
	require.NoError(t, err)

	return scope
}

// GetLogScopeAttrsE returns the settings Google Cloud holds for the given log scope.
// The ctx parameter supports cancellation and timeouts.
func GetLogScopeAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, scopeID string) (*logging.LogScope, error) {
	logger.Default.Logf(t, "Getting settings for log scope %s in %s in project %s", scopeID, location, projectID)

	service, err := NewLoggingServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetLogScopeAttrsWithClient(ctx, service, projectID, location, scopeID)
}

// GetLogScopeAttrsWithClient returns the settings Google Cloud holds for the given log scope using the supplied
// *logging.Service. Prefer this variant in unit tests where the service is backed by an httptest
// fake server (see logging_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetLogScopeAttrsWithClient(ctx context.Context, service *logging.Service, projectID string, location string, scopeID string) (*logging.LogScope, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/logScopes/%s", projectID, location, scopeID)

	scope, err := service.Projects.Locations.LogScopes.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the log scope %s does not exist in %s in project %s", scopeID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for log scope %s in %s in project %s: %w", scopeID, location, projectID, err)
	}

	return scope, nil
}

// GetLogSavedQueryAttrs returns the settings Google Cloud holds for the given saved query, so a test can assert on
// what was actually created rather than only that it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetLogSavedQueryAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, queryID string) *logging.SavedQuery {
	query, err := GetLogSavedQueryAttrsE(t, ctx, projectID, location, queryID)
	require.NoError(t, err)

	return query
}

// GetLogSavedQueryAttrsE returns the settings Google Cloud holds for the given saved query.
// The ctx parameter supports cancellation and timeouts.
func GetLogSavedQueryAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, queryID string) (*logging.SavedQuery, error) {
	logger.Default.Logf(t, "Getting settings for saved query %s in %s in project %s", queryID, location, projectID)

	service, err := NewLoggingServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetLogSavedQueryAttrsWithClient(ctx, service, projectID, location, queryID)
}

// GetLogSavedQueryAttrsWithClient returns the settings Google Cloud holds for the given saved query using the supplied
// *logging.Service. Prefer this variant in unit tests where the service is backed by an httptest
// fake server (see logging_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetLogSavedQueryAttrsWithClient(ctx context.Context, service *logging.Service, projectID string, location string, queryID string) (*logging.SavedQuery, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/savedQueries/%s", projectID, location, queryID)

	query, err := service.Projects.Locations.SavedQueries.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the saved query %s does not exist in %s in project %s", queryID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for saved query %s in %s in project %s: %w", queryID, location, projectID, err)
	}

	return query, nil
}

// GetLogViewIamPolicyAttrs returns the IAM policy Google Cloud holds for the given log view, so a
// test can assert on who was actually granted access to it. A view carries a policy of its own,
// which is how a bucket's logs are shared without sharing the bucket.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetLogViewIamPolicyAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, bucketID string, viewID string) *logging.Policy {
	policy, err := GetLogViewIamPolicyAttrsE(t, ctx, projectID, location, bucketID, viewID)
	require.NoError(t, err)

	return policy
}

// GetLogViewIamPolicyAttrsE returns the IAM policy Google Cloud holds for the given log view.
// The ctx parameter supports cancellation and timeouts.
func GetLogViewIamPolicyAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, bucketID string, viewID string) (*logging.Policy, error) {
	logger.Default.Logf(t, "Getting the IAM policy for log view %s on bucket %s in %s in project %s", viewID, bucketID, location, projectID)

	service, err := NewLoggingServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetLogViewIamPolicyAttrsWithClient(ctx, service, projectID, location, bucketID, viewID)
}

// GetLogViewIamPolicyAttrsWithClient returns the IAM policy Google Cloud holds for the given log
// view using the supplied *logging.Service. Prefer this variant in unit tests where the service is
// backed by an httptest fake server (see logging_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetLogViewIamPolicyAttrsWithClient(ctx context.Context, service *logging.Service, projectID string, location string, bucketID string, viewID string) (*logging.Policy, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/buckets/%s/views/%s", projectID, location, bucketID, viewID)

	// This call takes a request body rather than a plain resource name, which is why it is written
	// out here rather than generated like the reads around it. A policy carrying a conditional binding
	// is only returned in full at version 3, so that is what is asked for.
	policy, err := service.Projects.Locations.Buckets.Views.GetIamPolicy(name, &logging.GetIamPolicyRequest{
		Options: &logging.GetPolicyOptions{RequestedPolicyVersion: iamPolicyVersionWithConditions},
	}).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the log view %s does not exist on bucket %s in %s in project %s", viewID, bucketID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get the IAM policy for log view %s on bucket %s in %s in project %s: %w", viewID, bucketID, location, projectID, err)
	}

	return policy, nil
}
