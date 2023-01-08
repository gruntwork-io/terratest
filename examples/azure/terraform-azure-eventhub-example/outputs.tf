output "resource_group_name" {
  value = azurerm_resource_group.main.name
}

output "eventhub_namespace_name" {
  value = azurerm_eventhub_namespace.main.name
}

output "eventhub_name" {
  value = azurerm_eventhub.main.name
}

output "eventhub_partition_count" {
  value = azurerm_eventhub.main.partition_count
}

output "eventhub_message_retention" {
  value = azurerm_eventhub.main.message_retention
}
