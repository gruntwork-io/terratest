# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY AN ACTION GROUP
# This is an example of how to deploy an Azure Action Group to be used for Azure Alerts
# ---------------------------------------------------------------------------------------------------------------------

# ---------------------------------------------------------------------------------------------------------------------
# CONFIGURE OUR AZURE CONNECTION
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

resource "azurerm_resource_group" "rg" {
  name     = var.resource_group_name
  location = var.location
}

# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY AN AZURE APP SERVICE PLAN
# ---------------------------------------------------------------------------------------------------------------------

resource "azurerm_monitor_action_group" "actionGroup" {
  name                = var.app_name
  resource_group_name = azurerm_resource_group.rg.name
  short_name          = var.short_name
  tags                = azurerm_resource_group.rg.tags

  dynamic "email_receiver" {
    for_each = var.enable_email ? ["email_receiver"] : []
    content {
      name                    = var.email_name
      email_address           = var.email_address
      use_common_alert_schema = true
    }
  }

  dynamic "sms_receiver" {
    for_each = var.enable_sms ? ["sms_receiver"] : []
    content {
      name         = var.sms_name
      country_code = var.sms_country_code
      phone_number = var.sms_phone_number
    }
  }

  dynamic "webhook_receiver" {
    for_each = var.enable_webhook ? ["webhook_receiver"] : []
    content {
      name                    = var.webhook_name
      service_uri             = var.webhook_service_uri
      use_common_alert_schema = true
    }
  }

}