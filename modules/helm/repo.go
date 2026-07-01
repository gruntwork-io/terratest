package helm

import (
	"context"

	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
)

// AddRepo will setup the provided helm repository to the local helm client configuration. This will fail the test if
// there is an error.
//
// Deprecated: Use [AddRepoContext] instead.
func AddRepo(t testing.TestingT, options *Options, repoName string, repoURL string) {
	AddRepoContext(t, context.Background(), options, repoName, repoURL)
}

// AddRepoContext will setup the provided helm repository to the local helm client configuration. This will fail the
// test if there is an error. The ctx parameter supports cancellation and timeouts.
func AddRepoContext(t testing.TestingT, ctx context.Context, options *Options, repoName string, repoURL string) {
	require.NoError(t, AddRepoContextE(t, ctx, options, repoName, repoURL))
}

// AddRepoE will setup the provided helm repository to the local helm client configuration.
//
// Deprecated: Use [AddRepoContextE] instead.
func AddRepoE(t testing.TestingT, options *Options, repoName string, repoURL string) error {
	return AddRepoContextE(t, context.Background(), options, repoName, repoURL)
}

// AddRepoContextE will setup the provided helm repository to the local helm client configuration. The ctx parameter
// supports cancellation and timeouts.
func AddRepoContextE(t testing.TestingT, ctx context.Context, options *Options, repoName string, repoURL string) error {
	// Set required args
	args := []string{"add", repoName, repoURL}

	// Append helm repo add ExtraArgs if available
	if options.ExtraArgs != nil {
		if repoAddArgs, ok := options.ExtraArgs["repoAdd"]; ok {
			args = append(args, repoAddArgs...)
		}
	}

	_, err := RunHelmCommandAndGetOutputContextE(t, ctx, options, "repo", args...)

	return err
}

// RemoveRepo will remove the provided helm repository from the local helm client configuration. This will fail the test
// if there is an error.
//
// Deprecated: Use [RemoveRepoContext] instead.
func RemoveRepo(t testing.TestingT, options *Options, repoName string) {
	RemoveRepoContext(t, context.Background(), options, repoName)
}

// RemoveRepoContext will remove the provided helm repository from the local helm client configuration. This will fail
// the test if there is an error. The ctx parameter supports cancellation and timeouts.
func RemoveRepoContext(t testing.TestingT, ctx context.Context, options *Options, repoName string) {
	require.NoError(t, RemoveRepoContextE(t, ctx, options, repoName))
}

// RemoveRepoE will remove the provided helm repository from the local helm client configuration.
//
// Deprecated: Use [RemoveRepoContextE] instead.
func RemoveRepoE(t testing.TestingT, options *Options, repoName string) error {
	return RemoveRepoContextE(t, context.Background(), options, repoName)
}

// RemoveRepoContextE will remove the provided helm repository from the local helm client configuration. The ctx
// parameter supports cancellation and timeouts.
func RemoveRepoContextE(t testing.TestingT, ctx context.Context, options *Options, repoName string) error {
	_, err := RunHelmCommandAndGetOutputContextE(t, ctx, options, "repo", "remove", repoName)

	return err
}
