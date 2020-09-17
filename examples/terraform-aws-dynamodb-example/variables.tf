# ---------------------------------------------------------------------------------------------------------------------
# OPTIONAL PARAMETERS
# These parameters have reasonable defaults.
# ---------------------------------------------------------------------------------------------------------------------

variable "table_name" {
  description = "The name to set for the dynamoDB table."
  type        = string
  default     = "terratest-example"
}

