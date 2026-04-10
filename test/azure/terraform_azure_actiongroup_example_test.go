//go:build azure
// +build azure

// NOTE: We use build tags to differentiate azure testing because we currently do not have azure access setup for
// CircleCI.

package test_test

import (
	"strings"

	"testing"

	"github.com/gruntwork-io/terratest/modules/azure"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestTerraformAzureActionGroupExample(t *testing.T) {
	t.Parallel()

	_random := strings.ToLower(random.UniqueID())

	expectedResourceGroupName := "tmp-rg-" + _random
	expectedAppName := "tmp-asp-" + _random

	terraformOptions := &terraform.Options{
		TerraformDir: "../../examples/azure/terraform-azure-actiongroup-example",
		Vars: map[string]interface{}{
			"resource_group_name": expectedResourceGroupName,
			"app_name":            expectedAppName,
			"location":            "westus2",
			"short_name":          "blah",
			"enable_email":        true,
			"email_name":          "emailTestName",
			"email_address":       "sample@test.com",
			"enable_webhook":      true,
			"webhook_name":        "webhookTestName",
			"webhook_service_uri": "http://example.com/alert",
		},
	}
	// At the end of the test, run `terraform destroy` to clean up any resources that were created
	defer terraform.DestroyContext(t, t.Context(), terraformOptions)

	// This will run `terraform init` and `terraform apply` and fail the test if there are any errors
	terraform.InitAndApplyContext(t, t.Context(), terraformOptions)

	assert := assert.New(t)
	actionGroupID := terraform.OutputContext(t, t.Context(), terraformOptions, "action_group_id")
	assert.NotNil(actionGroupID)
	assert.Contains(actionGroupID, expectedAppName)

	actionGroup := azure.GetActionGroupResourceContext(t, t.Context(), expectedAppName, expectedResourceGroupName, "")

	assert.NotNil(actionGroup)
	assert.Len(actionGroup.Properties.EmailReceivers, 1)
	assert.Empty(actionGroup.Properties.SmsReceivers)
	assert.Len(actionGroup.Properties.WebhookReceivers, 1)
}
