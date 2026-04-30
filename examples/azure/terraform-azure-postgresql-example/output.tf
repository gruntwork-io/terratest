output "sku_name" {
  value = azurerm_postgresql_flexible_server.postgresqlserver.sku_name
}

output "servername" {
  value = azurerm_postgresql_flexible_server.postgresqlserver.name
}

output "rgname" {
  value = azurerm_resource_group.rg.name
}

output "fqdn" {
  value = azurerm_postgresql_flexible_server.postgresqlserver.fqdn
}

output "database_name" {
  value = azurerm_postgresql_flexible_server_database.postgresqldb.name
}
