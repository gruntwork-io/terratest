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
	"google.golang.org/api/sourcerepo/v1"
)

// GetSourceRepoRepositoryAttrs returns the settings Google Cloud holds for the given Cloud Source Repositories repository, so a test can assert on
// what was actually created rather than only that it exists.
// The repository's size and its clone URL are both filled in by Google, so a read is the only way to
// see them.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetSourceRepoRepositoryAttrs(t testing.TestingT, ctx context.Context, projectID string, repoID string) *sourcerepo.Repo {
	result, err := GetSourceRepoRepositoryAttrsE(t, ctx, projectID, repoID)
	require.NoError(t, err)

	return result
}

// GetSourceRepoRepositoryAttrsE returns the settings Google Cloud holds for the given Cloud Source Repositories repository.
// The ctx parameter supports cancellation and timeouts.
func GetSourceRepoRepositoryAttrsE(t testing.TestingT, ctx context.Context, projectID string, repoID string) (*sourcerepo.Repo, error) {
	logger.Default.Logf(t, "Getting settings for source repository %s in project %s", repoID, projectID)

	service, err := NewSourceRepoServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetSourceRepoRepositoryAttrsWithClient(ctx, service, projectID, repoID)
}

// GetSourceRepoRepositoryAttrsWithClient returns the settings Google Cloud holds for the given Cloud Source Repositories repository using the supplied
// *sourcerepo.Service. Prefer this variant in unit tests where the service is backed by an httptest fake server
// (see sourcerepo_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetSourceRepoRepositoryAttrsWithClient(ctx context.Context, service *sourcerepo.Service, projectID string, repoID string) (*sourcerepo.Repo, error) {
	name := fmt.Sprintf("projects/%s/repos/%s", projectID, repoID)

	result, err := service.Projects.Repos.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the source repository %s does not exist in project %s", repoID, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for source repository %s in project %s: %w", repoID, projectID, err)
	}

	return result, nil
}

// NewSourceRepoServiceE creates a Cloud Source Repositories service authenticated the same way every other client in
// this module is.
// The ctx parameter supports cancellation and timeouts.
func NewSourceRepoServiceE(t testing.TestingT, ctx context.Context) (*sourcerepo.Service, error) {
	return sourcerepo.NewService(ctx, append(withOptions(), option.WithScopes(sourcerepo.CloudPlatformScope))...)
}
