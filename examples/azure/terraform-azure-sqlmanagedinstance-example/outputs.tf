output "resource_group_name" {
  value = azurerm_resource_group.sqlmi_rg.name
}

output "network_security_group_name" {
  value = azurerm_network_security_group.sqlmi_nt_sec_grp.name
}

output "virtual_network_name" {
  value = azurerm_virtual_network.sqlmi_vm.name
}

output "subnet_name" {
  value = azurerm_subnet.sqlmi_sub.name
}

output "managed_instance_name" {
  value = azurerm_mssql_managed_instance.sqlmi_mi.name
}

output "managed_instance_db_name" {
  value = azurerm_mssql_managed_database.sqlmi_db.name
}