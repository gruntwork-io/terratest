package opa

import (
	"sync/atomic"

	"github.com/gruntwork-io/terratest/modules/testing"
)

// FormatOPAEvalArgs is an exported alias for formatOPAEvalArgs, used by external test packages.
var FormatOPAEvalArgs = formatOPAEvalArgs

// PolicyDirCache is an exported alias for policyDirCache, used by external test packages.
var PolicyDirCache = &policyDirCache

// SetDownloadPolicyToTempDirForTest swaps the package-level downloader for the duration of a test. The returned
// counter is incremented every time the swapped-in downloader is invoked, so tests can assert how many times the slow
// path actually ran. The returned restore function reinstalls the original downloader.
func SetDownloadPolicyToTempDirForTest(fn func() (string, error)) (callCount *int64, restore func()) {
	original := downloadPolicyToTempDirFn
	var counter int64
	downloadPolicyToTempDirFn = func(_ testing.TestingT, _, _ string) (string, error) {
		atomic.AddInt64(&counter, 1)
		return fn()
	}
	return &counter, func() {
		downloadPolicyToTempDirFn = original
	}
}

// ResetCachesForTest clears the package-level caches so tests can run from a clean slate.
func ResetCachesForTest() {
	policyDirCache.Range(func(k, _ interface{}) bool {
		policyDirCache.Delete(k)
		return true
	})
	inFlightDownloads.Range(func(k, _ interface{}) bool {
		inFlightDownloads.Delete(k)
		return true
	})
}
