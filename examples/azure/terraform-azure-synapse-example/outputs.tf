output "resource_group_name" {
  value = azurerm_resource_group.synapse_rg.name
}

output "synapse_storage_name" {
  value = azurerm_storage_account.storage_account.name
}

output "synapse_dlgen2_name" {
  value = azurerm_storage_data_lake_gen2_filesystem.dl_gen2.name
}

output "synapse_workspace_name" {
  value = azurerm_synapse_workspace.synapse_workspace.name
}

output "synapse_sqlpool_name" {
  value = azurerm_synapse_sql_pool.synapse_pool.name
}
