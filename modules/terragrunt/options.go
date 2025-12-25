package terragrunt

import (
	"io"
	"os"
	"time"

	"github.com/gruntwork-io/terratest/modules/logger/v2"
)

// Key concepts:
// - Options: Configure HOW the test framework executes tg (directories, retry logic, logging)
// - TerragruntArgs: Global terragrunt flags (e.g., --log-level, --no-color)
// - TerraformArgs: Command-specific terraform args (e.g., -upgrade for init, or the command itself for stack run)
// - Use Options.TerragruntDir to specify WHERE to run tg
//
// Example:
//
//	// For init with terraform flags
//	InitE(t, &Options{
//	    TerragruntDir: "/path/to/config",
//	    TerragruntArgs: []string{"--log-level", "info"},
//	    TerraformArgs: []string{"-upgrade=true"},
//	})
//
//	// For run-all with global flags
//	ApplyAllE(t, &Options{
//	    TerragruntDir: "/path/to/config",
//	    TerragruntArgs: []string{"--no-color"},
//	})
//
// Constants for test framework configuration and environment variables
const (
	DefaultTerragruntBinary = "terragrunt"
	NonInteractiveFlag      = "--non-interactive"
	TerragruntLogFormatKey  = "TG_LOG_FORMAT"
	TerragruntLogCustomKey  = "TG_LOG_CUSTOM_FORMAT"
	DefaultLogFormat        = "key-value"
	DefaultLogCustomFormat  = "%msg(color=disable)"
)

// Options represent the configuration options for tg test execution.
//
// This struct is divided into two clear categories:
//
// 1. TEST FRAMEWORK CONFIGURATION:
//   - Controls HOW the test framework executes tg
//   - Includes: binary paths, directories, retry logic, logging, environment
//   - These are NOT passed as command-line arguments to tg
//
// 2. TG COMMAND ARGUMENTS:
//   - TerragruntArgs: Global terragrunt flags (placed BEFORE the command)
//   - TerraformArgs: Command-specific flags (placed AFTER the command)
//   - These ARE passed directly to tg in the appropriate positions
//
// This separation eliminates confusion about which settings control the test
// framework vs which become tg command-line arguments.
type Options struct {
	// Test framework configuration (NOT passed to tg command line)
	TerragruntBinary string            // The tg binary to use (should be "terragrunt")
	TerragruntDir    string            // The directory containing the tg configuration
	EnvVars          map[string]string // Environment variables for command execution
	Logger           *logger.Logger    // Logger for command output

	// Test framework retry and error handling (NOT passed to tg command line)
	MaxRetries               int               // Maximum number of retries
	TimeBetweenRetries       time.Duration     // Time between retries
	RetryableTerraformErrors map[string]string // Retryable error patterns
	WarningsAsErrors         map[string]string // Warnings to treat as errors

	// Complex configuration that requires special formatting (NOT raw command-line args)
	BackendConfig map[string]interface{} // Backend configuration (formatted specially)
	PluginDir     string                 // Plugin directory (formatted specially)

	// Global terragrunt command-line flags (placed BEFORE the command)
	// Example: []string{"--log-level", "info", "--no-color"}
	TerragruntArgs []string

	// Command-specific terraform flags (placed AFTER the command)
	// Example: []string{"-upgrade=true"} for init, or []string{"plan"} for stack run
	TerraformArgs []string

	// Optional stdin to pass to Terraform commands
	Stdin io.Reader
}

// setTerragruntLogFormatting sets default log formatting for tg
// if it is not already set in options.EnvVars or OS environment vars
func setTerragruntLogFormatting(options *Options) {
	if options.EnvVars == nil {
		options.EnvVars = make(map[string]string)
	}

	_, inOpts := options.EnvVars[TerragruntLogFormatKey]
	if !inOpts {
		_, inEnv := os.LookupEnv(TerragruntLogFormatKey)
		if !inEnv {
			// key-value format for tg logs to avoid colors and have plain form
			// https://terragrunt.gruntwork.io/docs/reference/cli-options/#terragrunt-log-format
			options.EnvVars[TerragruntLogFormatKey] = DefaultLogFormat
		}
	}

	_, inOpts = options.EnvVars[TerragruntLogCustomKey]
	if !inOpts {
		_, inEnv := os.LookupEnv(TerragruntLogCustomKey)
		if !inEnv {
			options.EnvVars[TerragruntLogCustomKey] = DefaultLogCustomFormat
		}
	}
}
