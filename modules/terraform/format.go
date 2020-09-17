package terraform

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/gruntwork-io/terratest/modules/collections"
)

// TerraformDefaultLockingStatus - The terratest default command lock status (backwards compatibility)
var TerraformCommandsWithLockSupport = []string{
	"plan",
	"apply",
	"destroy",
	"init",
	"refresh",
	"taint",
	"untaint",
	"import",
}

// FormatArgs converts the inputs to a format palatable to terraform. This includes converting the given vars to the
// format the Terraform CLI expects (-var key=value).
func FormatArgs(options *Options, args ...string) []string {
	var terraformArgs []string
	commandType := args[0]
	lockSupported := collections.ListContains(TerraformCommandsWithLockSupport, commandType)

	terraformArgs = append(terraformArgs, args...)
	terraformArgs = append(terraformArgs, FormatTerraformVarsAsArgs(options.Vars)...)
	terraformArgs = append(terraformArgs, FormatTerraformArgs("-var-file", options.VarFiles)...)
	terraformArgs = append(terraformArgs, FormatTerraformArgs("-target", options.Targets)...)

	if lockSupported {
		// If command supports locking, handle lock arguments
		terraformArgs = append(terraformArgs, FormatTerraformLockAsArgs(options.Lock, options.LockTimeout)...)
	}

	return terraformArgs
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

// FormatTerraformArgs will format multiple args with the arg name (e.g. "-var-file", []string{"foo.tfvars", "bar.tfvars"})
// returns "-var-file foo.tfvars -var-file bar.tfvars"
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
