# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY AN AZURE CONTAINER APPS
# This is an example of how to deploy an Azure Container App and Azure Container App Job with the minimum set of options.
# ---------------------------------------------------------------------------------------------------------------------
# See test/azure/terraform_azure_container_apps_test.go for how to write automated tests for this code.
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
      version = "~> 3.103"
      source  = "hashicorp/azurerm"
    }
  }
}

# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY A RESOURCE GROUP
# ---------------------------------------------------------------------------------------------------------------------

resource "azurerm_resource_group" "aca" {
  name     = "terratest-rg-${var.postfix}"
  location = "East US"
}

# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY A AZURE APP ENVIRONMENT
# ---------------------------------------------------------------------------------------------------------------------

resource "azurerm_container_app_environment" "aca" {
  name                = "terratest-aca-env-${var.postfix}"
  location            = azurerm_resource_group.aca.location
  resource_group_name = azurerm_resource_group.aca.name
}

# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY A AZURE CONTAINER APP
# ---------------------------------------------------------------------------------------------------------------------

resource "azurerm_container_app" "aca" {
  name                         = "terratest-aca-${var.postfix}"
  resource_group_name          = azurerm_resource_group.aca.name
  container_app_environment_id = azurerm_container_app_environment.aca.id
  revision_mode                = "Single"
  template {
    container {
      name   = "terratest-aca-app-${var.postfix}"
      image  = "mcr.microsoft.com/azuredocs/containerapps-helloworld:latest"
      cpu    = "0.5"
      memory = "1.0Gi"
    }
  }
}

# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY A AZURE CONTAINER APP JOB
# ---------------------------------------------------------------------------------------------------------------------

resource "azurerm_container_app_job" "aca" {
  name                         = "terratest-aca-job-${var.postfix}"
  location                     = azurerm_resource_group.aca.location
  resource_group_name          = azurerm_resource_group.aca.name
  container_app_environment_id = azurerm_container_app_environment.aca.id
  replica_timeout_in_seconds   = 10
  template {
    container {
      name    = "terratest-aca-job-${var.postfix}"
      image   = "busybox:stable"
      command = ["echo", "Hello, World!"]
      cpu     = "0.5"
      memory  = "1.0Gi"
    }
  }
  manual_trigger_config {
    parallelism = 1
  }
}
