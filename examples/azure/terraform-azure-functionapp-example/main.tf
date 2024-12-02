# ---------------------------------------------------------------------------------------------------------------------
# Deploy an Azure storage account, service plan, function app, and application insights
# This is an example of how to deploy an Azure function app.
# See test/terraform_azure_functionapp_example_test.go for how to write automated tests for this code.
# ---------------------------------------------------------------------------------------------------------------------


# ---------------------------------------------------------------------------------------------------------------------
# CONFIGURE OUR AZURE CONNECTION
# ---------------------------------------------------------------------------------------------------------------------

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
# DEPLOY THE RESOURCE GROUP
# ---------------------------------------------------------------------------------------------------------------------

resource "azurerm_resource_group" "app_rg" {
  name     = "terratest-functionapp-rg-${var.postfix}"
  location = var.location
}

# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY THE AZURE STORAGE ACCOUNT
# ---------------------------------------------------------------------------------------------------------------------

resource "azurerm_storage_account" "storage" {
  name                     = "storageaccount${var.postfix}"
  resource_group_name      = azurerm_resource_group.app_rg.name
  location                 = azurerm_resource_group.app_rg.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY THE APP SERVICE PLAN
# ---------------------------------------------------------------------------------------------------------------------

resource "azurerm_app_service_plan" "app_service_plan" {
  name                = "appservice-plan-${var.postfix}"
  location            = azurerm_resource_group.app_rg.location
  resource_group_name = azurerm_resource_group.app_rg.name
  kind                = "FunctionApp"

  sku {
    tier = "Standard"
    size = "S1"
  }
}

# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY THE APPLICATION INSIGHTS
# ---------------------------------------------------------------------------------------------------------------------

resource "azurerm_application_insights" "application_insights" {
  name                = "appinsights-${var.postfix}"
  location            = azurerm_resource_group.app_rg.location
  resource_group_name = azurerm_resource_group.app_rg.name
  application_type    = "web"
}

# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY THE AZURE FUNCTION APP
# ---------------------------------------------------------------------------------------------------------------------

resource "azurerm_function_app" "function_app" {
  name                       = "functionapp-${var.postfix}"
  location                   = azurerm_resource_group.app_rg.location
  resource_group_name        = azurerm_resource_group.app_rg.name
  app_service_plan_id        = azurerm_app_service_plan.app_service_plan.id
  storage_account_name       = azurerm_storage_account.storage.name
  storage_account_access_key = azurerm_storage_account.storage.primary_access_key


  app_settings = {
    "APPINSIGHTS_INSTRUMENTATIONKEY"        = azurerm_application_insights.application_insights.instrumentation_key
    "APPLICATIONINSIGHTS_CONNECTION_STRING" = "InstrumentationKey=${azurerm_application_insights.application_insights.instrumentation_key}"
  }
}
