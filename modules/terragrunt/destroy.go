package terragrunt

import (
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// DestroyAll runs terragrunt destroy with the given options and return stdout.
func DestroyAll(t testing.TestingT, options *Options) string {
	out, err := DestroyAllE(t, options)
	require.NoError(t, err)
	return out
}

// DestroyAllE runs terragrunt destroy with the given options and return stdout.
func DestroyAllE(t testing.TestingT, options *Options) (string, error) {
	return runTerragruntCommandE(t, options, "run-all", formatArgs(options, prepend(options.ExtraArgs.Destroy, "destroy", "-auto-approve", "-input=false")...)...)
}
