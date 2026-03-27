package terraform

import (
	"context"
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
	_ = ctx // reserved for future use when opa.EvalE supports context

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
		defer func() { _ = os.RemoveAll(tmpDir) }()
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
	return opa.EvalE(t, opaEvalOptions, jsonFiles, resultQuery) //nolint:contextcheck // opa.EvalE does not have a context variant yet
}

// OPAEval runs `opa eval` with the given option on the terraform files identified in the TerraformDir directory of the
// Options struct. Note that since OPA does not natively support parsing HCL code, we first convert all the files to
// JSON prior to passing it through OPA. This function fails the test if there is an error.
//
// Deprecated: Use [OPAEvalContext] instead.
func OPAEval(
	t testing.TestingT,
	tfOptions *Options,
	opaEvalOptions *opa.EvalOptions,
	resultQuery string,
) {
	OPAEvalContext(t, context.Background(), tfOptions, opaEvalOptions, resultQuery)
}

// OPAEvalE runs `opa eval` with the given option on the terraform files identified in the TerraformDir directory of the
// Options struct. Note that since OPA does not natively support parsing HCL code, we first convert all the files to
// JSON prior to passing it through OPA.
//
// Deprecated: Use [OPAEvalContextE] instead.
func OPAEvalE(
	t testing.TestingT,
	tfOptions *Options,
	opaEvalOptions *opa.EvalOptions,
	resultQuery string,
) error {
	return OPAEvalContextE(t, context.Background(), tfOptions, opaEvalOptions, resultQuery)
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

	return os.WriteFile(jsonOutPath, converted, 0o600) //nolint:mnd // standard owner read-write permission
}
