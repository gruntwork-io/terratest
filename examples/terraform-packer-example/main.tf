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
# DEPLOY AN EC2 INSTANCE THAT RUNS A SIMPLE RUBY WEB APP BUILT USING A PACKER TEMPLATE
# See test/terraform_packer_example.go for how to write automated tests for this code.
# ---------------------------------------------------------------------------------------------------------------------

provider "aws" {
  region = var.aws_region
}

# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY THE EC2 INSTANCE
# ---------------------------------------------------------------------------------------------------------------------

resource "aws_instance" "example" {
  ami                    = var.ami_id
  instance_type          = var.instance_type
  user_data              = data.template_file.user_data.rendered
  vpc_security_group_ids = [aws_security_group.example.id]

  tags = {
    Name = var.instance_name
  }
}

# ---------------------------------------------------------------------------------------------------------------------
# CREATE A SECURITY GROUP TO CONTROL WHAT REQUESTS CAN GO IN AND OUT OF THE EC2 INSTANCE
# ---------------------------------------------------------------------------------------------------------------------

resource "aws_security_group" "example" {
  name = var.instance_name

  ingress {
    from_port = var.instance_port
    to_port   = var.instance_port
    protocol  = "tcp"

    # To keep this example simple, we allow incoming HTTP requests from any IP. In real-world usage, you may want to
    # lock this down to just the IPs of trusted servers (e.g., of a load balancer).
    cidr_blocks = ["0.0.0.0/0"]
  }
}

# ---------------------------------------------------------------------------------------------------------------------
# CREATE THE USER DATA SCRIPT THAT WILL RUN DURING BOOT ON THE EC2 INSTANCE
# ---------------------------------------------------------------------------------------------------------------------

data "template_file" "user_data" {
  template = file("${path.module}/user-data/user-data.sh")

  vars = {
    instance_text = var.instance_text
    instance_port = var.instance_port
  }
}

