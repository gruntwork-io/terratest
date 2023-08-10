# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY AN AZURE MySQL Database
# This is an example of how to deploy an Azure Mysql database.
# See test/terraform_azure_example_test.go for how to write automated tests for this code.
# ---------------------------------------------------------------------------------------------------------------------


# ---------------------------------------------------------------------------------------------------------------------
# CONFIGURE OUR AZURE CONNECTION
# ---------------------------------------------------------------------------------------------------------------------
terraform {
  # This module is now only being tested with Terraform 0.13.x. However, to make upgrading easier, we are setting
  # 0.12.26 as the minimum version, as that version added support for required_providers with source URLs, making it
  # forwards compatible with 0.13.x code.
  required_providers {
    azurerm = {
      version = "~>2.29.0"
      source  = "hashicorp/azurerm"
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
# DEPLOY AZURE MySQL SERVER
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

resource "azurerm_mysql_server" "mysqlserver" {
  name                = "mysqlserver-${var.postfix}"
  location            = azurerm_resource_group.mysql_rg.location
  resource_group_name = azurerm_resource_group.mysql_rg.name

  administrator_login          = var.mysqlserver_admin_login
  administrator_login_password = random_password.password.result

  sku_name   = var.mysqlserver_sku_name
  storage_mb = var.mysqlserver_storage_mb
  version    = "5.7"

  auto_grow_enabled                 = true
  geo_redundant_backup_enabled      = false
  infrastructure_encryption_enabled = true
  backup_retention_days             = 7
  public_network_access_enabled     = false
  ssl_enforcement_enabled           = true
  ssl_minimal_tls_version_enforced  = "TLS1_2"
}

# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY AZURE MySQL DATABASE
# ---------------------------------------------------------------------------------------------------------------------

resource "azurerm_mysql_database" "mysqldb" {
  name                = "mysqldb-${var.postfix}"
  resource_group_name = azurerm_resource_group.mysql_rg.name
  server_name         = azurerm_mysql_server.mysqlserver.name
  charset             = var.mysqldb_charset
  collation           = var.mysqldb_collation
}
