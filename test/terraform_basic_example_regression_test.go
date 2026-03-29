package test_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	test_structure "github.com/gruntwork-io/terratest/modules/test-structure"
)

// The tests in this folder are not example usage of Terratest. Instead, this is a regression test to ensure the
// formatting rules work with an actual Terraform call when using more complex structures.

func TestTerraformFormatNestedOneLevelList(t *testing.T) {
	t.Parallel()

	testList := [][]string{
		{random.UniqueID()},
	}

	options := getTerraformOptionsForFormatTests(t)
	options.Vars["example_any"] = testList

	defer terraform.DestroyContext(t, t.Context(), options)

	terraform.InitAndApplyContext(t, t.Context(), options)
	outputMap := terraform.OutputForKeysContext(t, t.Context(), options, []string{"example_any"})
	actualExampleList := outputMap["example_any"]

	assertEqualJSON(t, actualExampleList, testList)
}

func TestTerraformFormatNestedTwoLevelList(t *testing.T) {
	t.Parallel()

	testList := [][][]string{
		{{random.UniqueID()}},
	}

	options := getTerraformOptionsForFormatTests(t)
	options.Vars["example_any"] = testList

	defer terraform.DestroyContext(t, t.Context(), options)

	terraform.InitAndApplyContext(t, t.Context(), options)
	outputMap := terraform.OutputForKeysContext(t, t.Context(), options, []string{"example_any"})
	actualExampleList := outputMap["example_any"]

	assertEqualJSON(t, actualExampleList, testList)
}

func TestTerraformFormatNestedMultipleItems(t *testing.T) {
	t.Parallel()

	testList := [][]string{
		{random.UniqueID(), random.UniqueID()},
		{random.UniqueID(), random.UniqueID(), random.UniqueID()},
	}

	options := getTerraformOptionsForFormatTests(t)
	options.Vars["example_any"] = testList

	defer terraform.DestroyContext(t, t.Context(), options)

	terraform.InitAndApplyContext(t, t.Context(), options)
	outputMap := terraform.OutputForKeysContext(t, t.Context(), options, []string{"example_any"})
	actualExampleList := outputMap["example_any"]

	assertEqualJSON(t, actualExampleList, testList)
}

func TestTerraformFormatNestedOneLevelMap(t *testing.T) {
	t.Parallel()

	testMap := map[string]map[string]string{
		"test": {
			"foo": random.UniqueID(),
		},
	}

	options := getTerraformOptionsForFormatTests(t)
	options.Vars["example_any"] = testMap

	defer terraform.DestroyContext(t, t.Context(), options)

	terraform.InitAndApplyContext(t, t.Context(), options)
	outputMap := terraform.OutputForKeysContext(t, t.Context(), options, []string{"example_any"})
	actualExampleMap := outputMap["example_any"]

	assertEqualJSON(t, actualExampleMap, testMap)
}

func TestTerraformFormatNestedTwoLevelMap(t *testing.T) {
	t.Parallel()

	testMap := map[string]map[string]map[string]string{
		"test": {
			"foo": {
				"bar": random.UniqueID(),
			},
		},
	}

	options := getTerraformOptionsForFormatTests(t)
	options.Vars["example_any"] = testMap

	defer terraform.DestroyContext(t, t.Context(), options)

	terraform.InitAndApplyContext(t, t.Context(), options)
	outputMap := terraform.OutputForKeysContext(t, t.Context(), options, []string{"example_any"})
	actualExampleMap := outputMap["example_any"]

	assertEqualJSON(t, actualExampleMap, testMap)
}

func TestTerraformFormatNestedMultipleItemsMap(t *testing.T) {
	t.Parallel()

	testMap := map[string]map[string]string{
		"test": {
			"foo": random.UniqueID(),
			"bar": random.UniqueID(),
		},
		"other": {
			"baz": random.UniqueID(),
			"boo": random.UniqueID(),
		},
	}

	options := getTerraformOptionsForFormatTests(t)
	options.Vars["example_any"] = testMap

	defer terraform.DestroyContext(t, t.Context(), options)

	terraform.InitAndApplyContext(t, t.Context(), options)
	outputMap := terraform.OutputForKeysContext(t, t.Context(), options, []string{"example_any"})
	actualExampleMap := outputMap["example_any"]

	assertEqualJSON(t, actualExampleMap, testMap)
}

func TestTerraformFormatNestedListMap(t *testing.T) {
	t.Parallel()

	testMap := map[string][]string{
		"test": {random.UniqueID(), random.UniqueID()},
	}

	options := getTerraformOptionsForFormatTests(t)
	options.Vars["example_any"] = testMap

	defer terraform.DestroyContext(t, t.Context(), options)

	terraform.InitAndApplyContext(t, t.Context(), options)
	outputMap := terraform.OutputForKeysContext(t, t.Context(), options, []string{"example_any"})
	actualExampleMap := outputMap["example_any"]

	assertEqualJSON(t, actualExampleMap, testMap)
}

func getTerraformOptionsForFormatTests(t *testing.T) *terraform.Options {
	t.Helper()

	exampleFolder := test_structure.CopyTerraformFolderToTemp(t, "../", "examples/terraform-basic-example")

	// Set up terratest to retry on known failures
	maxTerraformRetries := 3
	sleepBetweenTerraformRetries := 5 * time.Second
	retryableTerraformErrors := map[string]string{
		// `terraform init` frequently fails in CI due to network issues accessing plugins. The reason is unknown, but
		// eventually these succeed after a few retries.
		".*unable to verify signature.*":             "Failed to retrieve plugin due to transient network error.",
		".*unable to verify checksum.*":              "Failed to retrieve plugin due to transient network error.",
		".*no provider exists with the given name.*": "Failed to retrieve plugin due to transient network error.",
		".*registry service is unreachable.*":        "Failed to retrieve plugin due to transient network error.",
		".*connection reset by peer.*":               "Failed to retrieve plugin due to transient network error.",
	}

	terraformOptions := &terraform.Options{
		TerraformDir:             exampleFolder,
		Vars:                     map[string]interface{}{},
		NoColor:                  true,
		RetryableTerraformErrors: retryableTerraformErrors,
		MaxRetries:               maxTerraformRetries,
		TimeBetweenRetries:       sleepBetweenTerraformRetries,
	}

	return terraformOptions
}

// The value of the output nested in the outputMap returned by OutputForKeys uses the interface{} type for nested
// structures. This can't be compared to actual types like [][]string{}, so we instead compare the json versions.
func assertEqualJSON(t *testing.T, actual interface{}, expected interface{}) {
	t.Helper()

	assert.JSONEq(t, toJSON(t, expected), toJSON(t, actual))
}

func toJSON(t *testing.T, v interface{}) string {
	t.Helper()

	data, err := json.Marshal(v)
	require.NoError(t, err)

	return string(data)
}
