package test

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
		[]string{random.UniqueId()},
	}

	options := GetTerraformOptionsForFormatTests(t)
	options.Vars["example_any"] = testList

	defer terraform.Destroy(t, options)
	terraform.InitAndApply(t, options)
	outputMap := terraform.OutputForKeys(t, options, []string{"example_any"})
	actualExampleList := outputMap["example_any"]
	AssertEqualJson(t, actualExampleList, testList)
}

func TestTerraformFormatNestedTwoLevelList(t *testing.T) {
	t.Parallel()

	testList := [][][]string{
		[][]string{[]string{random.UniqueId()}},
	}

	options := GetTerraformOptionsForFormatTests(t)
	options.Vars["example_any"] = testList

	defer terraform.Destroy(t, options)
	terraform.InitAndApply(t, options)
	outputMap := terraform.OutputForKeys(t, options, []string{"example_any"})
	actualExampleList := outputMap["example_any"]
	AssertEqualJson(t, actualExampleList, testList)
}

func TestTerraformFormatNestedMultipleItems(t *testing.T) {
	t.Parallel()

	testList := [][]string{
		[]string{random.UniqueId(), random.UniqueId()},
		[]string{random.UniqueId(), random.UniqueId(), random.UniqueId()},
	}

	options := GetTerraformOptionsForFormatTests(t)
	options.Vars["example_any"] = testList

	defer terraform.Destroy(t, options)
	terraform.InitAndApply(t, options)
	outputMap := terraform.OutputForKeys(t, options, []string{"example_any"})
	actualExampleList := outputMap["example_any"]
	AssertEqualJson(t, actualExampleList, testList)
}

func TestTerraformFormatNestedOneLevelMap(t *testing.T) {
	t.Parallel()

	testMap := map[string]map[string]string{
		"test": map[string]string{
			"foo": random.UniqueId(),
		},
	}

	options := GetTerraformOptionsForFormatTests(t)
	options.Vars["example_any"] = testMap

	defer terraform.Destroy(t, options)
	terraform.InitAndApply(t, options)
	outputMap := terraform.OutputForKeys(t, options, []string{"example_any"})
	actualExampleMap := outputMap["example_any"]
	AssertEqualJson(t, actualExampleMap, testMap)
}

func TestTerraformFormatNestedTwoLevelMap(t *testing.T) {
	t.Parallel()

	testMap := map[string]map[string]map[string]string{
		"test": map[string]map[string]string{
			"foo": map[string]string{
				"bar": random.UniqueId(),
			},
		},
	}

	options := GetTerraformOptionsForFormatTests(t)
	options.Vars["example_any"] = testMap

	defer terraform.Destroy(t, options)
	terraform.InitAndApply(t, options)
	outputMap := terraform.OutputForKeys(t, options, []string{"example_any"})
	actualExampleMap := outputMap["example_any"]
	AssertEqualJson(t, actualExampleMap, testMap)
}

func TestTerraformFormatNestedMultipleItemsMap(t *testing.T) {
	t.Parallel()

	testMap := map[string]map[string]string{
		"test": map[string]string{
			"foo": random.UniqueId(),
			"bar": random.UniqueId(),
		},
		"other": map[string]string{
			"baz": random.UniqueId(),
			"boo": random.UniqueId(),
		},
	}

	options := GetTerraformOptionsForFormatTests(t)
	options.Vars["example_any"] = testMap

	defer terraform.Destroy(t, options)
	terraform.InitAndApply(t, options)
	outputMap := terraform.OutputForKeys(t, options, []string{"example_any"})
	actualExampleMap := outputMap["example_any"]
	AssertEqualJson(t, actualExampleMap, testMap)
}

func TestTerraformFormatNestedListMap(t *testing.T) {
	t.Parallel()

	testMap := map[string][]string{
		"test": []string{random.UniqueId(), random.UniqueId()},
	}

	options := GetTerraformOptionsForFormatTests(t)
	options.Vars["example_any"] = testMap

	defer terraform.Destroy(t, options)
	terraform.InitAndApply(t, options)
	outputMap := terraform.OutputForKeys(t, options, []string{"example_any"})
	actualExampleMap := outputMap["example_any"]
	AssertEqualJson(t, actualExampleMap, testMap)
}

func GetTerraformOptionsForFormatTests(t *testing.T) *terraform.Options {
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
func AssertEqualJson(t *testing.T, actual interface{}, expected interface{}) {
	actualJson, err := json.Marshal(actual)
	require.NoError(t, err)
	expectedJson, err := json.Marshal(expected)
	require.NoError(t, err)
	assert.Equal(t, actualJson, expectedJson)
}
