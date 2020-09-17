# ---------------------------------------------------------------------------------------------------------------------
# ENVIRONMENT VARIABLES
# Define these secrets as environment variables
# ---------------------------------------------------------------------------------------------------------------------

# ARM_CLIENT_ID
# ARM_CLIENT_SECRET
# ARM_SUBSCRIPTION_ID
# ARM_TENANT_ID

# ---------------------------------------------------------------------------------------------------------------------
# REQUIRED PARAMETERS
# You must provide a value for each of these parameters.
# ---------------------------------------------------------------------------------------------------------------------

# ---------------------------------------------------------------------------------------------------------------------
# OPTIONAL PARAMETERS
# These parameters have reasonable defaults.
# ---------------------------------------------------------------------------------------------------------------------

variable "hostname" {
  description = "The hostname of the new VM to be configured"
  type        = string
  default     = "terratest-vm"
}

variable "password" {
  description = "The password to configure for SSH access"
  type        = string
  default     = "HorriblePassword1234!"
}

variable "prefix" {
  description = "The prefix that will be attached to all resources deployed"
  type        = string
  default     = "terratest-example"
}

variable "username" {
  description = "The username to be provisioned into your VM"
  type        = string
  default     = "testadmin"
}

