output "resource_group_name" {
  value = azurerm_resource_group.vmss_rg.name
}

output "public_ip_name" {
  value = azurerm_public_ip.pip.name
}

output "subnet_name" {
  value = azurerm_subnet.subnet.name
}

output "virtual_network_name" {
  value = azurerm_virtual_network.vnet.name
}

output "lb_public_name" {
  value = azurerm_lb.public.name
}

output "scaleset_name" {
  value = azurerm_virtual_machine_scale_set.scaleset.name
}

output "scaleset_vm_name_prefix" {
  sensitive = true
  value     = lookup(azurerm_virtual_machine_scale_set.scaleset.os_profile[0], "computer_name_prefix")
}

output "scaleset_vm_size" {
  value = lookup(azurerm_virtual_machine_scale_set.scaleset.sku[0], "name")
}

output "scaleset_capacity" {
  value = lookup(azurerm_virtual_machine_scale_set.scaleset.sku[0], "capacity")
}

output "scaleset_tags" {
  value = azurerm_virtual_machine_scale_set.scaleset.tags
}
