package helm

import (
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// Rollback will downgrade the release to the specified version. This will fail
// the test if there is an error.
func Rollback(t testing.TestingT, options *Options, releaseName string, revision string) {
	require.NoError(t, RollbackE(t, options, releaseName, revision))
}

// RollbackE will downgrade the release to the specified version
func RollbackE(t testing.TestingT, options *Options, releaseName string, revision string) error {
	var err error
	args := []string{}
	args = append(args, getNamespaceArgs(options)...)
	args = append(args, releaseName)
	if revision != "" {
		args = append(args, revision)
	}
	_, err = RunHelmCommandAndGetOutputE(t, options, "rollback", args...)
	return err
}
