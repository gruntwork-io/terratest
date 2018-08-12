package terraform

import (
	"github.com/justinbarrick/go-terraform-plan/plan"
	"io/ioutil"
	"os"
	"testing"
)

// InitAndPlan runs terraform init and plan with the given options and returns a parsed terraform plan.
func InitAndPlan(t *testing.T, options *Options) plan.Plan {
	out, err := InitAndPlanE(t, options)
	if err != nil {
		t.Fatal(err)
	}
	return out
}

// InitAndPlanE runs terraform init and plan with the given options and returns a parsed terraform plan and any errors.
func InitAndPlanE(t *testing.T, options *Options) (plan.Plan, error) {
	if _, err := InitE(t, options); err != nil {
		return plan.Plan{}, err
	}

	if _, err := GetE(t, options); err != nil {
		return plan.Plan{}, err
	}

	return PlanE(t, options)
}

// Plan runs terraform plan with the given options and returns a parsed terraform plan.
func Plan(t *testing.T, options *Options) plan.Plan {
	out, err := PlanE(t, options)
	if err != nil {
		t.Fatal(err)
	}
	return out
}

// PlanE runs terraform plan with the given options and returns a parsed terraform plan and any errors.
func PlanE(t *testing.T, options *Options) (plan.Plan, error) {
	file, err := ioutil.TempFile("", "plan")
	if err != nil {
		return plan.Plan{}, err
	}
	defer os.Remove(file.Name())

	_, err = RunTerraformCommandE(t, options, FormatArgs(options.Vars, "plan", "-input=false", "-lock=false", "-out="+file.Name())...)
	if err != nil {
		return plan.Plan{}, err
	}

	return plan.ReadPlanFile(file)
}
