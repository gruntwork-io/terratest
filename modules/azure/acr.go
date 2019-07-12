package azure

import (
	"context"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/containerregistry/mgmt/2017-10-01/containerregistry"
)

func registriesClientE(subscriptionID string) (*containerregistry.RegistriesClient, error) {
	authorizer, err := DeploymentServicePrincipalAuthorizer()
	if err != nil {
		return nil, err
	}
	client := containerregistry.NewRegistriesClient(subscriptionID)
	client.Authorizer = authorizer
	return &client, nil
}

func webhookClientE(subscriptionID string) (*containerregistry.WebhooksClient, error) {
	authorizer, err := DeploymentServicePrincipalAuthorizer()
	if err != nil {
		return nil, err
	}
	client := containerregistry.NewWebhooksClient(subscriptionID)
	client.Authorizer = authorizer
	return &client, nil
}

// ACRNetworkAclsE - Return the newtwork ACLs for an ACR instance
func ACRNetworkAclsE(subscriptionID string, resourceGroupName string, acrName string) (*containerregistry.NetworkRuleSet, error) {

	client, err := registriesClientE(subscriptionID)
	if err != nil {
		return nil, err
	}

	acr, err := client.Get(context.Background(), resourceGroupName, acrName)
	if err != nil {
		return nil, err
	}

	return acr.NetworkRuleSet, nil
}

// ACRNetworkAcls - Like ACRNetworkAclsE but fails in the case an error is returned
func ACRNetworkAcls(t *testing.T, subscriptionID string, resourceGroupName string, acrName string) *containerregistry.NetworkRuleSet {
	acls, err := ACRNetworkAclsE(subscriptionID, resourceGroupName, acrName)
	if err != nil {
		t.Fatal(err)
	}
	return acls
}

// ACRWebHookE - Return ACR Webhook definition
func ACRWebHookE(subscriptionID string, resourceGroupName string, acrName string, webhookName string) (*containerregistry.Webhook, error) {
	client, err := webhookClientE(subscriptionID)
	if err != nil {
		return nil, err
	}

	webhook, err := client.Get(context.Background(), resourceGroupName, acrName, webhookName)
	if err != nil {
		return nil, err
	}

	return &webhook, nil
}

// ACRWebHook - Like ACRWebHookE but fails in the case an error is returned
func ACRWebHook(t *testing.T, subscriptionID string, resourceGroupName string, acrName string, webhookName string) *containerregistry.Webhook {
	webhooks, err := ACRWebHookE(subscriptionID, resourceGroupName, acrName, webhookName)
	if err != nil {
		t.Fatal(err)
	}
	return webhooks
}

// ACRWebHookCallbackE - Get callback config for a webhook
func ACRWebHookCallbackE(subscriptionID string, resourceGroupName string, acrName string, webhookName string) (*containerregistry.CallbackConfig, error) {
	client, err := webhookClientE(subscriptionID)
	if err != nil {
		return nil, err
	}

	webhookCallback, err := client.GetCallbackConfig(context.Background(), resourceGroupName, acrName, webhookName)
	if err != nil {
		return nil, err
	}

	return &webhookCallback, nil
}

// ACRWebHookCallback - Like ACRWebHookCallbackE but fails in the case an error is returned
func ACRWebHookCallback(t *testing.T, subscriptionID string, resourceGroupName string, acrName string, webhookName string) *containerregistry.CallbackConfig {
	webhookCallback, err := ACRWebHookCallbackE(subscriptionID, resourceGroupName, acrName, webhookName)
	if err != nil {
		t.Fatal(err)
	}
	return webhookCallback
}
