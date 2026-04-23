package terraform

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

const skipJSONLogLine = " msg="

var (
	// ansiLineRegex matches lines starting with ANSI escape codes for text formatting (e.g., colors, styles).
	ansiLineRegex = regexp.MustCompile(`(?m)^\x1b\[[0-9;]*m.*`)
	// tgLogLevel matches log lines containing fields for time, level, prefix, binary, and message, each with non-whitespace values.
	tgLogLevel = regexp.MustCompile(`.*time=\S+ level=\S+ prefix=\S+ binary=\S+ msg=.*`)
)

// OutputContext calls terraform output for the given variable and returns its string value representation.
// It is only designed to work with primitive terraform types: string, number and bool.
// Please use OutputStructContext for anything else. The context argument can be used for cancellation or
// timeout control.
func OutputContext(t testing.TestingT, ctx context.Context, options *Options, key string) string {
	out, err := OutputContextE(t, ctx, options, key)
	require.NoError(t, err)

	return out
}

// OutputContextE calls terraform output for the given variable and returns its string value representation.
// It is only designed to work with primitive terraform types: string, number and bool.
// Please use OutputStructContextE for anything else. The context argument can be used for cancellation or
// timeout control.
func OutputContextE(t testing.TestingT, ctx context.Context, options *Options, key string) (string, error) {
	var val any

	err := OutputStructContextE(t, ctx, options, key, &val)

	return fmt.Sprintf("%v", val), err
}

// OutputRequiredContext calls terraform output for the given variable and returns its value. If the value is empty,
// the test is failed. The context argument can be used for cancellation or timeout control.
func OutputRequiredContext(t testing.TestingT, ctx context.Context, options *Options, key string) string {
	out, err := OutputRequiredContextE(t, ctx, options, key)
	require.NoError(t, err)

	return out
}

// OutputRequiredContextE calls terraform output for the given variable and returns its value. If the value is empty,
// an error is returned. The context argument can be used for cancellation or timeout control.
func OutputRequiredContextE(t testing.TestingT, ctx context.Context, options *Options, key string) (string, error) {
	out, err := OutputContextE(t, ctx, options, key)
	if err != nil {
		return "", err
	}

	if out == "" {
		return "", EmptyOutput(key)
	}

	return out, nil
}

// parseMap takes a map of interfaces and parses the types.
// It is recursive which allows it to support complex nested structures.
// At this time, this function uses https://golang.org/pkg/strconv/#ParseInt
// to determine if a number should be a float or an int. For this reason, if you are
// expecting a float with a zero as the "tenth" you will need to manually convert
// the return value to a float.
//
// This function exists to map return values of the terraform outputs to intuitive
// types. ie, if you are expecting a value of "1" you are implicitly expecting an int.
//
// This also allows the work to be executed recursively to support complex data types.
func parseMap(m map[string]any) (map[string]any, error) {
	result := make(map[string]any)

	for k, v := range m {
		switch vt := v.(type) {
		case map[string]any:
			nestedMap, err := parseMap(vt)
			if err != nil {
				return nil, err
			}

			result[k] = nestedMap
		case []any:
			nestedList, err := parseList(vt)
			if err != nil {
				return nil, err
			}

			result[k] = nestedList
		case float64:
			result[k] = parseFloat(vt)
		default:
			result[k] = vt
		}
	}

	return result, nil
}

func parseList(items []any) (_ []any, err error) {
	for i, v := range items {
		rv := reflect.ValueOf(v)

		switch rv.Kind() { //nolint:exhaustive // only map, slice/array, and float64 need special handling
		case reflect.Map:
			items[i], err = parseMap(rv.Interface().(map[string]any))
		case reflect.Slice, reflect.Array:
			items[i], err = parseList(rv.Interface().([]any))
		case reflect.Float64:
			items[i] = parseFloat(v)
		}

		if err != nil {
			return nil, err
		}
	}

	return items, nil
}

func parseFloat(v any) any {
	testInt, err := strconv.ParseInt((fmt.Sprintf("%v", v)), 10, 0)
	if err == nil {
		return int(testInt)
	}

	return v
}

// OutputMapOfObjectsContext calls terraform output for the given variable and returns its value as a map of
// lists/maps. If the output value is not a map of lists/maps, then it fails the test. The context argument
// can be used for cancellation or timeout control.
func OutputMapOfObjectsContext(t testing.TestingT, ctx context.Context, options *Options, key string) map[string]any {
	out, err := OutputMapOfObjectsContextE(t, ctx, options, key)
	require.NoError(t, err)

	return out
}

// OutputMapOfObjectsContextE calls terraform output for the given variable and returns its value as a map of
// lists/maps. Also returns an error object if an error was generated. If the output value is not a map of
// lists/maps, then it returns an error. The context argument can be used for cancellation or timeout control.
func OutputMapOfObjectsContextE(t testing.TestingT, ctx context.Context, options *Options, key string) (map[string]any, error) {
	out, err := OutputJSONContextE(t, ctx, options, key)
	if err != nil {
		return nil, err
	}

	var output map[string]any

	if err := json.Unmarshal([]byte(out), &output); err != nil {
		return nil, err
	}

	return parseMap(output)
}

// OutputListOfObjectsContext calls terraform output for the given variable and returns its value as a list of
// maps/lists. If the output value is not a list of maps/lists, then it fails the test. The context argument
// can be used for cancellation or timeout control.
func OutputListOfObjectsContext(t testing.TestingT, ctx context.Context, options *Options, key string) []map[string]any {
	out, err := OutputListOfObjectsContextE(t, ctx, options, key)
	require.NoError(t, err)

	return out
}

// OutputListOfObjectsContextE calls terraform output for the given variable and returns its value as a list of
// maps/lists. Also returns an error object if an error was generated. If the output value is not a list of
// maps/lists, then it returns an error. The context argument can be used for cancellation or timeout control.
func OutputListOfObjectsContextE(t testing.TestingT, ctx context.Context, options *Options, key string) ([]map[string]any, error) {
	out, err := OutputJSONContextE(t, ctx, options, key)
	if err != nil {
		return nil, err
	}

	var output []map[string]any

	if err := json.Unmarshal([]byte(out), &output); err != nil {
		return nil, err
	}

	var result []map[string]any

	for _, m := range output {
		newMap, err := parseMap(m)
		if err != nil {
			return nil, err
		}

		result = append(result, newMap)
	}

	return result, nil
}

// OutputListContext calls terraform output for the given variable and returns its value as a list.
// If the output value is not a list type, then it fails the test. The context argument can be used for
// cancellation or timeout control.
func OutputListContext(t testing.TestingT, ctx context.Context, options *Options, key string) []string {
	out, err := OutputListContextE(t, ctx, options, key)
	require.NoError(t, err)

	return out
}

// OutputListContextE calls terraform output for the given variable and returns its value as a list.
// If the output value is not a list type, then it returns an error. The context argument can be used for
// cancellation or timeout control.
func OutputListContextE(t testing.TestingT, ctx context.Context, options *Options, key string) ([]string, error) {
	out, err := OutputJSONContextE(t, ctx, options, key)
	if err != nil {
		return nil, err
	}

	var output any

	if err := json.Unmarshal([]byte(out), &output); err != nil {
		return nil, err
	}

	if outputList, isList := output.([]any); isList {
		return parseListOutputTerraform(outputList)
	}

	return nil, UnexpectedOutputType{Key: key, ExpectedType: "map or list", ActualType: reflect.TypeOf(output).String()}
}

// parseListOutputTerraform parses a list output in the format returned by Terraform 0.12 and newer versions.
func parseListOutputTerraform(outputList []any) ([]string, error) {
	list := make([]string, 0, len(outputList))

	for _, item := range outputList {
		list = append(list, fmt.Sprintf("%v", item))
	}

	return list, nil
}

// OutputMapContext calls terraform output for the given variable and returns its value as a map.
// If the output value is not a map type, then it fails the test. The context argument can be used for
// cancellation or timeout control.
func OutputMapContext(t testing.TestingT, ctx context.Context, options *Options, key string) map[string]string {
	out, err := OutputMapContextE(t, ctx, options, key)
	require.NoError(t, err)

	return out
}

// OutputMapContextE calls terraform output for the given variable and returns its value as a map.
// If the output value is not a map type, then it returns an error. The context argument can be used for
// cancellation or timeout control.
func OutputMapContextE(t testing.TestingT, ctx context.Context, options *Options, key string) (map[string]string, error) {
	out, err := OutputJSONContextE(t, ctx, options, key)
	if err != nil {
		return nil, err
	}

	outputMap := map[string]any{}

	if err := json.Unmarshal([]byte(out), &outputMap); err != nil {
		return nil, err
	}

	resultMap := make(map[string]string)

	for k, v := range outputMap {
		resultMap[k] = fmt.Sprintf("%v", v)
	}

	return resultMap, nil
}

// OutputForKeysContext calls terraform output for the given key list and returns values as a map.
// If keys are not found in the output, it fails the test. The context argument can be used for
// cancellation or timeout control.
func OutputForKeysContext(t testing.TestingT, ctx context.Context, options *Options, keys []string) map[string]any {
	out, err := OutputForKeysContextE(t, ctx, options, keys)
	require.NoError(t, err)

	return out
}

// OutputForKeysContextE calls terraform output for the given key list and returns values as a map.
// The returned values are of type any and need to be type casted as necessary. The context
// argument can be used for cancellation or timeout control.
func OutputForKeysContextE(t testing.TestingT, ctx context.Context, options *Options, keys []string) (map[string]any, error) {
	out, err := OutputJSONContextE(t, ctx, options, "")
	if err != nil {
		return nil, err
	}

	outputMap := map[string]map[string]any{}

	if err := json.Unmarshal([]byte(out), &outputMap); err != nil {
		return nil, err
	}

	if keys == nil {
		outputKeys := make([]string, 0, len(outputMap))

		for k := range outputMap {
			outputKeys = append(outputKeys, k)
		}

		keys = outputKeys
	}

	resultMap := make(map[string]any)

	for _, key := range keys {
		value, containsValue := outputMap[key]["value"]
		if !containsValue {
			return nil, OutputKeyNotFound(key)
		}

		resultMap[key] = value
	}

	return resultMap, nil
}

// OutputJSONContext calls terraform output for the given variable and returns the result as a JSON string.
// If key is an empty string, it will return all the output variables. The context argument can be used for
// cancellation or timeout control.
func OutputJSONContext(t testing.TestingT, ctx context.Context, options *Options, key string) string {
	str, err := OutputJSONContextE(t, ctx, options, key)
	require.NoError(t, err)

	return str
}

// OutputJSONContextE calls terraform output for the given variable and returns the result as a JSON string.
// If key is an empty string, it will return all the output variables. The context argument can be used for
// cancellation or timeout control.
func OutputJSONContextE(t testing.TestingT, ctx context.Context, options *Options, key string) (string, error) {
	args := []string{"output", "-no-color", "-json"}

	args = append(args, options.ExtraArgs.Output...)

	if key != "" {
		args = append(args, key)
	}

	rawJSON, err := RunTerraformCommandAndGetStdoutContextE(t, ctx, options, args...)
	if err != nil {
		return rawJSON, err
	}

	return cleanJSON(rawJSON)
}

// OutputStructContext calls terraform output for the given variable and stores the result in the value
// pointed to by v. If v is nil or not a pointer, or if the value returned by Terraform is not appropriate
// for a given target type, it fails the test. The context argument can be used for cancellation or timeout
// control.
func OutputStructContext(t testing.TestingT, ctx context.Context, options *Options, key string, v any) {
	err := OutputStructContextE(t, ctx, options, key, v)
	require.NoError(t, err)
}

// OutputStructContextE calls terraform output for the given variable and stores the result in the value
// pointed to by v. If v is nil or not a pointer, or if the value returned by Terraform is not appropriate
// for a given target type, it returns an error. The context argument can be used for cancellation or timeout
// control.
func OutputStructContextE(t testing.TestingT, ctx context.Context, options *Options, key string, v any) error {
	out, err := OutputJSONContextE(t, ctx, options, key)
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(out), &v)
}

// OutputAllContext calls terraform output and returns all values as a map.
// If there is an error fetching the output, it fails the test. The context argument can be used for
// cancellation or timeout control.
func OutputAllContext(t testing.TestingT, ctx context.Context, options *Options) map[string]any {
	out, err := OutputAllContextE(t, ctx, options)
	require.NoError(t, err)

	return out
}

// OutputAllContextE calls terraform output and returns all the outputs as a map. The context argument can
// be used for cancellation or timeout control.
func OutputAllContextE(t testing.TestingT, ctx context.Context, options *Options) (map[string]any, error) {
	return OutputForKeysContextE(t, ctx, options, nil)
}

// Output calls terraform output for the given variable and return its string value representation.
// It only designed to work with primitive terraform types: string, number and bool.
// Please use OutputStruct for anything else.
//
// Deprecated: Use [OutputContext] instead.
func Output(t testing.TestingT, options *Options, key string) string {
	return OutputContext(t, context.Background(), options, key)
}

// OutputE calls terraform output for the given variable and return its string value representation.
// It only designed to work with primitive terraform types: string, number and bool.
// Please use OutputStructE for anything else.
//
// Deprecated: Use [OutputContextE] instead.
func OutputE(t testing.TestingT, options *Options, key string) (string, error) {
	return OutputContextE(t, context.Background(), options, key)
}

// OutputRequired calls terraform output for the given variable and return its value. If the value is empty, fail the test.
//
// Deprecated: Use [OutputRequiredContext] instead.
func OutputRequired(t testing.TestingT, options *Options, key string) string {
	return OutputRequiredContext(t, context.Background(), options, key)
}

// OutputRequiredE calls terraform output for the given variable and return its value. If the value is empty, return an error.
//
// Deprecated: Use [OutputRequiredContextE] instead.
func OutputRequiredE(t testing.TestingT, options *Options, key string) (string, error) {
	return OutputRequiredContextE(t, context.Background(), options, key)
}

// OutputMapOfObjects calls terraform output for the given variable and returns its value as a map of lists/maps.
// If the output value is not a map of lists/maps, then it fails the test.
//
// Deprecated: Use [OutputMapOfObjectsContext] instead.
func OutputMapOfObjects(t testing.TestingT, options *Options, key string) map[string]any {
	return OutputMapOfObjectsContext(t, context.Background(), options, key)
}

// OutputMapOfObjectsE calls terraform output for the given variable and returns its value as a map of lists/maps.
// Also returns an error object if an error was generated.
// If the output value is not a map of lists/maps, then it fails the test.
//
// Deprecated: Use [OutputMapOfObjectsContextE] instead.
func OutputMapOfObjectsE(t testing.TestingT, options *Options, key string) (map[string]any, error) {
	return OutputMapOfObjectsContextE(t, context.Background(), options, key)
}

// OutputListOfObjects calls terraform output for the given variable and returns its value as a list of maps/lists.
// If the output value is not a list of maps/lists, then it fails the test.
//
// Deprecated: Use [OutputListOfObjectsContext] instead.
func OutputListOfObjects(t testing.TestingT, options *Options, key string) []map[string]any {
	return OutputListOfObjectsContext(t, context.Background(), options, key)
}

// OutputListOfObjectsE calls terraform output for the given variable and returns its value as a list of maps/lists.
// Also returns an error object if an error was generated.
// If the output value is not a list of maps/lists, then it fails the test.
//
// Deprecated: Use [OutputListOfObjectsContextE] instead.
func OutputListOfObjectsE(t testing.TestingT, options *Options, key string) ([]map[string]any, error) {
	return OutputListOfObjectsContextE(t, context.Background(), options, key)
}

// OutputList calls terraform output for the given variable and returns its value as a list.
// If the output value is not a list type, then it fails the test.
//
// Deprecated: Use [OutputListContext] instead.
func OutputList(t testing.TestingT, options *Options, key string) []string {
	return OutputListContext(t, context.Background(), options, key)
}

// OutputListE calls terraform output for the given variable and returns its value as a list.
// If the output value is not a list type, then it returns an error.
//
// Deprecated: Use [OutputListContextE] instead.
func OutputListE(t testing.TestingT, options *Options, key string) ([]string, error) {
	return OutputListContextE(t, context.Background(), options, key)
}

// OutputMap calls terraform output for the given variable and returns its value as a map.
// If the output value is not a map type, then it fails the test.
//
// Deprecated: Use [OutputMapContext] instead.
func OutputMap(t testing.TestingT, options *Options, key string) map[string]string {
	return OutputMapContext(t, context.Background(), options, key)
}

// OutputMapE calls terraform output for the given variable and returns its value as a map.
// If the output value is not a map type, then it returns an error.
//
// Deprecated: Use [OutputMapContextE] instead.
func OutputMapE(t testing.TestingT, options *Options, key string) (map[string]string, error) {
	return OutputMapContextE(t, context.Background(), options, key)
}

// OutputForKeys calls terraform output for the given key list and returns values as a map.
// If keys not found in the output, fails the test
//
// Deprecated: Use [OutputForKeysContext] instead.
func OutputForKeys(t testing.TestingT, options *Options, keys []string) map[string]any {
	return OutputForKeysContext(t, context.Background(), options, keys)
}

// OutputForKeysE calls terraform output for the given key list and returns values as a map.
// The returned values are of type any and need to be type casted as necessary. Refer to output_test.go
//
// Deprecated: Use [OutputForKeysContextE] instead.
func OutputForKeysE(t testing.TestingT, options *Options, keys []string) (map[string]any, error) {
	return OutputForKeysContextE(t, context.Background(), options, keys)
}

// OutputJSON calls terraform output for the given variable and returns the
// result as the JSON string.
// If key is an empty string, it will return all the output variables.
//
// Deprecated: Use [OutputJSONContext] instead.
func OutputJSON(t testing.TestingT, options *Options, key string) string {
	return OutputJSONContext(t, context.Background(), options, key)
}

// OutputJSONE calls terraform output for the given variable and returns the
// result as the JSON string.
// If key is an empty string, it will return all the output variables.
//
// Deprecated: Use [OutputJSONContextE] instead.
func OutputJSONE(t testing.TestingT, options *Options, key string) (string, error) {
	return OutputJSONContextE(t, context.Background(), options, key)
}

// OutputJson calls terraform output for the given variable and returns the
// result as the JSON string.
// If key is an empty string, it will return all the output variables.
//
// Deprecated: Use [OutputJSONContext] instead.
func OutputJson(t testing.TestingT, options *Options, key string) string { //nolint:revive,staticcheck // preserving deprecated function name
	return OutputJSONContext(t, context.Background(), options, key)
}

// OutputJsonE calls terraform output for the given variable and returns the
// result as the JSON string.
// If key is an empty string, it will return all the output variables.
//
// Deprecated: Use [OutputJSONContextE] instead.
func OutputJsonE(t testing.TestingT, options *Options, key string) (string, error) { //nolint:revive,staticcheck // preserving deprecated function name
	return OutputJSONContextE(t, context.Background(), options, key)
}

// OutputStruct calls terraform output for the given variable and stores the
// result in the value pointed to by v. If v is nil or not a pointer, or if
// the value returned by Terraform is not appropriate for a given target type,
// it fails the test.
//
// Deprecated: Use [OutputStructContext] instead.
func OutputStruct(t testing.TestingT, options *Options, key string, v any) {
	OutputStructContext(t, context.Background(), options, key, v)
}

// OutputStructE calls terraform output for the given variable and stores the
// result in the value pointed to by v. If v is nil or not a pointer, or if
// the value returned by Terraform is not appropriate for a given target type,
// it returns an error.
//
// Deprecated: Use [OutputStructContextE] instead.
func OutputStructE(t testing.TestingT, options *Options, key string, v any) error {
	return OutputStructContextE(t, context.Background(), options, key, v)
}

// OutputAll calls terraform output returns all values as a map.
// If there is error fetching the output, fails the test
//
// Deprecated: Use [OutputAllContext] instead.
func OutputAll(t testing.TestingT, options *Options) map[string]any {
	return OutputAllContext(t, context.Background(), options)
}

// OutputAllE calls terraform and returns all the outputs as a map
//
// Deprecated: Use [OutputAllContextE] instead.
func OutputAllE(t testing.TestingT, options *Options) (map[string]any, error) {
	return OutputAllContextE(t, context.Background(), options)
}

// cleanJSON removes ANSI characters from the JSON and normalizes formatting.
func cleanJSON(input string) (string, error) {
	// Remove ANSI escape codes
	cleaned := ansiLineRegex.ReplaceAllString(input, "")
	cleaned = tgLogLevel.ReplaceAllString(cleaned, "")

	lines := strings.Split(cleaned, "\n")

	var result []string

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		if trimmed != "" && !strings.Contains(trimmed, skipJSONLogLine) {
			result = append(result, trimmed)
		}
	}

	ansiClean := strings.Join(result, "\n")

	var jsonObj any

	if err := json.Unmarshal([]byte(ansiClean), &jsonObj); err != nil {
		return "", err
	}

	// Format JSON output with indentation
	normalized, err := json.MarshalIndent(jsonObj, "", "  ")
	if err != nil {
		return "", err
	}

	return string(normalized), nil
}
