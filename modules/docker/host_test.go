package docker_test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/docker/v2"
	"github.com/stretchr/testify/assert"
)

func TestGetDockerHostFromEnv(t *testing.T) {
	t.Parallel()

	tests := []struct {
		Input    string
		Expected string
	}{
		{
			Input:    "unix:///var/run/docker.sock",
			Expected: "localhost",
		},
		{
			Input:    "npipe:////./pipe/docker_engine",
			Expected: "localhost",
		},
		{
			Input:    "tcp://1.2.3.4:1234",
			Expected: "1.2.3.4",
		},
		{
			Input:    "tcp://1.2.3.4",
			Expected: "1.2.3.4",
		},
		{
			Input:    "ssh://1.2.3.4:22",
			Expected: "1.2.3.4",
		},
		{
			Input:    "fd://1.2.3.4:1234",
			Expected: "1.2.3.4",
		},
		{
			Input:    "",
			Expected: "localhost",
		},
		{
			Input:    "invalidValue",
			Expected: "localhost",
		},
		{
			Input:    "invalid::value::with::semicolons",
			Expected: "localhost",
		},
	}
	for _, test := range tests {
		t.Run("GetDockerHostFromEnv: "+test.Input, func(t *testing.T) {
			t.Parallel()

			testEnv := []string{
				"FOO=bar",
				"DOCKER_HOST=" + test.Input,
				"BAR=baz",
			}

			host := docker.GetDockerHostFromEnv(testEnv)
			assert.Equal(t, test.Expected, host)
		})
	}

	t.Run("GetDockerHostFromEnv: DOCKER_HOST unset", func(t *testing.T) {
		t.Parallel()

		testEnv := []string{
			"FOO=bar",
			"BAR=baz",
		}

		host := docker.GetDockerHostFromEnv(testEnv)
		assert.Equal(t, "localhost", host)
	})
}
