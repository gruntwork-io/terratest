package terraform

import (
	"errors"

	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// InitAndApply runs terraform init and apply with the given options and return stdout/stderr from the apply command. Note that this
// method does NOT call destroy and assumes the caller is responsible for cleaning up any resources created by running
// apply.
func InitAndApply(t testing.TestingT, options *Options) string {
	out, err := InitAndApplyE(t, options)
	require.NoError(t, err)
	return out
}

// InitAndApplyE runs terraform init and apply with the given options and return stdout/stderr from the apply command. Note that this
// method does NOT call destroy and assumes the caller is responsible for cleaning up any resources created by running
// apply.
func InitAndApplyE(t testing.TestingT, options *Options) (string, error) {
	if _, err := InitE(t, options); err != nil {
		return "", err
	}

	if _, err := GetE(t, options); err != nil {
		return "", err
	}

	return ApplyE(t, options)
}

// Apply runs terraform apply with the given options and return stdout/stderr. Note that this method does NOT call destroy and
// assumes the caller is responsible for cleaning up any resources created by running apply.
func Apply(t testing.TestingT, options *Options) string {
	out, err := ApplyE(t, options)
	require.NoError(t, err)
	return out
}

// TgApplyAll runs terragrunt apply with the given options and return stdout/stderr. Note that this method does NOT call destroy and
// assumes the caller is responsible for cleaning up any resources created by running apply.
func TgApplyAll(t testing.TestingT, options *Options) string {
	out, err := TgApplyAllE(t, options)
	require.NoError(t, err)
	return out
}

// ApplyE runs terraform apply with the given options and return stdout/stderr. Note that this method does NOT call destroy and
// assumes the caller is responsible for cleaning up any resources created by running apply.
func ApplyE(t testing.TestingT, options *Options) (string, error) {
	return RunTerraformCommandE(t, options, FormatArgs(options, "apply", "-input=false", "-auto-approve")...)
}

// TgApplyAllE runs terragrunt apply-all with the given options and return stdout/stderr. Note that this method does NOT call destroy and
// assumes the caller is responsible for cleaning up any resources created by running apply.
func TgApplyAllE(t testing.TestingT, options *Options) (string, error) {
	if options.TerraformBinary != "terragrunt" {
		return "", TgInvalidBinary(options.TerraformBinary)
	}

	return RunTerraformCommandE(t, options, FormatArgs(options, "apply-all", "-input=false", "-lock=false", "-auto-approve")...)
}

// ApplyAndIdempotent runs terraform apply with the given options and return stdout/stderr from the apply command. It then runs
// plan again and will fail the test if plan requires additional changes. Note that this method does NOT call destroy and assumes
// the caller is responsible for cleaning up any resources created by running apply.
func ApplyAndIdempotent(t testing.TestingT, options *Options) string {
	out, err := ApplyAndIdempotentE(t, options)
	require.NoError(t, err)

	return out
}

// ApplyAndIdempotentE runs terraform apply with the given options and return stdout/stderr from the apply command. It then runs
// plan again and will fail the test if plan requires additional changes. Note that this method does NOT call destroy and assumes
// the caller is responsible for cleaning up any resources created by running apply.
func ApplyAndIdempotentE(t testing.TestingT, options *Options) (string, error) {
	out, err := ApplyE(t, options)

	if err != nil {
		return out, err
	}

	exitCode, err := PlanExitCodeE(t, options)

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
func InitAndApplyAndIdempotent(t testing.TestingT, options *Options) string {
	out, err := InitAndApplyAndIdempotentE(t, options)
	require.NoError(t, err)

	return out
}

// InitAndApplyAndIdempotentE runs terraform init and apply with the given options and return stdout/stderr from the apply command. It then runs
// plan again and will fail the test if plan requires additional changes. Note that this method does NOT call destroy and assumes
// the caller is responsible for cleaning up any resources created by running apply.
func InitAndApplyAndIdempotentE(t testing.TestingT, options *Options) (string, error) {
	if _, err := InitE(t, options); err != nil {
		return "", err
	}

	return ApplyAndIdempotentE(t, options)
}
