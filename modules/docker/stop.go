package docker

import (
	"context"
	"strconv"

	"github.com/gruntwork-io/terratest/modules/core/v2/logger"
	"github.com/gruntwork-io/terratest/modules/shell/v2"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
)

// StopOptions defines the options that can be passed to the 'docker stop' command
type StopOptions struct {
	// Set a logger that should be used. See the logger package for more info.
	Logger *logger.Logger

	// Seconds to wait for stop before killing the container (default 10)
	Time int
}

// Stop runs the 'docker stop' command for the given containers and return the stdout/stderr. This method fails
// the test if there are any errors
//
// Deprecated: Use [StopContext] instead.
func Stop(t testing.TestingT, containers []string, options *StopOptions) string {
	return StopContext(t, context.Background(), containers, options)
}

// StopContext runs the 'docker stop' command for the given containers and returns stdout/stderr. This method
// fails the test if there are any errors. The ctx parameter supports cancellation and timeouts.
func StopContext(t testing.TestingT, ctx context.Context, containers []string, options *StopOptions) string {
	out, err := StopContextE(t, ctx, containers, options)
	require.NoError(t, err)

	return out
}

// StopE runs the 'docker stop' command for the given containers and returns any errors.
//
// Deprecated: Use [StopContextE] instead.
func StopE(t testing.TestingT, containers []string, options *StopOptions) (string, error) {
	return StopContextE(t, context.Background(), containers, options)
}

// StopContextE runs the 'docker stop' command for the given containers and returns stdout/stderr, or any
// error. The ctx parameter supports cancellation and timeouts.
func StopContextE(t testing.TestingT, ctx context.Context, containers []string, options *StopOptions) (string, error) {
	options.Logger.Logf(t, "Running 'docker stop' on containers '%s'", containers)

	args := formatDockerStopArgs(containers, options)

	cmd := &shell.Command{
		Command: "docker",
		Args:    args,
		Logger:  options.Logger,
	}

	return shell.RunCommandContextAndGetOutputE(t, ctx, cmd)
}

// formatDockerStopArgs formats the arguments for the 'docker stop' command
func formatDockerStopArgs(containers []string, options *StopOptions) []string {
	args := []string{"stop"}

	if options.Time != 0 {
		args = append(args, "--time", strconv.Itoa(options.Time))
	}

	args = append(args, containers...)

	return args
}
