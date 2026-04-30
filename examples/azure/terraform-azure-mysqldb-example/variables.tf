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
  description = "The supported azure location where the resource exists"
  type        = string
  default     = "West US2"
}

variable "mysqlserver_admin_login" {
  description = "The administrator login name for the mysql flexible server."
  type        = string
  default     = "mysqladmin"
}

variable "mysqlserver_sku_name" {
  description = "The SKU name for the mysql flexible server (Burstable / GeneralPurpose / MemoryOptimized)."
  type        = string
  default     = "B_Standard_B1ms"
}

variable "mysqlserver_storage_size_gb" {
  description = "The max storage allowed (in GB) for the mysql flexible server."
  type        = number
  default     = 20
}

variable "mysqlserver_version" {
  description = "The MySQL engine version for the flexible server."
  type        = string
  default     = "8.0.21"
}

variable "mysqldb_charset" {
  description = "The charset for the mysql database."
  type        = string
  default     = "utf8"
}

variable "mysqldb_collation" {
  description = "The collation for the mysql database."
  type        = string
  default     = "utf8_unicode_ci"
}

variable "postfix" {
  description = "A postfix string to centrally mitigate resource name collisions."
  type        = string
  default     = "resource"
}
