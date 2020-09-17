package docker

import (
	"strconv"

	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/shell"
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// StopOptions defines the options that can be passed to the 'docker stop' command
type StopOptions struct {
	// Seconds to wait for stop before killing the container (default 10)
	Time int

	// Set a logger that should be used. See the logger package for more info.
	Logger *logger.Logger
}

// Stop runs the 'docker stop' command for the given containers and return the stdout/stderr. This method fails
// the test if there are any errors
func Stop(t testing.TestingT, containers []string, options *StopOptions) string {
	out, err := StopE(t, containers, options)
	require.NoError(t, err)
	return out
}

// StopE runs the 'docker stop' command for the given containers and returns any errors.
func StopE(t testing.TestingT, containers []string, options *StopOptions) (string, error) {
	options.Logger.Logf(t, "Running 'docker stop' on containers '%s'", containers)

	args, err := formatDockerStopArgs(containers, options)
	if err != nil {
		return "", err
	}

	cmd := shell.Command{
		Command: "docker",
		Args:    args,
		Logger:  options.Logger,
	}

	return shell.RunCommandAndGetOutputE(t, cmd)

}

// formatDockerStopArgs formats the arguments for the 'docker stop' command
func formatDockerStopArgs(containers []string, options *StopOptions) ([]string, error) {
	args := []string{"stop"}

	if options.Time != 0 {
		args = append(args, "--time", strconv.Itoa(options.Time))
	}

	args = append(args, containers...)

	return args, nil
}
