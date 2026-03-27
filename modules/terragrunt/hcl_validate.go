package terragrunt //nolint:dupl // structural pattern for terragrunt command wrappers

import (
	"context"

	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// HclValidateContext runs terragrunt hcl validate to check terragrunt.hcl syntax.
// The provided context is passed through to the underlying command execution, allowing for timeout
// and cancellation control. This validates Terragrunt HCL configuration and can check for mis-aligned inputs.
// Use TerraformArgs to pass additional flags like "--inputs" or "--strict".
//
// Examples:
//
//	HclValidateContext(t, ctx, options)                                        // Basic syntax check
//	HclValidateContext(t, ctx, &Options{TerraformArgs: []string{"--inputs"}})  // Check input alignment
func HclValidateContext(t testing.TestingT, ctx context.Context, options *Options) string {
	out, err := HclValidateContextE(t, ctx, options)
	require.NoError(t, err)

	return out
}

// HclValidateContextE runs terragrunt hcl validate to check terragrunt.hcl syntax.
// The provided context is passed through to the underlying command execution, allowing for timeout
// and cancellation control. This validates Terragrunt HCL configuration and can check for mis-aligned inputs.
// Use TerraformArgs to pass additional flags like "--inputs" or "--strict".
func HclValidateContextE(t testing.TestingT, ctx context.Context, options *Options) (string, error) {
	return runTerragruntCommandE(t, ctx, options, "hcl", "validate")
}

// HclValidate runs terragrunt hcl validate to check terragrunt.hcl syntax.
// This validates Terragrunt HCL configuration and can check for mis-aligned inputs.
// Use TerraformArgs to pass additional flags like "--inputs" or "--strict".
//
// Deprecated: Use [HclValidateContext] instead.
func HclValidate(t testing.TestingT, options *Options) string {
	return HclValidateContext(t, context.Background(), options)
}

// HclValidateE runs terragrunt hcl validate to check terragrunt.hcl syntax.
// This validates Terragrunt HCL configuration and can check for mis-aligned inputs.
// Use TerraformArgs to pass additional flags like "--inputs" or "--strict".
//
// Deprecated: Use [HclValidateContextE] instead.
func HclValidateE(t testing.TestingT, options *Options) (string, error) {
	return HclValidateContextE(t, context.Background(), options)
}
