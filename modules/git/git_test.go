package git

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testGetCurrentBranchNameReturnsBranchName(t *testing.T) {
	err := exec.Command("git", "checkout", "master").Run()
	require.NoError(t, err)

	name := GetCurrentBranchName(t)

	assert.Equal(t, "master", name)
}

func testGetCurrentBranchNameReturnsEmptyForDetachedState(t *testing.T) {
	err := exec.Command("git", "checkout", "v0.0.1").Run()
	assert.Nil(t, err)

	name := GetCurrentBranchName(t)

	assert.Empty(t, name)
}

func testGetCurrentRefReturnsBranchName(t *testing.T) {
	err := exec.Command("git", "checkout", "master").Run()
	require.NoError(t, err)

	name := GetCurrentGitRef(t)

	assert.Equal(t, "master", name)
}

func testGetCurrentRefReturnsTagValue(t *testing.T) {
	err := exec.Command("git", "checkout", "v0.0.1").Run()
	require.NoError(t, err)

	name := GetCurrentGitRef(t)

	assert.Equal(t, "v0.0.1", name)
}

func testGetCurrentRefReturnsLightTagValue(t *testing.T) {
	err := exec.Command("git", "checkout", "58d3ea8").Run()
	require.NoError(t, err)

	name := GetCurrentGitRef(t)

	assert.Equal(t, "v0.0.1-1-g58d3ea8", name)
}

func TestGitRefChecks(t *testing.T) {
	t.Parallel()

	tmpdir, err := ioutil.TempDir("", "")
	require.NoError(t, err)
	defer os.RemoveAll(tmpdir)
	gitWorkDir := filepath.Join(tmpdir, "terratest")

	url := "https://github.com/gruntwork-io/terratest.git"
	err = exec.Command("git", "clone", url, gitWorkDir).Run()
	require.NoError(t, err)

	err = os.Chdir(gitWorkDir)
	require.NoError(t, err)

	t.Run("GetCurrentBranchNameReturnsBranchName", testGetCurrentBranchNameReturnsBranchName)
	t.Run("GetCurrentBranchNameReturnsEmptyForDetachedState", testGetCurrentBranchNameReturnsEmptyForDetachedState)
	t.Run("GetCurrentRefReturnsBranchName", testGetCurrentRefReturnsBranchName)
	t.Run("GetCurrentRefReturnsTagValue", testGetCurrentRefReturnsTagValue)
	t.Run("GetCurrentRefReturnsLightTagValue", testGetCurrentRefReturnsLightTagValue)
}
