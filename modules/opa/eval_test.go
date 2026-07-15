package opa_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/gruntwork-io/terratest/modules/opa/v2"
)

func TestFormatOPAEvalArgs(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		options  *opa.EvalOptions
		rulePath string
		jsonFile string
		query    string
		expected []string
	}{
		{
			name: "Basic args without extras",
			options: &opa.EvalOptions{
				FailMode: opa.NoFail,
			},
			rulePath: "/path/to/policy.rego",
			jsonFile: "/path/to/input.json",
			query:    "data.test.allow",
			expected: []string{"eval", "-i", "/path/to/input.json", "-d", "/path/to/policy.rego", "data.test.allow"},
		},
		{
			name: "With fail mode",
			options: &opa.EvalOptions{
				FailMode: opa.FailUndefined,
			},
			rulePath: "/path/to/policy.rego",
			jsonFile: "/path/to/input.json",
			query:    "data.test.allow",
			expected: []string{"eval", "--fail", "-i", "/path/to/input.json", "-d", "/path/to/policy.rego", "data.test.allow"},
		},
		{
			name: "With extra args",
			options: &opa.EvalOptions{
				FailMode:  opa.FailUndefined,
				ExtraArgs: []string{"--format", "json"},
			},
			rulePath: "/path/to/policy.rego",
			jsonFile: "/path/to/input.json",
			query:    "data.test.allow",
			expected: []string{"eval", "--format", "json", "--fail", "-i", "/path/to/input.json", "-d", "/path/to/policy.rego", "data.test.allow"},
		},
		{
			name: "With v0-compatible flag",
			options: &opa.EvalOptions{
				FailMode:  opa.FailUndefined,
				ExtraArgs: []string{"--v0-compatible"},
			},
			rulePath: "/path/to/policy.rego",
			jsonFile: "/path/to/input.json",
			query:    "data.test.allow",
			expected: []string{"eval", "--v0-compatible", "--fail", "-i", "/path/to/input.json", "-d", "/path/to/policy.rego", "data.test.allow"},
		},
		{
			name: "With multiple extra args",
			options: &opa.EvalOptions{
				FailMode:  opa.FailUndefined,
				ExtraArgs: []string{"--v0-compatible", "--format", "json"},
			},
			rulePath: "/path/to/policy.rego",
			jsonFile: "/path/to/input.json",
			query:    "data.test.allow",
			expected: []string{"eval", "--v0-compatible", "--format", "json", "--fail", "-i", "/path/to/input.json", "-d", "/path/to/policy.rego", "data.test.allow"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			actual := opa.FormatOPAEvalArgs(test.options, test.rulePath, test.jsonFile, test.query)
			assert.Equal(t, test.expected, actual)
		})
	}
}

func TestEvalWithOutput(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
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
				allow := true if {
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
				allow := true if {
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
		t.Helper()

		f, err := os.CreateTemp(t.TempDir(), name)
		require.NoError(t, err)
		t.Cleanup(func() { os.Remove(f.Name()) })

		_, err = f.WriteString(content)
		require.NoError(t, err)

		return f.Name()
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			policy := createTempFile(t, "policy-*.rego", test.policy)

			inputs := make([]string, len(test.inputs))

			for i, input := range test.inputs {
				f := createTempFile(t, "inputs-*.json", input)
				inputs[i] = f
			}

			options := &opa.EvalOptions{
				RulePath: policy,
			}

			outputs, err := opa.EvalWithOutputE(t, options, inputs, test.query)
			if test.isError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			for i, output := range test.outputs {
				require.JSONEq(t, output, outputs[i], "output for input: %d", i)
			}
		})
	}
}
