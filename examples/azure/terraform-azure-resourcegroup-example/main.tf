# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY A RESOURCE GROUP
# This is an example of how to deploy a Resource Group
# ---------------------------------------------------------------------------------------------------------------------
# See test/azure/terraform_azure_resourcegroup_example_test.go for how to write automated tests for this code.
# ---------------------------------------------------------------------------------------------------------------------

terraform {
  required_version = ">= 1.0"

  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "~> 4.0"
    }
  }
}

provider "azurerm" {
  features {}
}

# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY A RESOURCE GROUP
# ---------------------------------------------------------------------------------------------------------------------

resource "azurerm_resource_group" "resource_group" {
  name     = "terratest-rg-${var.postfix}"
  location = var.location
}
