# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY AN PostgreSQL Database
# This is an example of how to deploy an Azure PostgreSQL Flexible Server.
# Azure retired PostgreSQL Single Server (azurerm_postgresql_server) on 2025-03-28; this example uses the
# replacement Flexible Server resources.
# See test/terraform_azure_postgresql_example_test.go for how to write automated tests for this code.
# ---------------------------------------------------------------------------------------------------------------------


# ---------------------------------------------------------------------------------------------------------------------
# CONFIGURE OUR AZURE CONNECTION
# ---------------------------------------------------------------------------------------------------------------------
terraform {
  required_version = ">= 1.0"

  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "~> 4.0"
    }
    random = {
      source  = "hashicorp/random"
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
resource "azurerm_resource_group" "rg" {
  name     = "${var.resource_group_name}-${var.postfix}"
  location = var.location
}

# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY AZURE PostgreSQL FLEXIBLE SERVER
# ---------------------------------------------------------------------------------------------------------------------
resource "azurerm_postgresql_flexible_server" "postgresqlserver" {
  name                = "postgresqlserver-${var.postfix}"
  location            = azurerm_resource_group.rg.location
  resource_group_name = azurerm_resource_group.rg.name

  administrator_login    = "pgsqladmin"
  administrator_password = random_password.password.result

  sku_name = var.sku_name
  version  = var.postgresql_version

  storage_mb            = 32768
  backup_retention_days = 7

  zone = "1"

  # Public access only; no high-availability for cost reasons in this example.
  public_network_access_enabled = true
}

# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY AZURE PostgreSQL FLEXIBLE SERVER DATABASE
# ---------------------------------------------------------------------------------------------------------------------
resource "azurerm_postgresql_flexible_server_database" "postgresqldb" {
  name      = "postgresqldb"
  server_id = azurerm_postgresql_flexible_server.postgresqlserver.id
  charset   = "UTF8"
  collation = "en_US.utf8"
}

# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY A FIREWALL RULE
# Allow Azure-internal services to reach the server. Adjust ranges to taste.
# ---------------------------------------------------------------------------------------------------------------------
resource "azurerm_postgresql_flexible_server_firewall_rule" "allow_azure" {
  name             = "AllowAllAzureIPs"
  server_id        = azurerm_postgresql_flexible_server.postgresqlserver.id
  start_ip_address = "0.0.0.0"
  end_ip_address   = "0.0.0.0"
}

# ---------------------------------------------------------------------------------------------------------------------
# Use a random password generator
# ---------------------------------------------------------------------------------------------------------------------
resource "random_password" "password" {
  length  = 20
  special = true
  upper   = true
  lower   = true
  numeric = true
}
