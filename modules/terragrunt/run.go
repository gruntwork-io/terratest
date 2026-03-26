package terragrunt

import (
	"context"

	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// RunContext runs terragrunt run [tgArgs...] -- [tfArgs...] with the given options and returns stdout/stderr.
// The provided context is passed through to the underlying command execution, allowing for timeout
// and cancellation control. This is a generic wrapper that allows running any OpenTofu/Terraform command
// through terragrunt run. The -- separator disambiguates Terragrunt flags from OpenTofu/Terraform flags.
// The OpenTofu/Terraform command (e.g. "apply") should be the first element of tfArgs.
func RunContext(t testing.TestingT, ctx context.Context, options *Options, tgArgs []string, tfArgs []string) string {
	out, err := RunContextE(t, ctx, options, tgArgs, tfArgs)
	require.NoError(t, err)

	return out
}

// RunContextE runs terragrunt run [tgArgs...] -- [tfArgs...] with the given options and returns stdout/stderr.
// The provided context is passed through to the underlying command execution, allowing for timeout
// and cancellation control. This is a generic wrapper that allows running any OpenTofu/Terraform command
// through terragrunt run. The -- separator disambiguates Terragrunt flags from OpenTofu/Terraform flags.
// The OpenTofu/Terraform command (e.g. "apply") should be the first element of tfArgs.
func RunContextE(t testing.TestingT, ctx context.Context, options *Options, tgArgs []string, tfArgs []string) (string, error) {
	if len(tfArgs) == 0 {
		return "", ErrEmptyTfArgs
	}

	args := BuildRunArgs(tgArgs, tfArgs)

	return runTerragruntCommandE(t, ctx, options, "run", args...)
}

// Run runs terragrunt run [tgArgs...] -- [tfArgs...] with the given options and returns stdout/stderr.
// This is a generic wrapper that allows running any OpenTofu/Terraform command through terragrunt run.
// The -- separator disambiguates Terragrunt flags from OpenTofu/Terraform flags.
// The OpenTofu/Terraform command (e.g. "apply") should be the first element of tfArgs.
//
// Deprecated: Use [RunContext] instead.
func Run(t testing.TestingT, options *Options, tgArgs []string, tfArgs []string) string {
	return RunContext(t, context.Background(), options, tgArgs, tfArgs)
}

// RunE runs terragrunt run [tgArgs...] -- [tfArgs...] with the given options and returns stdout/stderr.
// This is a generic wrapper that allows running any OpenTofu/Terraform command through terragrunt run.
// The -- separator disambiguates Terragrunt flags from OpenTofu/Terraform flags.
// The OpenTofu/Terraform command (e.g. "apply") should be the first element of tfArgs.
//
// Deprecated: Use [RunContextE] instead.
func RunE(t testing.TestingT, options *Options, tgArgs []string, tfArgs []string) (string, error) {
	return RunContextE(t, context.Background(), options, tgArgs, tfArgs)
}
