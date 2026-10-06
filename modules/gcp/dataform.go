package gcp

import (
	"context"
	"errors"
	"fmt"

	"github.com/gruntwork-io/terratest/modules/core/v2/logger"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
	dataform "google.golang.org/api/dataform/v1beta1"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
)

// GetDataformFolderAttrs returns the settings Google Cloud holds for the Dataform folder, so a test can assert
// on what was actually created rather than only that it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetDataformFolderAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, folderID string) *dataform.Folder {
	result, err := GetDataformFolderAttrsE(t, ctx, projectID, location, folderID)
	require.NoError(t, err)

	return result
}

// GetDataformFolderAttrsE returns the settings Google Cloud holds for the Dataform folder.
// The ctx parameter supports cancellation and timeouts.
func GetDataformFolderAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, folderID string) (*dataform.Folder, error) {
	logger.Default.Logf(t, "Getting settings for Dataform folder %s in project %s", folderID, projectID)

	service, err := NewDataformServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetDataformFolderAttrsWithClient(ctx, service, projectID, location, folderID)
}

// GetDataformFolderAttrsWithClient returns the settings Google Cloud holds for the Dataform folder using the
// supplied *dataform.Service. Prefer this variant in unit tests where the service is backed by an
// httptest fake server (see dataform_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetDataformFolderAttrsWithClient(ctx context.Context, service *dataform.Service, projectID string, location string, folderID string) (*dataform.Folder, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/folders/%s", projectID, location, folderID)

	result, err := service.Projects.Locations.Folders.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Dataform folder %s in project %s does not exist", folderID, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Dataform folder %s in project %s: %w", folderID, projectID, err)
	}

	return result, nil
}

// NewDataformServiceE creates a Dataform service authenticated the same way every other client in
// this module is.
// The ctx parameter supports cancellation and timeouts.
func NewDataformServiceE(t testing.TestingT, ctx context.Context) (*dataform.Service, error) {
	return dataform.NewService(ctx, append(withOptions(), option.WithScopes(dataform.CloudPlatformScope))...)
}
