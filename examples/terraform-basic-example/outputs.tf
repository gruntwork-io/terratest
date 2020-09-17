output "example" {
  value = data.template_file.example.rendered
}

output "example2" {
  value = data.template_file.example2.rendered
}

output "example_list" {
  value = var.example_list
}

output "example_map" {
  value = var.example_map
}

output "example_any" {
  value = var.example_any
}
