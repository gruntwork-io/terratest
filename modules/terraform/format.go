package terraform

import (
	"fmt"

	"github.com/gruntwork-io/terratest/modules/collections"
	"github.com/gruntwork-io/terratest/modules/formatting"
)

const runAllCmd = "run-all"

// TerraformCommandsWithLockSupport is a list of all the Terraform commands that
// can obtain locks on Terraform state
var TerraformCommandsWithLockSupport = []string{
	"plan",
	"plan-all",
	"apply",
	"apply-all",
	"destroy",
	"destroy-all",
	"init",
	"refresh",
	"taint",
	"untaint",
	"import",
}

// TerraformCommandsWithPlanFileSupport is a list of all the Terraform commands that support interacting with plan
// files.
var TerraformCommandsWithPlanFileSupport = []string{
	"plan",
	"apply",
	"show",
	"graph",
}

// FormatArgs converts the inputs to a format palatable to terraform. This includes converting the given vars to the
// format the Terraform CLI expects (-var key=value).
func FormatArgs(options *Options, args ...string) []string {
	var terraformArgs []string
	commandType := args[0]
	// If the user is trying to run with run-all, then we need to make sure the command based args are based on the
	// actual terraform command. E.g., we want to base the logic on `plan` when `run-all plan` is passed in, not
	// `run-all`.
	if commandType == runAllCmd {
		commandType = args[1]
	}
	lockSupported := collections.ListContains(TerraformCommandsWithLockSupport, commandType)
	planFileSupported := collections.ListContains(TerraformCommandsWithPlanFileSupport, commandType)

	// Include -var and -var-file flags unless we're running 'apply' with a plan file
	includeVars := !(commandType == "apply" && len(options.PlanFilePath) > 0)

	terraformArgs = append(terraformArgs, args...)

	if includeVars {
		for _, v := range options.MixedVars {
			terraformArgs = append(terraformArgs, v.Args()...)
		}

		if options.SetVarsAfterVarFiles {
			terraformArgs = append(terraformArgs, FormatTerraformArgs("-var-file", options.VarFiles)...)
			terraformArgs = append(terraformArgs, FormatTerraformVarsAsArgs(options.Vars)...)
		} else {
			terraformArgs = append(terraformArgs, FormatTerraformVarsAsArgs(options.Vars)...)
			terraformArgs = append(terraformArgs, FormatTerraformArgs("-var-file", options.VarFiles)...)
		}
	}

	terraformArgs = append(terraformArgs, FormatTerraformArgs("-target", options.Targets)...)

	if options.NoColor {
		terraformArgs = append(terraformArgs, "-no-color")
	}

	if lockSupported {
		// If command supports locking, handle lock arguments
		terraformArgs = append(terraformArgs, FormatTerraformLockAsArgs(options.Lock, options.LockTimeout)...)
	}

	if planFileSupported {
		// The plan file arg should be last in the terraformArgs slice. Some commands use it as an input (e.g. show, apply)
		terraformArgs = append(terraformArgs, FormatTerraformPlanFileAsArg(commandType, options.PlanFilePath)...)
	}

	return terraformArgs
}

// FormatTerraformPlanFileAsArg formats the out variable as a command-line arg for Terraform (e.g. of the format
// -out=/some/path/to/plan.out or /some/path/to/plan.out). Only plan supports passing in the plan file as -out; the
// other commands expect it as the first positional argument. This returns an empty string if outPath is empty string.
func FormatTerraformPlanFileAsArg(commandType string, outPath string) []string {
	if outPath == "" {
		return nil
	}
	if commandType == "plan" {
		return []string{fmt.Sprintf("%s=%s", "-out", outPath)}
	}
	return []string{outPath}
}

// FormatTerraformVarsAsArgs formats the given variables as command-line args for Terraform (e.g. of the format
// -var key=value).
func FormatTerraformVarsAsArgs(vars map[string]interface{}) []string {
	return formatting.FormatTerraformArgs(vars, "-var", true, false)
}

// FormatTerraformLockAsArgs formats the lock and lock-timeout variables
// -lock, -lock-timeout
func FormatTerraformLockAsArgs(lockCheck bool, lockTimeout string) []string {
	lockArgs := []string{fmt.Sprintf("-lock=%v", lockCheck)}
	if lockTimeout != "" {
		lockTimeoutValue := fmt.Sprintf("%s=%s", "-lock-timeout", lockTimeout)
		lockArgs = append(lockArgs, lockTimeoutValue)
	}
	return lockArgs
}

// FormatTerraformPluginDirAsArgs formats the plugin-dir variable
// -plugin-dir
func FormatTerraformPluginDirAsArgs(pluginDir string) []string {
	return formatting.FormatPluginDirAsArgs(pluginDir)
}

// FormatTerraformArgs will format multiple args with the arg name (e.g. "-var-file", []string{"foo.tfvars", "bar.tfvars", "baz.tfvars.json"})
// returns "-var-file foo.tfvars -var-file bar.tfvars -var-file baz.tfvars.json"
func FormatTerraformArgs(argName string, args []string) []string {
	argsList := []string{}
	for _, argValue := range args {
		argsList = append(argsList, argName, argValue)
	}
	return argsList
}

// FormatTerraformBackendConfigAsArgs formats the given variables as backend config args for Terraform (e.g. of the
// format -backend-config=key=value).
func FormatTerraformBackendConfigAsArgs(vars map[string]interface{}) []string {
	return formatting.FormatBackendConfigAsArgs(vars)
}
