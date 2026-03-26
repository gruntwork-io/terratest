package terraform_test

import (
	"fmt"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestFormatTerraformPlanFileAsArgs(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		command  string
		out      string
		expected []string
	}{
		{command: "plan", out: "/some/plan/output", expected: []string{"-out=/some/plan/output"}},
		{command: "plan", out: "", expected: nil},
		{command: "apply", out: "/some/plan/output", expected: []string{"/some/plan/output"}},
		{command: "apply", out: "", expected: nil},
		{command: "show", out: "/some/plan/output", expected: []string{"/some/plan/output"}},
		{command: "show", out: "", expected: nil},
	}

	for _, testCase := range testCases {
		checkResultWithRetry(t, 100, testCase.expected, fmt.Sprintf("FormatTerraformPlanFileAsArgs(%v)", testCase.out), func() any {
			return terraform.FormatTerraformPlanFileAsArg(testCase.command, testCase.out)
		})
	}
}

func TestFormatTerraformVarsAsArgs(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		vars     map[string]any
		expected []string
	}{
		{vars: map[string]any{}, expected: nil},
		{vars: map[string]any{"foo": "bar"}, expected: []string{"-var", "foo=bar"}},
		{vars: map[string]any{"foo": 123}, expected: []string{"-var", "foo=123"}},
		{vars: map[string]any{"foo": true}, expected: []string{"-var", "foo=true"}},
		{vars: map[string]any{"foo": nil}, expected: []string{"-var", "foo=null"}},
		{vars: map[string]any{"foo": []int{1, 2, 3}}, expected: []string{"-var", "foo=[1, 2, 3]"}},
		{vars: map[string]any{"foo": map[string]string{"baz": "blah"}}, expected: []string{"-var", "foo={\"baz\" = \"blah\"}"}},
		{
			vars:     map[string]any{"str": "bar", "int": -1, "bool": false, "list": []string{"foo", "bar", "baz"}, "map": map[string]int{"foo": 0}},
			expected: []string{"-var", "str=bar", "-var", "int=-1", "-var", "bool=false", "-var", "list=[\"foo\", \"bar\", \"baz\"]", "-var", "map={\"foo\" = 0}"},
		},
	}

	for _, testCase := range testCases {
		checkResultWithRetry(t, 100, testCase.expected, fmt.Sprintf("FormatTerraformVarsAsArgs(%v)", testCase.vars), func() any {
			return terraform.FormatTerraformVarsAsArgs(testCase.vars)
		})
	}
}

// Some of our tests execute code that loops over a map to produce output. The problem is that the order of map
// iteration is generally unpredictable and, to make it even more unpredictable, Go intentionally randomizes the
// iteration order (https://blog.golang.org/go-maps-in-action#TOC_7). Therefore, the order of items in the output
// is unpredictable, and doing a simple assert.Equals call will intermittently fail.
//
// We have a few unsatisfactory ways to solve this problem:
//
//  1. Enforce iteration order. This is easy to do in other languages, where you have built-in sorted maps. In Go, no
//     such map exists, and if you create a custom one, you can't use the range keyword on it
//     (http://stackoverflow.com/a/35810932/483528). As a result, we'd have to modify our implementation code to take
//     iteration order into account which is a totally unnecessary feature that increases complexity.
//  2. We could parse the output string and do an order-independent comparison. However, that adds a bunch of parsing
//     logic into the test code which is a totally unnecessary feature that increases complexity.
//  3. We accept that Go is a shitty language and, if the test fails, we re-run it a bunch of times in the hope that, if
//     the bug is caused by key ordering, we will randomly get the proper order in a future run. The code being tested
//     here is tiny & fast, so doing a hundred retries is still sub millisecond, so while ugly, this provides a very
//     simple solution.
//
// Isn't it great that Go's designers built features into the language to prevent bugs that now force every Go
// developer to write thousands of lines of extra code like this, which is of course likely to contain bugs itself?
func checkResultWithRetry(t *testing.T, maxRetries int, expectedValue any, description string, generateValue func() any) {
	t.Helper()

	for i := 0; i < maxRetries; i++ {
		actualValue := generateValue()
		if assert.ObjectsAreEqual(expectedValue, actualValue) {
			return
		}

		t.Logf("Retry %d of %s failed: expected %v, got %v", i, description, expectedValue, actualValue)
	}

	assert.Fail(t, "checkResultWithRetry failed", "After %d retries, %s still not succeeding (see retries above)", description)
}

func TestFormatArgsAppliesLockCorrectly(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		command  []string
		expected []string
	}{
		{command: []string{"plan"}, expected: []string{"plan", "-lock=false"}},
		{command: []string{"validate"}, expected: []string{"validate"}},
		{command: []string{"validate", "--all"}, expected: []string{"validate", "--all"}},
		{command: []string{"plan", "--all"}, expected: []string{"plan", "--all", "-lock=false"}},
	}

	for _, testCase := range testCases {
		assert.Equal(t, testCase.expected, terraform.FormatArgs(&terraform.Options{}, testCase.command...))
	}
}

func TestFormatSetVarsAfterVarFilesFormatsCorrectly(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		vars                 map[string]any
		command              []string
		varFiles             []string
		expected             []string
		setVarsAfterVarFiles bool
	}{
		{command: []string{"plan"}, vars: map[string]any{"foo": "bar"}, varFiles: []string{"test.tfvars"}, setVarsAfterVarFiles: true, expected: []string{"plan", "-var-file", "test.tfvars", "-var", "foo=bar", "-lock=false"}},
		{command: []string{"plan"}, vars: map[string]any{"foo": "bar", "hello": "world"}, varFiles: []string{"test.tfvars"}, setVarsAfterVarFiles: true, expected: []string{"plan", "-var-file", "test.tfvars", "-var", "foo=bar", "-var", "hello=world", "-lock=false"}},
		{command: []string{"plan"}, vars: map[string]any{"foo": "bar", "hello": "world"}, varFiles: []string{"test.tfvars"}, setVarsAfterVarFiles: false, expected: []string{"plan", "-var", "foo=bar", "-var", "hello=world", "-var-file", "test.tfvars", "-lock=false"}},
		{command: []string{"plan"}, vars: map[string]any{"foo": "bar"}, varFiles: []string{"test.tfvars"}, setVarsAfterVarFiles: false, expected: []string{"plan", "-var", "foo=bar", "-var-file", "test.tfvars", "-lock=false"}},
	}

	for _, testCase := range testCases {
		result := terraform.FormatArgs(&terraform.Options{SetVarsAfterVarFiles: testCase.setVarsAfterVarFiles, Vars: testCase.vars, VarFiles: testCase.varFiles}, testCase.command...)

		// Make sure that -var and -var-file options are in the expected order relative to each other
		// Note that the order of the different -var and -var-file options may change
		// See this comment for more info: https://github.com/gruntwork-io/terratest/blob/6fb86056797e3e62ebdd9011ba26605e0976a6f8/modules/terraform/format_test.go#L123-L142
		for idx, arg := range result {
			if arg == "-var-file" || arg == "-var" {
				assert.Equal(t, testCase.expected[idx], arg)
			}
		}

		// Make sure that the order of other arguments hasn't been incorrectly modified
		assert.Equal(t, testCase.expected[0], result[0])
		assert.Equal(t, testCase.expected[len(testCase.expected)-1], result[len(result)-1])
	}
}

func TestMixedVars(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		vars                 map[string]any
		command              []string
		mixedVars            []terraform.Var
		varFiles             []string
		expected             []string
		setVarsAfterVarFiles bool
	}{
		{command: []string{"plan"}, mixedVars: []terraform.Var{terraform.VarFile("/path1"), terraform.VarInline("name", "value"), terraform.VarFile("/path2")}, vars: map[string]any{"foo": "bar"}, varFiles: []string{"test.tfvars"}, setVarsAfterVarFiles: true, expected: []string{"plan", "-var-file", "/path1", "-var", "name=value", "-var-file", "/path2", "-var-file", "test.tfvars", "-var", "foo=bar", "-lock=false"}},
		{command: []string{"plan"}, mixedVars: []terraform.Var{terraform.VarInline("name1", "value"), terraform.VarInline("name2", "value"), terraform.VarFile("/path")}, vars: map[string]any{"foo": "bar", "hello": "world"}, varFiles: []string{"test.tfvars"}, setVarsAfterVarFiles: true, expected: []string{"plan", "-var", "name1=value", "-var", "name2=value", "-var-file", "/path", "-var-file", "test.tfvars", "-var", "foo=bar", "-var", "hello=world", "-lock=false"}},
		{command: []string{"plan"}, mixedVars: []terraform.Var{terraform.VarFile("/path"), terraform.VarInline("name1", "value"), terraform.VarInline("name2", "value")}, vars: map[string]any{"foo": "bar", "hello": "world"}, varFiles: []string{"test.tfvars"}, setVarsAfterVarFiles: false, expected: []string{"plan", "-var-file", "path", "-var", "name1=value", "-var", "name2=value", "-var", "foo=bar", "-var", "hello=world", "-var-file", "test.tfvars", "-lock=false"}},
		{command: []string{"plan"}, mixedVars: []terraform.Var{terraform.VarFile("/path"), terraform.VarInline("name", "value")}, vars: map[string]any{"foo": "bar"}, varFiles: []string{"test.tfvars"}, setVarsAfterVarFiles: false, expected: []string{"plan", "-var-file", "/path", "-var", "name=value", "-var", "foo=bar", "-var-file", "test.tfvars", "-lock=false"}},
	}

	for _, testCase := range testCases {
		result := terraform.FormatArgs(&terraform.Options{SetVarsAfterVarFiles: testCase.setVarsAfterVarFiles, Vars: testCase.vars, VarFiles: testCase.varFiles, MixedVars: testCase.mixedVars}, testCase.command...)

		// Make sure that var defined in `MixedVars` are seriliazed in order and precede `Var`` and `VarFiles``
		// Make sure that -var and -var-file options are in the expected order relative to each other
		// Note that the order of the different -var and -var-file options may change
		// See this comment for more info: https://github.com/gruntwork-io/terratest/blob/6fb86056797e3e62ebdd9011ba26605e0976a6f8/modules/terraform/format_test.go#L123-L142
		for idx, arg := range result {
			if arg == "-var-file" || arg == "-var" {
				assert.Equal(t, testCase.expected[idx], arg)
			}
		}

		// Make sure that the order of other arguments hasn't been incorrectly modified
		assert.Equal(t, testCase.expected[0], result[0])
		assert.Equal(t, testCase.expected[len(testCase.expected)-1], result[len(result)-1])
	}
}
