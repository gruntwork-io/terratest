package test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

// An example of how to test information in a terraform plan.  This tests that
// a resource's name is correctly created from a value in a `for_each`
// expression.
func TestTerraformPlanDetailsExample(t *testing.T) {
	outputPath := "/tmp/terraform-plan.plan"
	terraformOptions := &terraform.Options{
		TerraformDir: "../examples/terraform-plan-details-example",

		// This will write a binary plan output file when plan is run
		PlanOutputFile: outputPath,
	}

	terraform.InitAndPlan(t, terraformOptions)

	// Outputs a JSON-formatted version of the terraform plan file
	output := terraform.Show(t, terraformOptions, terraform.JSON, outputPath)

	// Constructs a PlanInfo datastructure with information about the plan
	info, err := terraform.NewPlanInfo(output)
	assert.Nil(t, err)

	found := false

	for _, resource := range info.AllResources {
		if resource.Name == "my_null_resource" {
			triggers := resource.Attributes["triggers"].(map[string]interface{})

			// Test that the resource that we're expecting to be created is in
			// the plan with the expected name
			assert.Equal(t, triggers["some_attribute"], "attr-val1")

			found = true
		}
	}

	assert.True(t, found)
}
