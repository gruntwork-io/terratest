terraform {
  required_version = ">= 1.0"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

data "aws_availability_zones" "available" {
  state = "available"
}

provider "aws" {
  region = var.aws_region

  default_tags {
    tags = {
      "gw:repo"    = "https://github.com/gruntwork-io/terratest"
      "gw:example" = "terraform-aws-network-example"
    }
  }
}

# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY A SIMPLE NETWORK
# The network has an internet gateway and two subnets - private and public - in the same availability zone.
# ---------------------------------------------------------------------------------------------------------------------

resource "aws_vpc" "main" {
  cidr_block = var.main_vpc_cidr

  tags = {
    Name = var.tag_name
  }
}

resource "aws_internet_gateway" "main_gateway" {
  vpc_id = aws_vpc.main.id

  tags = {
    Name = var.tag_name
  }
}

resource "aws_subnet" "private" {
  vpc_id                  = aws_vpc.main.id
  cidr_block              = var.private_subnet_cidr
  map_public_ip_on_launch = false

  tags = {
    Name = var.tag_name
  }

  availability_zone = data.aws_availability_zones.available.names[0]
}

resource "aws_subnet" "public" {
  vpc_id                  = aws_vpc.main.id
  cidr_block              = var.public_subnet_cidr
  map_public_ip_on_launch = true

  tags = {
    Name = var.tag_name
  }

  availability_zone = data.aws_availability_zones.available.names[0]
}

# ---------------------------------------------------------------------------------------------------------------------
# CREATE AND ATTACH A ROUTING TABLE FOR THE PUBLIC NETWORK
# ---------------------------------------------------------------------------------------------------------------------

resource "aws_route_table" "public" {
  vpc_id = aws_vpc.main.id

  route {
    cidr_block = "91.189.0.0/24"
    gateway_id = aws_internet_gateway.main_gateway.id
  }

  tags = {
    Name = var.tag_name
  }
}

resource "aws_route_table_association" "public" {
  subnet_id      = aws_subnet.public.id
  route_table_id = aws_route_table.public.id
}

# ---------------------------------------------------------------------------------------------------------------------
# CREATE NAT GATEWAY FOR THE PRIVATE SUBNET
# ---------------------------------------------------------------------------------------------------------------------

resource "aws_eip" "nat" {
  domain = "vpc"
}

resource "aws_nat_gateway" "nat" {
  allocation_id = aws_eip.nat.id
  subnet_id     = aws_subnet.public.id
  depends_on    = [aws_internet_gateway.main_gateway]
}

# ---------------------------------------------------------------------------------------------------------------------
# CREATE AND ATTACH A ROUTING TABLE FOR THE PRIVATE NETWORK
# ---------------------------------------------------------------------------------------------------------------------

resource "aws_route_table" "private" {
  vpc_id = aws_vpc.main.id

  route {
    cidr_block     = "0.0.0.0/0"
    nat_gateway_id = aws_nat_gateway.nat.id
  }

  tags = {
    Name = var.tag_name
  }
}

resource "aws_route_table_association" "private" {
  subnet_id      = aws_subnet.private.id
  route_table_id = aws_route_table.private.id
}

