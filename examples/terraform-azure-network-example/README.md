# Terraform Azure Example

This folder contains a simple Terraform module that deploys network-related resources in [Azure](https://azure.microsoft.com/) to demonstrate
how you can use Terratest to write automated tests for your Azure Terraform code.

This module deploys the following resources :
  - a Resource Group
  - a [Virtual Network](https://azure.microsoft.com/en-us/services/virtual-network/)
  - 2 subnets within the Virtual Network
  - a [Network Security Group](https://docs.microsoft.com/en-us/azure/virtual-network/security-overview) for each subnet
  - a Network Security Rule associated with the first Network Security Group
  - a Public IP (with a public DNS name)

Check out [test/terraform_azure_network_example_test.go](/test/terraform_azure_network_example_test.go) to see how you can write
automated tests for this module.

Note that the Public IP in this module is just for demonstration purposes.
It is not attached to any resource (Virtual Machine, Load Balancer, etc...), so while this address is reachable from the internet, connecting to it will not result in anything.

**WARNING**: This module and the automated tests for it deploy real resources into your Azure account which can cost you
money. The resources are all part of the [Azure Free Account](https://azure.microsoft.com/en-us/free/), so if you haven't used that up,
it should be free, but you are completely responsible for all Azure charges.

## Running this module manually

1. Sign up for an [Azure](https://azure.microsoft.com/) account.
1. Configure your Azure credentials using [the Azure CLI](https://www.terraform.io/docs/providers/azurerm/auth/azure_cli.html), or environment variables supported by [the Azure Terraform provider](https://www.terraform.io/docs/providers/azurerm/index.html#argument-reference).
1. Install [Terraform](https://www.terraform.io/) and make sure it's on your `PATH`.
1. Run `terraform init`.
1. Run `terraform apply`.
1. When you're done, run `terraform destroy`.

## Running automated tests against this module

1. Sign up for an [Azure](https://azure.microsoft.com/) account.
1. Configure your Azure credentials using [the Azure CLI](https://www.terraform.io/docs/providers/azurerm/auth/azure_cli.html), or environment variables supported by [the Azure Terraform provider](https://www.terraform.io/docs/providers/azurerm/index.html#argument-reference).
1. Install [Terraform](https://www.terraform.io/) and make sure it's on your `PATH`.
1. Install [Golang](https://golang.org/) and make sure this code is checked out into your `GOPATH`.
1. `cd test`
1. `dep ensure`
1. `go test -v -run TestTerraformAzureNetworkExample -timeout 20m`
