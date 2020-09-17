package docker

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRun(t *testing.T) {
	t.Parallel()

	options := &RunOptions{
		Command:              []string{"-c", `echo "Hello, $NAME!"`},
		Entrypoint:           "sh",
		EnvironmentVariables: []string{"NAME=World"},
		Remove:               true,
	}

	out := Run(t, "alpine:3.7", options)
	require.Contains(t, out, "Hello, World!")
}
