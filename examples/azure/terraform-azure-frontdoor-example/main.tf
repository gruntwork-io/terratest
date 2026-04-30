# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY AN AZURE FRONT DOOR (CDN FRONT DOOR)
# This is an example of how to deploy an Azure Front Door using the modern azurerm_cdn_frontdoor_* resource family.
# The legacy `azurerm_frontdoor` resource was deprecated on April 1, 2025 and no longer permits creation of new
# Front Door resources, so this example uses the CDN Front Door resources that supersede it.
# ---------------------------------------------------------------------------------------------------------------------
# See test/azure/terraform_azure_frontdoor_example_test.go for how to write automated tests for this code.
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
  name     = "terratest-frontdoor-rg-${var.postfix}"
  location = var.location
}

# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY THE CDN FRONT DOOR PROFILE
# This is the top-level Front Door resource (replaces the classic `azurerm_frontdoor`).
# ---------------------------------------------------------------------------------------------------------------------

resource "azurerm_cdn_frontdoor_profile" "frontdoor" {
  name                = "terratest-afd-${var.postfix}"
  resource_group_name = azurerm_resource_group.rg.name
  sku_name            = "Standard_AzureFrontDoor"
}

# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY THE FRONT DOOR ENDPOINT
# A profile can host one or more endpoints (front-end host names).
# ---------------------------------------------------------------------------------------------------------------------

resource "azurerm_cdn_frontdoor_endpoint" "endpoint" {
  name                     = "terratest-ep-${var.postfix}"
  cdn_frontdoor_profile_id = azurerm_cdn_frontdoor_profile.frontdoor.id
}

# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY THE ORIGIN GROUP AND ORIGIN
# An origin group describes how the Front Door load-balances and health-checks a set of origins.
# ---------------------------------------------------------------------------------------------------------------------

resource "azurerm_cdn_frontdoor_origin_group" "origin_group" {
  name                     = "terratestOriginGroup"
  cdn_frontdoor_profile_id = azurerm_cdn_frontdoor_profile.frontdoor.id

  load_balancing {
    sample_size                 = 4
    successful_samples_required = 3
  }

  health_probe {
    interval_in_seconds = 100
    path                = "/"
    protocol            = "Https"
    request_type        = "HEAD"
  }
}

resource "azurerm_cdn_frontdoor_origin" "origin" {
  name                          = "terratestOrigin"
  cdn_frontdoor_origin_group_id = azurerm_cdn_frontdoor_origin_group.origin_group.id

  enabled                        = true
  host_name                      = var.backend_host
  http_port                      = 80
  https_port                     = 443
  origin_host_header             = var.backend_host
  priority                       = 1
  weight                         = 1000
  certificate_name_check_enabled = true
}

# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY THE ROUTE
# Routes wire an endpoint to an origin group with a forwarding/caching configuration.
# ---------------------------------------------------------------------------------------------------------------------

resource "azurerm_cdn_frontdoor_route" "route" {
  name                          = "terratestRoute"
  cdn_frontdoor_endpoint_id     = azurerm_cdn_frontdoor_endpoint.endpoint.id
  cdn_frontdoor_origin_group_id = azurerm_cdn_frontdoor_origin_group.origin_group.id
  cdn_frontdoor_origin_ids      = [azurerm_cdn_frontdoor_origin.origin.id]

  enabled                = true
  forwarding_protocol    = "MatchRequest"
  https_redirect_enabled = false
  patterns_to_match      = ["/*"]
  supported_protocols    = ["Http", "Https"]

  link_to_default_domain = true
}
