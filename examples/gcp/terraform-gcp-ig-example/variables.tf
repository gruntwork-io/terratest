# ---------------------------------------------------------------------------------------------------------------------
# ENVIRONMENT VARIABLES
# You must define the following environment variables.
# ---------------------------------------------------------------------------------------------------------------------

# GOOGLE_CREDENTIALS
# or
# GOOGLE_APPLICATION_CREDENTIALS

# ---------------------------------------------------------------------------------------------------------------------
# REQUIRED PARAMETERS
# You must provide a value for each of these parameters.
# ---------------------------------------------------------------------------------------------------------------------

variable "gcp_project_id" {
  description = "The ID of the GCP project in which these resources will be created."
  type        = string
}

variable "gcp_region" {
  description = "The region in which all GCP resources will be created."
  type        = string
}

# ---------------------------------------------------------------------------------------------------------------------
# OPTIONAL PARAMETERS
# These parameters have reasonable defaults.
# ---------------------------------------------------------------------------------------------------------------------

variable "cluster_name" {
  description = "The unique identifier for the resources created by this Terraform configuration."
  type        = string
  default     = "terratest-example"
}

variable "cluster_size" {
  description = "The number of Compute Instances to run in the Managed Instance Group."
  type        = number
  default     = 3
}

variable "machine_type" {
  description = "The Machine Type to use for the Compute Instances."
  type        = string
  default     = "f1-micro"
}

