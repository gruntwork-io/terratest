output "resource_group_name" {
  value = azurerm_resource_group.aca.name
}

output "container_app_env_name" {
  value = azurerm_container_app_environment.aca.name
}

output "container_app_name" {
  value = azurerm_container_app.aca.name
}

output "container_app_job_name" {
  value = azurerm_container_app_job.aca.name
}
