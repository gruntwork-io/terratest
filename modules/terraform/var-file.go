package terraform

import (
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/stretchr/testify/require"
)

// GetVariableAsStringFromVarFile Gets the string represention of a variable from a provided input file found in VarFile
// For list or map, use GetVariableAsListFromVarFile or GetVariableAsMapFromVarFile, respectively.
func GetVariableAsStringFromVarFile(t *testing.T, fileName string, key string) string {
	result, err := GetVariableAsStringFromVarFileE(t, fileName, key)
	require.NoError(t, err)

	return result
}

// GetVariableAsStringFromVarFileE Gets the string represention of a variable from a provided input file found in VarFile
// Will return an error if GetAllVariablesFromVarFileE returns an error or the key provided does not exist in the file.
// For list or map, use GetVariableAsListFromVarFile or GetVariableAsMapFromVarFile, respectively.
func GetVariableAsStringFromVarFileE(t *testing.T, fileName string, key string) (string, error) {
	var variables map[string]interface{}

	err := GetAllVariablesFromVarFileE(t, fileName, &variables)

	if err != nil {
		return "", err
	}

	variable, exists := variables[key]

	if !exists {
		return "", InputFileKeyNotFound{FilePath: fileName, Key: key}
	}

	return fmt.Sprintf("%v", variable), nil
}

// GetVariableAsMapFromVarFile Gets the map represention of a variable from a provided input file found in VarFile
// Note that this returns a map of strings. For maps containing complex types, use GetAllVariablesFromVarFile.
func GetVariableAsMapFromVarFile(t *testing.T, fileName string, key string) map[string]string {
	result, err := GetVariableAsMapFromVarFileE(t, fileName, key)
	require.NoError(t, err)

	return result
}

// GetVariableAsMapFromVarFileE Gets the map represention of a variable from a provided input file found in VarFile.
// Note that this returns a map of strings. For maps containing complex types, use GetAllVariablesFromVarFile
// Returns an error if GetAllVariablesFromVarFileE returns an error, the key provided does not exist, or the value associated with the key is not a map
func GetVariableAsMapFromVarFileE(t *testing.T, fileName string, key string) (map[string]string, error) {
	var variables map[string]interface{}

	resultMap := make(map[string]string)
	err := GetAllVariablesFromVarFileE(t, fileName, &variables)

	if err != nil {
		return nil, err
	}

	variable, exists := variables[key]

	if !exists {
		return nil, InputFileKeyNotFound{FilePath: fileName, Key: key}
	}

	if reflect.TypeOf(variable).String() != "[]map[string]interface {}" {
		return nil, UnexpectedOutputType{Key: key, ExpectedType: "[]map[string]interface {}", ActualType: reflect.TypeOf(variable).String()}
	}

	mapKeys := variable.([]map[string]interface{})

	if len(mapKeys) == 0 {
		return nil, errors.New("no map keys could be found for given map")
	}

	for mapKey, mapVal := range mapKeys[0] {
		resultMap[mapKey] = fmt.Sprintf("%v", mapVal)
	}

	return resultMap, nil
}

// GetVariableAsListFromVarFile Gets the string list represention of a variable from a provided input file found in VarFile
// Note that this returns a list of strings. For lists containing complex types, use GetAllVariablesFromVarFile.
func GetVariableAsListFromVarFile(t *testing.T, fileName string, key string) []string {
	result, err := GetVariableAsListFromVarFileE(t, fileName, key)
	require.NoError(t, err)

	return result
}

// GetVariableAsListFromVarFileE Gets the string list represention of a variable from a provided input file found in VarFile
// Note that this returns a list of strings. For lists containing complex types, use GetAllVariablesFromVarFile.
// Will return error if GetAllVariablesFromVarFileE returns an error, the key provided does not exist, or the value associated with the key is not a list
func GetVariableAsListFromVarFileE(t *testing.T, fileName string, key string) ([]string, error) {
	resultArray := []string{}

	var variables map[string]interface{}
	err := GetAllVariablesFromVarFileE(t, fileName, &variables)
	if err != nil {
		return nil, err
	}

	variable, exists := variables[key]
	if !exists {
		return nil, InputFileKeyNotFound{FilePath: fileName, Key: key}
	}

	if reflect.TypeOf(variable).String() != "[]interface {}" {
		return nil, UnexpectedOutputType{Key: key, ExpectedType: "[]interface {}", ActualType: reflect.TypeOf(variable).String()}
	}

	for _, item := range variable.([]interface{}) {
		resultArray = append(resultArray, fmt.Sprintf("%v", item))
	}

	return resultArray, nil
}

// GetAllVariablesFromVarFile Parses all data from a provided input file found in VarFile and stores the result in the value pointed to by out
func GetAllVariablesFromVarFile(t *testing.T, fileName string, out interface{}) {
	err := GetAllVariablesFromVarFileE(t, fileName, out)
	require.NoError(t, err)
}

// GetAllVariablesFromVarFileE Parses all data from a provided input file found ind in VarFile and stores the result in
// the value pointed to by out. Returns an error if the specified file does not exist, the specified file is not
// readable, or the specified file cannot be decoded from HCL
func GetAllVariablesFromVarFileE(t *testing.T, fileName string, out interface{}) error {
	fileContents, err := ioutil.ReadFile(fileName)
	if err != nil {
		return err
	}

	if err := parseAndDecodeHcl(string(fileContents), fileName, out); err != nil {
		return err
	}
	return nil
}

// parseAndDecodeHcl uses the HCL2 parser to parse the given string into an HCL file body, and then decode it into the
// given interface{}.
func parseAndDecodeHcl(hclContents string, filename string, out interface{}) (err error) {
	// The HCL2 parser and especially cty conversions will panic in many types of errors, so we have to recover from
	// those panics here and convert them to normal errors
	defer func() {
		if recovered := recover(); recovered != nil {
			err = PanicWhileParsingVarFile{RecoveredValue: recovered, ConfigFile: filename}
		}
	}()

	parser := hclparse.NewParser()

	var file *hcl.File
	if filepath.Ext(filename) == ".json" {
		rawFile, parseDiagnostics := parser.ParseJSON([]byte(hclContents), filename)
		if parseDiagnostics != nil && parseDiagnostics.HasErrors() {
			return parseDiagnostics
		}
		file = rawFile
	} else {
		rawFile, parseDiagnostics := parser.ParseHCL([]byte(hclContents), filename)
		if parseDiagnostics != nil && parseDiagnostics.HasErrors() {
			return parseDiagnostics
		}
		file = rawFile
	}

	decodeDiagnostics := gohcl.DecodeBody(file.Body, nil, out)
	if decodeDiagnostics != nil && decodeDiagnostics.HasErrors() {
		return decodeDiagnostics
	}
	return nil
}
