package terragrunt

import (
	"context"

	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// GraphContext runs terragrunt dag graph and returns the DOT-format dependency graph.
// The provided context is passed through to the underlying command execution, allowing for timeout
// and cancellation control. This is useful for verifying dependency relationships between terragrunt units.
// Log lines are stripped from the output so the result is clean DOT format.
func GraphContext(t testing.TestingT, ctx context.Context, options *Options) string {
	out, err := GraphContextE(t, ctx, options)
	require.NoError(t, err)

	return out
}

// GraphContextE runs terragrunt dag graph and returns the DOT-format dependency graph.
// The provided context is passed through to the underlying command execution, allowing for timeout
// and cancellation control. This is useful for verifying dependency relationships between terragrunt units.
// Log lines are stripped from the output so the result is clean DOT format.
func GraphContextE(t testing.TestingT, ctx context.Context, options *Options) (string, error) {
	rawOutput, err := runTerragruntCommandE(t, ctx, options, "dag", "graph")
	if err != nil {
		return "", err
	}

	return FilterLogLines(rawOutput), nil
}

// Graph runs terragrunt dag graph and returns the DOT-format dependency graph.
// This is useful for verifying dependency relationships between terragrunt units.
//
// Deprecated: Use [GraphContext] instead.
func Graph(t testing.TestingT, options *Options) string {
	return GraphContext(t, context.Background(), options)
}

// GraphE runs terragrunt dag graph and returns the DOT-format dependency graph.
// This is useful for verifying dependency relationships between terragrunt units.
// Log lines are stripped from the output so the result is clean DOT format.
//
// Deprecated: Use [GraphContextE] instead.
func GraphE(t testing.TestingT, options *Options) (string, error) {
	return GraphContextE(t, context.Background(), options)
}
