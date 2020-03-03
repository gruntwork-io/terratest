package terraform

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type Format int

const (
	JSON Format = iota
	Human
)

// Show runs terraform show with the given options and returns stdout/stderr from the show command.
// This will fail the test if there is an error in the command.
func Show(t *testing.T, options *Options, format Format, path string) string {
	out, err := ShowE(t, options, format, path)
	require.NoError(t, err)
	return out
}

// ShowE runs terraform show with the given options and returns stdout/stderr from the show command.
func ShowE(t *testing.T, options *Options, format Format, path string) (string, error) {
	args := []string{"show"}
	if format == JSON {
		args = append(args, "-json")
	}
	args = append(args, path)
	return RunTerraformCommandE(t, options, FormatArgs(options, args...)...)
}
