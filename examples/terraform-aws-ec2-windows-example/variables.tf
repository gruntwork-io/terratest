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

variable "region" {
  description = "The AWS region in which all resources will be created"
  type        = string
  default     = "us-west-2"
}

variable "ami" {
  description = "The ID of the AMI to run on the Windows instance."
  type        = string
}

# ---------------------------------------------------------------------------------------------------------------------
# OPTIONAL PARAMETERS
# These parameters have reasonable defaults.
# ---------------------------------------------------------------------------------------------------------------------

variable "name" {
  description = "The name of the Windows instance"
  type        = string
  default     = "windows_test_instance"
}

variable "instance_type" {
  description = "The instance type to deploy."
  type        = string
  default     = "t3.small"
}

variable "root_volume_size" {
  description = "The size in GiB of the root volume. Must match the root volume size of the target AMI."
  type        = number
  default     = 30
}

