package opa

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	getter "github.com/hashicorp/go-getter/v2"

	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/testing"
)

var (
	// A map that maps the go-getter base URL to the temporary directory where it is downloaded.
	policyDirCache sync.Map

	// A map of in-flight downloads keyed by baseDir. Used to deduplicate concurrent downloads of the same source so
	// that N parallel callers requesting the same rulePath result in a single actual download rather than N separate
	// downloads racing into N separate temp directories.
	inFlightDownloads sync.Map
)

// inFlightDownload represents a download operation that may currently be in progress. The first goroutine to request a
// given baseDir performs the download; all subsequent goroutines block on done and reuse its result.
type inFlightDownload struct {
	done   chan struct{}
	result string
	err    error
}

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

	downloadPath, hasDownloaded := policyDirCache.Load(baseDir)
	if hasDownloaded {
		logger.Default.Logf(t, "Previously downloaded %s: returning cached path", baseDir)
		return filepath.Join(downloadPath.(string), subDir), nil
	}

	// Cache miss. Coordinate with any other goroutines that may also be downloading this same baseDir so that we don't
	// end up with N parallel downloads racing into N separate temp directories. The first arrival performs the
	// download; everyone else waits on the in-flight entry and reuses the result.
	entry := &inFlightDownload{done: make(chan struct{})}

	actual, loaded := inFlightDownloads.LoadOrStore(baseDir, entry)
	if loaded {
		existing := actual.(*inFlightDownload)
		logger.Default.Logf(t, "Download of %s already in flight: waiting for it to complete", baseDir)
		<-existing.done
		if existing.err != nil {
			return "", existing.err
		}
		return filepath.Join(existing.result, subDir), nil
	}

	// We are the first arrival; perform the download and signal completion to any waiters.
	tempDir, err := downloadPolicyToTempDirFn(t, rulePath, baseDir)

	entry.result = tempDir
	entry.err = err
	close(entry.done)
	inFlightDownloads.Delete(baseDir)

	if err != nil {
		return "", err
	}

	policyDirCache.Store(baseDir, tempDir)

	return filepath.Join(tempDir, subDir), nil
}

// downloadPolicyToTempDirFn is the slow path of DownloadPolicyE: it actually downloads the given baseDir using
// go-getter into a fresh temp directory and returns the path to that directory. It is held as a package-level variable
// so that tests can swap it out with a stub to avoid hitting the network and to count invocations.
var downloadPolicyToTempDirFn = downloadPolicyToTempDir

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
