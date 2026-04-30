# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY AN AZURE MySQL FLEXIBLE SERVER + DATABASE
# This is an example of how to deploy an Azure MySQL Flexible Server and database.
# See test/terraform_azure_example_test.go for how to write automated tests for this code.
#
# NOTE: The legacy `azurerm_mysql_server` (MySQL Single Server) was retired by Azure on 2024-09-16.
# This example uses the modern flexible server resources, which are the supported replacement.
# ---------------------------------------------------------------------------------------------------------------------


# ---------------------------------------------------------------------------------------------------------------------
# CONFIGURE OUR AZURE CONNECTION
# ---------------------------------------------------------------------------------------------------------------------
terraform {
  required_version = ">= 1.0"

  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "~> 3.0"
    }
  }
}

provider "azurerm" {
  features {}
}

# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY A RESOURCE GROUP
# ---------------------------------------------------------------------------------------------------------------------

resource "azurerm_resource_group" "mysql_rg" {
  name     = "terratest-mysql-${var.postfix}"
  location = var.location
}

# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY AZURE MySQL FLEXIBLE SERVER
# ---------------------------------------------------------------------------------------------------------------------

# Random password is used as an example to simplify the deployment and improve the security of the database.
# This is not as a production recommendation as the password is stored in the Terraform state file.
resource "random_password" "password" {
  length           = 16
  override_special = "_%@"
  min_upper        = "1"
  min_lower        = "1"
  min_numeric      = "1"
  min_special      = "1"
}

resource "azurerm_mysql_flexible_server" "mysqlserver" {
  name                = "mysqlserver-${var.postfix}"
  location            = azurerm_resource_group.mysql_rg.location
  resource_group_name = azurerm_resource_group.mysql_rg.name

  administrator_login    = var.mysqlserver_admin_login
  administrator_password = random_password.password.result

  sku_name = var.mysqlserver_sku_name
  version  = var.mysqlserver_version

  backup_retention_days = 7

  storage {
    size_gb           = var.mysqlserver_storage_size_gb
    auto_grow_enabled = true
  }
}

# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY AZURE MySQL FLEXIBLE DATABASE
# ---------------------------------------------------------------------------------------------------------------------

resource "azurerm_mysql_flexible_database" "mysqldb" {
  name                = "mysqldb-${var.postfix}"
  resource_group_name = azurerm_resource_group.mysql_rg.name
  server_name         = azurerm_mysql_flexible_server.mysqlserver.name
  charset             = var.mysqldb_charset
  collation           = var.mysqldb_collation
}

# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY AZURE MySQL FLEXIBLE SERVER FIREWALL RULE
# Allow access from the Azure backbone (0.0.0.0) so the test runner can reach the server.
# ---------------------------------------------------------------------------------------------------------------------

resource "azurerm_mysql_flexible_server_firewall_rule" "allow_azure" {
  name                = "AllowAzureServices"
  resource_group_name = azurerm_resource_group.mysql_rg.name
  server_name         = azurerm_mysql_flexible_server.mysqlserver.name
  start_ip_address    = "0.0.0.0"
  end_ip_address      = "0.0.0.0"
}
