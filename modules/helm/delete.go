package helm

import (
	"context"

	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
)

// DeleteContext will delete the provided release. If you set purge to true, the release object will be deleted as well
// so that the release name can be reused. This will fail the test if there is an error. The ctx parameter supports
// cancellation and timeouts.
func DeleteContext(t testing.TestingT, ctx context.Context, options *Options, releaseName string, purge bool) {
	require.NoError(t, DeleteContextE(t, ctx, options, releaseName, purge))
}

// DeleteContextE will delete the provided release. If you set purge to true, the release object will be deleted as
// well so that the release name can be reused. The ctx parameter supports cancellation and timeouts.
func DeleteContextE(t testing.TestingT, ctx context.Context, options *Options, releaseName string, purge bool) error {
	args := []string{}
	if !purge {
		args = append(args, "--keep-history")
	}

	if options.ExtraArgs != nil {
		if deleteArgs, ok := options.ExtraArgs["delete"]; ok {
			args = append(args, deleteArgs...)
		}
	}

	args = append(args, releaseName)
	_, err := RunHelmCommandAndGetOutputContextE(t, ctx, options, "delete", args...)

	return err
}
