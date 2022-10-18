//go:build azure
// +build azure

// NOTE: We use build tags to differentiate azure testing because we currently do not have azure access setup for
// CircleCI.

package test

import (
	"fmt"
	"strings"

	"testing"

	"github.com/gruntwork-io/terratest/modules/azure"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestTerraformAzureActionGroupExample(t *testing.T) {
	t.Parallel()
	_random := strings.ToLower(random.UniqueId())

	expectedResourceGroupName := fmt.Sprintf("tmp-rg-%s", _random)
	expectedAppName := fmt.Sprintf("tmp-asp-%s", _random)

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
	defer terraform.Destroy(t, terraformOptions)

	// This will run `terraform init` and `terraform apply` and fail the test if there are any errors
	terraform.InitAndApply(t, terraformOptions)

	assert := assert.New(t)
	actionGroupId := terraform.Output(t, terraformOptions, "action_group_id")
	assert.NotNil(actionGroupId)
	assert.Contains(actionGroupId, expectedAppName)

	actionGroup := azure.GetActionGroupResource(t, expectedAppName, expectedResourceGroupName, "")

	assert.NotNil(actionGroup)
	assert.Equal(1, len(*actionGroup.EmailReceivers))
	assert.Equal(0, len(*actionGroup.SmsReceivers))
	assert.Equal(1, len(*actionGroup.WebhookReceivers))
}
