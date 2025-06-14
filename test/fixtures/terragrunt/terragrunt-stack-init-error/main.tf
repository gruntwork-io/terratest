# Simple Terraform configuration
resource "null_resource" "test" {
  provisioner "local-exec" {
    command = "echo 'Test resource'"
  }
} 