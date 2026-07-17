package helm

import (
	"context"

	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
)

// RollbackContext will downgrade the release to the specified version. This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func RollbackContext(t testing.TestingT, ctx context.Context, options *Options, releaseName string, revision string) {
	require.NoError(t, RollbackContextE(t, ctx, options, releaseName, revision))
}

// RollbackContextE will downgrade the release to the specified version. The ctx parameter supports cancellation and
// timeouts.
func RollbackContextE(t testing.TestingT, ctx context.Context, options *Options, releaseName string, revision string) error {
	args := []string{}

	if options.ExtraArgs != nil {
		if rollbackArgs, ok := options.ExtraArgs["rollback"]; ok {
			args = append(args, rollbackArgs...)
		}
	}

	args = append(args, releaseName)

	if revision != "" {
		args = append(args, revision)
	}

	_, err := RunHelmCommandAndGetOutputContextE(t, ctx, options, "rollback", args...)

	return err
}
