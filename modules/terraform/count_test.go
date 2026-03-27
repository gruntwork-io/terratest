package terraform_test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/files"
	"github.com/gruntwork-io/terratest/modules/terraform"
	ttesting "github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetResourceCount(t *testing.T) {
	t.Parallel()
	testFolder, err := files.CopyTerraformFolderToTemp("../../test/fixtures/terraform-basic-configuration", t.Name())
	require.NoError(t, err)

	terraformOptions := &terraform.Options{
		TerraformDir: testFolder,
		Vars: map[string]any{
			"cnt": 1,
		},
	}

	cnt := terraform.GetResourceCount(t, terraform.InitAndPlan(t, terraformOptions))
	assert.Equal(t, 1, cnt.Add)
	assert.Equal(t, 0, cnt.Change)
	assert.Equal(t, 0, cnt.Destroy)
}

func TestGetResourceCountEColor(t *testing.T) { //nolint:tparallel // subtests share mutable terraform options
	t.Parallel()

	runTestGetResourceCountE(t, false)
}

func TestGetResourceCountENoColor(t *testing.T) { //nolint:tparallel // subtests share mutable terraform options
	t.Parallel()

	runTestGetResourceCountE(t, true)
}

func runTestGetResourceCountE(t *testing.T, noColor bool) { //nolint:tparallel // subtests share mutable terraform options
	t.Helper()
	testCases := []struct {
		tfFuncToRun     func(t ttesting.TestingT, options *terraform.Options) string
		name            string
		cntValue        int
		expectedAdd     int
		expectedChange  int
		expectedDestroy int
	}{
		{name: "PlanZero", tfFuncToRun: terraform.InitAndPlan, cntValue: 0, expectedAdd: 0, expectedChange: 0, expectedDestroy: 0},
		{name: "ApplyZero", tfFuncToRun: terraform.InitAndApply, cntValue: 0, expectedAdd: 0, expectedChange: 0, expectedDestroy: 0},
		{name: "PlanAddResouce", tfFuncToRun: terraform.InitAndPlan, cntValue: 2, expectedAdd: 2, expectedChange: 0, expectedDestroy: 0},
		{name: "ApplyAddResouce", tfFuncToRun: terraform.InitAndApply, cntValue: 2, expectedAdd: 2, expectedChange: 0, expectedDestroy: 0},
		{name: "PlanNoOp", tfFuncToRun: terraform.InitAndApply, cntValue: 2, expectedAdd: 0, expectedChange: 0, expectedDestroy: 0},
		{name: "ApplyNoOp", tfFuncToRun: terraform.InitAndApply, cntValue: 2, expectedAdd: 0, expectedChange: 0, expectedDestroy: 0},
		{name: "PlanDestroyResource", tfFuncToRun: terraform.InitAndPlan, cntValue: 1, expectedAdd: 0, expectedChange: 0, expectedDestroy: 1},
		{name: "ApplyDestroyResource", tfFuncToRun: terraform.InitAndApply, cntValue: 1, expectedAdd: 0, expectedChange: 0, expectedDestroy: 1},
		{name: "Destroy", tfFuncToRun: terraform.Destroy, cntValue: 1, expectedAdd: 0, expectedChange: 0, expectedDestroy: 1},
		{name: "DestroyNoOp", tfFuncToRun: terraform.Destroy, cntValue: 1, expectedAdd: 0, expectedChange: 0, expectedDestroy: 0},
	}

	testFolder, err := files.CopyTerraformFolderToTemp("../../test/fixtures/terraform-basic-configuration", t.Name())
	require.NoError(t, err)

	terraformOptions := &terraform.Options{
		TerraformDir: testFolder,
		Vars: map[string]any{
			"cnt": 0,
		},
		NoColor: noColor,
	}

	for _, tc := range testCases {
		t.Run(tc.name,
			func(t *testing.T) {
				terraformOptions.Vars["cnt"] = tc.cntValue
				cnt, err := terraform.GetResourceCountE(t, tc.tfFuncToRun(t, terraformOptions))
				require.NoError(t, err)
				assert.Equal(t, tc.expectedAdd, cnt.Add)
				assert.Equal(t, tc.expectedChange, cnt.Change)
				assert.Equal(t, tc.expectedDestroy, cnt.Destroy)
			})
	}

	t.Run("InvalidInput",
		func(t *testing.T) {
			terraformOptions.Vars["cnt"] = "abc"
			cmdout, _ := terraform.PlanE(t, terraformOptions)
			cnt, err := terraform.GetResourceCountE(t, cmdout)
			require.EqualError(t, err, terraform.GetResourceCountErrMessage)
			assert.Nil(t, cnt)
		})
}
