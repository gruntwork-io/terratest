package test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/docker/v2"
	"github.com/stretchr/testify/assert"
)

func TestDockerComposeStdoutExample(t *testing.T) {
	t.Parallel()
	dockerComposeFile := "../examples/docker-compose-stdout-example/docker-compose.yml"

	// Run the build step first so that the build output doesn't go to stdout during the compose step.
	docker.RunDockerCompose(
		t,
		&docker.Options{},
		"-f",
		dockerComposeFile,
		"build",
	)

	// Run the Docker image, read the stdout from it, and make sure it contains the expected output.
	// The script must be run using `run bash_script` rather than `up`, so that the echo output from the script
	// is the only thing that outputs to stdout.
	output := docker.RunDockerComposeAndGetStdOut(
		t,
		&docker.Options{},
		"-f",
		dockerComposeFile,
		"run",
		"bash_script",
	)

	assert.Contains(t, output, "stdout: message")
	assert.NotContains(t, output, "stderr: error")
}
