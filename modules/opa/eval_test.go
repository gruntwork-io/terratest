package opa

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
