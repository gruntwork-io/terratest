package test

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/gruntwork-io/terratest/modules/aws"
	"github.com/gruntwork-io/terratest/modules/packer"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	test_structure "github.com/gruntwork-io/terratest/modules/test-structure"
)

func TestWindowsInstance(t *testing.T) {
	// Uncomment any of the following to skip that section during the test
	//os.Setenv("SKIP_setup", "true")
	//os.Setenv("SKIP_build_ami", "true")
	//os.Setenv("SKIP_deploy", "true")
	//os.Setenv("SKIP_validate", "true")
	//os.Setenv("SKIP_cleanup", "true")

	workingDir := filepath.Join(".", "stages", t.Name())
	testBasePath := test_structure.CopyTerraformFolderToTemp(t, "..", "examples/terraform-aws-ec2-windows-example")

	test_structure.RunTestStage(t, "setup", func() {
		uniqueID := random.UniqueId()
		region := aws.GetRandomRegion(t, []string{}, []string{})
		roleName := fmt.Sprintf("%s-test-role", uniqueID)

		instanceType := aws.GetRecommendedInstanceType(t, region, []string{"t2.micro, t3.micro", "t2.small", "t3.small"})
		test_structure.SaveString(t, workingDir, "region", region)
		test_structure.SaveString(t, workingDir, "uniqueID", uniqueID)
		test_structure.SaveString(t, workingDir, "instanceType", instanceType)
		test_structure.SaveString(t, workingDir, "roleName", roleName)
	})

	test_structure.RunTestStage(t, "build_ami", func() {
		region := test_structure.LoadString(t, workingDir, "region")
		instanceType := test_structure.LoadString(t, workingDir, "instanceType")
		roleName := test_structure.LoadString(t, workingDir, "roleName")

		varsMap := make(map[string]string)

		varsMap["instance_type"] = instanceType
		varsMap["region"] = region
		packerOptions := &packer.Options{
			Template: filepath.Join(testBasePath, "packer/build.pkr.hcl"),
			Vars:     varsMap,
		}

		amiID := packer.BuildArtifact(t, packerOptions)

		test_structure.SaveString(t, workingDir, "amiID", amiID)

		terratestOptions := &terraform.Options{
			TerraformDir: testBasePath,
			Vars:         make(map[string]interface{}),
		}

		terratestOptions.Vars["ami"] = amiID
		terratestOptions.Vars["region"] = region
		terratestOptions.Vars["iam_role_name"] = roleName
		test_structure.SaveTerraformOptions(t, workingDir, terratestOptions)
	})

	defer test_structure.RunTestStage(t, "cleanup", func() {
		terratestOptions := test_structure.LoadTerraformOptions(t, workingDir)
		terraform.Destroy(t, terratestOptions)
	})

	test_structure.RunTestStage(t, "deploy", func() {
		terratestOptions := test_structure.LoadTerraformOptions(t, workingDir)
		terraform.InitAndApply(t, terratestOptions)
	})

}
