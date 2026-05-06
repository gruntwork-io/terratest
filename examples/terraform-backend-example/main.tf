# ---------------------------------------------------------------------------------------------------------------------
# AN EXAMPLE OF HOW TO CONFIGURE A TERRAFORM BACKEND WITH TERRATEST
# Note that the example code here doesn't do anything other than set up a backend that Terratest will configure.
# ---------------------------------------------------------------------------------------------------------------------

terraform {
  required_version = ">= 1.0"

  # Leave the config for this backend unspecified so Terraform can fill it in. This is known as "partial configuration":
  # https://www.terraform.io/docs/backends/config.html#partial-configuration
  backend "s3" {}
}

variable "foo" {
  description = "Some data to store as an output of this module"
  type        = string
}

output "foo" {
  value = var.foo
}
