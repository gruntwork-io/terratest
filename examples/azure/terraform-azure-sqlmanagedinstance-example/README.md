# Terraform Azure SQL DB Example

This folder contains a Terraform module that deploys resources in [Azure](https://azure.microsoft.com/) to demonstrate how you can use Terratest to write automated tests for your Azure Terraform code. This module deploys a SQL Managed Instance, and a SQL Managed Instance database.

- A [SQL Managed Instance](https://azure.microsoft.com/en-us/products/azure-sql/managed-instance/).
- A SQL Managed Database.

Check out [test/azure/terraform_azure_sqlmanagedinstance_example_test.go](./../../../test/azure/terraform_azure_sqlmanagedinstance_example_test.go) to see how you can write automated tests for this module and validate the configuration of the parameters and options. 

**WARNING**: This module and the automated tests for it deploy real resources into your Azure account which can cost you money. The resources are all part of the [Azure Free Account](https://azure.microsoft.com/en-us/free/), so if you haven't used that up,
it should be free, but you are completely responsible for all Azure charges.

## Running this module manually
1. Sign up for [Azure](https://azure.microsoft.com/).
1. Configure your Azure credentials using one of the [supported methods for Azure CLI
   tools](https://docs.microsoft.com/en-us/cli/azure/azure-cli-configuration?view=azure-cli-latest)
1. Install [Terraform](https://www.terraform.io/) and make sure it's on your `PATH`.
1. Ensure [environment variables](../README.md#review-environment-variables) are available
1. Run `terraform init`
1. Run `terraform apply`
1. When you're done, run `terraform destroy`.


## Running automated tests against this module

**WARNING**: The deploymnet for this module usually takes more than 4-6 hours as stated in the [microsoft docs](https://learn.microsoft.com/en-us/azure/azure-sql/managed-instance/management-operations-overview?view=azuresql#duration), so please make sure to set the timeout accordingly in the below go test command.

1. Sign up for [Azure](https://azure.microsoft.com/)
2. Configure your Azure credentials using one of the [supported methods for Azure CLI
   tools](https://docs.microsoft.com/en-us/cli/azure/azure-cli-configuration?view=azure-cli-latest)
3. Install [Terraform](https://www.terraform.io/) and make sure it's on your `PATH`
4. Configure your Terratest [Go test environment](../README.md) 
5. `cd test/azure`
6. `go build terraform_azure_sqlmanagedinstance_example_test.go`
7. `go test -v -run TestTerraformAzureSQLManagedInstanceExample -timeout <in hours>`
