# Terraform Azure Automation Account Example

This folder contains a Terraform module that deploys resources in [Azure](https://azure.microsoft.com/) to demonstrate how you can use Terratest to write automated tests for your Azure Terraform code. This module deploys the following:

- An [Automation Account](https://azure.microsoft.com/services/automation/) that provides the module the following:
  - `Automation Account` with the name specified in the `automation_account_name` variable.
  - `Automation Account Connection Run As Account` with the name specified in the `automation_run_as_connection_name` variable.  See the section titled *_Example Service Principal and Certificate Setup_* below on how to configure the Automation Account RunAs Account credentials.
  - `Automation Account Connection Run As Certificate` with the name specified in the `automation_run_as_certificate_name` variable and the thumbprint in the `TF_VAR_automation_run_as_certificate_thumbprint` variable.
  - `Automation Account Connection Type` with type specified in the `automation_run_as_connection_type` variable.
  - [Desired State Configuration](https://docs.microsoft.com/powershell/scripting/dsc/getting-started/winGettingStarted?view=powershell-7#:~:text=Get%20started%20with%20Desired%20State%20Configuration%20%28DSC%29%20for,Windows%20PowerShell%20Desired%20State%20Configuration%20log%20files.%20) with the name specified in the `sample_dsc_name` variable and the path to the DSC specified in the `sample_dsc_path` variable.
- [Virtual Machine](https://docs.microsoft.com/azure/virtual-machines/) with the name specified in the `vm_name` output variable.
  - [Virtual Machine Extension](https://docs.microsoft.com/azure/virtual-machines/extensions/overview#:~:text=Troubleshoot%20extensions%20%20%20%20Namespace%20%20,Encryption%20for%20Windows%20%2012%20more%20rows%20), configured for DSC, with the `NodeConigurationName` set in the `sample_dsc_configuration_name` variable.
  - The VM includes a virtual network, a subnet, and a network interface with hard-coded configuration as it is not pertinent to this example.

Check out [test/azure/terraform_azure_automationaccount_example_test.go](/test/azure/terraform_azure_automationaccount_example_test.go) to see how you can write
automated tests for this module.

Note that the resources deployed in this module don't actually do anything; it just runs the resources for demonstration purposes.

**WARNING**: This module and the automated tests for it deploy real resources into your Azure account which can cost you
money. The resources are all part of the [Azure Free Account](https://azure.microsoft.com/free/), so if you haven't used that up,
it should be free, but you are completely responsible for all Azure charges.

## Example Service Principal and Certificate Setup

To run this example, you must create a service principal in Azure Active Directory as well as create a non-password protected self-signed certificate in .pfx format that you will need to upload into the Automation Account Run As service principal in Azure Active Directory as a secret for test purposes.  The self-signed certificate must be base64 encoded when passed into the `azurerm_automation_certificate` resource on the `base64` property.  Here are three different methods to pass in the certificate:

1. Generate, store, and retrieve the self-signed certificate as an output variable from a deployed Azure Key Vault, assigning it to the `base64` property using the [base64encode](https://www.terraform.io/docs/language/functions/base64encode.html)
2. Generate a self-signed certificate in the file system and pass the path to the certificate into the [filebase64](https://www.terraform.io/docs/language/functions/filebase64.html) terraform function when assigning to the `base64` property.
3. Manually base64 encode the self-signed certificate and make available to CI as a github.secret value, which is what this example does.  

In this case for the test example, the test self-signed certificate is base64 encoded and then set on the `TF_VAR_automation_run_as_certificate_base64` environment variable.  As an alternative, you could read the certificate from Key Vault as described above, which would be the recommended procedure in a production environment.

The documentation link [Manage an Azure Automation Run As account](https://docs.microsoft.com/azure/automation/manage-runas-account#:~:text=1%20Go%20to%20your%20Automation%20account%20and%20select,locate%20the%20role%20definition%20that%20is%20being%20used.) has additional background on the service principal configuration requirements.  

Some additional configuration is required.  The `TF_VAR_client_id` must be configured with the corresponding service principal client ID where the self-signed certificate was previously added as a secret in AAD. For the Automation Account Run As connection certificate, set the base64 encoded certificate value on the `TF_VAR_automation_run_as_certificate_base64` environment variable.  Also set the certificate thumbprint in the `TF_VAR_automation_run_as_certificate_thumbprint` variable.

In general when uploading a DSC to an Automation Account, you will need to kick off compilation of the DSC in the Automation Account prior to applying the DSC to a VM node, else it will fail to apply.  You can use PowerShell Core to compile the DSC in Terraform via a `null_resource`.  The `null_resource` named `azureSignInPWSH` first performs a sign-in to Azure from PowerShell Core. The `TF_VAR_client_id` and `TF_VAR_client_secret` environment variables must be configured with a service principal to sign-in to Azure from PowerShell core and compile the DSC.  The service principal can be the same service principal used to execute the terraform, but it is not required to use the same service principal. The `null_resource` named `compileSampleDSC` performs compilation of the DSC in the Azure Automation account. 

*_Warning: The `null_resource` `"azureSignInPWSH"` must sign-in to Azure using credentials, could be available in the logs if logging is enabled in Terraform. to avoid this, the `client_id` and `client_secret` variables have `sensitive = true` configured to prevent the values from being written to the log in order to mitigate this risk_*

*_Note: In a production system, you would store the service principal configuration in and create the certificate using Azure Key Vault and then configure a Terraform `azurerm_key_vault_secret` data source on the Key Vault instance to access the data securely directly from Terraform._*

## Running this module manually

1. Sign up for [Azure](https://azure.microsoft.com/)
1. Configure your Azure credentials using one of the [supported methods for Azure CLI
   tools](https://docs.microsoft.com/cli/azure/azure-cli-configuration?view=azure-cli-latest)
1. Install [Terraform](https://www.terraform.io/) and make sure it's on your `PATH`
1. Ensure [environment variables](../README.md#review-environment-variables) are available
1. Run `terraform init`
1. Run `terraform apply`
1. When you're done, run `terraform destroy`

## Running automated tests against this module

1. Sign up for [Azure](https://azure.microsoft.com/)
1. Configure your Azure credentials using one of the [supported methods for Azure CLI
   tools](https://docs.microsoft.com/cli/azure/azure-cli-configuration?view=azure-cli-latest)
1. Install [Terraform](https://www.terraform.io/) and make sure it's on your `PATH`
1. Configure your Terratest [Go test environment](../README.md)
1. `cd test/azure`
1. `go build terraform_azure_automationaccount_example_test.go`
1. `go test -v terraform_azure_automationaccount_example_test.go -timeout 20m` or `go test -tags=azureautomation -v -timeout 20m`


