package helm

import (
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// Delete will delete the provided release from Tiller. If you set purge to true, Tiller will delete the release object
// as well so that the release name can be reused. This will fail the test if there is an error.
func Delete(t testing.TestingT, options *Options, releaseName string, purge bool) {
	require.NoError(t, DeleteE(t, options, releaseName, purge))
}

// DeleteE will delete the provided release from Tiller. If you set purge to true, Tiller will delete the release object
// as well so that the release name can be reused.
func DeleteE(t testing.TestingT, options *Options, releaseName string, purge bool) error {
	args := []string{}
	if !purge {
		args = append(args, "--keep-history")
	}
	if options.ExtraArgs != nil {
		if deleteArgs, ok := options.ExtraArgs["delete"]; ok {
			args = append(args, deleteArgs...)
		}
	}
	args = append(args, releaseName)
	_, err := RunHelmCommandAndGetOutputE(t, options, "delete", args...)
	return err
}
