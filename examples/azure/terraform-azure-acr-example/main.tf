# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY AN AZURE CONTAINER REGISTRY
# This is an example of how to deploy an Azure Container Registry
# See test/terraform_azure_acr_example_test.go for how to write automated tests for this code.
# ---------------------------------------------------------------------------------------------------------------------

# ------------------------------------------------------------------------------
# CONFIGURE OUR AZURE CONNECTION
# ------------------------------------------------------------------------------

terraform {
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

resource "azurerm_resource_group" "rg" {
  name     = "terratest-acr-rg-${var.postfix}"
  location = var.location
}

# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY AN AZURE CONTAINER REGISTRY
# ---------------------------------------------------------------------------------------------------------------------

resource "azurerm_container_registry" "acr" {
  name                = "acr${var.postfix}"
  location            = azurerm_resource_group.rg.location
  resource_group_name = azurerm_resource_group.rg.name

  sku           = var.sku
  admin_enabled = true

  tags = {
    Environment = "Development"
  }
}
