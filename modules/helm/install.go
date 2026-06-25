package helm

import (
	"context"
	"path/filepath"

	"github.com/gruntwork-io/go-commons/errors"
	"github.com/stretchr/testify/require"

	"github.com/gruntwork-io/terratest/modules/files"
	"github.com/gruntwork-io/terratest/modules/testing"
)

// Install will install the selected helm chart with the provided options under the given release name. This will fail
// the test if there is an error.
//
// Deprecated: Use [InstallContext] instead.
func Install(t testing.TestingT, options *Options, chart string, releaseName string) {
	InstallContext(t, context.Background(), options, chart, releaseName)
}

// InstallContext will install the selected helm chart with the provided options under the given release name. This will
// fail the test if there is an error. The ctx parameter supports cancellation and timeouts.
func InstallContext(t testing.TestingT, ctx context.Context, options *Options, chart string, releaseName string) {
	require.NoError(t, InstallContextE(t, ctx, options, chart, releaseName))
}

// InstallE will install the selected helm chart with the provided options under the given release name.
//
// Deprecated: Use [InstallContextE] instead.
func InstallE(t testing.TestingT, options *Options, chart string, releaseName string) error {
	return InstallContextE(t, context.Background(), options, chart, releaseName)
}

// InstallContextE will install the selected helm chart with the provided options under the given release name. If
// releaseName is empty, the name is omitted so callers can rely on helm's --generate-name (passed through ExtraArgs).
// The ctx parameter supports cancellation and timeouts.
func InstallContextE(t testing.TestingT, ctx context.Context, options *Options, chart string, releaseName string) error {
	// If the chart refers to a path, convert to absolute path. Otherwise, pass straight through as it may be a remote
	// chart.
	if files.FileExists(chart) {
		absChartDir, err := filepath.Abs(chart)
		if err != nil {
			return errors.WithStackTrace(err)
		}

		chart = absChartDir
	}

	// build chart dependencies
	if options.BuildDependencies {
		if _, err := RunHelmCommandAndGetOutputContextE(t, ctx, options, "dependency", "build", chart); err != nil {
			return errors.WithStackTrace(err)
		}
	}

	// Now call out to helm install to install the charts with the provided options
	args, err := installArgs(options, chart, releaseName) //nolint:contextcheck // installArgs only formats values via context-free helpers
	if err != nil {
		return err
	}

	_, err = RunHelmCommandAndGetOutputContextE(t, ctx, options, "install", args...)

	return err
}

// installArgs builds the argument list passed to `helm install`. When releaseName is empty, the release name is omitted
// so callers can rely on helm's --generate-name (supplied through ExtraArgs["install"]) to name the release. Passing an
// empty release name as a positional argument would otherwise make helm fail with "expected at most two arguments".
func installArgs(options *Options, chart string, releaseName string) ([]string, error) {
	args := []string{}

	if options.ExtraArgs != nil {
		if extra, ok := options.ExtraArgs["install"]; ok {
			args = append(args, extra...)
		}
	}

	if options.Version != "" {
		args = append(args, "--version", options.Version)
	}

	args, err := getValuesArgsE(options, args...)
	if err != nil {
		return nil, err
	}

	if releaseName != "" {
		args = append(args, releaseName)
	}

	args = append(args, chart)

	return args, nil
}
