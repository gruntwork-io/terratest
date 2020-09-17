# ---------------------------------------------------------------------------------------------------------------------
# AN EXAMPLE OF HOW TO CONFIGURE A TERRAFORM BACKEND WITH TERRATEST
# Note that the example code here doesn't do anything other than set up a backend that Terratest will configure.
# ---------------------------------------------------------------------------------------------------------------------

terraform {
  # Leave the config for this backend unspecified so Terraform can fill it in. This is known as "partial configuration":
  # https://www.terraform.io/docs/backends/config.html#partial-configuration
  backend "s3" {}
  # This module is now only being tested with Terraform 0.13.x. However, to make upgrading easier, we are setting
  # 0.12.26 as the minimum version, as that version added support for required_providers with source URLs, making it
  # forwards compatible with 0.13.x code.
  required_version = ">= 0.12.26"
}

variable "foo" {
  description = "Some data to store as an output of this module"
  type        = string
}

output "foo" {
  value = var.foo
}
