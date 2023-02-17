output "sku_name" {
  value = azurerm_postgresql_server.postgresqlserver.sku_name
}

output "servername" {
  value = azurerm_postgresql_server.postgresqlserver.name

}

output "rgname" {
  value = azurerm_resource_group.rg.name
}