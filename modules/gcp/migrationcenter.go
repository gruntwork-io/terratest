package gcp

import (
	"context"
	"errors"
	"fmt"

	"github.com/gruntwork-io/terratest/modules/core/v2/logger"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/migrationcenter/v1"
	"google.golang.org/api/option"
)

// GetMigrationCenterGroupAttrs returns the settings Google Cloud holds for the Migration Center group, so a test can
// assert on what was actually created rather than only that it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetMigrationCenterGroupAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, groupID string) *migrationcenter.Group {
	result, err := GetMigrationCenterGroupAttrsE(t, ctx, projectID, location, groupID)
	require.NoError(t, err)

	return result
}

// GetMigrationCenterGroupAttrsE returns the settings Google Cloud holds for the Migration Center group.
// The ctx parameter supports cancellation and timeouts.
func GetMigrationCenterGroupAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, groupID string) (*migrationcenter.Group, error) {
	logger.Default.Logf(t, "Getting settings for group %s in %s in project %s", groupID, location, projectID)

	service, err := NewMigrationCenterServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetMigrationCenterGroupAttrsWithClient(ctx, service, projectID, location, groupID)
}

// GetMigrationCenterGroupAttrsWithClient returns the settings Google Cloud holds for the Migration Center group
// using the supplied *migrationcenter.Service. Prefer this variant in unit tests where the service
// is backed by an httptest fake server (see migrationcenter_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetMigrationCenterGroupAttrsWithClient(ctx context.Context, service *migrationcenter.Service, projectID string, location string, groupID string) (*migrationcenter.Group, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/groups/%s", projectID, location, groupID)

	result, err := service.Projects.Locations.Groups.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Migration Center group %s in %s in project %s does not exist", groupID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for group %s in %s in project %s: %w", groupID, location, projectID, err)
	}

	return result, nil
}

// GetMigrationCenterSourceAttrs returns the settings Google Cloud holds for the Migration Center source, so a test can
// assert on what was actually created rather than only that it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetMigrationCenterSourceAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, sourceID string) *migrationcenter.Source {
	result, err := GetMigrationCenterSourceAttrsE(t, ctx, projectID, location, sourceID)
	require.NoError(t, err)

	return result
}

// GetMigrationCenterSourceAttrsE returns the settings Google Cloud holds for the Migration Center source.
// The ctx parameter supports cancellation and timeouts.
func GetMigrationCenterSourceAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, sourceID string) (*migrationcenter.Source, error) {
	logger.Default.Logf(t, "Getting settings for source %s in %s in project %s", sourceID, location, projectID)

	service, err := NewMigrationCenterServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetMigrationCenterSourceAttrsWithClient(ctx, service, projectID, location, sourceID)
}

// GetMigrationCenterSourceAttrsWithClient returns the settings Google Cloud holds for the Migration Center source
// using the supplied *migrationcenter.Service. Prefer this variant in unit tests where the service
// is backed by an httptest fake server (see migrationcenter_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetMigrationCenterSourceAttrsWithClient(ctx context.Context, service *migrationcenter.Service, projectID string, location string, sourceID string) (*migrationcenter.Source, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/sources/%s", projectID, location, sourceID)

	result, err := service.Projects.Locations.Sources.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Migration Center source %s in %s in project %s does not exist", sourceID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for source %s in %s in project %s: %w", sourceID, location, projectID, err)
	}

	return result, nil
}

// GetMigrationCenterPreferenceSetAttrs returns the settings Google Cloud holds for the Migration Center preference set, so a test can
// assert on what was actually created rather than only that it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetMigrationCenterPreferenceSetAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, setID string) *migrationcenter.PreferenceSet {
	result, err := GetMigrationCenterPreferenceSetAttrsE(t, ctx, projectID, location, setID)
	require.NoError(t, err)

	return result
}

// GetMigrationCenterPreferenceSetAttrsE returns the settings Google Cloud holds for the Migration Center preference set.
// The ctx parameter supports cancellation and timeouts.
func GetMigrationCenterPreferenceSetAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, setID string) (*migrationcenter.PreferenceSet, error) {
	logger.Default.Logf(t, "Getting settings for preference set %s in %s in project %s", setID, location, projectID)

	service, err := NewMigrationCenterServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetMigrationCenterPreferenceSetAttrsWithClient(ctx, service, projectID, location, setID)
}

// GetMigrationCenterPreferenceSetAttrsWithClient returns the settings Google Cloud holds for the Migration Center preference set
// using the supplied *migrationcenter.Service. Prefer this variant in unit tests where the service
// is backed by an httptest fake server (see migrationcenter_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetMigrationCenterPreferenceSetAttrsWithClient(ctx context.Context, service *migrationcenter.Service, projectID string, location string, setID string) (*migrationcenter.PreferenceSet, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/preferenceSets/%s", projectID, location, setID)

	result, err := service.Projects.Locations.PreferenceSets.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Migration Center preference set %s in %s in project %s does not exist", setID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for preference set %s in %s in project %s: %w", setID, location, projectID, err)
	}

	return result, nil
}

// GetMigrationCenterReportConfigAttrs returns the settings Google Cloud holds for the Migration Center report config, so a test can
// assert on what was actually created rather than only that it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetMigrationCenterReportConfigAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, configID string) *migrationcenter.ReportConfig {
	result, err := GetMigrationCenterReportConfigAttrsE(t, ctx, projectID, location, configID)
	require.NoError(t, err)

	return result
}

// GetMigrationCenterReportConfigAttrsE returns the settings Google Cloud holds for the Migration Center report config.
// The ctx parameter supports cancellation and timeouts.
func GetMigrationCenterReportConfigAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, configID string) (*migrationcenter.ReportConfig, error) {
	logger.Default.Logf(t, "Getting settings for report config %s in %s in project %s", configID, location, projectID)

	service, err := NewMigrationCenterServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetMigrationCenterReportConfigAttrsWithClient(ctx, service, projectID, location, configID)
}

// GetMigrationCenterReportConfigAttrsWithClient returns the settings Google Cloud holds for the Migration Center report config
// using the supplied *migrationcenter.Service. Prefer this variant in unit tests where the service
// is backed by an httptest fake server (see migrationcenter_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetMigrationCenterReportConfigAttrsWithClient(ctx context.Context, service *migrationcenter.Service, projectID string, location string, configID string) (*migrationcenter.ReportConfig, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/reportConfigs/%s", projectID, location, configID)

	result, err := service.Projects.Locations.ReportConfigs.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Migration Center report config %s in %s in project %s does not exist", configID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for report config %s in %s in project %s: %w", configID, location, projectID, err)
	}

	return result, nil
}

// GetMigrationCenterImportJobAttrs returns the settings Google Cloud holds for the Migration Center import job, so a test can
// assert on what was actually created rather than only that it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetMigrationCenterImportJobAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, jobID string) *migrationcenter.ImportJob {
	result, err := GetMigrationCenterImportJobAttrsE(t, ctx, projectID, location, jobID)
	require.NoError(t, err)

	return result
}

// GetMigrationCenterImportJobAttrsE returns the settings Google Cloud holds for the Migration Center import job.
// The ctx parameter supports cancellation and timeouts.
func GetMigrationCenterImportJobAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, jobID string) (*migrationcenter.ImportJob, error) {
	logger.Default.Logf(t, "Getting settings for import job %s in %s in project %s", jobID, location, projectID)

	service, err := NewMigrationCenterServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetMigrationCenterImportJobAttrsWithClient(ctx, service, projectID, location, jobID)
}

// GetMigrationCenterImportJobAttrsWithClient returns the settings Google Cloud holds for the Migration Center import job
// using the supplied *migrationcenter.Service. Prefer this variant in unit tests where the service
// is backed by an httptest fake server (see migrationcenter_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetMigrationCenterImportJobAttrsWithClient(ctx context.Context, service *migrationcenter.Service, projectID string, location string, jobID string) (*migrationcenter.ImportJob, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/importJobs/%s", projectID, location, jobID)

	result, err := service.Projects.Locations.ImportJobs.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Migration Center import job %s in %s in project %s does not exist", jobID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for import job %s in %s in project %s: %w", jobID, location, projectID, err)
	}

	return result, nil
}

// GetMigrationCenterDiscoveryClientAttrs returns the settings Google Cloud holds for the Migration Center discovery client, so a test can
// assert on what was actually created rather than only that it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetMigrationCenterDiscoveryClientAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, clientID string) *migrationcenter.DiscoveryClient {
	result, err := GetMigrationCenterDiscoveryClientAttrsE(t, ctx, projectID, location, clientID)
	require.NoError(t, err)

	return result
}

// GetMigrationCenterDiscoveryClientAttrsE returns the settings Google Cloud holds for the Migration Center discovery client.
// The ctx parameter supports cancellation and timeouts.
func GetMigrationCenterDiscoveryClientAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, clientID string) (*migrationcenter.DiscoveryClient, error) {
	logger.Default.Logf(t, "Getting settings for discovery client %s in %s in project %s", clientID, location, projectID)

	service, err := NewMigrationCenterServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetMigrationCenterDiscoveryClientAttrsWithClient(ctx, service, projectID, location, clientID)
}

// GetMigrationCenterDiscoveryClientAttrsWithClient returns the settings Google Cloud holds for the Migration Center discovery client
// using the supplied *migrationcenter.Service. Prefer this variant in unit tests where the service
// is backed by an httptest fake server (see migrationcenter_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetMigrationCenterDiscoveryClientAttrsWithClient(ctx context.Context, service *migrationcenter.Service, projectID string, location string, clientID string) (*migrationcenter.DiscoveryClient, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/discoveryClients/%s", projectID, location, clientID)

	result, err := service.Projects.Locations.DiscoveryClients.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Migration Center discovery client %s in %s in project %s does not exist", clientID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for discovery client %s in %s in project %s: %w", clientID, location, projectID, err)
	}

	return result, nil
}

// GetMigrationCenterReportAttrs returns the settings Google Cloud holds for a Migration Center
// report, which lives under the report config that produced it rather than under the project.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetMigrationCenterReportAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, configID string, reportID string) *migrationcenter.Report {
	report, err := GetMigrationCenterReportAttrsE(t, ctx, projectID, location, configID, reportID)
	require.NoError(t, err)

	return report
}

// GetMigrationCenterReportAttrsE returns the settings Google Cloud holds for a Migration Center
// report.
// The ctx parameter supports cancellation and timeouts.
func GetMigrationCenterReportAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, configID string, reportID string) (*migrationcenter.Report, error) {
	logger.Default.Logf(t, "Getting settings for Migration Center report %s under config %s in %s in project %s", reportID, configID, location, projectID)

	service, err := NewMigrationCenterServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetMigrationCenterReportAttrsWithClient(ctx, service, projectID, location, configID, reportID)
}

// GetMigrationCenterReportAttrsWithClient returns the settings Google Cloud holds for a Migration
// Center report using the supplied *migrationcenter.Service. Prefer this variant in unit tests where
// the service is backed by an httptest fake server (see migrationcenter_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetMigrationCenterReportAttrsWithClient(ctx context.Context, service *migrationcenter.Service, projectID string, location string, configID string, reportID string) (*migrationcenter.Report, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/reportConfigs/%s/reports/%s", projectID, location, configID, reportID)

	report, err := service.Projects.Locations.ReportConfigs.Reports.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Migration Center report %s under config %s in %s in project %s does not exist", reportID, configID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Migration Center report %s under config %s in %s in project %s: %w", reportID, configID, location, projectID, err)
	}

	return report, nil
}

// NewMigrationCenterServiceE creates a Migration Center service authenticated the same way every
// other client in this module is.
// The ctx parameter supports cancellation and timeouts.
func NewMigrationCenterServiceE(t testing.TestingT, ctx context.Context) (*migrationcenter.Service, error) {
	return migrationcenter.NewService(ctx, append(withOptions(), option.WithScopes(migrationcenter.CloudPlatformScope))...)
}
