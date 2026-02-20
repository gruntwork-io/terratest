package terragrunt

import (
	"github.com/gruntwork-io/terratest/modules/testing"
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

// Validate runs terragrunt validate for a single unit and returns stdout/stderr.
func Validate(t testing.TestingT, options *Options) string {
	out, err := ValidateE(t, options)
	require.NoError(t, err)
	return out
}

// ValidateE runs terragrunt validate for a single unit and returns stdout/stderr.
func ValidateE(t testing.TestingT, options *Options) (string, error) {
	return runTerragruntCommandE(t, options, "validate")
}

// InitAndValidate runs terragrunt init followed by validate for a single unit and returns the validate stdout/stderr.
func InitAndValidate(t testing.TestingT, options *Options) string {
	out, err := InitAndValidateE(t, options)
	require.NoError(t, err)
	return out
}

// InitAndValidateE runs terragrunt init followed by validate for a single unit and returns the validate stdout/stderr.
func InitAndValidateE(t testing.TestingT, options *Options) (string, error) {
	if _, err := InitE(t, options); err != nil {
		return "", err
	}
	return ValidateE(t, options)
}
