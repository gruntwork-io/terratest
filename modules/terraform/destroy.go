package terraform

import (
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// Destroy runs terraform destroy with the given options and return stdout/stderr.
func Destroy(t testing.TestingT, options *Options) string {
	out, err := DestroyE(t, options)
	require.NoError(t, err)
	return out
}

// DestroyE runs terraform destroy with the given options and return stdout/stderr.
func DestroyE(t testing.TestingT, options *Options) (string, error) {
	return RunTerraformCommandE(t, options, FormatArgs(options, prepend(options.ExtraArgs.Destroy, "destroy", "-auto-approve", "-input=false")...)...)
}
