package podman

import (
	"fmt"

	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/shell"
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// StartOptions defines the options that can be passed to the 'podman start' command
type StartOptions struct {
	// Set a logger that should be used. See the logger package for more info.
	Logger *logger.Logger
}

// Start runs the 'podman start' command for the given containers and return the stdout/stderr. This method fails
// the test if there are any errors
func Start(t testing.TestingT, containers []string, options *StartOptions) string {
	out, err := StartE(t, containers, options)
	require.NoError(t, err)
	return out
}

// StartE runs the 'podman start' command for the given containers and returns any errors.
func StartE(t testing.TestingT, containers []string, options *StartOptions) (string, error) {
	options.Logger.Logf(t, "Running 'podman start' on containers '%s'", containers)

	args, err := formatPodmanStartArgs(containers, options)
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

// formatPodmanStartArgs formats the arguments for the 'podman start' command
func formatPodmanStartArgs(containers []string, options *StartOptions) ([]string, error) {
	args := []string{"start"}

	args = append(args, containers...)
	fmt.Println(containers)
	fmt.Println(args)
	return args, nil
}
