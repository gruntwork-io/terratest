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

variable "sqlmi_license_type" {
  description = "The license type for the sql managed instance"
  type        = string
  default     = "BasePrice"
}

variable "sku_name" {
  description = "The sku name for the sql managed instance"
  type        = string
  default     = "GP_Gen5"
}

variable "storage_size" {
  description = "The storage for the sql managed instance"
  type        = string
  default     = 32
}

variable "cores" {
  description = "The vcores for the sql managed instance"
  type        = string
  default     = 4
}

variable "admin_login" {
  description = "The login for the sql managed instance"
  type        = string
  default     = "sqlmiadmin"
}


variable "sqlmi_db_name" {
  description = "The Database for the sql managed instance"
  type        = string
  default     = "testdb"
}

variable "postfix" {
  description = "A postfix string to centrally mitigate resource name collisions."
  type        = string
  default     = "resource"
}