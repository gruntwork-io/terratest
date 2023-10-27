# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY AN AZURE NETWORK
# This is an example of how to deploy frequent Azure Networking Resources. Note this network doesn't actually do
# anything and is only created for the example to test their commonly needed and integrated properties.
# ---------------------------------------------------------------------------------------------------------------------
# See test/azure/terraform_azure_network_example_test.go for how to write automated tests for this code.
# ---------------------------------------------------------------------------------------------------------------------

provider "azurerm" {
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
  required_providers {
    azurerm = {
      version = "~>2.20"
      source  = "hashicorp/azurerm"
    }
  }
}

# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY A RESOURCE GROUP
# ---------------------------------------------------------------------------------------------------------------------

resource "azurerm_resource_group" "net" {
  name     = "terratest-network-rg-${var.postfix}"
  location = var.location
}

# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY VIRTUAL NETWORK
# ---------------------------------------------------------------------------------------------------------------------

resource "azurerm_virtual_network" "net" {
  name                = "vnet-${var.postfix}"
  location            = azurerm_resource_group.net.location
  resource_group_name = azurerm_resource_group.net.name
  address_space       = ["10.0.0.0/16"]
  dns_servers         = [var.dns_ip_01, var.dns_ip_02]
}

resource "azurerm_subnet" "net" {
  name                 = "subnet-${var.postfix}"
  resource_group_name  = azurerm_resource_group.net.name
  virtual_network_name = azurerm_virtual_network.net.name
  address_prefixes     = [var.subnet_prefix]
}

# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY PRIVATE NETWORK INTERFACE
# ---------------------------------------------------------------------------------------------------------------------

resource "azurerm_network_interface" "net01" {
  name                = "nic-private-${var.postfix}"
  location            = azurerm_resource_group.net.location
  resource_group_name = azurerm_resource_group.net.name

  ip_configuration {
    name                          = "terratestconfiguration1"
    subnet_id                     = azurerm_subnet.net.id
    private_ip_address_allocation = "Static"
    private_ip_address            = var.private_ip
  }
}

# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY PUBLIC ADDRESS AND NETWORK INTERFACE
# ---------------------------------------------------------------------------------------------------------------------

resource "azurerm_public_ip" "net" {
  name                    = "pip-${var.postfix}"
  resource_group_name     = azurerm_resource_group.net.name
  location                = azurerm_resource_group.net.location
  allocation_method       = "Static"
  ip_version              = "IPv4"
  sku                     = "Basic"
  idle_timeout_in_minutes = "4"
  domain_name_label       = var.domain_name_label
}

resource "azurerm_network_interface" "net02" {
  name                = "nic-public-${var.postfix}"
  location            = azurerm_resource_group.net.location
  resource_group_name = azurerm_resource_group.net.name

  ip_configuration {
    name                          = "terratestconfiguration1"
    subnet_id                     = azurerm_subnet.net.id
    private_ip_address_allocation = "Dynamic"
    public_ip_address_id          = azurerm_public_ip.net.id
  }
}

