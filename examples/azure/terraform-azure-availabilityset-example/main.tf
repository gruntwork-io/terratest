# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY AN AZURE AVAILABILITY SET
# This is an example of how to deploy an Azure Availability Set with a Virtual Machine in the availability set 
# and the minimum network resources for the VM.
# ---------------------------------------------------------------------------------------------------------------------
# See test/azure/terraform_azure_availabilityset_example_test.go for how to write automated tests for this code.
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
      version = "~> 2.50"
      source  = "hashicorp/azurerm"
    }
  }
}

# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY A RESOURCE GROUP
# ---------------------------------------------------------------------------------------------------------------------

resource "azurerm_resource_group" "avs" {
  name     = "terratest-avs-rg-${var.postfix}"
  location = var.location
}

# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY THE AVAILABILITY SET
# ---------------------------------------------------------------------------------------------------------------------

resource "azurerm_availability_set" "avs" {
  name                        = "avs-${var.postfix}"
  location                    = azurerm_resource_group.avs.location
  resource_group_name         = azurerm_resource_group.avs.name
  platform_fault_domain_count = var.avs_fault_domain_count
  managed                     = true
}

# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY MINIMAL NETWORK RESOURCES FOR VM
# ---------------------------------------------------------------------------------------------------------------------

resource "azurerm_virtual_network" "avs" {
  name                = "vnet-${var.postfix}"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.avs.location
  resource_group_name = azurerm_resource_group.avs.name
}

resource "azurerm_subnet" "avs" {
  name                 = "subnet-${var.postfix}"
  resource_group_name  = azurerm_resource_group.avs.name
  virtual_network_name = azurerm_virtual_network.avs.name
  address_prefixes     = ["10.0.17.0/24"]
}

resource "azurerm_network_interface" "avs" {
  name                = "nic-${var.postfix}"
  location            = azurerm_resource_group.avs.location
  resource_group_name = azurerm_resource_group.avs.name

  ip_configuration {
    name                          = "config-${var.postfix}-01"
    subnet_id                     = azurerm_subnet.avs.id
    private_ip_address_allocation = "Dynamic"
  }
}

# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY VIRTUAL MACHINE
# This VM does not actually do anything and is the smallest size VM available with an Ubuntu image
# ---------------------------------------------------------------------------------------------------------------------

resource "azurerm_virtual_machine" "avs" {
  name                             = "vm-${var.postfix}"
  location                         = azurerm_resource_group.avs.location
  resource_group_name              = azurerm_resource_group.avs.name
  network_interface_ids            = [azurerm_network_interface.avs.id]
  availability_set_id              = azurerm_availability_set.avs.id
  vm_size                          = "Standard_B1ls"
  delete_os_disk_on_termination    = true
  delete_data_disks_on_termination = true

  storage_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "18.04-LTS"
    version   = "latest"
  }

  storage_os_disk {
    name              = "osdisk-${var.postfix}"
    caching           = "None"
    create_option     = "FromImage"
    managed_disk_type = "Standard_LRS"
  }

  os_profile {
    computer_name  = "vm-${var.postfix}"
    admin_username = "testadmin"
    admin_password = random_password.avs.result
  }

  os_profile_linux_config {
    disable_password_authentication = false
  }

  depends_on = [random_password.avs]
}

resource "random_password" "avs" {
  length           = 16
  override_special = "-_%@"
  min_upper        = "1"
  min_lower        = "1"
  min_numeric      = "1"
  min_special      = "1"
}
