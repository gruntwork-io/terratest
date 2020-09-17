package docker

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBuild(t *testing.T) {
	t.Parallel()

	tag := "gruntwork-io/test-image:v1"
	text := "Hello, World!"

	options := &BuildOptions{
		Tags:      []string{tag},
		BuildArgs: []string{fmt.Sprintf("text=%s", text)},
	}

	Build(t, "../../test/fixtures/docker", options)

	out := Run(t, tag, &RunOptions{Remove: true})
	require.Contains(t, out, text)
}
