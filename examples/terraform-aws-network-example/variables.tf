# ---------------------------------------------------------------------------------------------------------------------
# ENVIRONMENT VARIABLES
# Define these secrets as environment variables
# ---------------------------------------------------------------------------------------------------------------------

# AWS_ACCESS_KEY_ID
# AWS_SECRET_ACCESS_KEY

# ---------------------------------------------------------------------------------------------------------------------
# REQUIRED PARAMETERS
# You must provide a value for each of these parameters.
# ---------------------------------------------------------------------------------------------------------------------

variable "main_vpc_cidr" {
  description = "The CIDR of the main VPC"
  type        = string
}

variable "public_subnet_cidr" {
  description = "The CIDR of public subnet"
  type        = string
}

variable "private_subnet_cidr" {
  description = "The CIDR of the private subnet"
  type        = string
}

# ---------------------------------------------------------------------------------------------------------------------
# OPTIONAL PARAMETERS
# These parameters have reasonable defaults.
# ---------------------------------------------------------------------------------------------------------------------

variable "aws_region" {
  description = "The AWS region to deploy into"
  type        = string
  default     = "us-east-1"
}

variable "tag_name" {
  description = "A name used to tag the resource"
  type        = string
  default     = "terraform-network-example"
}

