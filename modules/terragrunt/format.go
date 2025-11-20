package terragrunt

import (
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// FormatAll runs terragrunt hcl format to format all terragrunt.hcl files and returns stdout/stderr
func FormatAll(t testing.TestingT, options *Options) string {
	out, err := FormatAllE(t, options)
	require.NoError(t, err)
	return out
}

// FormatAllE runs terragrunt hcl format to format all terragrunt.hcl files and returns stdout/stderr
func FormatAllE(t testing.TestingT, options *Options) (string, error) {
	return runTerragruntCommandE(t, options, "hcl", "format")
}
