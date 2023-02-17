output "resource_group_name" {
  value = azurerm_resource_group.resource_group.name
}

output "loganalytics_workspace_name" {
  value = azurerm_log_analytics_workspace.log_analytics_workspace.name
}

output "loganalytics_workspace_sku" {
  value = azurerm_log_analytics_workspace.log_analytics_workspace.sku
}

output "loganalytics_workspace_retention" {
  value = azurerm_log_analytics_workspace.log_analytics_workspace.retention_in_days
}

