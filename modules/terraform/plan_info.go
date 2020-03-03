package terraform

import (
	"encoding/json"
	"fmt"
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
	Values       map[string]interface{}
	Changes      []Change
}

type Module struct {
	Address   string
	Resources []Resource
}

type PlannedValues struct {
	RootModule   Module   `json:"root_module"`
	ChildModules []Module `json:"child_modules"`
}

type PlanInfo struct {
	RawPlannedValues PlannedValues `json:"planned_values"`
	ChangedResources []Resource    `json:"resource_changes"`
	AllResources     []Resource
}

// NewPlanInfo returns a PlanInfo struct given the json-formatted output of a terraform plan.
func NewPlanInfo(jsonOutput string) PlanInfo {
	var v PlanInfo
	err := json.Unmarshal([]byte(jsonOutput), &v)

	if err != nil {
		fmt.Println("TODO: Couldn't parse json")
		fmt.Println(err)
	}

	allResources := []Resource{}

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

	return v
}
