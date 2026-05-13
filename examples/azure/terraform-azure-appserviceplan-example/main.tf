# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY AN AZURE APP SERVICE PLAN
# This is an example of how to deploy an Azure App Service Plan
# ---------------------------------------------------------------------------------------------------------------------

# ------------------------------------------------------------------------------
# CONFIGURE OUR AZURE CONNECTION
# ------------------------------------------------------------------------------

provider "azurerm" {
  features {}
  skip_provider_registration = true
  # To understand why ^ is here, see https://github.com/hashicorp/terraform/issues/18180
}

# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY A RESOURCE GROUP
# ---------------------------------------------------------------------------------------------------------------------

resource "azurerm_resource_group" "rg" {
  name     = var.resourceGroupName
  location = var.location
}

# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY AN AZURE APP SERVICE PLAN
# ---------------------------------------------------------------------------------------------------------------------

resource "azurerm_app_service_plan" "plan" {
  count = var.instanceCount

  name                = var.appName
  resource_group_name = azurerm_resource_group.rg.name
  location            = var.location

  tags     = azurerm_resource_group.rg.tags
  kind     = var.kind
  reserved = lower(var.kind) == "linux" ? true : lower(var.kind) == "windows" || lower(var.kind) == "app" ? false : var.reserved

  dynamic "sku" {
    for_each = lower(var.kind) == "functionapp" ? ["sku"] : []
    content {
      tier = "Dynamic"
      size = "Y1"
    }
  }

  dynamic "sku" {
    for_each = lower(var.kind) != "functionapp" ? ["sku"] : []
    content {
      tier     = var.skuTier
      size     = var.skuSize
      capacity = var.skuCapacity
    }
  }
}