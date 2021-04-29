
output "availability_set_name" {
  value = azurerm_availability_set.avs.name
}

output "managed_disk_name" {
  value = azurerm_managed_disk.disk.name
}

output "managed_disk_type" {
  value = azurerm_managed_disk.disk.storage_account_type
}

output "network_interface_name" {
  value = azurerm_network_interface.nic.name
}

output "os_disk_name" {
  value = azurerm_virtual_machine.vm_example.storage_os_disk[0].name
}

output "private_ip" {
  value = azurerm_network_interface.nic.ip_configuration[0].private_ip_address
}

output "public_ip_name" {
  value = azurerm_public_ip.pip.name
}

output "resource_group_name" {
  value = azurerm_resource_group.vm_rg.name
}

output "subnet_name" {
  value = azurerm_subnet.subnet.name
}

output "virtual_network_name" {
  value = azurerm_virtual_network.vnet.name
}

output "vm_admin_username" {
  value = nonsensitive(azurerm_virtual_machine.vm_example.os_profile[*].admin_username)
}

output "vm_image_sku" {
  value = azurerm_virtual_machine.vm_example.storage_image_reference[*].sku
}

output "vm_image_version" {
  value = azurerm_virtual_machine.vm_example.storage_image_reference[*].version
}

output "vm_name" {
  value = azurerm_virtual_machine.vm_example.name
}

output "vm_size" {
  value = azurerm_virtual_machine.vm_example.vm_size
}

output "vm_tags" {
  value = azurerm_virtual_machine.vm_example.tags
}
