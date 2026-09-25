package gcp

import (
	"context"
	"errors"
	"fmt"

	"github.com/gruntwork-io/terratest/modules/core/v2/logger"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/cloudresourcemanager/v3"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
)

// GetTagKeyAttrs returns the settings Google Cloud holds for the given tag key, so a test can assert
// on what was actually created rather than only that it exists. Google assigns a tag key its numeric
// id, so the caller passes the one it got back rather than a name it chose.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetTagKeyAttrs(t testing.TestingT, ctx context.Context, tagKeyID string) *cloudresourcemanager.TagKey {
	tagKey, err := GetTagKeyAttrsE(t, ctx, tagKeyID)
	require.NoError(t, err)

	return tagKey
}

// GetTagKeyAttrsE returns the settings Google Cloud holds for the given tag key.
// The ctx parameter supports cancellation and timeouts.
func GetTagKeyAttrsE(t testing.TestingT, ctx context.Context, tagKeyID string) (*cloudresourcemanager.TagKey, error) {
	logger.Default.Logf(t, "Getting settings for tag key %s", tagKeyID)

	service, err := NewResourceManagerServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetTagKeyAttrsWithClient(ctx, service, tagKeyID)
}

// GetTagKeyAttrsWithClient returns the settings Google Cloud holds for the given tag key using the
// supplied *cloudresourcemanager.Service. Prefer this variant in unit tests where the service is
// backed by an httptest fake server (see cloudresourcemanager_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetTagKeyAttrsWithClient(ctx context.Context, service *cloudresourcemanager.Service, tagKeyID string) (*cloudresourcemanager.TagKey, error) {
	tagKey, err := service.TagKeys.Get("tagKeys/" + tagKeyID).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the tag key %s does not exist", tagKeyID)
		}

		return nil, fmt.Errorf("failed to get settings for tag key %s: %w", tagKeyID, err)
	}

	return tagKey, nil
}

// GetTagValueAttrs returns the settings Google Cloud holds for the given tag value, so a test can
// assert on what was actually created rather than only that it exists. Google assigns a tag value its
// numeric id, so the caller passes the one it got back rather than a name it chose.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetTagValueAttrs(t testing.TestingT, ctx context.Context, tagValueID string) *cloudresourcemanager.TagValue {
	tagValue, err := GetTagValueAttrsE(t, ctx, tagValueID)
	require.NoError(t, err)

	return tagValue
}

// GetTagValueAttrsE returns the settings Google Cloud holds for the given tag value.
// The ctx parameter supports cancellation and timeouts.
func GetTagValueAttrsE(t testing.TestingT, ctx context.Context, tagValueID string) (*cloudresourcemanager.TagValue, error) {
	logger.Default.Logf(t, "Getting settings for tag value %s", tagValueID)

	service, err := NewResourceManagerServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetTagValueAttrsWithClient(ctx, service, tagValueID)
}

// GetTagValueAttrsWithClient returns the settings Google Cloud holds for the given tag value using
// the supplied *cloudresourcemanager.Service. Prefer this variant in unit tests where the service is
// backed by an httptest fake server (see cloudresourcemanager_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetTagValueAttrsWithClient(ctx context.Context, service *cloudresourcemanager.Service, tagValueID string) (*cloudresourcemanager.TagValue, error) {
	tagValue, err := service.TagValues.Get("tagValues/" + tagValueID).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the tag value %s does not exist", tagValueID)
		}

		return nil, fmt.Errorf("failed to get settings for tag value %s: %w", tagValueID, err)
	}

	return tagValue, nil
}

// GetTagBindingsAttrs returns the tag bindings Google Cloud holds for the given resource, so a test
// can assert on which values are actually attached to it. There is no call that reads one binding by
// name, so the resource's bindings are listed and the caller picks the one it made. The parent is a
// full resource name, such as //cloudresourcemanager.googleapis.com/projects/1234567890.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetTagBindingsAttrs(t testing.TestingT, ctx context.Context, parent string) []*cloudresourcemanager.TagBinding {
	bindings, err := GetTagBindingsAttrsE(t, ctx, parent)
	require.NoError(t, err)

	return bindings
}

// GetTagBindingsAttrsE returns the tag bindings Google Cloud holds for the given resource.
// The ctx parameter supports cancellation and timeouts.
func GetTagBindingsAttrsE(t testing.TestingT, ctx context.Context, parent string) ([]*cloudresourcemanager.TagBinding, error) {
	logger.Default.Logf(t, "Getting the tag bindings on %s", parent)

	service, err := NewResourceManagerServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetTagBindingsAttrsWithClient(ctx, service, parent)
}

// GetTagBindingsAttrsWithClient returns the tag bindings Google Cloud holds for the given resource
// using the supplied *cloudresourcemanager.Service. Prefer this variant in unit tests where the
// service is backed by an httptest fake server (see cloudresourcemanager_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetTagBindingsAttrsWithClient(ctx context.Context, service *cloudresourcemanager.Service, parent string) ([]*cloudresourcemanager.TagBinding, error) {
	response, err := service.TagBindings.List().Parent(parent).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("%s does not exist", parent)
		}

		return nil, fmt.Errorf("failed to list the tag bindings on %s: %w", parent, err)
	}

	return response.TagBindings, nil
}

// GetLienAttrs returns the settings Google Cloud holds for the given lien, so a test can assert on
// what was actually created rather than only that it exists. Google assigns a lien its id, and that
// is what the caller passes: the collection in front of it is added here.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetLienAttrs(t testing.TestingT, ctx context.Context, lienID string) *cloudresourcemanager.Lien {
	lien, err := GetLienAttrsE(t, ctx, lienID)
	require.NoError(t, err)

	return lien
}

// GetLienAttrsE returns the settings Google Cloud holds for the given lien.
// The ctx parameter supports cancellation and timeouts.
func GetLienAttrsE(t testing.TestingT, ctx context.Context, lienID string) (*cloudresourcemanager.Lien, error) {
	logger.Default.Logf(t, "Getting settings for lien %s", lienID)

	service, err := NewResourceManagerServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetLienAttrsWithClient(ctx, service, lienID)
}

// GetLienAttrsWithClient returns the settings Google Cloud holds for the given lien using the
// supplied *cloudresourcemanager.Service. Prefer this variant in unit tests where the service is
// backed by an httptest fake server (see cloudresourcemanager_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetLienAttrsWithClient(ctx context.Context, service *cloudresourcemanager.Service, lienID string) (*cloudresourcemanager.Lien, error) {
	lien, err := service.Liens.Get("liens/" + lienID).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the lien %s does not exist", lienID)
		}

		return nil, fmt.Errorf("failed to get settings for lien %s: %w", lienID, err)
	}

	return lien, nil
}

// NewResourceManagerServiceE creates a Cloud Resource Manager service authenticated the same way
// every other client in this module is. Tag keys, tag values, tag bindings and liens are all on it.
// The ctx parameter supports cancellation and timeouts.
func NewResourceManagerServiceE(t testing.TestingT, ctx context.Context) (*cloudresourcemanager.Service, error) {
	return cloudresourcemanager.NewService(ctx, append(withOptions(), option.WithScopes(cloudresourcemanager.CloudPlatformScope))...)
}
