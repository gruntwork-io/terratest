package helm

import (
	"context"
	"path/filepath"

	"github.com/gruntwork-io/go-commons/errors"
	"github.com/gruntwork-io/terratest/modules/files"
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// Upgrade will upgrade the release and chart will be deployed with the latest configuration. This will fail
// the test if there is an error.
//
// Deprecated: Use [UpgradeContext] instead.
func Upgrade(t testing.TestingT, options *Options, chart string, releaseName string) {
	UpgradeContext(t, context.Background(), options, chart, releaseName)
}

// UpgradeContext will upgrade the release and chart will be deployed with the latest configuration. This will fail
// the test if there is an error. The ctx parameter supports cancellation and timeouts.
func UpgradeContext(t testing.TestingT, ctx context.Context, options *Options, chart string, releaseName string) {
	require.NoError(t, UpgradeContextE(t, ctx, options, chart, releaseName))
}

// UpgradeE will upgrade the release and chart will be deployed with the latest configuration.
//
// Deprecated: Use [UpgradeContextE] instead.
func UpgradeE(t testing.TestingT, options *Options, chart string, releaseName string) error {
	return UpgradeContextE(t, context.Background(), options, chart, releaseName)
}

// UpgradeContextE will upgrade the release and chart will be deployed with the latest configuration. The ctx
// parameter supports cancellation and timeouts.
func UpgradeContextE(t testing.TestingT, ctx context.Context, options *Options, chart string, releaseName string) error {
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

	var err error

	args := []string{}

	if options.ExtraArgs != nil {
		if upgradeArgs, ok := options.ExtraArgs["upgrade"]; ok {
			args = append(args, upgradeArgs...)
		}
	}

	args, err = getValuesArgsE(options, args...) //nolint:contextcheck // getValuesArgsE is a local helper without context
	if err != nil {
		return err
	}

	args = append(args, "--install", releaseName, chart)

	if options.Version != "" {
		args = append(args, "--version", options.Version)
	}

	_, err = RunHelmCommandAndGetOutputContextE(t, ctx, options, "upgrade", args...)

	return err
}
