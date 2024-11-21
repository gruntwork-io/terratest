package docker

import (
	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/shell"
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

type LoginOptions struct {
	Registry string
	Login    string
	Password string
}

// Login runs the 'docker login' command to login the given registry. This will fail the test if there are any errors.
// Do not use production password, only for testing purpose.
func Login(t testing.TestingT, opt LoginOptions) {
	require.NoError(t, LoginE(t, opt))
}

// login runs the 'docker login' command to login the given registry.
func LoginE(t testing.TestingT, opt LoginOptions) error {
	logger.Logf(t, "Running 'docker login' for user %s", opt.Login)

	cmd := shell.Command{
		Command: "docker",
		Args:    []string{"login", opt.Registry, "-u", opt.Login, "-p", opt.Password},
	}
	return shell.RunCommandE(t, cmd)
}
