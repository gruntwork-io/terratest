package terragrunt

import (
	"encoding/json"
	"path"
	"runtime"
	"testing"

	"github.com/gruntwork-io/terratest/modules/files"
	"github.com/stretchr/testify/require"
)

func TestTerragruntStackOutput(t *testing.T) {
	t.Parallel()

	testFolder, err := files.CopyTerraformFolderToTemp("../../test/fixtures/terragrunt/terragrunt-stack-init", t.Name())
	require.NoError(t, err)

	// First initialize the stack
	_, err = TgStackInitE(t, &Options{
		TerragruntDir:    path.Join(testFolder, "live"),
		TerragruntBinary: "terragrunt",
	})
	require.NoError(t, err)

	// Generate with no-color option
	out, err := TgStackGenerateE(t, &Options{
		TerragruntDir:    path.Join(testFolder, "live"),
		TerragruntBinary: "terragrunt",
		NoColor:          true,
	})
	require.NoError(t, err)

	// Get the output of stack output command
	out, err = TgStackOutputE(t, &Options{
		TerragruntDir:    path.Join(testFolder, "live"),
		TerragruntBinary: "terragrunt",
	})
	require.NoError(t, err)
	require.Contains(t, out, ".terragrunt-stack")
	require.Contains(t, out, "has been successfully initialized!")

}

func TestTerragruntStackOutputError(t *testing.T) {
	t.Parallel()

	// Test with invalid terragrunt binary to ensure errors are caught
	_, err := TgStackOutputE(t, &Options{
		TerragruntDir:    "/nonexistent/path",
		TerragruntBinary: "nonexistent-binary",
	})
	require.Error(t, err)
	if runtime.GOOS == "linux" {
		require.Contains(t, err.Error(), "executable file not found in $PATH")
	}
}

func TestTerragruntStackOutputWithExtraArgs(t *testing.T) {
	t.Parallel()

	testFolder, err := files.CopyTerraformFolderToTemp("../../test/fixtures/terragrunt/terragrunt-stack-init", t.Name())
	require.NoError(t, err)

	baseOptions := &Options{
		TerragruntDir:    path.Join(testFolder, "live"),
		TerragruntBinary: "terragrunt",
		NoColor:          true,
	}

	_, err = TgStackGenerateE(t, baseOptions)
	require.NoError(t, err)

	_, err = TgStackRunE(t, &Options{
		TerragruntDir:    baseOptions.TerragruntDir,
		TerragruntBinary: baseOptions.TerragruntBinary,
		NoColor:          true,
		ExtraArgs: ExtraArgs{
			Apply: []string{"apply"},
		},
	})
	require.NoError(t, err)

	tests := []struct {
		name       string
		opts       *Options
		assertions func(t *testing.T, out string)
		expectErr  bool
	}{

		// Test : raw terragrunt output
		{
			name: "should apply stack without error",
			opts: &Options{
				TerragruntDir:    baseOptions.TerragruntDir,
				TerragruntBinary: baseOptions.TerragruntBinary,
			},
			assertions: func(t *testing.T, out string) {
				require.Contains(t, out, `output = "./test.txt"`)
			},
		},

		// Test: json output
		{
			name: "should return output in json format",
			opts: &Options{
				TerragruntDir:    baseOptions.TerragruntDir,
				TerragruntBinary: baseOptions.TerragruntBinary,
				OutputFormat:     "json",
			},
			assertions: func(t *testing.T, out string) {

				var parsed map[string]struct {
					Output string `json:"output"`
				}
				err := json.Unmarshal([]byte(out), &parsed)
				require.NoError(t, err, "Output should be valid JSON")

				require.Contains(t, parsed, "father", "JSON output should contain 'father' key")
				require.Equal(t, "./test.txt", parsed["father"].Output)
			},
		},

		// Test: skip generation
		{
			name: "should allow no stack generation",
			opts: &Options{
				TerragruntDir:    baseOptions.TerragruntDir,
				TerragruntBinary: baseOptions.TerragruntBinary,
				NoStackGenerate:  true,
			},
			assertions: func(t *testing.T, out string) {
				require.NotEmpty(t, out)
			},
		},

		// Test: raw output for specific key
		{
			name: "should return raw output for given key",
			opts: &Options{
				TerragruntDir:    baseOptions.TerragruntDir,
				TerragruntBinary: baseOptions.TerragruntBinary,
				OutputFormat:     "raw",
				OutputKey:        "father",
			},
			assertions: func(t *testing.T, out string) {
				require.Equal(t, "./test.txt", out)
			},
		},

		// Test: output for one unit
		{
			name: "should return output containing key and value",
			opts: &Options{
				TerragruntDir:    baseOptions.TerragruntDir,
				TerragruntBinary: baseOptions.TerragruntBinary,
				OutputKey:        "father",
			},
			assertions: func(t *testing.T, out string) {
				require.Contains(t, out, "father")
				require.Contains(t, out, "./test.txt")
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			out, err := TgStackOutputE(t, tc.opts)
			if tc.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				tc.assertions(t, out)
			}
		})
	}
}
