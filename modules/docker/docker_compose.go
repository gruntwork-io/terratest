package docker

import (
	"regexp"
	"strings"

	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/shell"
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
	"gotest.tools/v3/icmd"
)

// Options are Docker options.
type Options struct {
	WorkingDir string
	EnvVars    map[string]string

	// Whether ot not to enable buildkit. You can find more information about buildkit here https://docs.docker.com/build/buildkit/#getting-started.
	EnableBuildKit bool

	// Set a logger that should be used. See the logger package for more info.
	Logger      *logger.Logger
	ProjectName string
}

// RunDockerCompose runs docker compose with the given arguments and options and return stdout/stderr.
func RunDockerCompose(t testing.TestingT, options *Options, args ...string) string {
	out, err := runDockerComposeE(t, false, options, args...)
	if err != nil {
		t.Fatal(err)
	}
	return out
}

// RunDockerComposeAndGetStdout runs docker compose with the given arguments and options and returns only stdout.
func RunDockerComposeAndGetStdOut(t testing.TestingT, options *Options, args ...string) string {
	out, err := runDockerComposeE(t, true, options, args...)
	require.NoError(t, err)
	return out
}

// RunDockerComposeE runs docker compose with the given arguments and options and return stdout/stderr.
func RunDockerComposeE(t testing.TestingT, options *Options, args ...string) (string, error) {
	return runDockerComposeE(t, false, options, args...)
}

func runDockerComposeE(t testing.TestingT, stdout bool, options *Options, args ...string) (string, error) {
	var cmd shell.Command

	projectName := options.ProjectName
	if len(projectName) <= 0 {
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
		cmd = shell.Command{
			Command:    "docker",
			Args:       append([]string{"compose", "--project-name", generateValidDockerComposeProjectName(projectName)}, args...),
			WorkingDir: options.WorkingDir,
			Env:        options.EnvVars,
			Logger:     options.Logger,
		}
	} else {
		cmd = shell.Command{
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
		return shell.RunCommandAndGetStdOut(t, cmd), nil
	}

	return shell.RunCommandAndGetOutputE(t, cmd)
}

// Note: docker-compose command doesn't like lower case or special characters, other than -.
func generateValidDockerComposeProjectName(str string) string {
	lower_str := strings.ToLower(str)
	return regexp.MustCompile(`[^a-zA-Z0-9 ]+`).ReplaceAllString(lower_str, "-")
}
