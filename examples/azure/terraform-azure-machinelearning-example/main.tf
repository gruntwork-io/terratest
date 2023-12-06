# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY AN AZURE MACHINE LEARNING WORSPACE
# This is an example of how to deploy an Azure machine learning workspace with compute
# ---------------------------------------------------------------------------------------------------------------------
# See test/azure/terraform_azure_machine_learning_example_test.go for how to write automated tests for this code.
# ---------------------------------------------------------------------------------------------------------------------

provider "azurerm" {
  features {
    key_vault {
      purge_soft_delete_on_destroy = true
    }
  }
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

data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "ml_rg" {
  name     = "rg-ml-${var.postfix}"
  location = var.location
}

resource "azurerm_application_insights" "ml_appinsights" {
  name                = "ai-mlworspace-${var.postfix}"
  location            = azurerm_resource_group.ml_rg.location
  resource_group_name = azurerm_resource_group.ml_rg.name
  application_type    = "web"
}

resource "azurerm_key_vault" "ml_keyvault" {
  name                = "kv-mlworkspace-${var.postfix}"
  location            = azurerm_resource_group.ml_rg.location
  resource_group_name = azurerm_resource_group.ml_rg.name
  tenant_id           = data.azurerm_client_config.current.tenant_id
  sku_name            = "premium"
}

resource "azurerm_storage_account" "ml_storageaccount" {
  name                     = "samlwkspce${var.postfix}"
  location                 = azurerm_resource_group.ml_rg.location
  resource_group_name      = azurerm_resource_group.ml_rg.name
  account_tier             = "Standard"
  account_replication_type = "GRS"
}

resource "azurerm_machine_learning_workspace" "ml_workspace" {
  name                    = "mlworkspace-${var.postfix}"
  location                = azurerm_resource_group.ml_rg.location
  resource_group_name     = azurerm_resource_group.ml_rg.name
  application_insights_id = azurerm_application_insights.ml_appinsights.id
  key_vault_id            = azurerm_key_vault.ml_keyvault.id
  storage_account_id      = azurerm_storage_account.ml_storageaccount.id

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_machine_learning_compute_cluster" "ml_compute_cluster" {
  name                          = "mlcomputecluster-${var.postfix}"
  location                      = azurerm_resource_group.ml_rg.location
  vm_priority                   = "LowPriority"
  vm_size                       = "STANDARD_NC4AS_T4_V3"
  machine_learning_workspace_id = azurerm_machine_learning_workspace.ml_workspace.id
  scale_settings {
    min_node_count                       = 0
    max_node_count                       = 4
    scale_down_nodes_after_idle_duration = "PT2M"
  }
}
