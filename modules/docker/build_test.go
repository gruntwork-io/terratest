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

func TestBuildWithTarget(t *testing.T) {
	t.Parallel()

	tag := "gruntwork-io/test-image:target1"
	text := "Hello, World!"
	text1 := "Hello, World! This is build target 1!"

	options := &BuildOptions{
		Tags:      []string{tag},
		BuildArgs: []string{fmt.Sprintf("text=%s", text), fmt.Sprintf("text1=%s", text1)},
		Target:    "step1",
	}

	Build(t, "../../test/fixtures/docker", options)

	out := Run(t, tag, &RunOptions{Remove: true})
	require.Contains(t, out, text1)
}
