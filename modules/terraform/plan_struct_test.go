package terraform_test

import (
	"testing"

	http_helper "github.com/gruntwork-io/terratest/modules/http-helper"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	// NOTE: We pull down the json files from github during test runtime as opposed to checking it in as these source
	// files are licensed under MPL and we want to avoid a dual license scenario where some source files in terratest
	// are licensed under a different license.
	basicJSONURL      = "https://raw.githubusercontent.com/hashicorp/terraform-json/v0.8.0/testdata/basic/plan.json"
	deepModuleJSONURL = "https://raw.githubusercontent.com/hashicorp/terraform-json/v0.8.0/testdata/deep_module/plan.json"

	changesJSONURL = "https://raw.githubusercontent.com/hashicorp/terraform-json/v0.8.0/testdata/has_changes/plan.json"
)

func TestPlannedValuesMapWithBasicJson(t *testing.T) {
	t.Parallel()

	// Retrieve test data from the terraform-json project.
	_, jsonData := http_helper.HTTPGetContext(t, t.Context(), basicJSONURL, nil)

	plan, err := terraform.ParsePlanJSON(jsonData)
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
		terraform.RequirePlannedValuesMapKeyExists(t, plan, key)

		resource := plan.ResourcePlannedValuesMap[key]
		assert.Equal(t, key, resource.Address)
	}
}

func TestPlannedValuesMapWithDeepModuleJson(t *testing.T) {
	t.Parallel()

	// Retrieve test data from the terraform-json project.
	_, jsonData := http_helper.HTTPGetContext(t, t.Context(), deepModuleJSONURL, nil)

	plan, err := terraform.ParsePlanJSON(jsonData)
	require.NoError(t, err)

	query := []string{
		"module.foo.module.bar.null_resource.baz",
	}
	for _, key := range query {
		terraform.AssertPlannedValuesMapKeyExists(t, plan, key)
	}
}

func TestResourceChangesJson(t *testing.T) {
	t.Parallel()

	// Retrieve test data from the terraform-json project.
	_, jsonData := http_helper.HTTPGetContext(t, t.Context(), changesJSONURL, nil)

	plan, err := terraform.ParsePlanJSON(jsonData)
	require.NoError(t, err)

	// Spot check a few changes to make sure the right address was registered
	terraform.RequireResourceChangesMapKeyExists(t, plan, "module.foo.null_resource.foo")

	fooChanges := plan.ResourceChangesMap["module.foo.null_resource.foo"]
	require.NotNil(t, fooChanges.Change)
	assert.Equal(t, "bar", fooChanges.Change.After.(map[string]any)["triggers"].(map[string]any)["foo"].(string))

	terraform.RequireResourceChangesMapKeyExists(t, plan, "null_resource.bar")

	barChanges := plan.ResourceChangesMap["null_resource.bar"]
	require.NotNil(t, barChanges.Change)
	assert.Equal(t, "424881806176056736", barChanges.Change.After.(map[string]any)["triggers"].(map[string]any)["foo_id"].(string))
}
