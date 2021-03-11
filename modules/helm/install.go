package helm

import (
	"path/filepath"

	"github.com/gruntwork-io/go-commons/errors"
	"github.com/stretchr/testify/require"

	"github.com/gruntwork-io/terratest/modules/files"
	"github.com/gruntwork-io/terratest/modules/testing"
)

// Install will install the selected helm chart with the provided options under the given release name. This will fail
// the test if there is an error.
func Install(t testing.TestingT, options *Options, chart string, releaseName string) {
	require.NoError(t, InstallE(t, options, chart, releaseName))
}

// InstallE will install the selected helm chart with the provided options under the given release name.
func InstallE(t testing.TestingT, options *Options, chart string, releaseName string) error {
	// If the chart refers to a path, convert to absolute path. Otherwise, pass straight through as it may be a remote
	// chart.
	if files.FileExists(chart) {
		absChartDir, err := filepath.Abs(chart)
		if err != nil {
			return errors.WithStackTrace(err)
		}
		chart = absChartDir
	}

	// Now call out to helm install to install the charts with the provided options
	// Declare err here so that we can update args later
	var err error
	args := []string{}
	if options.ExtraArgs != nil {
		if installArgs, ok := options.ExtraArgs["install"]; ok {
			args = append(args, installArgs...)
		}
	}
	args = append(args, getNamespaceArgs(options)...)
	if options.Version != "" {
		args = append(args, "--version", options.Version)
	}
	args, err = getValuesArgsE(t, options, args...)
	if err != nil {
		return err
	}
	args = append(args, releaseName, chart)
	_, err = RunHelmCommandAndGetOutputE(t, options, "install", args...)
	return err
}
