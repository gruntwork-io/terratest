package docker

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDockerComposeWithBuildKit(t *testing.T) {
	t.Parallel()

	testToken := "testToken"
	dockerOptions := &Options{
		// Directory where docker-compose.yml lives
		WorkingDir: "../../test/fixtures/docker-compose-with-buildkit",

		// Configure the port the web app will listen on and the text it will return using environment variables
		EnvVars: map[string]string{
			"GITHUB_OAUTH_TOKEN": testToken,
		},
		EnableBuildKit: true,
	}
	out := RunDockerCompose(t, dockerOptions, "build", "--no-cache")
	out = RunDockerCompose(t, dockerOptions, "up")

	require.Contains(t, out, testToken)
}

func TestDockerComposeWithCustomProjectName(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		options  *Options
		expected string
	}{
		{
			name: "Testing ",
			options: &Options{
				WorkingDir: "../../test/fixtures/docker-compose-with-custom-project-name",
			},
			expected: "testdockercomposewithcustomprojectname",
		},
		{
			name: "Testing",
			options: &Options{
				WorkingDir:  "../../test/fixtures/docker-compose-with-custom-project-name",
				ProjectName: "testingProjectName",
			},
			expected: "testingprojectname",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Log(test.name)

			output := RunDockerCompose(t, test.options, "up", "-d")
			defer RunDockerCompose(t, test.options, "down", "--remove-orphans", "--timeout", "2")

			require.Contains(t, output, test.expected)
		})
	}
}
