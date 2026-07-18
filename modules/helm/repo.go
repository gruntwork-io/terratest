package helm

import (
	"context"

	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
)

// AddRepoContext will setup the provided helm repository to the local helm client configuration. This will fail the
// test if there is an error. The ctx parameter supports cancellation and timeouts.
func AddRepoContext(t testing.TestingT, ctx context.Context, options *Options, repoName string, repoURL string) {
	require.NoError(t, AddRepoContextE(t, ctx, options, repoName, repoURL))
}

// AddRepoContextE will setup the provided helm repository to the local helm client configuration. The ctx parameter
// supports cancellation and timeouts.
func AddRepoContextE(t testing.TestingT, ctx context.Context, options *Options, repoName string, repoURL string) error {

	args := []string{"add", repoName, repoURL}

	if options.ExtraArgs != nil {
		if repoAddArgs, ok := options.ExtraArgs["repoAdd"]; ok {
			args = append(args, repoAddArgs...)
		}
	}

	_, err := RunHelmCommandAndGetOutputContextE(t, ctx, options, "repo", args...)

	return err
}

// RemoveRepoContext will remove the provided helm repository from the local helm client configuration. This will fail
// the test if there is an error. The ctx parameter supports cancellation and timeouts.
func RemoveRepoContext(t testing.TestingT, ctx context.Context, options *Options, repoName string) {
	require.NoError(t, RemoveRepoContextE(t, ctx, options, repoName))
}

// RemoveRepoContextE will remove the provided helm repository from the local helm client configuration. The ctx
// parameter supports cancellation and timeouts.
func RemoveRepoContextE(t testing.TestingT, ctx context.Context, options *Options, repoName string) error {
	_, err := RunHelmCommandAndGetOutputContextE(t, ctx, options, "repo", "remove", repoName)

	return err
}
