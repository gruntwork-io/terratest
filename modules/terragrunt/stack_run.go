package terragrunt

import (
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// TgStackRun calls terragrunt stack run and returns stdout/stderr.
func TgStackRun(t testing.TestingT, options *Options) string {
	out, err := TgStackRunE(t, options)
	require.NoError(t, err)
	return out
}

// TgStackRunE calls terragrunt stack run and returns stdout/stderr.
func TgStackRunE(t testing.TestingT, options *Options) (string, error) {
	return runTerragruntStackCommandE(t, options, "run")
}
