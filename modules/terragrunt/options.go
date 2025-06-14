package terragrunt

import (
	"os"
	"time"

	"github.com/gruntwork-io/terratest/modules/logger"
)

// Options represents the configuration options for terragrunt commands
type Options struct {
	// Terragrunt-specific options
	TerragruntBinary string // The terragrunt binary to use (should be "terragrunt")
	TerragruntDir    string // The directory containing the terragrunt configuration

	// Command-specific options
	NoColor      bool // Whether to disable colored output
	Upgrade      bool // Whether to upgrade modules and plugins (init command)
	Reconfigure  bool // Whether to reconfigure the backend (init command)
	MigrateState bool // Whether to migrate state (init command)

	// Configuration options
	BackendConfig map[string]interface{} // Backend configuration (init command)
	PluginDir     string                 // Plugin directory (init command)
	EnvVars       map[string]string      // Environment variables for command execution

	// Execution options
	Logger                   *logger.Logger    // Logger for command output
	MaxRetries               int               // Maximum number of retries
	TimeBetweenRetries       time.Duration     // Time between retries
	RetryableTerraformErrors map[string]string // Retryable error patterns
	WarningsAsErrors         map[string]string // Warnings to treat as errors

	// Terragrunt-specific extra arguments
	ExtraArgs ExtraArgs
}

// ExtraArgs represents terragrunt-specific extra arguments for different commands
type ExtraArgs struct {
	Init     []string // Extra arguments for init command
	Apply    []string // Extra arguments for apply command (used by generate)
	Plan     []string // Extra arguments for plan command
	Destroy  []string // Extra arguments for destroy command
	Generate []string // Extra arguments for generate command
}

// GetCommonOptions extracts common terragrunt options and prepares arguments
// This is the terragrunt-specific version of terraform.GetCommonOptions
func GetCommonOptions(options *Options, args ...string) (*Options, []string) {
	// Set default binary if not specified
	if options.TerragruntBinary == "" {
		options.TerragruntBinary = "terragrunt"
	}

	// Add terragrunt-specific flags
	args = append(args, "--non-interactive")

	// Set terragrunt log formatting if not already set
	setTerragruntLogFormatting(options)

	return options, args
}

// setTerragruntLogFormatting sets default log formatting for terragrunt
// if it is not already set in options.EnvVars or OS environment vars
func setTerragruntLogFormatting(options *Options) {
	const (
		tgLogFormatKey       = "TG_LOG_FORMAT"
		tgLogCustomFormatKey = "TG_LOG_CUSTOM_FORMAT"
	)

	if options.EnvVars == nil {
		options.EnvVars = make(map[string]string)
	}

	_, inOpts := options.EnvVars[tgLogFormatKey]
	if !inOpts {
		_, inEnv := os.LookupEnv(tgLogFormatKey)
		if !inEnv {
			// key-value format for terragrunt logs to avoid colors and have plain form
			// https://terragrunt.gruntwork.io/docs/reference/cli-options/#terragrunt-log-format
			options.EnvVars[tgLogFormatKey] = "key-value"
		}
	}

	_, inOpts = options.EnvVars[tgLogCustomFormatKey]
	if !inOpts {
		_, inEnv := os.LookupEnv(tgLogCustomFormatKey)
		if !inEnv {
			options.EnvVars[tgLogCustomFormatKey] = "%msg(color=disable)"
		}
	}
}
