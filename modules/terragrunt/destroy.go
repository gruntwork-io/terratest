package terragrunt

import (
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// DestroyAll runs terragrunt run-all destroy with the given options and returns stdout.
func DestroyAll(t testing.TestingT, options *Options) string {
	out, err := DestroyAllE(t, options)
	require.NoError(t, err)
	return out
}

// DestroyAllE runs terragrunt run-all destroy with the given options and returns stdout.
func DestroyAllE(t testing.TestingT, options *Options) (string, error) {
	return runTerragruntCommandE(t, options, "run-all", "destroy", "-auto-approve", "-input=false")
}
