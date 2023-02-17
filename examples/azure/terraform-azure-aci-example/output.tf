output "resource_group_name" {
  value = azurerm_resource_group.rg.name
}

output "ip_address" {
  value = azurerm_container_group.aci.ip_address
}

output "fqdn" {
  value = azurerm_container_group.aci.fqdn
}

output "container_instance_name" {
  value = azurerm_container_group.aci.name
}
