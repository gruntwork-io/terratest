# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY AN AZURE NETWORK
# This is an example of how to deploy frequent Azure Networking Resources. Note this network doesn't actually do
# anything and is only created for the example to test their commonly needed and integrated properties.
# ---------------------------------------------------------------------------------------------------------------------
# See test/azure/terraform_azure_network_example_test.go for how to write automated tests for this code.
# ---------------------------------------------------------------------------------------------------------------------

terraform {
  required_version = ">= 1.0"

  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "~> 4.0"
    }
  }
}

provider "azurerm" {
  features {}
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
  sku                     = "Standard"
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

