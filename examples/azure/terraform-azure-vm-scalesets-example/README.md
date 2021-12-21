# Terraform Azure Virtual Machine Scale Sets Example

This folder contains a complete Terraform VM Scale Sets module that deploys resources in [Azure](https://azure.microsoft.com/)
to demonstrate how you can use Terratest to write automated tests for your Azure Virtual Machine Scale Sets Terraform code.
This module deploys these resources:

- A [Virtual Machine Scale Sets](https://azure.microsoft.com/en-us/services/virtual-machine-scale-sets) and gives VM instances the following resources:
  - [Virtual Machine](https://docs.microsoft.com/azure/virtual-machines)
  - [Managed Disk](https://docs.microsoft.com/azure/virtual-machines/managed-disks-overview)
  - [Network Interface](https://docs.microsoft.com/azure/virtual-network/virtual-network-network-interface)
- A [Virtual Network](https://azure.microsoft.com/services/virtual-network/) module that contains the following resources:
  - [Virtual Network](https://docs.microsoft.com/azure/virtual-network)
  - [Subnet](https://docs.microsoft.com/rest/api/virtualnetwork/subnets)
  - [Public Address](https://docs.microsoft.com/azure/virtual-network/public-ip-addresses)
  - [Load Balancer](https://docs.microsoft.com/en-us/azure/load-balancer)

Check out [test/azure/terraform_azure_vm_scalesets_example_test.go](/test/azure/terraform_azure_vm_scalesets_example_test.go)
to see how you can write automated tests for this module.

**WARNING**: This module and the automated tests for it deploy real resources into your Azure account which can cost you
money. The resources are all part of the [Azure Free Account](https://azure.microsoft.com/free/), so if you haven't used that
up, it should be free, but you are completely responsible for all Azure charges.

## Running this module manually

1. Sign up for [Azure](https://azure.microsoft.com/)
1. Configure your Azure credentials using one of the [supported methods for Azure CL
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
1. `go test -v -timeout 20m -run=TestTerraformAzureVmScaleSetsExample ./...`
