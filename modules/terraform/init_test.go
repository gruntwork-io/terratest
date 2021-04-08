package terraform

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/gruntwork-io/terratest/modules/files"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInitBackendConfig(t *testing.T) {
	t.Parallel()

	stateDirectory, err := ioutil.TempDir("", t.Name())
	if err != nil {
		t.Fatal(err)
	}

	remoteStateFile := filepath.Join(stateDirectory, "backend.tfstate")

	testFolder, err := files.CopyTerraformFolderToTemp("../../test/fixtures/terraform-backend", t.Name())
	if err != nil {
		t.Fatal(err)
	}

	options := &Options{
		TerraformDir: testFolder,
		BackendConfig: map[string]interface{}{
			"path": remoteStateFile,
		},
	}

	InitAndApply(t, options)

	assert.FileExists(t, remoteStateFile)
}

func TestInitPluginDir(t *testing.T) {
	t.Parallel()

	testingDir, err := ioutil.TempDir("", t.Name())
	require.NoError(t, err)
	defer os.RemoveAll(testingDir)

	terraformFixture := "../../test/fixtures/terraform-basic-configuration"

	initializedFolder, err := files.CopyTerraformFolderToTemp(terraformFixture, t.Name())
	require.NoError(t, err)
	defer os.RemoveAll(initializedFolder)

	testFolder, err := files.CopyTerraformFolderToTemp(terraformFixture, t.Name())
	require.NoError(t, err)
	defer os.RemoveAll(testFolder)

	terraformOptions := &Options{
		TerraformDir: initializedFolder,
	}

	terraformOptionsPluginDir := &Options{
		TerraformDir: testFolder,
		PluginDir:    testingDir,
	}

	Init(t, terraformOptions)

	_, err = InitE(t, terraformOptionsPluginDir)
	require.Error(t, err)

	// In Terraform 0.13, the directory is "plugins"
	initializedPluginDir := initializedFolder + "/.terraform/plugins"

	// In Terraform 0.14, the directory is "providers"
	initializedProviderDir := initializedFolder + "/.terraform/providers"

	files.CopyFolderContents(initializedPluginDir, testingDir)
	files.CopyFolderContents(initializedProviderDir, testingDir)

	initOutput := Init(t, terraformOptionsPluginDir)

	assert.Contains(t, initOutput, "(unauthenticated)")
}
