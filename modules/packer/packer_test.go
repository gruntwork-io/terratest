package packer_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/packer/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExtractAmiIdFromOneLine(t *testing.T) {
	t.Parallel()

	expectedAMIID := "ami-b481b3de"
	text := "1456332887,amazon-ebs,artifact,0,id,us-east-1:" + expectedAMIID

	actualAMIID, err := packer.ExtractArtifactID(text)
	if err != nil {
		t.Errorf("Did not expect to get an error when extracting a valid AMI ID: %s", err)
	}

	if actualAMIID != expectedAMIID {
		t.Errorf("Did not get expected AMI ID. Expected: %s. Actual: %s.", expectedAMIID, actualAMIID)
	}
}

func TestExtractImageIdFromOneLine(t *testing.T) {
	t.Parallel()

	expectedImageID := "terratest-packer-example-2018-08-09t12-02-58z"
	text := "1533816302,googlecompute,artifact,0,id," + expectedImageID

	actualImageID, err := packer.ExtractArtifactID(text)
	if err != nil {
		t.Errorf("Did not expect to get an error when extracting a valid Image ID: %s", err)
	}

	if actualImageID != expectedImageID {
		t.Errorf("Did not get expected Image ID. Expected: %s. Actual: %s.", expectedImageID, actualImageID)
	}
}

func TestExtractAmiIdFromMultipleLines(t *testing.T) {
	t.Parallel()

	expectedAMIID := "ami-b481b3de"
	text := `
	foo
	bar
	1456332887,amazon-ebs,artifact,0,id,us-east-1:` + expectedAMIID + `
	baz
	blah
	`

	actualAMIID, err := packer.ExtractArtifactID(text)
	if err != nil {
		t.Errorf("Did not expect to get an error when extracting a valid AMI ID: %s", err)
	}

	if actualAMIID != expectedAMIID {
		t.Errorf("Did not get expected AMI ID. Expected: %s. Actual: %s.", expectedAMIID, actualAMIID)
	}
}

func TestExtractImageIdFromMultipleLines(t *testing.T) {
	t.Parallel()

	expectedImageID := "terratest-packer-example-2018-08-09t12-02-58z"
	text := `
	foo
	bar
	1533816302,googlecompute,artifact,0,id,` + expectedImageID + `
	baz
	blah
	`

	actualImageID, err := packer.ExtractArtifactID(text)
	if err != nil {
		t.Errorf("Did not expect to get an error when extracting a valid Image ID: %s", err)
	}

	if actualImageID != expectedImageID {
		t.Errorf("Did not get the expected Image ID. Expected: %s. Actual: %s.", expectedImageID, actualImageID)
	}
}

func TestExtractAmiIdNoIdPresent(t *testing.T) {
	t.Parallel()

	text := `
	foo
	bar
	baz
	blah
	`

	_, err := packer.ExtractArtifactID(text)
	if err == nil {
		t.Error("Expected to get an error when extracting an AMI ID from text with no AMI in it, but got nil")
	}
}

func TestExtractArtifactINoIdPresent(t *testing.T) {
	t.Parallel()

	text := `
	foo
	bar
	baz
	blah
	`

	_, err := packer.ExtractArtifactID(text)
	if err == nil {
		t.Error("Expected to get an error when extracting an Artifact ID from text with no Artifact ID in it, but got nil")
	}
}

func TestFormatPackerArgs(t *testing.T) {
	t.Parallel()

	tests := []struct {
		option   *packer.Options
		expected string
	}{
		{
			option: &packer.Options{
				Template: "packer.json",
			},
			expected: "build -machine-readable packer.json",
		},
		{
			option: &packer.Options{
				Template: "packer.json",
				Vars: map[string]string{
					"foo": "bar",
				},
				Only: "onlythis",
			},
			expected: "build -machine-readable -var foo=bar -only=onlythis packer.json",
		},
		{
			option: &packer.Options{
				Template: "packer.json",
				Vars: map[string]string{
					"foo": "bar",
				},
				Only:   "onlythis",
				Except: "long-run-pp,artifact",
			},
			expected: "build -machine-readable -var foo=bar -only=onlythis -except=long-run-pp,artifact packer.json",
		},
		{
			option: &packer.Options{
				Template: "packer.json",
				Vars: map[string]string{
					"foo": "bar",
				},
				VarFiles: []string{
					"foofile.json",
				},
			},
			expected: "build -machine-readable -var foo=bar -var-file foofile.json packer.json",
		},
	}

	for _, test := range tests {
		args := packer.FormatPackerArgs(test.option)
		assert.Equal(t, test.expected, strings.Join(args, " "))
	}
}

func TestTrimPackerVersion(t *testing.T) {
	t.Parallel()

	tests := []struct {
		versionOutput string
		expected      string
	}{
		{
			// Pre 1.10 output
			versionOutput: "1.7.0",
			expected:      "1.7.0",
		},
		{
			// From 1.10 matches the output of packer version
			versionOutput: "Packer v1.10.0",
			expected:      "1.10.0",
		},
		{
			// From 1.10 matches the output of packer version
			versionOutput: "Packer v1.10.0\n\nYour version of Packer is out of date! The latest version\nis 1.10.3. You can update by downloading from www.packer.io/downloads\n",
			expected:      "1.10.0",
		},
	}

	for _, test := range tests {
		t.Run(test.versionOutput, func(t *testing.T) {
			t.Parallel()

			out := packer.TrimPackerVersion(test.versionOutput)
			assert.Equal(t, test.expected, out)
		})
	}
}

func TestGetArtifactIDFromManifestBuildNameE(t *testing.T) {
	t.Parallel()

	// example manifest from https://developer.hashicorp.com/packer/docs/post-processors/manifest
	manifest := `
{
  "builds": [
    {
      "name": "docker",
      "builder_type": "docker",
      "build_time": 1507245986,
      "files": [
        {
          "name": "packer_example",
          "size": 102219776
        }
      ],
      "artifact_id": "Container",
      "packer_run_uuid": "6d5d3185-fa95-44e1-8775-9e64fe2e2d8f",
      "custom_data": {
        "my_custom_data": "example"
      }
    }
  ],
  "last_run_uuid": "6d5d3185-fa95-44e1-8775-9e64fe2e2d8f"
}
`
	manifestPath := filepath.Join(t.TempDir(), "manifest.json")
	err := os.WriteFile(manifestPath, []byte(manifest), 0600)
	require.NoError(t, err)

	t.Run("Found", func(t *testing.T) {
		t.Parallel()

		artifactID, err := packer.GetArtifactIDFromManifestBuildNameE(t, manifestPath, "docker")
		require.NoError(t, err)
		assert.Equal(t, "Container", artifactID)

		artifactID2 := packer.GetArtifactIDFromManifestBuildName(t, manifestPath, "docker")
		assert.Equal(t, "Container", artifactID2)
	})

	t.Run("Not Found", func(t *testing.T) {
		t.Parallel()

		_, err := packer.GetArtifactIDFromManifestBuildNameE(t, manifestPath, "notfound")
		require.Error(t, err)
	})
}
