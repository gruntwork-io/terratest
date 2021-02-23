package terraform

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/stretchr/testify/require"
)

func TestGetVariablesFromVarFilesAsString(t *testing.T) {
	randomFileName := fmt.Sprintf("./%s.tfvars", random.UniqueId())

	testHcl := []byte(`
		aws_region     = "us-east-2"
		aws_account_id = "111111111111"
		number_type = 2
		boolean_type = true
		tags = {
			foo = "bar"
		}
		list = ["item1"]`)

	WriteFile(t, randomFileName, testHcl)
	defer os.Remove(randomFileName)

	stringVal := GetVariableAsStringFromVarFile(t, randomFileName, "aws_region")

	boolString := GetVariableAsStringFromVarFile(t, randomFileName, "boolean_type")

	numString := GetVariableAsStringFromVarFile(t, randomFileName, "number_type")

	require.Equal(t, "us-east-2", stringVal)
	require.Equal(t, "true", boolString)
	require.Equal(t, "2", numString)

}

func TestGetVariablesFromVarFilesAsStringKeyDoesNotExist(t *testing.T) {
	randomFileName := fmt.Sprintf("./%s.tfvars", random.UniqueId())

	testHcl := []byte(`
		aws_region     = "us-east-2"
		aws_account_id = "111111111111"
		tags = {
			foo = "bar"
		}
		list = ["item1"]`)

	WriteFile(t, randomFileName, testHcl)
	defer os.Remove(randomFileName)

	_, err := GetVariableAsStringFromVarFileE(t, randomFileName, "badkey")

	require.Error(t, err)
}

func TestGetVariableAsMapFromVarFile(t *testing.T) {
	randomFileName := fmt.Sprintf("./%s.tfvars", random.UniqueId())
	expected := make(map[string]string)
	expected["foo"] = "bar"

	testHcl := []byte(`
		aws_region     = "us-east-2"
		aws_account_id = "111111111111"
		tags = {
			foo = "bar"
		}
		list = ["item1"]`)

	WriteFile(t, randomFileName, testHcl)
	defer os.Remove(randomFileName)

	val := GetVariableAsMapFromVarFile(t, randomFileName, "tags")
	require.Equal(t, expected, val)
}

func TestGetVariableAsMapFromVarFileNotMap(t *testing.T) {
	randomFileName := fmt.Sprintf("./%s.tfvars", random.UniqueId())

	testHcl := []byte(`
		aws_region     = "us-east-2"
		aws_account_id = "111111111111"
		tags = {
			foo = "bar"
		}
		list = ["item1"]`)

	WriteFile(t, randomFileName, testHcl)
	defer os.Remove(randomFileName)

	_, err := GetVariableAsMapFromVarFileE(t, randomFileName, "aws_region")

	require.Error(t, err)
}

func TestGetVariableAsMapFromVarFileKeyDoesNotExist(t *testing.T) {
	randomFileName := fmt.Sprintf("./%s.tfvars", random.UniqueId())

	testHcl := []byte(`
		aws_region     = "us-east-2"
		aws_account_id = "111111111111"
		tags = {
			foo = "bar"
		}
		list = ["item1"]`)

	WriteFile(t, randomFileName, testHcl)
	defer os.Remove(randomFileName)

	_, err := GetVariableAsMapFromVarFileE(t, randomFileName, "badkey")

	require.Error(t, err)
}

func TestGetVariableAsListFromVarFile(t *testing.T) {
	randomFileName := fmt.Sprintf("./%s.tfvars", random.UniqueId())
	expected := []string{"item1"}

	testHcl := []byte(`
		aws_region     = "us-east-2"
		aws_account_id = "111111111111"
		tags = {
			foo = "bar"
		}
		list = ["item1"]`)

	WriteFile(t, randomFileName, testHcl)
	defer os.Remove(randomFileName)

	val := GetVariableAsListFromVarFile(t, randomFileName, "list")

	require.Equal(t, expected, val)
}

func TestGetVariableAsListNotList(t *testing.T) {
	randomFileName := fmt.Sprintf("./%s.tfvars", random.UniqueId())

	testHcl := []byte(`
		aws_region     = "us-east-2"
		aws_account_id = "111111111111"
		tags = {
			foo = "bar"
		}
		list = ["item1"]`)

	WriteFile(t, randomFileName, testHcl)
	defer os.Remove(randomFileName)

	_, err := GetVariableAsListFromVarFileE(t, randomFileName, "tags")

	require.Error(t, err)
}

func TestGetVariableAsListKeyDoesNotExist(t *testing.T) {
	randomFileName := fmt.Sprintf("./%s.tfvars", random.UniqueId())

	testHcl := []byte(`
		aws_region     = "us-east-2"
		aws_account_id = "111111111111"
		tags = {
			foo = "bar"
		}
		list = ["item1"]`)

	WriteFile(t, randomFileName, testHcl)
	defer os.Remove(randomFileName)

	_, err := GetVariableAsListFromVarFileE(t, randomFileName, "badkey")

	require.Error(t, err)
}
func TestGetAllVariablesFromVarFileEFileDoesNotExist(t *testing.T) {
	var variables map[string]interface{}
	err := GetAllVariablesFromVarFileE(t, "filea", &variables)
	require.Equal(t, "open filea: no such file or directory", err.Error())
}

func TestGetAllVariablesFromVarFileBadFile(t *testing.T) {
	randomFileName := fmt.Sprintf("./%s.tfvars", random.UniqueId())
	testHcl := []byte(`
		thiswillnotwork`)

	WriteFile(t, randomFileName, testHcl)
	defer os.Remove(randomFileName)

	var variables map[string]interface{}
	err := GetAllVariablesFromVarFileE(t, randomFileName, &variables)
	require.Error(t, err)

	// HCL library could change their error string, so we are only testing the error string contains what we add to it
	require.Regexp(t, fmt.Sprintf("^%s:2,3-18: ", randomFileName), err.Error())
}

func TestGetAllVariablesFromVarFile(t *testing.T) {
	randomFileName := fmt.Sprintf("./%s.tfvars", random.UniqueId())
	testHcl := []byte(`
	aws_region     = "us-east-2"
	`)

	WriteFile(t, randomFileName, testHcl)
	defer os.Remove(randomFileName)

	var variables map[string]interface{}
	err := GetAllVariablesFromVarFileE(t, randomFileName, &variables)
	require.NoError(t, err)

	expected := make(map[string]interface{})
	expected["aws_region"] = "us-east-2"

	require.Equal(t, expected, variables)
}

func TestGetAllVariablesFromVarFileStructOut(t *testing.T) {
	randomFileName := fmt.Sprintf("./%s.tfvars", random.UniqueId())
	testHcl := []byte(`
	aws_region     = "us-east-2"
	`)

	WriteFile(t, randomFileName, testHcl)
	defer os.Remove(randomFileName)

	var region struct {
		AwsRegion string `cty:"aws_region"`
	}
	err := GetAllVariablesFromVarFileE(t, randomFileName, &region)
	require.NoError(t, err)
	require.Equal(t, "us-east-2", region.AwsRegion)

}

// Helper function to write a file to the filesystem
// Will immediately fail the test if it could not write the file
func WriteFile(t *testing.T, fileName string, bytes []byte) {
	err := ioutil.WriteFile(fileName, bytes, 0644)

	require.NoError(t, err)
}
