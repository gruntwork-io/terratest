# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY AN AZURE VM ALONG WITH AN EXAMPLE NETWORK SECURITY GROUP (NSG)
# This is an example of how to deploy an NSG along with the minimum networking resources
# to support a basic virtual machine.
# ---------------------------------------------------------------------------------------------------------------------
# See test/azure/terraform_azure_nsg_example_test.go for how to write automated tests for this code.
# ---------------------------------------------------------------------------------------------------------------------

provider "azurerm" {
  features {}
}

# ---------------------------------------------------------------------------------------------------------------------
# PIN TERRAFORM VERSION TO >= 0.12
# The examples have been upgraded to 0.12 syntax
# ---------------------------------------------------------------------------------------------------------------------

terraform {
  required_version = ">= 0.12"
  required_providers {
    azurerm = {
      version = "~> 2.50"
      source  = "hashicorp/azurerm"
    }
  }
}

# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY A RESOURCE GROUP
# See test/terraform_azure_nsg_example_test.go for how to write automated tests for this code.
# ---------------------------------------------------------------------------------------------------------------------

resource "azurerm_resource_group" "nsg_rg" {
  name     = "${var.resource_group_name}-${var.postfix}"
  location = var.location
}

# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY VIRTUAL NETWORK RESOURCES
# ---------------------------------------------------------------------------------------------------------------------

resource "azurerm_virtual_network" "vnet" {
  name                = "${var.vnet_name}-${var.postfix}"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.nsg_rg.location
  resource_group_name = azurerm_resource_group.nsg_rg.name
}

resource "azurerm_subnet" "internal" {
  name                 = "${var.subnet_name}-${var.postfix}"
  resource_group_name  = azurerm_resource_group.nsg_rg.name
  virtual_network_name = azurerm_virtual_network.vnet.name
  address_prefixes     = ["10.0.17.0/24"]
}

resource "azurerm_network_interface" "main" {
  name                = "${var.vm_nic_name}-${var.postfix}"
  location            = azurerm_resource_group.nsg_rg.location
  resource_group_name = azurerm_resource_group.nsg_rg.name

  ip_configuration {
    name                          = "${var.vm_nic_ip_config_name}-${var.postfix}"
    subnet_id                     = azurerm_subnet.internal.id
    private_ip_address_allocation = "Dynamic"
  }
}

resource "azurerm_network_security_group" "nsg_example" {
  name                = "${var.nsg_name}-${var.postfix}"
  location            = azurerm_resource_group.nsg_rg.location
  resource_group_name = azurerm_resource_group.nsg_rg.name
}

resource "azurerm_network_interface_security_group_association" "main" {
  network_interface_id      = azurerm_network_interface.main.id
  network_security_group_id = azurerm_network_security_group.nsg_example.id
}

resource "azurerm_network_security_rule" "allow_ssh" {
  name                        = "${var.nsg_ssh_rule_name}-${var.postfix}"
  description                 = "${var.nsg_ssh_rule_name}-${var.postfix}"
  priority                    = 100
  direction                   = "Inbound"
  access                      = "Allow"
  protocol                    = "Tcp"
  source_port_range           = "*"
  destination_port_range      = 22
  source_address_prefix       = "*"
  destination_address_prefix  = "*"
  resource_group_name         = azurerm_resource_group.nsg_rg.name
  network_security_group_name = azurerm_network_security_group.nsg_example.name
}

resource "azurerm_network_security_rule" "block_http" {
  name                        = "${var.nsg_http_rule_name}-${var.postfix}"
  description                 = "${var.nsg_http_rule_name}-${var.postfix}"
  priority                    = 200
  direction                   = "Inbound"
  access                      = "Deny"
  protocol                    = "Tcp"
  source_port_range           = "*"
  destination_port_range      = 80
  source_address_prefix       = "*"
  destination_address_prefix  = "*"
  resource_group_name         = azurerm_resource_group.nsg_rg.name
  network_security_group_name = azurerm_network_security_group.nsg_example.name
}

# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY A VIRTUAL MACHINE RUNNING UBUNTU
# This VM does not actually do anything and is the smallest size VM available with an Ubuntu image
# ---------------------------------------------------------------------------------------------------------------------

resource "azurerm_virtual_machine" "vm_example" {
  name                             = "${var.vm_name}-${var.postfix}"
  location                         = azurerm_resource_group.nsg_rg.location
  resource_group_name              = azurerm_resource_group.nsg_rg.name
  network_interface_ids            = [azurerm_network_interface.main.id]
  vm_size                          = var.vm_size
  delete_os_disk_on_termination    = true
  delete_data_disks_on_termination = true

  storage_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }

  storage_os_disk {
    name              = "${var.os_disk_name}-${var.postfix}"
    caching           = "ReadWrite"
    create_option     = "FromImage"
    managed_disk_type = "Standard_LRS"
  }

  os_profile {
    computer_name  = var.hostname
    admin_username = var.username
    admin_password = random_password.nsg.result
  }

  os_profile_linux_config {
    disable_password_authentication = false
  }

  # Correctly setup the dependencies to make sure resources are correctly destroyed.
  depends_on = [
    azurerm_network_interface_security_group_association.main
  ]
}

resource "random_password" "nsg" {
  length           = 16
  override_special = "-_%@"
  min_upper        = "1"
  min_lower        = "1"
  min_numeric      = "1"
  min_special      = "1"
}

