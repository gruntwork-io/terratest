output "planids" {
  value = "${azurerm_app_service_plan.plan.*.id}"
}