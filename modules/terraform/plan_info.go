package terraform

import (
	"encoding/json"
	"fmt"
)

type Resource struct {
	Address string
	Name    string
	Index   string
	Values  map[string]interface{}
}

type Module struct {
	Resources []Resource
}

type PlannedValues struct {
	RootModule Module `json:"root_module"`
}

type PlanInfo struct {
	PlannedValues PlannedValues `json:"planned_values"`
}

func NewPlanInfo(jsonOutput string) PlanInfo {
	var v PlanInfo
	err := json.Unmarshal([]byte(jsonOutput), &v)

	if err != nil {
		fmt.Println("TODO: Couldn't parse json")
	}

	return v
}
