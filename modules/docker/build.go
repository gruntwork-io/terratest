package docker

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/gruntwork-io/terratest/modules/core/v2/logger"
	"github.com/gruntwork-io/terratest/modules/shell/v2"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/hashicorp/go-multierror"
	"github.com/stretchr/testify/require"
)

// BuildOptions defines options that can be passed to the 'docker build' command.
type BuildOptions struct {
	// Additional environment variables to pass in when running docker build command.
	Env map[string]string

	// Set a logger that should be used. See the logger package for more info.
	Logger *logger.Logger

	// Target build arg to pass to the 'docker build' command
	Target string

	// Tags for the Docker image
	Tags []string

	// Build args to pass the 'docker build' command
	BuildArgs []string

	// All architectures to target in a multiarch build. Configuring this variable will cause terratest to use docker
	// buildx to construct multiarch images.
	// You can read more about multiarch docker builds in the official documentation for buildx:
	// https://docs.docker.com/buildx/working-with-buildx/
	// NOTE: This list does not automatically include the current platform. For example, if you are building images on
	// an Apple Silicon based MacBook, and you configure this variable to []string{"linux/amd64"} to build an amd64
	// image, the buildx command will not automatically include linux/arm64 - you must include that explicitly.
	Architectures []string

	// Custom CLI options that will be passed as-is to the 'docker build' command. This is an "escape hatch" that allows
	// Terratest to not have to support every single command-line option offered by the 'docker build' command, and
	// solely focus on the most important ones.
	OtherOptions []string

	// Whether or not to push images directly to the registry on build. Note that for multiarch images (Architectures is
	// not empty), this must be true to ensure availability of all architectures - only the image for the current
	// platform will be loaded into the daemon (due to a limitation of the docker daemon), so you won't be able to run a
	// `docker push` command later to push the multiarch image.
	// See https://github.com/moby/moby/pull/38738 for more info on the limitation of multiarch images in docker daemon.
	Push bool

	// Whether or not to load the image into the docker daemon at the end of a multiarch build so that it can be used
	// locally. Note that this is only used when Architectures is set, and assumes the current architecture is already
	// included in the Architectures list.
	Load bool

	// Whether ot not to enable buildkit. You can find more information about buildkit here https://docs.docker.com/build/buildkit/#getting-started.
	EnableBuildKit bool
}

// Build runs the 'docker build' command at the given path with the given options and fails the test if there are any
// errors.
//
// Deprecated: Use [BuildContext] instead.
func Build(t testing.TestingT, path string, options *BuildOptions) {
	BuildContext(t, context.Background(), path, options)
}

// BuildContext runs the 'docker build' command at the given path with the given options and fails the test if
// there are any errors. The ctx parameter supports cancellation and timeouts.
func BuildContext(t testing.TestingT, ctx context.Context, path string, options *BuildOptions) {
	require.NoError(t, BuildContextE(t, ctx, path, options))
}

// BuildE runs the 'docker build' command at the given path with the given options and returns any errors.
//
// Deprecated: Use [BuildContextE] instead.
func BuildE(t testing.TestingT, path string, options *BuildOptions) error {
	return BuildContextE(t, context.Background(), path, options)
}

// BuildContextE runs the 'docker build' command at the given path with the given options and returns any errors.
// The ctx parameter supports cancellation and timeouts.
func BuildContextE(t testing.TestingT, ctx context.Context, path string, options *BuildOptions) error {
	options.Logger.Logf(t, "Running 'docker build' in %s", path)

	env := make(map[string]string)
	if options.Env != nil {
		env = options.Env
	}

	if options.EnableBuildKit {
		env["DOCKER_BUILDKIT"] = "1"
	}

	cmd := &shell.Command{
		Command: "docker",
		Args:    formatDockerBuildArgs(path, options),
		Logger:  options.Logger,
		Env:     env,
	}

	if err := shell.RunCommandContextE(t, ctx, cmd); err != nil {
		return err
	}

	// For non multiarch images, we need to call docker push for each tag since build does not have a push option like
	// buildx.
	if len(options.Architectures) == 0 && options.Push {
		errorsOccurred := new(multierror.Error)

		for _, tag := range options.Tags {
			if err := PushContextE(t, ctx, options.Logger, tag); err != nil {
				options.Logger.Logf(t, "ERROR: error pushing tag %s", tag)

				errorsOccurred = multierror.Append(errorsOccurred, err)
			}
		}

		return errorsOccurred.ErrorOrNil()
	}

	// For multiarch images, if a load is requested call the load command to export the built image into the daemon.
	if len(options.Architectures) > 0 && options.Load {
		loadCmd := &shell.Command{
			Command: "docker",
			Args:    formatDockerBuildxLoadArgs(path, options),
			Logger:  options.Logger,
		}

		return shell.RunCommandContextE(t, ctx, loadCmd)
	}

	return nil
}

// GitCloneAndBuild builds a new Docker image from a given Git repo. This function will clone the given repo at the
// specified ref, and call the docker build command on the cloned repo from the given relative path (relative to repo
// root). This will fail the test if there are any errors.
//
// Deprecated: Use [GitCloneAndBuildContext] instead.
func GitCloneAndBuild(
	t testing.TestingT,
	repo string,
	ref string,
	path string,
	dockerBuildOpts *BuildOptions,
) {
	GitCloneAndBuildContext(t, context.Background(), repo, ref, path, dockerBuildOpts)
}

// GitCloneAndBuildContext builds a new Docker image from a given Git repo. This function will clone the given
// repo at the specified ref, and call the docker build command on the cloned repo from the given relative path
// (relative to repo root). This will fail the test if there are any errors. The ctx parameter supports
// cancellation and timeouts.
func GitCloneAndBuildContext(
	t testing.TestingT,
	ctx context.Context,
	repo string,
	ref string,
	path string,
	dockerBuildOpts *BuildOptions,
) {
	require.NoError(t, GitCloneAndBuildContextE(t, ctx, repo, ref, path, dockerBuildOpts))
}

// GitCloneAndBuildE builds a new Docker image from a given Git repo. This function will clone the given repo at the
// specified ref, and call the docker build command on the cloned repo from the given relative path (relative to repo
// root).
//
// Deprecated: Use [GitCloneAndBuildContextE] instead.
func GitCloneAndBuildE(
	t testing.TestingT,
	repo string,
	ref string,
	path string,
	dockerBuildOpts *BuildOptions,
) error {
	return GitCloneAndBuildContextE(t, context.Background(), repo, ref, path, dockerBuildOpts)
}

// GitCloneAndBuildContextE builds a new Docker image from a given Git repo. This function will clone the given
// repo at the specified ref, and call the docker build command on the cloned repo from the given relative path
// (relative to repo root). The ctx parameter supports cancellation and timeouts.
func GitCloneAndBuildContextE(
	t testing.TestingT,
	ctx context.Context,
	repo string,
	ref string,
	path string,
	dockerBuildOpts *BuildOptions,
) error {
	workingDir, err := os.MkdirTemp("", "")
	if err != nil {
		return err
	}

	defer func() {
		if removeErr := os.RemoveAll(workingDir); removeErr != nil {
			dockerBuildOpts.Logger.Logf(t, "WARNING: failed to remove temp dir %s: %v", workingDir, removeErr)
		}
	}()

	cloneCmd := &shell.Command{
		Command: "git",
		Args:    []string{"clone", repo, workingDir},
	}

	if err := shell.RunCommandContextE(t, ctx, cloneCmd); err != nil {
		return err
	}

	checkoutCmd := &shell.Command{
		Command:    "git",
		Args:       []string{"checkout", ref},
		WorkingDir: workingDir,
	}

	if err := shell.RunCommandContextE(t, ctx, checkoutCmd); err != nil {
		return err
	}

	contextPath := filepath.Join(workingDir, path)

	return BuildContextE(t, ctx, contextPath, dockerBuildOpts)
}

// formatDockerBuildArgs formats the arguments for the 'docker build' command.
func formatDockerBuildArgs(path string, options *BuildOptions) []string {
	args := []string{}

	if len(options.Architectures) > 0 {
		args = append(
			args,
			"buildx",
			"build",
			"--platform",
			strings.Join(options.Architectures, ","),
		)

		if options.Push {
			args = append(args, "--push")
		}
	} else {
		args = append(args, "build")
	}

	return append(args, formatDockerBuildBaseArgs(path, options)...)
}

// formatDockerBuildxLoadArgs formats the arguments for calling load on the 'docker buildx' command.
func formatDockerBuildxLoadArgs(path string, options *BuildOptions) []string {
	base := formatDockerBuildBaseArgs(path, options)

	args := make([]string, 0, len(base)+3) //nolint:mnd // 3 = "buildx" + "build" + "--load"
	args = append(args, "buildx", "build", "--load")

	return append(args, base...)
}

// formatDockerBuildBaseArgs formats the common args for the build command, both for `build` and `buildx`.
func formatDockerBuildBaseArgs(path string, options *BuildOptions) []string {
	//nolint:mnd // 2 = pairs of (--flag, value); 3 = max of --target flag + value + path
	args := make([]string, 0, len(options.Tags)*2+len(options.BuildArgs)*2+len(options.OtherOptions)+3)

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

	return args
}
