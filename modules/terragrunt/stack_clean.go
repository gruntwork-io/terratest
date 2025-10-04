package terragrunt

import (
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// TgStackClean calls terragrunt stack clean to remove the .terragrunt-stack directory.
// This command cleans up the generated stack files created by stack generate or stack run.
func TgStackClean(t testing.TestingT, options *Options) string {
	out, err := TgStackCleanE(t, options)
	require.NoError(t, err)
	return out
}

// TgStackCleanE calls terragrunt stack clean to remove the .terragrunt-stack directory.
// This command cleans up the generated stack files created by stack generate or stack run.
func TgStackCleanE(t testing.TestingT, options *Options) (string, error) {
	return runTerragruntStackCommandE(t, options, "clean")
}