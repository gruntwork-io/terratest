output "resource_group_name" {
  value = "${azurerm_resource_group.example.name}"
}

output "virtual_network_name" {
  value = "${azurerm_virtual_network.example.name}"
}

output "first_subnet_address" {
  value = "${azurerm_subnet.subnet1.address_prefix}"
}

output "second_subnet_address" {
  value = "${azurerm_subnet.subnet2.address_prefix}"
}

output "public_ip_address" {
  value = "${azurerm_public_ip.example.ip_address}"
}

output "public_ip_fqdn" {
  value = "${azurerm_public_ip.example.fqdn}"
}
