package terragrunt

import (
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// StackClean calls terragrunt stack clean to remove the .terragrunt-stack directory
// This command cleans up the generated stack files created by stack generate or stack run
func StackClean(t testing.TestingT, options *Options) string {
	out, err := StackCleanE(t, options)
	require.NoError(t, err)
	return out
}

// StackCleanE calls terragrunt stack clean to remove the .terragrunt-stack directory
// This command cleans up the generated stack files created by stack generate or stack run
func StackCleanE(t testing.TestingT, options *Options) (string, error) {
	return runTerragruntStackCommandE(t, options, "clean")
}
