package formatting

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// FormatPluginDirAsArgs formats the plugin-dir variable
// -plugin-dir
func FormatPluginDirAsArgs(pluginDir string) []string {
	pluginArgs := []string{fmt.Sprintf("-plugin-dir=%v", pluginDir)}
	if pluginDir == "" {
		return nil
	}
	return pluginArgs
}

// FormatBackendConfigAsArgs formats the given variables as backend config args (e.g. of the
// format -backend-config=key=value). Used by both Terraform and Terragrunt.
func FormatBackendConfigAsArgs(vars map[string]interface{}) []string {
	return FormatTerraformArgs(vars, "-backend-config", false, true)
}

// FormatTerraformArgs formats the given vars into HCL command-line argument format, with each var being
// prefixed with the given prefix. This is used by both Terraform and Terragrunt commands.
// If useSpaceAsSeparator is true, a space will separate the prefix and each var (e.g., -var foo=bar).
// If useSpaceAsSeparator is false, an equals will separate the prefix and each var (e.g., -backend-config=foo=bar).
// If omitNil is false, then nil values will be included as "null" (e.g. -backend-config=foo=null).
// If omitNil is true, then nil values will be formatted without the value (e.g. -backend-config=foo).
func FormatTerraformArgs(vars map[string]interface{}, prefix string, useSpaceAsSeparator bool, omitNil bool) []string {
	var args []string

	for key, value := range vars {
		var argValue string
		if omitNil && value == nil {
			argValue = key
		} else {
			hclString := toHclString(value, false)
			argValue = fmt.Sprintf("%s=%s", key, hclString)
		}
		if useSpaceAsSeparator {
			args = append(args, prefix, argValue)
		} else {
			args = append(args, fmt.Sprintf("%s=%s", prefix, argValue))
		}
	}

	return args
}

// toHclString converts Go types to HCL syntax strings (e.g. []int{1,2,3} becomes "[1, 2, 3]").
// Terraform and Terragrunt accept command-line variables using HCL syntax (e.g. -var foo=[1,2,3]).
// Unfortunately, while the HCL library can parse HCL strings to Go types, it doesn't offer a way to
// convert arbitrary Go types back to HCL strings. This is a simple implementation that correctly
// handles ints, booleans, lists, and maps. Everything else is formatted as a string using Sprintf.
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
