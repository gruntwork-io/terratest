package test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestTerraformPlanDetailsExample(t *testing.T) {
	outputPath := "/tmp/terraform-plan.plan"
	terraformOptions := &terraform.Options{
		TerraformDir:   "../examples/terraform-plan-details-example",
		PlanOutputFile: outputPath,
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndPlan(t, terraformOptions)
	output := terraform.Show(t, terraformOptions, terraform.JSON, outputPath)
	info, err := terraform.NewPlanInfo(output)
	assert.Nil(t, err)

	for _, resource := range info.AllResources {
		if resource.Name == "my_null_resource" {
			triggers := resource.Attributes["triggers"].(map[string]interface{})
			assert.Equal(t, triggers["some_attribute"], "attr-val1")
		}
	}
}
