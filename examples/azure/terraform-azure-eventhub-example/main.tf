# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY AN AZURE EVENT HUB
# This is an example of how to deploy a event hub.
# ---------------------------------------------------------------------------------------------------------------------
# See test/azure/terraform_azure_eventhub_example_test.go for how to write automated tests for this code.
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
      source  = "hashicorp/azurerm"
      version = "~> 2.29"
    }
  }
}

# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY A RESOURCE GROUP
# ---------------------------------------------------------------------------------------------------------------------

resource "azurerm_resource_group" "eventhub_rg" {
  name     = "terratest-eventhub-rg-${var.postfix}"
  location = var.location
}

resource "azurerm_eventhub_namespace" "eventhub_namespace" {
    name                = "terratest-eventhub-namespace-${var.postfix}"
    location            = azurerm_resource_group.eventhub_rg.location
    resource_group_name = azurerm_resource_group.eventhub_rg.name
    sku                 = "Standard"
}

resource "azurerm_eventhub" "eventhub" {
    name                = "terratest-eventhub-${var.postfix}"
    namespace_name      = azurerm_eventhub_namespace.eventhub_namespace.name
    resource_group_name = azurerm_resource_group.eventhub_rg.name
    partition_count     = 2
    message_retention   = 1
}