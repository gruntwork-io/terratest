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

// InstallContextE will install the selected helm chart with the provided options under the given release name. The ctx
// parameter supports cancellation and timeouts.
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
	var err error

	args := []string{}

	if options.ExtraArgs != nil {
		if installArgs, ok := options.ExtraArgs["install"]; ok {
			args = append(args, installArgs...)
		}
	}

	if options.Version != "" {
		args = append(args, "--version", options.Version)
	}

	args, err = getValuesArgsE(options, args...) //nolint:contextcheck // getValuesArgsE is a local helper without context
	if err != nil {
		return err
	}

	args = append(args, releaseName, chart)
	_, err = RunHelmCommandAndGetOutputContextE(t, ctx, options, "install", args...)

	return err
}
