output "resource_group_name" {
  value = azurerm_resource_group.automation_account_dsc_rg.name
}
output "automation_account_name" {
  value = azurerm_automation_account.automation_account.name
}
output "sample_dsc_name" {
  value = azurerm_automation_dsc_configuration.sample_dsc.name
}
