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
