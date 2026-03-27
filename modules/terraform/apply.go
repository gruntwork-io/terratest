package terraform

import (
	"context"
	"errors"

	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// InitAndApply runs terraform init and apply with the given options and return stdout/stderr from the apply command. Note that this
// method does NOT call destroy and assumes the caller is responsible for cleaning up any resources created by running
// apply.
//
// Deprecated: Use [InitAndApplyContext] instead.
func InitAndApply(t testing.TestingT, options *Options) string {
	return InitAndApplyContext(t, context.Background(), options)
}

// InitAndApplyContext runs terraform init and apply with the given options and returns stdout/stderr from the apply
// command. The provided context is passed through to the underlying command execution, allowing for timeout and
// cancellation control. Note that this method does NOT call destroy and assumes the caller is responsible for cleaning
// up any resources created by running apply.
func InitAndApplyContext(t testing.TestingT, ctx context.Context, options *Options) string {
	out, err := InitAndApplyContextE(t, ctx, options)
	require.NoError(t, err)

	return out
}

// InitAndApplyE runs terraform init and apply with the given options and return stdout/stderr from the apply command. Note that this
// method does NOT call destroy and assumes the caller is responsible for cleaning up any resources created by running
// apply.
//
// Deprecated: Use [InitAndApplyContextE] instead.
func InitAndApplyE(t testing.TestingT, options *Options) (string, error) {
	return InitAndApplyContextE(t, context.Background(), options)
}

// InitAndApplyContextE runs terraform init and apply with the given options and returns stdout/stderr from the apply
// command. The provided context is passed through to the underlying command execution, allowing for timeout and
// cancellation control. Note that this method does NOT call destroy and assumes the caller is responsible for cleaning
// up any resources created by running apply.
func InitAndApplyContextE(t testing.TestingT, ctx context.Context, options *Options) (string, error) {
	if _, err := InitContextE(t, ctx, options); err != nil {
		return "", err
	}

	return ApplyContextE(t, ctx, options)
}

// Apply runs terraform apply with the given options and return stdout/stderr. Note that this method does NOT call destroy and
// assumes the caller is responsible for cleaning up any resources created by running apply.
//
// Deprecated: Use [ApplyContext] instead.
func Apply(t testing.TestingT, options *Options) string {
	return ApplyContext(t, context.Background(), options)
}

// ApplyContext runs terraform apply with the given options and returns stdout/stderr. The provided context is passed
// through to the underlying command execution, allowing for timeout and cancellation control. Note that this method
// does NOT call destroy and assumes the caller is responsible for cleaning up any resources created by running apply.
func ApplyContext(t testing.TestingT, ctx context.Context, options *Options) string {
	out, err := ApplyContextE(t, ctx, options)
	require.NoError(t, err)

	return out
}

// ApplyE runs terraform apply with the given options and return stdout/stderr. Note that this method does NOT call destroy and
// assumes the caller is responsible for cleaning up any resources created by running apply.
//
// Deprecated: Use [ApplyContextE] instead.
func ApplyE(t testing.TestingT, options *Options) (string, error) {
	return ApplyContextE(t, context.Background(), options)
}

// ApplyContextE runs terraform apply with the given options and returns stdout/stderr. The provided context is passed
// through to the underlying command execution, allowing for timeout and cancellation control. Note that this method
// does NOT call destroy and assumes the caller is responsible for cleaning up any resources created by running apply.
func ApplyContextE(t testing.TestingT, ctx context.Context, options *Options) (string, error) {
	return RunTerraformCommandContextE(t, ctx, options, FormatArgs(options, prepend(options.ExtraArgs.Apply, "apply", "-input=false", "-auto-approve")...)...)
}

// ApplyAndIdempotent runs terraform apply with the given options and return stdout/stderr from the apply command. It then runs
// plan again and will fail the test if plan requires additional changes. Note that this method does NOT call destroy and assumes
// the caller is responsible for cleaning up any resources created by running apply.
//
// Deprecated: Use [ApplyAndIdempotentContext] instead.
func ApplyAndIdempotent(t testing.TestingT, options *Options) string {
	return ApplyAndIdempotentContext(t, context.Background(), options)
}

// ApplyAndIdempotentContext runs terraform apply with the given options and returns stdout/stderr from the apply
// command. It then runs plan again and will fail the test if plan requires additional changes. The provided context is
// passed through to the underlying command execution, allowing for timeout and cancellation control. Note that this
// method does NOT call destroy and assumes the caller is responsible for cleaning up any resources created by running
// apply.
func ApplyAndIdempotentContext(t testing.TestingT, ctx context.Context, options *Options) string {
	out, err := ApplyAndIdempotentContextE(t, ctx, options)
	require.NoError(t, err)

	return out
}

// ApplyAndIdempotentE runs terraform apply with the given options and return stdout/stderr from the apply command. It then runs
// plan again and will fail the test if plan requires additional changes. Note that this method does NOT call destroy and assumes
// the caller is responsible for cleaning up any resources created by running apply.
//
// Deprecated: Use [ApplyAndIdempotentContextE] instead.
func ApplyAndIdempotentE(t testing.TestingT, options *Options) (string, error) {
	return ApplyAndIdempotentContextE(t, context.Background(), options)
}

// ApplyAndIdempotentContextE runs terraform apply with the given options and returns stdout/stderr from the apply
// command. It then runs plan again and will fail the test if plan requires additional changes. The provided context is
// passed through to the underlying command execution, allowing for timeout and cancellation control. Note that this
// method does NOT call destroy and assumes the caller is responsible for cleaning up any resources created by running
// apply.
func ApplyAndIdempotentContextE(t testing.TestingT, ctx context.Context, options *Options) (string, error) {
	out, err := ApplyContextE(t, ctx, options)
	if err != nil {
		return out, err
	}

	exitCode, err := PlanExitCodeContextE(t, ctx, options)
	if err != nil {
		return out, err
	}

	if exitCode != 0 {
		return out, errors.New("terraform configuration not idempotent")
	}

	return out, nil
}

// InitAndApplyAndIdempotent runs terraform init and apply with the given options and return stdout/stderr from the apply command. It then runs
// plan again and will fail the test if plan requires additional changes. Note that this method does NOT call destroy and assumes
// the caller is responsible for cleaning up any resources created by running apply.
//
// Deprecated: Use [InitAndApplyAndIdempotentContext] instead.
func InitAndApplyAndIdempotent(t testing.TestingT, options *Options) string {
	return InitAndApplyAndIdempotentContext(t, context.Background(), options)
}

// InitAndApplyAndIdempotentContext runs terraform init, apply, and then plan with the given options and returns
// stdout/stderr from the apply command. It will fail the test if plan requires additional changes after the apply. The
// provided context is passed through to the underlying command execution, allowing for timeout and cancellation
// control. Note that this method does NOT call destroy and assumes the caller is responsible for cleaning up any
// resources created by running apply.
func InitAndApplyAndIdempotentContext(t testing.TestingT, ctx context.Context, options *Options) string {
	out, err := InitAndApplyAndIdempotentContextE(t, ctx, options)
	require.NoError(t, err)

	return out
}

// InitAndApplyAndIdempotentE runs terraform init and apply with the given options and return stdout/stderr from the apply command. It then runs
// plan again and will fail the test if plan requires additional changes. Note that this method does NOT call destroy and assumes
// the caller is responsible for cleaning up any resources created by running apply.
//
// Deprecated: Use [InitAndApplyAndIdempotentContextE] instead.
func InitAndApplyAndIdempotentE(t testing.TestingT, options *Options) (string, error) {
	return InitAndApplyAndIdempotentContextE(t, context.Background(), options)
}

// InitAndApplyAndIdempotentContextE runs terraform init, apply, and then plan with the given options and returns
// stdout/stderr from the apply command. It will fail the test if plan requires additional changes after the apply. The
// provided context is passed through to the underlying command execution, allowing for timeout and cancellation
// control. Note that this method does NOT call destroy and assumes the caller is responsible for cleaning up any
// resources created by running apply.
func InitAndApplyAndIdempotentContextE(t testing.TestingT, ctx context.Context, options *Options) (string, error) {
	if _, err := InitContextE(t, ctx, options); err != nil {
		return "", err
	}

	return ApplyAndIdempotentContextE(t, ctx, options)
}
