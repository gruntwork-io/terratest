package opa

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/gruntwork-io/terratest/modules/files"
	"github.com/gruntwork-io/terratest/modules/git"
)

// Test to make sure the DownloadPolicyE function returns a local path without processing it.
func TestDownloadPolicyReturnsLocalPath(t *testing.T) {
	t.Parallel()

	localPath := "../../examples/terraform-opa-example/policy/enforce_source.rego"
	path, err := DownloadPolicyE(t, localPath)
	require.NoError(t, err)
	assert.Equal(t, localPath, path)
}

// Test to make sure the DownloadPolicyE function returns a remote path to a temporary directory.
func TestDownloadPolicyDownloadsRemote(t *testing.T) {
	t.Parallel()

	curRef := git.GetCurrentGitRef(t)
	baseDir := fmt.Sprintf("git::https://github.com/gruntwork-io/terratest.git?ref=%s", curRef)
	localPath := "../../examples/terraform-opa-example/policy/enforce_source.rego"
	remotePath := fmt.Sprintf("git::https://github.com/gruntwork-io/terratest.git//examples/terraform-opa-example/policy/enforce_source.rego?ref=%s", curRef)

	// Make sure we clean up the downloaded file, while simultaneously asserting that the download dir was stored in the
	// cache.
	defer func() {
		downloadPathRaw, inCache := policyDirCache.Load(baseDir)
		require.True(t, inCache)
		downloadPath := downloadPathRaw.(string)
		if strings.HasSuffix(downloadPath, "/getter") {
			downloadPath = filepath.Dir(downloadPath)
		}
		assert.NoError(t, os.RemoveAll(downloadPath))
	}()

	path, err := DownloadPolicyE(t, remotePath)
	require.NoError(t, err)

	absPath, err := filepath.Abs(localPath)
	require.NoError(t, err)
	assert.NotEqual(t, absPath, path)

	localContents, err := os.ReadFile(localPath)
	require.NoError(t, err)
	remoteContents, err := os.ReadFile(path)
	require.NoError(t, err)
	assert.Equal(t, localContents, remoteContents)
}

// Test to make sure the DownloadPolicyE function uses the cache if it has already downloaded an existing base path.
func TestDownloadPolicyReusesCachedDir(t *testing.T) {
	t.Parallel()

	baseDir := "git::https://github.com/gruntwork-io/terratest.git?ref=master"
	remotePath := "git::https://github.com/gruntwork-io/terratest.git//examples/terraform-opa-example/policy/enforce_source.rego?ref=master"
	remotePathAltSubPath := "git::https://github.com/gruntwork-io/terratest.git//modules/opa/eval.go?ref=master"

	// Make sure we clean up the downloaded file, while simultaneously asserting that the download dir was stored in the
	// cache.
	defer func() {
		downloadPathRaw, inCache := policyDirCache.Load(baseDir)
		require.True(t, inCache)
		downloadPath := downloadPathRaw.(string)

		if strings.HasSuffix(downloadPath, "/getter") {
			downloadPath = filepath.Dir(downloadPath)
		}
		assert.NoError(t, os.RemoveAll(downloadPath))
	}()

	path, err := DownloadPolicyE(t, remotePath)
	require.NoError(t, err)
	files.FileExists(path)

	downloadPathRaw, inCache := policyDirCache.Load(baseDir)
	require.True(t, inCache)
	downloadPath := downloadPathRaw.(string)

	// make sure the second call is exactly equal to the first call
	newPath, err := DownloadPolicyE(t, remotePath)
	require.NoError(t, err)
	assert.Equal(t, path, newPath)

	// Also make sure the cache is reused for alternative sub dirs.
	newAltPath, err := DownloadPolicyE(t, remotePathAltSubPath)
	require.NoError(t, err)
	assert.True(t, strings.HasPrefix(path, downloadPath))
	assert.True(t, strings.HasPrefix(newAltPath, downloadPath))
}
