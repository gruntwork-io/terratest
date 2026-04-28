package opa

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	getter "github.com/hashicorp/go-getter/v2"
	"golang.org/x/sync/singleflight"

	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/testing"
)

var (
	// A map that maps the go-getter base URL to the temporary directory where it is downloaded.
	policyDirCache sync.Map

	// downloadGroup deduplicates concurrent downloads for the same baseDir so that N parallel callers requesting the
	// same rulePath result in a single underlying download rather than N separate downloads racing into N temp dirs.
	downloadGroup singleflight.Group
)

// DownloadPolicyE takes in a rule path written in go-getter syntax and downloads it to a temporary directory so that it
// can be passed to opa. The temporary directory that is used is cached based on the go-getter base path, and reused
// across calls.
// For example, if you call DownloadPolicyE with the go-getter URL multiple times:
//
//	git::https://github.com/gruntwork-io/terratest.git//policies/foo.rego?ref=main
//
// The first time the gruntwork-io/terratest repo will be downloaded to a new temp directory. All subsequent calls will
// reuse that first temporary dir where the repo was cloned. This is preserved even if a different subdir is requested
// later, e.g.: git::https://github.com/gruntwork-io/terratest.git//examples/bar.rego?ref=main
// Note that the query parameters are always included in the base URL. This means that if you use a different ref (e.g.,
// git::https://github.com/gruntwork-io/terratest.git//examples/bar.rego?ref=v0.39.3), then that will be cloned to a new
// temporary directory rather than the cached dir.
func DownloadPolicyE(t testing.TestingT, rulePath string) (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("getting current working directory: %w", err)
	}

	// File getters are assumed to be a local path reference, so pass through the original path.
	var fileGetter getter.FileGetter
	if ok, _ := fileGetter.Detect(&getter.Request{
		Src:     rulePath,
		Pwd:     cwd,
		GetMode: getter.ModeAny,
	}); ok {
		return rulePath, nil
	}

	// At this point we assume the getter URL is a remote URL, so we start the process of downloading it to a temp dir.

	// First, check if we had already downloaded the source and it is in our cache.
	baseDir, subDir := getter.SourceDirSubdir(rulePath)

	if downloadPath, hasDownloaded := policyDirCache.Load(baseDir); hasDownloaded {
		logger.Default.Logf(t, "Previously downloaded %s: returning cached path", baseDir)
		return filepath.Join(downloadPath.(string), subDir), nil
	}

	// Cache miss. Use singleflight to ensure that only one goroutine actually performs the download for a given
	// baseDir; any concurrent callers block on the same call and reuse its result.
	v, err, _ := downloadGroup.Do(baseDir, func() (any, error) {
		// Re-check the cache in case another goroutine populated it while we were waiting to enter the singleflight.
		if downloadPath, hasDownloaded := policyDirCache.Load(baseDir); hasDownloaded {
			return downloadPath.(string), nil
		}

		tempDir, err := downloadPolicyToTempDir(t, rulePath, baseDir)
		if err != nil {
			return "", err
		}

		policyDirCache.Store(baseDir, tempDir)

		return tempDir, nil
	})
	if err != nil {
		return "", err
	}

	return filepath.Join(v.(string), subDir), nil
}

// downloadPolicyToTempDir downloads the given baseDir using go-getter into a fresh temp directory and returns the path
// to the directory containing the downloaded source.
func downloadPolicyToTempDir(t testing.TestingT, rulePath, baseDir string) (string, error) {
	tempDir, err := os.MkdirTemp("", "terratest-opa-policy-*")
	if err != nil {
		return "", fmt.Errorf("creating temp directory for policy download: %w", err)
	}

	// go-getter doesn't work if you give it a directory that already exists, so we add an additional path in the
	// tempDir to make sure we feed a directory that doesn't exist yet.
	tempDir = filepath.Join(tempDir, "getter")

	logger.Default.Logf(t, "Downloading %s to temp dir %s", rulePath, tempDir)

	if _, err := getter.GetAny(context.Background(), tempDir, baseDir); err != nil {
		return "", fmt.Errorf("downloading policy from %s: %w", baseDir, err)
	}

	return tempDir, nil
}
