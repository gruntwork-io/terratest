package docker

import (
	"context"

	"github.com/gruntwork-io/terratest/modules/core/v2/logger"
	"github.com/gruntwork-io/terratest/modules/core/v2/shell"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
)

// PushContext runs the 'docker push' command to push the given tag. This will fail the test if there are any
// errors. The ctx parameter supports cancellation and timeouts.
func PushContext(t testing.TestingT, ctx context.Context, logger *logger.Logger, tag string) {
	require.NoError(t, PushContextE(t, ctx, logger, tag))
}

// PushContextE runs the 'docker push' command to push the given tag. The ctx parameter supports cancellation
// and timeouts.
func PushContextE(t testing.TestingT, ctx context.Context, logger *logger.Logger, tag string) error {
	logger.Logf(t, "Running 'docker push' for tag %s", tag)

	cmd := &shell.Command{
		Command: "docker",
		Args:    []string{"push", tag},
		Logger:  logger,
	}

	return shell.RunCommandContextE(t, ctx, cmd)
}
