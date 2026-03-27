package terragrunt

import (
	"context"

	"github.com/gruntwork-io/terratest/modules/testing"
)

// RunAllContext runs terragrunt run --all -- <command> with the given options and returns stdout/stderr.
// The provided context is passed through to the underlying command execution, allowing for timeout
// and cancellation control.
//
// Deprecated: Use [RunContext] with the --all flag in tgArgs instead.
func RunAllContext(t testing.TestingT, ctx context.Context, options *Options, command string) string {
	return RunContext(t, ctx, options, []string{"--all"}, []string{command})
}

// RunAllContextE runs terragrunt run --all -- <command> with the given options and returns stdout/stderr.
// The provided context is passed through to the underlying command execution, allowing for timeout
// and cancellation control.
//
// Deprecated: Use [RunContextE] with the --all flag in tgArgs instead.
func RunAllContextE(t testing.TestingT, ctx context.Context, options *Options, command string) (string, error) {
	return RunContextE(t, ctx, options, []string{"--all"}, []string{command})
}

// RunAll runs terragrunt run --all -- <command> with the given options and returns stdout/stderr.
//
// Deprecated: Use [RunContext] with the --all flag in tgArgs instead.
func RunAll(t testing.TestingT, options *Options, command string) string {
	return RunAllContext(t, context.Background(), options, command)
}

// RunAllE runs terragrunt run --all -- <command> with the given options and returns stdout/stderr.
//
// Deprecated: Use [RunContextE] with the --all flag in tgArgs instead.
func RunAllE(t testing.TestingT, options *Options, command string) (string, error) {
	return RunAllContextE(t, context.Background(), options, command)
}
