package docker

import (
	"github.com/gruntwork-io/terratest/modules/logger/v2"
	"github.com/gruntwork-io/terratest/modules/shell/v2"
	"github.com/gruntwork-io/terratest/modules/testing/v2"
	"github.com/stretchr/testify/require"
)

// Push runs the 'docker push' command to push the given tag. This will fail the test if there are any errors.
func Push(t testing.TestingT, logger *logger.Logger, tag string) {
	require.NoError(t, PushE(t, logger, tag))
}

// PushE runs the 'docker push' command to push the given tag.
func PushE(t testing.TestingT, logger *logger.Logger, tag string) error {
	logger.Logf(t, "Running 'docker push' for tag %s", tag)

	cmd := shell.Command{
		Command: "docker",
		Args:    []string{"push", tag},
		Logger:  logger,
	}
	return shell.RunCommandE(t, cmd)
}
