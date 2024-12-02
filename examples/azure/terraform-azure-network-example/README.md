# Terraform Azure Network Example

This folder contains a simple Terraform module that deploys resources in [Azure](https://azure.microsoft.com/) to demonstrate
how you can use Terratest to write automated tests for your Azure Terraform code. This module deploys to a Virtual Network two Network Interface Cards, one with an internal only IP and another with an internal and external Public IP.

- A [Virtual Network](https://azure.microsoft.com/en-us/services/virtual-network/) module that includes the following resources:
  - [Virtual Network](https://docs.microsoft.com/en-us/azure/virtual-network/) with the name specified in the `virtual_network_name` variable.
  - [Subnet](https://docs.microsoft.com/en-us/rest/api/virtualnetwork/subnets) with the name specified in the `subnet_name` variable.
  - [Public Address](https://docs.microsoft.com/en-us/azure/virtual-network/public-ip-addresses) with the name specified in the `public_ip_name` variable.
  - [Internal Network Interface](https://docs.microsoft.com/en-us/azure/virtual-network/virtual-network-network-interface) with the name specified in the `network_interface_internal` variable.
  - [ExternalNetwork Interface](https://docs.microsoft.com/en-us/azure/virtual-network/virtual-network-network-interface) with the name specified in the `network_interface_external` variable.

Check out [test/azure/terraform_azure_network_example_test.go](/test/azure/terraform_azure_network_example_test.go) to see how you can write
automated tests for this module.

Note that the Azure Virtual Network, Subnet, Network Interface and Public IP resources in this module don't actually do anything; it just runs the resources for
demonstration purposes.

**WARNING**: This module and the automated tests for it deploy real resources into your Azure account which can cost you
money. The resources are all part of the [Azure Free Account](https://azure.microsoft.com/en-us/free/), so if you haven't used that up,
it should be free, but you are completely responsible for all Azure charges.

## Running this module manually

1. Sign up for [Azure](https://azure.microsoft.com/)
1. Configure your Azure credentials using one of the [supported methods for Azure CLI
   tools](https://docs.microsoft.com/en-us/cli/azure/azure-cli-configuration?view=azure-cli-latest)
1. Install [Terraform](https://www.terraform.io/) and make sure it's on your `PATH`
1. Ensure [environment variables](../README.md#review-environment-variables) are available
1. Run `terraform init`
1. Run `terraform apply`
1. When you're done, run `terraform destroy`

## Running automated tests against this module

1. Sign up for [Azure](https://azure.microsoft.com/)
1. Configure your Azure credentials using one of the [supported methods for Azure CLI
   tools](https://docs.microsoft.com/en-us/cli/azure/azure-cli-configuration?view=azure-cli-latest)
1. Install [Terraform](https://www.terraform.io/) and make sure it's on your `PATH`
1. Configure your Terratest [Go test environment](../README.md)
1. `cd test/azure`
1. `go build terraform_azure_network_example_test.go`
1. `go test -v -run TestTerraformAzureNetworkExample`
