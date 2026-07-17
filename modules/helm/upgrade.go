package helm

import (
	"context"
	"path/filepath"

	"github.com/gruntwork-io/go-commons/errors"
	"github.com/gruntwork-io/terratest/modules/core/v2/files"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
)

// UpgradeContext will upgrade the release and chart will be deployed with the latest configuration. This will fail
// the test if there is an error. The ctx parameter supports cancellation and timeouts.
func UpgradeContext(t testing.TestingT, ctx context.Context, options *Options, chart string, releaseName string) {
	require.NoError(t, UpgradeContextE(t, ctx, options, chart, releaseName))
}

// UpgradeContextE will upgrade the release and chart will be deployed with the latest configuration. The ctx
// parameter supports cancellation and timeouts.
func UpgradeContextE(t testing.TestingT, ctx context.Context, options *Options, chart string, releaseName string) error {

	if files.FileExists(chart) {
		absChartDir, err := filepath.Abs(chart)
		if err != nil {
			return errors.WithStackTrace(err)
		}

		chart = absChartDir
	}

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

	args, err = getValuesArgsE(options, args...)
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
