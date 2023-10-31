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

variable "postfix" {
  description = "A postfix string to centrally mitigate resource name collisions."
  type        = string
  default     = "example"
}

variable "location" {
  description = "The Azure Region where the MySQL Flexible Server should exist."
  type        = string
  default     = "West Europe"
}

variable "mysql_flexible_server_administrator_login" {
  description = "The Administrator login for the MySQL Flexible Server."
  type        = string
  default     = "mysqladmin"
}

variable "mysql_flexible_server_backup_retention_days" {
  description = "The backup retention days for the MySQL Flexible Server."
  type        = number
  default     = 7

  validation {
    condition     = var.mysql_flexible_server_backup_retention_days >= 1 && var.mysql_flexible_server_backup_retention_days <= 35
    error_message = "MySQL Flexible Server retention days should be between 1 and 35."
  }
}

variable "mysql_flexible_server_sku_name" {
  description = "The SKU Name for the MySQL Flexible Server."
  type        = string
  default     = "B_Standard_B1ms"
}

variable "mysql_flexible_server_storage_size_gb" {
  description = "The max storage allowed for the MySQL Flexible Server."
  type        = number
  default     = 32

  validation {
    condition     = var.mysql_flexible_server_storage_size_gb >= 20 && var.mysql_flexible_server_storage_size_gb <= 16384
    error_message = "MySQL Flexible Server storage size (GB) should be a value between 20 and 16384."
  }
}

variable "mysql_flexible_server_zone" {
  description = "Specifies the Availability Zone in which this MySQL Flexible Server should be located."
  type        = number
  default     = 1

  validation {
    condition     = var.mysql_flexible_server_zone == 1 || var.mysql_flexible_server_zone == 2 || var.mysql_flexible_server_zone == 3
    error_message = "MySQL Flexible Server possible Availability Zone are 1, 2 or 3."
  }
}

variable "mysql_version" {
  description = "The version of the MySQL Flexible Server to use."
  type        = string
  default     = "5.7"

  validation {
    condition     = var.mysql_version == "5.7" || var.mysql_version == "8.0.21"
    error_message = "MySQL version for Flexbile Server instance should be 5.7 or 8.0.21."
  }
}

variable "mysql_flexible_server_db_charset" {
  description = "Specifies the Charset for the MySQL Database, which needs to be a valid MySQL Charset."
  type        = string
  default     = "utf8"
}

variable "mysql_flexible_server_db_collation" {
  description = "Specifies the Collation for the MySQL Database, which needs to be a valid MySQL Collation."
  type        = string
  default     = "utf8_unicode_ci"
}