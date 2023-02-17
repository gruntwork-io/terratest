# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY AN PostgreSQL Database
# This is an example of how to deploy an Azure PostgreSQL database.
# See test/terraform_azure_example_test.go for how to write automated tests for this code.
# ---------------------------------------------------------------------------------------------------------------------


# ---------------------------------------------------------------------------------------------------------------------
# CONFIGURE OUR AZURE CONNECTION
# ---------------------------------------------------------------------------------------------------------------------
provider "azurerm" {
  features {}
}

# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY A RESOURCE GROUP
# ---------------------------------------------------------------------------------------------------------------------
resource "azurerm_resource_group" "rg" {
  name     = "${var.resource_group_name}-${var.postfix}"
  location = var.location
}

# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY AZURE PostgreSQL SERVER
# ---------------------------------------------------------------------------------------------------------------------
resource "azurerm_postgresql_server" "postgresqlserver" {
  name                = "postgresqlserver-${var.postfix}"
  location            = azurerm_resource_group.rg.location
  resource_group_name = azurerm_resource_group.rg.name

  sku_name = "B_Gen5_2"

  storage_mb                   = 5120
  backup_retention_days        = 7
  geo_redundant_backup_enabled = false
  auto_grow_enabled            = true

  administrator_login          = "pgsqladmin"
  administrator_login_password = random_password.password.result
  version                      = "11"
  ssl_enforcement_enabled      = true
}

# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY AZURE PostgreSQL Database
# ---------------------------------------------------------------------------------------------------------------------
resource "azurerm_postgresql_database" "postgresqldb" {
  name                = "postgresqldb"
  resource_group_name = azurerm_resource_group.rg.name
  server_name         = azurerm_postgresql_server.postgresqlserver.name
  charset             = "UTF8"
  collation           = "English_United States.1252"
}

# ---------------------------------------------------------------------------------------------------------------------
# Use a random password geneerator
# ---------------------------------------------------------------------------------------------------------------------
resource "random_password" "password" {
  length  = 20
  special = true
  upper   = true
  lower   = true
  number  = true
}
