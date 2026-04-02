# ---------------------------------------------------------------------------------------------------------------------
# ENVIRONMENT VARIABLES
# You must define the following environment variables.
# ---------------------------------------------------------------------------------------------------------------------

# GOOGLE_CREDENTIALS
# or
# GOOGLE_APPLICATION_CREDENTIALS

variable "gcp_project_id" {
  description = "The ID of the GCP project in which these resources will be created."
}

# ---------------------------------------------------------------------------------------------------------------------
# REQUIRED PARAMETERS
# You must provide a value for each of these parameters.
# ---------------------------------------------------------------------------------------------------------------------
# (none)

# ---------------------------------------------------------------------------------------------------------------------
# OPTIONAL PARAMETERS
# These parameters have reasonable defaults.
# ---------------------------------------------------------------------------------------------------------------------

variable "topic_name" {
  description = "The name of the Pub/Sub topic to create."
  type        = string
  default     = "terratest-example-topic"
}

variable "subscription_name" {
  description = "The name of the Pub/Sub subscription to create."
  type        = string
  default     = "terratest-example-sub"
}
