package packer

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/files"

	"github.com/stretchr/testify/assert"
)

func TestExtractAmiIdFromOneLine(t *testing.T) {
	t.Parallel()

	expectedAMIID := "ami-b481b3de"
	text := fmt.Sprintf("1456332887,amazon-ebs,artifact,0,id,us-east-1:%s", expectedAMIID)
	actualAMIID, err := extractArtifactID(text)

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
	text := fmt.Sprintf("1533816302,googlecompute,artifact,0,id,%s", expectedImageID)
	actualImageID, err := extractArtifactID(text)

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
	text := fmt.Sprintf(`
	foo
	bar
	1456332887,amazon-ebs,artifact,0,id,us-east-1:%s
	baz
	blah
	`, expectedAMIID)

	actualAMIID, err := extractArtifactID(text)

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
	text := fmt.Sprintf(`
	foo
	bar
	1533816302,googlecompute,artifact,0,id,%s
	baz
	blah
	`, expectedImageID)

	actualImageID, err := extractArtifactID(text)

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

	_, err := extractArtifactID(text)

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

	_, err := extractArtifactID(text)

	if err == nil {
		t.Error("Expected to get an error when extracting an Artifact ID from text with no Artifact ID in it, but got nil")
	}
}

func TestFormatPackerArgs(t *testing.T) {
	t.Parallel()

	tests := []struct {
		option   *Options
		expected string
	}{
		{
			option: &Options{
				Template: "packer.json",
			},
			expected: "build -machine-readable packer.json",
		},
		{
			option: &Options{
				Template: "packer.json",
				Vars: map[string]string{
					"foo": "bar",
				},
				Only: "onlythis",
			},
			expected: "build -machine-readable -var foo=bar -only=onlythis packer.json",
		},
		{
			option: &Options{
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
			option: &Options{
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
		args := formatPackerArgs(test.option)
		assert.Equal(t, strings.Join(args, " "), test.expected)
	}
}

func TestPackertInitWithFile(t *testing.T) {
	t.Parallel()
	tests := []struct {
		packerTestFile string
	}{
		{
			packerTestFile: "packer.pkr.hcl",
		}, {
			packerTestFile: "packer.json",
		},
	}
	for _, test := range tests {
		packerDirectory := t.TempDir()
		packerTestFile := filepath.Join("../../test/fixtures/packer-basic", test.packerTestFile)
		packerFile := filepath.Join(packerDirectory, test.packerTestFile)
		err := files.CopyFile(packerTestFile, packerFile)
		if err != nil {
			t.Fatal(err)
		}

		options := &Options{
			Template:   packerFile,
			WorkingDir: packerDirectory,
		}
		err = packerInit(t, options)
		assert.NoError(t, err)
	}
}

func TestPackertInitWithDirectory(t *testing.T) {
	t.Parallel()
	tests := []struct {
		packerTestDirectory string
	}{
		{
			packerTestDirectory: "packer-basic",
		},
	}
	for _, test := range tests {
		packerTestDirectory := filepath.Join("../../test/fixtures", test.packerTestDirectory)
		packerDirectory, err := files.CopyFolderToTemp(packerTestDirectory, t.Name(), func(path string) bool { return true })
		if err != nil {
			t.Fatal(err)
		}

		options := &Options{
			Template:   ".",
			WorkingDir: packerDirectory,
		}

		err = packerInit(t, options)
		assert.NoError(t, err)
	}
}
