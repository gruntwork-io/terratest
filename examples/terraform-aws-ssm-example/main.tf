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

provider "aws" {
  region = var.region
}

# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY AN INSTANCE WITH SSM SUPPORT
# ---------------------------------------------------------------------------------------------------------------------

data "aws_iam_policy_document" "example" {
  version = "2012-10-17"

  statement {
    sid = "1"

    actions = [
      "sts:AssumeRole",
    ]

    principals {
      type        = "Service"
      identifiers = ["ec2.amazonaws.com"]
    }
  }
}

resource "aws_iam_role" "example" {
  name_prefix        = "example"
  assume_role_policy = data.aws_iam_policy_document.example.json
}

resource "aws_iam_role_policy_attachment" "example_ssm" {
  role       = aws_iam_role.example.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonEC2RoleforSSM"
}

resource "aws_iam_instance_profile" "example" {
  name_prefix = "example"
  role        = aws_iam_role.example.name
}

data "aws_ami" "amazon_linux_2" {
  most_recent = true
  owners      = ["amazon"]

  filter {
    name   = "name"
    values = ["amzn2-ami-hvm*"]
  }
}

# ---------------------------------------------------------------------------------------------------------------------
# The instance must have a public ip to be able to contact AWS SSM
# ---------------------------------------------------------------------------------------------------------------------

resource "aws_instance" "example" {
  ami                         = data.aws_ami.amazon_linux_2.id
  instance_type               = var.instance_type
  associate_public_ip_address = true
  iam_instance_profile        = aws_iam_instance_profile.example.name
}
