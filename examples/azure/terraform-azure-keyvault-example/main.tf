# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY AN AZURE KEY VAULT
# This is an example of how to deploy a Key Vault 
# ---------------------------------------------------------------------------------------------------------------------
# See test/azure/terraform_azure_keyvault_example_test.go for how to write automated tests for this code.
# ---------------------------------------------------------------------------------------------------------------------

provider "azurerm" {
  features {
    key_vault {
      purge_soft_delete_on_destroy = false
    }
  }
}

# ---------------------------------------------------------------------------------------------------------------------
# PIN TERRAFORM VERSION TO >= 0.12
# The examples have been upgraded to 0.12 syntax
# ---------------------------------------------------------------------------------------------------------------------

terraform {
  # This module is now only being tested with Terraform 0.13.x. However, to make upgrading easier, we are setting
  # 0.12.26 as the minimum version, as that version added support for required_providers with source URLs, making it
  # forwards compatible with 0.13.x code.
  required_version = ">= 0.12.26"
  required_providers {
    azurerm = {
      version = "~>3.0"
      source  = "hashicorp/azurerm"
    }
  }
}

# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY A RESOURCE GROUP
# ---------------------------------------------------------------------------------------------------------------------

resource "azurerm_resource_group" "resource_group" {
  name     = "terratest-kv-rg-${var.postfix}"
  location = var.location
}

# ---------------------------------------------------------------------------------------------------------------------
# CONFIGURE A CLIENT FOR KEY VAULT ACCESS
# ---------------------------------------------------------------------------------------------------------------------

data "azurerm_client_config" "current" {}

# ---------------------------------------------------------------------------------------------------------------------
# CONFIGURE AN ACCESS POLICY TO MANAGE THE SECRET, KEY, AND CERTIFICATE
# ---------------------------------------------------------------------------------------------------------------------

data "azurerm_key_vault_access_policy" "contributor" {
  name = "Key, Secret, & Certificate Management"
}

# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY A KEY VAULT
# ---------------------------------------------------------------------------------------------------------------------

resource "azurerm_key_vault" "key_vault" {
  name                        = "keyvault-${var.postfix}"
  location                    = azurerm_resource_group.resource_group.location
  resource_group_name         = azurerm_resource_group.resource_group.name
  enabled_for_disk_encryption = true
  tenant_id                   = data.azurerm_client_config.current.tenant_id

  soft_delete_retention_days = 7
  purge_protection_enabled   = false

  sku_name = "standard"

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.object_id

    key_permissions = [
      "Create",
      "Get",
      "List",
      "Delete",
      "Purge",
      "SetRotationPolicy",
      "GetRotationPolicy"
    ]

    secret_permissions = [
      "Set",
      "Get",
      "List",
      "Delete",
      "Purge",
    ]

    certificate_permissions = [
      "Create",
      "Delete",
      "DeleteIssuers",
      "Get",
      "GetIssuers",
      "Import",
      "List",
      "ListIssuers",
      "ManageContacts",
      "ManageIssuers",
      "SetIssuers",
      "Update",
      "Purge",
    ]
  }
}

# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY A SECRET TO THE KEY VAULT
# ---------------------------------------------------------------------------------------------------------------------

resource "azurerm_key_vault_secret" "key_vault_secret" {
  name         = "${var.secret_name}-${var.postfix}"
  value        = "mysecret"
  key_vault_id = azurerm_key_vault.key_vault.id
}

# ---------------------------------------------------------------------------------------------------------------------
#  DEPLOY A KEY TO THE KEY VAULT
# ---------------------------------------------------------------------------------------------------------------------

resource "azurerm_key_vault_key" "key_vault_key" {
  name         = "${var.key_name}-${var.postfix}"
  key_vault_id = azurerm_key_vault.key_vault.id
  key_type     = "RSA"
  key_size     = 2048

  key_opts = [
    "decrypt",
    "encrypt",
    "sign",
    "unwrapKey",
    "verify",
    "wrapKey",
  ]
}

# ---------------------------------------------------------------------------------------------------------------------
#  DEPLOY A CERTIFICATE TO THE KEY VAULT
#  The example uses a sample pfx file with plain text password to make it easier to test. However, in production modules 
#  should use a more secure mechanisms for transferring these files.
# ---------------------------------------------------------------------------------------------------------------------
resource "azurerm_key_vault_certificate" "key_vault_certificate" {
  name         = "${var.certificate_name}-${var.postfix}"
  key_vault_id = azurerm_key_vault.key_vault.id

  certificate {
    contents = filebase64("example.pfx")
    password = "password"
  }

  certificate_policy {
    issuer_parameters {
      name = "Self"
    }

    key_properties {
      exportable = true
      key_size   = 2048
      key_type   = "RSA"
      reuse_key  = false
    }

    secret_properties {
      content_type = "application/x-pkcs12"
    }
  }
}
