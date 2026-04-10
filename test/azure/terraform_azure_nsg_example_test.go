//go:build azure
// +build azure

// NOTE: We use build tags to differentiate azure testing because we currently do not have azure access setup for
// CircleCI.

package test_test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/azure"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTerraformAzureNsgExample(t *testing.T) {
	t.Parallel()

	randomPostfixValue := random.UniqueID()

	// Construct options for TF apply
	terraformOptions := &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: "../../examples/azure/terraform-azure-nsg-example",
		Vars: map[string]interface{}{
			"postfix": randomPostfixValue,
		},
	}

	defer terraform.DestroyContext(t, t.Context(), terraformOptions)

	terraform.InitAndApplyContext(t, t.Context(), terraformOptions)

	resourceGroupName := terraform.OutputContext(t, t.Context(), terraformOptions, "resource_group_name")
	nsgName := terraform.OutputContext(t, t.Context(), terraformOptions, "nsg_name")
	sshRuleName := terraform.OutputContext(t, t.Context(), terraformOptions, "ssh_rule_name")
	httpRuleName := terraform.OutputContext(t, t.Context(), terraformOptions, "http_rule_name")

	// A default NSG has 6 rules, and we have two custom rules for a total of 8
	rules, err := azure.GetAllNSGRulesContextE(t.Context(), resourceGroupName, nsgName, "")
	require.NoError(t, err)
	assert.Len(t, rules.SummarizedRules, 8)

	// We should have a rule for allowing ssh
	sshRule := rules.FindRuleByName(sshRuleName)

	// That rule should allow port 22 inbound
	assert.True(t, sshRule.AllowsDestinationPort(t, "22"))

	// But should not allow 80 inbound
	assert.False(t, sshRule.AllowsDestinationPort(t, "80"))

	// SSh is allowed from any port
	assert.True(t, sshRule.AllowsSourcePort(t, "*"))

	// We should have a rule for blocking HTTP
	httpRule := rules.FindRuleByName(httpRuleName)

	// This rule should BLOCK port 80 inbound
	assert.False(t, httpRule.AllowsDestinationPort(t, "80"))
}
