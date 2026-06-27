package helm

import (
	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/gruntwork-io/terratest/modules/core/v2/logger"
)

// Options represents the options for a Helm command.
type Options struct {
	// Set a non-default logger that should be used. See the logger package for more info.
	// Use logger.Discard to not print the output while executing the command.
	Logger *logger.Logger

	// Values that should be set via the command line.
	SetValues map[string]string

	// Values that should be set via the command line explicitly as `string` types.
	SetStrValues map[string]string

	// SetJSONValues are values that should be set via the command line in JSON format.
	SetJSONValues map[string]string

	// Deprecated: Use [SetJSONValues] instead.
	SetJsonValues map[string]string //nolint:revive,staticcheck // Deprecated field kept for backwards compatibility.

	// Values that should be set from a file. These should be file paths. Use to avoid logging secrets.
	SetFiles map[string]string

	// KubectlOptions to control how to authenticate to kubernetes cluster. `nil` => use defaults.
	KubectlOptions *k8s.KubectlOptions

	// Environment variables to set when running helm.
	EnvVars map[string]string

	// Extra arguments to pass to the helm install/upgrade/rollback/delete and helm repo add commands.
	// The key signals the command (e.g., install) while the values are the extra arguments to pass through.
	ExtraArgs map[string][]string

	// The path to the helm home to use when calling out to helm.
	// Empty string means use default ($HOME/.helm).
	HomePath string

	// Version of chart.
	Version string

	// The path to the snapshot directory when using snapshot based testing.
	// Empty string means use default ($PWD/__snapshot__).
	SnapshotPath string

	// List of values files to render.
	ValuesFiles []string

	// If true, helm dependencies will be built before rendering template, installing or upgrading the chart.
	BuildDependencies bool
}
