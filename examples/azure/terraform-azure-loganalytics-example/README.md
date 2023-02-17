# Terraform Azure Log Analytics Example

This folder contains a Terraform module that deploys resources in [Azure](https://azure.microsoft.com/) to demonstrate
how you can use TerraTest to write automated tests for your Azure Terraform code. This module deploys a Log Analytics Workspace.

- A [Log Analytics Workspace](https://docs.microsoft.com/azure/azure-monitor/platform/log-analytics-agent) that gives the module the following:
  - [Name](https://docs.microsoft.com/azure/azure-monitor/learn/quick-create-workspace#:~:text=%20Create%20a%20Log%20Analytics%20workspace%20in%20the,and%20region%20as%20in%20the%20deleted...%20More%20)  with the value specified in the `loganalytics_workspace_name`  output variable.
  - [Sku](https://docs.microsoft.com/azure/azure-monitor/learn/quick-create-workspace#:~:text=%20Create%20a%20Log%20Analytics%20workspace%20in%20the,and%20region%20as%20in%20the%20deleted...%20More%20)  with the value specified in the `loganalytics_workspace_sku`  output variable.
  - [RetentionPeriodInDays](https://docs.microsoft.com/azure/azure-monitor/learn/quick-create-workspace#:~:text=%20Create%20a%20Log%20Analytics%20workspace%20in%20the,and%20region%20as%20in%20the%20deleted...%20More%20)  with the value specified in the `loganalytics_workspace_retention`  output variable.

Check out [test/azure/terraform_azure_loganalytics_example_test.go](/test/azure/terraform_azure_loganalytics_example_test.go) to see how you can write
automated tests for this module.

Note that the Log Analytics Workspace in this module don't actually do anything; it just runs the resources for
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
1. Configure your TerraTest [Go test environment](../README.md)
1. `cd test/azure`
1. `go build terraform_azure_loganalytics_example_test.go`
1. `go test -v -run TestTerraformAzureLogAnalyticsExample`
