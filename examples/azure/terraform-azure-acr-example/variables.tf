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

variable "sku" {
  description = "SKU tier for the ACR."
  default     = "Premium"
}


variable "location" {
  description = "The supported azure location where the resource exists"
  type        = string
  default     = "West US2"
}

variable "postfix" {
  description = "A postfix string to centrally mitigate resource name collisions."
  type        = string
  default     = "1276"
}
