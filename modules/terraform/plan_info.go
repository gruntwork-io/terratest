package terraform

import (
	"encoding/json"
)

type Attributes map[string]interface{}

type Change struct {
	Actions      []string
	Before       Attributes
	After        Attributes
	AfterUnknown Attributes `json:"after_unknown"`
}

// TODO: The `Index` field can be a string or an int depending on if `count` or
// `for_each` was use.  Doesn't parse as an int right now.
type Resource struct {
	Module       string
	Address      string
	Mode         string
	Type         string
	Name         string
	Index        string
	ProviderName string `json:"provider_name"`
}

type ChangedResource struct {
	Resource
	Changes []Change
}

type KnownResource struct {
	Resource
	Attributes map[string]interface{} `json:"values"`
}

type Module struct {
	Address   string
	Resources []KnownResource
}

type PlannedValues struct {
	RootModule   Module   `json:"root_module"`
	ChildModules []Module `json:"child_modules"`
}

// PlanInfo contains information about a terraform plan.  The info in this data
// structure is a (very) slight simplication of a JSON formatted terraform
// plan, described here:
// https://www.terraform.io/docs/internals/json-format.html#plan-representation
//
// ChangedResources is a list of resources that describe the changes that
// terraform will make.  These changes are represented as `Change` structs in
// the resource's `Changes` field.  If a resource would not be changed by a
// plan, it will not show up in the `ChangedResources` field.
//
// AllResources is a list of all of the KNOWN project resources in the state
// after config in the plan would be applied.  The attributes of these
// resources are in the `Attributes` field on the resources.  You can make
// assertions about these attributes.
type PlanInfo struct {
	RawPlannedValues PlannedValues     `json:"planned_values"`
	ChangedResources []ChangedResource `json:"resource_changes"`
	AllResources     []KnownResource
}

// NewPlanInfo returns a PlanInfo struct given the json-formatted output of a terraform plan.
func NewPlanInfo(jsonOutput string) (PlanInfo, error) {
	var v PlanInfo
	err := json.Unmarshal([]byte(jsonOutput), &v)

	allResources := []KnownResource{}

	// Flatten the root module and child module planned resources
	for _, resource := range v.RawPlannedValues.RootModule.Resources {
		resource.Module = "root"
		allResources = append(allResources, resource)
	}

	for _, module := range v.RawPlannedValues.ChildModules {
		for _, resource := range module.Resources {
			resource.Module = module.Address
			allResources = append(allResources, resource)
		}
	}

	v.AllResources = allResources

	return v, err
}
