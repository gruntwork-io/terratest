package gcp

import (
	"context"
	"errors"
	"fmt"

	"github.com/gruntwork-io/terratest/modules/core/v2/logger"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/netapp/v1"
	"google.golang.org/api/option"
)

// GetNetAppHostGroupAttrs returns the settings Google Cloud holds for the given NetApp host group, so a test can assert on
// what was actually created rather than only that it exists. A host group names the initiators a volume may be exported to; it exports nothing by itself.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetNetAppHostGroupAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, groupID string) *netapp.HostGroup {
	result, err := GetNetAppHostGroupAttrsE(t, ctx, projectID, location, groupID)
	require.NoError(t, err)

	return result
}

// GetNetAppHostGroupAttrsE returns the settings Google Cloud holds for the given NetApp host group.
// The ctx parameter supports cancellation and timeouts.
func GetNetAppHostGroupAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, groupID string) (*netapp.HostGroup, error) {
	logger.Default.Logf(t, "Getting settings for NetApp host group %s in project %s", groupID, projectID)

	service, err := NewNetAppServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetNetAppHostGroupAttrsWithClient(ctx, service, projectID, location, groupID)
}

// GetNetAppHostGroupAttrsWithClient returns the settings Google Cloud holds for the given NetApp host group using the
// supplied *netapp.Service. Prefer this variant in unit tests where the service is backed by an
// httptest fake server (see netapp_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetNetAppHostGroupAttrsWithClient(ctx context.Context, service *netapp.Service, projectID string, location string, groupID string) (*netapp.HostGroup, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/hostGroups/%s", projectID, location, groupID)

	result, err := service.Projects.Locations.HostGroups.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the NetApp host group %s in project %s does not exist", groupID, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for NetApp host group %s in project %s: %w", groupID, projectID, err)
	}

	return result, nil
}

// GetNetAppKmsConfigAttrs returns the settings Google Cloud holds for the given NetApp KMS config, so a test can assert on
// what was actually created rather than only that it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetNetAppKmsConfigAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, configID string) *netapp.KmsConfig {
	result, err := GetNetAppKmsConfigAttrsE(t, ctx, projectID, location, configID)
	require.NoError(t, err)

	return result
}

// GetNetAppKmsConfigAttrsE returns the settings Google Cloud holds for the given NetApp KMS config.
// The ctx parameter supports cancellation and timeouts.
func GetNetAppKmsConfigAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, configID string) (*netapp.KmsConfig, error) {
	logger.Default.Logf(t, "Getting settings for NetApp KMS config %s in project %s", configID, projectID)

	service, err := NewNetAppServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetNetAppKmsConfigAttrsWithClient(ctx, service, projectID, location, configID)
}

// GetNetAppKmsConfigAttrsWithClient returns the settings Google Cloud holds for the given NetApp KMS config using the
// supplied *netapp.Service. Prefer this variant in unit tests where the service is backed by an
// httptest fake server (see netapp_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetNetAppKmsConfigAttrsWithClient(ctx context.Context, service *netapp.Service, projectID string, location string, configID string) (*netapp.KmsConfig, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/kmsConfigs/%s", projectID, location, configID)

	result, err := service.Projects.Locations.KmsConfigs.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the NetApp KMS config %s in project %s does not exist", configID, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for NetApp KMS config %s in project %s: %w", configID, projectID, err)
	}

	return result, nil
}

// NewNetAppServiceE creates a NetApp Volumes service authenticated the same way every other client in this
// module is.
// The ctx parameter supports cancellation and timeouts.
func NewNetAppServiceE(t testing.TestingT, ctx context.Context) (*netapp.Service, error) {
	return netapp.NewService(ctx, append(withOptions(), option.WithScopes(netapp.CloudPlatformScope))...)
}
