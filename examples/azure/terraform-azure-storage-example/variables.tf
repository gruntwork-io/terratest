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

variable "location" {
  description = "The location to set for the storage account."
  type        = string
  default     = "East US"
}

variable "storage_account_kind" {
  description = "The kind of storage account to set"
  type        = string
  default     = "StorageV2"
}

variable "storage_account_tier" {
  description = "The tier of storage account to set"
  type        = string
  default     = "Standard"
}

variable "storage_replication_type" {
  description = "The replication type of storage account to set"
  type        = string
  default     = "GRS"
}

variable "container_access_type" {
  description = "The replication type of storage account to set"
  type        = string
  default     = "private"
}

variable "postfix" {
  description = "A postfix string to centrally mitigate resource name collisions"
  type        = string
  default     = "resource"
}

