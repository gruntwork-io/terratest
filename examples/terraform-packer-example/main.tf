terraform {
  required_version = ">= 1.0"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
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
  ami           = var.ami_id
  instance_type = var.instance_type

  user_data = templatefile("${path.module}/user-data/user-data.sh", {
    instance_text = var.instance_text
    instance_port = var.instance_port
  })

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
