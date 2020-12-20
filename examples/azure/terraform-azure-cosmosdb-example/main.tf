# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY AN AZURE VIRTUAL MACHINE
# This is an example of how to deploy an Azure Virtual Machine with the minimum network resources.
# ---------------------------------------------------------------------------------------------------------------------
# See test/azure/terraform_azure_example_test.go for how to write automated tests for this code.
# ---------------------------------------------------------------------------------------------------------------------

provider "azurerm" {
  version = "~> 2.29"
  features {}
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
}

# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY A RESOURCE GROUP
# ---------------------------------------------------------------------------------------------------------------------

resource "azurerm_resource_group" "rg" {
  name     = "terratest-cosmos-rg-${var.postfix}"
  location = var.location
}

# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY A COSMOSDB ACCOUNT
# ---------------------------------------------------------------------------------------------------------------------

resource "azurerm_cosmosdb_account" "test" {
  name                = "terratest-${var.postfix}"
  location            = azurerm_resource_group.rg.location
  resource_group_name = azurerm_resource_group.rg.name
  offer_type          = "Standard"
  kind                = "GlobalDocumentDB"

  consistency_policy {
    consistency_level       = "Session"
    max_interval_in_seconds = 5
    max_staleness_prefix    = 100
  }

  geo_location {
    location          = azurerm_resource_group.rg.location
    failover_priority = 0
  }
}

resource "azurerm_cosmosdb_sql_database" "testdb" {
  name                = "testdb"
  throughput          = var.throughput
  resource_group_name = azurerm_resource_group.rg.name
  account_name        = azurerm_cosmosdb_account.test.name
}

resource "azurerm_cosmosdb_sql_container" "container1" {
  name                = "test-container-1"
  throughput          = var.throughput
  partition_key_path  = "/key1"
  resource_group_name = azurerm_cosmosdb_account.test.resource_group_name
  account_name        = azurerm_cosmosdb_account.test.name
  database_name       = azurerm_cosmosdb_sql_database.testdb.name
}

resource "azurerm_cosmosdb_sql_container" "container2" {
  name                = "test-container-2"
  partition_key_path  = "/key2"
  resource_group_name = azurerm_cosmosdb_account.test.resource_group_name
  account_name        = azurerm_cosmosdb_account.test.name
  database_name       = azurerm_cosmosdb_sql_database.testdb.name
}

resource "azurerm_cosmosdb_sql_container" "container3" {
  name                = "test-container-3"
  partition_key_path  = "/key3"
  resource_group_name = azurerm_cosmosdb_account.test.resource_group_name
  account_name        = azurerm_cosmosdb_account.test.name
  database_name       = azurerm_cosmosdb_sql_database.testdb.name
}
