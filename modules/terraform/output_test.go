package terraform_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gruntwork-io/terratest/modules/files"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/require"
)

func TestOutputString(t *testing.T) {
	t.Parallel()

	testFolder, err := files.CopyTerraformFolderToTemp("../../test/fixtures/terraform-output", t.Name())
	require.NoError(t, err)

	options := &terraform.Options{
		TerraformDir: testFolder,
	}

	terraform.InitAndApply(t, options)

	b := terraform.Output(t, options, "bool")
	require.Equal(t, "true", b, "Bool %q should match %q", "true", b)

	str := terraform.Output(t, options, "string")
	require.Equal(t, "This is a string.", str, "String %q should match %q", "This is a string.", str)

	num := terraform.Output(t, options, "number")
	require.Equal(t, "3.14", num, "Number %q should match %q", "3.14", num)

	num1 := terraform.Output(t, options, "number1")
	require.Equal(t, "3", num1, "Number %q should match %q", "3", num1)

	unicodeString := terraform.Output(t, options, "unicode_string")
	require.Equal(t, "söme chäräcter", unicodeString)
}

func TestOutputList(t *testing.T) {
	t.Parallel()

	testFolder, err := files.CopyTerraformFolderToTemp("../../test/fixtures/terraform-output-list", t.Name())
	require.NoError(t, err)

	options := &terraform.Options{
		TerraformDir: testFolder,
	}

	terraform.InitAndApply(t, options)
	out := terraform.OutputList(t, options, "giant_steps")

	expectedLen := 4
	expectedItem := "John Coltrane"
	expectedArray := []string{"John Coltrane", "Tommy Flanagan", "Paul Chambers", "Art Taylor"}

	require.Len(t, out, expectedLen, "Output should contain %d items", expectedLen)
	require.Contains(t, out, expectedItem, "Output should contain %q item", expectedItem)
	require.Equal(t, expectedItem, out[0], "First item should be %q, got %q", expectedItem, out[0])
	require.Equal(t, expectedArray, out, "Array %q should match %q", expectedArray, out)
}

func TestOutputNotListError(t *testing.T) {
	t.Parallel()

	testFolder, err := files.CopyTerraformFolderToTemp("../../test/fixtures/terraform-output-list", t.Name())
	if err != nil {
		t.Fatal(err)
	}

	options := &terraform.Options{
		TerraformDir: testFolder,
	}

	terraform.InitAndApply(t, options)

	_, err = terraform.OutputListE(t, options, "not_a_list")
	require.Error(t, err)
}

func TestOutputMap(t *testing.T) {
	t.Parallel()

	testFolder, err := files.CopyTerraformFolderToTemp("../../test/fixtures/terraform-output-map", t.Name())
	require.NoError(t, err)

	options := &terraform.Options{
		TerraformDir: testFolder,
	}

	terraform.InitAndApply(t, options)
	out := terraform.OutputMap(t, options, "mogwai")

	t.Log(out)

	expectedLen := 4
	expectedMap := map[string]string{
		"guitar_1": "Stuart Braithwaite",
		"guitar_2": "Barry Burns",
		"bass":     "Dominic Aitchison",
		"drums":    "Martin Bulloch",
	}

	require.Len(t, out, expectedLen, "Output should contain %d items", expectedLen)
	require.Equal(t, expectedMap, out, "Map %q should match %q", expectedMap, out)
}

func TestOutputNotMapError(t *testing.T) {
	t.Parallel()

	testFolder, err := files.CopyTerraformFolderToTemp("../../test/fixtures/terraform-output-map", t.Name())
	require.NoError(t, err)

	options := &terraform.Options{
		TerraformDir: testFolder,
	}

	terraform.InitAndApply(t, options)

	_, err = terraform.OutputMapE(t, options, "not_a_map")
	require.Error(t, err)
}

func TestOutputMapOfObjects(t *testing.T) {
	t.Parallel()

	testFolder, err := files.CopyTerraformFolderToTemp("../../test/fixtures/terraform-output-mapofobjects", t.Name())
	require.NoError(t, err)

	options := &terraform.Options{
		TerraformDir: testFolder,
	}

	terraform.InitAndApply(t, options)
	out := terraform.OutputMapOfObjects(t, options, "map_of_objects")

	nestedMap1 := map[string]any{
		"four": 4,
		"five": "five",
	}

	nestedList1 := []any{
		map[string]any{
			"six":   6,
			"seven": "seven",
		},
	}

	expectedMap1 := map[string]any{
		"somebool":  true,
		"somefloat": 1.1,
		"one":       1,
		"two":       "two",
		"three":     "three",
		"nest":      nestedMap1,
		"nest_list": nestedList1,
	}

	require.Equal(t, expectedMap1, out)
}

func TestOutputNotMapOfObjectsError(t *testing.T) {
	t.Parallel()

	testFolder, err := files.CopyTerraformFolderToTemp("../../test/fixtures/terraform-output-mapofobjects", t.Name())
	require.NoError(t, err)

	options := &terraform.Options{
		TerraformDir: testFolder,
	}

	terraform.InitAndApply(t, options)

	_, err = terraform.OutputMapOfObjectsE(t, options, "not_map_of_objects")
	require.Error(t, err)
}

func TestOutputListOfObjects(t *testing.T) {
	t.Parallel()

	testFolder, err := files.CopyTerraformFolderToTemp("../../test/fixtures/terraform-output-listofobjects", t.Name())
	require.NoError(t, err)

	options := &terraform.Options{
		TerraformDir: testFolder,
	}

	terraform.InitAndApply(t, options)
	out := terraform.OutputListOfObjects(t, options, "list_of_maps")

	expectedLen := 3
	nestedMap1 := map[string]any{
		"four": 4,
		"five": "five",
	}
	nestedList1 := []any{
		map[string]any{
			"four": 4,
			"five": "five",
		},
	}
	expectedMap1 := map[string]any{
		"one":   1,
		"two":   "two",
		"three": "three",
		"more":  nestedMap1,
	}

	expectedMap2 := map[string]any{
		"one":   "one",
		"two":   2,
		"three": 3,
		"more":  nestedList1,
	}

	expectedMap3 := map[string]any{
		"one":   "one",
		"two":   2,
		"three": 3,
		"more": []any{
			"one",
			2,
			3.4,
			[]any{"one", 2, 3.4},
			map[string]any{"one": 2, "three": 3.4},
		},
	}

	require.Len(t, out, expectedLen, "Output should contain %d items", expectedLen)
	assert.Equal(t, expectedMap1, out[0], "First map should be %q, got %q", expectedMap1, out[0])
	assert.Equal(t, expectedMap2, out[1], "Second map should be %q, got %q", expectedMap2, out[1])
	assert.Equal(t, expectedMap3, out[2], "Third map should be %q, got %q", expectedMap3, out[1])
}

func TestOutputNotListOfObjectsError(t *testing.T) {
	t.Parallel()

	testFolder, err := files.CopyTerraformFolderToTemp("../../test/fixtures/terraform-output-listofobjects", t.Name())
	require.NoError(t, err)

	options := &terraform.Options{
		TerraformDir: testFolder,
	}

	terraform.InitAndApply(t, options)

	_, err = terraform.OutputListOfObjectsE(t, options, "not_list_of_maps")
	require.Error(t, err)
}

func TestOutputsForKeys(t *testing.T) {
	t.Parallel()

	testFolder, err := files.CopyTerraformFolderToTemp("../../test/fixtures/terraform-output-all", t.Name())
	require.NoError(t, err)

	options := &terraform.Options{
		TerraformDir: testFolder,
	}

	keys := []string{"our_star", "stars", "magnitudes"}

	terraform.InitAndApply(t, options)
	out := terraform.OutputForKeys(t, options, keys)

	expectedLen := 3
	require.Len(t, out, expectedLen, "Output should contain %d items", expectedLen)

	// String value
	expectedString := "Sun"
	str, ok := out["our_star"].(string)
	require.Truef(t, ok, "Wrong data type for 'our_star', expected string, got %T", out["our_star"])
	require.Equal(t, expectedString, str, "String %q should match %q", expectedString, str)

	// List value
	expectedListLen := 3
	outputInterfaceList, ok := out["stars"].([]any)
	require.Truef(t, ok, "Wrong data type for 'stars', expected [], got %T", out["stars"])

	expectedListItem := "Sirius"

	require.Len(t, outputInterfaceList, expectedListLen, "Output list should contain %d items", expectedListLen)
	require.Equal(t, expectedListItem, outputInterfaceList[0].(string), "List Item %q should match %q",
		expectedListItem, outputInterfaceList[0].(string))

	// Map value
	outputInterfaceMap, ok := out["magnitudes"].(map[string]any)
	require.Truef(t, ok, "Wrong data type for 'magnitudes', expected map[string], got %T", out["magnitudes"])

	expectedMapLen := 3
	expectedMapItem := -1.46

	require.Len(t, outputInterfaceMap, expectedMapLen, "Output map should contain %d items", expectedMapLen)
	require.InEpsilon(t, expectedMapItem, outputInterfaceMap["Sirius"].(float64), 0.0001, "Map Item %q should match %q",
		expectedMapItem, outputInterfaceMap["Sirius"].(float64))

	// Key not in the parameter list
	outputNotPresentMap, ok := out["constellations"].(map[string]any)
	require.False(t, ok)
	require.Nil(t, outputNotPresentMap)
}

func TestOutputJson(t *testing.T) {
	t.Parallel()

	testFolder, err := files.CopyTerraformFolderToTemp("../../test/fixtures/terraform-output", t.Name())
	require.NoError(t, err)

	options := &terraform.Options{
		TerraformDir: testFolder,
	}

	terraform.InitAndApply(t, options)

	expected := `{
  "bool": {
    "sensitive": false,
    "type": "bool",
    "value": true
  },
  "number": {
    "sensitive": false,
    "type": "number",
    "value": 3.14
  },
  "number1": {
    "sensitive": false,
    "type": "number",
    "value": 3
  },
  "string": {
    "sensitive": false,
    "type": "string",
    "value": "This is a string."
  },
  "unicode_string": {
    "sensitive": false,
    "type": "string",
    "value": "söme chäräcter"
  }
}`

	str := terraform.OutputJSON(t, options, "")
	require.Equal(t, expected, str, "JSON %q should match %q", expected, str)
}

func TestOutputStruct(t *testing.T) {
	t.Parallel()

	type TestStruct struct {
		Somemap     map[string]any
		Somestring  string
		Listmaps    []map[string]any
		Liststrings []string
		Somefloat   float64
		Someint     int
		Somebool    bool
	}

	testFolder, err := files.CopyTerraformFolderToTemp("../../test/fixtures/terraform-output-struct", t.Name())
	if err != nil {
		t.Fatal(err)
	}

	options := &terraform.Options{
		TerraformDir: testFolder,
		// Let's test order or ExtraArgs while we are at it
		ExtraArgs: terraform.ExtraArgs{
			Output: []string{"-state=terraform.tfstate"},
		},
	}

	terraform.InitAndApply(t, options)

	expectedObject := TestStruct{
		Somebool:    true,
		Somefloat:   0.1,
		Someint:     1,
		Somestring:  "two",
		Somemap:     map[string]any{"three": 3.0, "four": "four"},
		Listmaps:    []map[string]any{{"five": 5.0, "six": "six"}},
		Liststrings: []string{"seven", "eight", "nine"},
	}
	actualObject := TestStruct{}

	terraform.OutputStruct(t, options, "object", &actualObject)

	expectedList := []TestStruct{
		{
			Somebool:   true,
			Somefloat:  0.1,
			Someint:    1,
			Somestring: "two",
		},
		{
			Somebool:   false,
			Somefloat:  0.3,
			Someint:    4,
			Somestring: "five",
		},
	}
	actualList := []TestStruct{}

	terraform.OutputStruct(t, options, "list_of_objects", &actualList)

	require.Equal(t, expectedObject, actualObject, "Object should be %q, got %q", expectedObject, actualObject)
	require.Equal(t, expectedList, actualList, "List should be %q, got %q", expectedList, actualList)
}

func TestOutputsAll(t *testing.T) {
	t.Parallel()

	testFolder, err := files.CopyTerraformFolderToTemp("../../test/fixtures/terraform-output-all", t.Name())
	if err != nil {
		t.Fatal(err)
	}

	options := &terraform.Options{
		TerraformDir: testFolder,
	}

	terraform.InitAndApply(t, options)
	out := terraform.OutputAll(t, options)

	expectedLen := 4
	require.Len(t, out, expectedLen, "Output should contain %d items", expectedLen)

	// String Value
	expectedString := "Sun"
	str, ok := out["our_star"].(string)
	require.Truef(t, ok, "Wrong data type for 'our_star', expected string, got %T", out["our_star"])
	require.Equal(t, expectedString, str, "String %q should match %q", expectedString, str)

	// List Value
	expectedListLen := 3
	outputInterfaceList, ok := out["stars"].([]any)
	require.Truef(t, ok, "Wrong data type for 'stars', expected [], got %T", out["stars"])

	expectedListItem := "Betelgeuse"

	require.Len(t, outputInterfaceList, expectedListLen, "Output list should contain %d items", expectedListLen)
	require.Equal(t, expectedListItem, outputInterfaceList[2].(string), "List item %q should match %q",
		expectedListItem, outputInterfaceList[0].(string))

	// Map Value
	expectedMapLen := 4
	outputInterfaceMap, ok := out["constellations"].(map[string]any)
	require.Truef(t, ok, "Wrong data type for 'constellations', expected map[string], got %T", out["constellations"])

	expectedMapItem := "Aldebaran"

	require.Len(t, outputInterfaceMap, expectedMapLen, "Output map should contain 4 items")
	require.Equal(t, expectedMapItem, outputInterfaceMap["Taurus"].(string), "Map item %q should match %q",
		expectedMapItem, outputInterfaceMap["Taurus"].(string))
}

func TestOutputsForKeysError(t *testing.T) {
	t.Parallel()

	testFolder, err := files.CopyTerraformFolderToTemp("../../test/fixtures/terraform-output-map", t.Name())
	require.NoError(t, err)

	options := &terraform.Options{
		TerraformDir: testFolder,
	}

	terraform.InitAndApply(t, options)

	_, err = terraform.OutputForKeysE(t, options, []string{"random_key"})
	require.Error(t, err)
}

func TestOutputWithDebugLogLevel(t *testing.T) {
	t.Parallel()

	testFolder, err := files.CopyTerraformFolderToTemp("../../test/fixtures/terraform-output-mapofobjects", t.Name())
	require.NoError(t, err)

	options := &terraform.Options{
		TerraformDir: testFolder,
	}

	terraform.InitAndApply(t, options)

	_, err = terraform.OutputMapOfObjectsE(t, &terraform.Options{
		TerraformDir: options.TerraformDir,
		EnvVars:      map[string]string{"TF_LOG": "DEBUG"},
	}, "map_of_objects")
	require.NoError(t, err)
}
