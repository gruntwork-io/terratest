package test_structure_test //nolint:staticcheck // package name determined by directory

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/aws"
	"github.com/gruntwork-io/terratest/modules/files"
	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/ssh"
	"github.com/gruntwork-io/terratest/modules/terraform"
	teststructure "github.com/gruntwork-io/terratest/modules/test-structure"
	gotesting "github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testData struct {
	Baz map[string]interface{}
	Foo string
	Bar bool
}

func TestSaveAndLoadTestData(t *testing.T) {
	t.Parallel()

	isTestDataPresent := teststructure.IsTestDataPresent(t, "/file/that/does/not/exist")
	assert.False(t, isTestDataPresent, "Expected no test data would be present because no test data file exists.")

	tmpDir := t.TempDir()
	tmpFile := teststructure.FormatTestDataPath(tmpDir, "test-save-load.json")

	expectedData := testData{
		Foo: "foo",
		Bar: true,
		Baz: map[string]interface{}{"abc": "def", "ghi": 1.0, "klm": false},
	}

	isTestDataPresent = teststructure.IsTestDataPresent(t, tmpFile)
	assert.False(t, isTestDataPresent, "Expected no test data would be present because file exists but no data has been written yet.")

	overwrite := true
	teststructure.SaveTestData(t, tmpFile, overwrite, expectedData)

	isTestDataPresent = teststructure.IsTestDataPresent(t, tmpFile)
	assert.True(t, isTestDataPresent, "Expected test data would be present because file exists and data has been written to file.")

	actualData := testData{}
	teststructure.LoadTestData(t, tmpFile, &actualData)
	assert.Equal(t, expectedData, actualData)

	overwritingData := testData{
		Foo: "foo",
		Bar: false,
		Baz: map[string]interface{}{"123": "456", "789": 1.0, "0": false},
	}
	teststructure.SaveTestData(t, tmpFile, !overwrite, overwritingData)
	teststructure.LoadTestData(t, tmpFile, &actualData)
	assert.Equal(t, expectedData, actualData)

	teststructure.CleanupTestData(t, tmpFile)
	assert.False(t, files.FileExists(tmpFile))
}

func TestIsEmptyJson(t *testing.T) {
	t.Parallel()

	var jsonValue []byte

	var isEmpty bool

	jsonValue = []byte("null")
	isEmpty = teststructure.IsEmptyJSON(t, jsonValue)
	assert.True(t, isEmpty, `The JSON literal "null" should be treated as an empty value.`)

	jsonValue = []byte("false")
	isEmpty = teststructure.IsEmptyJSON(t, jsonValue)
	assert.True(t, isEmpty, `The JSON literal "false" should be treated as an empty value.`)

	jsonValue = []byte("true")
	isEmpty = teststructure.IsEmptyJSON(t, jsonValue)
	assert.False(t, isEmpty, `The JSON literal "true" should be treated as a non-empty value.`)

	jsonValue = []byte("0")
	isEmpty = teststructure.IsEmptyJSON(t, jsonValue)
	assert.True(t, isEmpty, `The JSON literal "0" should be treated as an empty value.`)

	jsonValue = []byte("1")
	isEmpty = teststructure.IsEmptyJSON(t, jsonValue)
	assert.False(t, isEmpty, `The JSON literal "1" should be treated as a non-empty value.`)

	jsonValue = []byte("{}")
	isEmpty = teststructure.IsEmptyJSON(t, jsonValue)
	assert.True(t, isEmpty, `The JSON value "{}" should be treated as an empty value.`)

	jsonValue = []byte(`{ "key": "val" }`)
	isEmpty = teststructure.IsEmptyJSON(t, jsonValue)
	assert.False(t, isEmpty, `The JSON value { "key": "val" } should be treated as a non-empty value.`)

	jsonValue = []byte(`[]`)
	isEmpty = teststructure.IsEmptyJSON(t, jsonValue)
	assert.True(t, isEmpty, `The JSON value "[]" should be treated as an empty value.`)

	jsonValue = []byte(`[{ "key": "val" }]`)
	isEmpty = teststructure.IsEmptyJSON(t, jsonValue)
	assert.False(t, isEmpty, `The JSON value [{ "key": "val" }] should be treated as a non-empty value.`)
}

func TestSaveAndLoadTerraformOptions(t *testing.T) {
	t.Parallel()

	tmpFolder := t.TempDir()

	expectedData := &terraform.Options{
		TerraformDir: "/abc/def/ghi",
		Vars:         map[string]interface{}{},
	}
	teststructure.SaveTerraformOptions(t, tmpFolder, expectedData)

	actualData := teststructure.LoadTerraformOptions(t, tmpFolder)
	assert.Equal(t, expectedData, actualData)
}

func TestSaveTerraformOptionsIfNotPresent(t *testing.T) {
	t.Parallel()

	tmpFolder := t.TempDir()

	expectedData := &terraform.Options{
		TerraformDir: "/abc/def/ghi",
		Vars:         map[string]interface{}{},
	}
	teststructure.SaveTerraformOptionsIfNotPresent(t, tmpFolder, expectedData)

	overwritingData := &terraform.Options{
		TerraformDir: "/123/456/789",
		Vars:         map[string]interface{}{},
	}
	teststructure.SaveTerraformOptionsIfNotPresent(t, tmpFolder, overwritingData)

	actualData := teststructure.LoadTerraformOptions(t, tmpFolder)
	assert.Equal(t, expectedData, actualData)
}

func TestSaveTerraformOptionsOverwrite(t *testing.T) {
	t.Parallel()

	tmpFolder := t.TempDir()

	originaData := &terraform.Options{
		TerraformDir: "/abc/def/ghi",
		Vars:         map[string]interface{}{},
	}
	teststructure.SaveTerraformOptions(t, tmpFolder, originaData)

	overwritingData := &terraform.Options{
		TerraformDir: "/123/456/789",
		Vars:         map[string]interface{}{},
	}
	teststructure.SaveTerraformOptions(t, tmpFolder, overwritingData)

	actualData := teststructure.LoadTerraformOptions(t, tmpFolder)
	assert.Equal(t, overwritingData, actualData)
}

func TestSaveAndLoadAmiId(t *testing.T) {
	t.Parallel()

	tmpFolder := t.TempDir()

	expectedData := "ami-abcd1234"
	teststructure.SaveArtifactID(t, tmpFolder, expectedData)

	actualData := teststructure.LoadArtifactID(t, tmpFolder)
	assert.Equal(t, expectedData, actualData)
}

func TestSaveAndLoadArtifactID(t *testing.T) {
	t.Parallel()

	tmpFolder := t.TempDir()

	expectedData := "terratest-packer-example-2018-08-08t15-35-19z"
	teststructure.SaveArtifactID(t, tmpFolder, expectedData)

	actualData := teststructure.LoadArtifactID(t, tmpFolder)
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

	teststructure.SaveString(t, tmpFolder, name1, expectedData1)
	teststructure.SaveString(t, tmpFolder, name2, expectedData2)
	teststructure.SaveString(t, tmpFolder, name3, expectedData3)
	teststructure.SaveString(t, tmpFolder, name4, expectedData4)

	actualData1 := teststructure.LoadString(t, tmpFolder, name1)
	actualData2 := teststructure.LoadString(t, tmpFolder, name2)
	actualData3 := teststructure.LoadString(t, tmpFolder, name3)
	actualData4 := teststructure.LoadString(t, tmpFolder, name4)

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

	teststructure.SaveString(t, tmpFolder, name, val1)
	teststructure.SaveString(t, tmpFolder, name, val2)

	actualVal := teststructure.LoadString(t, tmpFolder, name)

	assert.Equal(t, val2, actualVal, "Actual test data should use overwritten values")
}

func TestSaveAndLoadNamedInts(t *testing.T) {
	t.Parallel()

	tmpFolder := t.TempDir()

	name1 := "test-int1"
	expectedData1 := 23842834

	name2 := "test-int2"
	expectedData2 := 52

	teststructure.SaveInt(t, tmpFolder, name1, expectedData1)
	teststructure.SaveInt(t, tmpFolder, name2, expectedData2)

	actualData1 := teststructure.LoadInt(t, tmpFolder, name1)
	actualData2 := teststructure.LoadInt(t, tmpFolder, name2)

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
	teststructure.SaveKubectlOptions(t, tmpFolder, expectedData)

	actualData := teststructure.LoadKubectlOptions(t, tmpFolder)
	assert.Equal(t, expectedData, actualData)
}

type tStringLogger struct {
	sb strings.Builder
}

func (l *tStringLogger) Logf(t gotesting.TestingT, format string, args ...interface{}) {
	t.Helper()
	fmt.Fprintf(&l.sb, format, args...)
	l.sb.WriteRune('\n')
}

func TestSaveAndLoadEC2KeyPair(t *testing.T) {
	t.Parallel()

	def, slogger := logger.Default, &tStringLogger{}
	logger.Default = logger.New(slogger)

	t.Cleanup(func() {
		logger.Default = def
	})

	keyPair, err := ssh.GenerateRSAKeyPairE(t, 2048) //nolint:mnd // RSA key size for testing
	require.NoError(t, err)

	ec2KeyPair := &aws.Ec2Keypair{
		KeyPair: keyPair,
		Name:    "test-ec2-key-pair",
		Region:  "us-east-1",
	}

	storedEC2KeyPair, err := json.Marshal(ec2KeyPair) //nolint:musttag // aws.Ec2Keypair does not have json tags
	require.NoError(t, err)

	tmpFolder := t.TempDir()
	teststructure.SaveEc2KeyPair(t, tmpFolder, ec2KeyPair)
	loadedEC2KeyPair := teststructure.LoadEc2KeyPair(t, tmpFolder)
	assert.Equal(t, ec2KeyPair, loadedEC2KeyPair)

	assert.NotContains(t, slogger.sb.String(), string(storedEC2KeyPair), "stored ec2 key pair should not be logged")
}
