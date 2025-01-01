package test_structure

import (
	"os"
	"testing"

	"github.com/gruntwork-io/terratest/modules/files"
	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

type testData struct {
	Foo string
	Bar bool
	Baz map[string]interface{}
}

func TestSaveAndLoadTestData(t *testing.T) {
	t.Parallel()

	isTestDataPresent := IsTestDataPresent(t, "/file/that/does/not/exist")
	assert.False(t, isTestDataPresent, "Expected no test data would be present because no test data file exists.")

	tmpFile, err := os.CreateTemp("", "save-and-load-test-data")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	expectedData := testData{
		Foo: "foo",
		Bar: true,
		Baz: map[string]interface{}{"abc": "def", "ghi": 1.0, "klm": false},
	}

	isTestDataPresent = IsTestDataPresent(t, tmpFile.Name())
	assert.False(t, isTestDataPresent, "Expected no test data would be present because file exists but no data has been written yet.")

	overwrite := true
	SaveTestData(t, tmpFile.Name(), overwrite, expectedData)

	isTestDataPresent = IsTestDataPresent(t, tmpFile.Name())
	assert.True(t, isTestDataPresent, "Expected test data would be present because file exists and data has been written to file.")

	actualData := testData{}
	LoadTestData(t, tmpFile.Name(), &actualData)
	assert.Equal(t, expectedData, actualData)

	overwritingData := testData{
		Foo: "foo",
		Bar: false,
		Baz: map[string]interface{}{"123": "456", "789": 1.0, "0": false},
	}
	SaveTestData(t, tmpFile.Name(), !overwrite, overwritingData)
	LoadTestData(t, tmpFile.Name(), &actualData)
	assert.Equal(t, expectedData, actualData)

	CleanupTestData(t, tmpFile.Name())
	assert.False(t, files.FileExists(tmpFile.Name()))
}

func TestIsEmptyJson(t *testing.T) {
	t.Parallel()

	var jsonValue []byte
	var isEmpty bool

	jsonValue = []byte("null")
	isEmpty = isEmptyJSON(t, jsonValue)
	assert.True(t, isEmpty, `The JSON literal "null" should be treated as an empty value.`)

	jsonValue = []byte("false")
	isEmpty = isEmptyJSON(t, jsonValue)
	assert.True(t, isEmpty, `The JSON literal "false" should be treated as an empty value.`)

	jsonValue = []byte("true")
	isEmpty = isEmptyJSON(t, jsonValue)
	assert.False(t, isEmpty, `The JSON literal "true" should be treated as a non-empty value.`)

	jsonValue = []byte("0")
	isEmpty = isEmptyJSON(t, jsonValue)
	assert.True(t, isEmpty, `The JSON literal "0" should be treated as an empty value.`)

	jsonValue = []byte("1")
	isEmpty = isEmptyJSON(t, jsonValue)
	assert.False(t, isEmpty, `The JSON literal "1" should be treated as a non-empty value.`)

	jsonValue = []byte("{}")
	isEmpty = isEmptyJSON(t, jsonValue)
	assert.True(t, isEmpty, `The JSON value "{}" should be treated as an empty value.`)

	jsonValue = []byte(`{ "key": "val" }`)
	isEmpty = isEmptyJSON(t, jsonValue)
	assert.False(t, isEmpty, `The JSON value { "key": "val" } should be treated as a non-empty value.`)

	jsonValue = []byte(`[]`)
	isEmpty = isEmptyJSON(t, jsonValue)
	assert.True(t, isEmpty, `The JSON value "[]" should be treated as an empty value.`)

	jsonValue = []byte(`[{ "key": "val" }]`)
	isEmpty = isEmptyJSON(t, jsonValue)
	assert.False(t, isEmpty, `The JSON value [{ "key": "val" }] should be treated as a non-empty value.`)
}

func TestSaveAndLoadTerraformOptions(t *testing.T) {
	t.Parallel()

	tmpFolder := t.TempDir()

	expectedData := &terraform.Options{
		TerraformDir: "/abc/def/ghi",
		Vars:         map[string]interface{}{},
	}
	SaveTerraformOptions(t, tmpFolder, expectedData)

	actualData := LoadTerraformOptions(t, tmpFolder)
	assert.Equal(t, expectedData, actualData)
}

func TestSaveTerraformOptionsIfNotPresent(t *testing.T) {
	t.Parallel()

	tmpFolder := t.TempDir()

	expectedData := &terraform.Options{
		TerraformDir: "/abc/def/ghi",
		Vars:         map[string]interface{}{},
	}
	SaveTerraformOptionsIfNotPresent(t, tmpFolder, expectedData)

	overwritingData := &terraform.Options{
		TerraformDir: "/123/456/789",
		Vars:         map[string]interface{}{},
	}
	SaveTerraformOptionsIfNotPresent(t, tmpFolder, overwritingData)

	actualData := LoadTerraformOptions(t, tmpFolder)
	assert.Equal(t, expectedData, actualData)
}

func TestSaveTerraformOptionsOverwrite(t *testing.T) {
	t.Parallel()

	tmpFolder := t.TempDir()

	originaData := &terraform.Options{
		TerraformDir: "/abc/def/ghi",
		Vars:         map[string]interface{}{},
	}
	SaveTerraformOptions(t, tmpFolder, originaData)

	overwritingData := &terraform.Options{
		TerraformDir: "/123/456/789",
		Vars:         map[string]interface{}{},
	}
	SaveTerraformOptions(t, tmpFolder, overwritingData)

	actualData := LoadTerraformOptions(t, tmpFolder)
	assert.Equal(t, overwritingData, actualData)
}

func TestSaveAndLoadAmiId(t *testing.T) {
	t.Parallel()

	tmpFolder := t.TempDir()

	expectedData := "ami-abcd1234"
	SaveAmiId(t, tmpFolder, expectedData)

	actualData := LoadAmiId(t, tmpFolder)
	assert.Equal(t, expectedData, actualData)
}

func TestSaveAndLoadArtifactID(t *testing.T) {
	t.Parallel()

	tmpFolder := t.TempDir()

	expectedData := "terratest-packer-example-2018-08-08t15-35-19z"
	SaveArtifactID(t, tmpFolder, expectedData)

	actualData := LoadArtifactID(t, tmpFolder)
	assert.Equal(t, expectedData, actualData)
}

func TestSaveAndLoadNamedStrings(t *testing.T) {
	t.Parallel()

	tmpFolder := t.TempDir()

	name1 := "test-ami"
	expectedData1 := "ami-abcd1234"

	name2 := "test-ami2"
	expectedData2 := "ami-xyz98765"

	name3 := "test-image"
	expectedData3 := "terratest-packer-example-2018-08-08t15-35-19z"

	name4 := "test-image2"
	expectedData4 := "terratest-packer-example-2018-01-03t12-35-00z"

	SaveString(t, tmpFolder, name1, expectedData1)
	SaveString(t, tmpFolder, name2, expectedData2)
	SaveString(t, tmpFolder, name3, expectedData3)
	SaveString(t, tmpFolder, name4, expectedData4)

	actualData1 := LoadString(t, tmpFolder, name1)
	actualData2 := LoadString(t, tmpFolder, name2)
	actualData3 := LoadString(t, tmpFolder, name3)
	actualData4 := LoadString(t, tmpFolder, name4)

	assert.Equal(t, expectedData1, actualData1)
	assert.Equal(t, expectedData2, actualData2)
	assert.Equal(t, expectedData3, actualData3)
	assert.Equal(t, expectedData4, actualData4)
}

func TestSaveDuplicateTestData(t *testing.T) {
	t.Parallel()

	tmpFolder := t.TempDir()

	name := "hello-world"
	val1 := "hello world"
	val2 := "buenos dias, mundo"

	SaveString(t, tmpFolder, name, val1)
	SaveString(t, tmpFolder, name, val2)

	actualVal := LoadString(t, tmpFolder, name)

	assert.Equal(t, val2, actualVal, "Actual test data should use overwritten values")
}

func TestSaveAndLoadNamedInts(t *testing.T) {
	t.Parallel()

	tmpFolder := t.TempDir()

	name1 := "test-int1"
	expectedData1 := 23842834

	name2 := "test-int2"
	expectedData2 := 52

	SaveInt(t, tmpFolder, name1, expectedData1)
	SaveInt(t, tmpFolder, name2, expectedData2)

	actualData1 := LoadInt(t, tmpFolder, name1)
	actualData2 := LoadInt(t, tmpFolder, name2)

	assert.Equal(t, expectedData1, actualData1)
	assert.Equal(t, expectedData2, actualData2)
}

func TestSaveAndLoadKubectlOptions(t *testing.T) {
	t.Parallel()

	tmpFolder := t.TempDir()

	expectedData := &k8s.KubectlOptions{
		ContextName: "terratest-context",
		ConfigPath:  "~/.kube/config",
		Namespace:   "default",
		Env: map[string]string{
			"TERRATEST_ENV_VAR": "terratest",
		},
	}
	SaveKubectlOptions(t, tmpFolder, expectedData)

	actualData := LoadKubectlOptions(t, tmpFolder)
	assert.Equal(t, expectedData, actualData)
}
