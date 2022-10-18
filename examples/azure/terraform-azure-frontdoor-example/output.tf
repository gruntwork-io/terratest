output "resource_group_name" {
  value = azurerm_resource_group.rg.name
}

output "front_door_name" {
  description = "Specifies the name of the Front Door service."
  value       = azurerm_frontdoor.frontdoor.name
}

output "front_door_url" {
  description = "Specifies the host name of the frontend_endpoint. Must be a domain name."
  value       = azurerm_frontdoor.frontdoor.frontend_endpoint[0].host_name
}

output "front_door_endpoint_name" {
  description = "Specifies the friendly name of the frontend_endpoint"
  value       = azurerm_frontdoor.frontdoor.frontend_endpoint[0].name
}