# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY AN AZURE MONITOR DIAGNOSTIC SETTING
# This is an example of how to deploy an Azure Monitor Diagnostic Setting
# for a key vault with a storage account.
# ---------------------------------------------------------------------------------------------------------------------
# See test/azure/terraform_azure_monitor_example_test.go for how to write automated tests for this code.
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
    azuread = {
      version = "=0.7.0"
      source  = "hashicorp/azuread"
    }
  }
}

resource "random_string" "short" {
  length  = 3
  lower   = true
  upper   = false
  number  = false
  special = false
}

resource "random_string" "long" {
  length  = 6
  lower   = true
  upper   = false
  number  = false
  special = false
}

# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY A RESOURCE GROUP
# ---------------------------------------------------------------------------------------------------------------------

resource "azurerm_resource_group" "monitor" {
  name     = "terratest-monitor-rg-${var.postfix}"
  location = var.location
}

data "azurerm_client_config" "current" {}

# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY A STORAGE ACCOUNT
# ---------------------------------------------------------------------------------------------------------------------

resource "azurerm_storage_account" "monitor" {
  name                     = format("%s%s", "storage", random_string.long.result)
  resource_group_name      = azurerm_resource_group.monitor.name
  location                 = azurerm_resource_group.monitor.location
  account_tier             = "Standard"
  account_replication_type = "GRS"

  tags = {
    environment = "staging"
  }
}

# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY A KEY VAULT
# ---------------------------------------------------------------------------------------------------------------------

resource "azurerm_key_vault" "monitor" {
  name                        = "kv-${var.postfix}"
  location                    = azurerm_resource_group.monitor.location
  resource_group_name         = azurerm_resource_group.monitor.name
  enabled_for_disk_encryption = true
  tenant_id                   = data.azurerm_client_config.current.tenant_id
  soft_delete_enabled         = true
  purge_protection_enabled    = false

  sku_name = "standard"

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.object_id

    key_permissions = [
      "create",
      "get",
      "list",
      "delete",
    ]

    secret_permissions = [
      "set",
      "get",
      "list",
      "delete",
    ]

    certificate_permissions = [
      "create",
      "delete",
      "deleteissuers",
      "get",
      "getissuers",
      "import",
      "list",
      "listissuers",
      "managecontacts",
      "manageissuers",
      "setissuers",
      "update",
    ]
  }

  network_acls {
    default_action = "Deny"
    bypass         = "AzureServices"
  }

  tags = {
    environment = "Testing"
  }
}

# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY A DIAGNOSTIC SETTING
# https://www.terraform.io/docs/providers/azurerm/r/monitor_diagnostic_setting.html
# ---------------------------------------------------------------------------------------------------------------------

resource "azurerm_monitor_diagnostic_setting" "monitor" {
  name               = var.diagnosticSettingName
  target_resource_id = azurerm_key_vault.monitor.id
  storage_account_id = azurerm_storage_account.monitor.id

  log {
    category = "AuditEvent"
    enabled  = false

    retention_policy {
      enabled = false
    }
  }

  metric {
    category = "AllMetrics"

    retention_policy {
      enabled = false
    }
  }
}
