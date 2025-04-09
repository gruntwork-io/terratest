# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY AN AZURE FRONT DOOR
# This is an example of how to deploy an Azure Front Door with the minimum resources.
# ---------------------------------------------------------------------------------------------------------------------
# See test/azure/terraform_azure_frontdoor_example_test.go for how to write automated tests for this code.
# ---------------------------------------------------------------------------------------------------------------------

terraform {
  required_version = ">=0.14.0"
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
# DEPLOY FRONT DOOR
# ---------------------------------------------------------------------------------------------------------------------

resource "azurerm_frontdoor" "frontdoor" {
  name                = "terratest-afd-${var.postfix}"
  resource_group_name = azurerm_resource_group.rg.name

  backend_pool_settings {
    enforce_backend_pools_certificate_name_check = false
  }

  routing_rule {
    name               = "terratestRoutingRule1"
    accepted_protocols = ["Http", "Https"]
    patterns_to_match  = ["/*"]
    frontend_endpoints = ["terratestEndpoint"]
    forwarding_configuration {
      forwarding_protocol = "MatchRequest"
      backend_pool_name   = "terratestBackend"
    }
  }

  backend_pool_load_balancing {
    name = "terratestLoadBalanceSetting"
  }

  backend_pool_health_probe {
    name = "terratestHealthProbeSetting"
  }

  backend_pool {
    name = "terratestBackend"
    backend {
      host_header = var.backend_host
      address     = var.backend_host
      http_port   = 80
      https_port  = 443
    }

    load_balancing_name = "terratestLoadBalanceSetting"
    health_probe_name   = "terratestHealthProbeSetting"
  }

  frontend_endpoint {
    name      = "terratestEndpoint"
    host_name = "terratest-afd-${var.postfix}.azurefd.net"
  }
}
