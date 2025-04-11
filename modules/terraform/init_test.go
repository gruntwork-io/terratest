package terraform

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/gruntwork-io/terratest/modules/files"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInitBackendConfig(t *testing.T) {
	t.Parallel()

	testFolderPath := "../../test/fixtures/terraform-backend"
	testFolder, err := files.CopyTerraformFolderToTemp(testFolderPath, t.Name())
	if err != nil {
		t.Fatal(err)
	}

	tmpStateFile := filepath.Join(t.TempDir(), "backend.tfstate")
	ttable := []struct {
		name    string
		path    string
		options *Options
	}{
		{
			name: "KeyValue",
			path: tmpStateFile,
			options: &Options{
				TerraformDir: testFolder,
				BackendConfig: map[string]interface{}{
					"path": tmpStateFile,
				},
			},
		},
		{
			name: "File",
			path: filepath.Join(testFolder, "backend.tfstate"),
			options: &Options{
				TerraformDir: testFolder,
				Reconfigure:  true,
				BackendConfig: map[string]interface{}{
					"backend.hcl": nil,
				},
			},
		},
	}

	for _, tt := range ttable {
		t.Run(tt.name, func(t *testing.T) {
			InitAndApply(t, tt.options)
			assert.FileExists(t, tt.path)
		})
	}
}

func TestInitPluginDir(t *testing.T) {
	t.Parallel()

	testingDir := t.TempDir()

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

func TestInitReconfigureBackend(t *testing.T) {
	t.Parallel()

	stateDirectory := t.TempDir()

	testFolder, err := files.CopyTerraformFolderToTemp("../../test/fixtures/terraform-backend", t.Name())
	require.NoError(t, err)
	defer os.RemoveAll(testFolder)

	options := &Options{
		TerraformDir: testFolder,
		BackendConfig: map[string]interface{}{
			"path":          filepath.Join(stateDirectory, "backend.tfstate"),
			"workspace_dir": "current",
		},
	}

	Init(t, options)

	options.BackendConfig["workspace_dir"] = "new"
	_, err = InitE(t, options)
	assert.Error(t, err, "Backend initialization with changed configuration should fail without -reconfigure option")

	options.Reconfigure = true
	_, err = InitE(t, options)
	assert.NoError(t, err, "Backend initialization with changed configuration should success with -reconfigure option")
}

func TestInitBackendMigration(t *testing.T) {
	t.Parallel()

	stateDirectory := t.TempDir()

	testFolder, err := files.CopyTerraformFolderToTemp("../../test/fixtures/terraform-backend", t.Name())
	require.NoError(t, err)
	defer os.RemoveAll(testFolder)

	options := &Options{
		TerraformDir: testFolder,
		BackendConfig: map[string]interface{}{
			"path":          filepath.Join(stateDirectory, "backend.tfstate"),
			"workspace_dir": "current",
		},
	}

	Init(t, options)

	options.BackendConfig["workspace_dir"] = "new"
	_, err = InitE(t, options)
	assert.Error(t, err, "Backend initialization with changed configuration should fail without -migrate-state option")

	options.MigrateState = true
	_, err = InitE(t, options)
	assert.NoError(t, err, "Backend initialization with changed configuration should success with -migrate-state option")
}

func TestInitNoColorOption(t *testing.T) {
	t.Parallel()

	testFolder, err := files.CopyTerraformFolderToTemp("../../test/fixtures/terraform-no-error", t.Name())
	require.NoError(t, err)

	options := WithDefaultRetryableErrors(t, &Options{
		TerraformDir: testFolder,
		NoColor:      true,
	})

	out := InitAndApply(t, options)

	require.Contains(t, out, "Hello, World")

	// Check that NoColor correctly doesn't output the colour escape codes which look like [0m,[1m or [32m
	require.NotRegexp(t, `\[\d*m`, out, "Output should not contain color escape codes")
}
