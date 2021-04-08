package helm

import (
	"path/filepath"

	"github.com/gruntwork-io/go-commons/errors"
	"github.com/gruntwork-io/terratest/modules/files"
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// Upgrade will upgrade the release and chart will be deployed with the lastest configuration. This will fail
// the test if there is an error.
func Upgrade(t testing.TestingT, options *Options, chart string, releaseName string) {
	require.NoError(t, UpgradeE(t, options, chart, releaseName))
}

// UpgradeE will upgrade the release and chart will be deployed with the lastest configuration.
func UpgradeE(t testing.TestingT, options *Options, chart string, releaseName string) error {
	// If the chart refers to a path, convert to absolute path. Otherwise, pass straight through as it may be a remote
	// chart.
	if files.FileExists(chart) {
		absChartDir, err := filepath.Abs(chart)
		if err != nil {
			return errors.WithStackTrace(err)
		}
		chart = absChartDir
	}

	var err error
	args := []string{}
	if options.ExtraArgs != nil {
		if upgradeArgs, ok := options.ExtraArgs["upgrade"]; ok {
			args = append(args, upgradeArgs...)
		}
	}
	args = append(args, getNamespaceArgs(options)...)
	args, err = getValuesArgsE(t, options, args...)
	if err != nil {
		return err
	}

	args = append(args, "--install", releaseName, chart)
	_, err = RunHelmCommandAndGetOutputE(t, options, "upgrade", args...)
	return err
}
