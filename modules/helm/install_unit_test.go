package helm //nolint:testpackage // white-box test for the unexported installArgs helper

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// When a release name is provided, it is passed as a positional argument right before the chart.
func TestInstallArgsIncludesReleaseName(t *testing.T) {
	t.Parallel()

	args, err := installArgs(&Options{}, "stable/chart", "my-release")
	require.NoError(t, err)
	assert.Equal(t, []string{"my-release", "stable/chart"}, args)
}

// When the release name is empty, it must not be passed as a positional argument. Otherwise helm fails with
// "expected at most two arguments". This lets callers use helm's --generate-name through ExtraArgs.
func TestInstallArgsOmitsEmptyReleaseName(t *testing.T) {
	t.Parallel()

	options := &Options{ExtraArgs: map[string][]string{"install": {"--generate-name"}}}
	args, err := installArgs(options, "stable/chart", "")
	require.NoError(t, err)

	assert.Equal(t, []string{"--generate-name", "stable/chart"}, args)
	assert.NotContains(t, args, "")
}

// The full argument order must be preserved: ExtraArgs, then --version, then values args, then the release name, with
// the chart always last.
func TestInstallArgsOrdering(t *testing.T) {
	t.Parallel()

	options := &Options{
		Version:   "1.2.3",
		SetValues: map[string]string{"foo": "bar"},
		ExtraArgs: map[string][]string{"install": {"--atomic"}},
	}
	args, err := installArgs(options, "repo/chart", "my-release")
	require.NoError(t, err)

	assert.Equal(t, []string{"--atomic", "--version", "1.2.3", "--set", "foo=bar", "my-release", "repo/chart"}, args)
	// The chart must be the final positional argument, with the release name immediately before it.
	assert.Equal(t, "repo/chart", args[len(args)-1])
	assert.Equal(t, "my-release", args[len(args)-2])
}
