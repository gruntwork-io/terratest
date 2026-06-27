package docker_test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/docker/v2"
	"github.com/stretchr/testify/require"
)

func TestDockerComposeWithBuildKit(t *testing.T) {
	t.Parallel()

	testToken := "testToken"
	dockerOptions := &docker.Options{
		// Directory where docker-compose.yml lives
		WorkingDir: "../../test/fixtures/docker-compose-with-buildkit",

		// Configure the port the web app will listen on and the text it will return using environment variables
		EnvVars: map[string]string{
			"GITHUB_OAUTH_TOKEN": testToken,
		},
		EnableBuildKit: true,
	}

	ctx := t.Context()
	docker.RunDockerComposeContext(t, ctx, dockerOptions, "build", "--no-cache")

	out := docker.RunDockerComposeContext(t, ctx, dockerOptions, "up")

	require.Contains(t, out, testToken)
}

func TestDockerComposeWithCustomProjectName(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		options  *docker.Options
		expected string
	}{
		{
			name: "Testing ",
			options: &docker.Options{
				WorkingDir: "../../test/fixtures/docker-compose-with-custom-project-name",
			},
			expected: "testdockercomposewithcustomprojectname",
		},
		{
			name: "Testing",
			options: &docker.Options{
				WorkingDir:  "../../test/fixtures/docker-compose-with-custom-project-name",
				ProjectName: "testingProjectName",
			},
			expected: "testingprojectname",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			t.Log(test.name)

			ctx := t.Context()

			output := docker.RunDockerComposeContext(t, ctx, test.options, "up", "-d")
			defer docker.RunDockerComposeContext(t, ctx, test.options, "down", "--remove-orphans", "--timeout", "2")

			require.Contains(t, output, test.expected)
		})
	}
}
