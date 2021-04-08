package helm

import (
	"github.com/gruntwork-io/go-commons/errors"
	"github.com/gruntwork-io/terratest/modules/shell"
	"github.com/gruntwork-io/terratest/modules/testing"
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

// getValuesArgsE computes the args to pass in for setting values
func getValuesArgsE(t testing.TestingT, options *Options, args ...string) ([]string, error) {
	args = append(args, formatSetValuesAsArgs(options.SetValues, "--set")...)
	args = append(args, formatSetValuesAsArgs(options.SetStrValues, "--set-string")...)

	valuesFilesArgs, err := formatValuesFilesAsArgsE(t, options.ValuesFiles)
	if err != nil {
		return args, errors.WithStackTrace(err)
	}
	args = append(args, valuesFilesArgs...)

	setFilesArgs, err := formatSetFilesAsArgsE(t, options.SetFiles)
	if err != nil {
		return args, errors.WithStackTrace(err)
	}
	args = append(args, setFilesArgs...)
	return args, nil
}

// RunHelmCommandAndGetOutputE runs helm with the given arguments and options and returns combined, interleaved stdout/stderr.
func RunHelmCommandAndGetOutputE(t testing.TestingT, options *Options, cmd string, additionalArgs ...string) (string, error) {
	helmCmd := prepareHelmCommand(t, options, cmd, additionalArgs...)
	return shell.RunCommandAndGetOutputE(t, helmCmd)
}

// RunHelmCommandAndGetStdOutE runs helm with the given arguments and options and returns stdout.
func RunHelmCommandAndGetStdOutE(t testing.TestingT, options *Options, cmd string, additionalArgs ...string) (string, error) {
	helmCmd := prepareHelmCommand(t, options, cmd, additionalArgs...)
	return shell.RunCommandAndGetStdOutE(t, helmCmd)
}

func prepareHelmCommand(t testing.TestingT, options *Options, cmd string, additionalArgs ...string) shell.Command {
	args := []string{cmd}
	args = getCommonArgs(options, args...)
	args = append(args, additionalArgs...)

	helmCmd := shell.Command{
		Command:    "helm",
		Args:       args,
		WorkingDir: ".",
		Env:        options.EnvVars,
		Logger:     options.Logger,
	}
	return helmCmd
}
