//go:build azure
// +build azure

// NOTE: We use build tags to differentiate azure testing because we currently do not have azure access setup for
// CircleCI.

package test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/azure"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestTerraformAzureNsgExample(t *testing.T) {
	t.Parallel()

	randomPostfixValue := random.UniqueId()

	// Construct options for TF apply
	terraformOptions := &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: "../../examples/azure/terraform-azure-nsg-example",
		Vars: map[string]interface{}{
			"postfix": randomPostfixValue,
		},
	}

	defer terraform.Destroy(t, terraformOptions)
	terraform.InitAndApply(t, terraformOptions)

	resourceGroupName := terraform.Output(t, terraformOptions, "resource_group_name")
	nsgName := terraform.Output(t, terraformOptions, "nsg_name")
	sshRuleName := terraform.Output(t, terraformOptions, "ssh_rule_name")
	httpRuleName := terraform.Output(t, terraformOptions, "http_rule_name")

	// A default NSG has 6 rules, and we have two custom rules for a total of 8
	rules, err := azure.GetAllNSGRulesE(resourceGroupName, nsgName, "")
	assert.NoError(t, err)
	assert.Equal(t, 8, len(rules.SummarizedRules))

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
