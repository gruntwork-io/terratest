package test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestTerraformPlanDetailsExample(t *testing.T) {
	terraformOptions := &terraform.Options{
		// website::tag::1:: Set the path to the Terraform code that will be tested.
		TerraformDir:   "../examples/terraform-plan-details-example",
		PlanOutputFile: "/tmp/terraform-plan.plan",
	}
	terraform.InitAndPlan(t, terraformOptions)
	output := terraform.Show(t, terraformOptions, terraform.JSON, terraformOptions.PlanOutputFile)
	pinfo := terraform.NewPlanInfo(output)
	t.Log(pinfo.PlannedValues.RootModule.Resources[0].Values["triggers"])
	assert.True(t, true)
}
