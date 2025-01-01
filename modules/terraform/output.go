package terraform

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"

	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// Output calls terraform output for the given variable and return its string value representation.
// It only designed to work with primitive terraform types: string, number and bool.
// Please use OutputStruct for anything else.
func Output(t testing.TestingT, options *Options, key string) string {
	out, err := OutputE(t, options, key)
	require.NoError(t, err)
	return out
}

// OutputE calls terraform output for the given variable and return its string value representation.
// It only designed to work with primitive terraform types: string, number and bool.
// Please use OutputStructE for anything else.
func OutputE(t testing.TestingT, options *Options, key string) (string, error) {
	var val interface{}
	err := OutputStructE(t, options, key, &val)
	return fmt.Sprintf("%v", val), err
}

// OutputRequired calls terraform output for the given variable and return its value. If the value is empty, fail the test.
func OutputRequired(t testing.TestingT, options *Options, key string) string {
	out, err := OutputRequiredE(t, options, key)
	require.NoError(t, err)
	return out
}

// OutputRequiredE calls terraform output for the given variable and return its value. If the value is empty, return an error.
func OutputRequiredE(t testing.TestingT, options *Options, key string) (string, error) {
	out, err := OutputE(t, options, key)

	if err != nil {
		return "", err
	}
	if out == "" {
		return "", EmptyOutput(key)
	}

	return out, nil
}

// parseListOfMaps takes a list of maps and parses the types.
// It is mainly a wrapper for parseMap to support lists.
func parseListOfMaps(l []interface{}) ([]map[string]interface{}, error) {
	var result []map[string]interface{}

	for _, v := range l {

		asMap, isMap := v.(map[string]interface{})
		if !isMap {
			err := errors.New("Type switching to map[string]interface{} failed.")
			return nil, err
		}

		m, err := parseMap(asMap)

		if err != nil {
			return nil, err
		}
		result = append(result, m)
	}

	return result, nil

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
func parseMap(m map[string]interface{}) (map[string]interface{}, error) {

	result := make(map[string]interface{})

	for k, v := range m {
		switch vt := v.(type) {
		case map[string]interface{}:
			nestedMap, err := parseMap(vt)
			if err != nil {
				return nil, err
			}
			result[k] = nestedMap
		case []interface{}:
			nestedList, err := parseListOfMaps(vt)
			if err != nil {
				return nil, err
			}
			result[k] = nestedList
		case float64:
			testInt, err := strconv.ParseInt((fmt.Sprintf("%v", vt)), 10, 0)
			if err == nil {
				result[k] = int(testInt)
			} else {
				result[k] = vt
			}
		default:
			result[k] = vt
		}

	}

	return result, nil
}

// OutputMapOfObjects calls terraform output for the given variable and returns its value as a map of lists/maps.
// If the output value is not a map of lists/maps, then it fails the test.
func OutputMapOfObjects(t testing.TestingT, options *Options, key string) map[string]interface{} {
	out, err := OutputMapOfObjectsE(t, options, key)
	require.NoError(t, err)
	return out
}

// OutputMapOfObjectsE calls terraform output for the given variable and returns its value as a map of lists/maps.
// Also returns an error object if an error was generated.
// If the output value is not a map of lists/maps, then it fails the test.
func OutputMapOfObjectsE(t testing.TestingT, options *Options, key string) (map[string]interface{}, error) {
	out, err := OutputJsonE(t, options, key)

	if err != nil {
		return nil, err
	}

	var output map[string]interface{}

	if err := json.Unmarshal([]byte(out), &output); err != nil {
		return nil, err
	}

	return parseMap(output)
}

// OutputListOfObjects calls terraform output for the given variable and returns its value as a list of maps/lists.
// If the output value is not a list of maps/lists, then it fails the test.
func OutputListOfObjects(t testing.TestingT, options *Options, key string) []map[string]interface{} {
	out, err := OutputListOfObjectsE(t, options, key)
	require.NoError(t, err)
	return out
}

// OutputListOfObjectsE calls terraform output for the given variable and returns its value as a list of maps/lists.
// Also returns an error object if an error was generated.
// If the output value is not a list of maps/lists, then it fails the test.
func OutputListOfObjectsE(t testing.TestingT, options *Options, key string) ([]map[string]interface{}, error) {
	out, err := OutputJsonE(t, options, key)

	if err != nil {
		return nil, err
	}

	var output []map[string]interface{}

	if err := json.Unmarshal([]byte(out), &output); err != nil {
		return nil, err
	}

	var result []map[string]interface{}

	for _, m := range output {
		newMap, err := parseMap(m)

		if err != nil {
			return nil, err
		}

		result = append(result, newMap)
	}

	return result, nil
}

// OutputList calls terraform output for the given variable and returns its value as a list.
// If the output value is not a list type, then it fails the test.
func OutputList(t testing.TestingT, options *Options, key string) []string {
	out, err := OutputListE(t, options, key)
	require.NoError(t, err)
	return out
}

// OutputListE calls terraform output for the given variable and returns its value as a list.
// If the output value is not a list type, then it returns an error.
func OutputListE(t testing.TestingT, options *Options, key string) ([]string, error) {
	out, err := OutputJsonE(t, options, key)
	if err != nil {
		return nil, err
	}

	var output interface{}
	if err := json.Unmarshal([]byte(out), &output); err != nil {
		return nil, err
	}

	if outputList, isList := output.([]interface{}); isList {
		return parseListOutputTerraform(outputList, key)
	}

	return nil, UnexpectedOutputType{Key: key, ExpectedType: "map or list", ActualType: reflect.TypeOf(output).String()}
}

// Parse a list output in the format it is returned by Terraform 0.12 and newer versions
func parseListOutputTerraform(outputList []interface{}, key string) ([]string, error) {
	list := []string{}

	for _, item := range outputList {
		list = append(list, fmt.Sprintf("%v", item))
	}

	return list, nil
}

// OutputMap calls terraform output for the given variable and returns its value as a map.
// If the output value is not a map type, then it fails the test.
func OutputMap(t testing.TestingT, options *Options, key string) map[string]string {
	out, err := OutputMapE(t, options, key)
	require.NoError(t, err)
	return out
}

// OutputMapE calls terraform output for the given variable and returns its value as a map.
// If the output value is not a map type, then it returns an error.
func OutputMapE(t testing.TestingT, options *Options, key string) (map[string]string, error) {
	out, err := OutputJsonE(t, options, key)
	if err != nil {
		return nil, err
	}

	outputMap := map[string]interface{}{}
	if err := json.Unmarshal([]byte(out), &outputMap); err != nil {
		return nil, err
	}

	resultMap := make(map[string]string)
	for k, v := range outputMap {
		resultMap[k] = fmt.Sprintf("%v", v)
	}
	return resultMap, nil
}

// OutputForKeys calls terraform output for the given key list and returns values as a map.
// If keys not found in the output, fails the test
func OutputForKeys(t testing.TestingT, options *Options, keys []string) map[string]interface{} {
	out, err := OutputForKeysE(t, options, keys)
	require.NoError(t, err)
	return out
}

// OutputJson calls terraform output for the given variable and returns the
// result as the json string.
// If key is an empty string, it will return all the output variables.
func OutputJson(t testing.TestingT, options *Options, key string) string {
	str, err := OutputJsonE(t, options, key)
	require.NoError(t, err)
	return str
}

// OutputJsonE calls terraform output for the given variable and returns the
// result as the json string.
// If key is an empty string, it will return all the output variables.
func OutputJsonE(t testing.TestingT, options *Options, key string) (string, error) {
	args := []string{"output", "-no-color", "-json"}
	if key != "" {
		args = append(args, key)
	}

	return RunTerraformCommandAndGetStdoutE(t, options, args...)
}

// OutputStruct calls terraform output for the given variable and stores the
// result in the value pointed to by v. If v is nil or not a pointer, or if
// the value returned by Terraform is not appropriate for a given target type,
// it fails the test.
func OutputStruct(t testing.TestingT, options *Options, key string, v interface{}) {
	err := OutputStructE(t, options, key, v)
	require.NoError(t, err)
}

// OutputStructE calls terraform output for the given variable and stores the
// result in the value pointed to by v. If v is nil or not a pointer, or if
// the value returned by Terraform is not appropriate for a given target type,
// it returns an error.
func OutputStructE(t testing.TestingT, options *Options, key string, v interface{}) error {
	out, err := OutputJsonE(t, options, key)
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(out), &v)
}

// OutputForKeysE calls terraform output for the given key list and returns values as a map.
// The returned values are of type interface{} and need to be type casted as necessary. Refer to output_test.go
func OutputForKeysE(t testing.TestingT, options *Options, keys []string) (map[string]interface{}, error) {
	out, err := OutputJsonE(t, options, "")
	if err != nil {
		return nil, err
	}

	outputMap := map[string]map[string]interface{}{}
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

	resultMap := make(map[string]interface{})
	for _, key := range keys {
		value, containsValue := outputMap[key]["value"]
		if !containsValue {
			return nil, OutputKeyNotFound(string(key))
		}
		resultMap[key] = value
	}
	return resultMap, nil
}

// OutputAll calls terraform output returns all values as a map.
// If there is error fetching the output, fails the test
func OutputAll(t testing.TestingT, options *Options) map[string]interface{} {
	out, err := OutputAllE(t, options)
	require.NoError(t, err)
	return out
}

// OutputAllE calls terraform and returns all the outputs as a map
func OutputAllE(t testing.TestingT, options *Options) (map[string]interface{}, error) {
	return OutputForKeysE(t, options, nil)
}
