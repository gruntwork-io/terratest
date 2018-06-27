package terragrunt

import (
	"fmt"
	"reflect"
	"strings"
)

// FormatArgs converts the inputs to a format palatable to terragrunt. This includes converting the given vars to the
// format the Terragrunt CLI expects (-var key=value).
func FormatArgs(customVars map[string]interface{}, args ...string) []string {
	varsAsArgs := FormatTerragruntVarsAsArgs(customVars)
	return append(args, varsAsArgs...)
}

// FormatTerragruntVarsAsArgs formats the given variables as command-line args for Terragrunt (e.g. of the format
// -var key=value).
func FormatTerragruntVarsAsArgs(vars map[string]interface{}) []string {
	return formatTerragruntArgs(vars, "-var")
}

// FormatTerragruntBackendConfigAsArgs formats the given variables as backend config args for Terragrunt (e.g. of the
// format -backend-config key=value).
func FormatTerragruntBackendConfigAsArgs(vars map[string]interface{}) []string {
	return formatTerragruntArgs(vars, "-backend-config")
}

// Format the given vars into 'Terragrunt' format, with each var being prefixed with the given prefix.
func formatTerragruntArgs(vars map[string]interface{}, prefix string) []string {
	var args []string

	for key, value := range vars {
		hclString := toHclString(value)
		argValue := fmt.Sprintf("%s=%s", key, hclString)
		args = append(args, prefix, argValue)
	}

	return args
}

// Terragrunt allows you to pass in command-line variables using HCL syntax (e.g. -var foo=[1,2,3]). Unfortunately,
// while their golang hcl library can convert an HCL string to a Go type, they don't seem to offer a library to convert
// arbitrary Go types to an HCL string. Therefore, this method is a simple implementation that correctly handles
// ints, booleans, lists, and maps. Everything else is forced into a string using Sprintf. Hopefully, this approach is
// good enough for the type of variables we deal with in Terratest.
func toHclString(value interface{}) string {
	// Ideally, we'd use a type switch here to identify slices and maps, but we can't do that, because Go doesn't
	// support generics, and the type switch only matches concrete types. So we could match []interface{}, but if
	// a user passes in []string{}, that would NOT match (the same logic applies to maps). Therefore, we have to
	// use reflection and manually convert into []interface{} and map[string]interface{}.

	if slice, isSlice := tryToConvertToGenericSlice(value); isSlice {
		return sliceToHclString(slice)
	} else if m, isMap := tryToConvertToGenericMap(value); isMap {
		return mapToHclString(m)
	} else {
		return primitiveToHclString(value)
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
		hclValue := toHclString(value)
		hclValues = append(hclValues, hclValue)
	}

	return fmt.Sprintf("[%s]", strings.Join(hclValues, ", "))
}

// Convert a map to an HCL string. See ToHclString for details.
func mapToHclString(m map[string]interface{}) string {
	keyValuePairs := []string{}

	for key, value := range m {
		keyValuePair := fmt.Sprintf("%s = %s", key, toHclString(value))
		keyValuePairs = append(keyValuePairs, keyValuePair)
	}

	return fmt.Sprintf("{%s}", strings.Join(keyValuePairs, ", "))
}

// Convert a primitive, such as a bool, int, or string, to an HCL string. If this isn't a primitive, force its value
// using Sprintf. See ToHclString for details.
func primitiveToHclString(value interface{}) string {
	switch v := value.(type) {

	// Terragrunt treats a boolean true as a 1 and a boolean false as a 0. It's best to convert to these ints when
	// passing booleans as -var parameters. Moreover, due to a Terragrunt bug
	// (https://github.com/hashicorp/terragrunt/issues/7962), all ints must be wrapped as strings.
	case bool:
		if v {
			return "\"1\""
		}
		return "\"0\""

	// Note: due to a Terragrunt bug (https://github.com/hashicorp/terragrunt/issues/7962), we can't use proper HCL
	// syntax for ints have to wrap them as strings by falling through to the default case
	//case int: return strconv.Itoa(v)

	default:
		return fmt.Sprintf("\"%v\"", v)
	}
}
