package helm

import (
	"context"
	"slices"

	"github.com/gruntwork-io/terratest/modules/shell"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
)

// getCommonArgs extracts common helm options. In this case, these are:
// - kubeconfig path
// - kubeconfig context
// - helm home path
func getCommonArgs(options *Options, args ...string) []string {
	if options.KubectlOptions != nil && options.KubectlOptions.ContextName != "" {
		args = append(args, "--kube-context", options.KubectlOptions.ContextName)
	}

	if options.KubectlOptions != nil && options.KubectlOptions.ConfigPath != "" {
		args = append(args, "--kubeconfig", options.KubectlOptions.ConfigPath)
	}

	if options.HomePath != "" {
		args = append(args, "--home", options.HomePath)
	}

	return args
}

// getNamespaceArgs returns the args to append for the namespace, if set in the helm Options struct.
func getNamespaceArgs(options *Options) []string {
	if options.KubectlOptions != nil && options.KubectlOptions.Namespace != "" {
		return []string{"--namespace", options.KubectlOptions.Namespace}
	}

	return []string{}
}

// mergeSetJSONValues merges the deprecated SetJsonValues field into SetJSONValues. Values from SetJSONValues take
// precedence when the same key appears in both maps.
func mergeSetJSONValues(options *Options) map[string]string {
	if len(options.SetJsonValues) == 0 { //nolint:staticcheck // Reading deprecated field for backwards compatibility.
		return options.SetJSONValues
	}

	if len(options.SetJSONValues) == 0 {
		return options.SetJsonValues //nolint:staticcheck // Reading deprecated field for backwards compatibility.
	}

	merged := make(map[string]string, len(options.SetJsonValues)+len(options.SetJSONValues)) //nolint:staticcheck // Reading deprecated field for backwards compatibility.

	for k, v := range options.SetJsonValues { //nolint:staticcheck // Reading deprecated field for backwards compatibility.
		merged[k] = v
	}

	// SetJSONValues takes precedence over the deprecated SetJsonValues.
	for k, v := range options.SetJSONValues {
		merged[k] = v
	}

	return merged
}

// getValuesArgsE computes the args to pass in for setting values.
func getValuesArgsE(options *Options, args ...string) ([]string, error) {
	args = append(args, FormatSetValuesAsArgs(options.SetValues, "--set")...)
	args = append(args, FormatSetValuesAsArgs(options.SetStrValues, "--set-string")...)
	args = append(args, FormatSetValuesAsArgs(mergeSetJSONValues(options), "--set-json")...)

	valuesFilesArgs, err := FormatValuesFilesAsArgsE(options.ValuesFiles)
	if err != nil {
		return args, err
	}

	args = append(args, valuesFilesArgs...)

	setFilesArgs, err := FormatSetFilesAsArgsE(options.SetFiles)
	if err != nil {
		return args, err
	}

	args = append(args, setFilesArgs...)

	return args, nil
}

// RunHelmCommandAndGetOutputE runs helm with the given arguments and options and returns combined, interleaved
// stdout/stderr.
//
// Deprecated: Use [RunHelmCommandAndGetOutputContextE] instead.
func RunHelmCommandAndGetOutputE(t testing.TestingT, options *Options, cmd string, additionalArgs ...string) (string, error) {
	return RunHelmCommandAndGetOutputContextE(t, context.Background(), options, cmd, additionalArgs...)
}

// RunHelmCommandAndGetOutputContextE runs helm with the given arguments and options and returns combined, interleaved
// stdout/stderr. The ctx parameter supports cancellation and timeouts.
func RunHelmCommandAndGetOutputContextE(t testing.TestingT, ctx context.Context, options *Options, cmd string, additionalArgs ...string) (string, error) {
	helmCmd := PrepareHelmCommand(options, cmd, additionalArgs...)

	return shell.RunCommandContextAndGetOutputE(t, ctx, helmCmd)
}

// RunHelmCommandAndGetStdOutE runs helm with the given arguments and options and returns stdout.
//
// Deprecated: Use [RunHelmCommandAndGetStdOutContextE] instead.
func RunHelmCommandAndGetStdOutE(t testing.TestingT, options *Options, cmd string, additionalArgs ...string) (string, error) {
	return RunHelmCommandAndGetStdOutContextE(t, context.Background(), options, cmd, additionalArgs...)
}

// RunHelmCommandAndGetStdOutContextE runs helm with the given arguments and options and returns stdout. The ctx
// parameter supports cancellation and timeouts.
func RunHelmCommandAndGetStdOutContextE(t testing.TestingT, ctx context.Context, options *Options, cmd string, additionalArgs ...string) (string, error) {
	helmCmd := PrepareHelmCommand(options, cmd, additionalArgs...)

	return shell.RunCommandContextAndGetStdOutE(t, ctx, helmCmd)
}

// RunHelmCommandAndGetStdOutErrE runs helm with the given arguments and options and returns stdout and stderr
// separately.
//
// Deprecated: Use [RunHelmCommandAndGetStdOutErrContextE] instead.
func RunHelmCommandAndGetStdOutErrE(t testing.TestingT, options *Options, cmd string, additionalArgs ...string) (string, string, error) {
	return RunHelmCommandAndGetStdOutErrContextE(t, context.Background(), options, cmd, additionalArgs...)
}

// RunHelmCommandAndGetStdOutErrContextE runs helm with the given arguments and options and returns stdout and stderr
// separately. The ctx parameter supports cancellation and timeouts.
func RunHelmCommandAndGetStdOutErrContextE(t testing.TestingT, ctx context.Context, options *Options, cmd string, additionalArgs ...string) (string, string, error) {
	helmCmd := PrepareHelmCommand(options, cmd, additionalArgs...)

	return shell.RunCommandContextAndGetStdOutErrE(t, ctx, helmCmd)
}

// PrepareHelmCommand builds a shell.Command for running helm with the given options, subcommand, and additional
// arguments.
func PrepareHelmCommand(options *Options, cmd string, additionalArgs ...string) *shell.Command {
	args := []string{cmd}
	args = getCommonArgs(options, args...)

	// namespace arg only appended if it is not already present
	if !slices.Contains(additionalArgs, "--namespace") {
		args = append(args, getNamespaceArgs(options)...)
	}

	args = append(args, additionalArgs...)

	helmCmd := &shell.Command{
		Command:    "helm",
		Args:       args,
		WorkingDir: ".",
		Env:        options.EnvVars,
		Logger:     options.Logger,
	}

	return helmCmd
}
