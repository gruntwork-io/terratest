# Terraform Azure Eventhub Example

This folder contains a simple Terraform module that deploys resources in [Azure](https://azure.microsoft.com/) to demonstrate
how you can use Terratest to write automated tests for your Azure Terraform code. This module deploys an [Azure Event Hubs](https://azure.microsoft.com/en-us/products/event-hubs/).

Check out [test/azure/terraform_azure_eventhub_example_test.go](/test/azure/terraform_azure_eventhub_example_test.go) to see how you can write automated tests for this module.

Note that the resources deployed in this module don't actually do anything; it just runs the resources for demonstration purposes.

**WARNING**: This module and the automated tests for it deploy real resources into your Azure account which can cost you money. The resources are all part of the [Azure Free Account](https://azure.microsoft.com/free/), so if you haven't used that up, it should be free, but you are completely responsible for all Azure charges.

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

1. Sign up for [Azure](https://azure.microsoft.com/).
2. Configure your Azure credentials using one of the [supported methods for Azure CLI tools](https://docs.microsoft.com/cli/azure/azure-cli-configuration?view=azure-cli-latest).

3. Install [Terraform](https://www.terraform.io/) and make sure it's on your `PATH`.
4. [Review environment variables](#review-environment-variables).
5. Install [Golang](https://golang.org/) and make sure this code is checked out into your `GOPATH`.
6. `cd test`
7. Make sure [the azure-sdk-for-go versions match](#check-go-dependencies) in [/test/go.mod](/test/go.mod) and in [test/azure/terraform_azure_eventhub_example_test.go](/test/azure/terraform_azure_eventhub_example_test.go).
8. `go test -v -run TestTerraformAzureEventhubkExample`
