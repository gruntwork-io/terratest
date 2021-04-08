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

variable "key_pair_name" {
  description = "The EC2 Key Pair to associate with the EC2 Instance for SSH access."
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

variable "instance_name" {
  description = "The Name tag to set for the EC2 Instance."
  type        = string
  default     = "terratest-example"
}

variable "ssh_port" {
  description = "The port the EC2 Instance should listen on for SSH requests."
  type        = number
  default     = 22
}

variable "ssh_user" {
  description = "SSH user name to use for remote exec connections,"
  type        = string
  default     = "ubuntu"
}

variable "instance_type" {
  description = "Instance type to use for EC2 Instance"
  type        = string
  default     = "t2.micro"
}

