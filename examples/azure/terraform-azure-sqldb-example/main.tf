# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY AN AZURE SQL Database
# This is an example of how to deploy an Azure sql database.
# ---------------------------------------------------------------------------------------------------------------------


# ------------------------------------------------------------------------------
# CONFIGURE OUR AZURE CONNECTION
# ------------------------------------------------------------------------------

provider "azurerm" {
  version = "~>2.8.0"
  features {}
}

# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY A RESOURCE GROUP
# See test/terraform_azure_example_test.go for how to write automated tests for this code.
# ---------------------------------------------------------------------------------------------------------------------

resource "azurerm_resource_group" "sql" {
  name     = "terratest-rg-${var.postfix}"
  location = var.location
}

# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY AZURE SQL SERVER
# ---------------------------------------------------------------------------------------------------------------------

resource "random_password" "password" {
  length = 16
  special = true
  override_special = "_%@"
}

resource "azurerm_sql_server" "sqlserver" {
  name                         = "mssqlserver-${var.postfix}"
  resource_group_name          = azurerm_resource_group.sql.name
  location                     = azurerm_resource_group.sql.location
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
  resource_group_name = azurerm_resource_group.sql.name
  location            = azurerm_resource_group.sql.location
  server_name         = azurerm_sql_server.sqlserver.name
  tags = {
    environment = var.tags
  }
}