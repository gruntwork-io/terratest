# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY AN AZURE MANAGED DISK
# This is an example of how to deploy a managed disk.
# ---------------------------------------------------------------------------------------------------------------------
# See test/azure/terraform_azure_disk_example_test.go for how to write automated tests for this code.
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
