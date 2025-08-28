package terragrunt

import (
	"io"
	"os"
	"time"

	"github.com/gruntwork-io/terratest/modules/logger"
)

// Key concepts:
// - Options: Configure HOW the test framework executes tg (directories, retry logic, logging)
// - TerragruntArgs: Arguments for tg itself (e.g., --no-color for tg output)
// - TerraformArgs: Arguments passed to underlying terraform commands after -- separator
// - Use Options.TerragruntDir to specify WHERE to run tg
//
// Example:
//
//	// For init with terraform-specific flags
//	TgInitE(t, &Options{
//	    TerragruntDir: "/path/to/config",
//	    TerragruntArgs: []string{"--no-color"},
//	    TerraformArgs: []string{"-upgrade=true"},
//	})
//
//	// For stack run with terraform plan
//	TgStackRunE(t, &Options{
//	    TerragruntDir: "/path/to/config",
//	    TerragruntArgs: []string{"--no-color"},
//	    TerraformArgs: []string{"plan", "-out=tfplan"},
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
	ArgSeparator            = "--"
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
//   - All actual tg command-line arguments go in ExtraArgs []string
//   - This includes flags like -no-color, -upgrade, -reconfigure, etc.
//   - These ARE passed directly to the specific tg command being executed
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

	// Tg-specific command-line arguments (e.g., --no-color for tg itself)
	TerragruntArgs []string

	// Terraform command-line arguments to be passed after -- separator
	// These are passed directly to the underlying terraform commands
	TerraformArgs []string

	// Optional stdin to pass to Terraform commands
	Stdin io.Reader

	// The vars to pass to Terraform commands using the -var option. Note that terraform does not support passing `null`
	// as a variable value through the command line. That is, if you use `map[string]interface{}{"foo": nil}` as `Vars`,
	// this will translate to the string literal `"null"` being assigned to the variable `foo`. However, nulls in
	// lists and maps/objects are supported. E.g., the following var will be set as expected (`{ bar = null }`:
	// map[string]interface{}{
	//     "foo": map[string]interface{}{"bar": nil},
	// }
	Vars                 map[string]interface{}
	VarFiles             []string  // The var file paths to pass to Terraform commands using -var-file option.
	SetVarsAfterVarFiles bool      // Pass -var options after -var-file options to Terraform commands
	MixedVars            []Var     // Mix of `-var` and `-var-file` in arbritrary order, use `VarInline()` `VarFile()` to set the value.
	PlanFilePath         string    // The path to output a plan file to (for the plan command) or read one from (for the apply command)
	Targets              []string  // The target resources to pass to the terraform command with -target
	Lock                 bool      // The lock option to pass to the terraform command with -lock
	LockTimeout          string    // The lock timeout option to pass to the terraform command with -lock-timeout
	NoColor              bool      // Whether the -no-color flag will be set for any Terraform command or not
	ExtraArgs            ExtraArgs // Extra arguments passed to Terraform commands
}

type ExtraArgs struct {
	Apply   []string
	Destroy []string
	Plan    []string
}

func prepend(args []string, arg ...string) []string {
	return append(arg, args...)
}

// GetCommonOptions extracts common tg options and prepares arguments
// This is the tg-specific version of terraform.GetCommonOptions
func GetCommonOptions(options *Options, args ...string) (*Options, []string) {
	// Set default binary if not specified
	if options.TerragruntBinary == "" {
		options.TerragruntBinary = DefaultTerragruntBinary
	}

	// Add tg-specific flags
	args = append(args, NonInteractiveFlag)

	// Set tg log formatting if not already set
	setTerragruntLogFormatting(options)

	return options, args
}

// GetArgsForCommand returns the appropriate arguments based on the command type
// It handles the separation of tg and terraform arguments
func GetArgsForCommand(options *Options, useArgSeparator bool) []string {
	var args []string

	// First add tg-specific arguments
	args = append(args, options.TerragruntArgs...)

	// Then add terraform arguments with separator if needed
	if len(options.TerraformArgs) > 0 {
		if useArgSeparator {
			args = append(args, ArgSeparator)
		}
		args = append(args, options.TerraformArgs...)
	}

	return args
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
