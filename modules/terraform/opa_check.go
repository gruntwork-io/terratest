package terraform

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/stretchr/testify/require"
	"github.com/tmccombs/hcl2json/convert"

	"github.com/gruntwork-io/terratest/modules/files"
	"github.com/gruntwork-io/terratest/modules/opa"
	"github.com/gruntwork-io/terratest/modules/testing"
)

// OPAEval runs `opa eval` with the given option on the terraform files identified in the TerraformDir directory of the
// Options struct. Note that since OPA does not natively support parsing HCL code, we first convert all the files to
// JSON prior to passing it through OPA. This function fails the test if there is an error.
func OPAEval(
	t testing.TestingT,
	tfOptions *Options,
	opaEvalOptions *opa.EvalOptions,
	resultQuery string,
) {
	require.NoError(t, OPAEvalE(t, tfOptions, opaEvalOptions, resultQuery))
}

// OPAEvalE runs `opa eval` with the given option on the terraform files identified in the TerraformDir directory of the
// Options struct. Note that since OPA does not natively support parsing HCL code, we first convert all the files to
// JSON prior to passing it through OPA.
func OPAEvalE(
	t testing.TestingT,
	tfOptions *Options,
	opaEvalOptions *opa.EvalOptions,
	resultQuery string,
) error {
	tfOptions.Logger.Logf(t, "Running terraform files in %s through `opa eval` on policy %s", tfOptions.TerraformDir, opaEvalOptions.RulePath)

	// Find all the tf files in the terraform dir to process.
	tfFiles, err := files.FindTerraformSourceFilesInDir(tfOptions.TerraformDir)
	if err != nil {
		return err
	}

	// Create a temporary dir to store all the json files
	tmpDir, err := os.MkdirTemp("", "terratest-opa-hcl2json-*")
	if err != nil {
		return err
	}
	if !opaEvalOptions.DebugKeepTempFiles {
		defer os.RemoveAll(tmpDir)
	}
	tfOptions.Logger.Logf(t, "Using temporary folder %s for json representation of terraform module %s", tmpDir, tfOptions.TerraformDir)

	// Convert all the found tf files to json format so OPA works.
	jsonFiles := make([]string, len(tfFiles))
	errorsOccurred := new(multierror.Error)
	for i, tfFile := range tfFiles {
		tfFileBase := filepath.Base(tfFile)
		tfFileBaseName := strings.TrimSuffix(tfFileBase, filepath.Ext(tfFileBase))
		outPath := filepath.Join(tmpDir, tfFileBaseName+".json")
		tfOptions.Logger.Logf(t, "Converting %s to json %s", tfFile, outPath)
		if err := HCLFileToJSONFile(tfFile, outPath); err != nil {
			errorsOccurred = multierror.Append(errorsOccurred, err)
		}
		jsonFiles[i] = outPath
	}
	if err := errorsOccurred.ErrorOrNil(); err != nil {
		return err
	}

	// Run OPA checks on each of the converted json files.
	return opa.EvalE(t, opaEvalOptions, jsonFiles, resultQuery)
}

// HCLFileToJSONFile is a function that takes a path containing HCL code, and converts it to JSON representation and
// writes out the contents to the given path.
func HCLFileToJSONFile(hclPath, jsonOutPath string) error {
	fileBytes, err := os.ReadFile(hclPath)
	if err != nil {
		return err
	}
	converted, err := convert.Bytes(fileBytes, hclPath, convert.Options{})
	if err != nil {
		return err
	}
	return os.WriteFile(jsonOutPath, converted, 0600)
}
