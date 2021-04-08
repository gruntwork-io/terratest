package docker

import (
	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/shell"
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// BuildOptions defines options that can be passed to the 'docker build' command.
type BuildOptions struct {
	// Tags for the Docker image
	Tags []string

	// Build args to pass the 'docker build' command
	BuildArgs []string

	// Target build arg to pass to the 'docker build' command
	Target string

	// Custom CLI options that will be passed as-is to the 'docker build' command. This is an "escape hatch" that allows
	// Terratest to not have to support every single command-line option offered by the 'docker build' command, and
	// solely focus on the most important ones.
	OtherOptions []string

	// Set a logger that should be used. See the logger package for more info.
	Logger *logger.Logger
}

// Build runs the 'docker build' command at the given path with the given options and fails the test if there are any
// errors.
func Build(t testing.TestingT, path string, options *BuildOptions) {
	require.NoError(t, BuildE(t, path, options))
}

// BuildE runs the 'docker build' command at the given path with the given options and returns any errors.
func BuildE(t testing.TestingT, path string, options *BuildOptions) error {
	options.Logger.Logf(t, "Running 'docker build' in %s", path)

	args, err := formatDockerBuildArgs(path, options)
	if err != nil {
		return err
	}

	cmd := shell.Command{
		Command: "docker",
		Args:    args,
		Logger:  options.Logger,
	}

	_, buildErr := shell.RunCommandAndGetOutputE(t, cmd)
	return buildErr
}

// formatDockerBuildArgs formats the arguments for the 'docker build' command.
func formatDockerBuildArgs(path string, options *BuildOptions) ([]string, error) {
	args := []string{"build"}

	for _, tag := range options.Tags {
		args = append(args, "--tag", tag)
	}

	for _, arg := range options.BuildArgs {
		args = append(args, "--build-arg", arg)
	}

	if len(options.Target) > 0 {
		args = append(args, "--target", options.Target)
	}

	args = append(args, options.OtherOptions...)

	args = append(args, path)

	return args, nil
}
