//go:build aws

package test_test

import (
	"path/filepath"
	"testing"

	"github.com/gruntwork-io/terratest/modules/aws"
	"github.com/gruntwork-io/terratest/modules/core/v2/random"
	"github.com/gruntwork-io/terratest/modules/packer"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/gruntwork-io/terratest/modules/teststructure"
)

func TestWindowsInstance(t *testing.T) {
	t.Parallel()

	// Uncomment any of the following to skip that section during the test
	// os.Setenv("SKIP_setup", "true")
	// os.Setenv("SKIP_build_ami", "true")
	// os.Setenv("SKIP_deploy", "true")
	// os.Setenv("SKIP_validate", "true")
	// os.Setenv("SKIP_cleanup", "true")

	workingDir := filepath.Join(".", "stages", t.Name())
	testBasePath := teststructure.CopyTerraformFolderToTemp(t, "..", "examples/terraform-aws-ec2-windows-example")

	teststructure.RunTestStage(t, "setup", func() {
		ctx := t.Context()
		uniqueID := random.UniqueID()
		region := aws.GetRandomRegionContext(t, ctx, []string{}, []string{})
		roleName := uniqueID + "-test-role"

		instanceType := aws.GetRecommendedInstanceTypeContext(t, ctx, region, []string{"t2.micro, t3.micro", "t2.small", "t3.small"})
		teststructure.SaveString(t, workingDir, "region", region)
		teststructure.SaveString(t, workingDir, "uniqueID", uniqueID)
		teststructure.SaveString(t, workingDir, "instanceType", instanceType)
		teststructure.SaveString(t, workingDir, "roleName", roleName)
	})

	teststructure.RunTestStage(t, "build_ami", func() {
		region := teststructure.LoadString(t, workingDir, "region")
		instanceType := teststructure.LoadString(t, workingDir, "instanceType")
		roleName := teststructure.LoadString(t, workingDir, "roleName")

		varsMap := make(map[string]string)

		varsMap["instance_type"] = instanceType
		varsMap["region"] = region
		packerOptions := &packer.Options{
			Template: filepath.Join(testBasePath, "packer/build.pkr.hcl"),
			Vars:     varsMap,
		}

		amiID := packer.BuildArtifactContext(t, t.Context(), packerOptions)

		teststructure.SaveString(t, workingDir, "amiID", amiID)

		terratestOptions := &terraform.Options{
			TerraformDir: testBasePath,
			Vars:         make(map[string]interface{}),
		}

		terratestOptions.Vars["ami"] = amiID
		terratestOptions.Vars["region"] = region
		terratestOptions.Vars["iam_role_name"] = roleName
		teststructure.SaveTerraformOptions(t, workingDir, terratestOptions)
	})

	defer teststructure.RunTestStage(t, "cleanup", func() {
		terratestOptions := teststructure.LoadTerraformOptions(t, workingDir)
		terraform.DestroyContext(t, t.Context(), terratestOptions)
	})

	teststructure.RunTestStage(t, "deploy", func() {
		terratestOptions := teststructure.LoadTerraformOptions(t, workingDir)
		terraform.InitAndApplyContext(t, t.Context(), terratestOptions)
	})
}
