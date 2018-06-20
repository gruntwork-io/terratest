package docker

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/shell"
)

// Options are Docker options.
type Options struct {
	WorkingDir string
	EnvVars    map[string]string
}

// RunDockerCompose runs docker-compose with the given arguments and options and return stdout/stderr.
func RunDockerCompose(t *testing.T, options *Options, args ...string) string {
	out, err := RunDockerComposeE(t, options, args...)
	if err != nil {
		t.Fatal(err)
	}
	return out
}

// RunDockerComposeE runs docker-compose with the given arguments and options and return stdout/stderr.
func RunDockerComposeE(t *testing.T, options *Options, args ...string) (string, error) {
	cmd := shell.Command{
		Command: "docker-compose",
		// We append --project-name to ensure containers from multiple different tests using Docker Compose don't end
		// up in the same project and end up conflicting with each other.
		Args:       append([]string{"--project-name", t.Name()}, args...),
		WorkingDir: options.WorkingDir,
		Env:        options.EnvVars,
	}

	return shell.RunCommandAndGetOutputE(t, cmd)
}
