package docker

import (
	"fmt"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/git"
	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/random"
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

func TestBuildWithBuildKit(t *testing.T) {
	t.Parallel()

	tag := "gruntwork-io/test-image-with-buildkit:v1"
	testToken := "testToken"
	options := &BuildOptions{
		Tags:           []string{tag},
		EnableBuildKit: true,
		OtherOptions:   []string{"--secret", fmt.Sprintf("id=github-token,env=%s", "GITHUB_OAUTH_TOKEN")},
		Env:            map[string]string{"GITHUB_OAUTH_TOKEN": testToken},
	}

	Build(t, "../../test/fixtures/docker-with-buildkit", options)
	out := Run(t, tag, &RunOptions{Remove: false})
	require.Contains(t, out, testToken)
}

func TestBuildMultiArch(t *testing.T) {
	t.Parallel()

	tag := "gruntwork-io/test-image:v1"
	text := "Hello, World!"

	options := &BuildOptions{
		Tags:          []string{tag},
		BuildArgs:     []string{fmt.Sprintf("text=%s", text)},
		Architectures: []string{"linux/arm64", "linux/amd64"},
		Load:          true,
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

func TestGitCloneAndBuild(t *testing.T) {
	t.Parallel()

	uniqueID := strings.ToLower(random.UniqueId())
	imageTag := "gruntwork-io-foo-test:" + uniqueID
	text := "Hello, World!"

	buildOpts := &BuildOptions{
		Tags:      []string{imageTag},
		BuildArgs: []string{fmt.Sprintf("text=%s", text)},
	}
	gitBranchName := git.GetCurrentBranchName(t)
	if gitBranchName == "" {
		logger.Logf(t, "WARNING: git.GetCurrentBranchName returned an empty string; falling back to master")
		gitBranchName = "master"
	}
	GitCloneAndBuild(t, "git@github.com:gruntwork-io/terratest.git", gitBranchName, "test/fixtures/docker", buildOpts)

	out := Run(t, imageTag, &RunOptions{Remove: true})
	require.Contains(t, out, text)
}
