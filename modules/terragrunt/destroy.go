package terragrunt

import (
	"github.com/gruntwork-io/terratest/modules/testing/v2"
	"github.com/stretchr/testify/require"
)

// DestroyAll runs terragrunt destroy --all with the given options and returns stdout.
func DestroyAll(t testing.TestingT, options *Options) string {
	out, err := DestroyAllE(t, options)
	require.NoError(t, err)
	return out
}

// DestroyAllE runs terragrunt destroy --all with the given options and returns stdout.
func DestroyAllE(t testing.TestingT, options *Options) (string, error) {
	return runTerragruntCommandE(t, options, "destroy", "--all", "-auto-approve", "-input=false")
}
