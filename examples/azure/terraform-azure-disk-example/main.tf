# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY AN AZURE MANAGED DISK
# This is an example of how to deploy a managed disk.
# ---------------------------------------------------------------------------------------------------------------------
# See test/azure/terraform_azure_disk_example_test.go for how to write automated tests for this code.
# ---------------------------------------------------------------------------------------------------------------------

provider "azurerm" {
  features {}
}

terraform {
  # This module is now only being tested with Terraform 0.13.x. However, to make upgrading easier, we are setting
  # 0.12.26 as the minimum version, as that version added support for required_providers with source URLs, making it
  # forwards compatible with 0.13.x code.
  required_version = ">= 0.12.26"
  required_providers {
    azurerm = {
      version = "~> 2.29"
      source  = "hashicorp/azurerm"
    }
  }
}

# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY A RESOURCE GROUP
# ---------------------------------------------------------------------------------------------------------------------

resource "azurerm_resource_group" "disk_rg" {
  name     = "terratest-disk-rg-${var.postfix}"
  location = var.location
}

# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY THE DISK
# ---------------------------------------------------------------------------------------------------------------------

resource "azurerm_managed_disk" "disk" {
  name                 = "disk-${var.postfix}"
  location             = azurerm_resource_group.disk_rg.location
  resource_group_name  = azurerm_resource_group.disk_rg.name
  storage_account_type = var.disk_type
  create_option        = "Empty"
  disk_size_gb         = 10
}
