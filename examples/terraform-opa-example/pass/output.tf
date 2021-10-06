output "recommended_instance_type" {
  description = "The recommended instance type to use in this AWS region. This will be the first instance type in var.instance_types which is available in all AZs in this region."
  value       = module.instance_types.recommended_instance_type
}
