package terragrunt

import (
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// ApplyAll runs terragrunt apply --all with the given options and returns stdout/stderr. Note that this method does NOT call destroy and
// assumes the caller is responsible for cleaning up any resources created by running apply.
func ApplyAll(t testing.TestingT, options *Options) string {
	out, err := ApplyAllE(t, options)
	require.NoError(t, err)
	return out
}

// ApplyAllE runs terragrunt apply --all with the given options and returns stdout/stderr. Note that this method does NOT call destroy and
// assumes the caller is responsible for cleaning up any resources created by running apply.
func ApplyAllE(t testing.TestingT, options *Options) (string, error) {
	return runTerragruntCommandE(t, options, "apply", "--all", "-input=false", "-auto-approve")
}

// Apply runs terragrunt apply for a single unit and returns stdout/stderr.
func Apply(t testing.TestingT, options *Options) string {
	out, err := ApplyE(t, options)
	require.NoError(t, err)
	return out
}

// ApplyE runs terragrunt apply for a single unit and returns stdout/stderr.
func ApplyE(t testing.TestingT, options *Options) (string, error) {
	return runTerragruntCommandE(t, options, "apply", "-input=false", "-auto-approve")
}

// InitAndApply runs terragrunt init followed by apply for a single unit and returns the apply stdout/stderr.
func InitAndApply(t testing.TestingT, options *Options) string {
	out, err := InitAndApplyE(t, options)
	require.NoError(t, err)
	return out
}

// InitAndApplyE runs terragrunt init followed by apply for a single unit and returns the apply stdout/stderr.
func InitAndApplyE(t testing.TestingT, options *Options) (string, error) {
	if _, err := InitE(t, options); err != nil {
		return "", err
	}
	return ApplyE(t, options)
}
