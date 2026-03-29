// Package formatting provides internal utilities for formatting Terraform/Terragrunt CLI arguments.
package formatting

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// FormatBackendConfigAsArgs formats backend configuration as Terraform CLI args.
// Example: {"bucket": "my-bucket"} -> ["-backend-config=bucket=my-bucket"]
func FormatBackendConfigAsArgs(vars map[string]interface{}) []string {
	return formatTerraformArgs(vars, "-backend-config", false, true)
}

// FormatPluginDirAsArgs formats plugin directory as a Terraform CLI arg.
// Returns nil if pluginDir is empty.
func FormatPluginDirAsArgs(pluginDir string) []string {
	if pluginDir == "" {
		return nil
	}

	return []string{"-plugin-dir=" + pluginDir}
}

// formatTerraformArgs formats vars as CLI args with the given prefix.
func formatTerraformArgs(vars map[string]interface{}, prefix string, useSpaceAsSeparator bool, omitNil bool) []string {
	var args []string

	for key, value := range vars {
		var argValue string
		if omitNil && value == nil {
			argValue = key
		} else {
			hclString := ToHCLString(value, false)
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

// ToHCLString converts Go values to HCL-formatted strings for Terraform CLI arguments.
// Handles primitives, slices, and maps. Example: []int{1,2,3} -> "[1, 2, 3]"
func ToHCLString(value interface{}, isNested bool) string {
	if slice, isSlice := tryToConvertToGenericSlice(value); isSlice {
		return sliceToHclString(slice)
	} else if m, isMap := tryToConvertToGenericMap(value); isMap {
		return mapToHclString(m)
	} else {
		return primitiveToHclString(value, isNested)
	}
}

// tryToConvertToGenericSlice converts any slice type to []interface{} using reflection.
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

// tryToConvertToGenericMap converts any map[string]T to map[string]interface{} using reflection.
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

func sliceToHclString(slice []interface{}) string {
	hclValues := make([]string, 0, len(slice))

	for _, value := range slice {
		hclValue := ToHCLString(value, true)
		hclValues = append(hclValues, hclValue)
	}

	return fmt.Sprintf("[%s]", strings.Join(hclValues, ", "))
}

func mapToHclString(m map[string]interface{}) string {
	keyValuePairs := make([]string, 0, len(m))

	for key, value := range m {
		keyValuePair := fmt.Sprintf(`"%s" = %s`, key, ToHCLString(value, true))
		keyValuePairs = append(keyValuePairs, keyValuePair)
	}

	return fmt.Sprintf("{%s}", strings.Join(keyValuePairs, ", "))
}

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
			return "\"" + v + "\""
		}

		return v

	default:
		return fmt.Sprintf("%v", v)
	}
}
