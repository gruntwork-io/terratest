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
func GraphE(t testing.TestingT, options *Options) (string, error) {
	return runTerragruntCommandE(t, options, "dag", "graph")
}
