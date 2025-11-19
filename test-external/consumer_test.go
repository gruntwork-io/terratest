package testexternal

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/terragrunt"
	test_structure "github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/assert"
)

// TestConsumerSimulation validates that external consumers can import and use
// Terratest modules without ambiguous import errors. This test imports a
// representative mix of modules across different dependency tiers.
func TestConsumerSimulation(t *testing.T) {
	t.Parallel()

	// Test tier 0: modules/testing
	_ = test_structure.TestingT(t)

	// Test tier 1: modules/logger
	log := logger.Default
	assert.NotNil(t, log)

	// Test tier 3: modules/terragrunt
	options := &terragrunt.Options{
		TerragruntDir: "/path/to/terragrunt",
	}
	assert.NotNil(t, options)
}
