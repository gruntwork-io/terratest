package terraform

import (
	"github.com/gruntwork-io/terratest/modules/testing/v2"
	"github.com/stretchr/testify/require"
)

// Validate calls terraform validate and returns stdout/stderr.
func Validate(t testing.TestingT, options *Options) string {
	out, err := ValidateE(t, options)
	require.NoError(t, err)
	return out
}

// ValidateInputs calls terragrunt hcl validate and returns stdout/stderr.
// This validates Terragrunt HCL syntax and checks for mis-aligned inputs.
//
// Deprecated: Use github.com/gruntwork-io/terratest/modules/terragrunt package instead.
func ValidateInputs(t testing.TestingT, options *Options) string {
	out, err := ValidateInputsE(t, options)
	require.NoError(t, err)
	return out
}

// ValidateE calls terraform validate and returns stdout/stderr.
func ValidateE(t testing.TestingT, options *Options) (string, error) {
	return RunTerraformCommandE(t, options, FormatArgs(options, prepend(options.ExtraArgs.Validate, "validate")...)...)
}

// ValidateInputsE calls terragrunt hcl validate and returns stdout/stderr.
// This validates Terragrunt HCL syntax and checks for mis-aligned inputs.
//
// Deprecated: Use github.com/gruntwork-io/terratest/modules/terragrunt package instead.
func ValidateInputsE(t testing.TestingT, options *Options) (string, error) {
	if options.TerraformBinary != "terragrunt" {
		return "", TgInvalidBinary(options.TerraformBinary)
	}
	return RunTerraformCommandE(t, options, FormatArgs(options, prepend(options.ExtraArgs.ValidateInputs, "hcl", "validate")...)...)
}

// InitAndValidate runs terraform init and validate with the given options and returns stdout/stderr from the validate command.
// This will fail the test if there is an error in the command.
func InitAndValidate(t testing.TestingT, options *Options) string {
	out, err := InitAndValidateE(t, options)
	require.NoError(t, err)
	return out
}

// InitAndValidateInputs runs terragrunt init and hcl validate with the given options and returns stdout/stderr from the validate command.
// This validates Terragrunt HCL syntax and checks for mis-aligned inputs.
//
// Deprecated: Use github.com/gruntwork-io/terratest/modules/terragrunt package instead.
func InitAndValidateInputs(t testing.TestingT, options *Options) string {
	out, err := InitAndValidateInputsE(t, options)
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

// InitAndValidateInputsE runs terragrunt init and hcl validate with the given options and returns stdout/stderr.
// This validates Terragrunt HCL syntax and checks for mis-aligned inputs.
//
// Deprecated: Use github.com/gruntwork-io/terratest/modules/terragrunt package instead.
func InitAndValidateInputsE(t testing.TestingT, options *Options) (string, error) {
	if _, err := InitE(t, options); err != nil {
		return "", err
	}
	return ValidateInputsE(t, options)
}
