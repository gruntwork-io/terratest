package opa

import (
	"os"
	"path/filepath"
	"strings"
	"sync"

	getter "github.com/hashicorp/go-getter"

	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/testing"
)

var (
	// A map that maps the go-getter base URL to the temporary directory where it is downloaded.
	policyDirCache sync.Map
)

// DownloadPolicyE takes in a rule path written in go-getter syntax and downloads it to a temporary directory so that it
// can be passed to opa. The temporary directory that is used is cached based on the go-getter base path, and reused
// across calls.
// For example, if you call DownloadPolicyE with the go-getter URL multiple times:
//
//	git::https://github.com/gruntwork-io/terratest.git//policies/foo.rego?ref=master
//
// The first time the gruntwork-io/terratest repo will be downloaded to a new temp directory. All subsequent calls will
// reuse that first temporary dir where the repo was cloned. This is preserved even if a different subdir is requested
// later, e.g.: git::https://github.com/gruntwork-io/terratest.git//examples/bar.rego?ref=master.
// Note that the query parameters are always included in the base URL. This means that if you use a different ref (e.g.,
// git::https://github.com/gruntwork-io/terratest.git//examples/bar.rego?ref=v0.39.3), then that will be cloned to a new
// temporary directory rather than the cached dir.
func DownloadPolicyE(t testing.TestingT, rulePath string) (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	detected, err := getter.Detect(rulePath, cwd, getter.Detectors)
	if err != nil {
		return "", err
	}

	// File getters are assumed to be a local path reference, so pass through the original path.
	if strings.HasPrefix(detected, "file") {
		return rulePath, nil
	}

	// At this point we assume the getter URL is a remote URL, so we start the process of downloading it to a temp dir.

	// First, check if we had already downloaded the source and it is in our cache.
	baseDir, subDir := getter.SourceDirSubdir(rulePath)
	downloadPath, hasDownloaded := policyDirCache.Load(baseDir)
	if hasDownloaded {
		logger.Logf(t, "Previously downloaded %s: returning cached path", baseDir)
		return filepath.Join(downloadPath.(string), subDir), nil
	}

	// Not downloaded, so use go-getter to download the remote source to a temp dir.
	tempDir, err := os.MkdirTemp("", "terratest-opa-policy-*")
	if err != nil {
		return "", err
	}
	// go-getter doesn't work if you give it a directory that already exists, so we add an additional path in the
	// tempDir to make sure we feed a directory that doesn't exist yet.
	tempDir = filepath.Join(tempDir, "getter")

	logger.Logf(t, "Downloading %s to temp dir %s", rulePath, tempDir)
	if err := getter.GetAny(tempDir, baseDir); err != nil {
		return "", err
	}
	policyDirCache.Store(baseDir, tempDir)
	return filepath.Join(tempDir, subDir), nil
}
