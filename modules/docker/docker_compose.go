package docker

import (
	"context"
	"regexp"
	"strings"

	"github.com/gruntwork-io/terratest/modules/core/v2/logger"
	"github.com/gruntwork-io/terratest/modules/shell"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
	"gotest.tools/v3/icmd"
)

// Options are Docker options.
type Options struct {
	EnvVars    map[string]string
	Logger     *logger.Logger
	WorkingDir string

	ProjectName string

	// Whether ot not to enable buildkit. You can find more information about buildkit here https://docs.docker.com/build/buildkit/#getting-started.
	EnableBuildKit bool
}

// RunDockerCompose runs docker compose with the given arguments and options and return stdout/stderr.
//
// Deprecated: Use [RunDockerComposeContext] instead.
func RunDockerCompose(t testing.TestingT, options *Options, args ...string) string {
	return RunDockerComposeContext(t, context.Background(), options, args...)
}

// RunDockerComposeContext runs docker compose with the given arguments and options and returns stdout/stderr.
// This method fails the test if there are any errors. The ctx parameter supports cancellation and timeouts.
func RunDockerComposeContext(t testing.TestingT, ctx context.Context, options *Options, args ...string) string {
	out, err := runDockerComposeE(t, ctx, false, options, args...)
	if err != nil {
		t.Fatal(err)
	}

	return out
}

// RunDockerComposeAndGetStdOut runs docker compose with the given arguments and options and returns only stdout.
//
// Deprecated: Use [RunDockerComposeAndGetStdOutContext] instead.
func RunDockerComposeAndGetStdOut(t testing.TestingT, options *Options, args ...string) string {
	return RunDockerComposeAndGetStdOutContext(t, context.Background(), options, args...)
}

// RunDockerComposeAndGetStdOutContext runs docker compose with the given arguments and options and returns only
// stdout. This method fails the test if there are any errors. The ctx parameter supports cancellation and
// timeouts.
func RunDockerComposeAndGetStdOutContext(t testing.TestingT, ctx context.Context, options *Options, args ...string) string {
	out, err := runDockerComposeE(t, ctx, true, options, args...)
	require.NoError(t, err)

	return out
}

// RunDockerComposeE runs docker compose with the given arguments and options and return stdout/stderr.
//
// Deprecated: Use [RunDockerComposeContextE] instead.
func RunDockerComposeE(t testing.TestingT, options *Options, args ...string) (string, error) {
	return RunDockerComposeContextE(t, context.Background(), options, args...)
}

// RunDockerComposeContextE runs docker compose with the given arguments and options and returns stdout/stderr,
// or any error. The ctx parameter supports cancellation and timeouts.
func RunDockerComposeContextE(t testing.TestingT, ctx context.Context, options *Options, args ...string) (string, error) {
	return runDockerComposeE(t, ctx, false, options, args...)
}

func runDockerComposeE(t testing.TestingT, ctx context.Context, stdout bool, options *Options, args ...string) (string, error) {
	var cmd *shell.Command

	projectName := options.ProjectName
	if len(projectName) == 0 {
		projectName = strings.ToLower(t.Name())
	}

	dockerComposeVersionCmd := icmd.Command("docker", "compose", "version")
	result := icmd.RunCmd(dockerComposeVersionCmd)

	if options.EnableBuildKit {
		if options.EnvVars == nil {
			options.EnvVars = make(map[string]string)
		}

		options.EnvVars["DOCKER_BUILDKIT"] = "1"
		options.EnvVars["COMPOSE_DOCKER_CLI_BUILD"] = "1"
	}

	if result.ExitCode == 0 {
		cmd = &shell.Command{
			Command:    "docker",
			Args:       append([]string{"compose", "--project-name", generateValidDockerComposeProjectName(projectName)}, args...),
			WorkingDir: options.WorkingDir,
			Env:        options.EnvVars,
			Logger:     options.Logger,
		}
	} else {
		cmd = &shell.Command{
			Command: "docker-compose",
			// We append --project-name to ensure containers from multiple different tests using Docker Compose don't end
			// up in the same project and end up conflicting with each other.
			Args:       append([]string{"--project-name", generateValidDockerComposeProjectName(projectName)}, args...),
			WorkingDir: options.WorkingDir,
			Env:        options.EnvVars,
			Logger:     options.Logger,
		}
	}

	if stdout {
		return shell.RunCommandContextAndGetStdOut(t, ctx, cmd), nil
	}

	return shell.RunCommandContextAndGetOutputE(t, ctx, cmd)
}

// generateValidDockerComposeProjectName generates a valid project name for docker-compose.
// Note: docker-compose command doesn't like lower case or special characters, other than -.
func generateValidDockerComposeProjectName(str string) string {
	lowerStr := strings.ToLower(str)

	return regexp.MustCompile(`[^a-zA-Z0-9 ]+`).ReplaceAllString(lowerStr, "-")
}
