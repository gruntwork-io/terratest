# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY AN AZURE Synapse Analytics
# This is an example of how to deploy an AZURE Synapse Analytics
# See test/terraform_azure_example_test.go for how to write automated tests for this code.
# ---------------------------------------------------------------------------------------------------------------------


# ---------------------------------------------------------------------------------------------------------------------
# CONFIGURE OUR AZURE CONNECTION
# ---------------------------------------------------------------------------------------------------------------------

provider "azurerm" {
  features {}
}

terraform {
  required_providers {
    azurerm = {
      version = "~>2.93.0"
      source  = "hashicorp/azurerm"
    }
  }
}
# ---------------------------------------------------------------------------------------------------------------------
# CREATE RANDOM PASSWORD
# ---------------------------------------------------------------------------------------------------------------------

# Random password is used as an example to simplify the deployment and improve the security of the database.
# This is not as a production recommendation as the password is stored in the Terraform state file.
resource "random_password" "password" {
  length           = 16
  override_special = "-_%@"
  min_upper        = "1"
  min_lower        = "1"
  min_numeric      = "1"
  min_special      = "1"
}

# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY A RESOURCE GROUP
# ---------------------------------------------------------------------------------------------------------------------

resource "azurerm_resource_group" "synapse_rg" {
  name     = "terratest-synapse-${var.postfix}"
  location = var.location
}

# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY A STORAGE ACCOUNT
# ---------------------------------------------------------------------------------------------------------------------

resource "azurerm_storage_account" "storage_account" {
  name                     = "storage${var.postfix}"
  resource_group_name      = azurerm_resource_group.synapse_rg.name
  location                 = azurerm_resource_group.synapse_rg.location
  account_kind             = var.storage_account_kind
  account_tier             = var.storage_account_tier
  account_replication_type = var.storage_account_replication_type
}


# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY A DATA LAKE GEN2
# ---------------------------------------------------------------------------------------------------------------------

resource "azurerm_storage_data_lake_gen2_filesystem" "dl_gen2" {
  name               = "dlgen2-${var.postfix}"
  storage_account_id = azurerm_storage_account.storage_account.id
}

# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY A SYNAPSE WORKSPACE
# ---------------------------------------------------------------------------------------------------------------------

resource "azurerm_synapse_workspace" "synapse_workspace" {
  name                                 = "mysynapse${var.postfix}"
  resource_group_name                  = azurerm_resource_group.synapse_rg.name
  location                             = azurerm_resource_group.synapse_rg.location
  storage_data_lake_gen2_filesystem_id = azurerm_storage_data_lake_gen2_filesystem.dl_gen2.id
  sql_administrator_login              = var.synapse_sql_user
  sql_administrator_login_password     = random_password.password.result
  managed_virtual_network_enabled      = true
}

# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY A SYNAPSE SQL POOL
# ---------------------------------------------------------------------------------------------------------------------

resource "azurerm_synapse_sql_pool" "synapse_pool" {
  name                 = "sqlpool${var.postfix}"
  synapse_workspace_id = azurerm_synapse_workspace.synapse_workspace.id
  sku_name             = var.synapse_sqlpool_sku_name
  create_mode          = var.synapse_sqlpool_create_mode
}