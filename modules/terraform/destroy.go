package terraform

import (
	"github.com/gruntwork-io/terratest/modules/testing/v2"
	"github.com/stretchr/testify/require"
)

// Destroy runs terraform destroy with the given options and return stdout/stderr.
func Destroy(t testing.TestingT, options *Options) string {
	out, err := DestroyE(t, options)
	require.NoError(t, err)
	return out
}

// TgDestroyAll runs terragrunt destroy with the given options and return stdout.
//
// Deprecated: Use github.com/gruntwork-io/terratest/modules/terragrunt.DestroyAll instead.
func TgDestroyAll(t testing.TestingT, options *Options) string {
	out, err := TgDestroyAllE(t, options)
	require.NoError(t, err)
	return out
}

// DestroyE runs terraform destroy with the given options and return stdout/stderr.
func DestroyE(t testing.TestingT, options *Options) (string, error) {
	return RunTerraformCommandE(t, options, FormatArgs(options, prepend(options.ExtraArgs.Destroy, "destroy", "-auto-approve", "-input=false")...)...)
}

// TgDestroyAllE runs terragrunt destroy --all with the given options and return stdout.
//
// Deprecated: Use github.com/gruntwork-io/terratest/modules/terragrunt.DestroyAllE instead.
func TgDestroyAllE(t testing.TestingT, options *Options) (string, error) {
	if options.TerraformBinary != "terragrunt" {
		return "", TgInvalidBinary(options.TerraformBinary)
	}

	return RunTerraformCommandE(t, options, FormatArgs(options, prepend(options.ExtraArgs.Destroy, "destroy", "--all", "-auto-approve", "-input=false")...)...)
}
