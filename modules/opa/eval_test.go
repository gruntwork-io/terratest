package opa

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFormatOPAEvalArgs(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		options  *EvalOptions
		rulePath string
		jsonFile string
		query    string
		expected []string
	}{
		{
			name: "Basic args without extras",
			options: &EvalOptions{
				FailMode: NoFail,
			},
			rulePath: "/path/to/policy.rego",
			jsonFile: "/path/to/input.json",
			query:    "data.test.allow",
			expected: []string{"eval", "-i", "/path/to/input.json", "-d", "/path/to/policy.rego", "data.test.allow"},
		},
		{
			name: "With fail mode",
			options: &EvalOptions{
				FailMode: FailUndefined,
			},
			rulePath: "/path/to/policy.rego",
			jsonFile: "/path/to/input.json",
			query:    "data.test.allow",
			expected: []string{"eval", "--fail", "-i", "/path/to/input.json", "-d", "/path/to/policy.rego", "data.test.allow"},
		},
		{
			name: "With extra args",
			options: &EvalOptions{
				FailMode:  FailUndefined,
				ExtraArgs: []string{"--v0-compatible", "--format", "json"},
			},
			rulePath: "/path/to/policy.rego",
			jsonFile: "/path/to/input.json",
			query:    "data.test.allow",
			expected: []string{"eval", "--fail", "--v0-compatible", "--format", "json", "-i", "/path/to/input.json", "-d", "/path/to/policy.rego", "data.test.allow"},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			actual := formatOPAEvalArgs(test.options, test.rulePath, test.jsonFile, test.query)
			assert.Equal(t, test.expected, actual)
		})
	}
}

func TestEvalWithOutput(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string

		policy  string
		query   string
		inputs  []string
		outputs []string
		isError bool
	}{
		{
			name: "Success",
			policy: `
				package test
				allow {
					startswith(input.user, "admin")
				}
			`,
			query: "data.test.allow",
			inputs: []string{
				`{"user": "admin-1"}`,
				`{"user": "admin-2"}`,
				`{"user": "admin-3"}`,
			},
			outputs: []string{
				`{
					"result": [{
						"expressions": [{
							"value": true,
							"text": "data.test.allow",
							"location": {
								"row": 1,
								"col": 1
							}
						}]
					}]
				}`,
				`{
					"result": [{
						"expressions": [{
							"value": true,
							"text": "data.test.allow",
							"location": {
								"row": 1,
								"col": 1
							}
						}]
					}]
				}`,
				`{
					"result": [{
						"expressions": [{
							"value": true,
							"text": "data.test.allow",
							"location": {
								"row": 1,
								"col": 1
							}
						}]
					}]
				}`,
			},
		},
		{
			name: "ContainsError",
			policy: `
				package test
				allow {
					input.user == "admin"
				}
			`,
			query:   "data.test.allow",
			isError: true,
			inputs: []string{
				`{"user": "admin"}`,
				`{"user": "nobody"}`,
			},
			outputs: []string{
				`{
					"result": [{
						"expressions": [{
							"value": true,
							"text": "data.test.allow",
							"location": {
								"row": 1,
								"col": 1
							}
						}]
					}]
				}`,
				`{
					"result": [{
						"expressions": [{
							"value": {
								"test": {}
							},
							"text": "data",
							"location": {
								"row": 1,
								"col": 1
							}
						}]
					}]
				}`,
			},
		},
	}

	createTempFile := func(t *testing.T, name string, content string) string {
		f, err := os.CreateTemp(t.TempDir(), name)
		require.NoError(t, err)
		t.Cleanup(func() { os.Remove(f.Name()) })
		_, err = f.WriteString(content)
		require.NoError(t, err)
		return f.Name()
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			policy := createTempFile(t, "policy-*.rego", test.policy)
			inputs := make([]string, len(test.inputs))
			for i, input := range test.inputs {
				f := createTempFile(t, "inputs-*.json", input)
				inputs[i] = f
			}

			options := &EvalOptions{
				RulePath: policy,
			}

			outputs, err := EvalWithOutputE(t, options, inputs, test.query)
			if test.isError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			for i, output := range test.outputs {
				require.JSONEq(t, output, outputs[i], "output for input: %d", i)
			}
		})
	}
}
