package terraform

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/files"
	ttesting "github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetResourceCount(t *testing.T) {
	t.Parallel()
	testFolder, err := files.CopyTerraformFolderToTemp("../../test/fixtures/terraform-basic-configuration", t.Name())
	require.NoError(t, err)

	terraformOptions := &Options{
		TerraformDir: testFolder,
		Vars: map[string]interface{}{
			"cnt": 1,
		},
	}

	cnt := GetResourceCount(t, InitAndPlan(t, terraformOptions))
	assert.Equal(t, 1, cnt.Add)
	assert.Equal(t, 0, cnt.Change)
	assert.Equal(t, 0, cnt.Destroy)
}

func TestGetResourceCountEColor(t *testing.T) {
	t.Parallel()
	runTestGetResourceCountE(t, false)
}

func TestGetResourceCountENoColor(t *testing.T) {
	t.Parallel()
	runTestGetResourceCountE(t, true)
}

func runTestGetResourceCountE(t *testing.T, noColor bool) {
	testCases := []struct {
		Name                                         string
		tfFuncToRun                                  func(t ttesting.TestingT, options *Options) string
		cntValue                                     int
		expectedAdd, expectedChange, expectedDestroy int
	}{
		{"PlanZero", InitAndPlan, 0, 0, 0, 0},
		{"ApplyZero", InitAndApply, 0, 0, 0, 0},
		{"PlanAddResouce", InitAndPlan, 2, 2, 0, 0},
		{"ApplyAddResouce", InitAndApply, 2, 2, 0, 0},
		{"PlanNoOp", InitAndApply, 2, 0, 0, 0},
		{"ApplyNoOp", InitAndApply, 2, 0, 0, 0},
		{"PlanDestroyResource", InitAndPlan, 1, 0, 0, 1},
		{"ApplyDestroyResource", InitAndApply, 1, 0, 0, 1},
		{"Destroy", Destroy, 1, 0, 0, 1},
		{"DestroyNoOp", Destroy, 1, 0, 0, 0},
	}

	testFolder, err := files.CopyTerraformFolderToTemp("../../test/fixtures/terraform-basic-configuration", t.Name())
	require.NoError(t, err)

	terraformOptions := &Options{
		TerraformDir: testFolder,
		Vars: map[string]interface{}{
			"cnt": 0,
		},
		NoColor: noColor,
	}

	for _, tc := range testCases {
		t.Run(tc.Name,
			func(t *testing.T) {
				terraformOptions.Vars["cnt"] = tc.cntValue
				cnt, err := GetResourceCountE(t, tc.tfFuncToRun(t, terraformOptions))
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedAdd, cnt.Add)
				assert.Equal(t, tc.expectedChange, cnt.Change)
				assert.Equal(t, tc.expectedDestroy, cnt.Destroy)
			})
	}

	t.Run("InvalidInput",
		func(t *testing.T) {
			terraformOptions.Vars["cnt"] = "abc"
			cmdout, _ := PlanE(t, terraformOptions)
			cnt, err := GetResourceCountE(t, cmdout)
			assert.EqualError(t, err, getResourceCountErrMessage)
			assert.Nil(t, cnt)
		})

}
