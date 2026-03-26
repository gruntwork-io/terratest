package terraform_test

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/files"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/require"
)

func TestShowWithInlinePlan(t *testing.T) {
	t.Parallel()

	testFolder, err := files.CopyTerraformFolderToTemp("../../test/fixtures/terraform-basic-configuration", t.Name())
	require.NoError(t, err)

	planFilePath := filepath.Join(testFolder, "plan.out")

	options := &terraform.Options{
		TerraformDir: testFolder,
		PlanFilePath: planFilePath,
		Vars: map[string]any{
			"cnt": 1,
		},
	}

	out := terraform.InitAndPlan(t, options)
	out = strings.ReplaceAll(out, "\n", "")
	require.Contains(t, out, "Saved the plan to:"+planFilePath)
	require.FileExists(t, planFilePath, "Plan file was not saved to expected location:", planFilePath)

	// show command does not accept Vars
	showOptions := &terraform.Options{
		TerraformDir: testFolder,
		PlanFilePath: planFilePath,
	}

	// Test the JSON string
	planJSON := terraform.Show(t, showOptions)
	require.Contains(t, planJSON, "null_resource.test[0]")
}

func TestShowWithStructInlinePlan(t *testing.T) {
	t.Parallel()

	testFolder, err := files.CopyTerraformFolderToTemp("../../test/fixtures/terraform-basic-configuration", t.Name())
	require.NoError(t, err)

	planFilePath := filepath.Join(testFolder, "plan.out")

	options := &terraform.Options{
		TerraformDir: testFolder,
		PlanFilePath: planFilePath,
		Vars: map[string]any{
			"cnt": 1,
		},
	}

	out := terraform.InitAndPlan(t, options)
	out = strings.ReplaceAll(out, "\n", "")
	require.Contains(t, out, "Saved the plan to:"+planFilePath)
	require.FileExists(t, planFilePath, "Plan file was not saved to expected location:", planFilePath)

	// show command does not accept Vars
	showOptions := &terraform.Options{
		TerraformDir: testFolder,
		PlanFilePath: planFilePath,
	}

	// Test the JSON string
	plan := terraform.ShowWithStruct(t, showOptions)
	require.Contains(t, plan.ResourcePlannedValuesMap, "null_resource.test[0]")
}
