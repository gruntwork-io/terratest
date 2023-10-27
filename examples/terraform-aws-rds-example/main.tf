# ---------------------------------------------------------------------------------------------------------------------
# PIN TERRAFORM VERSION TO >= 0.12
# The examples have been upgraded to 0.12 syntax
# ---------------------------------------------------------------------------------------------------------------------

provider "aws" {
  region = var.region
}

terraform {
  # This module is now only being tested with Terraform 0.13.x. However, to make upgrading easier, we are setting
  # 0.12.26 as the minimum version, as that version added support for required_providers with source URLs, making it
  # forwards compatible with 0.13.x code.
  required_version = ">= 0.12.26"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = ">= 4.61.0, < 5.0.0"
    }
  }
}

# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY INTO THE DEFAULT VPC AND SUBNETS
# To keep this example simple, we are deploying into the Default VPC and its subnets. In real-world usage, you should
# deploy into a custom VPC and private subnets. Given the subnet group needs to span multiple AZs and hence subnets we
# have deployed it across all the subnets of the default VPC.
# ---------------------------------------------------------------------------------------------------------------------

data "aws_vpc" "default" {
  default = true
}

data "aws_subnets" "all" {
  filter {
    name   = "vpc-id"
    values = [data.aws_vpc.default.id]
  }
}

# ---------------------------------------------------------------------------------------------------------------------
# CREATE AN SUBNET GROUP ACROSS ALL THE SUBNETS OF THE DEFAULT ASG TO HOST THE RDS INSTANCE
# ---------------------------------------------------------------------------------------------------------------------

resource "aws_db_subnet_group" "example" {
  name       = var.name
  subnet_ids = data.aws_subnets.all.ids

  tags = {
    Name = var.name
  }
}

# ---------------------------------------------------------------------------------------------------------------------
# CREATE A CUSTOM PARAMETER GROUP AND AN OPTION GROUP FOR CONFIGURABILITY
# ---------------------------------------------------------------------------------------------------------------------

resource "aws_db_option_group" "example" {
  name                 = var.name
  engine_name          = var.engine_name
  major_engine_version = var.major_engine_version

  tags = {
    Name = var.name
  }

  option {
    option_name = "MARIADB_AUDIT_PLUGIN"

    option_settings {
      name  = "SERVER_AUDIT_EVENTS"
      value = "CONNECT"
    }
  }
}

resource "aws_db_parameter_group" "example" {
  name   = var.name
  family = var.family

  tags = {
    Name = var.name
  }

  parameter {
    name  = "general_log"
    value = "0"
  }
}

# ---------------------------------------------------------------------------------------------------------------------
# CREATE A SECURITY GROUP TO ALLOW ACCESS TO THE RDS INSTANCE
# ---------------------------------------------------------------------------------------------------------------------

resource "aws_security_group" "db_instance" {
  name   = var.name
  vpc_id = data.aws_vpc.default.id
}

resource "aws_security_group_rule" "allow_db_access" {
  type              = "ingress"
  from_port         = var.port
  to_port           = var.port
  protocol          = "tcp"
  security_group_id = aws_security_group.db_instance.id
  cidr_blocks       = ["0.0.0.0/0"]
}

# ---------------------------------------------------------------------------------------------------------------------
# CREATE THE DATABASE INSTANCE
# ---------------------------------------------------------------------------------------------------------------------

resource "aws_db_instance" "example" {
  identifier             = var.name
  engine                 = var.engine_name
  engine_version         = var.engine_version
  port                   = var.port
  db_name                = var.database_name
  username               = var.username
  password               = var.password
  instance_class         = var.instance_class
  allocated_storage      = var.allocated_storage
  skip_final_snapshot    = true
  license_model          = var.license_model
  db_subnet_group_name   = aws_db_subnet_group.example.id
  vpc_security_group_ids = [aws_security_group.db_instance.id]
  publicly_accessible    = true
  parameter_group_name   = aws_db_parameter_group.example.id
  option_group_name      = aws_db_option_group.example.id

  tags = {
    Name = var.name
  }
}
