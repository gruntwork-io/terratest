package terraform

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/gruntwork-io/terratest/modules/collections"
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
		terraformArgs = append(terraformArgs, FormatTerraformArgs("-var-file", options.VarFiles)...)
		terraformArgs = append(terraformArgs, FormatTerraformVarsAsArgs(options.Vars)...)
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
	return formatTerraformArgs(vars, "-var", true)
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
	pluginArgs := []string{fmt.Sprintf("-plugin-dir=%v", pluginDir)}
	if pluginDir == "" {
		return nil
	}
	return pluginArgs
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
	return formatTerraformArgs(vars, "-backend-config", false)
}

// Format the given vars into 'Terraform' format, with each var being prefixed with the given prefix. If
// useSpaceAsSeparator is true, a space will separate the prefix and each var (e.g., -var foo=bar). If
// useSpaceAsSeparator is false, an equals will separate the prefix and each var (e.g., -backend-config=foo=bar).
func formatTerraformArgs(vars map[string]interface{}, prefix string, useSpaceAsSeparator bool) []string {
	var args []string

	for key, value := range vars {
		hclString := toHclString(value, false)
		argValue := fmt.Sprintf("%s=%s", key, hclString)
		if useSpaceAsSeparator {
			args = append(args, prefix, argValue)
		} else {
			args = append(args, fmt.Sprintf("%s=%s", prefix, argValue))
		}
	}

	return args
}

// Terraform allows you to pass in command-line variables using HCL syntax (e.g. -var foo=[1,2,3]). Unfortunately,
// while their golang hcl library can convert an HCL string to a Go type, they don't seem to offer a library to convert
// arbitrary Go types to an HCL string. Therefore, this method is a simple implementation that correctly handles
// ints, booleans, lists, and maps. Everything else is forced into a string using Sprintf. Hopefully, this approach is
// good enough for the type of variables we deal with in Terratest.
func toHclString(value interface{}, isNested bool) string {
	// Ideally, we'd use a type switch here to identify slices and maps, but we can't do that, because Go doesn't
	// support generics, and the type switch only matches concrete types. So we could match []interface{}, but if
	// a user passes in []string{}, that would NOT match (the same logic applies to maps). Therefore, we have to
	// use reflection and manually convert into []interface{} and map[string]interface{}.

	if slice, isSlice := tryToConvertToGenericSlice(value); isSlice {
		return sliceToHclString(slice)
	} else if m, isMap := tryToConvertToGenericMap(value); isMap {
		return mapToHclString(m)
	} else {
		return primitiveToHclString(value, isNested)
	}
}

// Try to convert the given value to a generic slice. Return the slice and true if the underlying value itself was a
// slice and an empty slice and false if it wasn't. This is necessary because Go is a shitty language that doesn't
// have generics, nor useful utility methods built-in. For more info, see: http://stackoverflow.com/a/12754757/483528
func tryToConvertToGenericSlice(value interface{}) ([]interface{}, bool) {
	reflectValue := reflect.ValueOf(value)
	if reflectValue.Kind() != reflect.Slice {
		return []interface{}{}, false
	}

	genericSlice := make([]interface{}, reflectValue.Len())

	for i := 0; i < reflectValue.Len(); i++ {
		genericSlice[i] = reflectValue.Index(i).Interface()
	}

	return genericSlice, true
}

// Try to convert the given value to a generic map. Return the map and true if the underlying value itself was a
// map and an empty map and false if it wasn't. This is necessary because Go is a shitty language that doesn't
// have generics, nor useful utility methods built-in. For more info, see: http://stackoverflow.com/a/12754757/483528
func tryToConvertToGenericMap(value interface{}) (map[string]interface{}, bool) {
	reflectValue := reflect.ValueOf(value)
	if reflectValue.Kind() != reflect.Map {
		return map[string]interface{}{}, false
	}

	reflectType := reflect.TypeOf(value)
	if reflectType.Key().Kind() != reflect.String {
		return map[string]interface{}{}, false
	}

	genericMap := make(map[string]interface{}, reflectValue.Len())

	mapKeys := reflectValue.MapKeys()
	for _, key := range mapKeys {
		genericMap[key.String()] = reflectValue.MapIndex(key).Interface()
	}

	return genericMap, true
}

// Convert a slice to an HCL string. See ToHclString for details.
func sliceToHclString(slice []interface{}) string {
	hclValues := []string{}

	for _, value := range slice {
		hclValue := toHclString(value, true)
		hclValues = append(hclValues, hclValue)
	}

	return fmt.Sprintf("[%s]", strings.Join(hclValues, ", "))
}

// Convert a map to an HCL string. See ToHclString for details.
func mapToHclString(m map[string]interface{}) string {
	keyValuePairs := []string{}

	for key, value := range m {
		keyValuePair := fmt.Sprintf(`"%s" = %s`, key, toHclString(value, true))
		keyValuePairs = append(keyValuePairs, keyValuePair)
	}

	return fmt.Sprintf("{%s}", strings.Join(keyValuePairs, ", "))
}

// Convert a primitive, such as a bool, int, or string, to an HCL string. If this isn't a primitive, force its value
// using Sprintf. See ToHclString for details.
func primitiveToHclString(value interface{}, isNested bool) string {
	if value == nil {
		return "null"
	}

	switch v := value.(type) {

	case bool:
		return strconv.FormatBool(v)

	case string:
		// If string is nested in a larger data structure (e.g. list of string, map of string), ensure value is quoted
		if isNested {
			return fmt.Sprintf("\"%v\"", v)
		}

		return fmt.Sprintf("%v", v)

	default:
		return fmt.Sprintf("%v", v)
	}
}
