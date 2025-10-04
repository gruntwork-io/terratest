package terragrunt

import (
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// TgStackGenerate calls terragrunt stack generate and returns stdout/stderr.
func TgStackGenerate(t testing.TestingT, options *Options) string {
	out, err := TgStackGenerateE(t, options)
	require.NoError(t, err)
	return out
}

// TgStackGenerateE calls terragrunt stack generate and returns stdout/stderr.
func TgStackGenerateE(t testing.TestingT, options *Options) (string, error) {
	return runTerragruntStackCommandE(t, options, "generate")
}
