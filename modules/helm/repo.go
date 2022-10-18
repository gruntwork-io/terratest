package helm

import (
	"github.com/stretchr/testify/require"

	"github.com/gruntwork-io/terratest/modules/testing"
)

// AddRepo will setup the provided helm repository to the local helm client configuration. This will fail the test if
// there is an error.
func AddRepo(t testing.TestingT, options *Options, repoName string, repoURL string) {
	require.NoError(t, AddRepoE(t, options, repoName, repoURL))
}

// AddRepoE will setup the provided helm repository to the local helm client configuration.
func AddRepoE(t testing.TestingT, options *Options, repoName string, repoURL string) error {
	// Set required args
	args := []string{"add", repoName, repoURL}

	// Append helm repo add ExtraArgs if available
	if options.ExtraArgs != nil {
		if repoAddArgs, ok := options.ExtraArgs["repoAdd"]; ok {
			args = append(args, repoAddArgs...)
		}
	}
	_, err := RunHelmCommandAndGetOutputE(t, options, "repo", args...)
	return err
}

// RemoveRepo will remove the provided helm repository from the local helm client configuration. This will fail the test
// if there is an error.
func RemoveRepo(t testing.TestingT, options *Options, repoName string) {
	require.NoError(t, RemoveRepoE(t, options, repoName))
}

// RemoveRepoE will remove the provided helm repository from the local helm client configuration.
func RemoveRepoE(t testing.TestingT, options *Options, repoName string) error {
	_, err := RunHelmCommandAndGetOutputE(t, options, "repo", "remove", repoName)
	return err
}
