# Terratest Configuration and Setup

Terratest uses Go to make calls to Azure through the azure-sdk-for-go library and independently confirm the actual Azure resource property matches the expected state provided by Terraform output variables.

- Instructions for running each Azure Terratest module are included in each Terraform example sub-folder:
  - examples/azure/terraform-azure-\*-example/README.md
- Tests which assert against expected Terraform output values are located in the the respective go files of the folder:
  - [test/azure/terraform-azure-\*-example_test.go](../../test/azure)
- Test APIs which provide the actual Azure resource property values via the azure-sdk-for-go are located in the folder:
  - [modules/azure](../../modules/azure)

## Go Dependencies

Install [Golang](https://golang.org/) and make sure this code is checked out into your `GOPATH`

These modules are currently using the latest version of Go and was tested with **go1.14.4**.

## Azure-sdk-for-go version

Let's make sure [go.mod](https://github.com/gruntwork-io/terratest/blob/master/go.mod) includes the appropriate [azure-sdk-for-go version](https://github.com/Azure/azure-sdk-for-go/releases/tag/v46.1.0):

```go
require (
    ...
    github.com/Azure/azure-sdk-for-go v46.1.0+incompatible
    ...
)
```

If we make changes to either the **go.mod** or the **go test file**, we should make sure that the go build command works still.

```powershell
go build terraform_azure_*_test.go
```

## Review Environment Variables

As part of configuring terraform for Azure, we'll want to check that we have set the appropriate [credentials](https://docs.microsoft.com/azure/terraform/terraform-install-configure?toc=https%3A%2F%2Fdocs.microsoft.com%2Fen-us%2Fazure%2Fterraform%2Ftoc.json&bc=https%3A%2F%2Fdocs.microsoft.com%2Fen-us%2Fazure%2Fbread%2Ftoc.json#set-up-terraform-access-to-azure) and also that we set the [environment variables](https://docs.microsoft.com/azure/terraform/terraform-install-configure?toc=https%3A%2F%2Fdocs.microsoft.com%2Fen-us%2Fazure%2Fterraform%2Ftoc.json&bc=https%3A%2F%2Fdocs.microsoft.com%2Fen-us%2Fazure%2Fbread%2Ftoc.json#configure-terraform-environment-variables) on the testing host.

For non-commercial cloud deployments, set the "AZURE_ENVIRONMENT" environment variable to the appropriate cloud environment (this is used by the [client_factory](../../modules/azure/client_factory.go) to set the correct azure endpoints):

```bash
export ARM_CLIENT_ID=your_app_id
export ARM_CLIENT_SECRET=your_password
export ARM_SUBSCRIPTION_ID=your_subscription_id
export ARM_TENANT_ID=your_tenant_id

# AZURE_ENVIRONMENT is the name of the Azure environment to use. Set to one of the following:
export AZURE_ENVIRONMENT=AzureUSGovernmentCloud
export AZURE_ENVIRONMENT=AzureChinaCloud
export AZURE_ENVIRONMENT=AzureGermanCloud
export AZURE_ENVIRONMENT=AzurePublicCloud
export AZURE_ENVIRONMENT=AzureStackCloud
```

Note, in a Windows environment, these should be set as **system environment variables**. We can use a PowerShell console with administrative rights to update these environment variables:

```powershell
[System.Environment]::SetEnvironmentVariable("ARM_CLIENT_ID",$your_app_id,[System.EnvironmentVariableTarget]::Machine)
[System.Environment]::SetEnvironmentVariable("ARM_CLIENT_SECRET",$your_password,[System.EnvironmentVariableTarget]::Machine)
[System.Environment]::SetEnvironmentVariable("ARM_SUBSCRIPTION_ID",$your_subscription_id,[System.EnvironmentVariableTarget]::Machine)
[System.Environment]::SetEnvironmentVariable("ARM_TENANT_ID",$your_tenant_id,[System.EnvironmentVariableTarget]::Machine)
[System.Environment]::SetEnvironmentVariable("AZURE_ENVIRONMENT",$your_azure_env,[System.EnvironmentVariableTarget]::Machine)
```
