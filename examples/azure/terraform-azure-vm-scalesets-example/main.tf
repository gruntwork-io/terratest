# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY AN AZURE VIRTUAL MACHINE SCALE SET
# This is an example of how to deploy an Azure Virtual Machine Scale set.
# ---------------------------------------------------------------------------------------------------------------------
# See test/azure/terraform_azure_vm_scalesets_example_test.go for how to write automated tests for this code.
# ---------------------------------------------------------------------------------------------------------------------

provider "azurerm" {
  version = "~> 2.50"
  features {}
}

# ---------------------------------------------------------------------------------------------------------------------
# PIN TERRAFORM VERSION
# ---------------------------------------------------------------------------------------------------------------------

terraform {
  # This module is now only being tested with Terraform 0.13.x. However, to make upgrading easier, we are setting
  # 0.12.26 as the minimum version, as that version added support for required_providers with source URLs, making it
  # forwards compatible with 0.13.x code.
  required_version = ">= 0.12.26"
}

# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY A RESOURCE GROUP
# ---------------------------------------------------------------------------------------------------------------------

resource "azurerm_resource_group" "vmss_rg" {
  name     = "terratest-vmss-rg-${var.postfix}"
  location = var.location
}

# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY NETWORK RESOURCES
# This network includes a public address for integration test demonstration purposes
# ---------------------------------------------------------------------------------------------------------------------

resource "azurerm_virtual_network" "vnet" {
  name                = "vnet-${var.postfix}"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.vmss_rg.location
  resource_group_name = azurerm_resource_group.vmss_rg.name
}

resource "azurerm_subnet" "subnet" {
  name                 = "subnet-${var.postfix}"
  resource_group_name  = azurerm_resource_group.vmss_rg.name
  virtual_network_name = azurerm_virtual_network.vnet.name
  address_prefixes     = [var.subnet_prefix]
}

resource "azurerm_public_ip" "pip" {
  name                    = "pip-${var.postfix}"
  location                = azurerm_resource_group.vmss_rg.location
  resource_group_name     = azurerm_resource_group.vmss_rg.name
  allocation_method       = "Static"
  ip_version              = "IPv4"
  sku                     = "Basic"
  idle_timeout_in_minutes = "4"
}

resource "azurerm_lb" "public" {
  name                = "lb-public-${var.postfix}"
  location            = azurerm_resource_group.vmss_rg.location
  resource_group_name = azurerm_resource_group.vmss_rg.name
  sku                 = "Basic"

  frontend_ip_configuration {
    name                 = "config-public"
    public_ip_address_id = azurerm_public_ip.pip.id
  }
}

resource "azurerm_lb_backend_address_pool" "backendpool" {
  loadbalancer_id = azurerm_lb.public.id
  name            = "VMScaleSetBackendPool"
}

# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY VIRTUAL MACHINE SCALE SET
# The VM in a scale set does not actually do anything and is the smallest size VM available with a Windows image
# ---------------------------------------------------------------------------------------------------------------------

resource "azurerm_virtual_machine_scale_set" "scaleset" {
  name                = "vmss-${var.postfix}"
  location            = azurerm_resource_group.vmss_rg.location
  resource_group_name = azurerm_resource_group.vmss_rg.name
  upgrade_policy_mode = "Manual"
  license_type        = var.vm_license_type

  sku {
    name     = var.vm_size
    tier     = var.vm_tier
    capacity = var.vmss_capacity
  }

  storage_profile_image_reference {
    publisher = var.vm_image_publisher
    offer     = var.vm_image_offer
    sku       = var.vm_image_sku
    version   = var.vm_image_version
  }

  storage_profile_os_disk {
    name              = ""
    caching           = "ReadWrite"
    create_option     = "FromImage"
    managed_disk_type = var.disk_type
  }

  os_profile {
    computer_name_prefix = "vmss"
    admin_username       = var.user_name
    admin_password       = random_password.rand.result
  }

  os_profile_windows_config {
    provision_vm_agent = true
  }

  network_profile {
    name    = "networkprofile"
    primary = true

    ip_configuration {
      name                                   = "terratestconfiguration1"
      primary                                = true
      subnet_id                              = azurerm_subnet.subnet.id
      load_balancer_backend_address_pool_ids = [azurerm_lb_backend_address_pool.backendpool.id]
    }
  }

  tags = {
    "Version"     = "0.0.1"
    "Environment" = "dev"
  }

  depends_on = [random_password.rand]
}

# Random password is used as an example to simplify the deployment and improve the security of the remote VM.
# This is not as a production recommendation as the password is stored in the Terraform state file.
resource "random_password" "rand" {
  length           = 16
  override_special = "-_%@"
  min_upper        = "1"
  min_lower        = "1"
  min_numeric      = "1"
  min_special      = "1"
}
