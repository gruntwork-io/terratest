package docker_test

import (
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/docker/v2"
	"github.com/gruntwork-io/terratest/modules/core/v2/git"
	"github.com/gruntwork-io/terratest/modules/core/v2/logger"
	"github.com/gruntwork-io/terratest/modules/core/v2/random"
	"github.com/stretchr/testify/require"
)

func TestBuild(t *testing.T) {
	t.Parallel()

	tag := "gruntwork-io/test-image:v1"
	text := "Hello, World!"

	options := &docker.BuildOptions{
		Tags:      []string{tag},
		BuildArgs: []string{"text=" + text},
	}

	ctx := t.Context()
	docker.BuildContext(t, ctx, "../../test/fixtures/docker", options)

	out := docker.RunContext(t, ctx, tag, &docker.RunOptions{Remove: true})
	require.Contains(t, out, text)
}

func TestBuildWithBuildKit(t *testing.T) {
	t.Parallel()

	tag := "gruntwork-io/test-image-with-buildkit:v1"
	testToken := "testToken"
	options := &docker.BuildOptions{
		Tags:           []string{tag},
		EnableBuildKit: true,
		OtherOptions:   []string{"--secret", "id=github-token,env=" + "GITHUB_OAUTH_TOKEN"},
		Env:            map[string]string{"GITHUB_OAUTH_TOKEN": testToken},
	}

	ctx := t.Context()
	docker.BuildContext(t, ctx, "../../test/fixtures/docker-with-buildkit", options)

	out := docker.RunContext(t, ctx, tag, &docker.RunOptions{Remove: false})
	require.Contains(t, out, testToken)
}

func TestBuildMultiArch(t *testing.T) {
	t.Parallel()

	tag := "gruntwork-io/test-image:v1"
	text := "Hello, World!"

	options := &docker.BuildOptions{
		Tags:          []string{tag},
		BuildArgs:     []string{"text=" + text},
		Architectures: []string{"linux/arm64", "linux/amd64"},
		Load:          true,
	}

	ctx := t.Context()
	docker.BuildContext(t, ctx, "../../test/fixtures/docker", options)

	out := docker.RunContext(t, ctx, tag, &docker.RunOptions{Remove: true})
	require.Contains(t, out, text)
}

func TestBuildWithTarget(t *testing.T) {
	t.Parallel()

	tag := "gruntwork-io/test-image:target1"
	text := "Hello, World!"
	text1 := "Hello, World! This is build target 1!"

	options := &docker.BuildOptions{
		Tags:      []string{tag},
		BuildArgs: []string{"text=" + text, "text1=" + text1},
		Target:    "step1",
	}

	ctx := t.Context()
	docker.BuildContext(t, ctx, "../../test/fixtures/docker", options)

	out := docker.RunContext(t, ctx, tag, &docker.RunOptions{Remove: true})
	require.Contains(t, out, text1)
}

func TestGitCloneAndBuild(t *testing.T) {
	t.Parallel()

	uniqueID := strings.ToLower(random.UniqueID())
	imageTag := "gruntwork-io-foo-test:" + uniqueID
	text := "Hello, World!"

	buildOpts := &docker.BuildOptions{
		Tags:      []string{imageTag},
		BuildArgs: []string{"text=" + text},
	}

	ctx := t.Context()

	gitBranchName := git.GetCurrentBranchNameContext(t, ctx, "")
	if gitBranchName == "" {
		logger.Default.Logf(t, "WARNING: git.GetCurrentBranchNameContext returned an empty string; falling back to main")

		gitBranchName = "main"
	}

	docker.GitCloneAndBuildContext(t, ctx, "git@github.com:gruntwork-io/terratest.git", gitBranchName, "test/fixtures/docker", buildOpts)

	out := docker.RunContext(t, ctx, imageTag, &docker.RunOptions{Remove: true})
	require.Contains(t, out, text)
}
