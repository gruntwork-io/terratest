output "resource_group_name" {
  value = azurerm_resource_group.datafactory_rg.name
}

output "datafactory_name" {
  value = azurerm_data_factory.data_factory.name
}
