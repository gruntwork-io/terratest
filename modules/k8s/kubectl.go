package k8s

import (
	"context"
	"net/url"
	"os"

	"github.com/stretchr/testify/require"

	"github.com/gruntwork-io/terratest/modules/shell"
	"github.com/gruntwork-io/terratest/modules/testing"
)

// RunKubectlContext calls kubectl using the provided context, options, and args, failing the test on error.
func RunKubectlContext(t testing.TestingT, ctx context.Context, options *KubectlOptions, args ...string) {
	require.NoError(t, RunKubectlContextE(t, ctx, options, args...))
}

// RunKubectl calls kubectl using the provided options and args, failing the test on error.
//
// Deprecated: Use RunKubectlContext instead.
func RunKubectl(t testing.TestingT, options *KubectlOptions, args ...string) {
	RunKubectlContext(t, context.Background(), options, args...)
}

// RunKubectlContextE calls kubectl using the provided context, options, and args.
func RunKubectlContextE(t testing.TestingT, ctx context.Context, options *KubectlOptions, args ...string) error {
	_, err := RunKubectlAndGetOutputContextE(t, ctx, options, args...)
	return err
}

// RunKubectlE calls kubectl using the provided options and args.
//
// Deprecated: Use RunKubectlContextE instead.
func RunKubectlE(t testing.TestingT, options *KubectlOptions, args ...string) error {
	return RunKubectlContextE(t, context.Background(), options, args...)
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

// RunKubectlAndGetOutputE calls kubectl using the provided options and args, returning the combined output of
// stdout and stderr.
//
// Deprecated: Use RunKubectlAndGetOutputContextE instead.
func RunKubectlAndGetOutputE(t testing.TestingT, options *KubectlOptions, args ...string) (string, error) {
	return RunKubectlAndGetOutputContextE(t, context.Background(), options, args...)
}

// KubectlDelete will take in a file path and delete it from the cluster targeted by KubectlOptions. If there are any
// errors, fail the test immediately.
//
// Deprecated: Use [KubectlDeleteContext] instead.
func KubectlDelete(t testing.TestingT, options *KubectlOptions, configPath string) {
	KubectlDeleteContext(t, context.Background(), options, configPath)
}

// KubectlDeleteE will take in a file path and delete it from the cluster targeted by KubectlOptions.
//
// Deprecated: Use [KubectlDeleteContextE] instead.
func KubectlDeleteE(t testing.TestingT, options *KubectlOptions, configPath string) error {
	return KubectlDeleteContextE(t, context.Background(), options, configPath)
}

// KubectlDeleteContext deletes the resource at configPath from the cluster, using the provided context. Fails the test on error.
func KubectlDeleteContext(t testing.TestingT, ctx context.Context, options *KubectlOptions, configPath string) {
	require.NoError(t, KubectlDeleteContextE(t, ctx, options, configPath))
}

// KubectlDeleteContextE deletes the resource at configPath from the cluster, using the provided context.
func KubectlDeleteContextE(t testing.TestingT, ctx context.Context, options *KubectlOptions, configPath string) error {
	return RunKubectlContextE(t, ctx, options, "delete", "-f", configPath)
}

// KubectlDeleteFromKustomize will take in a kustomization directory path and delete it from the cluster targeted by KubectlOptions. If there are any
// errors, fail the test immediately.
//
// Deprecated: Use [KubectlDeleteFromKustomizeContext] instead.
func KubectlDeleteFromKustomize(t testing.TestingT, options *KubectlOptions, configPath string) {
	KubectlDeleteFromKustomizeContext(t, context.Background(), options, configPath)
}

// KubectlDeleteFromKustomizeE will take in a kustomization directory path and delete it from the cluster targeted by KubectlOptions.
//
// Deprecated: Use [KubectlDeleteFromKustomizeContextE] instead.
func KubectlDeleteFromKustomizeE(t testing.TestingT, options *KubectlOptions, configPath string) error {
	return KubectlDeleteFromKustomizeContextE(t, context.Background(), options, configPath)
}

// KubectlDeleteFromKustomizeContext deletes the kustomization at configPath from the cluster, using the provided context. Fails the test on error.
func KubectlDeleteFromKustomizeContext(t testing.TestingT, ctx context.Context, options *KubectlOptions, configPath string) {
	require.NoError(t, KubectlDeleteFromKustomizeContextE(t, ctx, options, configPath))
}

// KubectlDeleteFromKustomizeContextE deletes the kustomization at configPath from the cluster, using the provided context.
func KubectlDeleteFromKustomizeContextE(t testing.TestingT, ctx context.Context, options *KubectlOptions, configPath string) error {
	return RunKubectlContextE(t, ctx, options, "delete", "-k", configPath)
}

// KubectlDeleteFromString will take in a kubernetes resource config as a string and delete it on the cluster specified
// by the provided kubectl options.
//
// Deprecated: Use [KubectlDeleteFromStringContext] instead.
func KubectlDeleteFromString(t testing.TestingT, options *KubectlOptions, configData string) {
	KubectlDeleteFromStringContext(t, context.Background(), options, configData)
}

// KubectlDeleteFromStringE will take in a kubernetes resource config as a string and delete it on the cluster specified
// by the provided kubectl options. If it fails, this will return the error.
//
// Deprecated: Use [KubectlDeleteFromStringContextE] instead.
func KubectlDeleteFromStringE(t testing.TestingT, options *KubectlOptions, configData string) error {
	return KubectlDeleteFromStringContextE(t, context.Background(), options, configData)
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

// KubectlApply will take in a file path and apply it to the cluster targeted by KubectlOptions. If there are any
// errors, fail the test immediately.
//
// Deprecated: Use [KubectlApplyContext] instead.
func KubectlApply(t testing.TestingT, options *KubectlOptions, configPath string) {
	KubectlApplyContext(t, context.Background(), options, configPath)
}

// KubectlApplyE will take in a file path and apply it to the cluster targeted by KubectlOptions.
//
// Deprecated: Use [KubectlApplyContextE] instead.
func KubectlApplyE(t testing.TestingT, options *KubectlOptions, configPath string) error {
	return KubectlApplyContextE(t, context.Background(), options, configPath)
}

// KubectlApplyContext applies the resource at configPath to the cluster, using the provided context. Fails the test on error.
func KubectlApplyContext(t testing.TestingT, ctx context.Context, options *KubectlOptions, configPath string) {
	require.NoError(t, KubectlApplyContextE(t, ctx, options, configPath))
}

// KubectlApplyContextE applies the resource at configPath to the cluster, using the provided context.
func KubectlApplyContextE(t testing.TestingT, ctx context.Context, options *KubectlOptions, configPath string) error {
	return RunKubectlContextE(t, ctx, options, "apply", "-f", configPath)
}

// KubectlApplyFromKustomize will take in a kustomization directory path and apply it to the cluster targeted by KubectlOptions. If there are any
// errors, fail the test immediately.
//
// Deprecated: Use [KubectlApplyFromKustomizeContext] instead.
func KubectlApplyFromKustomize(t testing.TestingT, options *KubectlOptions, configPath string) {
	KubectlApplyFromKustomizeContext(t, context.Background(), options, configPath)
}

// KubectlApplyFromKustomizeE will take in a kustomization directory path and apply it to the cluster targeted by KubectlOptions.
//
// Deprecated: Use [KubectlApplyFromKustomizeContextE] instead.
func KubectlApplyFromKustomizeE(t testing.TestingT, options *KubectlOptions, configPath string) error {
	return KubectlApplyFromKustomizeContextE(t, context.Background(), options, configPath)
}

// KubectlApplyFromKustomizeContext applies the kustomization at configPath to the cluster, using the provided context. Fails the test on error.
func KubectlApplyFromKustomizeContext(t testing.TestingT, ctx context.Context, options *KubectlOptions, configPath string) {
	require.NoError(t, KubectlApplyFromKustomizeContextE(t, ctx, options, configPath))
}

// KubectlApplyFromKustomizeContextE applies the kustomization at configPath to the cluster, using the provided context.
func KubectlApplyFromKustomizeContextE(t testing.TestingT, ctx context.Context, options *KubectlOptions, configPath string) error {
	return RunKubectlContextE(t, ctx, options, "apply", "-k", configPath)
}

// KubectlApplyFromString will take in a kubernetes resource config as a string and apply it on the cluster specified
// by the provided kubectl options.
//
// Deprecated: Use [KubectlApplyFromStringContext] instead.
func KubectlApplyFromString(t testing.TestingT, options *KubectlOptions, configData string) {
	KubectlApplyFromStringContext(t, context.Background(), options, configData)
}

// KubectlApplyFromStringE will take in a kubernetes resource config as a string and apply it on the cluster specified
// by the provided kubectl options. If it fails, this will return the error.
//
// Deprecated: Use [KubectlApplyFromStringContextE] instead.
func KubectlApplyFromStringE(t testing.TestingT, options *KubectlOptions, configData string) error {
	return KubectlApplyFromStringContextE(t, context.Background(), options, configData)
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
