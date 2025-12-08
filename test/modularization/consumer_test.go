package modularization_test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/collections"
	dns_helper "github.com/gruntwork-io/terratest/modules/dns-helper"
	"github.com/gruntwork-io/terratest/modules/environment"
	"github.com/gruntwork-io/terratest/modules/git"
	http_helper "github.com/gruntwork-io/terratest/modules/http-helper"
	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/oci"
	"github.com/gruntwork-io/terratest/modules/ssh"
	"github.com/gruntwork-io/terratest/modules/terragrunt"
	test_structure "github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/assert"
)

// TestConsumerSimulation validates that external consumers can import and use
// Terratest modules without ambiguous import errors. This test imports a
// representative mix of modules across different dependency tiers.
func TestConsumerSimulation(t *testing.T) {
	t.Parallel()

	// Test tier 0: modules/testing and modules/collections
	_ = test_structure.TestingT(t)
	result := collections.ListContains([]string{"foo", "bar"}, "foo")
	assert.True(t, result)

	// Test tier 1: modules/logger, modules/environment, modules/git
	log := logger.Default
	assert.NotNil(t, log)
	_ = environment.GetFirstNonEmptyEnvVarOrEmptyString(t, []string{"PATH"})
	_ = git.GetCurrentBranchName

	// Test tier 2: modules/oci
	_ = oci.GetRootCompartmentID

	// Test tier 3: modules/terragrunt, modules/http-helper, modules/dns-helper, modules/ssh
	options := &terragrunt.Options{
		TerragruntDir: "/path/to/terragrunt",
	}
	assert.NotNil(t, options)
	_ = http_helper.HttpGet
	_ = dns_helper.DNSLookupAuthoritative
	_ = ssh.CheckSshCommand
}
