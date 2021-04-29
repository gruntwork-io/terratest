package terraform

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/gruntwork-io/terratest/modules/files"
	"github.com/stretchr/testify/require"
)

func TestShowWithInlinePlan(t *testing.T) {
	t.Parallel()

	testFolder, err := files.CopyTerraformFolderToTemp("../../test/fixtures/terraform-basic-configuration", t.Name())
	require.NoError(t, err)
	planFilePath := filepath.Join(testFolder, "plan.out")

	options := &Options{
		TerraformDir: testFolder,
		PlanFilePath: planFilePath,
		Vars: map[string]interface{}{
			"cnt": 1,
		},
	}

	out := InitAndPlan(t, options)
	require.Contains(t, out, fmt.Sprintf("Saved the plan to: %s", planFilePath))
	require.FileExists(t, planFilePath, "Plan file was not saved to expected location:", planFilePath)

	// show command does not accept Vars
	showOptions := &Options{
		TerraformDir: testFolder,
		PlanFilePath: planFilePath,
	}

	// Test the JSON string
	planJSON := Show(t, showOptions)
	require.Contains(t, planJSON, "null_resource.test[0]")
}
