package modularization_test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/aws/v2"
	"github.com/gruntwork-io/terratest/modules/azure/v2"
	"github.com/gruntwork-io/terratest/modules/collections/v2"
	dns_helper "github.com/gruntwork-io/terratest/modules/dns-helper/v2"
	"github.com/gruntwork-io/terratest/modules/environment/v2"
	"github.com/gruntwork-io/terratest/modules/gcp/v2"
	"github.com/gruntwork-io/terratest/modules/git/v2"
	"github.com/gruntwork-io/terratest/modules/helm/v2"
	http_helper "github.com/gruntwork-io/terratest/modules/http-helper/v2"
	"github.com/gruntwork-io/terratest/modules/k8s/v2"
	"github.com/gruntwork-io/terratest/modules/logger/v2"
	"github.com/gruntwork-io/terratest/modules/oci/v2"
	"github.com/gruntwork-io/terratest/modules/ssh/v2"
	"github.com/gruntwork-io/terratest/modules/terraform/v2"
	"github.com/gruntwork-io/terratest/modules/terragrunt/v2"
	test_structure "github.com/gruntwork-io/terratest/modules/testing/v2"
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

	// Test tier 2: modules/oci, modules/aws, modules/azure, modules/gcp
	_ = oci.GetRootCompartmentID
	_ = aws.GetRandomRegion                 // AWS module
	_ = azure.GetSubscriptionClientE        // Azure module
	_ = gcp.GetGoogleProjectIDFromEnvVar    // GCP module

	// Test tier 3: modules/terragrunt, modules/http-helper, modules/dns-helper, modules/ssh
	options := &terragrunt.Options{
		TerragruntDir: "/path/to/terragrunt",
	}
	assert.NotNil(t, options)
	_ = http_helper.HttpGet
	_ = dns_helper.DNSLookupAuthoritative
	_ = ssh.CheckSshCommand

	// Test tier 4: modules/terraform, modules/k8s, modules/helm (high-level modules)
	tfOptions := &terraform.Options{
		TerraformDir: "/path/to/terraform",
	}
	assert.NotNil(t, tfOptions)
	_ = k8s.GetKubeConfigPathE              // K8s module
	_ = helm.Install                        // Helm module
}
