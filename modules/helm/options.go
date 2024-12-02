package helm

import (
	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/gruntwork-io/terratest/modules/logger"
)

type Options struct {
	ValuesFiles       []string            // List of values files to render.
	SetValues         map[string]string   // Values that should be set via the command line.
	SetStrValues      map[string]string   // Values that should be set via the command line explicitly as `string` types.
	SetJsonValues     map[string]string   // Values that should be set via the command line in JSON format.
	SetFiles          map[string]string   // Values that should be set from a file. These should be file paths. Use to avoid logging secrets.
	KubectlOptions    *k8s.KubectlOptions // KubectlOptions to control how to authenticate to kubernetes cluster. `nil` => use defaults.
	HomePath          string              // The path to the helm home to use when calling out to helm. Empty string means use default ($HOME/.helm).
	EnvVars           map[string]string   // Environment variables to set when running helm
	Version           string              // Version of chart
	Logger            *logger.Logger      // Set a non-default logger that should be used. See the logger package for more info. Use logger.Discard to not print the output while executing the command.
	ExtraArgs         map[string][]string // Extra arguments to pass to the helm install/upgrade/rollback/delete and helm repo add commands. The key signals the command (e.g., install) while the values are the extra arguments to pass through.
	BuildDependencies bool                // If true, helm dependencies will be built before rendering template, installing or upgrade the chart.
	SnapshotPath      string              // The path to the snapshot directory when using snapshot based testing. Empty string means use default ($PWD/__snapshot__).
}
