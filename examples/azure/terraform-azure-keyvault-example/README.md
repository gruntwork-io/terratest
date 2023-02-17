# Terraform Azure Keyvault Example

This folder contains a simple Terraform module that deploys resources in [Azure](https://azure.microsoft.com/) to demonstrate
how you can use Terratest to write automated tests for your Azure Terraform code. This module deploys a Key Vault with one secret, key, and certificate each.

- A [Key Vault](https://azure.microsoft.com/services/key-vault/) that gives the module the following:
  - [Secret](https://docs.microsoft.com/azure/key-vault/general/about-keys-secrets-certificates)  with the value specified in the `secret_name`  output variable.
  - [Key](https://docs.microsoft.com/azure/key-vault/general/about-keys-secrets-certificates)  with the value specified in the `key_name`  output variable.
  - [Certificate](https://docs.microsoft.com/azure/key-vault/general/about-keys-secrets-certificates)  with the value specified in the `certificate_name`  output variable.

Check out [test/azure/terraform_azure_keyvault_example_test.go](/test/azure/terraform_azure_keyvault_example_test.go) to see how you can write
automated tests for this module.

Note that the Key Vault in this module don't actually do anything; it just runs the resources for
demonstration purposes.

**WARNING**: This module and the automated tests for it deploy real resources into your Azure account which can cost you
money. The resources are all part of the [Azure Free Account](https://azure.microsoft.com/free/), so if you haven't used that up,
it should be free, but you are completely responsible for all Azure charges.

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
1. `go build terraform_azure_keyvault_example_test.go`
1. `go test -v -run TestTerraformAzureKeyVaultExample`
