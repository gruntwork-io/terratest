package docker_test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/docker/v2"
	"github.com/stretchr/testify/require"
)

func TestRun(t *testing.T) {
	t.Parallel()

	options := &docker.RunOptions{
		Command:              []string{"-c", `echo "Hello, $NAME!"`},
		Entrypoint:           "sh",
		EnvironmentVariables: []string{"NAME=World"},
		Remove:               true,
	}

	out := docker.RunContext(t, t.Context(), "alpine:3.7", options)
	require.Contains(t, out, "Hello, World!")
}
