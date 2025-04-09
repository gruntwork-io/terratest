# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY AN AZURE SQL Database
# This is an example of how to deploy an Azure sql database.
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
      version = "~>2.29"
      source  = "hashicorp/azurerm"
    }
  }
}
# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY A RESOURCE GROUP
# ---------------------------------------------------------------------------------------------------------------------

resource "azurerm_resource_group" "sql_rg" {
  name     = "terratest-rg-${var.postfix}"
  location = var.location
}

# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY AZURE SQL SERVER
# ---------------------------------------------------------------------------------------------------------------------

resource "random_password" "password" {
  length           = 16
  override_special = "_%@"
  min_upper        = "1"
  min_lower        = "1"
  min_numeric      = "1"
  min_special      = "1"
}

resource "azurerm_sql_server" "sqlserver" {
  name                         = "mssqlserver-${var.postfix}"
  resource_group_name          = azurerm_resource_group.sql_rg.name
  location                     = azurerm_resource_group.sql_rg.location
  version                      = "12.0"
  administrator_login          = var.sqlserver_admin_login
  administrator_login_password = random_password.password.result

  tags = {
    environment = var.tags
  }
}

# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY AZURE SQL DATA BASE
# ---------------------------------------------------------------------------------------------------------------------

resource "azurerm_sql_database" "sqldb" {
  name                = "sqldb-${var.postfix}"
  resource_group_name = azurerm_resource_group.sql_rg.name
  location            = azurerm_resource_group.sql_rg.location
  server_name         = azurerm_sql_server.sqlserver.name
  tags = {
    environment = var.tags
  }
}
