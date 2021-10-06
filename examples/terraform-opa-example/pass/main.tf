terraform {
  # This module is now only being tested with Terraform 1.0.x. However, to make upgrading easier, we are setting
  # 0.12.26 as the minimum version, as that version added support for required_providers with source URLs, making it
  # forwards compatible with 1.0.x code.
  required_version = ">= 0.12.26"
}

provider "aws" {
  region = var.aws_region
}

module "instance_types" {
  source         = "git::git@github.com:gruntwork-io/terraform-aws-utilities.git//modules/instance-type?ref=v0.6.0"
  instance_types = ["t2.micro", "t3.micro"]
}
