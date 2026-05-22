package git_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/gruntwork-io/terratest/modules/core/v2/git"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//nolint:paralleltest,tparallel // subtests mutate shared git state (checkout) and must run sequentially
func TestGitRefChecks(t *testing.T) {
	t.Parallel()

	tmpdir := t.TempDir()
	gitWorkDir := filepath.Join(tmpdir, "terratest")

	url := "https://github.com/gruntwork-io/terratest.git"
	err := exec.CommandContext(t.Context(), "git", "clone", url, gitWorkDir).Run()
	require.NoError(t, err)

	t.Run("GetCurrentBranchNameReturnsBranchName", func(t *testing.T) {
		err := exec.CommandContext(t.Context(), "git", "-C", gitWorkDir, "checkout", "main").Run()
		require.NoError(t, err)

		name := git.GetCurrentBranchNameContext(t, t.Context(), gitWorkDir)

		assert.Equal(t, "main", name)
	})

	t.Run("GetCurrentBranchNameReturnsEmptyForDetachedState", func(t *testing.T) {
		err := exec.CommandContext(t.Context(), "git", "-C", gitWorkDir, "checkout", "v0.0.1").Run()
		require.NoError(t, err)

		name := git.GetCurrentBranchNameContext(t, t.Context(), gitWorkDir)

		assert.Empty(t, name)
	})

	t.Run("GetCurrentRefReturnsBranchName", func(t *testing.T) {
		err := exec.CommandContext(t.Context(), "git", "-C", gitWorkDir, "checkout", "main").Run()
		require.NoError(t, err)

		name := git.GetCurrentGitRefContext(t, t.Context(), gitWorkDir)

		assert.Equal(t, "main", name)
	})

	t.Run("GetCurrentRefReturnsTagValue", func(t *testing.T) {
		err := exec.CommandContext(t.Context(), "git", "-C", gitWorkDir, "checkout", "v0.0.1").Run()
		require.NoError(t, err)

		name := git.GetCurrentGitRefContext(t, t.Context(), gitWorkDir)

		assert.Equal(t, "v0.0.1", name)
	})

	t.Run("GetCurrentRefReturnsLightTagValue", func(t *testing.T) {
		err := exec.CommandContext(t.Context(), "git", "-C", gitWorkDir, "checkout", "58d3ea8").Run()
		require.NoError(t, err)

		name := git.GetCurrentGitRefContext(t, t.Context(), gitWorkDir)

		assert.Equal(t, "v0.0.1-1-g58d3ea8f", name)
	})
}

func TestGetRepoRoot(t *testing.T) {
	t.Parallel()

	cwd, err := os.Getwd()
	require.NoError(t, err)

	expectedRepoRoot, err := filepath.Abs(filepath.Join(cwd, "..", "..", ".."))
	require.NoError(t, err)

	repoRoot := git.GetRepoRootContext(t, t.Context(), cwd)
	assert.Equal(t, expectedRepoRoot, repoRoot)
}
