package docker

import (
	"context"

	"github.com/gruntwork-io/terratest/modules/core/v2/logger"
	"github.com/gruntwork-io/terratest/modules/shell"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
)

// RunOptions defines options that can be passed to the 'docker run' command.
type RunOptions struct {
	// Set a logger that should be used. See the logger package for more info.
	Logger *logger.Logger

	// Override the default ENTRYPOINT of the Docker image
	Entrypoint string

	// Assign a name to the container
	Name string

	// Set platform
	Platform string

	// Username or UID
	User string

	// Override the default COMMAND of the Docker image
	Command []string

	// Set environment variables
	EnvironmentVariables []string

	// Bind mount these volume(s) when running the container
	Volumes []string

	// Custom CLI options that will be passed as-is to the 'docker run' command. This is an "escape hatch" that allows
	// Terratest to not have to support every single command-line option offered by the 'docker run' command, and
	// solely focus on the most important ones.
	OtherOptions []string

	// If set to true, pass the --detach flag to 'docker run' to run the container in the background
	Detach bool

	// If set to true, pass the --init flag to 'docker run' to run an init inside the container that forwards signals
	// and reaps processes
	Init bool

	// If set to true, pass the --privileged flag to 'docker run' to give extended privileges to the container
	Privileged bool

	// If set to true, pass the --rm flag to 'docker run' to automatically remove the container when it exits
	Remove bool

	// If set to true, pass the -tty flag to 'docker run' to allocate a pseudo-TTY
	Tty bool
}

// Run runs the 'docker run' command on the given image with the given options and return stdout/stderr. This method
// fails the test if there are any errors.
//
// Deprecated: Use [RunContext] instead.
func Run(t testing.TestingT, image string, options *RunOptions) string {
	return RunContext(t, context.Background(), image, options)
}

// RunContext runs the 'docker run' command on the given image with the given options and returns stdout/stderr.
// This method fails the test if there are any errors. The ctx parameter supports cancellation and timeouts.
func RunContext(t testing.TestingT, ctx context.Context, image string, options *RunOptions) string {
	out, err := RunContextE(t, ctx, image, options)
	require.NoError(t, err)

	return out
}

// RunE runs the 'docker run' command on the given image with the given options and return stdout/stderr, or any error.
//
// Deprecated: Use [RunContextE] instead.
func RunE(t testing.TestingT, image string, options *RunOptions) (string, error) {
	return RunContextE(t, context.Background(), image, options)
}

// RunContextE runs the 'docker run' command on the given image with the given options and returns
// stdout/stderr, or any error. The ctx parameter supports cancellation and timeouts.
func RunContextE(t testing.TestingT, ctx context.Context, image string, options *RunOptions) (string, error) {
	options.Logger.Logf(t, "Running 'docker run' on image '%s'", image)

	args := formatDockerRunArgs(image, options)

	cmd := &shell.Command{
		Command: "docker",
		Args:    args,
		Logger:  options.Logger,
	}

	return shell.RunCommandContextAndGetOutputE(t, ctx, cmd)
}

// RunAndGetID runs the 'docker run' command on the given image with the given options and returns the container ID
// that is returned in stdout. This method fails the test if there are any errors.
//
// Deprecated: Use [RunAndGetIDContext] instead.
func RunAndGetID(t testing.TestingT, image string, options *RunOptions) string {
	return RunAndGetIDContext(t, context.Background(), image, options)
}

// RunAndGetIDContext runs the 'docker run' command on the given image with the given options and returns the
// container ID that is returned in stdout. This method fails the test if there are any errors. The ctx
// parameter supports cancellation and timeouts.
func RunAndGetIDContext(t testing.TestingT, ctx context.Context, image string, options *RunOptions) string {
	out, err := RunAndGetIDContextE(t, ctx, image, options)
	require.NoError(t, err)

	return out
}

// RunAndGetIDE runs the 'docker run' command on the given image with the given options and returns the container ID
// that is returned in stdout, or any error.
//
// Deprecated: Use [RunAndGetIDContextE] instead.
func RunAndGetIDE(t testing.TestingT, image string, options *RunOptions) (string, error) {
	return RunAndGetIDContextE(t, context.Background(), image, options)
}

// RunAndGetIDContextE runs the 'docker run' command on the given image with the given options and returns the
// container ID that is returned in stdout, or any error. The ctx parameter supports cancellation and timeouts.
func RunAndGetIDContextE(t testing.TestingT, ctx context.Context, image string, options *RunOptions) (string, error) {
	options.Logger.Logf(t, "Running 'docker run' on image '%s', returning stdout", image)

	args := formatDockerRunArgs(image, options)

	cmd := &shell.Command{
		Command: "docker",
		Args:    args,
		Logger:  options.Logger,
	}

	return shell.RunCommandContextAndGetStdOutE(t, ctx, cmd)
}

// formatDockerRunArgs formats the arguments for the 'docker run' command.
func formatDockerRunArgs(image string, options *RunOptions) []string {
	args := []string{"run"}

	if options.Detach {
		args = append(args, "--detach")
	}

	if options.Entrypoint != "" {
		args = append(args, "--entrypoint", options.Entrypoint)
	}

	for _, envVar := range options.EnvironmentVariables {
		args = append(args, "--env", envVar)
	}

	if options.Init {
		args = append(args, "--init")
	}

	if options.Name != "" {
		args = append(args, "--name", options.Name)
	}

	if options.Platform != "" {
		args = append(args, "--platform", options.Platform)
	}

	if options.Privileged {
		args = append(args, "--privileged")
	}

	if options.Remove {
		args = append(args, "--rm")
	}

	if options.Tty {
		args = append(args, "--tty")
	}

	if options.User != "" {
		args = append(args, "--user", options.User)
	}

	for _, volume := range options.Volumes {
		args = append(args, "--volume", volume)
	}

	args = append(args, options.OtherOptions...)
	args = append(args, image)
	args = append(args, options.Command...)

	return args
}
