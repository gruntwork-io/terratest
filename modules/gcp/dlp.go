package gcp

import (
	"context"
	"errors"
	"fmt"

	"github.com/gruntwork-io/terratest/modules/core/v2/logger"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/dlp/v2"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
)

// GetDLPInspectTemplateAttrs returns the settings Google Cloud holds for the given inspection template, so a test can assert on
// what was actually created rather than only that it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetDLPInspectTemplateAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, templateID string) *dlp.GooglePrivacyDlpV2InspectTemplate {
	template, err := GetDLPInspectTemplateAttrsE(t, ctx, projectID, location, templateID)
	require.NoError(t, err)

	return template
}

// GetDLPInspectTemplateAttrsE returns the settings Google Cloud holds for the given inspection template.
// The ctx parameter supports cancellation and timeouts.
func GetDLPInspectTemplateAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, templateID string) (*dlp.GooglePrivacyDlpV2InspectTemplate, error) {
	logger.Default.Logf(t, "Getting settings for inspection template %s in %s in project %s", templateID, location, projectID)

	service, err := NewDLPServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetDLPInspectTemplateAttrsWithClient(ctx, service, projectID, location, templateID)
}

// GetDLPInspectTemplateAttrsWithClient returns the settings Google Cloud holds for the given inspection template using the supplied
// *dlp.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see dlp_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetDLPInspectTemplateAttrsWithClient(ctx context.Context, service *dlp.Service, projectID string, location string, templateID string) (*dlp.GooglePrivacyDlpV2InspectTemplate, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/inspectTemplates/%s", projectID, location, templateID)

	template, err := service.Projects.Locations.InspectTemplates.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the inspection template %s does not exist in %s in project %s", templateID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for inspection template %s in %s in project %s: %w", templateID, location, projectID, err)
	}

	return template, nil
}

// GetDLPDeidentifyTemplateAttrs returns the settings Google Cloud holds for the given de-identification template, so a test can assert on
// what was actually created rather than only that it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetDLPDeidentifyTemplateAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, templateID string) *dlp.GooglePrivacyDlpV2DeidentifyTemplate {
	template, err := GetDLPDeidentifyTemplateAttrsE(t, ctx, projectID, location, templateID)
	require.NoError(t, err)

	return template
}

// GetDLPDeidentifyTemplateAttrsE returns the settings Google Cloud holds for the given de-identification template.
// The ctx parameter supports cancellation and timeouts.
func GetDLPDeidentifyTemplateAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, templateID string) (*dlp.GooglePrivacyDlpV2DeidentifyTemplate, error) {
	logger.Default.Logf(t, "Getting settings for de-identification template %s in %s in project %s", templateID, location, projectID)

	service, err := NewDLPServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetDLPDeidentifyTemplateAttrsWithClient(ctx, service, projectID, location, templateID)
}

// GetDLPDeidentifyTemplateAttrsWithClient returns the settings Google Cloud holds for the given de-identification template using the supplied
// *dlp.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see dlp_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetDLPDeidentifyTemplateAttrsWithClient(ctx context.Context, service *dlp.Service, projectID string, location string, templateID string) (*dlp.GooglePrivacyDlpV2DeidentifyTemplate, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/deidentifyTemplates/%s", projectID, location, templateID)

	template, err := service.Projects.Locations.DeidentifyTemplates.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the de-identification template %s does not exist in %s in project %s", templateID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for de-identification template %s in %s in project %s: %w", templateID, location, projectID, err)
	}

	return template, nil
}

// GetDLPStoredInfoTypeAttrs returns the settings Google Cloud holds for the given stored info type, so a test can assert on
// what was actually created rather than only that it exists. What was asked for is on the pending or current version
// rather than on the type itself.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetDLPStoredInfoTypeAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, infoTypeID string) *dlp.GooglePrivacyDlpV2StoredInfoType {
	infoType, err := GetDLPStoredInfoTypeAttrsE(t, ctx, projectID, location, infoTypeID)
	require.NoError(t, err)

	return infoType
}

// GetDLPStoredInfoTypeAttrsE returns the settings Google Cloud holds for the given stored info type.
// The ctx parameter supports cancellation and timeouts.
func GetDLPStoredInfoTypeAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, infoTypeID string) (*dlp.GooglePrivacyDlpV2StoredInfoType, error) {
	logger.Default.Logf(t, "Getting settings for stored info type %s in %s in project %s", infoTypeID, location, projectID)

	service, err := NewDLPServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetDLPStoredInfoTypeAttrsWithClient(ctx, service, projectID, location, infoTypeID)
}

// GetDLPStoredInfoTypeAttrsWithClient returns the settings Google Cloud holds for the given stored info type using the supplied
// *dlp.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see dlp_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetDLPStoredInfoTypeAttrsWithClient(ctx context.Context, service *dlp.Service, projectID string, location string, infoTypeID string) (*dlp.GooglePrivacyDlpV2StoredInfoType, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/storedInfoTypes/%s", projectID, location, infoTypeID)

	infoType, err := service.Projects.Locations.StoredInfoTypes.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the stored info type %s does not exist in %s in project %s", infoTypeID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for stored info type %s in %s in project %s: %w", infoTypeID, location, projectID, err)
	}

	return infoType, nil
}

// GetDLPJobTriggerAttrs returns the settings Google Cloud holds for the given job trigger, so a test can assert on
// what was actually created rather than only that it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetDLPJobTriggerAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, triggerID string) *dlp.GooglePrivacyDlpV2JobTrigger {
	trigger, err := GetDLPJobTriggerAttrsE(t, ctx, projectID, location, triggerID)
	require.NoError(t, err)

	return trigger
}

// GetDLPJobTriggerAttrsE returns the settings Google Cloud holds for the given job trigger.
// The ctx parameter supports cancellation and timeouts.
func GetDLPJobTriggerAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, triggerID string) (*dlp.GooglePrivacyDlpV2JobTrigger, error) {
	logger.Default.Logf(t, "Getting settings for job trigger %s in %s in project %s", triggerID, location, projectID)

	service, err := NewDLPServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetDLPJobTriggerAttrsWithClient(ctx, service, projectID, location, triggerID)
}

// GetDLPJobTriggerAttrsWithClient returns the settings Google Cloud holds for the given job trigger using the supplied
// *dlp.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see dlp_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetDLPJobTriggerAttrsWithClient(ctx context.Context, service *dlp.Service, projectID string, location string, triggerID string) (*dlp.GooglePrivacyDlpV2JobTrigger, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/jobTriggers/%s", projectID, location, triggerID)

	trigger, err := service.Projects.Locations.JobTriggers.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the job trigger %s does not exist in %s in project %s", triggerID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for job trigger %s in %s in project %s: %w", triggerID, location, projectID, err)
	}

	return trigger, nil
}

// NewDLPServiceE creates a Sensitive Data Protection service authenticated the same way every other client in this
// module is.
// The ctx parameter supports cancellation and timeouts.
func NewDLPServiceE(t testing.TestingT, ctx context.Context) (*dlp.Service, error) {
	return dlp.NewService(ctx, append(withOptions(), option.WithScopes(dlp.CloudPlatformScope))...)
}
