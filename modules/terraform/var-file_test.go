package terraform_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/require"
)

func TestGetVariablesFromVarFilesAsString(t *testing.T) {
	t.Parallel()

	randomFileName := fmt.Sprintf("./%s.tfvars", random.UniqueID())

	testHcl := []byte(`
		aws_region     = "us-east-2"
		aws_account_id = "111111111111"
		number_type = 2
		boolean_type = true
		tags = {
			foo = "bar"
		}
		list = ["item1"]`)

	writeFile(t, randomFileName, testHcl)
	defer os.Remove(randomFileName)

	stringVal := terraform.GetVariableAsStringFromVarFile(t, randomFileName, "aws_region")

	boolString := terraform.GetVariableAsStringFromVarFile(t, randomFileName, "boolean_type")

	numString := terraform.GetVariableAsStringFromVarFile(t, randomFileName, "number_type")

	require.Equal(t, "us-east-2", stringVal)
	require.Equal(t, "true", boolString)
	require.Equal(t, "2", numString)
}

func TestGetVariablesFromVarFilesAsStringKeyDoesNotExist(t *testing.T) {
	t.Parallel()

	randomFileName := fmt.Sprintf("./%s.tfvars", random.UniqueID())

	testHcl := []byte(`
		aws_region     = "us-east-2"
		aws_account_id = "111111111111"
		tags = {
			foo = "bar"
		}
		list = ["item1"]`)

	writeFile(t, randomFileName, testHcl)
	defer os.Remove(randomFileName)

	_, err := terraform.GetVariableAsStringFromVarFileE(t, randomFileName, "badkey")

	require.Error(t, err)
}

func TestGetVariableAsMapFromVarFile(t *testing.T) {
	t.Parallel()

	randomFileName := fmt.Sprintf("./%s.tfvars", random.UniqueID())
	expected := make(map[string]string)
	expected["foo"] = "bar"

	testHcl := []byte(`
		aws_region     = "us-east-2"
		aws_account_id = "111111111111"
		tags = {
			foo = "bar"
		}
		list = ["item1"]`)

	writeFile(t, randomFileName, testHcl)
	defer os.Remove(randomFileName)

	val := terraform.GetVariableAsMapFromVarFile(t, randomFileName, "tags")
	require.Equal(t, expected, val)
}

func TestGetVariableAsMapFromVarFileNotMap(t *testing.T) {
	t.Parallel()

	randomFileName := fmt.Sprintf("./%s.tfvars", random.UniqueID())

	testHcl := []byte(`
		aws_region     = "us-east-2"
		aws_account_id = "111111111111"
		tags = {
			foo = "bar"
		}
		list = ["item1"]`)

	writeFile(t, randomFileName, testHcl)
	defer os.Remove(randomFileName)

	_, err := terraform.GetVariableAsMapFromVarFileE(t, randomFileName, "aws_region")

	require.Error(t, err)
}

func TestGetVariableAsMapFromVarFileKeyDoesNotExist(t *testing.T) {
	t.Parallel()

	randomFileName := fmt.Sprintf("./%s.tfvars", random.UniqueID())

	testHcl := []byte(`
		aws_region     = "us-east-2"
		aws_account_id = "111111111111"
		tags = {
			foo = "bar"
		}
		list = ["item1"]`)

	writeFile(t, randomFileName, testHcl)
	defer os.Remove(randomFileName)

	_, err := terraform.GetVariableAsMapFromVarFileE(t, randomFileName, "badkey")

	require.Error(t, err)
}

func TestGetVariableAsListFromVarFile(t *testing.T) {
	t.Parallel()

	randomFileName := fmt.Sprintf("./%s.tfvars", random.UniqueID())
	expected := []string{"item1"}

	testHcl := []byte(`
		aws_region     = "us-east-2"
		aws_account_id = "111111111111"
		tags = {
			foo = "bar"
		}
		list = ["item1"]`)

	writeFile(t, randomFileName, testHcl)
	defer os.Remove(randomFileName)

	val := terraform.GetVariableAsListFromVarFile(t, randomFileName, "list")

	require.Equal(t, expected, val)
}

func TestGetVariableAsListNotList(t *testing.T) {
	t.Parallel()

	randomFileName := fmt.Sprintf("./%s.tfvars", random.UniqueID())

	testHcl := []byte(`
		aws_region     = "us-east-2"
		aws_account_id = "111111111111"
		tags = {
			foo = "bar"
		}
		list = ["item1"]`)

	writeFile(t, randomFileName, testHcl)
	defer os.Remove(randomFileName)

	_, err := terraform.GetVariableAsListFromVarFileE(t, randomFileName, "tags")

	require.Error(t, err)
}

func TestGetVariableAsListKeyDoesNotExist(t *testing.T) {
	t.Parallel()

	randomFileName := fmt.Sprintf("./%s.tfvars", random.UniqueID())

	testHcl := []byte(`
		aws_region     = "us-east-2"
		aws_account_id = "111111111111"
		tags = {
			foo = "bar"
		}
		list = ["item1"]`)

	writeFile(t, randomFileName, testHcl)
	defer os.Remove(randomFileName)

	_, err := terraform.GetVariableAsListFromVarFileE(t, randomFileName, "badkey")

	require.Error(t, err)
}

func TestGetAllVariablesFromVarFileEFileDoesNotExist(t *testing.T) {
	t.Parallel()

	var variables map[string]any

	err := terraform.GetAllVariablesFromVarFileE(t, "filea", &variables)
	require.Equal(t, "open filea: no such file or directory", err.Error())
}

func TestGetAllVariablesFromVarFileBadFile(t *testing.T) {
	t.Parallel()

	randomFileName := fmt.Sprintf("./%s.tfvars", random.UniqueID())
	testHcl := []byte(`
		thiswillnotwork`)

	writeFile(t, randomFileName, testHcl)
	defer os.Remove(randomFileName)

	var variables map[string]any

	err := terraform.GetAllVariablesFromVarFileE(t, randomFileName, &variables)
	require.Error(t, err)

	// HCL library could change their error string, so we are only testing the error string contains what we add to it
	require.Regexp(t, fmt.Sprintf("^%s:2,3-18: ", randomFileName), err.Error())
}

func TestGetAllVariablesFromVarFile(t *testing.T) {
	t.Parallel()

	randomFileName := fmt.Sprintf("./%s.tfvars", random.UniqueID())
	testHcl := []byte(`
	aws_region     = "us-east-2"
	`)

	writeFile(t, randomFileName, testHcl)
	defer os.Remove(randomFileName)

	var variables map[string]any

	err := terraform.GetAllVariablesFromVarFileE(t, randomFileName, &variables)
	require.NoError(t, err)

	expected := make(map[string]any)
	expected["aws_region"] = "us-east-2"

	require.Equal(t, expected, variables)
}

func TestGetAllVariablesFromVarFileStructOut(t *testing.T) {
	t.Parallel()

	randomFileName := fmt.Sprintf("./%s.tfvars", random.UniqueID())
	testHcl := []byte(`
	aws_region     = "us-east-2"
	`)

	writeFile(t, randomFileName, testHcl)
	defer os.Remove(randomFileName)

	var region struct {
		AwsRegion string `cty:"aws_region"`
	}

	err := terraform.GetAllVariablesFromVarFileE(t, randomFileName, &region)
	require.NoError(t, err)
	require.Equal(t, "us-east-2", region.AwsRegion)
}

func TestGetVariablesFromVarFilesAsStringJSON(t *testing.T) {
	t.Parallel()

	randomFileName := fmt.Sprintf("./%s.tfvars.json", random.UniqueID())

	testJSON := []byte(`
		{
			"aws_region": "us-east-2",
			"aws_account_id": "111111111111",
			"number_type": 2,
			"boolean_type": true,
			"tags": {
				"foo": "bar"
			},
			"list": ["item1"]
		}`)

	writeFile(t, randomFileName, testJSON)
	defer os.Remove(randomFileName)

	stringVal := terraform.GetVariableAsStringFromVarFile(t, randomFileName, "aws_region")

	boolString := terraform.GetVariableAsStringFromVarFile(t, randomFileName, "boolean_type")

	numString := terraform.GetVariableAsStringFromVarFile(t, randomFileName, "number_type")

	require.Equal(t, "us-east-2", stringVal)
	require.Equal(t, "true", boolString)
	require.Equal(t, "2", numString)
}

func TestGetVariablesFromVarFilesAsStringKeyDoesNotExistJSON(t *testing.T) {
	t.Parallel()

	randomFileName := fmt.Sprintf("./%s.tfvars.json", random.UniqueID())

	testJSON := []byte(`
		{
			"aws_region": "us-east-2",
			"aws_account_id": "111111111111",
			"tags": {
				"foo": "bar"
			},
			"list": ["item1"]
		}`)

	writeFile(t, randomFileName, testJSON)
	defer os.Remove(randomFileName)

	_, err := terraform.GetVariableAsStringFromVarFileE(t, randomFileName, "badkey")

	require.Error(t, err)
}

func TestGetVariableAsMapFromVarFileJSON(t *testing.T) {
	t.Parallel()

	randomFileName := fmt.Sprintf("./%s.tfvars.json", random.UniqueID())
	expected := make(map[string]string)
	expected["foo"] = "bar"

	testJSON := []byte(`
		{
			"aws_region": "us-east-2",
			"aws_account_id": "111111111111",
			"tags": {
				"foo": "bar"
			},
			"list": ["item1"]
		}`)

	writeFile(t, randomFileName, testJSON)
	defer os.Remove(randomFileName)

	val := terraform.GetVariableAsMapFromVarFile(t, randomFileName, "tags")
	require.Equal(t, expected, val)
}

func TestGetVariableAsMapFromVarFileNotMapJSON(t *testing.T) {
	t.Parallel()

	randomFileName := fmt.Sprintf("./%s.tfvars.json", random.UniqueID())

	testJSON := []byte(`
		{
			"aws_region": "us-east-2",
			"aws_account_id": "111111111111",
			"tags": {
				"foo": "bar"
			},
			"list": ["item1"]
		}`)

	writeFile(t, randomFileName, testJSON)
	defer os.Remove(randomFileName)

	_, err := terraform.GetVariableAsMapFromVarFileE(t, randomFileName, "aws_region")

	require.Error(t, err)
}

func TestGetVariableAsMapFromVarFileKeyDoesNotExistJSON(t *testing.T) {
	t.Parallel()

	randomFileName := fmt.Sprintf("./%s.tfvars.json", random.UniqueID())

	testJSON := []byte(`
		{
			"aws_region": "us-east-2",
			"aws_account_id": "111111111111",
			"tags": {
				"foo": "bar"
			},
			"list": ["item1"]
		}`)

	writeFile(t, randomFileName, testJSON)
	defer os.Remove(randomFileName)

	_, err := terraform.GetVariableAsMapFromVarFileE(t, randomFileName, "badkey")

	require.Error(t, err)
}

func TestGetVariableAsListFromVarFileJSON(t *testing.T) {
	t.Parallel()

	randomFileName := fmt.Sprintf("./%s.tfvars.json", random.UniqueID())
	expected := []string{"item1"}

	testJSON := []byte(`
		{
			"aws_region": "us-east-2",
			"aws_account_id": "111111111111",
			"tags": {
				"foo": "bar"
			},
			"list": ["item1"]
		}`)

	writeFile(t, randomFileName, testJSON)
	defer os.Remove(randomFileName)

	val := terraform.GetVariableAsListFromVarFile(t, randomFileName, "list")

	require.Equal(t, expected, val)
}

func TestGetVariableAsListNotListJSON(t *testing.T) {
	t.Parallel()

	randomFileName := fmt.Sprintf("./%s.tfvars.json", random.UniqueID())

	testJSON := []byte(`
		{
			"aws_region": "us-east-2",
			"aws_account_id": "111111111111",
			"tags": {
				"foo": "bar"
			},
			"list": ["item1"]
		}`)

	writeFile(t, randomFileName, testJSON)
	defer os.Remove(randomFileName)

	_, err := terraform.GetVariableAsListFromVarFileE(t, randomFileName, "tags")

	require.Error(t, err)
}

func TestGetVariableAsListKeyDoesNotExistJSON(t *testing.T) {
	t.Parallel()

	randomFileName := fmt.Sprintf("./%s.tfvars.json", random.UniqueID())

	testJSON := []byte(`
		{
			"aws_region": "us-east-2",
			"aws_account_id": "111111111111",
			"tags": {
				"foo": "bar"
			},
			"list": ["item1"]
		}`)

	writeFile(t, randomFileName, testJSON)
	defer os.Remove(randomFileName)

	_, err := terraform.GetVariableAsListFromVarFileE(t, randomFileName, "badkey")

	require.Error(t, err)
}

func TestGetAllVariablesFromVarFileBadFileJSON(t *testing.T) {
	t.Parallel()

	randomFileName := fmt.Sprintf("./%s.tfvars.json", random.UniqueID())
	testJSON := []byte(`
		{
			thiswillnotwork
		}`)

	writeFile(t, randomFileName, testJSON)
	defer os.Remove(randomFileName)

	var variables map[string]any

	err := terraform.GetAllVariablesFromVarFileE(t, randomFileName, &variables)
	require.Error(t, err)

	// HCL library could change their error string, so we are only testing the error string contains what we add to it
	require.Regexp(t, fmt.Sprintf("^%s:3,7-22: ", randomFileName), err.Error())
}

func TestGetAllVariablesFromVarFileJSON(t *testing.T) {
	t.Parallel()

	randomFileName := fmt.Sprintf("./%s.tfvars.json", random.UniqueID())
	testJSON := []byte(`
	{
		"aws_region": "us-east-2"
	}
	`)

	writeFile(t, randomFileName, testJSON)
	defer os.Remove(randomFileName)

	var variables map[string]any

	err := terraform.GetAllVariablesFromVarFileE(t, randomFileName, &variables)
	require.NoError(t, err)

	expected := make(map[string]any)
	expected["aws_region"] = "us-east-2"

	require.Equal(t, expected, variables)
}

func TestGetAllVariablesFromVarFileStructOutJSON(t *testing.T) {
	t.Parallel()

	randomFileName := fmt.Sprintf("./%s.tfvars.json", random.UniqueID())
	testJSON := []byte(`
	{
		"aws_region": "us-east-2"
	}
	`)

	writeFile(t, randomFileName, testJSON)
	defer os.Remove(randomFileName)

	var region struct {
		AwsRegion string `cty:"aws_region"`
	}

	err := terraform.GetAllVariablesFromVarFileE(t, randomFileName, &region)
	require.NoError(t, err)
	require.Equal(t, "us-east-2", region.AwsRegion)
}

// writeFile is a helper function to write a file to the filesystem.
// It will immediately fail the test if it could not write the file.
func writeFile(t *testing.T, fileName string, bytes []byte) {
	t.Helper()

	err := os.WriteFile(fileName, bytes, 0644)
	require.NoError(t, err)
}
