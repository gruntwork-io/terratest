package terraform

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/stretchr/testify/require"
	"github.com/tmccombs/hcl2json/convert"

	"github.com/gruntwork-io/terratest/modules/core/v2/files"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/gruntwork-io/terratest/modules/opa/v2"
)

// OPAEvalContext runs `opa eval` with the given options on the terraform files identified in the TerraformDir
// directory of the Options struct. Note that since OPA does not natively support parsing HCL code, we first
// convert all the files to JSON prior to passing it through OPA. The context argument can be used for
// cancellation or timeout control. This function fails the test if there is an error.
func OPAEvalContext(
	t testing.TestingT,
	ctx context.Context,
	tfOptions *Options,
	opaEvalOptions *opa.EvalOptions,
	resultQuery string,
) {
	require.NoError(t, OPAEvalContextE(t, ctx, tfOptions, opaEvalOptions, resultQuery))
}

// OPAEvalContextE runs `opa eval` with the given options on the terraform files identified in the TerraformDir
// directory of the Options struct. Note that since OPA does not natively support parsing HCL code, we first
// convert all the files to JSON prior to passing it through OPA. The context argument can be used for
// cancellation or timeout control.
func OPAEvalContextE(
	t testing.TestingT,
	ctx context.Context,
	tfOptions *Options,
	opaEvalOptions *opa.EvalOptions,
	resultQuery string,
) error {
	_ = ctx

	tfOptions.Logger.Logf(t, "Running terraform files in %s through `opa eval` on policy %s", tfOptions.TerraformDir, opaEvalOptions.RulePath)

	tfFiles, err := files.FindTerraformSourceFilesInDir(tfOptions.TerraformDir)
	if err != nil {
		return err
	}

	tmpDir, err := os.MkdirTemp("", "terratest-opa-hcl2json-*")
	if err != nil {
		return err
	}

	if !opaEvalOptions.DebugKeepTempFiles {
		defer func() { _ = os.RemoveAll(tmpDir) }()
	}

	tfOptions.Logger.Logf(t, "Using temporary folder %s for json representation of terraform module %s", tmpDir, tfOptions.TerraformDir)

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

	return os.WriteFile(jsonOutPath, converted, 0o600)
}
