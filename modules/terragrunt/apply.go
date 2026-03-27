package terragrunt

import (
	"context"

	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// ApplyAllContext runs terragrunt run --all apply with the given options and returns stdout/stderr.
// The provided context is passed through to the underlying command execution, allowing for timeout
// and cancellation control. Note that this method does NOT call destroy and assumes the caller is
// responsible for cleaning up any resources created by running apply.
func ApplyAllContext(t testing.TestingT, ctx context.Context, options *Options) string {
	out, err := ApplyAllContextE(t, ctx, options)
	require.NoError(t, err)

	return out
}

// ApplyAllContextE runs terragrunt run --all apply with the given options and returns stdout/stderr.
// The provided context is passed through to the underlying command execution, allowing for timeout
// and cancellation control. Note that this method does NOT call destroy and assumes the caller is
// responsible for cleaning up any resources created by running apply.
func ApplyAllContextE(t testing.TestingT, ctx context.Context, options *Options) (string, error) {
	args := BuildRunArgs([]string{"--all"}, []string{"apply", "-input=false", "-auto-approve"})

	return runTerragruntCommandE(t, ctx, options, "run", args...)
}

// ApplyAll runs terragrunt run --all apply with the given options and returns stdout/stderr. Note that this method does NOT call destroy and
// assumes the caller is responsible for cleaning up any resources created by running apply.
//
// Deprecated: Use [ApplyAllContext] instead.
func ApplyAll(t testing.TestingT, options *Options) string {
	return ApplyAllContext(t, context.Background(), options)
}

// ApplyAllE runs terragrunt run --all -- apply with the given options and returns stdout/stderr. Note that this method does NOT call destroy and
// assumes the caller is responsible for cleaning up any resources created by running apply.
//
// Deprecated: Use [ApplyAllContextE] instead.
func ApplyAllE(t testing.TestingT, options *Options) (string, error) {
	return ApplyAllContextE(t, context.Background(), options)
}

// ApplyContext runs terragrunt run apply for a single unit and returns stdout/stderr.
// The provided context is passed through to the underlying command execution, allowing for timeout
// and cancellation control.
func ApplyContext(t testing.TestingT, ctx context.Context, options *Options) string {
	out, err := ApplyContextE(t, ctx, options)
	require.NoError(t, err)

	return out
}

// ApplyContextE runs terragrunt run -- apply for a single unit and returns stdout/stderr.
// The provided context is passed through to the underlying command execution, allowing for timeout
// and cancellation control.
func ApplyContextE(t testing.TestingT, ctx context.Context, options *Options) (string, error) {
	args := BuildRunArgs([]string{}, []string{"apply", "-input=false", "-auto-approve"})

	return runTerragruntCommandE(t, ctx, options, "run", args...)
}

// Apply runs terragrunt run apply for a single unit and returns stdout/stderr.
//
// Deprecated: Use [ApplyContext] instead.
func Apply(t testing.TestingT, options *Options) string {
	return ApplyContext(t, context.Background(), options)
}

// ApplyE runs terragrunt run -- apply for a single unit and returns stdout/stderr.
//
// Deprecated: Use [ApplyContextE] instead.
func ApplyE(t testing.TestingT, options *Options) (string, error) {
	return ApplyContextE(t, context.Background(), options)
}

// InitAndApplyContext runs terragrunt init followed by apply for a single unit and returns the apply stdout/stderr.
// The provided context is passed through to both the init and apply command executions.
func InitAndApplyContext(t testing.TestingT, ctx context.Context, options *Options) string {
	out, err := InitAndApplyContextE(t, ctx, options)
	require.NoError(t, err)

	return out
}

// InitAndApplyContextE runs terragrunt init followed by apply for a single unit and returns the apply stdout/stderr.
// The provided context is passed through to both the init and apply command executions.
func InitAndApplyContextE(t testing.TestingT, ctx context.Context, options *Options) (string, error) {
	if _, err := InitContextE(t, ctx, options); err != nil {
		return "", err
	}

	return ApplyContextE(t, ctx, options)
}

// InitAndApply runs terragrunt init followed by apply for a single unit and returns the apply stdout/stderr.
//
// Deprecated: Use [InitAndApplyContext] instead.
func InitAndApply(t testing.TestingT, options *Options) string {
	return InitAndApplyContext(t, context.Background(), options)
}

// InitAndApplyE runs terragrunt init followed by apply for a single unit and returns the apply stdout/stderr.
//
// Deprecated: Use [InitAndApplyContextE] instead.
func InitAndApplyE(t testing.TestingT, options *Options) (string, error) {
	return InitAndApplyContextE(t, context.Background(), options)
}
