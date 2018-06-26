package terragrunt

import "time"

// Options for running Terragrunt commands
type Options struct {
	TerragruntDir             string                 // The path to the folder where the Terragrunt code is defined.
	Vars                      map[string]interface{} // The vars to pass to Terragrunt commands using the -var option.
	EnvVars                   map[string]string      // Environment variables to set when running Terragrunt
	BackendConfig             map[string]interface{} // The vars to pass to the terragrunt init command for extra configuration for the backend
	RetryableTerragruntErrors map[string]string      // If Terragrunt apply fails with one of these (transient) errors, retry. The keys are text to look for in the error and the message is what to display to a user if that error is found.
	MaxRetries                int                    // Maximum number of times to retry errors matching RetryableTerragruntErrors
	TimeBetweenRetries        time.Duration          // The amount of time to wait between retries
	Upgrade                   bool                   // Whether the -upgrade flag of the terragrunt init command should be set to true or not
	NoColor                   bool                   // Whether the -no-color flag will be set for any Terragrunt command or not
}
