package podman

import (
	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/shell"
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// RemoveOptions defines the options that can be passed to the 'podman rm' command
type RemoveOptions struct {
	// Set a logger that should be used. See the logger package for more info.
	Logger *logger.Logger
}

// Remove runs the 'podman rm' command for the given containers and return the stdout/stderr. This method fails
// the test if there are any errors
func Remove(t testing.TestingT, containers []string, options *RemoveOptions) string {
	out, err := RemoveE(t, containers, options)
	require.NoError(t, err)
	return out
}

// RemoveE runs the 'podman rm' command for the given containers and returns any errors.
func RemoveE(t testing.TestingT, containers []string, options *RemoveOptions) (string, error) {
	options.Logger.Logf(t, "Running 'podman rm' on containers '%s'", containers)

	args, err := formatDockerRemoveArgs(containers, options)
	if err != nil {
		return "", err
	}

	cmd := shell.Command{
		Command: "podman",
		Args:    args,
		Logger:  options.Logger,
	}

	return shell.RunCommandAndGetOutputE(t, cmd)
}

// formatDockerRemoveArgs formats the arguments for the 'podman rm' command
func formatDockerRemoveArgs(containers []string, options *RemoveOptions) ([]string, error) {
	args := []string{"rm"}

	args = append(args, containers...)

	return args, nil
}
