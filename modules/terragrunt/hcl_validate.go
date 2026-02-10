package terragrunt

import (
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// HclValidate runs terragrunt hcl validate to check terragrunt.hcl syntax.
// This validates Terragrunt HCL configuration and can check for mis-aligned inputs.
// Use TerraformArgs to pass additional flags like "--inputs" or "--strict".
//
// Examples:
//
//	HclValidate(t, options)                                        // Basic syntax check
//	HclValidate(t, &Options{TerraformArgs: []string{"--inputs"}})  // Check input alignment
func HclValidate(t testing.TestingT, options *Options) string {
	out, err := HclValidateE(t, options)
	require.NoError(t, err)
	return out
}

// HclValidateE runs terragrunt hcl validate to check terragrunt.hcl syntax.
// This validates Terragrunt HCL configuration and can check for mis-aligned inputs.
// Use TerraformArgs to pass additional flags like "--inputs" or "--strict".
func HclValidateE(t testing.TestingT, options *Options) (string, error) {
	return runTerragruntCommandE(t, options, "hcl", "validate")
}
