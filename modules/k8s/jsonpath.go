package k8s

import (
	"bytes"
	"encoding/json"

	"k8s.io/client-go/util/jsonpath"

	"github.com/gruntwork-io/terratest/modules/testing/v2"
	"github.com/stretchr/testify/require"
)

// UnmarshalJSONPath allows you to use an arbitrary JSONPath string to query a json blob and unmarshal the resulting
// output into a go object. Note that the output will always be a list. That means that if you query a single object,
// the output will be a list of single element, not the element itself. However, if the json path maps to a list, then
// the output will be that list.
// Example:
//
// jsonBlob := []byte(`{"key": {"data": [1,2,3]}}`)
// jsonPath := "{.key.data[*]}"
// var output []int
// UnmarshalJSONPath(t, jsonBlob, jsonPath, &output)
// // output is []int{1,2,3}
//
// This will fail the test if there is an error.
func UnmarshalJSONPath(t testing.TestingT, jsonData []byte, jsonpathStr string, output interface{}) {
	err := UnmarshalJSONPathE(t, jsonData, jsonpathStr, output)
	require.NoError(t, err)
}

// UnmarshalJSONPathE allows you to use an arbitrary JSONPath string to query a json blob and unmarshal the resulting
// output into a go object. Note that the output will always be a list. That means that if you query a single object,
// the output will be a list of single element, not the element itself. However, if the json path maps to a list, then
// the output will be that list.
// Example:
//
// jsonBlob := []byte(`{"key": {"data": [1,2,3]}}`)
// jsonPath := "{.key.data[*]}"
// var output []int
// UnmarshalJSONPathE(t, jsonBlob, jsonPath, &output)
// => output = []int{1,2,3}
func UnmarshalJSONPathE(t testing.TestingT, jsonData []byte, jsonpathStr string, output interface{}) error {
	// First, unmarshal the full json object. We use interface{} to avoid the type conversions, as jsonpath will handle
	// it for us.
	var blob interface{}
	if err := json.Unmarshal(jsonData, &blob); err != nil {
		return JSONPathMalformedJSONErr{err}
	}

	// Then, query the json object with the given jsonpath to get the output string.
	jsonpathParser := jsonpath.New(t.Name())
	jsonpathParser.EnableJSONOutput(true)
	if err := jsonpathParser.Parse(jsonpathStr); err != nil {
		return JSONPathMalformedJSONPathErr{err}
	}
	outputJSONBuffer := new(bytes.Buffer)
	if err := jsonpathParser.Execute(outputJSONBuffer, blob); err != nil {
		return JSONPathExtractJSONPathErr{err}
	}
	outputJSON := outputJSONBuffer.Bytes()

	// Finally, we need to unmarshal the output object into the given output var.
	if err := json.Unmarshal(outputJSON, output); err != nil {
		return JSONPathMalformedJSONPathResultErr{err}
	}
	return nil
}
