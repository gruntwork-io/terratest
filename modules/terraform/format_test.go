package terraform

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatTerraformVarsAsArgs(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		vars     map[string]interface{}
		expected []string
	}{
		{map[string]interface{}{}, nil},
		{map[string]interface{}{"foo": "bar"}, []string{"-var", "foo=bar"}},
		{map[string]interface{}{"foo": 123}, []string{"-var", "foo=123"}},
		{map[string]interface{}{"foo": true}, []string{"-var", "foo=true"}},
		{map[string]interface{}{"foo": nil}, []string{"-var", "foo=null"}},
		{map[string]interface{}{"foo": []int{1, 2, 3}}, []string{"-var", "foo=[1, 2, 3]"}},
		{map[string]interface{}{"foo": map[string]string{"baz": "blah"}}, []string{"-var", "foo={\"baz\" = \"blah\"}"}},
		{
			map[string]interface{}{"str": "bar", "int": -1, "bool": false, "list": []string{"foo", "bar", "baz"}, "map": map[string]int{"foo": 0}},
			[]string{"-var", "str=bar", "-var", "int=-1", "-var", "bool=false", "-var", "list=[\"foo\", \"bar\", \"baz\"]", "-var", "map={\"foo\" = 0}"},
		},
	}

	for _, testCase := range testCases {
		checkResultWithRetry(t, 100, testCase.expected, fmt.Sprintf("FormatTerraformVarsAsArgs(%v)", testCase.vars), func() interface{} {
			return FormatTerraformVarsAsArgs(testCase.vars)
		})
	}
}

func TestPrimitiveToHclString(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		value    interface{}
		expected string
	}{
		{"", ""},
		{"foo", "foo"},
		{"true", "true"},
		{true, "true"},
		{3, "3"},
		{nil, "null"},
		{[]int{1, 2, 3}, "[1 2 3]"}, // Anything that isn't a primitive is forced into a string
	}

	for _, testCase := range testCases {
		actual := primitiveToHclString(testCase.value, false)
		assert.Equal(t, testCase.expected, actual, "Value: %v", testCase.value)
	}
}

func TestMapToHclString(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		value    map[string]interface{}
		expected string
	}{
		{map[string]interface{}{}, "{}"},
		{map[string]interface{}{"key1": "value1"}, "{\"key1\" = \"value1\"}"},
		{map[string]interface{}{"key1": 123}, "{\"key1\" = 123}"},
		{map[string]interface{}{"key1": true}, "{\"key1\" = true}"},
		{map[string]interface{}{"key1": []int{1, 2, 3}}, "{\"key1\" = [1, 2, 3]}"}, // Any value that isn't a primitive is forced into a string
		{map[string]interface{}{"key1": "value1", "key2": 0, "key3": false}, "{\"key1\" = \"value1\", \"key2\" = 0, \"key3\" = false}"},
		{map[string]interface{}{"key1.a.b.c": "value1"}, "{\"key1.a.b.c\" = \"value1\"}"},
	}

	for _, testCase := range testCases {
		checkResultWithRetry(t, 100, testCase.expected, fmt.Sprintf("mapToHclString(%v)", testCase.value), func() interface{} {
			return mapToHclString(testCase.value)
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
// 1. Enforce iteration order. This is easy to do in other languages, where you have built-in sorted maps. In Go, no
//    such map exists, and if you create a custom one, you can't use the range keyword on it
//    (http://stackoverflow.com/a/35810932/483528). As a result, we'd have to modify our implementation code to take
//    iteration order into account which is a totally unnecessary feature that increases complexity.
// 2. We could parse the output string and do an order-independent comparison. However, that adds a bunch of parsing
//    logic into the test code which is a totally unnecessary feature that increases complexity.
// 3. We accept that Go is a shitty language and, if the test fails, we re-run it a bunch of times in the hope that, if
//    the bug is caused by key ordering, we will randomly get the proper order in a future run. The code being tested
//    here is tiny & fast, so doing a hundred retries is still sub millisecond, so while ugly, this provides a very
//    simple solution.
//
// Isn't it great that Go's designers built features into the language to prevent bugs that now force every Go
// developer to write thousands of lines of extra code like this, which is of course likely to contain bugs itself?
func checkResultWithRetry(t *testing.T, maxRetries int, expectedValue interface{}, description string, generateValue func() interface{}) {
	for i := 0; i < maxRetries; i++ {
		actualValue := generateValue()
		if assert.ObjectsAreEqual(expectedValue, actualValue) {
			return
		}
		t.Logf("Retry %d of %s failed: expected %v, got %v", i, description, expectedValue, actualValue)
	}

	assert.Fail(t, "checkResultWithRetry failed", "After %d retries, %s still not succeeding (see retries above)", description)
}

func TestSliceToHclString(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		value    []interface{}
		expected string
	}{
		{[]interface{}{}, "[]"},
		{[]interface{}{"foo"}, "[\"foo\"]"},
		{[]interface{}{123}, "[123]"},
		{[]interface{}{true}, "[true]"},
		{[]interface{}{[]int{1, 2, 3}}, "[[1, 2, 3]]"}, // Any value that isn't a primitive is forced into a string
		{[]interface{}{"foo", 0, false}, "[\"foo\", 0, false]"},
		{[]interface{}{map[string]interface{}{"foo": "bar"}}, "[{\"foo\" = \"bar\"}]"},
		{[]interface{}{map[string]interface{}{"foo": "bar"}, map[string]interface{}{"foo": "bar"}}, "[{\"foo\" = \"bar\"}, {\"foo\" = \"bar\"}]"},
	}

	for _, testCase := range testCases {
		actual := sliceToHclString(testCase.value)
		assert.Equal(t, testCase.expected, actual, "Value: %v", testCase.value)
	}
}

func TestToHclString(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		value    interface{}
		expected string
	}{
		{"", ""},
		{"foo", "foo"},
		{123, "123"},
		{true, "true"},
		{[]int{1, 2, 3}, "[1, 2, 3]"},
		{[]string{"foo", "bar", "baz"}, "[\"foo\", \"bar\", \"baz\"]"},
		{map[string]string{"key1": "value1"}, "{\"key1\" = \"value1\"}"},
		{map[string]int{"key1": 123}, "{\"key1\" = 123}"},
	}

	for _, testCase := range testCases {
		actual := toHclString(testCase.value, false)
		assert.Equal(t, testCase.expected, actual, "Value: %v", testCase.value)
	}
}

func TestTryToConvertToGenericSlice(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		value           interface{}
		expectedSlice   []interface{}
		expectedIsSlice bool
	}{
		{"", []interface{}{}, false},
		{"foo", []interface{}{}, false},
		{true, []interface{}{}, false},
		{531, []interface{}{}, false},
		{map[string]string{"foo": "bar"}, []interface{}{}, false},
		{[]string{}, []interface{}{}, true},
		{[]int{}, []interface{}{}, true},
		{[]bool{}, []interface{}{}, true},
		{[]interface{}{}, []interface{}{}, true},
		{[]string{"foo", "bar", "baz"}, []interface{}{"foo", "bar", "baz"}, true},
		{[]int{1, 2, 3}, []interface{}{1, 2, 3}, true},
		{[]bool{true, true, false}, []interface{}{true, true, false}, true},
		{[]interface{}{"foo", "bar", "baz"}, []interface{}{"foo", "bar", "baz"}, true},
	}

	for _, testCase := range testCases {
		actualSlice, actualIsSlice := tryToConvertToGenericSlice(testCase.value)
		assert.Equal(t, testCase.expectedSlice, actualSlice, "Value: %v", testCase.value)
		assert.Equal(t, testCase.expectedIsSlice, actualIsSlice, "Value: %v", testCase.value)
	}
}

func TestTryToConvertToGenericMap(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		value         interface{}
		expectedMap   map[string]interface{}
		expectedIsMap bool
	}{
		{"", map[string]interface{}{}, false},
		{"foo", map[string]interface{}{}, false},
		{true, map[string]interface{}{}, false},
		{531, map[string]interface{}{}, false},
		{[]string{"foo", "bar"}, map[string]interface{}{}, false},
		{map[int]int{}, map[string]interface{}{}, false},
		{map[bool]string{}, map[string]interface{}{}, false},
		{map[string]string{}, map[string]interface{}{}, true},
		{map[string]int{}, map[string]interface{}{}, true},
		{map[string]bool{}, map[string]interface{}{}, true},
		{map[string]interface{}{}, map[string]interface{}{}, true},
		{map[string]string{"key1": "value1", "key2": "value2"}, map[string]interface{}{"key1": "value1", "key2": "value2"}, true},
		{map[string]int{"key1": 1, "key2": 2, "key3": 3}, map[string]interface{}{"key1": 1, "key2": 2, "key3": 3}, true},
		{map[string]bool{"key1": true}, map[string]interface{}{"key1": true}, true},
		{map[string]interface{}{"key1": "value1"}, map[string]interface{}{"key1": "value1"}, true},
	}

	for _, testCase := range testCases {
		actualMap, actualIsMap := tryToConvertToGenericMap(testCase.value)
		assert.Equal(t, testCase.expectedMap, actualMap, "Value: %v", testCase.value)
		assert.Equal(t, testCase.expectedIsMap, actualIsMap, "Value: %v", testCase.value)
	}
}
