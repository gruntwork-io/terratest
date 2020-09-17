# ---------------------------------------------------------------------------------------------------------------------
# PIN TERRAFORM VERSION TO >= 0.12
# The examples have been upgraded to 0.12 syntax
# ---------------------------------------------------------------------------------------------------------------------

terraform {
  # This module is now only being tested with Terraform 0.13.x. However, to make upgrading easier, we are setting
  # 0.12.26 as the minimum version, as that version added support for required_providers with source URLs, making it
  # forwards compatible with 0.13.x code.
  required_version = ">= 0.12.26"
}

# ---------------------------------------------------------------------------------------------------------------------
# BASIC TERRAFORM EXAMPLE
# See test/terraform_aws_example.go for how to write automated tests for this code.
# ---------------------------------------------------------------------------------------------------------------------

data "template_file" "example" {
  template = var.example
}

data "template_file" "example2" {
  template = var.example2
}

resource "local_file" "example" {
  content  = "${data.template_file.example.rendered} + ${data.template_file.example2.rendered}"
  filename = "example.txt"
}

resource "local_file" "example2" {
  content  = data.template_file.example2.rendered
  filename = "example2.txt"
}

