package k8s_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/gruntwork-io/terratest/modules/k8s/v2"
)

func TestUnmarshalJSONPath(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		expectedOut any
		name        string
		jsonBlob    string
		jsonPath    string
	}{
		{
			name:        "boolField",
			jsonBlob:    `{"key": true}`,
			jsonPath:    "{ .key }",
			expectedOut: []bool{true},
		},
		{
			name:     "nestedObject",
			jsonBlob: `{"key": {"data": [1,2,3]}}`,
			jsonPath: "{ .key }",
			expectedOut: []map[string][]int{
				{
					"data": {1, 2, 3},
				},
			},
		},
		{
			name:        "nestedArray",
			jsonBlob:    `{"key": {"data": [1,2,3]}}`,
			jsonPath:    "{ .key.data[*] }",
			expectedOut: []int{1, 2, 3},
		},
	}

	for _, testCase := range testCases {
		// capture range variable so that it doesn't update when the subtest goroutine swaps.
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			var output any
			k8s.UnmarshalJSONPath(t, []byte(testCase.jsonBlob), testCase.jsonPath, &output)
			// NOTE: we have to do equality check on the marshalled json data to allow equality checks over dynamic
			// types in this table driven test.
			expectedOutJSON, err := json.Marshal(testCase.expectedOut)
			require.NoError(t, err)
			actualOutJSON, err := json.Marshal(output)
			require.NoError(t, err)
			assert.JSONEq(t, string(expectedOutJSON), string(actualOutJSON))
		})
	}
}
