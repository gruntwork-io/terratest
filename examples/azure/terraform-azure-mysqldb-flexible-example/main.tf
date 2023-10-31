# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY AN AZURE Database For MySQL - Flexible Server
# This is an example of how to deploy an Azure Database Mysql - Flexible Server .
# See test/terraform_azure_mysqldb_flexible_example_test.go for how to write automated tests for this code.
# ---------------------------------------------------------------------------------------------------------------------

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
      source  = "hashicorp/azurerm"
      version = "~> 3.41.0"
    }

    random = {
      source  = "hashicorp/random"
      version = "~> 3.4.3"
    }
  }
}

# ---------------------------------------------------------------------------------------------------------------------
# CONFIGURE OUR AZURE CONNECTION
# ---------------------------------------------------------------------------------------------------------------------

provider "azurerm" {
  features {}
}

# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY A RESOURCE GROUP
# ---------------------------------------------------------------------------------------------------------------------

resource "azurerm_resource_group" "example" {
  location = var.location
  name     = "rg-flexible-${var.postfix}"
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

resource "azurerm_mysql_flexible_server" "example" {
  name                = "mysql-flexible-${var.postfix}"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  administrator_login    = var.mysql_flexible_server_administrator_login
  administrator_password = random_password.password.result

  backup_retention_days = var.mysql_flexible_server_backup_retention_days
  sku_name              = var.mysql_flexible_server_sku_name
  version               = var.mysql_version
  # Azure automatically deploy the instance in an Availability Zone.
  # By providing this, when updating your configuration, you avoid error:
  #     Error: `zone` cannot be changed independently.
  zone = var.mysql_flexible_server_zone

  storage {
    auto_grow_enabled = false
    size_gb           = var.mysql_flexible_server_storage_size_gb
  }
}

# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY AZURE MySQL DATABASE
# ---------------------------------------------------------------------------------------------------------------------

resource "azurerm_mysql_flexible_database" "example" {
  name                = "mysqldb-flexible-${var.postfix}"
  resource_group_name = azurerm_resource_group.example.name
  server_name         = azurerm_mysql_flexible_server.example.name

  charset   = var.mysql_flexible_server_db_charset
  collation = var.mysql_flexible_server_db_collation
}