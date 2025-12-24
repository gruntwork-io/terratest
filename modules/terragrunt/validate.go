package terragrunt

import (
	"github.com/gruntwork-io/terratest/modules/testing/v2"
	"github.com/stretchr/testify/require"
)

// ValidateAll runs terragrunt validate --all with the given options and returns stdout/stderr
func ValidateAll(t testing.TestingT, options *Options) string {
	out, err := ValidateAllE(t, options)
	require.NoError(t, err)
	return out
}

// ValidateAllE runs terragrunt validate --all with the given options and returns stdout/stderr
func ValidateAllE(t testing.TestingT, options *Options) (string, error) {
	return runTerragruntCommandE(t, options, "validate", "--all")
}
