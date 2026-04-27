package opa_test

import (
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/gruntwork-io/terratest/modules/files"
	"github.com/gruntwork-io/terratest/modules/git"
	"github.com/gruntwork-io/terratest/modules/opa"
)

// TestDownloadPolicyReturnsLocalPath makes sure the DownloadPolicyE function returns a local path without processing it.
func TestDownloadPolicyReturnsLocalPath(t *testing.T) {
	t.Parallel()

	localPath := "../../examples/terraform-opa-example/policy/enforce_source.rego"
	path, err := opa.DownloadPolicyE(t, localPath)
	require.NoError(t, err)
	assert.Equal(t, localPath, path)
}

// TestDownloadPolicyDownloadsRemote makes sure the DownloadPolicyE function returns a remote path to a temporary
// directory.
func TestDownloadPolicyDownloadsRemote(t *testing.T) {
	t.Parallel()

	curRef := git.GetCurrentGitRefContext(t, t.Context(), "")
	baseDir := "git::https://github.com/gruntwork-io/terratest.git?ref=" + curRef
	localPath := "../../examples/terraform-opa-example/policy/enforce_source.rego"
	remotePath := "git::https://github.com/gruntwork-io/terratest.git//examples/terraform-opa-example/policy/enforce_source.rego?ref=" + curRef

	// Make sure we clean up the downloaded file, while simultaneously asserting that the download dir was stored in the
	// cache.
	defer func() {
		downloadPathRaw, inCache := opa.PolicyDirCache.Load(baseDir)
		require.True(t, inCache)

		downloadPath := downloadPathRaw.(string)

		if strings.HasSuffix(downloadPath, "/getter") {
			downloadPath = filepath.Dir(downloadPath)
		}

		assert.NoError(t, os.RemoveAll(downloadPath))
	}()

	path, err := opa.DownloadPolicyE(t, remotePath)
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

// TestDownloadPolicyDeduplicatesConcurrentDownloads makes sure that when many goroutines simultaneously request the
// same rulePath, only a single underlying download is performed. Without the in-flight deduplication this would
// trigger one download per goroutine, each into its own temp directory.
func TestDownloadPolicyDeduplicatesConcurrentDownloads(t *testing.T) {
	// Not Parallel: this test mutates the package-level downloader and caches.
	opa.ResetCachesForTest()
	defer opa.ResetCachesForTest()

	// Use a temp directory as a stand-in for a downloaded source so we don't need network access.
	fakeDownloadDir := t.TempDir()

	// Block the first arrival inside the slow path until all goroutines have had a chance to enter DownloadPolicyE.
	// This maximizes the chance that the in-flight dedup path is exercised on the other 49 goroutines.
	releaseDownload := make(chan struct{})

	callCount, restore := opa.SetDownloadPolicyToTempDirForTest(func() (string, error) {
		<-releaseDownload
		return fakeDownloadDir, nil
	})
	defer restore()

	const numGoroutines = 50
	rulePath := "git::https://example.invalid/repo.git//policy.rego?ref=test"

	var wg sync.WaitGroup
	results := make([]string, numGoroutines)
	errs := make([]error, numGoroutines)
	started := make(chan struct{}, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			started <- struct{}{}
			path, err := opa.DownloadPolicyE(t, rulePath)
			results[idx] = path
			errs[idx] = err
		}(i)
	}

	// Wait until every goroutine has at least entered DownloadPolicyE before releasing the slow path.
	for i := 0; i < numGoroutines; i++ {
		<-started
	}
	// Give the goroutines a moment to actually progress to the LoadOrStore call before unblocking the first arrival.
	time.Sleep(50 * time.Millisecond)
	close(releaseDownload)

	wg.Wait()

	expectedPath := filepath.Join(fakeDownloadDir, "policy.rego")
	for i := 0; i < numGoroutines; i++ {
		require.NoError(t, errs[i], "goroutine %d failed", i)
		assert.Equal(t, expectedPath, results[i], "goroutine %d got unexpected path", i)
	}

	assert.Equal(t, int64(1), atomic.LoadInt64(callCount),
		"expected exactly one underlying download, got %d", atomic.LoadInt64(callCount))
}

// TestDownloadPolicyReusesCachedDir makes sure the DownloadPolicyE function uses the cache if it has already downloaded
// an existing base path.
func TestDownloadPolicyReusesCachedDir(t *testing.T) {
	t.Parallel()

	baseDir := "git::https://github.com/gruntwork-io/terratest.git?ref=main"
	remotePath := "git::https://github.com/gruntwork-io/terratest.git//examples/terraform-opa-example/policy/enforce_source.rego?ref=main"
	remotePathAltSubPath := "git::https://github.com/gruntwork-io/terratest.git//modules/opa/eval.go?ref=main"

	// Make sure we clean up the downloaded file, while simultaneously asserting that the download dir was stored in the
	// cache.
	defer func() {
		downloadPathRaw, inCache := opa.PolicyDirCache.Load(baseDir)
		require.True(t, inCache)

		downloadPath := downloadPathRaw.(string)

		if strings.HasSuffix(downloadPath, "/getter") {
			downloadPath = filepath.Dir(downloadPath)
		}

		assert.NoError(t, os.RemoveAll(downloadPath))
	}()

	path, err := opa.DownloadPolicyE(t, remotePath)
	require.NoError(t, err)
	files.FileExists(path)

	downloadPathRaw, inCache := opa.PolicyDirCache.Load(baseDir)
	require.True(t, inCache)

	downloadPath := downloadPathRaw.(string)

	// make sure the second call is exactly equal to the first call
	newPath, err := opa.DownloadPolicyE(t, remotePath)
	require.NoError(t, err)
	assert.Equal(t, path, newPath)

	// Also make sure the cache is reused for alternative sub dirs.
	newAltPath, err := opa.DownloadPolicyE(t, remotePathAltSubPath)
	require.NoError(t, err)
	assert.True(t, strings.HasPrefix(path, downloadPath))
	assert.True(t, strings.HasPrefix(newAltPath, downloadPath))
}
