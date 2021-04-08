package terraform

import (
	"encoding/json"

	"github.com/gruntwork-io/terratest/modules/testing"
	tfjson "github.com/hashicorp/terraform-json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// PlanStruct is a Go Struct representation of the plan object returned from Terraform (after running `terraform show`).
// Unlike the raw plan representation returned by terraform-json, this struct provides a map that maps the resource
// addresses to the changes and planned values to make it easier to navigate the raw plan struct.
type PlanStruct struct {
	// The raw representation of the plan. See
	// https://www.terraform.io/docs/internals/json-format.html#plan-representation for details on the structure of the
	// plan output.
	RawPlan tfjson.Plan

	// A map that maps full resource addresses (e.g., module.foo.null_resource.test) to the planned values of that
	// resource.
	ResourcePlannedValuesMap map[string]*tfjson.StateResource

	// A map that maps full resource addresses (e.g., module.foo.null_resource.test) to the planned actions terraform
	// will take on that resource.
	ResourceChangesMap map[string]*tfjson.ResourceChange
}

// parsePlanJson takes in the json string representation of the terraform plan and returns a go struct representation
// for easy introspection.
func parsePlanJson(jsonStr string) (*PlanStruct, error) {
	plan := &PlanStruct{}

	if err := json.Unmarshal([]byte(jsonStr), &plan.RawPlan); err != nil {
		return nil, err
	}

	plan.ResourcePlannedValuesMap = parsePlannedValues(plan)
	plan.ResourceChangesMap = parseResourceChanges(plan)
	return plan, nil
}

// parseResourceChanges takes a plan and returns a map that maps resource addresses to the planned changes for that
// resource. If there are no changes, this returns an empty map instead of erroring.
func parseResourceChanges(plan *PlanStruct) map[string]*tfjson.ResourceChange {
	out := map[string]*tfjson.ResourceChange{}
	for _, change := range plan.RawPlan.ResourceChanges {
		out[change.Address] = change
	}
	return out
}

// parsePlannedValues takes a plan and walks through the planned values to return a map that maps the full resource
// addresses to the planned resources. If there are no planned values, this returns an empty map instead of erroring.
func parsePlannedValues(plan *PlanStruct) map[string]*tfjson.StateResource {
	plannedValues := plan.RawPlan.PlannedValues
	if plannedValues == nil {
		// No planned values, so return empty map.
		return map[string]*tfjson.StateResource{}
	}

	rootModule := plannedValues.RootModule
	if rootModule == nil {
		// No module resources, so return empty map.
		return map[string]*tfjson.StateResource{}
	}
	return parseModulePlannedValues(rootModule)
}

// parseModulePlannedValues will recursively walk through the modules in the planned_values of the plan struct to
// construct a map that maps the full resource addresses to the planned resource.
func parseModulePlannedValues(module *tfjson.StateModule) map[string]*tfjson.StateResource {
	out := map[string]*tfjson.StateResource{}
	for _, resource := range module.Resources {
		// NOTE: the Address attribute of the module resource always returns the full address, even when the resource is
		// nested within sub modules.
		out[resource.Address] = resource
	}

	// NOTE: base case of recursion is when ChildModules is empty list.
	for _, child := range module.ChildModules {
		// Recurse in to the child module. We take a recursive approach here despite limitations of the recursion stack
		// in golang due to the fact that it is rare to have heavily deep module calls in Terraform. So we optimize for
		// code readability as opposed to performance.
		childMap := parseModulePlannedValues(child)
		for k, v := range childMap {
			out[k] = v
		}
	}
	return out
}

// AssertPlannedValuesMapKeyExists checks if the given key exists in the map, failing the test if it does not.
func AssertPlannedValuesMapKeyExists(t testing.TestingT, plan *PlanStruct, keyQuery string) {
	_, hasKey := plan.ResourcePlannedValuesMap[keyQuery]
	assert.Truef(t, hasKey, "Given planned values map does not have key %s", keyQuery)
}

// RequirePlannedValuesMapKeyExists checks if the given key exists in the map, failing and halting the test if it does not.
func RequirePlannedValuesMapKeyExists(t testing.TestingT, plan *PlanStruct, keyQuery string) {
	_, hasKey := plan.ResourcePlannedValuesMap[keyQuery]
	require.Truef(t, hasKey, "Given planned values map does not have key %s", keyQuery)
}

// AssertResourceChangesMapKeyExists checks if the given key exists in the map, failing the test if it does not.
func AssertResourceChangesMapKeyExists(t testing.TestingT, plan *PlanStruct, keyQuery string) {
	_, hasKey := plan.ResourceChangesMap[keyQuery]
	assert.Truef(t, hasKey, "Given resource changes map does not have key %s", keyQuery)
}

// RequireResourceChangesMapKeyExists checks if the given key exists in the map, failing the test if it does not.
func RequireResourceChangesMapKeyExists(t testing.TestingT, plan *PlanStruct, keyQuery string) {
	_, hasKey := plan.ResourceChangesMap[keyQuery]
	require.Truef(t, hasKey, "Given resource changes map does not have key %s", keyQuery)
}
