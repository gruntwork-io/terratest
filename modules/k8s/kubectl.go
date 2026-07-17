package k8s

import (
	"context"
	"net/url"
	"os"

	"github.com/stretchr/testify/require"

	"github.com/gruntwork-io/terratest/modules/core/v2/shell"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
)

// RunKubectlContext calls kubectl using the provided context, options, and args, failing the test on error.
func RunKubectlContext(t testing.TestingT, ctx context.Context, options *KubectlOptions, args ...string) {
	require.NoError(t, RunKubectlContextE(t, ctx, options, args...))
}

// RunKubectlContextE calls kubectl using the provided context, options, and args.
func RunKubectlContextE(t testing.TestingT, ctx context.Context, options *KubectlOptions, args ...string) error {
	_, err := RunKubectlAndGetOutputContextE(t, ctx, options, args...)
	return err
}

// RunKubectlAndGetOutputContextE calls kubectl using the provided context, options, and args, returning the
// combined output of stdout and stderr.
func RunKubectlAndGetOutputContextE(t testing.TestingT, ctx context.Context, options *KubectlOptions, args ...string) (string, error) {
	cmdArgs := []string{}
	if options.ContextName != "" {
		cmdArgs = append(cmdArgs, "--context", options.ContextName)
	}

	if options.ConfigPath != "" {
		cmdArgs = append(cmdArgs, "--kubeconfig", options.ConfigPath)
	}

	if options.Namespace != "" {
		cmdArgs = append(cmdArgs, "--namespace", options.Namespace)
	}

	if options.RequestTimeout > 0 {
		cmdArgs = append(cmdArgs, "--request-timeout", options.RequestTimeout.String())
	}

	cmdArgs = append(cmdArgs, args...)
	command := &shell.Command{
		Command: "kubectl",
		Args:    cmdArgs,
		Env:     options.Env,
		Logger:  options.Logger,
	}

	return shell.RunCommandContextAndGetOutputE(t, ctx, command)
}

// KubectlDeleteContext deletes the resource at configPath from the cluster, using the provided context. Fails the test on error.
func KubectlDeleteContext(t testing.TestingT, ctx context.Context, options *KubectlOptions, configPath string) {
	require.NoError(t, KubectlDeleteContextE(t, ctx, options, configPath))
}

// KubectlDeleteContextE deletes the resource at configPath from the cluster, using the provided context.
func KubectlDeleteContextE(t testing.TestingT, ctx context.Context, options *KubectlOptions, configPath string) error {
	return RunKubectlContextE(t, ctx, options, "delete", "-f", configPath)
}

// KubectlDeleteFromKustomizeContext deletes the kustomization at configPath from the cluster, using the provided context. Fails the test on error.
func KubectlDeleteFromKustomizeContext(t testing.TestingT, ctx context.Context, options *KubectlOptions, configPath string) {
	require.NoError(t, KubectlDeleteFromKustomizeContextE(t, ctx, options, configPath))
}

// KubectlDeleteFromKustomizeContextE deletes the kustomization at configPath from the cluster, using the provided context.
func KubectlDeleteFromKustomizeContextE(t testing.TestingT, ctx context.Context, options *KubectlOptions, configPath string) error {
	return RunKubectlContextE(t, ctx, options, "delete", "-k", configPath)
}

// KubectlDeleteFromStringContext deletes the kubernetes resource from configData on the cluster, using the provided context. Fails the test on error.
func KubectlDeleteFromStringContext(t testing.TestingT, ctx context.Context, options *KubectlOptions, configData string) {
	require.NoError(t, KubectlDeleteFromStringContextE(t, ctx, options, configData))
}

// KubectlDeleteFromStringContextE deletes the kubernetes resource from configData on the cluster, using the provided context.
func KubectlDeleteFromStringContextE(t testing.TestingT, ctx context.Context, options *KubectlOptions, configData string) error {
	tmpfile, err := StoreConfigToTempFileE(t, configData)
	if err != nil {
		return err
	}

	defer func() { _ = os.Remove(tmpfile) }()

	return KubectlDeleteContextE(t, ctx, options, tmpfile)
}

// KubectlApplyContext applies the resource at configPath to the cluster, using the provided context. Fails the test on error.
func KubectlApplyContext(t testing.TestingT, ctx context.Context, options *KubectlOptions, configPath string) {
	require.NoError(t, KubectlApplyContextE(t, ctx, options, configPath))
}

// KubectlApplyContextE applies the resource at configPath to the cluster, using the provided context.
func KubectlApplyContextE(t testing.TestingT, ctx context.Context, options *KubectlOptions, configPath string) error {
	return RunKubectlContextE(t, ctx, options, "apply", "-f", configPath)
}

// KubectlApplyFromKustomizeContext applies the kustomization at configPath to the cluster, using the provided context. Fails the test on error.
func KubectlApplyFromKustomizeContext(t testing.TestingT, ctx context.Context, options *KubectlOptions, configPath string) {
	require.NoError(t, KubectlApplyFromKustomizeContextE(t, ctx, options, configPath))
}

// KubectlApplyFromKustomizeContextE applies the kustomization at configPath to the cluster, using the provided context.
func KubectlApplyFromKustomizeContextE(t testing.TestingT, ctx context.Context, options *KubectlOptions, configPath string) error {
	return RunKubectlContextE(t, ctx, options, "apply", "-k", configPath)
}

// KubectlApplyFromStringContext applies the kubernetes resource from configData to the cluster, using the provided context. Fails the test on error.
func KubectlApplyFromStringContext(t testing.TestingT, ctx context.Context, options *KubectlOptions, configData string) {
	require.NoError(t, KubectlApplyFromStringContextE(t, ctx, options, configData))
}

// KubectlApplyFromStringContextE applies the kubernetes resource from configData to the cluster, using the provided context.
func KubectlApplyFromStringContextE(t testing.TestingT, ctx context.Context, options *KubectlOptions, configData string) error {
	tmpfile, err := StoreConfigToTempFileE(t, configData)
	if err != nil {
		return err
	}

	defer func() { _ = os.Remove(tmpfile) }()

	return KubectlApplyContextE(t, ctx, options, tmpfile)
}

// StoreConfigToTempFile will store the provided config data to a temporary file created on the os and return the
// filename.
func StoreConfigToTempFile(t testing.TestingT, configData string) string {
	out, err := StoreConfigToTempFileE(t, configData)
	require.NoError(t, err)

	return out
}

// StoreConfigToTempFileE will store the provided config data to a temporary file created on the os and return the
// filename, or error.
func StoreConfigToTempFileE(t testing.TestingT, configData string) (string, error) {
	escapedTestName := url.PathEscape(t.Name())

	tmpfile, err := os.CreateTemp("", escapedTestName)
	if err != nil {
		return "", err
	}
	defer tmpfile.Close()

	_, err = tmpfile.WriteString(configData)

	return tmpfile.Name(), err
}
