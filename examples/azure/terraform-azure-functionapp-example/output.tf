output "resource_group_name" {
  value = azurerm_resource_group.app_rg.name
}

output "function_app_id" {
  value = azurerm_function_app.function_app.id
}

output "default_hostname" {
  value = azurerm_function_app.function_app.default_hostname
}

output "function_app_kind" {
  value = azurerm_function_app.function_app.kind
}

output "function_app_name" {
  value = azurerm_function_app.function_app.name
}
