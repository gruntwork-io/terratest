package terraform_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/gruntwork-io/terratest/modules/files"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInitBackendConfig(t *testing.T) {
	t.Parallel()

	testFolderPath := "../../test/fixtures/terraform-backend"

	ttable := []struct {
		setup func(t *testing.T, testFolder string) (*terraform.Options, string)
		name  string
	}{
		{
			name: "KeyValue",
			setup: func(t *testing.T, testFolder string) (*terraform.Options, string) {
				t.Helper()
				tmpStateFile := filepath.Join(t.TempDir(), "backend.tfstate")

				return &terraform.Options{
					TerraformDir: testFolder,
					BackendConfig: map[string]any{
						"path": tmpStateFile,
					},
				}, tmpStateFile
			},
		},
		{
			name: "File",
			setup: func(t *testing.T, testFolder string) (*terraform.Options, string) {
				t.Helper()

				return &terraform.Options{
					TerraformDir: testFolder,
					Reconfigure:  true,
					BackendConfig: map[string]any{
						"backend.hcl": nil,
					},
				}, filepath.Join(testFolder, "backend.tfstate")
			},
		},
	}

	for _, tt := range ttable {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			testFolder, err := files.CopyTerraformFolderToTemp(testFolderPath, tt.name)
			require.NoError(t, err)

			options, expectedPath := tt.setup(t, testFolder)
			terraform.InitAndApply(t, options)
			assert.FileExists(t, expectedPath)
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

	terraformOptions := &terraform.Options{
		TerraformDir: initializedFolder,
	}

	terraformOptionsPluginDir := &terraform.Options{
		TerraformDir: testFolder,
		PluginDir:    testingDir,
	}

	terraform.Init(t, terraformOptions)

	_, err = terraform.InitE(t, terraformOptionsPluginDir)
	require.Error(t, err)

	// In Terraform 0.13, the directory is "plugins"
	initializedPluginDir := initializedFolder + "/.terraform/plugins"

	// In Terraform 0.14, the directory is "providers"
	initializedProviderDir := initializedFolder + "/.terraform/providers"

	files.CopyFolderContents(initializedPluginDir, testingDir)
	files.CopyFolderContents(initializedProviderDir, testingDir)

	initOutput := terraform.Init(t, terraformOptionsPluginDir)

	assert.Contains(t, initOutput, "(unauthenticated)")
}

func TestInitReconfigureBackend(t *testing.T) {
	t.Parallel()

	stateDirectory := t.TempDir()

	testFolder, err := files.CopyTerraformFolderToTemp("../../test/fixtures/terraform-backend", t.Name())
	require.NoError(t, err)

	defer os.RemoveAll(testFolder)

	options := &terraform.Options{
		TerraformDir: testFolder,
		BackendConfig: map[string]any{
			"path":          filepath.Join(stateDirectory, "backend.tfstate"),
			"workspace_dir": "current",
		},
	}

	terraform.Init(t, options)

	options.BackendConfig["workspace_dir"] = "new"
	_, err = terraform.InitE(t, options)
	require.Error(t, err, "Backend initialization with changed configuration should fail without -reconfigure option")

	options.Reconfigure = true
	_, err = terraform.InitE(t, options)
	require.NoError(t, err, "Backend initialization with changed configuration should success with -reconfigure option")
}

func TestInitBackendMigration(t *testing.T) {
	t.Parallel()

	stateDirectory := t.TempDir()

	testFolder, err := files.CopyTerraformFolderToTemp("../../test/fixtures/terraform-backend", t.Name())
	require.NoError(t, err)

	defer os.RemoveAll(testFolder)

	options := &terraform.Options{
		TerraformDir: testFolder,
		BackendConfig: map[string]any{
			"path":          filepath.Join(stateDirectory, "backend.tfstate"),
			"workspace_dir": "current",
		},
	}

	terraform.Init(t, options)

	options.BackendConfig["workspace_dir"] = "new"
	_, err = terraform.InitE(t, options)
	require.Error(t, err, "Backend initialization with changed configuration should fail without -migrate-state option")

	options.MigrateState = true
	_, err = terraform.InitE(t, options)
	require.NoError(t, err, "Backend initialization with changed configuration should success with -migrate-state option")
}

func TestInitNoColorOption(t *testing.T) {
	t.Parallel()

	testFolder, err := files.CopyTerraformFolderToTemp("../../test/fixtures/terraform-no-error", t.Name())
	require.NoError(t, err)

	options := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: testFolder,
		NoColor:      true,
	})

	out := terraform.InitAndApply(t, options)

	require.Contains(t, out, "Hello, World")

	// Check that NoColor correctly doesn't output the colour escape codes which look like [0m,[1m or [32m
	require.NotRegexp(t, `\[\d*m`, out, "Output should not contain color escape codes")
}
