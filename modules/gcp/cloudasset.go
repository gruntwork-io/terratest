package gcp

import (
	"context"
	"errors"
	"fmt"

	"github.com/gruntwork-io/terratest/modules/core/v2/logger"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/cloudasset/v1"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
)

// GetAssetFeedAttrs returns the settings Google Cloud holds for the given Cloud Asset feed, so a
// test can assert on which asset changes it reports and where it sends them. Google names a feed by
// project number rather than by project id, so the caller passes the number: a project id in its
// place asks for a feed that cannot exist.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetAssetFeedAttrs(t testing.TestingT, ctx context.Context, projectNumber string, feedID string) *cloudasset.Feed {
	feed, err := GetAssetFeedAttrsE(t, ctx, projectNumber, feedID)
	require.NoError(t, err)

	return feed
}

// GetAssetFeedAttrsE returns the settings Google Cloud holds for the given Cloud Asset feed.
// The ctx parameter supports cancellation and timeouts.
func GetAssetFeedAttrsE(t testing.TestingT, ctx context.Context, projectNumber string, feedID string) (*cloudasset.Feed, error) {
	logger.Default.Logf(t, "Getting settings for Cloud Asset feed %s in project %s", feedID, projectNumber)

	service, err := NewCloudAssetServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetAssetFeedAttrsWithClient(ctx, service, projectNumber, feedID)
}

// GetAssetFeedAttrsWithClient returns the settings Google Cloud holds for the given Cloud Asset feed
// using the supplied *cloudasset.Service. Prefer this variant in unit tests where the service is
// backed by an httptest fake server (see cloudasset_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetAssetFeedAttrsWithClient(ctx context.Context, service *cloudasset.Service, projectNumber string, feedID string) (*cloudasset.Feed, error) {
	name := fmt.Sprintf("projects/%s/feeds/%s", projectNumber, feedID)

	feed, err := service.Feeds.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Cloud Asset feed %s in project %s does not exist", feedID, projectNumber)
		}

		return nil, fmt.Errorf("failed to get settings for Cloud Asset feed %s in project %s: %w", feedID, projectNumber, err)
	}

	return feed, nil
}

// NewCloudAssetServiceE creates a Cloud Asset service authenticated the same way every other client
// in this module is.
// The ctx parameter supports cancellation and timeouts.
func NewCloudAssetServiceE(t testing.TestingT, ctx context.Context) (*cloudasset.Service, error) {
	return cloudasset.NewService(ctx, append(withOptions(), option.WithScopes(cloudasset.CloudPlatformScope))...)
}
