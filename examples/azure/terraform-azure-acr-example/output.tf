output "resource_group_name" {
  value = azurerm_resource_group.rg.name
}

output "container_registry_name" {
  value = azurerm_container_registry.acr.name
}

output "login_server" {
  value = azurerm_container_registry.acr.login_server
}

output "admin_username" {
  value     = azurerm_container_registry.acr.admin_username
  sensitive = true
}

output "admin_password" {
  value     = azurerm_container_registry.acr.admin_password
  sensitive = true
}
