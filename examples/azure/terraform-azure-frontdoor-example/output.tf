output "resource_group_name" {
  value = azurerm_resource_group.rg.name
}

output "front_door_name" {
  description = "Specifies the name of the CDN Front Door profile (replaces the legacy Front Door service name)."
  value       = azurerm_cdn_frontdoor_profile.frontdoor.name
}

output "front_door_url" {
  description = "Specifies the host name (FQDN) of the Front Door endpoint."
  value       = azurerm_cdn_frontdoor_endpoint.endpoint.host_name
}

output "front_door_endpoint_name" {
  description = "Specifies the name of the Front Door endpoint."
  value       = azurerm_cdn_frontdoor_endpoint.endpoint.name
}
