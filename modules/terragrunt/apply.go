package terragrunt

import (
	"github.com/gruntwork-io/terratest/modules/testing/v2"
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
