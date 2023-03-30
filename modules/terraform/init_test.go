package terraform

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/gruntwork-io/terratest/modules/files"
	"github.com/gruntwork-io/terratest/modules/logger"
	ttest "github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInitBackendConfig(t *testing.T) {
	t.Parallel()

	stateDirectory := t.TempDir()

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

type testLog struct {
	w io.Writer
}

func (l testLog) Logf(t ttest.TestingT, format string, args ...interface{}) {
	fmt.Fprintf(l.w, format, args...)
}

func TestInitAdditionalFlags(t *testing.T) {
	t.Parallel()

	ttests := map[string]struct {
		// returns logged out buffer, opts, expect string, cleanup
		setup func(t *testing.T) (*bytes.Buffer, *Options, string, func())
	}{
		"backend set to false": {
			func(t *testing.T) (*bytes.Buffer, *Options, string, func()) {
				b := &bytes.Buffer{}
				l := testLog{b}
				stateDirectory := t.TempDir()
				testFolder, err := files.CopyTerraformFolderToTemp("../../test/fixtures/terraform-backend", "be-set-to-false")
				require.NoError(t, err)
				backendPath := filepath.Join(stateDirectory, "backend.tfstate")

				return b,
					&Options{
						Logger:       logger.New(l),
						TerraformDir: testFolder,
						Reconfigure:  true,
						BackendConfig: map[string]any{
							"path": backendPath,
						},
						AdditionalInitFlags: []string{"-backend=false"},
					},
					fmt.Sprintf("[init -upgrade=false -reconfigure -backend-config=path=%s -backend=false]", backendPath),
					func() {
						os.RemoveAll(testFolder)
					}
			},
		},
		"backend set to true": {
			func(t *testing.T) (*bytes.Buffer, *Options, string, func()) {
				b := &bytes.Buffer{}
				l := testLog{b}
				stateDirectory := t.TempDir()
				testFolder, err := files.CopyTerraformFolderToTemp("../../test/fixtures/terraform-backend", "be-set-to-true")
				require.NoError(t, err)
				backendPath := filepath.Join(stateDirectory, "backend.tfstate")

				return b, &Options{
						Logger:       logger.New(l),
						TerraformDir: testFolder,
						Reconfigure:  true,
						BackendConfig: map[string]any{
							"path": backendPath,
						},
						AdditionalInitFlags: []string{"-backend=true"},
					},
					fmt.Sprintf("[init -upgrade=false -reconfigure -backend-config=path=%s -backend=true]", backendPath),
					func() {
						os.RemoveAll(testFolder)
					}
			},
		},
		"backend not set via additional args": {
			func(t *testing.T) (*bytes.Buffer, *Options, string, func()) {
				b := &bytes.Buffer{}
				l := testLog{b}
				stateDirectory := t.TempDir()
				testFolder, err := files.CopyTerraformFolderToTemp("../../test/fixtures/terraform-backend", "not-set-via-args")
				require.NoError(t, err)
				backendPath := filepath.Join(stateDirectory, "backend.tfstate")

				return b, &Options{
						Logger:       logger.New(l),
						TerraformDir: testFolder,
						Reconfigure:  true,
						BackendConfig: map[string]any{
							"path": backendPath,
						},
						AdditionalInitFlags: []string{},
					},
					fmt.Sprintf("[init -upgrade=false -reconfigure -backend-config=path=%s]", backendPath),
					func() {
						os.RemoveAll(testFolder)
					}
			},
		},
		// should ignore the
		"backend set to false and protected flags re-specified": {
			setup: func(t *testing.T) (*bytes.Buffer, *Options, string, func()) {
				b := &bytes.Buffer{}
				l := testLog{b}
				stateDirectory := t.TempDir()
				testFolder, err := files.CopyTerraformFolderToTemp("../../test/fixtures/terraform-backend", "be-no-specified")
				require.NoError(t, err)
				backendPath := filepath.Join(stateDirectory, "backend.tfstate")
				return b,
					&Options{
						Logger:       logger.New(l),
						TerraformDir: testFolder,
						Reconfigure:  false,
						MigrateState: false,
						BackendConfig: map[string]any{
							"path": backendPath,
						},
						AdditionalInitFlags: []string{"-backend=false", "-reconfigure", "-migrate-state"},
					},
					fmt.Sprintf("[init -upgrade=false -backend-config=path=%s -backend=false]", backendPath),
					func() {
						os.RemoveAll(testFolder)
					}
			},
		},
	}
	for name, tt := range ttests {
		t.Run(name, func(t *testing.T) {

			b, options, expect, cleanUp := tt.setup(t)
			defer cleanUp()
			_, err := InitE(t, options)
			assert.Contains(t, b.String(), expect)
			assert.NoError(t, err)
		})
	}
}
