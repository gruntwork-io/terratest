package terraform

import (
	"testing"

	http_helper "github.com/gruntwork-io/terratest/modules/http-helper"
	tfjson "github.com/hashicorp/terraform-json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	// NOTE: We pull down the json files from github during test runtime as opposed to checking it in as these source
	// files are licensed under MPL and we want to avoid a dual license scenario where some source files in terratest
	// are licensed under a different license.
	basicJsonUrl      = "https://raw.githubusercontent.com/hashicorp/terraform-json/v0.13.0/testdata/basic/plan.json"
	deepModuleJsonUrl = "https://raw.githubusercontent.com/hashicorp/terraform-json/v0.13.0/testdata/deep_module/plan.json"

	changesJsonUrl = "https://raw.githubusercontent.com/hashicorp/terraform-json/v0.13.0/testdata/has_changes/plan.json"
)

func validateHTTPSuccess(statusCode int) bool {
	return statusCode >= 200 && statusCode < 300
}

func TestPlannedValuesMapWithBasicJson(t *testing.T) {
	t.Parallel()

	// Retrieve test data from the terraform-json project.
	statusCode, jsonData := http_helper.HttpGet(t, basicJsonUrl, nil)
	require.True(t, validateHTTPSuccess(statusCode))

	plan, err := ParsePlanJSON(jsonData)
	require.NoError(t, err)

	query := []string{
		"data.null_data_source.baz",
		"null_resource.bar",
		"null_resource.baz[0]",
		"null_resource.baz[1]",
		"null_resource.baz[2]",
		"null_resource.foo",
		"module.foo.null_resource.aliased",
		"module.foo.null_resource.foo",
	}
	for _, key := range query {
		RequirePlannedValuesMapKeyExists(t, plan, key)
		resource := plan.ResourcePlannedValuesMap[key]
		assert.Equal(t, resource.Address, key)
	}
}

func TestPlannedValuesMapWithDeepModuleJson(t *testing.T) {
	t.Parallel()

	// Retrieve test data from the terraform-json project.
	statusCode, jsonData := http_helper.HttpGet(t, deepModuleJsonUrl, nil)
	require.True(t, validateHTTPSuccess(statusCode))

	plan, err := ParsePlanJSON(jsonData)
	require.NoError(t, err)

	query := []string{
		"module.foo.module.bar.null_resource.baz",
	}
	for _, key := range query {
		AssertPlannedValuesMapKeyExists(t, plan, key)
	}
}

func TestResourceChangesJson(t *testing.T) {
	t.Parallel()

	// Retrieve test data from the terraform-json project.
	statusCode, jsonData := http_helper.HttpGet(t, changesJsonUrl, nil)
	require.True(t, validateHTTPSuccess(statusCode))

	plan, err := ParsePlanJSON(jsonData)
	require.NoError(t, err)

	// Spot check a few changes to make sure the right address was registered
	RequireResourceChangesMapKeyExists(t, plan, "module.foo.null_resource.foo")
	fooChanges := plan.ResourceChangesMap["module.foo.null_resource.foo"]
	require.NotNil(t, fooChanges.Change)
	assert.Equal(t, fooChanges.Change.After.(map[string]interface{})["triggers"].(map[string]interface{})["foo"].(string), "bar")

	RequireResourceChangesMapKeyExists(t, plan, "null_resource.bar")
	barChanges := plan.ResourceChangesMap["null_resource.bar"]
	require.NotNil(t, barChanges.Change)
	assert.Equal(t, barChanges.Change.After.(map[string]interface{})["triggers"].(map[string]interface{})["foo_id"].(string), "424881806176056736")
}

func TestOutputChangesJson(t *testing.T) {
	t.Parallel()

	// Retrieve test data from the terraform-json project.
	statusCode, jsonData := http_helper.HttpGet(t, changesJsonUrl, nil)
	require.True(t, validateHTTPSuccess(statusCode))

	plan, err := ParsePlanJSON(jsonData)
	require.NoError(t, err)

	// Spot check a few changes to make sure the right address was registered
	RequireOutputChangesMapKeyExists(t, plan, "foo")
	fooChanges := plan.OutputChangesMap["foo"]
	require.NotNil(t, fooChanges)
	assert.Equal(t, fooChanges.Actions, tfjson.Actions{"create"})
	assert.Equal(t, fooChanges.After, "bar")

	AssertOutputChangesMapKeyExists(t, plan, "map")
	mapChanges := plan.OutputChangesMap["map"]
	require.NotNil(t, mapChanges)
	assert.Equal(t, mapChanges.Actions, tfjson.Actions{"create"})
	assert.Equal(t, mapChanges.After.(map[string]interface{}), map[string]interface{}{"foo": "bar", "number": float64(42)})
}
