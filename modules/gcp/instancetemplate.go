package gcp

import (
	"context"
	"fmt"

	"github.com/gruntwork-io/terratest/modules/core/v2/logger"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/compute/v1"
)

// FetchInstanceTemplate queries GCP to return the settings it holds for the given global
// instance template, so a test can assert on what was actually created rather than only that it
// exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func FetchInstanceTemplate(t testing.TestingT, ctx context.Context, projectID string, name string) *compute.InstanceTemplate {
	template, err := FetchInstanceTemplateE(t, ctx, projectID, name)
	require.NoError(t, err)

	return template
}

// FetchInstanceTemplateE queries GCP to return the settings it holds for the given global
// instance template.
// The ctx parameter supports cancellation and timeouts.
func FetchInstanceTemplateE(t testing.TestingT, ctx context.Context, projectID string, name string) (*compute.InstanceTemplate, error) {
	logger.Default.Logf(t, "Getting global instance template %s", name)

	service, err := NewComputeServiceContextE(t, ctx)
	if err != nil {
		return nil, err
	}

	return FetchInstanceTemplateWithClient(ctx, service, projectID, name)
}

// FetchInstanceTemplateWithClient queries GCP to return the settings it holds for the given global
// instance template using the supplied *compute.Service. Prefer this variant in unit tests where
// the service is backed by an httptest fake server (see instancetemplate_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func FetchInstanceTemplateWithClient(ctx context.Context, service *compute.Service, projectID string, name string) (*compute.InstanceTemplate, error) {
	template, err := service.InstanceTemplates.Get(projectID, name).Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("InstanceTemplates.Get(%s, %s) got error: %w", projectID, name, err)
	}

	return template, nil
}

// FetchRegionalInstanceGroupManager queries GCP to return the settings it holds for the given
// regional managed instance group, so a test can assert on what was actually created rather than
// only that it exists. It returns the manager's own settings, such as its template, target size and
// update policy. The instances it runs are read with FetchRegionalInstanceGroupContext.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func FetchRegionalInstanceGroupManager(t testing.TestingT, ctx context.Context, projectID string, region string, name string) *compute.InstanceGroupManager {
	manager, err := FetchRegionalInstanceGroupManagerE(t, ctx, projectID, region, name)
	require.NoError(t, err)

	return manager
}

// FetchRegionalInstanceGroupManagerE queries GCP to return the settings it holds for the given
// regional managed instance group.
// The ctx parameter supports cancellation and timeouts.
func FetchRegionalInstanceGroupManagerE(t testing.TestingT, ctx context.Context, projectID string, region string, name string) (*compute.InstanceGroupManager, error) {
	logger.Default.Logf(t, "Getting regional managed instance group %s in region %s", name, region)

	service, err := NewComputeServiceContextE(t, ctx)
	if err != nil {
		return nil, err
	}

	return FetchRegionalInstanceGroupManagerWithClient(ctx, service, projectID, region, name)
}

// FetchRegionalInstanceGroupManagerWithClient queries GCP to return the settings it holds for the
// given regional managed instance group using the supplied *compute.Service. Prefer this variant in
// unit tests where the service is backed by an httptest fake server (see instancetemplate_test.go
// for the pattern).
// The ctx parameter supports cancellation and timeouts.
func FetchRegionalInstanceGroupManagerWithClient(ctx context.Context, service *compute.Service, projectID string, region string, name string) (*compute.InstanceGroupManager, error) {
	manager, err := service.RegionInstanceGroupManagers.Get(projectID, region, name).Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("RegionInstanceGroupManagers.Get(%s, %s, %s) got error: %w", projectID, region, name, err)
	}

	return manager, nil
}

// FetchRegionInstanceTemplate queries GCP to return the settings it holds for the given regional instance template, so a test can
// assert on what was actually created rather than only that it exists. A regional template is stored in one region rather than globally, so it is read by region.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func FetchRegionInstanceTemplate(t testing.TestingT, ctx context.Context, projectID string, region string, name string) *compute.InstanceTemplate {
	template, err := FetchRegionInstanceTemplateE(t, ctx, projectID, region, name)
	require.NoError(t, err)

	return template
}

// FetchRegionInstanceTemplateE queries GCP to return the settings it holds for the given regional instance template.
// The ctx parameter supports cancellation and timeouts.
func FetchRegionInstanceTemplateE(t testing.TestingT, ctx context.Context, projectID string, region string, name string) (*compute.InstanceTemplate, error) {
	logger.Default.Logf(t, "Getting regional instance template %s in region %s", name, region)

	service, err := NewComputeServiceContextE(t, ctx)
	if err != nil {
		return nil, err
	}

	return FetchRegionInstanceTemplateWithClient(ctx, service, projectID, region, name)
}

// FetchRegionInstanceTemplateWithClient queries GCP to return the settings it holds for the given regional instance template using the
// supplied *compute.Service. Prefer this variant in unit tests where the service is backed by an
// httptest fake server (see instancetemplate_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func FetchRegionInstanceTemplateWithClient(ctx context.Context, service *compute.Service, projectID string, region string, name string) (*compute.InstanceTemplate, error) {
	template, err := service.RegionInstanceTemplates.Get(projectID, region, name).Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("RegionInstanceTemplates.Get(%s, %s, %s) got error: %w", projectID, region, name, err)
	}

	return template, nil
}

// FetchInstanceGroupManager queries GCP to return the settings it holds for the given zonal managed instance group, so a test can
// assert on what was actually created rather than only that it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func FetchInstanceGroupManager(t testing.TestingT, ctx context.Context, projectID string, zone string, name string) *compute.InstanceGroupManager {
	manager, err := FetchInstanceGroupManagerE(t, ctx, projectID, zone, name)
	require.NoError(t, err)

	return manager
}

// FetchInstanceGroupManagerE queries GCP to return the settings it holds for the given zonal managed instance group.
// The ctx parameter supports cancellation and timeouts.
func FetchInstanceGroupManagerE(t testing.TestingT, ctx context.Context, projectID string, zone string, name string) (*compute.InstanceGroupManager, error) {
	logger.Default.Logf(t, "Getting zonal managed instance group %s in zone %s", name, zone)

	service, err := NewComputeServiceContextE(t, ctx)
	if err != nil {
		return nil, err
	}

	return FetchInstanceGroupManagerWithClient(ctx, service, projectID, zone, name)
}

// FetchInstanceGroupManagerWithClient queries GCP to return the settings it holds for the given zonal managed instance group using the
// supplied *compute.Service. Prefer this variant in unit tests where the service is backed by an
// httptest fake server (see instancetemplate_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func FetchInstanceGroupManagerWithClient(ctx context.Context, service *compute.Service, projectID string, zone string, name string) (*compute.InstanceGroupManager, error) {
	manager, err := service.InstanceGroupManagers.Get(projectID, zone, name).Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("InstanceGroupManagers.Get(%s, %s, %s) got error: %w", projectID, zone, name, err)
	}

	return manager, nil
}

// FetchAutoscaler queries GCP to return the settings it holds for the given autoscaler, so a test can
// assert on what was actually created rather than only that it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func FetchAutoscaler(t testing.TestingT, ctx context.Context, projectID string, zone string, name string) *compute.Autoscaler {
	autoscaler, err := FetchAutoscalerE(t, ctx, projectID, zone, name)
	require.NoError(t, err)

	return autoscaler
}

// FetchAutoscalerE queries GCP to return the settings it holds for the given autoscaler.
// The ctx parameter supports cancellation and timeouts.
func FetchAutoscalerE(t testing.TestingT, ctx context.Context, projectID string, zone string, name string) (*compute.Autoscaler, error) {
	logger.Default.Logf(t, "Getting autoscaler %s in zone %s", name, zone)

	service, err := NewComputeServiceContextE(t, ctx)
	if err != nil {
		return nil, err
	}

	return FetchAutoscalerWithClient(ctx, service, projectID, zone, name)
}

// FetchAutoscalerWithClient queries GCP to return the settings it holds for the given autoscaler using the
// supplied *compute.Service. Prefer this variant in unit tests where the service is backed by an
// httptest fake server (see instancetemplate_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func FetchAutoscalerWithClient(ctx context.Context, service *compute.Service, projectID string, zone string, name string) (*compute.Autoscaler, error) {
	autoscaler, err := service.Autoscalers.Get(projectID, zone, name).Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("Autoscalers.Get(%s, %s, %s) got error: %w", projectID, zone, name, err)
	}

	return autoscaler, nil
}
