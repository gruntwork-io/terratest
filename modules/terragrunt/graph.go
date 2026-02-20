package terragrunt

import (
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// Graph runs terragrunt dag graph and returns the DOT-format dependency graph.
// This is useful for verifying dependency relationships between terragrunt units.
func Graph(t testing.TestingT, options *Options) string {
	out, err := GraphE(t, options)
	require.NoError(t, err)
	return out
}

// GraphE runs terragrunt dag graph and returns the DOT-format dependency graph.
// This is useful for verifying dependency relationships between terragrunt units.
// Log lines are stripped from the output so the result is clean DOT format.
func GraphE(t testing.TestingT, options *Options) (string, error) {
	rawOutput, err := runTerragruntCommandE(t, options, "dag", "graph")
	if err != nil {
		return "", err
	}

	return filterLogLines(rawOutput), nil
}
