# This is a test resource that echoes the message specified by var.message
resource "null_resource" "greet" {
  count = 5

  triggers = {
    run_always = timestamp()
  }

  provisioner "local-exec" {
    command = "echo ${var.message}"
  }
}
