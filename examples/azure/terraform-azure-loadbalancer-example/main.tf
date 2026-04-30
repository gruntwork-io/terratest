# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY AN AZURE AVAILABILITY SET
# This is an example of how to deploy an Azure Availability Set with a Virtual Machine in the availability set 
# and the minimum network resources for the VM.
# ---------------------------------------------------------------------------------------------------------------------
# See test/azure/terraform_azure_loadbalancer_example_test.go for how to write automated tests for this code.
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

resource "azurerm_resource_group" "lb_rg" {
  name     = "terratest-lb-rg-${var.postfix}"
  location = var.location
}

# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY VIRTUAL NETWORK
# ---------------------------------------------------------------------------------------------------------------------

resource "azurerm_virtual_network" "vnet" {
  name                = "vnet-${var.postfix}"
  location            = azurerm_resource_group.lb_rg.location
  resource_group_name = azurerm_resource_group.lb_rg.name
  address_space       = ["10.200.0.0/21"]
}

resource "azurerm_subnet" "subnet" {
  name                 = "subnet-${var.postfix}"
  resource_group_name  = azurerm_resource_group.lb_rg.name
  virtual_network_name = azurerm_virtual_network.vnet.name
  address_prefixes     = ["10.200.2.0/25"]
}

# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY LOAD BALANCER WITH PUBLIC IP 
# ---------------------------------------------------------------------------------------------------------------------

resource "azurerm_public_ip" "pip" {
  name                    = "pip-${var.postfix}"
  location                = azurerm_resource_group.lb_rg.location
  resource_group_name     = azurerm_resource_group.lb_rg.name
  allocation_method       = "Static"
  ip_version              = "IPv4"
  sku                     = "Standard"
  idle_timeout_in_minutes = "4"
}

resource "azurerm_lb" "public" {
  name                = "lb-public-${var.postfix}"
  location            = azurerm_resource_group.lb_rg.location
  resource_group_name = azurerm_resource_group.lb_rg.name
  sku                 = "Standard"

  frontend_ip_configuration {
    name                 = "config-public"
    public_ip_address_id = azurerm_public_ip.pip.id
  }
}

# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY LOAD BALANCER WITH PRIVATE IPs 
# ---------------------------------------------------------------------------------------------------------------------

resource "azurerm_lb" "private" {
  name                = "lb-private-${var.postfix}"
  location            = azurerm_resource_group.lb_rg.location
  resource_group_name = azurerm_resource_group.lb_rg.name
  sku                 = "Standard"

  frontend_ip_configuration {
    name                          = "config-private-static"
    subnet_id                     = azurerm_subnet.subnet.id
    private_ip_address            = var.lb_private_ip
    private_ip_address_allocation = "Static"
  }

  frontend_ip_configuration {
    name                          = "config-private-dynamic"
    subnet_id                     = azurerm_subnet.subnet.id
    private_ip_address_allocation = "Dynamic"
  }
}

# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY LOAD BALANCER WITH NO FRONTEND CONFIGURATION
# ---------------------------------------------------------------------------------------------------------------------

resource "azurerm_lb" "default" {
  name                = "lb-no-frontend-${var.postfix}"
  location            = azurerm_resource_group.lb_rg.location
  resource_group_name = azurerm_resource_group.lb_rg.name
  sku                 = "Standard"
}
