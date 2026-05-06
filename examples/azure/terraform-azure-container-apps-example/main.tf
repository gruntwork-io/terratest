# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY AN AZURE CONTAINER APPS
# This is an example of how to deploy an Azure Container App and Azure Container App Job with the minimum set of options.
# ---------------------------------------------------------------------------------------------------------------------
# See test/azure/terraform_azure_container_apps_test.go for how to write automated tests for this code.
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
