package terraform

import (
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// Validate calls terraform validate and returns stdout/stderr.
func Validate(t testing.TestingT, options *Options) string {
	out, err := ValidateE(t, options)
	require.NoError(t, err)
	return out
}

// ValidateE calls terraform validate and returns stdout/stderr.
func ValidateE(t testing.TestingT, options *Options) (string, error) {
	return RunTerraformCommandE(t, options, FormatArgs(options, "validate")...)
}

// InitAndValidate runs terraform init and validate with the given options and returns stdout/stderr from the validate command.
// This will fail the test if there is an error in the command.
func InitAndValidate(t testing.TestingT, options *Options) string {
	out, err := InitAndValidateE(t, options)
	require.NoError(t, err)
	return out
}

// InitAndValidateE runs terraform init and validate with the given options and returns stdout/stderr from the validate command.
func InitAndValidateE(t testing.TestingT, options *Options) (string, error) {
	if _, err := InitE(t, options); err != nil {
		return "", err
	}

	return ValidateE(t, options)
}
