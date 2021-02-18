package terraform

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"reflect"
	"testing"

	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/stretchr/testify/require"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/gocty"
	ctyjson "github.com/zclconf/go-cty/cty/json"
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
	variables, err := GetAllVariablesFromVarFileE(t, fileName)
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
	variables, err := GetAllVariablesFromVarFileE(t, fileName)
	if err != nil {
		return nil, err
	}

	variable, exists := variables[key]

	if !exists {
		return nil, InputFileKeyNotFound{FilePath: fileName, Key: key}
	}

	if reflect.TypeOf(variable).String() != "map[string]interface {}" {
		return nil, UnexpectedOutputType{Key: key, ExpectedType: "map[string]interface {}", ActualType: reflect.TypeOf(variable).String()}
	}

	resultMap := make(map[string]string)
	for mapKey, mapVal := range variable.(map[string]interface{}) {
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
	variables, err := GetAllVariablesFromVarFileE(t, fileName)
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

	resultArray := []string{}
	for _, item := range variable.([]interface{}) {
		resultArray = append(resultArray, fmt.Sprintf("%v", item))
	}

	return resultArray, nil
}

// GetAllVariablesFromVarFileE Parses all data from a provided input file found in VarFile and returns a map of the
// variable values where the keys are the var names and the value are the parsed var values.
func GetAllVariablesFromVarFile(t *testing.T, fileName string) map[string]interface{} {
	out, err := GetAllVariablesFromVarFileE(t, fileName)
	require.NoError(t, err)
	return out
}

// GetAllVariablesFromVarFileE Parses all data from a provided input file found in VarFile and returns a map of the
// variable values where the keys are the var names and the value are the parsed var values. Returns an error if the
// specified file does not exist, the specified file is not readable, or the specified file cannot be decoded from HCL
func GetAllVariablesFromVarFileE(t *testing.T, fileName string) (map[string]interface{}, error) {
	fileContents, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	out, err := parseAndDecodeVarFile(string(fileContents), fileName)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// parseAndDecodeHcl uses the HCL2 parser to parse the given varfile string into an HCL file body, and then decode it
// into a map that maps var names to values.
func parseAndDecodeVarFile(hclContents string, filename string) (out map[string]interface{}, err error) {
	// The HCL2 parser and especially cty conversions will panic in many types of errors, so we have to recover from
	// those panics here and convert them to normal errors
	defer func() {
		if recovered := recover(); recovered != nil {
			err = PanicWhileParsingVarFile{RecoveredValue: recovered, ConfigFile: filename}
		}
	}()

	parser := hclparse.NewParser()

	file, parseDiagnostics := parser.ParseHCL([]byte(hclContents), filename)
	if parseDiagnostics != nil && parseDiagnostics.HasErrors() {
		return nil, parseDiagnostics
	}

	// VarFiles should only have attributes, so extract the attributes and decode the expressions into the return map.
	attrs, hclDiags := file.Body.JustAttributes()
	if hclDiags != nil && hclDiags.HasErrors() {
		return nil, hclDiags
	}

	valMap := map[string]cty.Value{}
	for name, attr := range attrs {
		val, hclDiags := attr.Expr.Value(nil) // nil because no function calls or variable references are allowed here
		if hclDiags != nil && hclDiags.HasErrors() {
			return nil, hclDiags
		}
		valMap[name] = val
	}
	varMap, err := parseCtyValueMapToGoMap(valMap)
	if err != nil {
		return nil, err
	}
	return varMap, nil
}

// The following functions were lifted from terragrunt -
// https://github.com/gruntwork-io/terragrunt/blob/874d2db4618afc7dc8e153ab0219282147865d4a/config/cty_helpers.go

// This is a hacky workaround to convert a cty Value map to a Go map[string]interface{}. cty does not support this directly
// (https://github.com/hashicorp/hcl2/issues/108) and doing it with gocty.FromCtyValue is nearly impossible, as cty
// requires you to specify all the output types and will error out when it hits interface{}. So, as an ugly workaround,
// we convert the given value to JSON using cty's JSON library and then convert the JSON back to a
// map[string]interface{} using the Go json library.
func parseCtyValueMapToGoMap(valMap map[string]cty.Value) (map[string]interface{}, error) {
	value, err := convertValuesMapToCtyVal(valMap)
	if err != nil {
		return nil, err
	}

	jsonBytes, err := ctyjson.Marshal(value, cty.DynamicPseudoType)
	if err != nil {
		return nil, err
	}

	var ctyJsonOutput CtyJsonOutput
	if err := json.Unmarshal(jsonBytes, &ctyJsonOutput); err != nil {
		return nil, err
	}

	return ctyJsonOutput.Value, nil
}

// When you convert a cty value to JSON, if any of that types are not yet known (i.e., are labeled as
// DynamicPseudoType), cty's Marshall method will write the type information to a type field and the actual value to
// a value field. This struct is used to capture that information so when we parse the JSON back into a Go struct, we
// can pull out just the Value field we need.
type CtyJsonOutput struct {
	Value map[string]interface{}
	Type  interface{}
}

// convertValuesMapToCtyVal takes a map of name - cty.Value pairs and converts to a single cty.Value object so that it
// can be fed through the cty json parser.
func convertValuesMapToCtyVal(valMap map[string]cty.Value) (cty.Value, error) {
	valMapAsCty := cty.NilVal
	if valMap != nil && len(valMap) > 0 {
		var err error
		valMapAsCty, err = gocty.ToCtyValue(valMap, generateTypeFromValuesMap(valMap))
		if err != nil {
			return valMapAsCty, err
		}
	}
	return valMapAsCty, nil
}

// generateTypeFromValuesMap takes a values map and returns an object type that has the same number of fields, but
// bound to each type of the underlying evaluated expression. This is the only way the HCL decoder will be happy, as
// object type is the only map type that allows different types for each attribute (cty.Map requires all attributes to
// have the same type.
func generateTypeFromValuesMap(valMap map[string]cty.Value) cty.Type {
	outType := map[string]cty.Type{}
	for k, v := range valMap {
		outType[k] = v.Type()
	}
	return cty.Object(outType)
}
